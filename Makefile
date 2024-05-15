# Set the binary name
BINARY_NAME=cloudblocks

# Set the source directory
SRC_DIR=.

# Set the output directory for the binary
OUTPUT_DIR=$(HOME)

# Set the installation directory
INSTALL_DIR=$(HOME)/.cloudblocks

.PHONY: all build install clean uninstall

all: build

build:
	go build -o $(OUTPUT_DIR)/$(BINARY_NAME) $(SRC_DIR)

install:
	mkdir -p $(INSTALL_DIR)
	cp $(OUTPUT_DIR)/$(BINARY_NAME) $(INSTALL_DIR)
	chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@if [ "$$(uname)" = "Darwin" ]; then \
		if ! grep 'export PATH=$$PATH:$(INSTALL_DIR)' ~/.zshrc; then \
			echo '\n\nexport PATH=$$PATH:$(INSTALL_DIR)' >> ~/.zshrc; \
		fi; \
		$(INSTALL_DIR)/$(BINARY_NAME) init --workdir=$(INSTALL_DIR)/work/ --modulesdir=$(INSTALL_DIR)/modules; \
		source ~/.zshrc; \
	else \
		if ! grep 'export PATH=$$PATH:$(INSTALL_DIR)' ~/.bashrc; then \
			echo '\n\nexport PATH=$$PATH:$(INSTALL_DIR)' >> ~/.bashrc; \
		fi; \
		$(INSTALL_DIR)/$(BINARY_NAME) init --workdir=$(INSTALL_DIR)/work/ --modulesdir=$(INSTALL_DIR)/modules; \
		. ~/.bashrc; \
	fi

clean:
	rm -f $(OUTPUT_DIR)/$(BINARY_NAME)

uninstall:
	rm -rf $(INSTALL_DIR)
	@if [ "$$(uname)" = "Darwin" ]; then \
		sed -i '' '/export PATH=\$\$PATH:$(INSTALL_DIR)/d' ~/.zshrc; \
		source ~/.zshrc; \
	else \
		sed -i '/export PATH=\$\$PATH:$(INSTALL_DIR)/d' ~/.bashrc; \
		. ~/.bashrc; \
	fi