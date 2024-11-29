.PHONY: proto clean build run

# Generate protobuf code
# proto:
# 	@echo "Generating protobuf code..."
# 	@for service in product-service inventory-service; do \
# 		protoc --proto_path=proto \
# 			--go_out=$$service/proto --go_opt=paths=source_relative \
# 			--go-grpc_out=$$service/proto --go-grpc_opt=paths=source_relative \
# 			proto/*.proto; \
# 	done

proto:
	@echo "Generating protobuf code..."
	@for service in product-service inventory-service order-service; do \
		protoc --proto_path=proto \
			--go_out=$$service \
			--go-grpc_out=$$service \
			proto/*.proto; \
	done

# Clean generated code
clean:
	@echo "Cleaning generated code..."
	@rm -rf */proto/*

# Build all services
build: proto
	@echo "Building services..."
	docker-compose build

# Run all services
run: build
	@echo "Starting services..."
	docker-compose up

# Stop all services
stop:
	@echo "Stopping services..."
	docker-compose down

# Generate proto and run services in development mode
dev: proto
	@echo "Starting services in development mode..."
	docker-compose up --build