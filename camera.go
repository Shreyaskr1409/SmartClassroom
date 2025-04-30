package main

import (
	"fmt"
	"image"
	"os"
	"os/exec"
)

// TakePicture captures an image from the Raspberry Pi camera and returns it
func TakePicture() (image.Image, error) {
	cmd := exec.Command("libcamera-jpeg", "-o", "/tmp/image.jpg")
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to take picture: %v", err)
	}

	file, err := os.Open("/tmp/image.jpg")
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}
	return img, nil
}
