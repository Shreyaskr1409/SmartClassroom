package main

import (
	"fmt"
	"testing"
	"time"
)

// TestUltrasonicSensor tests if the ultrasonic sensor is working correctly
func TestUltrasonicSensor(t *testing.T) {
	// Take multiple readings to account for potential fluctuations
	const numReadings = 5

	fmt.Println("Testing ultrasonic sensor...")
	fmt.Println("Please ensure there are no objects within 2 meters of the sensor")
	time.Sleep(3 * time.Second) // Give time to clear the area

	var readings []float64
	var sum float64

	// Take multiple readings
	for i := 0; i < numReadings; i++ {
		distance := GetDistance()
		readings = append(readings, distance)
		sum += distance
		fmt.Printf("Reading %d: %.2f cm\n", i+1, distance)
		time.Sleep(500 * time.Millisecond) // Wait between readings
	}

	// Calculate average
	avg := sum / float64(numReadings)
	fmt.Printf("Average distance: %.2f cm\n", avg)

	// Check if readings are within reasonable range
	// For an open area, we expect readings to be greater than a minimum threshold
	const minExpectedDistance = 20.0 // cm
	if avg < minExpectedDistance {
		t.Errorf("Average distance reading (%.2f cm) is less than expected minimum (%.2f cm)", avg, minExpectedDistance)
	}

	// Check for consistency in readings (standard deviation)
	var sumSquaredDiffs float64
	for _, reading := range readings {
		diff := reading - avg
		sumSquaredDiffs += diff * diff
	}
	stdDev := (sumSquaredDiffs / float64(numReadings))

	// Allowing for some variance in readings, but not too much
	const maxAcceptableStdDev = 10.0 // cm
	if stdDev > maxAcceptableStdDev {
		t.Errorf("Sensor readings have high variance (std dev: %.2f cm)", stdDev)
	} else {
		fmt.Printf("Sensor readings are consistent (variance: %.2f)\n", stdDev)
	}
}

// TestObjectDetection tests if the sensor can detect an object placed in front of it
func TestObjectDetection(t *testing.T) {
	fmt.Println("\nTesting object detection...")
	fmt.Println("Starting with no object, then please place an object about 1 meter from the sensor when prompted")

	// First reading without object
	time.Sleep(2 * time.Second)
	initialDistance := GetDistance()
	fmt.Printf("Initial distance (no object): %.2f cm\n", initialDistance)

	// Prompt to place object
	fmt.Println("Please place an object about 1 meter from the sensor now")
	time.Sleep(5 * time.Second) // Give time to place object

	// Take reading with object
	objectDistance := GetDistance()
	fmt.Printf("Distance with object: %.2f cm\n", objectDistance)

	// Check if there's a significant change in distance
	if objectDistance >= initialDistance || (initialDistance-objectDistance) < 20 {
		t.Errorf("Object detection failed: Change in distance not significant (%.2f cm -> %.2f cm)", initialDistance, objectDistance)
	} else {
		fmt.Printf("Object detected! Distance changed by %.2f cm\n", initialDistance-objectDistance)
	}

	fmt.Println("Please remove the object")
	time.Sleep(3 * time.Second)
}

// TestDistanceThreshold tests if the sensor correctly determines when an object is within the defined threshold
func TestDistanceThreshold(t *testing.T) {
	fmt.Println("\nTesting distance threshold...")

	// Get current distance
	currentDistance := GetDistance()
	fmt.Printf("Current distance: %.2f cm (Threshold: %.2f cm)\n", currentDistance, distanceThreshold)

	// Check if current reading would trigger the detection
	isWithinThreshold := currentDistance < distanceThreshold
	fmt.Printf("Is object within threshold? %v\n", isWithinThreshold)

	// Instructions for manual verification
	if isWithinThreshold {
		fmt.Println("An object is currently within the detection threshold")
		fmt.Println("Please remove any objects and test again if this is unexpected")
	} else {
		fmt.Println("No object is currently within the detection threshold")
		fmt.Printf("Place an object closer than %.2f cm to trigger detection\n", distanceThreshold)
	}
}
