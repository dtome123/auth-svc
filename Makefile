# ===============================
# Global Config
# ===============================
BIN_DIR := bin

# ===============================
# Tool: genrsa
# ===============================
GENRSA_NAME := genrsa
GENRSA_MAIN := ./tools/rsa/main.go
GENRSA_BIN := $(BIN_DIR)/$(GENRSA_NAME)

## Build the genrsa tool
build-genrsa:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(GENRSA_BIN) $(GENRSA_MAIN)
## Run the genrsa tool with output directory
run-genrsa:
	@$(GENRSA_BIN) --out-dir=$(out)

# ===============================
# Clean all
# ===============================
.PHONY: clean

clean:
	rm -rf $(BIN_DIR)

