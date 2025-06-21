package main

import (
	"fmt"
	"sort"
	"time"
)

type DataPoint struct {
	Timestamp time.Time
	Value     float64
}

type TimeSeries struct {
	Points []DataPoint
}

func NewTimeSeries() *TimeSeries {
	return &TimeSeries{
		Points: make([]DataPoint, 0),
	}
}

func (ts *TimeSeries) AddPoint(timestamp time.Time, value float64) {
	newPoint := DataPoint{Timestamp: timestamp, Value: value}

	i := sort.Search(len(ts.Points), func(j int) bool {
		return ts.Points[j].Timestamp.After(newPoint.Timestamp)
	})

	ts.Points = append(ts.Points, DataPoint{})
	copy(ts.Points[i+1:], ts.Points[i:])
	ts.Points[i] = newPoint
}

func (ts *TimeSeries) Mean() float64 {
	if len(ts.Points) == 0 {
		return 0.0
	}
	sum := 0.0
	for _, p := range ts.Points {
		sum += p.Value
	}
	return sum / float64(len(ts.Points))
}

func (ts *TimeSeries) MinMax() (min float64, max float64) {
	if len(ts.Points) == 0 {
		return 0.0, 0.0
	}
	min = ts.Points[0].Value
	max = ts.Points[0].Value
	for _, p := range ts.Points {
		if p.Value < min {
			min = p.Value
		}
		if p.Value > max {
			max = p.Value
		}
	}
	return min, max
}

func (ts *TimeSeries) SimpleTrend() string {
	if len(ts.Points) < 2 {
		return "Not enough data for trend analysis"
	}
	firstValue := ts.Points[0].Value
	lastValue := ts.Points[len(ts.Points)-1].Value

	if lastValue > firstValue {
		return "Upward Trend"
	} else if lastValue < firstValue {
		return "Downward Trend"
	} else {
		return "Stable Trend"
	}
}

func main() {
	ts := NewTimeSeries()

	ts.AddPoint(time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC), 10.5)
	ts.AddPoint(time.Date(2023, time.January, 3, 0, 0, 0, 0, time.UTC), 12.1)
	ts.AddPoint(time.Date(2023, time.January, 2, 0, 0, 0, 0, time.UTC), 11.2)
	ts.AddPoint(time.Date(2023, time.January, 4, 0, 0, 0, 0, time.UTC), 13.0)
	ts.AddPoint(time.Date(2023, time.January, 5, 0, 0, 0, 0, time.UTC), 12.8)
	ts.AddPoint(time.Date(2023, time.January, 6, 0, 0, 0, 0, time.UTC), 13.5)

	fmt.Println("Time Series Data:")
	for _, p := range ts.Points {
		fmt.Printf("  %s: %.2f\n", p.Timestamp.Format("2006-01-02"), p.Value)
	}

	fmt.Println("\n--- Analysis ---")
	mean := ts.Mean()
	fmt.Printf("Mean Value: %.2f\n", mean)

	minVal, maxVal := ts.MinMax()
	fmt.Printf("Min Value: %.2f, Max Value: %.2f\n", minVal, maxVal)

	trend := ts.SimpleTrend()
	fmt.Printf("Simple Trend: %s\n", trend)

	ts2 := NewTimeSeries()
	ts2.AddPoint(time.Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC), 20.0)
	ts2.AddPoint(time.Date(2023, time.February, 2, 0, 0, 0, 0, time.UTC), 19.5)
	ts2.AddPoint(time.Date(2023, time.February, 3, 0, 0, 0, 0, time.UTC), 18.0)
	fmt.Println("\n--- Second Time Series Analysis ---")
	fmt.Printf("Simple Trend (ts2): %s\n", ts2.SimpleTrend())

	ts3 := NewTimeSeries()
	ts3.AddPoint(time.Date(2023, time.March, 1, 0, 0, 0, 0, time.UTC), 15.0)
	ts3.AddPoint(time.Date(2023, time.March, 2, 0, 0, 0, 0, time.UTC), 15.0)
	fmt.Println("\n--- Third Time Series Analysis ---")
	fmt.Printf("Simple Trend (ts3): %s\n", ts3.SimpleTrend())

	ts4 := NewTimeSeries()
	fmt.Println("\n--- Empty Time Series Analysis ---")
	fmt.Printf("Mean (ts4): %.2f\n", ts4.Mean())
	min4, max4 := ts4.MinMax()
	fmt.Printf("Min/Max (ts4): %.2f / %.2f\n", min4, max4)
	fmt.Printf("Simple Trend (ts4): %s\n", ts4.SimpleTrend())
}

// Additional implementation at 2025-06-21 03:51:22
