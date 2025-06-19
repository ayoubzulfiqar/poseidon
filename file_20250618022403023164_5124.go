package main

import (
	"fmt"
	"math/rand"
	"time"
)

var currentBatteryLevel = 100

func getBatteryLevel() int {
	if currentBatteryLevel > 0 {
		currentBatteryLevel -= rand.Intn(5) + 1
		if currentBatteryLevel < 0 {
			currentBatteryLevel = 0
		}
	}
	return currentBatteryLevel
}

func monitorBattery(lowThreshold, criticalThreshold int, checkInterval time.Duration) {
	fmt.Println("Battery Monitor Started.")
	fmt.Printf("Low Alert: %d%%\n", lowThreshold)
	fmt.Printf("Critical Alert: %d%%\n", criticalThreshold)
	fmt.Printf("Check Interval: %s\n", checkInterval)

	for {
		level := getBatteryLevel()
		fmt.Printf("Current Level: %d%%\n", level)

		if level <= criticalThreshold {
			fmt.Printf("ALERT: CRITICAL BATTERY! %d%%\n", level)
		} else if level <= lowThreshold {
			fmt.Printf("ALERT: LOW BATTERY! %d%%\n", level)
		} else if level == 100 {
			fmt.Println("Battery is fully charged.")
		} else if level == 0 {
			fmt.Println("Battery is drained.")
		}

		time.Sleep(checkInterval)

		if level == 0 {
			fmt.Println("Simulating recharge...")
			currentBatteryLevel = 100
			time.Sleep(5 * time.Second)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	low := 20
	critical := 10
	interval := 3 * time.Second

	monitorBattery(low, critical, interval)
}

// Additional implementation at 2025-06-18 02:24:44
