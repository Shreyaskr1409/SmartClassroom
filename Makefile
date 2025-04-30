GO_VERSION=1.20.5
LIBCAMERA_PACKAGE=libcamera-apps

install:
	@echo "Installing Python and Go..."
	sudo apt-get update
	# Install Python and pip
	sudo apt-get install -y python3 python3-pip
	# Install Go
	sudo apt-get install -y golang-go
	# Install libcamera-jpeg package
	sudo apt-get install -y $(LIBCAMERA_PACKAGE)

	@echo "Installing Python dependencies from requirements.txt..."
	pip3 install -r mlapi/requirements.txt

	@echo "Tidying Go modules..."
	go mod tidy

	@echo "Installation complete."

run:
	sudo raspi-config nonint do_camera 0
	@echo "Running servers..."

	cd mlapi && nohup python3 app.py > ../mlapi.log 2>&1 &

	go run main.go

test:
	chmod +x run_tests.sh cleanup.sh
	./run_tests.sh
