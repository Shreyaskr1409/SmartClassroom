#!/bin/bash

# Script to clean up temporary files and processes created during testing

echo "Cleaning up temporary files..."

# Remove temporary image files
echo "Removing temporary image files..."
rm -f /tmp/image.jpg
rm -f /tmp/test_face.jpg
echo "Temporary image files removed."

# Check and kill any test processes that might still be running
echo "Checking for running test processes..."

# Check for Python Flask API process
FLASK_PID=$(pgrep -f "python3 app.py")
if [ ! -z "$FLASK_PID" ]; then
    echo "Killing Flask API process (PID: $FLASK_PID)..."
    kill $FLASK_PID
    echo "Flask API process terminated."
else
    echo "No Flask API process found."
fi

# Check for Go test processes
GO_TEST_PID=$(pgrep -f "go test")
if [ ! -z "$GO_TEST_PID" ]; then
    echo "Killing Go test process (PID: $GO_TEST_PID)..."
    kill $GO_TEST_PID
    echo "Go test process terminated."
else
    echo "No Go test process found."
fi

echo "Cleanup completed successfully!"
