#!/bin/bash

# Set the directory for the issuing files
OUTPUT_DIR="build"
mkdir -p $OUTPUT_DIR

go mod tidy

# Build the Linux Executable
echo "Build Linux Executable..."
if ! GOOS=linux GOARCH=amd64 go build -v -o $OUTPUT_DIR/or2b; then
    echo "Error building the Linux Executable"
    exit 1
fi
echo '----------------------------------------------------------------'

# Build the Windows Executable
echo "Build Windows Executable..."
if ! GOOS=windows GOARCH=amd64 go build -v -o $OUTPUT_DIR/or2b.exe; then
    echo "Error building the Windows Executable"
    exit 1
fi
echo '----------------------------------------------------------------'
if [ "$HOSTNAME" == "victnix" ]; then
    sudo cp -v build/or2b /usr/local/bin/
else
    cp -v build/or2b.exe /gast-drive-d/Daten/
fi

echo "Build successfully completed!"
