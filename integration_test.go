package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

// TestIntegration tests the full integration flow: sensor -> camera -> ML API
func TestIntegration(t *testing.T) {
	fmt.Println("Starting integration test...")
	fmt.Println("This test will simulate the full flow:")
	fmt.Println("1. Ultrasonic sensor detection")
	fmt.Println("2. Camera capture")
	fmt.Println("3. ML API prediction")

	// Check if ML API server is running
	fmt.Println("\nChecking if ML API server is running...")
	if !isMLServerRunning() {
		t.Fatal("ML API server is not running. Please start it with 'cd mlapi && python3 app.py'")
	}
	fmt.Println("ML API server is running!")

	// Take a distance reading
	fmt.Println("\nTaking distance reading...")
	distance := GetDistance()
	fmt.Printf("Current distance: %.2f cm (Threshold: %d cm)\n", distance, distanceThreshold)

	// Simulate object detection
	fmt.Println("\nPlease place an object in front of the sensor (within threshold distance)")
	fmt.Println("Waiting 5 seconds...")
	time.Sleep(5 * time.Second)

	// Take new distance reading after object placement
	newDistance := GetDistance()
	fmt.Printf("New distance: %.2f cm\n", newDistance)

	if newDistance >= distanceThreshold {
		fmt.Println("Object not detected within threshold. Continuing test anyway for demonstration...")
	} else {
		fmt.Println("Object detected within threshold!")
	}

	// Take picture
	fmt.Println("\nCapturing image...")
	img, err := TakePicture()
	if err != nil {
		t.Fatalf("Failed to take picture: %v", err)
	}
	fmt.Printf("Image captured successfully with dimensions: %dx%d\n", img.Bounds().Dx(), img.Bounds().Dy())

	// Send image to ML API
	fmt.Println("\nSending image to ML API for face recognition...")
	response, err := SendImage(img)
	if err != nil {
		t.Fatalf("Failed to send image to ML API: %v", err)
	}

	// Parse and display the response
	fmt.Println("Response received from ML API:")
	fmt.Println(response)

	var responseData map[string]interface{}
	err = json.Unmarshal([]byte(response), &responseData)
	if err != nil {
		t.Errorf("Failed to parse API response as JSON: %v", err)
	} else {
		// Check if faces were detected
		if results, ok := responseData["results"].([]interface{}); ok {
			fmt.Printf("ML API detected %d face(s)\n", len(results))
			for i, result := range results {
				resultMap := result.(map[string]interface{})
				if recognized, ok := resultMap["recognized"].(bool); ok && recognized {
					name := resultMap["name"].(string)
					fmt.Printf("Face %d recognized as: %s\n", i+1, name)
				} else {
					fmt.Printf("Face %d not recognized\n", i+1)
				}
			}
		} else if errMsg, ok := responseData["error"].(string); ok && errMsg == "No face detected." {
			fmt.Println("No faces detected in the image")
		}
	}

	// Clean up
	fmt.Println("\nCleaning up...")
	cleanupTestFiles()

	fmt.Println("\nIntegration test completed!")
}

// Helper function to check if ML API server is running
func isMLServerRunning() bool {
	// Simple HTTP HEAD request to check if server is up
	_, err := http.Head("http://localhost:3300/")
	return err == nil
}

// Helper function to clean up temporary files created during testing
func cleanupTestFiles() {
	filesToCleanup := []string{
		"/tmp/image.jpg",
		"/tmp/test_face.jpg",
	}

	for _, file := range filesToCleanup {
		if _, err := os.Stat(file); err == nil {
			err = os.Remove(file)
			if err != nil {
				fmt.Printf("Failed to remove temporary file %s: %v\n", file, err)
			} else {
				fmt.Printf("Removed temporary file: %s\n", file)
			}
		}
	}
}

// TestCleanup is a standalone test just for cleaning up temporary files
func TestCleanup(t *testing.T) {
	fmt.Println("Running cleanup test...")
	cleanupTestFiles()
	fmt.Println("Cleanup completed")
}
