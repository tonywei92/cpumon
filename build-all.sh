#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Set the project name
PROJECT_NAME="cpumon"

# Define the output directory
BUILD_DIR="build"
mkdir -p "$BUILD_DIR"

# Array of targets: "OS ARCH EXTENSION COMPRESSION"
# EXTENSION is for binary files (e.g., .exe for Windows)
# COMPRESSION is either "zip" or "tar.gz"
targets=(
  "windows amd64 .exe zip"
  "windows 386 .exe zip"
  "linux amd64 '' tar.gz"
  "linux 386 '' tar.gz"
  "linux arm '' tar.gz"
  "linux arm64 '' tar.gz"
  "darwin amd64 '' tar.gz"
  "darwin arm64 '' tar.gz"
)

# Function to build and package for a single target
build_and_package() {
  local OS=$1
  local ARCH=$2
  local EXT=$3
  local COMPRESS=$4

  # Determine the output binary name
  if [ "$OS" == "windows" ]; then
    BINARY_NAME="${PROJECT_NAME}${EXT}"
  else
    BINARY_NAME="${PROJECT_NAME}"
  fi

  # Determine the final archive name
  if [ "$OS" == "windows" ]; then
    ARCHIVE_NAME="${PROJECT_NAME}-${OS}-${ARCH}.zip"
  else
    ARCHIVE_NAME="${PROJECT_NAME}-${OS}-${ARCH}.tar.gz"
  fi

  echo "----------------------------------------"
  echo "Building for OS: $OS, ARCH: $ARCH"

  # Set environment variables for cross-compilation
  export GOOS="$OS"
  export GOARCH="$ARCH"

  # Handle specific architectures if needed (e.g., GOARM for ARM)
  if [ "$ARCH" == "arm" ]; then
    export GOARM=7  # You can adjust this based on your requirements
  fi

  # Build the binary
  if [ "$OS" == "windows" ]; then
    OUTPUT_PATH="$BUILD_DIR/$BINARY_NAME"
  else
    OUTPUT_PATH="$BUILD_DIR/$BINARY_NAME"
  fi

  echo "Compiling $OUTPUT_PATH..."
  go build -o "$OUTPUT_PATH"

  # Compress the binary
  if [ "$COMPRESS" == "zip" ]; then
    echo "Compressing $BINARY_NAME into $ARCHIVE_NAME..."
    (cd "$BUILD_DIR" && zip -q "$ARCHIVE_NAME" "$BINARY_NAME")
  else
    echo "Compressing $BINARY_NAME into $ARCHIVE_NAME..."
    (cd "$BUILD_DIR" && tar -czf "$ARCHIVE_NAME" "$BINARY_NAME")
  fi

  # Remove the raw binary after compression
  rm "$OUTPUT_PATH"

  echo "Packaged $ARCHIVE_NAME successfully."
}

# Iterate over each target and build/package
for target in "${targets[@]}"; do
  # Read target components
  read -r OS ARCH EXT COMPRESS <<< "$target"
  build_and_package "$OS" "$ARCH" "$EXT" "$COMPRESS"
done

echo "----------------------------------------"
echo "All builds and packaging completed successfully."
echo "Artifacts are located in the '$BUILD_DIR' directory."
