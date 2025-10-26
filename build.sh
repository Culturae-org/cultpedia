#!/bin/bash

#---------------------------------------------------
#         Cultpedia on Linux and macOS             | 
#---------------------------------------------------

set -e

echo "Building Cultpedia..."
go build -o cultpedia ./cmd

echo "✔ Build successful!"
echo ""
echo "You can now run: ./cultpedia"
echo ""
echo "Usage:"
echo "  Interactive mode:  ./cultpedia"
echo "  Commands:          ./cultpedia help"
