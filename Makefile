BINARY_NAME=fuse
BUILD_DIR=bin
CLI_BINARY=fuse
CLI_BUILD_DIR=$(BUILD_DIR)/cli
MAKEFLAGS += --silent

.DEFAULT_GOAL := run

.PHONY: run build cli cli-run cli-local cli-help tidy clean swagger swagger-fmt air

run:
	@go run -buildvcs=false cmd/api/main.go


build: 
	@echo "🔨 Building..."
	@mkdir -p $(BUILD_DIR)
	@go build -buildvcs=false -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/api/main.go
	@echo "✓ Build completed!"

cli:
	@echo "🔨 Building Fuse CLI..."
	@mkdir -p $(CLI_BUILD_DIR)
	@go build -buildvcs=false -o $(CLI_BUILD_DIR)/$(CLI_BINARY) ./cli/cmd/fuse
	@echo "✓ CLI built at $(CLI_BUILD_DIR)/$(CLI_BINARY)"

cli-run: cli
	@$(CLI_BUILD_DIR)/$(CLI_BINARY) $(ARGS)

cli-local: cli
	@$(CLI_BUILD_DIR)/$(CLI_BINARY) --api-url http://localhost:3000 $(ARGS)

cli-help: cli
	@$(CLI_BUILD_DIR)/$(CLI_BINARY) --help

tidy:
	@echo "📦 Tidying up dependencies..."
	@go mod tidy
	@echo "✓ Dependencies tidied!"

clean:
	@echo "🧹 Cleaning..."
	@go clean
	@rm -rf $(BUILD_DIR) tmp/ docs/api/
	@echo "✓ Cleanup complete!"

swagger:
	@echo "🔄 Generating swagger documentation..."
	@mkdir -p docs/api
	@go run github.com/swaggo/swag/cmd/swag init -g cmd/api/docs.go -o docs/api --parseDependency --parseInternal
	@echo "✓ Swagger docs generated!"

swagger-fmt:
	@echo "🔧 Formatting swagger comments..."
	@go run github.com/swaggo/swag/cmd/swag fmt -g cmd/api/docs.go
	@echo "✓ Swagger comments formatted!"

air:
	@echo "🚀 Starting live reload with Air..."
	@air
