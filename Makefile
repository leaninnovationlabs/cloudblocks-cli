# Set the binary name
BINARY_NAME=cloudblocks

# Set the source directory
SRC_DIR=.

# Set the output directory for the binary
OUTPUT_DIR=$(HOME)

# Set the installation directory
INSTALL_DIR=$(HOME)/.cloudblocks

.PHONY: all build install clean

all: build

build:
	go build -o $(OUTPUT_DIR)/$(BINARY_NAME) $(SRC_DIR)

install:
	mkdir -p $(INSTALL_DIR)
	cp $(OUTPUT_DIR)/$(BINARY_NAME) $(INSTALL_DIR)
	chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	echo 'export PATH=$$PATH:$(INSTALL_DIR)' >> ~/.bashrc
	sudo mkdir -p /opt/work
	sudo mkdir -p /opt/modules
	sudo chown -R $$USER:$$USER /opt/work
	sudo chown -R $$USER:$$USER /opt/modules
	$(INSTALL_DIR)/$(BINARY_NAME) init --workdir=/opt/work/ --modulesdir=/opt/modules
	. ~/.bashrc

clean:
	rm -f $(OUTPUT_DIR)/$(BINARY_NAME)
uninstall:
	rm -rf $(INSTALL_DIR)
	sed -i '/export PATH=\$\$PATH:\/usr\/local\/bin\/cloudblocks/d' ~/.bashrc
	sed -i '/export PATH=\$\$PATH:$(INSTALL_DIR)/d' ~/.bashrc
	. ~/.bashrc