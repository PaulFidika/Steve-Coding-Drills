#!/bin/bash

# Get the directory of the script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Change to the script's directory
cd "$SCRIPT_DIR" || exit

# Ensure pb/chat directory exists and is empty
if [ ! -d "pb/chat" ]; then
    mkdir -p pb/chat
else
    rm -rf pb/chat/*
fi

# Ensure pb/coffeeshop directory exists and is empty
if [ ! -d "pb/coffeeshop" ]; then
    mkdir -p pb/coffeeshop
else
    rm -rf pb/coffeeshop/*
fi

# Generate protobuf files for chat.proto
protoc --proto_path=proto --go_out=pb/chat --go_opt=paths=source_relative \
    --go-grpc_out=pb/chat --go-grpc_opt=paths=source_relative \
    proto/chat.proto

# Generate protobuf files for coffeeshop.proto
protoc --proto_path=proto --go_out=pb/coffeeshop --go_opt=paths=source_relative \
    --go-grpc_out=pb/coffeeshop --go-grpc_opt=paths=source_relative \
    proto/coffeeshop.proto

echo "Proto files generated successfully in pb directory"
