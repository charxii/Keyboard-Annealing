from PIL import Image, ImageDraw, ImageFont
import json
from typing import Dict


def create_keyboard_image(title: str, layout_string: str):
    rows = [10, 11, 10]
    key_width = 50
    key_height = 50
    spacing = 10
    corner_radius = 10
    margin = 20
    font_size = 20
    title_font_size = 24
    title_height = 40

    image_width = (key_width + spacing) * max(rows) - spacing + 2 * margin
    image_height = (
        (key_height + spacing) * len(rows) - spacing + 2 * margin + title_height
    )

    image = Image.new("RGB", (image_width, image_height), "white")
    draw = ImageDraw.Draw(image)

    try:
        font = ImageFont.truetype("arial.ttf", font_size)
        title_font = ImageFont.truetype("arial.ttf", title_font_size)
    except IOError:
        font = ImageFont.load_default()
        title_font = ImageFont.load_default()

    # Draw title
    title_width, title_height = draw.textsize(title, font=title_font)
    title_x = (image_width - title_width) // 2
    title_y = margin // 2
    draw.text((title_x, title_y), title, font=title_font, fill="black")

    index = 0
    for row in range(len(rows)):
        for col in range(rows[row]):
            x = margin + col * (key_width + spacing)
            y = margin + row * (key_height + spacing) + title_height
            draw.rounded_rectangle(
                [x, y, x + key_width, y + key_height],
                radius=corner_radius,
                outline="black",
                fill="lightgray",
            )

            key_text = layout_string[index]
            text_width, text_height = draw.textsize(key_text, font=font)
            text_x = x + (key_width - text_width) // 2
            text_y = y + (key_height - text_height) // 2
            draw.text((text_x, text_y), key_text, font=font, fill="black")
            index += 1

    return image


def stack_images_vertically(images):
    max_width = max(image.width for image in images)
    total_height = sum(image.height for image in images)
    stacked_image = Image.new("RGB", (max_width, total_height))
    y_offset = 0
    for image in images:
        stacked_image.paste(image, (0, y_offset))
        y_offset += image.height

    return stacked_image


def load_stats(path: str) -> Dict[str, Dict[str, int]]:
    with open(path) as f:
        data = json.load(f)

    replace = [name for name in data.keys() if name.startswith("000 ")]
    for name in replace:
        data[name.replace("000 ", "")] = data[name]
        del data[name]
    return data


if __name__ == "__main__":
    keyboards = load_stats("layouts.json")
    keyboards = {i: j for i, j in keyboards.items() if i.startswith("optimized")}

    images = [create_keyboard_image(*i) for i in keyboards.items()]

    combined_image = stack_images_vertically(images)

    combined_image.save("combined.png")
