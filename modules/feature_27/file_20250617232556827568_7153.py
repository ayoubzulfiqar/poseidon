import turtle
import random

def draw_circle(t, x, y, radius, color):
    t.penup()
    t.goto(x, y - radius)
    t.pendown()
    t.fillcolor(color)
    t.begin_fill()
    t.circle(radius)
    t.end_fill()

def draw_polygon(t, x, y, side_length, num_sides, color):
    t.penup()
    t.goto(x, y)
    t.pendown()
    t.fillcolor(color)
    t.begin_fill()
    angle = 360 / num_sides
    for _ in range(num_sides):
        t.forward(side_length)
        t.right(angle)
    t.end_fill()

def generate_random_shape(t, screen_width, screen_height):
    shape_types = ["circle", "square", "triangle", "pentagon", "hexagon"]
    chosen_shape = random.choice(shape_types)

    x = random.randint(-screen_width // 2 + 70, screen_width // 2 - 70)
    y = random.randint(-screen_height // 2 + 70, screen_height // 2 - 70)

    r = random.random()
    g = random.random()
    b = random.random()
    color = (r, g, b)

    t.pencolor(color)

    if chosen_shape == "circle":
        radius = random.randint(20, 80)
        draw_circle(t, x, y, radius, color)
    else:
        side_length = random.randint(30, 100)
        num_sides = 0
        if chosen_shape == "square":
            num_sides = 4
        elif chosen_shape == "triangle":
            num_sides = 3
        elif chosen_shape == "pentagon":
            num_sides = 5
        elif chosen_shape == "hexagon":
            num_sides = 6
        draw_polygon(t, x, y, side_length, num_sides, color)

def main():
    screen = turtle.Screen()
    screen.setup(width=800, height=600)
    screen.colormode(1.0)
    screen.tracer(0)

    t = turtle.Turtle()
    t.speed(0)
    t.hideturtle()

    num_shapes = random.randint(10, 30)

    for _ in range(num_shapes):
        generate_random_shape(t, screen.window_width(), screen.window_height())

    screen.update()
    screen.exitonclick()

if __name__ == "__main__":
    main()

# Additional implementation at 2025-06-17 23:26:43
from PIL import Image, ImageDraw
import random
import math

def get_random_color():
    return (random.randint(0, 255), random.randint(0, 255), random.randint(0, 255))

def generate_shapes_image(width=800, height=600, min_shapes=5, max_shapes=20, output_filename="random_shapes.png"):
    image = Image.new("RGB", (width, height), (255, 255, 255))
    draw = ImageDraw.Draw(image)

    num_shapes = random.randint(min_shapes, max_shapes)

    for _ in range(num_shapes):
        shape_type = random.choice(["circle", "rectangle", "triangle", "polygon"])
        fill_color = get_random_color()
        outline_color = get_random_color() if random.random() > 0.5 else None

        if shape_type == "circle":
            radius = random.randint(10, min(width, height) // 8)
            center_x = random.randint(radius, width - radius)
            center_y = random.randint(radius, height - radius)
            bbox = (center_x - radius, center_y - radius, center_x + radius, center_y + radius)
            draw.ellipse(bbox, fill=fill_color, outline=outline_color)

        elif shape_type == "rectangle":
            x1 = random.randint(0, width - 20)
            y1 = random.randint(0, height - 20)
            x2 = random.randint(x1 + 20, min(x1 + width // 4, width))
            y2 = random.randint(y1 + 20, min(y1 + height // 4, height))
            bbox = (x1, y1, x2, y2)
            draw.rectangle(bbox, fill=fill_color, outline=outline_color)

        elif shape_type == "triangle":
            center_x = random.randint(width // 4, width * 3 // 4)
            center_y = random.randint(height // 4, height * 3 // 4)
            max_offset = min(width, height) // 6

            p1_x = random.randint(max(0, center_x - max_offset), min(width, center_x + max_offset))
            p1_y = random.randint(max(0, center_y - max_offset), min(height, center_y + max_offset))
            p2_x = random.randint(max(0, center_x - max_offset), min(width, center_x + max_offset))
            p2_y = random.randint(max(0, center_y - max_offset), min(height, center_y + max_offset))
            p3_x = random.randint(max(0, center_x - max_offset), min(width, center_x + max_offset))
            p3_y = random.randint(max(0, center_y - max_offset), min(height, center_y + max_offset))

            points = [(p1_x, p1_y), (p2_x, p2_y), (p3_x, p3_y)]
            draw.polygon(points, fill=fill_color, outline=outline_color)

        elif shape_type == "polygon":
            num_sides = random.randint(3, 8)
            center_x = random.randint(width // 4, width * 3 // 4)
            center_y = random.randint(height // 4, height * 3 // 4)
            radius = random.randint(20, min(width, height) // 5)

            points = []
            start_angle = random.uniform(0, 2 * math.pi)
            for i in range(num_sides):
                angle = start_angle + (2 * math.pi * i / num_sides) + random.uniform(-0.5, 0.5)
                x = center_x + radius * math.cos(angle)
                y = center_y + radius * math.sin(angle)
                points.append((x, y))
            draw.polygon(points, fill=fill_color, outline=outline_color)

    image.save(output_filename)

if __name__ == "__main__":
    generate_shapes_image()