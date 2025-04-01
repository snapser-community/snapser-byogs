import sys
import re
import subprocess
import os
import json

# Helper


def load_env(parent_folder, env_file='.env'):
    """Manually load a .env file into environment variables."""
    file_path = os.path.join(parent_folder, env_file)
    if os.path.exists(file_path):
        with open(file_path, 'r') as file:
            for line in file:
                # Strip whitespace and ignore comments and empty lines
                line = line.strip()
                if line and not line.startswith('#'):
                    key, value = line.split('=', 1)
                    os.environ[key.strip()] = value.strip()


# Check if enough arguments have been passed
if len(sys.argv) < 4:
    print("Usage: python create_snapend.py $companyId $gameId $byogsTag")
    sys.exit(1)

# Accessing user arguments
company_id = sys.argv[1]
game_id = sys.argv[2]
tag = sys.argv[3]

# Constants
CURRENT_DIR = os.path.dirname(os.path.realpath(__file__))
RESOURCES_DIR = os.path.join(CURRENT_DIR, 'snapser-resources')
load_env(RESOURCES_DIR)
SOURCE_SNAPEND_MANIFEST = os.getenv(
    'SNAPEND_MANIFEST_FILE', 'snapser-base-snapend-manifest.json')
DEST_SNAPEND_MANIFEST = os.getenv(
    'SNAPEND_MANIFEST_FILE', 'snapser-snapend-manifest.json')
SNAPS_VERSION = os.getenv('SNAPS_VERSION', 'production')
SNAPS_VERSION_SUFFIX = ''
if SNAPS_VERSION is not None and SNAPS_VERSION != 'production':
    SNAPS_VERSION_PREFIX = f'-{SNAPS_VERSION}'
AUTH_SNAP_VERSION = os.getenv(
    'AUTH_SNAP_VERSION', 'v0.47.0') + SNAPS_VERSION_SUFFIX
GSF_SNAP_VERSION = os.getenv(
    'GSF_SNAP_VERSION', 'v0.47.0') + SNAPS_VERSION_SUFFIX
LANGUAGE = os.getenv('LANGUAGE', 'python')
SNAPEND_ENV = os.getenv('SNAPEND_ENV', 'development')


# Validation
# Check if SNAPEND_ENV is either development or staging
if SNAPEND_ENV not in ['development', 'staging']:
    print("Invalid SNAPEND_ENV. Should be either development or staging.")
    sys.exit(1)

# Check version
# Check if the SNAPEND_MANIFEST file exists
if not os.path.exists(os.path.join(RESOURCES_DIR, SOURCE_SNAPEND_MANIFEST)):
    print(f"{SOURCE_SNAPEND_MANIFEST} file not found at {RESOURCES_DIR}/.")
    sys.exit(1)

# Action starts
# First lets update the Snapend manifest file
try:
    snapend_manifest = None
    with open(os.path.join(RESOURCES_DIR, SOURCE_SNAPEND_MANIFEST), 'r') as f:
        snapend_manifest = json.load(f)
        snapend_manifest['environment'] = 'DEVELOPMENT'
        found = False
        for snap in snapend_manifest['service_definitions']:
            if snap['id'] == 'auth':
                snap['version'] = AUTH_SNAP_VERSION
            if snap['id'] == 'game-server-fleets':
                snap['version'] = GSF_SNAP_VERSION
                found = True
        if not found:
            snapend_manifest['service_definitions'].append({
            "id": "game-server-fleets",
            "language": "go",
            "version": GSF_SNAP_VERSION,
            "author_id": "snapser",
            "category": "SNAPSER",
            "subcategory": "FLEET",
            "data_dependencies": [
                {
                "type": "POSTGRES",
                "postgres": {
                    "database": "game_server_fleets_settings"
                }
                },
                {
                "type": "REDIS",
                "redis": {}
                }
            ]
            })
    with open(os.path.join(RESOURCES_DIR, DEST_SNAPEND_MANIFEST), 'w') as f:
        json.dump(snapend_manifest, f, indent=4)
    print(f"Updated {DEST_SNAPEND_MANIFEST} file.")
except Exception as e:
    print(f"Failed to load {SOURCE_SNAPEND_MANIFEST} file.")
    print(e)
    sys.exit(1)

# Now lets create the Snapend
# snapctl snapend clone --game-id 6e0d68e9-a679-461d-a4c5-f5e29e655a87 --name byosnap-python-aj-demo --env development --manifest-path-filename "C:\Users\name\Downloads\snapser-ox1bcyim-manifest.json" --blocking
snapend_name = 'byogs-tutorial-basic'
print(
    f"Creating Snapend with name: {snapend_name}-demo in environment: {SNAPEND_ENV}")
cmd = [
    'snapctl', 'snapend', 'clone', '--game-id', game_id,
    '--name', f"{snapend_name}-demo", '--env', SNAPEND_ENV,
    '--manifest-path-filename',
    os.path.join(RESOURCES_DIR, DEST_SNAPEND_MANIFEST),
    '--blocking'
]

response = subprocess.run(cmd)

if response.returncode == 0:
    print(
        f"Your Snapend is created successfully with Snaps: auth and game-server-fleets.")
else:
    print(f"Failed to create a Snapend.")
    print('IMPORTANT: Check if you have already created a Snapend with the same name.')
    sys.exit(1)
