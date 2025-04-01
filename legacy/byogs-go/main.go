// Copyright 2020 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main is a very simple server with UDP (default), TCP, or both
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"agones.dev/agones/pkg/util/signals"
	sdk "agones.dev/agones/sdks/go"
	inventorypb "github.com/snapser/simplegs/snapserpb/inventory"
	statspb "github.com/snapser/simplegs/snapserpb/statistics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type GameServer struct {
	Ctx              context.Context
	S                *sdk.SDK
	InventoryClient  inventorypb.InventoryServiceClient
	StatisticsClient statspb.StatisticsServiceClient
	Matches          map[string]*Match
	AddrMatchID      map[string]string
	AddrPlayerID     map[string]string
}

type Match struct {
	MatchId string
	Players map[string]*Player
}

type Player struct {
	UserId     string
	Connection net.Addr
}

func New() (*GameServer, error) {
	sigCtx, _ := signals.NewSigKillContext()
	log.Print("Creating SDK instance")
	s, err := sdk.NewSDK()
	if err != nil {
		log.Fatalf("Could not connect to sdk: %v", err)
	}

	gs := &GameServer{
		Ctx: sigCtx,
		S:   s,
	}
	inventoryURL := os.Getenv("SNAPEND_INVENTORY_GRPC_URL")
	if inventoryURL == "" {
		log.Print("SNAPEND_INVENTORY_GRPC_URL not set")
	} else {
		inventoryClient, err := createInventoryClient(inventoryURL)
		if err != nil {
			log.Print("Error creating inventory client")
		} else {
			gs.InventoryClient = inventoryClient
		}
	}
	statisticsUrl := os.Getenv("SNAPEND_STATISTICS_GRPC_URL")
	if statisticsUrl == "" {
		log.Printf("SNAPEND_STATISTICS_GRPC_URL not set")
	} else {
		statisticsClient, err := createStatisticsClient(statisticsUrl)
		if err != nil {
			log.Printf("Error creating statistics client: %s", err.Error())
		} else {
			gs.StatisticsClient = statisticsClient
		}
	}
	return gs, nil
}

// main starts a UDP or TCP server
func main() {
	port := flag.String("port", "7654", "The port to listen to traffic on")
	udp := flag.Bool("udp", true, "Server will listen on UDP")

	flag.Parse()
	if ep := os.Getenv("PORT"); ep != "" {
		port = &ep
	}
	if eudp := os.Getenv("UDP"); eudp != "" {
		u := strings.ToUpper(eudp) == "TRUE"
		udp = &u
	}
	gs, err := New()
	if err != nil {
		log.Fatalf("Could not create GameServer: %v", err)
		return
	}

	if *udp {
		go gs.UdpListener(port)
		// _ = gs.S.WatchGameServer(func(baseGS *basesdk.GameServer) {
		// 	//gs.Logger.Info().Str("status", baseGS.Status.String()).Interface("Annotations", baseGS.ObjectMeta.Annotations).Interface("labels", baseGS.ObjectMeta.Labels).Msg("Agones SDK Event: " + baseGS.Status.State)
		// 	if baseGS.Status.State == "Allocated" {
		// 		gs.updateGameServerState(baseGS.ObjectMeta.Name, "READY")
		// 	}
		// })
	}
	log.Println("Readying")
	gs.Ready()
	log.Println("Ready")

	<-gs.Ctx.Done()
	os.Exit(0)
}

func (gs *GameServer) handleResponse(txt string) (response string, addACK bool, responseError error) {
	parts := strings.Split(strings.TrimSpace(txt), " ")
	response = txt
	addACK = true
	responseError = nil

	switch parts[0] {
	// shuts down the gameserver
	case "EXIT":
		// handle elsewhere, as we respond before exiting
		return
	case "CRASH":
		log.Print("Crashing.")
		os.Exit(1)
		return "", false, nil
	case "WIN":
		if len(parts) < 2 {
			return "", true, fmt.Errorf("no user id provided")
		}
		ctx := metadata.AppendToOutgoingContext(gs.Ctx, "gateway", "internal")
		_, err := gs.StatisticsClient.IncrementUserStatistic(ctx, &statspb.IncrementUserStatisticRequest{
			UserId: parts[1],
			Key:    "wins",
			Delta:  1,
		})
		if err != nil {
			log.Printf("Error win statistic: %s", err.Error())
			return "", true, fmt.Errorf("could not increment wins: %w", err)
		}

		_, err = gs.InventoryClient.UpdateUserVirtualCurrency(ctx, &inventorypb.UpdateUserVirtualCurrencyRequest{
			UserId:       parts[1],
			CurrencyName: "coins",
			Amount:       100,
		})
		if err != nil {
			log.Printf("Error currency: %s", err.Error())
			return "", true, fmt.Errorf("could not update user virtual currency: %w", err)
		}
		return fmt.Sprintf("%s winner\n", parts[1]), false, nil
	case "LOSE":
		log.Print("Losing.")
		if len(parts) < 2 {
			return "", true, fmt.Errorf("no user id provided")
		}
		ctx := metadata.AppendToOutgoingContext(gs.Ctx, "gateway", "internal")
		_, err := gs.StatisticsClient.IncrementUserStatistic(ctx, &statspb.IncrementUserStatisticRequest{
			UserId: parts[1],
			Key:    "losses",
			Delta:  1,
		})
		if err != nil {
			log.Printf("Error lose statistic: %s", err.Error())
			return "", false, fmt.Errorf("could not increment losses: %w", err)
		}
		return fmt.Sprintf("%s loser\n", parts[1]), false, nil
	}
	return
}

func (gs *GameServer) UdpListener(port *string) {
	log.Printf("Starting UDP server, listening on port %s", *port)
	conn, err := net.ListenPacket("udp", ":"+*port)
	if err != nil {
		log.Fatalf("Could not start UDP server: %v", err)
	}
	defer conn.Close() // nolint: errcheck
	gs.udpReadWriteLoop(conn)
}

func (gs *GameServer) udpReadWriteLoop(conn net.PacketConn) {
	b := make([]byte, 1024)
	for {
		sender, txt := readPacket(conn, b)

		log.Printf("Received UDP: %v", txt)

		response, addACK, err := gs.handleResponse(txt)
		if err != nil {
			response = "ERROR: " + response + "\n"
		} else if addACK {
			response = "ACK: " + response + "\n"
		}

		gs.udpRespond(conn, sender, response)

		if txt == "EXIT" {
			gs.exit()
		}
	}
}

// respond responds to a given sender.
func (gs *GameServer) udpRespond(conn net.PacketConn, sender net.Addr, txt string) {
	if _, err := conn.WriteTo([]byte(txt), sender); err != nil {
		log.Fatalf("Could not write to udp stream: %v", err)
	}
}

// readPacket reads a string from the connection
func readPacket(conn net.PacketConn, b []byte) (net.Addr, string) {
	n, sender, err := conn.ReadFrom(b)
	if err != nil {
		log.Fatalf("Could not read from udp stream: %v", err)
	}
	txt := strings.TrimSpace(string(b[:n]))
	log.Printf("Received packet from %v: %v", sender.String(), txt)
	return sender, txt
}

// exit shutdowns the server
func (gs *GameServer) exit() {
	log.Printf("Received EXIT command. Exiting.")
	// This tells Agones to shutdown this Game Server
	shutdownErr := gs.S.Shutdown()
	if shutdownErr != nil {
		log.Printf("Could not shutdown")
	}
	// The process will exit when Agones removes the pod and the
	// container receives the SIGTERM signal
}

func (gs *GameServer) Ready() {
	err := gs.S.Ready()
	if err != nil {
		log.Fatalf("Could not send ready message")
	}
}

func (gs *GameServer) updateGameServerState(name string, state string) {
	log.Printf("Updating GameServer state: %v", state)

	url := os.Getenv("GATEWAY_URL")
	apiKey := os.Getenv("API_KEY")

	if url == "" || apiKey == "" {
		log.Printf("Gateway URL or API key not set")
		return
	}

	url = fmt.Sprintf(url, name)

	log.Printf("Url: %v Api-Key: %v", url, apiKey)

	jsonPayload := []byte(`{
		"state": "READY",
		"game_server_state_metadata": {
			"map": "desolation_sound",
			"difficulty": "nightmare"
		}
	}`)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("App-Key", apiKey)

	client := &http.Client{}

	res, err := client.Do(req)
	defer func() {
		if res != nil {
			err := res.Body.Close()
			if err != nil {
				log.Printf("Error closing response body")
			}
		}
	}()

	if err != nil {
		log.Printf("Error sending request")
		return
	}

	if res.StatusCode != 200 {
		log.Printf("Error updating gameserver state: %v", res.StatusCode)
		return
	}
	log.Printf("GameServer state updated: %v", res.Body)
}

func createInventoryClient(url string) (inventorypb.InventoryServiceClient, error) {
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := inventorypb.NewInventoryServiceClient(conn)
	return client, nil
}

func createStatisticsClient(url string) (statspb.StatisticsServiceClient, error) {
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := statspb.NewStatisticsServiceClient(conn)
	return client, nil
}
