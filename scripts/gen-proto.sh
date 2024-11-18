#!/bin/bash

# Create proto directories if they don't exist
mkdir -p product-service/proto
mkdir -p inventory-service/proto

# Generate proto files
make proto

# Check if generation was successful
# if [ $? -eq 0 ]; then
#     echo "Proto files generated successfully!"
    
#     # Build and start services
#     make build
# else
#     echo "Failed to generate proto files!"
#     exit 1
# fi