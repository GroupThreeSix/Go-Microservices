#!/bin/bash

# Configuration
DOCKER_REGISTRY="tuilakhanh"
VERSION="latest"

# Colors for output
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Function to build and push a service
build_and_push() {
    local service=$1
    echo -e "${GREEN}Building ${service}...${NC}"
    
    # Build the image
    podman build -t ${DOCKER_REGISTRY}/${service}:${VERSION} .
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}Pushing ${service}...${NC}"
        podman push ${DOCKER_REGISTRY}/${service}:${VERSION}
    else
        echo "Failed to build ${service}"
        exit 1
    fi
}

# Build and push each service
echo -e "${GREEN}Starting build and push process...${NC}"

# Build and push inventory service
cd inventory-service
build_and_push "inventory-service"
cd ..

# Build and push product service
cd product-service
build_and_push "product-service"
cd ..

# Build and push order service
cd order-service
build_and_push "order-service"
cd ..

# Build and push api gateway
cd api-gateway
build_and_push "api-gateway"
cd ..

echo -e "${GREEN}All services have been built and pushed successfully!${NC}"

# List all images
echo -e "${GREEN}Local images:${NC}"
podman images | grep ${DOCKER_REGISTRY}