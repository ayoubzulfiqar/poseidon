package main

import (
	"fmt"
	"time"
)

// ConvertTime converts a given time.Time object from its current location
// to a new specified time zone.
// The input 't' must have a location associated with it (e.g., from time.Now() or time.ParseInLocation).
// The 'toZone' parameter is the name of the target time zone (e.g., "America/New_York", "Europe/London").
// It returns the converted time.Time object in the target zone or an error if the target zone is invalid.
func ConvertTime(t time.Time, toZone string) (time.Time, error) {
	loc, err := time.LoadLocation(toZone)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to load target time zone %q: %w", toZone, err)
	}
	return t.In(loc), nil
}

func main() {
	// Example 1: Convert current UTC time to America/New_York
	nowUTC := time.Now().UTC()
	fmt.Printf("Current UTC time: %s\n", nowUTC.Format(time.RFC3339))

	nyTime, err := ConvertTime(nowUTC, "America/New_York")
	if err != nil {
		fmt.Printf("Error converting to New York time: %v\n", err)
	} else {
		fmt.Printf("Current New York time: %s\n", nyTime.Format(time.RFC3339))
	}

	fmt.Println("---")

	// Example 2: Convert a specific time from America/Los_Angeles to Europe/London
	// First, parse the time string into a time.Time object with the correct location.
	laLoc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		fmt.Printf("Error loading LA location: %v\n", err)
		return
	}

	// A specific time in LA, e.g., 2023-10-27 10:00:00 PST
	timeStrLA := "2023-10-27 10:00:00"
	layout := "2006-01-02 15:04:05"
	specificTimeLA, err := time.ParseInLocation(layout, timeStrLA, laLoc)
	if err != nil {
		fmt.Printf("Error parsing time in LA: %v\n", err)
		return
	}
	fmt.Printf("Specific time in Los Angeles: %s\n", specificTimeLA.Format(time.RFC3339))

	londonTime, err := ConvertTime(specificTimeLA, "Europe/London")
	if err != nil {
		fmt.Printf("Error converting to London time: %v\n", err)
	} else {
		fmt.Printf("Converted time in London: %s\n", londonTime.Format(time.RFC3339))
	}

	fmt.Println("---")

	// Example 3: Convert a specific time from Europe/Berlin to Asia/Tokyo
	berlinLoc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		fmt.Printf("Error loading Berlin location: %v\n", err)
		return
	}

	timeStrBerlin := "2024-01-15 14:30:00" // 2:30 PM in Berlin
	specificTimeBerlin, err := time.ParseInLocation(layout, timeStrBerlin, berlinLoc)
	if err != nil {
		fmt.Printf("Error parsing time in Berlin: %v\n", err)
		return
	}
	fmt.Printf("Specific time in Berlin: %s\n", specificTimeBerlin.Format(time.RFC3339))

	tokyoTime, err := ConvertTime(specificTimeBerlin, "Asia/Tokyo")
	if err != nil {
		fmt.Printf("Error converting to Tokyo time: %v\n", err)
	} else {
		fmt.Printf("Converted time in Tokyo: %s\n", tokyoTime.Format(time.RFC3339))
	}

	fmt.Println("---")

	// Example 4: Invalid time zone
	_, err = ConvertTime(time.Now(), "Invalid/TimeZone")
	if err != nil {
		fmt.Printf("Expected error for invalid time zone: %v\n", err)
	}
}

// Additional implementation at 2025-06-20 22:35:52
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// getInput reads a line of text from stdin, trims whitespace, and returns it.
func getInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
	return strings.TrimSpace(input)
}

func main() {
	// Define common time formats for parsing
	const (
		dateTimeFormat = "2006-01-02 15:04:05"
		timeOnlyFormat = "15:04:05"
	)

	fmt.Println("Time Zone Converter")
	fmt.Println("-------------------")

	// Get input time string from user
	timeStr := getInput("Enter time (e.g., '2023-10-27 15:30:00' or '15:30:00'): ")
	if timeStr == "" {
		log.Fatalf("Error: Time input cannot be empty.")
	}

	// Get source time zone from user
	sourceTZ := getInput("Enter source time zone (e.g., 'America/New_York', 'UTC', 'Europe/London'): ")
	if sourceTZ == "" {
		log.Fatalf("Error: Source time zone cannot be empty.")
	}

	// Get target time zone from user
	targetTZ := getInput("Enter target time zone (e.g., 'Asia/Tokyo', 'Australia/Sydney'): ")
	if targetTZ == "" {
		log.Fatalf("Error: Target time zone cannot be empty.")
	}

	// Load source time zone location
	sourceLoc, err := time.LoadLocation(sourceTZ)
	if err != nil {
		log.Fatalf("Error loading source time zone '%s': %v", sourceTZ, err)
	}

	// Load target time zone location
	targetLoc, err := time.LoadLocation(targetTZ)
	if err != nil {
		log.Fatalf("Error loading target time zone '%s': %v", targetTZ, err)
	}

	var parsedTime time.Time

	// Attempt to parse the input as a full date-time string first
	parsedTime, err = time.ParseInLocation(dateTimeFormat, timeStr, sourceLoc)
	if err != nil {
		// If full date-time parsing fails, try parsing as time-only
		// time.Parse defaults to Jan 1, year 0, UTC if no date/timezone info is present
		tempTime, errTimeOnly := time.Parse(timeOnlyFormat, timeStr)
		if errTimeOnly != nil {
			log.Fatalf("Error parsing time '%s'. Please use format 'YYYY-MM-DD HH:MM:SS' or 'HH:MM:SS': %v", timeStr, err)
		}

		// If time-only parsing succeeded, combine it with today's date in the source timezone
		nowInSourceLoc := time.Now().In(sourceLoc)
		parsedTime = time.Date(
			nowInSourceLoc.Year(),
			nowInSourceLoc.Month(),
			nowInSourceLoc.Day(),
			tempTime.Hour(),
			tempTime.Minute(),
			tempTime.Second(),
			tempTime.Nanosecond(),
			sourceLoc,
		)
	}

	fmt.Printf("\nOriginal Time: %s in %s\n", parsedTime.Format(dateTimeFormat), sourceLoc.String())

	// Convert the time to the target time zone
	convertedTime := parsedTime.In(targetLoc)

	fmt.Printf("Converted Time: %s in %s\n", convertedTime.Format(dateTimeFormat), targetLoc.String())
	fmt.Println("-------------------")

	// Additional functionality: Displaying the offset difference
	// Zone() returns the time zone name and offset (in seconds east of UTC) for the time t.
	_, sourceOffset := parsedTime.Zone()
	_, targetOffset := convertedTime.Zone()
	offsetDiffHours := float64(targetOffset-sourceOffset) / 3600.0

	fmt.Printf("Time Zone Offset Difference (Target - Source): %.1f hours\n", offsetDiffHours)
	fmt.Println("Note: This offset difference reflects the specific time and may vary due to Daylight Saving Time (DST).")
}

// Additional implementation at 2025-06-20 22:36:35
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func getTimeZone(zoneName string) (*time.Location, error) {
	loc, err := time.LoadLocation(zoneName)
	if err != nil {
		return nil, fmt.Errorf("invalid time zone: %s", zoneName)
	}
	return loc, nil
}

func convertTime(inputTimeStr, sourceZoneName, targetZoneName string) (time.Time, time.Time, error) {
	sourceLoc, err := getTimeZone(sourceZoneName)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	targetLoc, err := getTimeZone(targetZoneName)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	const layout = "2006-01-02 15:04:05"

	parsedTime, err := time.ParseInLocation(layout, inputTimeStr, sourceLoc)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("could not parse time '%s' with format '%s': %w", inputTimeStr, layout, err)
	}

	convertedTime := parsedTime.In(targetLoc)

	return parsedTime, convertedTime, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter time (e.g., 2023-10-27 15:30:00): ")
	inputTimeStr, _ := reader.ReadString('\n')
	inputTimeStr = strings.TrimSpace(inputTimeStr)

	fmt.Print("Enter source time zone (e.g., America/New_York): ")
	sourceZoneName, _ := reader.ReadString('\n')
	sourceZoneName = strings.TrimSpace(sourceZoneName)

	fmt.Print("Enter target time zone (e.g., Europe/London): ")
	targetZoneName, _ := reader.ReadString('\n')
	targetZoneName = strings.TrimSpace(targetZoneName)

	originalTime, convertedTime, err := convertTime(inputTimeStr, sourceZoneName, targetZoneName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nOriginal Time: %s (%s)\n", originalTime.Format("2006-01-02 15:04:05 MST"), originalTime.Location().String())
	fmt.Printf("Converted Time: %s (%s)\n", convertedTime.Format("2006-01-02 15:04:05 MST"), convertedTime.Location().String())
}