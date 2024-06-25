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
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	basesdk "agones.dev/agones/pkg/sdk"
	"agones.dev/agones/pkg/util/signals"
	sdk "agones.dev/agones/sdks/go"
	"github.com/rs/zerolog"
)

type GameServer struct {
	Ctx    context.Context
	S      *sdk.SDK
	Logger zerolog.Logger
}

func New() (*GameServer, error) {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	sigCtx, _ := signals.NewSigKillContext()
	log.Print("Creating SDK instance")
	s, err := sdk.NewSDK()
	if err != nil {
		logger.Fatal().Err(err).Msg("Could not connect to sdk")
	}

	gs := &GameServer{
		Ctx:    sigCtx,
		S:      s,
		Logger: logger,
	}
	return gs, nil
}

func main() {
	port := flag.String("port", "7654", "The port to listen to traffic on")

	flag.Parse()
	if ep := os.Getenv("PORT"); ep != "" {
		port = &ep
	}
	gs, err := New()
	if err != nil {
		gs.Logger.Fatal().Err(err).Msg("Could not create GameServer")
		return
	}

	go gs.UdpListener(port)
	gs.S.WatchGameServer(func(baseGS *basesdk.GameServer) {
		gs.Logger.Info().Str("status", baseGS.Status.String()).Interface("Annotations", baseGS.ObjectMeta.Annotations).Interface("labels", baseGS.ObjectMeta.Labels).Msg("Watching gameserver")
	})
	gs.Logger.Info().Msg("Readying")
	err = gs.S.Ready()
	if err != nil {
		gs.Logger.Fatal().Err(err).Msg("Could not send ready message")
	}
	gs.Logger.Info().Msg("Ready")

	<-gs.Ctx.Done()
	os.Exit(0)
}

func (gs *GameServer) UdpListener(port *string) {
	gs.Logger.Info().Str("port", *port).Msg("Starting UDP server")
	conn, err := net.ListenPacket("udp", ":"+*port)
	if err != nil {
		gs.Logger.Error().Err(err).Msg("Could not start UDP server")
	}
	defer conn.Close()
	gs.udpReadWriteLoop(conn)
}

func (gs *GameServer) udpReadWriteLoop(conn net.PacketConn) {
	b := make([]byte, 1024)
	for {
		n, sender, err := conn.ReadFrom(b)
		if err != nil {
			gs.Logger.Error().Err(err).Msg("Could not read from udp stream")
		}
		txt := strings.TrimSpace(string(b[:n]))
		gs.Logger.Info().Str("sender", sender.String()).Str("txt", txt).Msg("Received UDP packet")

		response := gs.handleInput(txt, sender)
		if _, err := conn.WriteTo([]byte(response), sender); err != nil {
			gs.Logger.Error().Err(err).Msg("Could not write to udp stream")
		}

		if txt == "EXIT" {
			gs.Logger.Info().Msg("Shutting down")
			err := gs.S.Shutdown()
			if err != nil {
				gs.Logger.Error().Err(err).Msg("Could not shutdown")
			}
		}
	}
}

func (gs *GameServer) handleInput(txt string, sender net.Addr) string {
	switch txt {
	case "STATUS":
		return "OK"
	case "CRASH":
		gs.Logger.Info().Msg("Crashing")
		os.Exit(1)
	}
	return fmt.Sprintf("ACK: %s", txt)
}
