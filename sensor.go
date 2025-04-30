package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
)

// Sensor setup
var (
	trig = rpi.P1_16 // GPIO23
	echo = rpi.P1_18 // GPIO24
)

// Distance threshold
const distanceThreshold = 150 // in cm (1.5 meters)

func init() {
	// Initialize GPIO
	_, err := host.Init()
	if err != nil {
		log.Fatalf("Failed to initialize host: %v", err)
	}

	if err := trig.Out(gpio.Low); err != nil {
		log.Fatalf("Failed to set trig pin: %v", err)
	}

	if err := echo.In(gpio.PullDown, gpio.FallingEdge); err != nil { // Using FallingEdge for detection
		log.Fatalf("Failed to set echo pin: %v", err)
	}
}

// GetDistance reads the ultrasonic sensor and returns the distance in centimeters
func GetDistance() float64 {
	// Send trigger pulse
	trig.Out(gpio.High)
	time.Sleep(10 * time.Microsecond)
	trig.Out(gpio.Low)

	// Wait for the echo to go high and record the start time
	for echo.Read() == gpio.Low {
	}
	start := time.Now()

	// Wait for the echo to go low and record the end time
	for echo.Read() == gpio.High {
	}
	duration := time.Since(start)

	// Calculate distance
	distance := duration.Seconds() * 34300 / 2
	return distance
}

// SendImage sends the captured image to a URL
func SendImage(img image.Image) (string, error) {
	var buf bytes.Buffer
	// Encode the image into a buffer
	err := jpeg.Encode(&buf, img, nil)
	if err != nil {
		return "", fmt.Errorf("failed to encode image: %v", err)
	}

	// Send the image as HTTP POST request
	resp, err := http.Post("http://example.com/upload", "image/jpeg", &buf)
	if err != nil {
		return "", fmt.Errorf("failed to send image: %v", err)
	}
	defer resp.Body.Close()

	// Read and return the response
	body, err := io.ReadAll(resp.Body) // Replaced ioutil.ReadAll with io.ReadAll
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}
	return string(body), nil
}

// MonitorSensor continuously monitors the ultrasonic sensor and triggers picture taking and sending
func MonitorSensor() {
	for {
		// Get the distance reading
		distance := GetDistance()
		fmt.Printf("Distance: %.2f cm\n", distance)

		if distance < distanceThreshold {
			// Take a picture if the distance is below threshold
			fmt.Println("Object detected within range. Taking picture...")
			img, err := TakePicture()
			if err != nil {
				log.Printf("Error taking picture: %v", err)
				continue
			}

			// Send the image to the server
			fmt.Println("Sending image to server...")
			response, err := SendImage(img)
			if err != nil {
				log.Printf("Error sending image: %v", err)
				continue
			}

			// Print the response from the server
			fmt.Printf("Server response: %s\n", response)
		}

		// Wait before the next reading
		time.Sleep(1 * time.Second)
	}
}
