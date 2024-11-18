#!/bin/bash

# Copy example env files to actual env files if they don't exist
for service in product-service order-service inventory-service api-gateway; do
    if [ ! -f "./$service/.env" ]; then
        cp "./$service/.env.example" "./$service/.env"
        echo "Created .env file for $service"
    fi
done