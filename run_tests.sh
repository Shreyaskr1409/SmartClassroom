#!/bin/bash

# Script to run all tests for the face recognition system

echo "====================================================="
echo "Running Tests for Face Recognition System"
echo "====================================================="

# Function to check if a command exists
command_exists() {
    command -v "$1" &> /dev/null
}

# Check for required commands
if ! command_exists go; then
    echo "Error: Go is not installed or not in PATH"
    echo "Please run 'make install' first"
    exit 1
fi

if ! command_exists python3; then
    echo "Error: Python3 is not installed or not in PATH"
    echo "Please run 'make install' first"
    exit 1
fi

# Check if libcamera-jpeg is installed
if ! command_exists libcamera-jpeg; then
    echo "Error: libcamera-jpeg is not installed"
    echo "Please run 'make install' first"
    exit 1
fi

# Make sure camera is enabled
echo "Ensuring camera is enabled..."
sudo raspi-config nonint do_camera 0

# Start the ML API server in the background
echo -e "\n====================================================="
echo "Starting ML API server for testing..."
echo "====================================================="
cd mlapi && python3 app.py > ../mlapi_test.log 2>&1 &
ML_API_PID=$!
cd ..

# Give the server time to start
echo "Waiting for ML API server to start..."
sleep 3

# Run camera tests
echo -e "\n====================================================="
echo "Running Camera Tests..."
echo "====================================================="
go test -v camera_test.go camera.go
CAMERA_TEST_RESULT=$?

# Run sensor tests
echo -e "\n====================================================="
echo "Running Ultrasonic Sensor Tests..."
echo "====================================================="
go test -v sensor_test.go sensor.go
SENSOR_TEST_RESULT=$?

# Run ML API tests
echo -e "\n====================================================="
echo "Running ML API Tests..."
echo "====================================================="
cd mlapi && python3 -m unittest test_api.py
ML_API_TEST_RESULT=$?
cd ..

# Run integration test
echo -e "\n====================================================="
echo "Running Integration Test..."
echo "====================================================="
go test -v integration_test.go camera.go sensor.go

# Run cleanup test
echo -e "\n====================================================="
echo "Running Cleanup Test..."
echo "====================================================="
go test -v -run TestCleanup integration_test.go

# Stop the ML API server
echo -e "\n====================================================="
echo "Stopping ML API server..."
echo "====================================================="
kill $ML_API_PID

# Final cleanup
echo -e "\n====================================================="
echo "Running final cleanup..."
echo "====================================================="
bash cleanup.sh

# Test results summary
echo -e "\n====================================================="
echo "Test Results Summary:"
echo "====================================================="
echo "Camera Tests: $([ $CAMERA_TEST_RESULT -eq 0 ] && echo 'PASSED' || echo 'FAILED')"
echo "Sensor Tests: $([ $SENSOR_TEST_RESULT -eq 0 ] && echo 'PASSED' || echo 'FAILED')"
echo "ML API Tests: $([ $ML_API_TEST_RESULT -eq 0 ] && echo 'PASSED' || echo 'FAILED')"

# Overall result
if [ $CAMERA_TEST_RESULT -eq 0 ] && [ $SENSOR_TEST_RESULT -eq 0 ] && [ $ML_API_TEST_RESULT -eq 0 ]; then
    echo -e "\nAll tests PASSED!"
    exit 0
else
    echo -e "\nSome tests FAILED. Please check the logs."
    exit 1
fi
