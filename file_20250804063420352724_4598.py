import tkinter as tk
import random

def generate_random_color():
    r = random.randint(0, 255)
    g = random.randint(0, 255)
    b = random.randint(0, 255)
    return f'#{r:02x}{g:02x}{b:02x}'

def draw_random_shape(canvas, width, height):
    canvas.delete("all")

    num_shapes = random.randint(5, 20)

    for _ in range(num_shapes):
        shape_type = random.choice(['rectangle', 'oval', 'line'])
        fill_color = generate_random_color()
        outline_color = generate_random_color()
        line_width = random.randint(1, 5)

        if shape_type == 'rectangle':
            x1 = random.randint(0, width - 50)
            y1 = random.randint(0, height - 50)
            x2 = random.randint(x1 + 10, min(x1 + 100, width))
            y2 = random.randint(y1 + 10, min(y1 + 100, height))
            canvas.create_rectangle(x1, y1, x2, y2, fill=fill_color, outline=outline_color, width=line_width)
        elif shape_type == 'oval':
            x1 = random.randint(0, width - 50)
            y1 = random.randint(0, height - 50)
            x2 = random.randint(x1 + 10, min(x1 + 100, width))
            y2 = random.randint(y1 + 10, min(y1 + 100, height))
            canvas.create_oval(x1, y1, x2, y2, fill=fill_color, outline=outline_color, width=line_width)
        elif shape_type == 'line':
            x1 = random.randint(0, width)
            y1 = random.randint(0, height)
            x2 = random.randint(0, width)
            y2 = random.randint(0, height)
            canvas.create_line(x1, y1, x2, y2, fill=outline_color, width=line_width)

def main():
    root = tk.Tk()
    root.title("Random Geometric Shapes")

    canvas_width = 600
    canvas_height = 400

    canvas = tk.Canvas(root, width=canvas_width, height=canvas_height, bg="white")
    canvas.pack(pady=10)

    generate_button = tk.Button(root, text="Generate New Shapes", command=lambda: draw_random_shape(canvas, canvas_width, canvas_height))
    generate_button.pack(pady=5)

    draw_random_shape(canvas, canvas_width, canvas_height) # Draw initial shapes

    root.mainloop()

if __name__ == "__main__":
    main()

# Additional implementation at 2025-08-04 06:34:57
import random
from PIL import Image, ImageDraw

def generate_random_color():
    """Generates a random RGB color tuple."""
    return (random.randint(0, 255), random.randint(0, 255), random.randint(0, 255))

def draw_random_shape(draw, width, height):
    """Draws a single random geometric shape on the given ImageDraw object."""
    shape_type = random.choice(['rectangle', 'circle', 'triangle', 'line'])
    fill_color = generate_random_color()
    outline_color = generate_random_color() if random.random() > 0.5 else None # Optional outline

    # Define a minimum size for shapes to be visible
    min_dim = min(width, height) // 20
    max_dim_factor = 0.8 # Max dimension as a factor of image size

    if shape_type == 'rectangle':
        x1 = random.randint(0, width - min_dim)
        y1 = random.randint(0, height - min_dim)
        x2 = random.randint(x1 + min_dim, min(width, x1 + int(width * max_dim_factor)))
        y2 = random.randint(y1 + min_dim, min(height, y1 + int(height * max_dim_factor)))
        draw.rectangle([x1, y1, x2, y2], fill=fill_color, outline=outline_color)
    elif shape_type == 'circle':
        x1 = random.randint(0, width - min_dim)
        y1 = random.randint(0, height - min_dim)
        # Ensure circle fits and has a minimum size
        max_radius = min(width - x1, height - y1, int(min(width, height) * max_dim_factor / 2))
        if max_radius < min_dim // 2: # Not enough space for min radius
            return # Skip drawing if too small
        radius = random.randint(min_dim // 2, max_radius)
        x2 = x1 + 2 * radius
        y2 = y1 + 2 * radius
        draw.ellipse([x1, y1, x2, y2], fill=fill_color, outline=outline_color)
    elif shape_type == 'triangle':
        points = []
        for _ in range(3):
            points.append((random.randint(0, width), random.randint(0, height)))
        draw.polygon(points, fill=fill_color, outline=outline_color)
    elif shape_type == 'line':
        x1 = random.randint(0, width)
        y1 = random.randint(0, height)
        x2 = random.randint(0, width)
        y2 = random.randint(0, height)
        line_width = random.randint(1, min(width, height) // 50) # Random line width
        draw.line([x1, y1, x2, y2], fill=fill_color, width=line_width)

def generate_random_shapes_image(
    image_width=800,
    image_height=600,
    num_shapes=50,
    background_color=(255, 255, 255), # White background
    output_filename="random_shapes.png"
):
    """
    Generates an image filled with random geometric shapes.

    Args:
        image_width (int): The width of the output image.
        image_height (int): The height of the output image.
        num_shapes (int): The number of random shapes to draw.
        background_color (tuple): RGB tuple for the background color.
        output_filename (str): The name of the file to save the image.
    """
    # Create a new blank image with the specified background color
    image = Image.new("RGB", (image_width, image_height), background_color)
    draw = ImageDraw.Draw(image)

    # Draw multiple random shapes
    for _ in range(num_shapes):
        draw_random_shape(draw, image_width, image_height)

    # Save the generated image
    image.save(output_filename)
    print(f"Generated '{output_filename}' with {num_shapes} shapes.")

if __name__ == "__main__":
    # Example usage:
    generate_random_shapes_image(
        image_width=1024,
        image_height=768,
        num_shapes=75,
        background_color=(0, 0, 0), # Black background
        output_filename="random_shapes_example.png"
    )

    generate_random_shapes_image(
        image_width=600,
        image_height=600,
        num_shapes=100,
        background_color=(255, 255, 255), # White background
        output_filename="random_shapes_square.png"
    )

    generate_random_shapes_image(
        image_width=1920,
        image_height=1080,
        num_shapes=150,
        background_color=(50, 50, 50), # Dark gray background
        output_filename="random_shapes_hd.png"
    )

# Additional implementation at 2025-08-04 06:35:48
from PIL import Image, ImageDraw
import random

WIDTH = 800
HEIGHT = 600
BACKGROUND_COLOR = (255, 255, 255)

img = Image.new('RGB', (WIDTH, HEIGHT), color=BACKGROUND_COLOR)
draw = ImageDraw.Draw(img)

num_shapes = random.randint(5, 20)

SHAPE_TYPES = ['rectangle', 'circle', 'line', 'triangle']

for _ in range(num_shapes):
    shape_type = random.choice(SHAPE_TYPES)

    fill_color = (random.randint(0, 255), random.randint(0, 255), random.randint(0, 255))
    outline_color = (random.randint(0, 255), random.randint(0, 255), random.randint(0, 255))

    draw_mode = random.choice(['fill', 'outline', 'both'])
    
    current_fill = fill_color if draw_mode in ['fill', 'both'] else None
    current_outline = outline_color if draw_mode in ['outline', 'both'] else None
    
    current_line_width = random.randint(1, 5) if current_outline else 0

    if shape_type == 'rectangle':
        x1 = random.randint(0, WIDTH - 50)
        y1 = random.randint(0, HEIGHT - 50)
        x2 = random.randint(x1 + 20, min(x1 + 200, WIDTH))
        y2 = random.randint(y1 + 20, min(y1 + 200, HEIGHT))
        draw.rectangle([x1, y1, x2, y2], fill=current_fill, outline=current_outline, width=current_line_width)

    elif shape_type == 'circle':
        x1 = random.randint(0, WIDTH - 50)
        y1 = random.randint(0, HEIGHT - 50)
        radius = random.randint(10, 100)
        x2 = x1 + radius * 2
        y2 = y1 + radius * 2
        draw.ellipse([x1, y1, x2, y2], fill=current_fill, outline=current_outline, width=current_line_width)

    elif shape_type == 'line':
        x1 = random.randint(0, WIDTH)
        y1 = random.randint(0, HEIGHT)
        x2 = random.randint(0, WIDTH)
        y2 = random.randint(0, HEIGHT)
        draw.line([x1, y1, x2, y2], fill=outline_color, width=random.randint(1, 7))

    elif shape_type == 'triangle':
        p1_x = random.randint(0, WIDTH)
        p1_y = random.randint(0, HEIGHT)
        p2_x = random.randint(0, WIDTH)
        p2_y = random.randint(0, HEIGHT)
        p3_x = random.randint(0, WIDTH)
        p3_y = random.randint(0, HEIGHT)
        
        points = [(p1_x, p1_y), (p2_x, p2_y), (p3_x, p3_y)]
        draw.polygon(points, fill=current_fill, outline=current_outline, width=current_line_width)

output_filename = "random_shapes.png"
img.save(output_filename)

# Additional implementation at 2025-08-04 06:36:34
from PIL import Image, ImageDraw
import random

def generate_random_color():
    return (random.randint(0, 255), random.randint(0, 255), random.randint(0, 255))

def generate_random_shape_image(width=800, height=600, num_shapes=20, output_filename="random_shapes.png"):
    image = Image.new("RGB", (width, height), (255, 255, 255))
    draw = ImageDraw.Draw(image)

    shape_types = ["rectangle", "circle", "triangle", "line"]

    for _ in range(num_shapes):
        shape_type = random.choice(shape_types)
        color = generate_random_color()

        min_dimension = 20
        max_dimension_x = width // 3
        max_dimension_y = height // 3

        if shape_type == "rectangle":
            x1 = random.randint(0, width - min_dimension)
            y1 = random.randint(0, height - min_dimension)
            x2 = random.randint(x1 + min_dimension, min(x1 + max_dimension_x, width))
            y2 = random.randint(y1 + min_dimension, min(y1 + max_dimension_y, height))
            draw.rectangle((x1, y1, x2, y2), fill=color)
        elif shape_type == "circle":
            diameter = random.randint(min_dimension, min(max_dimension_x, max_dimension_y))
            x1 = random.randint(0, width - diameter)
            y1 = random.randint(0, height - diameter)
            x2 = x1 + diameter
            y2 = y1 + diameter
            draw.ellipse((x1, y1, x2, y2), fill=color)
        elif shape_type == "triangle":
            center_x = random.randint(0, width)
            center_y = random.randint(0, height)
            
            max_dev_x = min(max_dimension_x // 2, center_x, width - center_x)
            max_dev_y = min(max_dimension_y // 2, center_y, height - center_y)
            
            if max_dev_x < min_dimension // 2: max_dev_x = min_dimension // 2
            if max_dev_y < min_dimension // 2: max_dev_y = min_dimension // 2

            p1 = (random.randint(center_x - max_dev_x, center_x + max_dev_x),
                  random.randint(center_y - max_dev_y, center_y + max_dev_y))
            p2 = (random.randint(center_x - max_dev_x, center_x + max_dev_x),
                  random.randint(center_y - max_dev_y, center_y + max_dev_y))
            p3 = (random.randint(center_x - max_dev_x, center_x + max_dev_x),
                  random.randint(center_y - max_dev_y, center_y + max_dev_y))
            
            p1 = (max(0, min(width-1, p1[0])), max(0, min(height-1, p1[1])))
            p2 = (max(0, min(width-1, p2[0])), max(0, min(height-1, p2[1])))
            p3 = (max(0, min(width-1, p3[0])), max(0, min(height-1, p3[1])))

            draw.polygon([p1, p2, p3], fill=color)
        elif shape_type == "line":
            x1 = random.randint(0, width)
            y1 = random.randint(0, height)
            x2 = random.randint(0, width)
            y2 = random.randint(0, height)
            line_width = random.randint(1, 5)
            draw.line((x1, y1, x2, y2), fill=color, width=line_width)

    image.save(output_filename)

if __name__ == "__main__":
    generate_random_shape_image(width=1024, height=768, num_shapes=50, output_filename="my_random_shapes.png")
    generate_random_shape_image(width=640, height=480, num_shapes=30, output_filename="another_shapes_image.png")
    generate_random_shape_image(width=300, height=200, num_shapes=10, output_filename="small_shapes.png")