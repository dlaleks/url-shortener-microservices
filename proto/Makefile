# Proto Generation Makefile for URL Shortener Microservices
# Dependencies: protoc, protoc-gen-go, protoc-gen-go-grpc

.PHONY: all clean deps gen-common gen-url gen-user gen-analytics install-tools

# Default target
all: deps gen

# Install required protoc plugins
install-tools:
	@echo "Installing protoc plugins..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Tools installed successfully"

# Check if required tools are installed
deps:
	@echo "Checking dependencies..."
	@which protoc > /dev/null || (echo "protoc is required but not installed"; exit 1)
	@which protoc-gen-go > /dev/null || (echo "protoc-gen-go is required. Run 'make install-tools'"; exit 1)
	@which protoc-gen-go-grpc > /dev/null || (echo "protoc-gen-go-grpc is required. Run 'make install-tools'"; exit 1)
	@echo "All dependencies are installed"

# Generate all proto files
gen: gen-common gen-url gen-user gen-analytics
	@echo "All proto files generated successfully"

# Generate common types
gen-common:
	@echo "Generating common proto files..."
	@mkdir -p gen/common
	protoc \
		--proto_path=. \
		--go_out=gen \
		--go_opt=paths=source_relative \
		common/types.proto

# Generate URL service proto
gen-url:
	@echo "Generating URL service proto files..."
	@mkdir -p gen/url
	protoc \
		--proto_path=. \
		--go_out=gen \
		--go_opt=paths=source_relative \
		--go-grpc_out=gen \
		--go-grpc_opt=paths=source_relative \
		url/url_service.proto

# Generate User service proto
gen-user:
	@echo "Generating User service proto files..."
	@mkdir -p gen/user
	protoc \
		--proto_path=. \
		--go_out=gen \
		--go_opt=paths=source_relative \
		--go-grpc_out=gen \
		--go-grpc_opt=paths=source_relative \
		user/user_service.proto

# Generate Analytics service proto
gen-analytics:
	@echo "Generating Analytics service proto files..."
	@mkdir -p gen/analytics
	protoc \
		--proto_path=. \
		--go_out=gen \
		--go_opt=paths=source_relative \
		--go-grpc_out=gen \
		--go-grpc_opt=paths=source_relative \
		analytics/analytics_service.proto

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	rm -rf gen/
	@echo "Generated files cleaned"

# Help
help:
	@echo "Available targets:"
	@echo "  all            - Install dependencies and generate all proto files"
	@echo "  install-tools  - Install required protoc plugins"
	@echo "  deps           - Check if dependencies are installed"
	@echo "  gen            - Generate all proto files"
	@echo "  gen-common     - Generate common types only"
	@echo "  gen-url        - Generate URL service proto only"
	@echo "  gen-user       - Generate User service proto only"
	@echo "  gen-analytics  - Generate Analytics service proto only"
	@echo "  clean          - Remove generated files"
	@echo "  help           - Show this help"