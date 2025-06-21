package main

import (
	"fmt"
	"math"
	"time"
)

const (
	gravity = 9.81 // m/s^2
	dt      = 0.01 // time step in seconds
)

type Projectile struct {
	x  float64
	y  float64
	vx float64
	vy float64
}

func NewProjectile(initialVelocity, angleDegrees float64) *Projectile {
	angleRadians := angleDegrees * math.Pi / 180.0
	return &Projectile{
		x:  0.0,
		y:  0.0,
		vx: initialVelocity * math.Cos(angleRadians),
		vy: initialVelocity * math.Sin(angleRadians),
	}
}

func (p *Projectile) Update(dt float64) bool {
	if p.y <= 0 && p.vy <= 0 && p.vx == 0 {
		p.y = 0
		p.vx = 0
		p.vy = 0
		return false
	}

	p.vy -= gravity * dt

	p.x += p.vx * dt
	p.y += p.vy * dt

	if p.y < 0 {
		p.y = 0
		p.vx = 0
		p.vy = 0
		return false
	}
	return true
}

func main() {
	initialVelocity := 50.0
	angleDegrees := 45.0

	p := NewProjectile(initialVelocity, angleDegrees)

	fmt.Printf("Time (s)\tX (m)\tY (m)\n")
	fmt.Printf("--------\t-----\t-----\n")

	t := 0.0
	fmt.Printf("%.2f\t\t%.2f\t%.2f\n", t, p.x, p.y)

	for {
		t += dt
		inMotion := p.Update(dt)

		fmt.Printf("%.2f\t\t%.2f\t%.2f\n", t, p.x, p.y)

		if !inMotion {
			break
		}

		time.Sleep(1 * time.Millisecond)
	}
}

// Additional implementation at 2025-06-21 00:50:47
package main

import (
	"fmt"
	"math"
)

// Vector represents a 2D vector for position and velocity
type Vector struct {
	X, Y float64
}

// Add returns the sum of two vectors
func (v Vector) Add(other Vector) Vector {
	return Vector{v.X + other.X, v.Y + other.Y}
}

// Scale returns the vector scaled by a scalar
func (v Vector) Scale(s float64) Vector {
	return Vector{v.X * s, v.Y * s}
}

// Magnitude returns the length of the vector
func (v Vector) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Normalize returns a unit vector in the same direction
func (v Vector) Normalize() Vector {
	mag := v.Magnitude()
	if mag == 0 {
		return Vector{}
	}
	return Vector{v.X / mag, v.Y / mag}
}

// Projectile represents an object in the simulation
type Projectile struct {
	Position Vector
	Velocity Vector
	Mass     float64 // kg
	Radius   float64 // m, for air resistance calculation (assumes spherical)
	Color    string  // For identification in output
}

// Constants for the simulation
const (
	G   = 9.81  // Acceleration due to gravity (m/s^2)
	DT  = 0.01  // Time step (s)
	RHO = 1.225 // Air density (kg/m^3) at sea level, 15°C
	CD  = 0.47  // Drag coefficient for a smooth sphere
)

// SimulateStep updates the projectile's state for one time step.
// It incorporates gravity and optionally air resistance.
// It also handles ground collision by stopping vertical movement at Y=0.
func (p *Projectile) SimulateStep(dt float64, applyAirResistance bool) {
	// Calculate gravitational force
	gravityForce := Vector{0, -p.Mass * G}
	netForce := gravityForce

	// Calculate and apply air resistance if enabled
	if applyAirResistance {
		// Drag force = 0.5 * rho * v^2 * Cd * A
		// A = pi * r^2 (cross-sectional area of a sphere)
		area := math.Pi * p.Radius * p.Radius
		speed := p.Velocity.Magnitude()
		dragMagnitude := 0.5 * RHO * speed * speed * CD * area

		// Drag force acts opposite to velocity
		if speed > 0 {
			dragDirection := p.Velocity.Normalize().Scale(-1)
			dragForce := dragDirection.Scale(dragMagnitude)
			netForce = netForce.Add(dragForce)
		}
	}

	// Calculate acceleration (F = ma => a = F/m)
	acceleration := netForce.Scale(1 / p.Mass)

	// Update velocity and position using Euler integration
	p.Velocity = p.Velocity.Add(acceleration.Scale(dt))
	p.Position = p.Position.Add(p.Velocity.Scale(dt))

	// Ground collision detection and response
	if p.Position.Y < 0 {
		p.Position.Y = 0      // Set position to ground level
		p.Velocity.Y = 0      // Stop vertical movement
		p.Velocity.X = 0      // Stop horizontal movement (no bounce implemented)
	}
}

// SimulateTrajectory runs the simulation for a projectile until it hits the ground (Y <= 0).
// It returns a slice of Vector representing the trajectory points.
func SimulateTrajectory(p Projectile, applyAirResistance bool) []Vector {
	trajectory := []Vector{p.Position}
	// Continue simulation as long as the projectile is above ground and moving
	for p.Position.Y >= 0 && p.Velocity.Magnitude() > 0.01 { // Small threshold for stopping
		p.SimulateStep(DT, applyAirResistance)
		trajectory = append(trajectory, p.Position)
		// Safety break to prevent infinite loops for very long or edge-case trajectories
		if len(trajectory) > 500000 { // Max 5000 seconds of simulation at DT=0.01
			break
		}
	}
	return trajectory
}

func main() {
	fmt.Println("--- Go Projectile Motion Simulator ---")

	// Scenario 1: Basic Projectile (no air resistance)
	fmt.Println("\n--- Scenario 1: Projectile without Air Resistance ---")
	p1 := Projectile{
		Position: Vector{0, 100}, // Start at (0, 100) meters
		Velocity: Vector{50, 0},  // Initial velocity 50 m/s horizontally
		Mass:     1.0,            // Mass (kg) - irrelevant without air resistance
		Radius:   0.1,            // Radius (m) - irrelevant without air resistance
		Color:    "Blue",
	}
	trajectory1 := SimulateTrajectory(p1, false) // false for no air resistance
	fmt.Printf("Projectile %s (No Air Resistance) - Initial Pos: (%.2f, %.2f), Vel: (%.2f, %.2f)\n",
		p1.Color, p1.Position.X, p1.Position.Y, p1.Velocity.X, p1.Velocity.Y)
	fmt.Printf("  Total points: %d. Final position: (%.2f, %.2f)\n", len(trajectory1), trajectory1[len(trajectory1)-1].X, trajectory1[len(trajectory1)-1].Y)
	fmt.Print("  Sample Trajectory Points (X, Y): ")
	for i := 0; i < 3 && i < len(trajectory1); i++ {
		fmt.Printf("(%.2f, %.2f) ", trajectory1[i].X, trajectory1[i].Y)
	}
	if len(trajectory1) > 6 {
		fmt.Print("... ")
		for i := len(trajectory1) - 3; i < len(trajectory1); i++ {
			fmt.Printf("(%.2f, %.2f) ", trajectory1[i].X, trajectory1[i].Y)
		}
	}
	fmt.Println()

	// Scenario 2: Projectile with Air Resistance
	fmt.Println("\n--- Scenario 2: Projectile with Air Resistance ---")
	p2 := Projectile{
		Position: Vector{0, 100}, // Start at (0, 100) meters
		Velocity: Vector{50, 0},  // Initial velocity 50 m/s horizontally
		Mass:     1.0,            // Mass (kg)
		Radius:   0.1,            // Radius (m) - matters for air resistance
		Color:    "Red",
	}
	trajectory2 := SimulateTrajectory(p2, true) // true for air resistance
	fmt.Printf("Projectile %s (With Air Resistance) - Initial Pos: (%.2f, %.2f), Vel: (%.2f, %.2f)\n",
		p2.Color, p2.Position.X, p2.Position.Y, p2.Velocity.X, p2.Velocity.Y)
	fmt.Printf("  Total points: %d. Final position: (%.2f, %.2f)\n", len(trajectory2), trajectory2[len(trajectory2)-1].X, trajectory2[len(trajectory2)-1].Y)
	fmt.Print("  Sample Trajectory Points (X, Y): ")
	for i := 0; i < 3 && i < len(trajectory2); i++ {
		fmt.Printf("(%.2f, %.2f) ", trajectory2[i].X, trajectory2[i].Y)
	}
	if len(trajectory2) > 6 {
		fmt.Print("... ")
		for i := len(trajectory2) - 3; i < len(trajectory2); i++ {
			fmt.Printf("(%.2f, %.2f) ", trajectory2[i].X, trajectory2[i].Y)
		}
	}
	fmt.Println()

	// Scenario 3: Multiple Projectiles with varying properties and air resistance settings
	fmt.Println("\n--- Scenario 3: Multiple Projectiles ---")
	projectiles := []Projectile{
		{Position: Vector{0, 50}, Velocity: Vector{30, 30}, Mass: 1.0, Radius: 0.05, Color: "Green"},  // Light, small, with AR
		{Position: Vector{0, 50}, Velocity: Vector{30, 30}, Mass: 10.0, Radius: 0.2, Color: "Yellow"}, // Heavy, large, with AR
		{Position: Vector{0, 50}, Velocity: Vector{30, 30}, Mass: 1.0, Radius: 0.05, Color: "Purple"}, // Same as Green, but NO AR
	}

	for i, p := range projectiles {
		applyAR := true
		if p.Color == "Purple" { // Example: explicitly disable AR for the "Purple" projectile
			applyAR = false
		}
		fmt.Printf("\nSimulating Projectile %d (%s, Mass: %.1fkg, Radius: %.2fm, Air Resistance: %t)\n",
			i+1, p.Color, p.Mass, p.Radius, applyAR)
		traj := SimulateTrajectory(p, applyAR)
		fmt.Printf("  Initial Pos: (%.2f, %.2f), Vel: (%.2f, %.2f)\n", p.Position.X, p.Position.Y, p.Velocity.X, p.Velocity.Y)
		fmt.Printf("  Total points: %d. Final position: (%.2f, %.2f)\n", len(traj), traj[len(traj)-1].X, traj[len(traj)-1].Y)
		fmt.Print("  Sample Trajectory Points: ")
		for j := 0; j < 3 && j < len(traj); j++ {
			fmt.Printf("(%.2f, %.2f) ", traj[j].X, traj[j].Y)
		}
		if len(traj) > 6 {
			fmt.Print("... ")
			for j := len(traj) - 3; j < len(traj); j++ {
				fmt.Printf("(%.2f, %.2f) ", traj[j].X, traj[j].Y)
			}
		}
		fmt.Println()
	}
}

// Additional implementation at 2025-06-21 00:51:46
package main

import (
	"fmt"
	"math"
)

// Point represents a 2D coordinate.
type Point struct {
	X, Y float64
}

// Projectile represents the state of a projectile.
type Projectile struct {
	Position Point
	Velocity Point // Velocity components (Vx, Vy)
}

// simulateProjectile simulates the trajectory of a projectile over time.
// It returns a slice of trajectory points.
func simulateProjectile(initialVelocity, angleDegrees, gravity, timeStep float64) []Point {
	angleRad := angleDegrees * math.Pi / 180.0

	// Initial velocity components
	initialVx := initialVelocity * math.Cos(angleRad)
	initialVy := initialVelocity * math.Sin(angleRad)

	p := Projectile{
		Position: Point{X: 0, Y: 0},
		Velocity: Point{X: initialVx, Y: initialVy},
	}

	trajectory := []Point{p.Position} // Start with the initial position

	// Simulate until the projectile hits or goes below the ground (Y < 0)
	for {
		// Calculate the next state
		nextX := p.Position.X + p.Velocity.X*timeStep
		nextY := p.Position.Y + p.Velocity.Y*timeStep
		nextVy := p.Velocity.Y - gravity*timeStep // Vx remains constant

		// If the next Y position is negative, the projectile has hit the ground.
		// We stop the simulation here. The analytical formulas will give the exact range.
		if nextY < 0 {
			break
		}

		// Update projectile state
		p.Position.X = nextX
		p.Position.Y = nextY
		p.Velocity.Y = nextVy

		trajectory = append(trajectory, p.Position)
	}

	return trajectory
}

func main() {
	const gravity = 9.81  // m/s^2 (acceleration due to gravity)
	const timeStep = 0.01 // seconds (time increment for simulation)

	var initialVelocity float64
	var angleDegrees float64

	fmt.Print("Enter initial velocity (m/s): ")
	_, err := fmt.Scanln(&initialVelocity)
	if err != nil {
		fmt.Println("Invalid input for velocity. Using default 50 m/s.")
		initialVelocity = 50.0
	}

	fmt.Print("Enter launch angle (degrees from horizontal, 0-90): ")
	_, err = fmt.Scanln(&angleDegrees)
	if err != nil {
		fmt.Println("Invalid input for angle. Using default 45 degrees.")
		angleDegrees = 45.0
	}

	// Input validation and adjustment
	if initialVelocity < 0 {
		fmt.Println("Initial velocity cannot be negative. Setting to 0.")
		initialVelocity = 0
	}
	if angleDegrees < 0 {
		fmt.Println("Launch angle cannot be negative. Setting to 0 degrees.")
		angleDegrees = 0
	} else if angleDegrees > 90 {
		fmt.Println("Launch angle cannot exceed 90 degrees. Setting to 90 degrees.")
		angleDegrees = 90
	}

	// Perform the simulation to get trajectory points
	trajectory := simulateProjectile(initialVelocity, angleDegrees, gravity, timeStep)

	// Calculate analytical results for precision
	angleRad := angleDegrees * math.Pi / 180.0

	// Time of Flight (T = 2 * v0 * sin(theta) / g)
	timeOfFlight := (2 * initialVelocity * math.Sin(angleRad)) / gravity

	// Maximum Height (H_max = (v0 * sin(theta))^2 / (2 * g))
	maxHeight := math.Pow(initialVelocity*math.Sin(angleRad), 2) / (2 * gravity)

	// Horizontal Range (R = v0^2 * sin(2*theta) / g)
	horizontalRange := math.Pow(initialVelocity, 2) * math.Sin(2*angleRad) / gravity

	fmt.Println("\n--- Simulation Parameters ---")
	fmt.Printf("Initial Velocity: %.2f m/s\n", initialVelocity)
	fmt.Printf("Launch Angle: %.2f degrees\n", angleDegrees)
	fmt.Printf("Gravity: %.2f m/s^2\n", gravity)
	fmt.Printf("Time Step: %.4f s\n", timeStep)

	fmt.Println("\n--- Calculated Analytical Results ---")
	fmt.Printf("Time of Flight: %.4f seconds\n", timeOfFlight)
	fmt.Printf("Maximum Height: %.4f meters\n", maxHeight)
	fmt.Printf("Horizontal Range: %.4f meters\n", horizontalRange)

	fmt.Println("\n--- Trajectory Points (X, Y) ---")
	// Print only a subset of points to avoid excessive output
	// Prints approximately every 0.1 seconds of simulation time
	printInterval := int(0.1 / timeStep)
	if printInterval == 0 { // Ensure at least one point is printed if timeStep is large
		printInterval = 1
	}

	for i, p := range trajectory {
		// Print points at regular intervals and always the first and last points
		if i == 0 || i%printInterval == 0 || i == len(trajectory)-1 {
			fmt.Printf("Time: %.4f s, Position: (%.4f, %.4f)\n", float64(i)*timeStep, p.X, p.Y)
		}
	}
	// Ensure the very last point is printed if it wasn't caught by the interval or was the only point
	if len(trajectory) > 0 && (len(trajectory)-1)%printInterval != 0 && len(trajectory)-1 != 0 {
		fmt.Printf("Time: %.4f s, Position: (%.4f, %.4f)\n", float64(len(trajectory)-1)*timeStep, trajectory[len(trajectory)-1].X, trajectory[len(trajectory)-1].Y)
	}
}

// Additional implementation at 2025-06-21 00:52:51
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// Projectile represents the state of a projectile
type Projectile struct {
	Mass   float64 // kg
	Radius float64 // m (for calculating cross-sectional area for drag)
	X, Y   float64 // current position (m)
	Vx, Vy float64 // current velocity (m/s)
}

// Simulator holds simulation parameters and controls the simulation
type Simulator struct {
	Gravity             float64 // m/s^2
	AirDensity          float64 // kg/m^3
	TimeStep            float64 // seconds
	MaxTime             float64 // seconds
	EnableAirResistance bool    // Flag to enable/disable air resistance
}

// NewProjectile creates a new projectile with given initial conditions
func NewProjectile(mass, radius, initialX, initialY, initialSpeed, initialAngleDegrees float64) *Projectile {
	angleRad := initialAngleDegrees * math.Pi / 180.0 // Convert angle to radians
	return &Projectile{
		Mass:   mass,
		Radius: radius,
		X:      initialX,
		Y:      initialY,
		Vx:     initialSpeed * math.Cos(angleRad),
		Vy:     initialSpeed * math.Sin(angleRad),
	}
}

// Update advances the projectile's state by one time step
func (s *Simulator) Update(p *Projectile) {
	ax := 0.0
	ay := -s.Gravity // Acceleration due to gravity

	if s.EnableAirResistance {
		// Calculate drag force
		// Fd = 0.5 * rho * Cd * A * v^2
		// Cd (Drag Coefficient) for a sphere is approximately 0.47
		// A (Cross-sectional Area) = pi * r^2
		Cd := 0.47                                 // Drag coefficient for a sphere
		A := math.Pi * p.Radius * p.Radius         // Cross-sectional area

		vTotal := math.Sqrt(p.Vx*p.Vx + p.Vy*p.Vy) // Total velocity magnitude

		if vTotal > 0 { // Avoid division by zero if velocity is zero
			dragForceMagnitude := 0.5 * s.AirDensity * Cd * A * vTotal * vTotal

			// Drag acceleration components
			ax -= dragForceMagnitude * (p.Vx / vTotal) / p.Mass
			ay -= dragForceMagnitude * (p.Vy / vTotal) / p.Mass
		}
	}

	// Update velocity using current accelerations
	p.Vx += ax * s.TimeStep
	p.Vy += ay * s.TimeStep

	// Update position using new velocities
	p.X += p.Vx * s.TimeStep
	p.Y += p.Vy * s.TimeStep
}

// Simulate runs the projectile motion simulation and prints trajectory points
func (s *Simulator) Simulate(p *Projectile) {
	fmt.Println("Time (s), X (m), Y (m), Vx (m/s), Vy (m/s)")
	currentTime := 0.0
	for p.Y >= 0 && currentTime <= s.MaxTime { // Continue until projectile hits ground or max time is reached
		fmt.Printf("%.4f, %.4f, %.4f, %.4f, %.4f\n", currentTime, p.X, p.Y, p.Vx, p.Vy)
		s.Update(p)
		currentTime += s.TimeStep
	}
	// Print the final state when it hits the ground or max time is reached
	// Note: If Y goes slightly negative, this prints the last state before or at ground hit.
	fmt.Printf("%.4f, %.4f, %.4f, %.4f, %.4f\n", currentTime, p.X, p.Y, p.Vx, p.Vy)
}

// getUserInput reads a float64 value from the console with a given prompt
func getUserInput(prompt string) (float64, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	val, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid input: %v", err)
	}
	return val, nil
}

// getUserBoolInput reads a boolean value (y/n) from the console with a given prompt
func getUserBoolInput(prompt string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt + " (y/n): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "y" || input == "yes" {
		return true, nil
	} else if input == "n" || input == "no" {
		return false, nil
	}
	return false, fmt.Errorf("invalid input, please enter 'y' or 'n'")
}

func main() {
	fmt.Println("Go Projectile Motion Simulator")
	fmt.Println("-------------------------------")

	mass, err := getUserInput("Enter projectile mass (kg): ")
	if err != nil {
		fmt.Println(err)
		return
	}
	radius, err := getUserInput("Enter projectile radius (m): ")
	if err != nil {
		fmt.Println(err)
		return
	}
	initialSpeed, err := getUserInput("Enter initial speed (m/s): ")
	if err != nil {
		fmt.Println(err)
		return
	}
	initialAngle, err := getUserInput("Enter launch angle (degrees from horizontal): ")
	if err != nil {
		fmt.Println(err)
		return
	}
	initialX, err := getUserInput("Enter initial X position (m): ")
	if err != nil {
		fmt.Println(err)
		return
	}
	initialY, err := getUserInput("Enter initial Y position (m): ")
	if err != nil {
		fmt.Println(err)
		return
	}

	timeStep, err := getUserInput("Enter simulation time step (seconds, e.g., 0.01): ")
	if err != nil {
		fmt.Println(err)
		return
	}
	maxTime, err := getUserInput("Enter maximum simulation time (seconds, e.g., 100): ")
	if err != nil {
		fmt.Println(err)
		return
	}
	enableAirResistance, err := getUserBoolInput("Enable air resistance?")
	if err != nil {
		fmt.Println(err)
		return
	}

	projectile := NewProjectile(mass, radius, initialX, initialY, initialSpeed, initialAngle)
	simulator := &Simulator{
		Gravity:             9.81,  // Standard gravity on Earth
		AirDensity:          1.225, // Standard air density at sea level (kg/m^3)
		TimeStep:            timeStep,
		MaxTime:             maxTime,
		EnableAirResistance: enableAirResistance,
	}

	fmt.Println("\nStarting Simulation...")
	simulator.Simulate(projectile)
	fmt.Println("Simulation Complete.")
}