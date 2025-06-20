package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type RGB struct {
	R, G, B uint8
}

type HEX string

type HSL struct {
	H, S, L float64
}

func (c RGB) ToHEX() HEX {
	return HEX(fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B))
}

func (c RGB) ToHSL() HSL {
	r := float64(c.R) / 255.0
	g := float64(c.G) / 255.0
	b := float64(c.B) / 255.0

	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))

	h, s, l := 0.0, 0.0, 0.0

	l = (max + min) / 2.0

	if max == min {
		h = 0.0
		s = 0.0
	} else {
		delta := max - min
		if l > 0.5 {
			s = delta / (2.0 - max - min)
		} else {
			s = delta / (max + min)
		}

		switch max {
		case r:
			h = (g - b) / delta
			if g < b {
				h += 6.0
			}
		case g:
			h = (b - r) / delta + 2.0
		case b:
			h = (r - g) / delta + 4.0
		}
		h *= 60.0
	}

	return HSL{H: h, S: s * 100.0, L: l * 100.0}
}

func (h HEX) ToRGB() (RGB, error) {
	hexStr := string(h)
	if !strings.HasPrefix(hexStr, "#") || len(hexStr) != 7 {
		return RGB{}, fmt.Errorf("invalid HEX format: %s", hexStr)
	}

	r, err := strconv.ParseUint(hexStr[1:3], 16, 8)
	if err != nil {
		return RGB{}, fmt.Errorf("invalid red component: %w", err)
	}
	g, err := strconv.ParseUint(hexStr[3:5], 16, 8)
	if err != nil {
		return RGB{}, fmt.Errorf("invalid green component: %w", err)
	}
	b, err := strconv.ParseUint(hexStr[5:7], 16, 8)
	if err != nil {
		return RGB{}, fmt.Errorf("invalid blue component: %w", err)
	}

	return RGB{R: uint8(r), G: uint8(g), B: uint8(b)}, nil
}

func (h HEX) ToHSL() (HSL, error) {
	rgb, err := h.ToRGB()
	if err != nil {
		return HSL{}, err
	}
	return rgb.ToHSL(), nil
}

func (c HSL) ToRGB() RGB {
	h := c.H
	s := c.S / 100.0
	l := c.L / 100.0

	var r, g, b float64

	if s == 0 {
		r, g, b = l, l, l
	} else {
		hue2rgb := func(p, q, t float64) float64 {
			if t < 0 {
				t += 1
			}
			if t > 1 {
				t -= 1
			}
			if t < 1.0/6.0 {
				return p + (q-p)*6.0*t
			}
			if t < 1.0/2.0 {
				return q
			}
			if t < 2.0/3.0 {
				return p + (q-p)*(2.0/3.0-t)*6.0
			}
			return p
		}

		var q float64
		if l < 0.5 {
			q = l * (1.0 + s)
		} else {
			q = l + s - l*s
		}
		p := 2.0*l - q

		r = hue2rgb(p, q, h/360.0 + 1.0/3.0)
		g = hue2rgb(p, q, h/360.0)
		b = hue2rgb(p, q, h/360.0 - 1.0/3.0)
	}

	return RGB{
		R: uint8(math.Round(r * 255.0)),
		G: uint8(math.Round(g * 255.0)),
		B: uint8(math.Round(b * 255.0)),
	}
}

func (c HSL) ToHEX() HEX {
	rgb := c.ToRGB()
	return rgb.ToHEX()
}

func main() {
	rgbColor := RGB{R: 255, G: 128, B: 0}
	fmt.Printf("Original RGB: %+v\n", rgbColor)
	hexFromRGB := rgbColor.ToHEX()
	fmt.Printf("RGB to HEX: %s\n", hexFromRGB)
	hslFromRGB := rgbColor.ToHSL()
	fmt.Printf("RGB to HSL: %+v\n", hslFromRGB)

	fmt.Println("---")

	hexColor := HEX("#FF8000")
	fmt.Printf("Original HEX: %s\n", hexColor)
	rgbFromHEX, err := hexColor.ToRGB()
	if err != nil {
		fmt.Printf("Error converting HEX to RGB: %v\n", err)
	} else {
		fmt.Printf("HEX to RGB: %+v\n", rgbFromHEX)
	}
	hslFromHEX, err := hexColor.ToHSL()
	if err != nil {
		fmt.Printf("Error converting HEX to HSL: %v\n", err)
	} else {
		fmt.Printf("HEX to HSL: %+v\n", hslFromHEX)
	}

	fmt.Println("---")

	hslColor := HSL{H: 30, S: 100, L: 50}
	fmt.Printf("Original HSL: %+v\n", hslColor)
	rgbFromHSL := hslColor.ToRGB()
	fmt.Printf("HSL to RGB: %+v\n", rgbFromHSL)
	hexFromHSL := hslColor.ToHEX()
	fmt.Printf("HSL to HEX: %s\n", hexFromHSL)

	fmt.Println("---")

	rgbGray := RGB{R: 128, G: 128, B: 128}
	fmt.Printf("Original RGB (Gray): %+v\n", rgbGray)
	hslGray := rgbGray.ToHSL()
	fmt.Printf("RGB to HSL (Gray): %+v\n", hslGray)

	hslGrayTest := HSL{H: 0, S: 0, L: 50}
	rgbGrayTest := hslGrayTest.ToRGB()
	fmt.Printf("HSL to RGB (Gray): %+v\n", rgbGrayTest)

	rgbRed := RGB{R: 255, G: 0, B: 0}
	fmt.Printf("Original RGB (Red): %+v\n", rgbRed)
	hslRed := rgbRed.ToHSL()
	fmt.Printf("RGB to HSL (Red): %+v\n", hslRed)

	hslRedTest := HSL{H: 0, S: 100, L: 50}
	rgbRedTest := hslRedTest.ToRGB()
	fmt.Printf("HSL to RGB (Red): %+v\n", rgbRedTest)
}

// Additional implementation at 2025-06-20 01:22:10
package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Color struct {
	R, G, B uint8
	H, S, L float64
	Hex     string
}

func NewColorFromRGB(r, g, b uint8) Color {
	c := Color{R: r, G: g, B: b}
	c.toHex()
	c.toHSL()
	return c
}

func NewColorFromHex(hex string) (Color, error) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return Color{}, fmt.Errorf("invalid hex string length: %s", hex)
	}

	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return Color{}, fmt.Errorf("invalid red hex component: %w", err)
	}
	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return Color{}, fmt.Errorf("invalid green hex component: %w", err)
	}
	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return Color{}, fmt.Errorf("invalid blue hex component: %w", err)
	}

	c := Color{R: uint8(r), G: uint8(g), B: uint8(b), Hex: "#" + hex}
	c.toHSL()
	return c, nil
}

func NewColorFromHSL(h, s, l float64) Color {
	c := Color{H: h, S: s, L: l}
	c.toRGB()
	c.toHex()
	return c
}

func (c *Color) toHex() {
	c.Hex = fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}

func (c *Color) toRGB() {
	h := c.H
	s := c.S
	l := c.L

	if s == 0 {
		c.R = uint8(l * 255)
		c.G = uint8(l * 255)
		c.B = uint8(l * 255)
		return
	}

	var hue2rgb func(p, q, t float64) float64
	hue2rgb = func(p, q, t float64) float64 {
		if t < 0 {
			t += 1
		}
		if t > 1 {
			t -= 1
		}
		if t < 1.0/6.0 {
			return p + (q-p)*6*t
		}
		if t < 1.0/2.0 {
			return q
		}
		if t < 2.0/3.0 {
			return p + (q-p)*(2.0/3.0-t)*6
		}
		return p
	}

	q := 0.0
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q

	h /= 360.0

	c.R = uint8(math.Round(hue2rgb(p, q, h+1.0/3.0) * 255))
	c.G = uint8(math.Round(hue2rgb(p, q, h) * 255))
	c.B = uint8(math.Round(hue2rgb(p, q, h-1.0/3.0) * 255))
}

func (c *Color) toHSL() {
	r := float64(c.R) / 255.0
	g := float64(c.G) / 255.0
	b := float64(c.B) / 255.0

	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))

	c.L = (max + min) / 2.0

	if max == min {
		c.H = 0
		c.S = 0
	} else {
		d := max - min
		if c.L > 0.5 {
			c.S = d / (2 - max - min)
		} else {
			c.S = d / (max + min)
		}

		switch max {
		case r:
			c.H = (g - b) / d
			if g < b {
				c.H += 6
			}
		case g:
			c.H = (b-r)/d + 2
		case b:
			c.H = (r-g)/d + 4
		}
		c.H *= 60
	}
}

func (c Color) GetComplementaryHSL() (float64, float64, float64) {
	compH := math.Mod(c.H+180, 360)
	return compH, c.S, c.L
}

func (c Color) GetGrayscaleRGB() (uint8, uint8, uint8) {
	gray := uint8(0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B))
	return gray, gray, gray
}

func (c Color) String() string {
	return fmt.Sprintf("RGB(%d, %d, %d) HSL(%.0f, %.2f, %.2f) Hex(%s)",
		c.R, c.G, c.B, c.H, c.S, c.L, c.Hex)
}

func main() {
	color1 := NewColorFromRGB(255, 0, 0)
	fmt.Println("Color 1 (Red):", color1)
	compH, compS, compL := color1.GetComplementaryHSL()
	fmt.Printf("  Complementary HSL: HSL(%.0f, %.2f, %.2f)\n", compH, compS, compL)
	grayR, grayG, grayB := color1.GetGrayscaleRGB()
	fmt.Printf("  Grayscale RGB: RGB(%d, %d, %d)\n", grayR, grayG, grayB)

	color2, err := NewColorFromHex("#00FF00")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Color 2 (Green):", color2)
		compH, compS, compL = color2.GetComplementaryHSL()
		fmt.Printf("  Complementary HSL: HSL(%.0f, %.2f, %.2f)\n", compH, compS, compL)
		grayR, grayG, grayB = color2.GetGrayscaleRGB()
		fmt.Printf("  Grayscale RGB: RGB(%d, %d, %d)\n", grayR, grayG, grayB)
	}

	color3 := NewColorFromHSL(240, 1.0, 0.5)
	fmt.Println("Color 3 (Blue):", color3)
	compH, compS, compL = color3.GetComplementaryHSL()
	fmt.Printf("  Complementary HSL: HSL(%.0f, %.2f, %.2f)\n", compH, compS, compL)
	grayR, grayG, grayB = color3.GetGrayscaleRGB()
	fmt.Printf("  Grayscale RGB: RGB(%d, %d, %d)\n", grayR, grayG, grayB)

	_, err = NewColorFromHex("invalid")
	if err != nil {
		fmt.Println("Error for invalid hex:", err)
	}
}

// Additional implementation at 2025-06-20 01:22:54
package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// RGB represents a color in RGB format.
type RGB struct {
	R, G, B uint8
}

// HSL represents a color in HSL format.
// H is in [0, 360), S and L are in [0, 1].
type HSL struct {
	H, S, L float64
}

// HEX represents a color in HEX format.
type HEX string

// String returns the RGB color as a string.
func (c RGB) String() string {
	return fmt.Sprintf("rgb(%d, %d, %d)", c.R, c.G, c.B)
}

// String returns the HSL color as a string.
func (c HSL) String() string {
	return fmt.Sprintf("hsl(%.0f, %.0f%%, %.0f%%)", c.H, c.S*100, c.L*100)
}

// String returns the HEX color as a string.
func (h HEX) String() string {
	return string(h)
}

// ToHEX converts an RGB color to its HEX representation.
func (c RGB) ToHEX() HEX {
	return HEX(fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B))
}

// ToHSL converts an RGB color to its HSL representation.
func (c RGB) ToHSL() HSL {
	r := float64(c.R) / 255
	g := float64(c.G) / 255
	b := float64(c.B) / 255

	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))

	h, s, l := 0.0, 0.0, (max+min)/2

	if max == min {
		h = 0 // achromatic
	} else {
		d := max - min
		s = d / (1 - math.Abs(2*l-1))
		switch max {
		case r:
			h = math.Mod((g-b)/d+6, 6) * 60
		case g:
			h = ((b-r)/d + 2) * 60
		case b:
			h = ((r-g)/d + 4) * 60
		}
	}
	return HSL{H: h, S: s, L: l}
}

// ToRGB converts a HEX color to its RGB representation.
func (h HEX) ToRGB() (RGB, error) {
	hexStr := strings.TrimPrefix(string(h), "#")
	if len(hexStr) != 6 && len(hexStr) != 3 {
		return RGB{}, fmt.Errorf("invalid HEX format: %s", h)
	}

	if len(hexStr) == 3 { // Expand 3-digit to 6-digit
		hexStr = string(hexStr[0]) + string(hexStr[0]) +
			string(hexStr[1]) + string(hexStr[1]) +
			string(hexStr[2]) + string(hexStr[2])
	}

	r, err := strconv.ParseUint(hexStr[0:2], 16, 8)
	if err != nil {
		return RGB{}, fmt.Errorf("invalid R component in HEX: %s", hexStr[0:2])
	}
	g, err := strconv.ParseUint(hexStr[2:4], 16, 8)
	if err != nil {
		return RGB{}, fmt.Errorf("invalid G component in HEX: %s", hexStr[2:4])
	}
	b, err := strconv.ParseUint(hexStr[4:6], 16, 8)
	if err != nil {
		return RGB{}, fmt.Errorf("invalid B component in HEX: %s", hexStr[4:6])
	}

	return RGB{R: uint8(r), G: uint8(g), B: uint8(b)}, nil
}

// ToHSL converts a HEX color to its HSL representation.
func (h HEX) ToHSL() (HSL, error) {
	rgb, err := h.ToRGB()
	if err != nil {
		return HSL{}, err
	}
	return rgb.ToHSL(), nil
}

// ToRGB converts an HSL color to its RGB representation.
func (c HSL) ToRGB() RGB {
	var r, g, b float64

	if c.S == 0 {
		r = c.L
		g = c.L
		b = c.L // achromatic
	} else {
		hue2rgb := func(p, q, t float64) float64 {
			if t < 0 {
				t += 1
			}
			if t > 1 {
				t -= 1
			}
			if t < 1/6.0 {
				return p + (q-p)*6*t
			}
			if t < 1/2.0 {
				return q
			}
			if t < 2/3.0 {
				return p + (q-p)*(2/3.0-t)*6
			}
			return p
		}

		q := 0.0
		if c.L < 0.5 {
			q = c.L * (1 + c.S)
		} else {
			q = c.L + c.S - c.L*c.S
		}
		p := 2*c.L - q

		r = hue2rgb(p, q, c.H/360+1/3.0)
		g = hue2rgb(p, q, c.H/360)
		b = hue2rgb(p, q, c.H/360-1/3.0)
	}

	return RGB{R: uint8(r * 255), G: uint8(g * 255), B: uint8(b * 255)}
}

// ToHEX converts an HSL color to its HEX representation.
func (c HSL) ToHEX() HEX {
	rgb := c.ToRGB()
	return rgb.ToHEX()
}

// ParseRGB parses an RGB string (e.g., "rgb(255, 0, 128)") into an RGB struct.
func ParseRGB(s string) (RGB, error) {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "rgb(") || !strings.HasSuffix(s, ")") {
		return RGB{}, fmt.Errorf("invalid RGB format: %s", s)
	}
	parts := strings.Split(s[4:len(s)-1], ",")
	if len(parts) != 3 {
		return RGB{}, fmt.Errorf("invalid RGB components count: %s", s)
	}

	r, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil || r < 0 || r > 255 {
		return RGB{}, fmt.Errorf("invalid R component: %s", parts[0])
	}
	g, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil || g < 0 || g > 255 {
		return RGB{}, fmt.Errorf("invalid G component: %s", parts[1])
	}
	b, err := strconv.Atoi(strings.TrimSpace(parts[2]))
	if err != nil || b < 0 || b > 255 {
		return RGB{}, fmt.Errorf("invalid B component: %s", parts[2])
	}
	return RGB{R: uint8(r), G: uint8(g), B: uint8(b)}, nil
}

// ParseHEX parses a HEX string (e.g., "#FF0080" or "FF0080") into a HEX type.
func ParseHEX(s string) (HEX, error) {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "#") {
		s = "#" + s // Assume # if missing
	}
	if len(s) != 7 && len(s) != 4 { // #RRGGBB or #RGB
		return "", fmt.Errorf("invalid HEX format length: %s", s)
	}
	// Validate hex characters
	for i := 1; i < len(s); i++ {
		c := s[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return "", fmt.Errorf("invalid HEX character '%c' in %s", c, s)
		}
	}
	return HEX(strings.ToUpper(s)), nil
}

// ParseHSL parses an HSL string (e.g., "hsl(120, 100%, 50%)") into an HSL struct.
func ParseHSL(s string) (HSL, error) {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "hsl(") || !strings.HasSuffix(s, ")") {
		return HSL{}, fmt.Errorf("invalid HSL format: %s", s)
	}
	parts := strings.Split(s[4:len(s)-1], ",")
	if len(parts) != 3 {
		return HSL{}, fmt.Errorf("invalid HSL components count: %s", s)
	}

	hStr := strings.TrimSpace(parts[0])
	sStr := strings.TrimSuffix(strings.TrimSpace(parts[1]), "%")
	lStr := strings.TrimSuffix(strings.TrimSpace(parts[2]), "%")

	h, err := strconv.ParseFloat(hStr, 64)
	if err != nil || h < 0 || h >= 360 { // H is [0, 360)
		return HSL{}, fmt.Errorf("invalid H component: %s", parts[0])
	}
	sVal, err := strconv.ParseFloat(sStr, 64)
	if err != nil || sVal < 0 || sVal > 100 {
		return HSL{}, fmt.Errorf("invalid S component: %s", parts[1])
	}
	lVal, err := strconv.

// Additional implementation at 2025-06-20 01:24:03
package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// RGB represents a color in RGB format.
type RGB struct {
	R, G, B uint8
}

// String returns the CSS-like string representation of the RGB color.
func (c RGB) String() string {
	return fmt.Sprintf("rgb(%d, %d, %d)", c.R, c.G, c.B)
}

// Hex represents a color in hexadecimal format.
type Hex string

// String returns the string representation of the Hex color.
func (h Hex) String() string {
	return string(h)
}

// HSL represents a color in HSL format.
// H is in degrees [0, 360], S and L are in [0, 1].
type HSL struct {
	H, S, L float64
}

// String returns the CSS-like string representation of the HSL color.
func (c HSL) String() string {
	return fmt.Sprintf("hsl(%.0f, %.0f%%, %.0f%%)", c.H, c.S*100, c.L*100)
}

// RGBToHex converts an RGB color to its hexadecimal representation.
func RGBToHex(c RGB) Hex {
	return Hex(fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B))
}

// HexToRGB converts a hexadecimal color string to an RGB color.
// It supports formats like "#RRGGBB", "RRGGBB", "#RGB", "RGB".
func HexToRGB(hexStr string) (RGB, error) {
	hexStr = strings.TrimPrefix(hexStr, "#")
	if len(hexStr) == 3 {
		hexStr = string([]byte{hexStr[0], hexStr[0], hexStr[1], hexStr[1], hexStr[2], hexStr[2]})
	}
	if len(hexStr) != 6 {
		return RGB{}, errors.New("invalid hex string length, expected 3 or 6 characters (excluding '#')")
	}

	r, err := strconv.ParseUint(hexStr[0:2], 16, 8)
	if err != nil {
		return RGB{}, fmt.Errorf("invalid red component in hex: %w", err)
	}
	g, err := strconv.ParseUint(hexStr[2:4], 16, 8)
	if err != nil {
		return RGB{}, fmt.Errorf("invalid green component in hex: %w", err)
	}
	b, err := strconv.ParseUint(hexStr[4:6], 16, 8)
	if err != nil {
		return RGB{}, fmt.Errorf("invalid blue component in hex: %w", err)
	}

	return RGB{uint8(r), uint8(g), uint8(b)}, nil
}

// RGBToHSL converts an RGB color to its HSL representation.
// Algorithm based on https://en.wikipedia.org/wiki/HSL_and_HSV#From_RGB
func RGBToHSL(c RGB) HSL {
	r := float64(c.R) / 255.0
	g := float64(c.G) / 255.0
	b := float64(c.B) / 255.0

	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))

	h, s, l := 0.0, 0.0, (max+min)/2.0

	if max == min {
		h = 0.0 // achromatic
		s = 0.0
	} else {
		d := max - min
		if l > 0.5 {
			s = d / (2.0 - max - min)
		} else {
			s = d / (max + min)
		}

		switch max {
		case r:
			h = (g - b) / d
			if g < b {
				h += 6.0
			}
		case g:
			h = (b-r)/d + 2.0
		case b:
			h = (r-g)/d + 4.0
		}
		h *= 60.0
	}
	return HSL{H: h, S: s, L: l}
}

// HSLToRGB converts an HSL color to its RGB representation.
// Algorithm based on https://en.wikipedia.org/wiki/HSL_and_HSV#To_RGB
func HSLToRGB(c HSL) RGB {
	var r, g, b float64

	if c.S == 0 {
		r, g, b = c.L, c.L, c.L // achromatic
	} else {
		hue2rgb := func(p, q, t float64) float64 {
			if t < 0 {
				t += 1
			}
			if t > 1 {
				t -= 1
			}
			if t < 1.0/6.0 {
				return p + (q-p)*6*t
			}
			if t < 1.0/2.0 {
				return q
			}
			if t < 2.0/3.0 {
				return p + (q-p)*(2.0/3.0-t)*6
			}
			return p
		}

		q := c.L
		if c.L < 0.5 {
			q = c.L * (1 + c.S)
		} else {
			q = c.L + c.S - c.L*c.S
		}
		p := 2*c.L - q
		h := c.H / 360.0

		r = hue2rgb(p, q, h+1.0/3.0)
		g = hue2rgb(p, q, h)
		b = hue2rgb(p, q, h-1.0/3.0)
	}

	return RGB{uint8(r * 255), uint8(g * 255), uint8(b * 255)}
}

// HexToHSL converts a hexadecimal color string to an HSL color.
func HexToHSL(hexStr string) (HSL, error) {
	rgb, err := HexToRGB(hexStr)
	if err != nil {
		return HSL{}, err
	}
	return RGBToHSL(rgb), nil
}

// HSLToHex converts an HSL color to its hexadecimal representation.
func HSLToHex(c HSL) Hex {
	rgb := HSLToRGB(c)
	return RGBToHex(rgb)
}

// ParseRGB parses a CSS-like RGB string (e.g., "rgb(255, 0, 128)") into an RGB struct.
func ParseRGB(rgbStr string) (RGB, error) {
	rgbStr = strings.TrimSpace(rgbStr)
	if !strings.HasPrefix(rgbStr, "rgb(") || !strings.HasSuffix(rgbStr, ")") {
		return RGB{}, errors.New("invalid RGB string format, expected 'rgb(r,g,b)'")
	}
	content := strings.TrimPrefix(rgbStr, "rgb(")
	content = strings.TrimSuffix(content, ")")
	parts := strings.Split(content, ",")
	if len(parts) != 3 {
		return RGB{}, errors.New("invalid RGB string format, expected 3 components")
	}

	var r, g, b uint64
	var err error

	r, err = strconv.ParseUint(strings.TrimSpace(parts[0]), 10, 8)
	if err != nil {
		return RGB{}, fmt.Errorf("invalid red component: %w", err)
	}
	g, err = strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 8)
	if err != nil {
		return RGB{}, fmt.Errorf("invalid green component: %w", err)
	}
	b, err = strconv.ParseUint(strings.TrimSpace(parts[2]), 10, 8)
	if err != nil {
		return RGB{}, fmt.Errorf("invalid blue component: %w", err)
	}

	return RGB{uint8(r), uint8(g), uint8(b)}, nil
}

// ParseHSL parses a CSS-like HSL string (e.g., "hsl(120, 100%, 50%)") into an HSL struct.
func ParseHSL(hslStr string) (HSL, error) {
	hslStr = strings.TrimSpace(hslStr)
	if !strings.HasPrefix(hslStr, "hsl(") || !strings.HasSuffix(hslStr, ")") {
		return HSL{}, errors.New("invalid HSL string format, expected 'hsl(h,s%,l%)'")
	}
	content := strings.TrimPrefix(hslStr, "hsl(")
	content = strings.TrimSuffix(content, ")")
	parts := strings.Split(content, ",")
	if len(parts) != 3 {
		return HSL{}, errors.New("invalid HSL string format, expected 3 components")
	}

	var h, s, l float64
	var err error

	h, err = strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return HSL{}, fmt.Errorf("invalid hue component: %w", err)
	}

	sStr := strings.TrimSpace(parts[1])
	if !strings.HasSuffix(sStr, "%") {
		return HSL{}, errors.New("invalid saturation component: missing '%'")
	}
	s, err = strconv.ParseFloat(strings.TrimSuffix(sStr, "%"), 64)
	if err != nil {
		return HSL{}, fmt.Errorf("invalid saturation component: %w", err)
	}
	s /= 100.0 // Convert percentage to [0, 1]

	lStr := strings.TrimSpace(parts[2])
	if !strings.HasSuffix(lStr, "%") {
		return HSL{}, errors.New("invalid lightness component: missing '%'")
	}
	l, err = strconv.ParseFloat(strings.TrimSuffix(lStr, "%"), 64)
	if err != nil {
		return HSL{}, fmt.Errorf("invalid lightness component: %w", err)
	}
	l /= 100.0 // Convert percentage to [0, 1]

	return HSL{H: h, S: s, L: l}, nil
}

// Invert returns the inverted RGB color.
func (c RGB) Invert() RGB {
	return RGB{R: 255 - c.R, G: 255 - c.G, B: 255 - c.B}
}

// Grayscale converts the RGB color to its grayscale equivalent using the luminosity method.
func (c RGB) Grayscale() RGB {
	gray := uint8