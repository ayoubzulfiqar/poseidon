package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func convertLength() {
	fmt.Println("\n--- Length Conversion ---")
	fmt.Println("Units: m (meters), ft (feet), in (inches), cm (centimeters), km (kilometers), mi (miles)")
	valueStr := getInput("Enter value to convert: ")
	value, err := parseFloat(valueStr)
	if err != nil {
		fmt.Println("Invalid value. Please enter a number.")
		return
	}

	fromUnit := strings.ToLower(getInput("Convert from unit (e.g., m, ft): "))
	toUnit := strings.ToLower(getInput("Convert to unit (e.g., cm, mi): "))

	var valueInMeters float64
	switch fromUnit {
	case "m":
		valueInMeters = value
	case "ft":
		valueInMeters = value * 0.3048
	case "in":
		valueInMeters = value * 0.0254
	case "cm":
		valueInMeters = value * 0.01
	case "km":
		valueInMeters = value * 1000
	case "mi":
		valueInMeters = value * 1609.34
	default:
		fmt.Println("Invalid 'from' unit.")
		return
	}

	var result float64
	var resultUnit string
	switch toUnit {
	case "m":
		result = valueInMeters
		resultUnit = "m"
	case "ft":
		result = valueInMeters / 0.3048
		resultUnit = "ft"
	case "in":
		result = valueInMeters / 0.0254
		resultUnit = "in"
	case "cm":
		result = valueInMeters / 0.01
		resultUnit = "cm"
	case "km":
		result = valueInMeters / 1000
		resultUnit = "km"
	case "mi":
		result = valueInMeters / 1609.34
		resultUnit = "mi"
	default:
		fmt.Println("Invalid 'to' unit.")
		return
	}

	fmt.Printf("%.4f %s is %.4f %s\n", value, fromUnit, result, resultUnit)
}

func convertWeight() {
	fmt.Println("\n--- Weight Conversion ---")
	fmt.Println("Units: kg (kilograms), lb (pounds), g (grams), oz (ounces)")
	valueStr := getInput("Enter value to convert: ")
	value, err := parseFloat(valueStr)
	if err != nil {
		fmt.Println("Invalid value. Please enter a number.")
		return
	}

	fromUnit := strings.ToLower(getInput("Convert from unit (e.g., kg, lb): "))
	toUnit := strings.ToLower(getInput("Convert to unit (e.g., g, oz): "))

	var valueInKg float64
	switch fromUnit {
	case "kg":
		valueInKg = value
	case "lb":
		valueInKg = value * 0.453592
	case "g":
		valueInKg = value * 0.001
	case "oz":
		valueInKg = value * 0.0283495
	default:
		fmt.Println("Invalid 'from' unit.")
		return
	}

	var result float64
	var resultUnit string
	switch toUnit {
	case "kg":
		result = valueInKg
		resultUnit = "kg"
	case "lb":
		result = valueInKg / 0.453592
		resultUnit = "lb"
	case "g":
		result = valueInKg / 0.001
		resultUnit = "g"
	case "oz":
		result = valueInKg / 0.0283495
		resultUnit = "oz"
	default:
		fmt.Println("Invalid 'to' unit.")
		return
	}

	fmt.Printf("%.4f %s is %.4f %s\n", value, fromUnit, result, resultUnit)
}

func convertTemperature() {
	fmt.Println("\n--- Temperature Conversion ---")
	fmt.Println("Units: C (Celsius), F (Fahrenheit), K (Kelvin)")
	valueStr := getInput("Enter value to convert: ")
	value, err := parseFloat(valueStr)
	if err != nil {
		fmt.Println("Invalid value. Please enter a number.")
		return
	}

	fromUnit := strings.ToLower(getInput("Convert from unit (e.g., C, F): "))
	toUnit := strings.ToLower(getInput("Convert to unit (e.g., F, K): "))

	var valueInC float64
	switch fromUnit {
	case "c":
		valueInC = value
	case "f":
		valueInC = (value - 32) * 5 / 9
	case "k":
		valueInC = value - 273.15
	default:
		fmt.Println("Invalid 'from' unit.")
		return
	}

	var result float64
	var resultUnit string
	switch toUnit {
	case "c":
		result = valueInC
		resultUnit = "C"
	case "f":
		result = (valueInC * 9 / 5) + 32
		resultUnit = "F"
	case "k":
		result = valueInC + 273.15
		resultUnit = "K"
	default:
		fmt.Println("Invalid 'to' unit.")
		return
	}

	fmt.Printf("%.2f %s is %.2f %s\n", value, strings.ToUpper(fromUnit), result, strings.ToUpper(resultUnit))
}

func main() {
	for {
		fmt.Println("\n--- CLI Unit Converter ---")
		fmt.Println("1. Length Conversion")
		fmt.Println("2. Weight Conversion")
		fmt.Println("3. Temperature Conversion")
		fmt.Println("4. Exit")

		choice := getInput("Enter your choice: ")

		switch choice {
		case "1":
			convertLength()
		case "2":
			convertWeight()
		case "3":
			convertTemperature()
		case "4":
			fmt.Println("Exiting converter. Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please select a valid option (1-4).")
		}
	}
}

// Additional implementation at 2025-06-20 01:39:29
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var lengthToBase = map[string]float64{
	"m":  1.0,        // meters
	"km": 1000.0,     // kilometers
	"cm": 0.01,       // centimeters
	"mm": 0.001,      // millimeters
	"mi": 1609.34,    // miles
	"yd": 0.9144,     // yards
	"ft": 0.3048,     // feet
	"in": 0.0254,     // inches
}

var weightToBase = map[string]float64{
	"kg": 1.0,         // kilograms
	"g":  0.001,      // grams
	"mg": 0.000001,   // milligrams
	"t":  1000.0,     // metric tons
	"lb": 0.453592,   // pounds
	"oz": 0.0283495,  // ounces
	"st": 6.35029,    // stones
}

func isLengthUnit(unit string) bool {
	_, ok := lengthToBase[unit]
	return ok
}

func isWeightUnit(unit string) bool {
	_, ok := weightToBase[unit]
	return ok
}

func isTemperatureUnit(unit string) bool {
	tempUnits := map[string]bool{
		"c": true, "f": true, "k": true,
	}
	_, ok := tempUnits[strings.ToLower(unit)]
	return ok
}

func convertLength(value float64, fromUnit, toUnit string) (float64, error) {
	fromFactor, ok := lengthToBase[fromUnit]
	if !ok {
		return 0, fmt.Errorf("unsupported length unit: %s", fromUnit)
	}
	toFactor, ok := lengthToBase[toUnit]
	if !ok {
		return 0, fmt.Errorf("unsupported length unit: %s", toUnit)
	}

	meters := value * fromFactor
	result := meters / toFactor
	return result, nil
}

func convertWeight(value float64, fromUnit, toUnit string) (float64, error) {
	fromFactor, ok := weightToBase[fromUnit]
	if !ok {
		return 0, fmt.Errorf("unsupported weight unit: %s", fromUnit)
	}
	toFactor, ok := weightToBase[toUnit]
	if !ok {
		return 0, fmt.Errorf("unsupported weight unit: %s", toUnit)
	}

	kilograms := value * fromFactor
	result := kilograms / toFactor
	return result, nil
}

func convertTemperature(value float64, fromUnit, toUnit string) (float64, error) {
	fromUnit = strings.ToUpper(fromUnit)
	toUnit = strings.ToUpper(toUnit)

	var celsius float64

	switch fromUnit {
	case "C":
		celsius = value
	case "F":
		celsius = (value - 32) * 5 / 9
	case "K":
		celsius = value - 273.15
	default:
		return 0, fmt.Errorf("unsupported temperature unit: %s", fromUnit)
	}

	switch toUnit {
	case "C":
		return celsius, nil
	case "F":
		return (celsius * 9 / 5) + 32, nil
	case "K":
		return celsius + 273.15, nil
	default:
		return 0, fmt.Errorf("unsupported temperature unit: %s", toUnit)
	}
}

func main() {
	args := os.Args[1:]

	if len(args) != 4 || strings.ToLower(args[2]) != "to" {
		fmt.Println("Usage: converter <value> <from_unit> to <to_unit>")
		fmt.Println("Supported Length Units: m, km, cm, mm, mi, yd, ft, in")
		fmt.Println("Supported Weight Units: kg, g, mg, t (metric ton), lb, oz, st")
		fmt.Println("Supported Temperature Units: C, F, K")
		fmt.Println("Example: converter 100 m to ft")
		fmt.Println("Example: converter 25 C to F")
		fmt.Println("Example: converter 10 kg to lb")
		os.Exit(1)
	}

	valueStr := args[0]
	fromUnitRaw := args[1]
	toUnitRaw := args[3]

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		fmt.Printf("Error: Invalid value '%s'. Must be a number.\n", valueStr)
		os.Exit(1)
	}

	fromUnitLower := strings.ToLower(fromUnitRaw)
	toUnitLower := strings.ToLower(toUnitRaw)

	var result float64
	var conversionErr error

	if isLengthUnit(fromUnitLower) && isLengthUnit(toUnitLower) {
		result, conversionErr = convertLength(value, fromUnitLower, toUnitLower)
	} else if isWeightUnit(fromUnitLower) && isWeightUnit(toUnitLower) {
		result, conversionErr = convertWeight(value, fromUnitLower, toUnitLower)
	} else if isTemperatureUnit(fromUnitLower) && isTemperatureUnit(toUnitLower) {
		result, conversionErr = convertTemperature(value, fromUnitLower, toUnitLower)
	} else {
		conversionErr = fmt.Errorf("unsupported or mixed unit types: %s to %s", fromUnitRaw, toUnitRaw)
	}

	if conversionErr != nil {
		fmt.Printf("Error: %v\n", conversionErr)
		os.Exit(1)
	}

	fmt.Printf("%.4f %s is %.4f %s\n", value, fromUnitRaw, result, toUnitRaw)
}

// Additional implementation at 2025-06-20 01:40:05
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var lengthFactors = map[string]float64{
	"m":          1.0,
	"meter":      1.0,
	"km":         1000.0,
	"kilometer":  1000.0,
	"cm":         0.01,
	"centimeter": 0.01,
	"mm":         0.001,
	"millimeter": 0.001,
	"mi":         1609.34,
	"mile":       1609.34,
	"yd":         0.9144,
	"yard":       0.9144,
	"ft":         0.3048,
	"foot":       0.3048,
	"in":         0.0254,
	"inch":       0.0254,
}

var weightFactors = map[string]float64{
	"kg":        1.0,
	"kilogram":  1.0,
	"g":         0.001,
	"gram":      0.001,
	"mg":        0.000001,
	"milligram": 0.000001,
	"lb":        0.453592,
	"pound":     0.453592,
	"oz":        0.0283495,
	"ounce":     0.0283495,
}

func toCelsius(value float64, unit string) (float64, error) {
	switch strings.ToLower(unit) {
	case "c", "celsius":
		return value, nil
	case "f", "fahrenheit":
		return (value - 32) * 5 / 9, nil
	case "k", "kelvin":
		return value - 273.15, nil
	default:
		return 0, fmt.Errorf("unsupported temperature unit: %s", unit)
	}
}

func fromCelsius(value float64, unit string) (float64, error) {
	switch strings.ToLower(unit) {
	case "c", "celsius":
		return value, nil
	case "f", "fahrenheit":
		return (value * 9 / 5) + 32, nil
	case "k", "kelvin":
		return value + 273.15, nil
	default:
		return 0, fmt.Errorf("unsupported temperature unit: %s", unit)
	}
}

func convertLength(value float64, fromUnit, toUnit string) (float64, error) {
	fromUnit = strings.ToLower(fromUnit)
	toUnit = strings.ToLower(toUnit)

	fromFactor, ok := lengthFactors[fromUnit]
	if !ok {
		return 0, fmt.Errorf("unsupported length unit: %s", fromUnit)
	}
	toFactor, ok := lengthFactors[toUnit]
	if !ok {
		return 0, fmt.Errorf("unsupported length unit: %s", toUnit)
	}

	meters := value * fromFactor
	result := meters / toFactor
	return result, nil
}

func convertWeight(value float64, fromUnit, toUnit string) (float64, error) {
	fromUnit = strings.ToLower(fromUnit)
	toUnit = strings.ToLower(toUnit)

	fromFactor, ok := weightFactors[fromUnit]
	if !ok {
		return 0, fmt.Errorf("unsupported weight unit: %s", fromUnit)
	}
	toFactor, ok := weightFactors[toUnit]
	if !ok {
		return 0, fmt.Errorf("unsupported weight unit: %s", toUnit)
	}

	kilograms := value * fromFactor
	result := kilograms / toFactor
	return result, nil
}

func convertTemperature(value float64, fromUnit, toUnit string) (float64, error) {
	celsiusValue, err := toCelsius(value, fromUnit)
	if err != nil {
		return 0, err
	}

	result, err := fromCelsius(celsiusValue, toUnit)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func printUsage() {
	fmt.Println("Usage: converter <value> <from_unit> <to_unit>")
	fmt.Println("Example: converter 10 km miles")
	fmt.Println("Supported Length Units: m, km, cm, mm, mi, yd, ft, in (and full names)")
	fmt.Println("Supported Weight Units: kg, g, mg, lb, oz (and full names)")
	fmt.Println("Supported Temperature Units: c, f, k (and full names)")
}

func main() {
	if len(os.Args) != 4 {
		printUsage()
		os.Exit(1)
	}

	valueStr := os.Args[1]
	fromUnit := os.Args[2]
	toUnit := os.Args[3]

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		fmt.Printf("Error: Invalid value '%s'. Please provide a number.\n", valueStr)
		os.Exit(1)
	}

	var result float64
	var conversionErr error

	fromUnitLower := strings.ToLower(fromUnit)
	toUnitLower := strings.ToLower(toUnit)

	_, fromIsLength := lengthFactors[fromUnitLower]
	_, toIsLength := lengthFactors[toUnitLower]

	_, fromIsWeight := weightFactors[fromUnitLower]
	_, toIsWeight := weightFactors[toUnitLower]

	isTempUnit := func(unit string) bool {
		u := strings.ToLower(unit)
		return u == "c" || u == "celsius" ||
			u == "f" || u == "fahrenheit" ||
			u == "k" || u == "kelvin"
	}
	fromIsTemp := isTempUnit(fromUnit)
	toIsTemp := isTempUnit(toUnit)

	if fromIsLength && toIsLength {
		result, conversionErr = convertLength(value, fromUnit, toUnit)
	} else if fromIsWeight && toIsWeight {
		result, conversionErr = convertWeight(value, fromUnit, toUnit)
	} else if fromIsTemp && toIsTemp {
		result, conversionErr = convertTemperature(value, fromUnit, toUnit)
	} else {
		conversionErr = fmt.Errorf("incompatible units or unsupported unit type: %s to %s", fromUnit, toUnit)
	}

	if conversionErr != nil {
		fmt.Printf("Error: %v\n", conversionErr)
		printUsage()
		os.Exit(1)
	}

	fmt.Printf("%.4f %s is %.4f %s\n", value, fromUnit, result, toUnit)
}