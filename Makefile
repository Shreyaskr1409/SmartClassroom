GO_VERSION=1.20.5
GO_BIN_DIR=/usr/local/go
LIBCAMERA_PACKAGE=libcamera-apps

install:
	@echo "Installing Go version $(GO_VERSION)..."
	# Download and install Go
	wget https://golang.org/dl/go$(GO_VERSION).linux-armv6l.tar.gz -O /tmp/go.tar.gz
	sudo tar -C /usr/local -xvzf /tmp/go.tar.gz
	rm /tmp/go.tar.gz
	# Set up Go path
	echo "export PATH=$(GO_BIN_DIR)/bin:\$$PATH" >> ~/.bashrc
	source ~/.bashrc

	@echo "Installing necessary Go packages..."
	# Install Go packages (periph.io for GPIO and other dependencies)
	go get -u periph.io/x/host/v3
	go get -u periph.io/x/conn/v3/gpio

	@echo "Installing libcamera-jpeg..."
	# Install libcamera-jpeg package
	sudo apt-get update
	sudo apt-get install -y $(LIBCAMERA_PACKAGE)

	@echo "Go, necessary packages, and libcamera-jpeg have been installed."

# Run the program (initialize camera and run main.go)
run:
	@echo "Initializing camera..."
	# Make sure the camera is enabled
	sudo raspi-config nonint do_camera 0
	# Run the main Go application
	go run main.go
