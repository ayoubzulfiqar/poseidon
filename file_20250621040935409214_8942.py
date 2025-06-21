import math
import matplotlib.pyplot as plt

def simulate_projectile(initial_velocity, launch_angle_deg, time_step=0.01):
    g = 9.81  # Acceleration due to gravity (m/s^2)
    
    launch_angle_rad = math.radians(launch_angle_deg)
    
    vx0 = initial_velocity * math.cos(launch_angle_rad)
    vy0 = initial_velocity * math.sin(launch_angle_rad)
    
    x_coords = [0.0]
    y_coords = [0.0]
    
    time = 0.0
    
    while True:
        time += time_step
        
        x = vx0 * time
        y = vy0 * time - 0.5 * g * time**2
        
        x_coords.append(x)
        y_coords.append(y)
        
        if y <= 0 and time > 0:
            if y < 0:
                y_coords[-1] = 0.0  # Ensure the last point is exactly on the ground
            break
            
    return x_coords, y_coords

def main():
    initial_velocity = 50.0  # m/s
    launch_angle_deg = 45.0  # degrees
    
    x_trajectory, y_trajectory = simulate_projectile(initial_velocity, launch_angle_deg)
    
    plt.figure(figsize=(10, 6))
    plt.plot(x_trajectory, y_trajectory, label=f'Trajectory (v0={initial_velocity} m/s, angle={launch_angle_deg}°)')
    plt.xlabel('Horizontal Distance (m)')
    plt.ylabel('Vertical Distance (m)')
    plt.title('Projectile Motion Simulation')
    plt.grid(True)
    plt.axhline(0, color='black', linewidth=0.5)
    plt.axvline(0, color='black', linewidth=0.5)
    plt.legend()
    plt.gca().set_aspect('equal', adjustable='box')
    plt.show()

if __name__ == "__main__":
    main()

# Additional implementation at 2025-06-21 04:10:00
import math
import matplotlib.pyplot as plt

GRAVITY = 9.81
TIME_STEP = 0.01

class Projectile:
    def __init__(self, initial_velocity, launch_angle_deg, mass, drag_coefficient, color='blue', label='Projectile'):
        self.initial_velocity = initial_velocity
        self.launch_angle_rad = math.radians(launch_angle_deg)
        self.mass = mass
        self.drag_coefficient = drag_coefficient

        self.x = 0.0
        self.y = 0.0
        self.vx = initial_velocity * math.cos(self.launch_angle_rad)
        self.vy = initial_velocity * math.sin(self.launch_angle_rad)

        self.trajectory_x = [self.x]
        self.trajectory_y = [self.y]
        self.time_in_air = 0.0
        self.max_height = 0.0
        self.range = 0.0
        self.landed = False

        self.color = color
        self.label = label

    def update(self, dt):
        if self.landed:
            return

        speed = math.sqrt(self.vx**2 + self.vy**2)
        
        ax_drag = 0.0
        ay_drag = 0.0

        if speed > 0 and self.drag_coefficient > 0:
            # Quadratic drag: Fd = -k * v^2 (direction opposite to velocity)
            # Fd_magnitude = self.drag_coefficient * speed**2
            # ax_drag = -Fd_magnitude * (self.vx / speed) / self.mass
            # ay_drag = -Fd_magnitude * (self.vy / speed) / self.mass
            
            # Simplified form for quadratic drag components
            ax_drag = -self.drag_coefficient * self.vx * speed / self.mass
            ay_drag = -self.drag_coefficient * self.vy * speed / self.mass

        ax_net = ax_drag
        ay_net = -GRAVITY + ay_drag

        self.vx += ax_net * dt
        self.vy += ay_net * dt

        self.x += self.vx * dt
        self.y += self.vy * dt

        self.time_in_air += dt

        if self.y <= 0:
            self.y = 0
            self.landed = True
            self.range = self.x

        if self.y > self.max_height:
            self.max_height = self.y

        self.trajectory_x.append(self.x)
        self.trajectory_y.append(self.y)

def simulate_projectiles(projectiles):
    all_landed = False
    while not all_landed:
        all_landed = True
        for proj in projectiles:
            if not proj.landed:
                proj.update(TIME_STEP)
                if not proj.landed:
                    all_landed = False

if __name__ == "__main__":
    projectiles_to_simulate = []

    projectiles_to_simulate.append(
        Projectile(
            initial_velocity=50,
            launch_angle_deg=45,
            mass=1.0,
            drag_coefficient=0.0,
            color='blue',
            label='No Air Resistance (50m/s, 45°)'
        )
    )

    projectiles_to_simulate.append(
        Projectile(
            initial_velocity=50,
            launch_angle_deg=45,
            mass=1.0,
            drag_coefficient=0.01,
            color='red',
            label='With Air Resistance (50m/s, 45°)'
        )
    )

    projectiles_to_simulate.append(
        Projectile(
            initial_velocity=70,
            launch_angle_deg=30,
            mass=1.0,
            drag_coefficient=0.005,
            color='green',
            label='Faster, Lower Angle (70m/s, 30°)'
        )
    )

    simulate_projectiles(projectiles_to_simulate)

    plt.figure(figsize=(10, 6))
    max_x_plot = 0
    max_y_plot = 0

    for proj in projectiles_to_simulate:
        plt.plot(proj.trajectory_x, proj.trajectory_y, color=proj.color, label=proj.label)
        max_x_plot = max(max_x_plot, proj.range)
        max_y_plot = max(max_y_plot, proj.max_height)

        print(f"--- {proj.label} ---")
        print(f"  Initial Velocity: {proj.initial_velocity:.2f} m/s")
        print(f"  Launch Angle: {math.degrees(proj.launch_angle_rad):.2f}°")
        print(f"  Mass: {proj.mass:.2f} kg")
        print(f"  Drag Coefficient: {proj.drag_coefficient:.4f}")
        print(f"  Max Height: {proj.max_height:.2f} m")
        print(f"  Range: {proj.range:.2f} m")
        print(f"  Time of Flight: {proj.time_in_air:.2f} s")
        print("-" * 30)

    plt.title('Projectile Motion Simulation')
    plt.xlabel('Horizontal Distance (m)')
    plt.ylabel('Vertical Distance (m)')
    plt.grid(True)
    plt.axhline(0, color='black', linewidth=0.5)
    plt.xlim(0, max_x_plot * 1.1)
    plt.ylim(0, max_y_plot * 1.1)
    plt.legend()
    plt.gca().set_aspect('equal', adjustable='box')
    plt.show()

# Additional implementation at 2025-06-21 04:11:14
import math
import matplotlib.pyplot as plt

G = 9.81
RHO = 1.225

class Projectile:
    def __init__(self, mass, radius, initial_speed, launch_angle_deg, initial_height=0.0, drag_coefficient=0.47):
        self.mass = mass
        self.radius = radius
        self.area = math.pi * (self.radius ** 2)
        self.drag_coefficient = drag_coefficient

        self.x = 0.0
        self.y = initial_height
        self.vx = initial_speed * math.cos(math.radians(launch_angle_deg))
        self.vy = initial_speed * math.sin(math.radians(launch_angle_deg))

        self.time = 0.0
        self.path_x = [self.x]
        self.path_y = [self.y]
        self.max_height = self.y

    def _calculate_drag_force(self):
        speed = math.sqrt(self.vx**2 + self.vy**2)
        if speed == 0:
            return 0.0, 0.0
        
        drag_magnitude = 0.5 * RHO * speed**2 * self.drag_coefficient * self.area
        
        drag_fx = -drag_magnitude * (self.vx / speed)
        drag_fy = -drag_magnitude * (self.vy / speed)
        
        return drag_fx, drag_fy

    def update(self, dt):
        drag_fx, drag_fy = self._calculate_drag_force()
        
        net_fx = drag_fx
        net_fy = -self.mass * G + drag_fy

        ax = net_fx / self.mass
        ay = net_fy / self.mass

        self.vx += ax * dt
        self.vy += ay * dt

        self.x += self.vx * dt
        self.y += self.vy * dt
        self.time += dt

        self.path_x.append(self.x)
        self.path_y.append(self.y)
        if self.y > self.max_height:
            self.max_height = self.y

        if self.y <= 0:
            self.y = 0
            return False
        return True

def run_simulation(projectile, dt=0.01):
    print(f"Simulating projectile with mass={projectile.mass}kg, radius={projectile.radius}m, drag_coeff={projectile.drag_coefficient}")
    print(f"Initial speed={math.sqrt(projectile.vx**2 + projectile.vy**2):.2f} m/s, Angle={math.degrees(math.atan2(projectile.vy, projectile.vx)):.2f} deg, Initial height={projectile.path_y[0]:.2f} m")

    while projectile.update(dt):
        pass

    print("\n--- Simulation Results ---")
    print(f"Time of Flight: {projectile.time:.2f} seconds")
    print(f"Range: {projectile.x:.2f} meters")
    print(f"Maximum Height: {projectile.max_height:.2f} meters")

    plt.figure(figsize=(10, 6))
    plt.plot(projectile.path_x, projectile.path_y, label='Projectile Trajectory')
    plt.xlabel('Horizontal Distance (m)')
    plt.ylabel('Vertical Distance (m)')
    plt.title('Projectile Motion Simulation')
    plt.grid(True)
    plt.axhline(0, color='black', linestyle='--', linewidth=0.8)
    plt.legend()
    plt.axis('equal')
    plt.ylim(bottom=0)
    plt.show()

if __name__ == "__main__":
    try:
        mass = float(input("Enter projectile mass (kg): "))
        radius = float(input("Enter projectile radius (m): "))
        initial_speed = float(input("Enter initial speed (m/s): "))
        launch_angle_deg = float(input("Enter launch angle (degrees from horizontal): "))
        initial_height_str = input("Enter initial height (m, default 0): ")
        initial_height = float(initial_height_str) if initial_height_str else 0.0
        
        drag_coeff_input = input("Enter drag coefficient (e.g., 0.47 for sphere, leave blank for default): ")
        drag_coefficient = float(drag_coeff_input) if drag_coeff_input else 0.47

        dt = 0.01

        my_projectile = Projectile(mass, radius, initial_speed, launch_angle_deg, initial_height, drag_coefficient)
        run_simulation(my_projectile, dt)

    except ValueError:
        print("Invalid input. Please enter numeric values.")
    except Exception as e:
        print(f"An error occurred: {e}")