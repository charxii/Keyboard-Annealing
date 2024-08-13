import json
import matplotlib.pyplot as plt
from typing import Dict
import numpy as np


def load_stats(path: str) -> Dict[str, Dict[str, int]]:
    with open(path) as f:
        data = json.load(f)

    replace = [name for name in data.keys() if name.startswith("000 ")]
    for name in replace:
        data[name.replace("000 ", "")] = data[name]
        del data[name]
    return data


if __name__ == "__main__":
    # Load stats
    data = load_stats("stats.json")

    # Prepare categories and stats
    categories = list(data.keys())
    stats = sorted(next(iter(data.values())).keys())
    values = {stat: [data[category][stat] for category in categories] for stat in stats}
    x = np.arange(len(categories))
    width = 0.15

    # Manually set the background color
    background_color = "#2e3440"  # A dark gray-blue color
    text_color = "#d8dee9"  # A light gray color for text

    # Create the figure and axis with the desired background color
    fig, ax = plt.subplots(figsize=(16, 8))
    fig.patch.set_facecolor(background_color)  # Background color for the figure
    ax.set_facecolor(background_color)  # Background color for the plot area

    # Define a color palette suitable for colorblind users
    colors = plt.cm.viridis(np.linspace(0, 1, len(stats)))

    # Plot the bars with the new color palette
    bars = []
    for i, stat in enumerate(stats):
        bar = ax.bar(x + i * width, values[stat], width, label=stat, color=colors[i])
        bars.append(bar)

    # Annotate bars with proportional values
    for stat, bar_group in zip(stats, bars):
        max_value = max(values[stat])
        for bar in bar_group:
            height = bar.get_height()
            percentage = height / max_value
            ax.text(
                bar.get_x() + bar.get_width() / 2,
                height + 0.05,
                "x" + f"{percentage:.3f}".lstrip("0"),
                ha="center",
                va="bottom",
                fontsize=9,
                color=text_color,  # Text color for better contrast
            )

    # Customize labels and title with the specified text color
    ax.set_xlabel("Layout", color=text_color)
    ax.set_ylabel("Percentage", color=text_color)
    ax.set_title("Metric Comparison", color=text_color)
    ax.set_xticks(x + width * (len(stats) - 1) / 2)
    ax.set_xticklabels(categories, rotation=45, ha="right", color=text_color)

    # Set the axes (spines) to white
    ax.spines["bottom"].set_color(text_color)
    ax.spines["left"].set_color(text_color)
    ax.spines["top"].set_color(text_color)
    ax.spines["right"].set_color(text_color)

    # Set the tick parameters (colors)
    ax.tick_params(axis="x", colors=text_color)
    ax.tick_params(axis="y", colors=text_color)

    # Adjust legend and layout
    ax.legend(
        loc="center left",
        bbox_to_anchor=(1, 0.5),
        fontsize=10,
        facecolor=background_color,
        edgecolor="none",
        labelcolor=text_color,
    )
    fig.tight_layout()

    # Show the plot
    plt.show()
