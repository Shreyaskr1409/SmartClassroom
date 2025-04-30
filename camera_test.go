package main

import (
	"fmt"
	"os"
	"testing"
)

// TestCameraCapture tests if the camera can capture an image successfully
func TestCameraCapture(t *testing.T) {
	// Take a picture using the camera function
	img, err := TakePicture()
	if err != nil {
		t.Fatalf("Failed to take picture: %v", err)
	}

	// Check if the returned image is not nil
	if img == nil {
		t.Error("TakePicture returned nil image")
	}

	// Check if the dimensions of the image are reasonable
	bounds := img.Bounds()
	if bounds.Dx() <= 0 || bounds.Dy() <= 0 {
		t.Errorf("Image has invalid dimensions: %dx%d", bounds.Dx(), bounds.Dy())
	}

	fmt.Printf("Successfully captured image with dimensions: %dx%d\n", bounds.Dx(), bounds.Dy())

	// Check if the temporary file was created
	_, err = os.Stat("/tmp/image.jpg")
	if os.IsNotExist(err) {
		t.Error("Temporary image file was not created")
	}
}

// TestImageCleanup tests if the temporary image file gets cleaned up
func TestImageCleanup(t *testing.T) {
	// Make sure there's an image to clean up
	_, err := TakePicture()
	if err != nil {
		t.Skipf("Skipping cleanup test as camera failed: %v", err)
	}

	// Check if temporary file exists
	_, err = os.Stat("/tmp/image.jpg")
	if os.IsNotExist(err) {
		t.Error("Expected temporary image file not found")
	}

	// Clean up the temporary file
	err = os.Remove("/tmp/image.jpg")
	if err != nil {
		t.Errorf("Failed to clean up temporary image file: %v", err)
	}

	// Verify the file was removed
	_, err = os.Stat("/tmp/image.jpg")
	if !os.IsNotExist(err) {
		t.Error("Temporary image file was not removed properly")
	}
}
