#!/usr/bin/env python3
"""
visualize-comparison.py
Generate visualization charts for language expressiveness comparison
"""

import json
import sys
import os
from pathlib import Path

try:
    import matplotlib.pyplot as plt
    import numpy as np
    HAS_MATPLOTLIB = True
except ImportError:
    HAS_MATPLOTLIB = False
    print("Warning: matplotlib not installed. Visualization will be skipped.")
    print("Install with: pip install matplotlib numpy")

SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent
DATA_FILE = PROJECT_ROOT / "docs" / "comparison-loc-data.json"
OUTPUT_DIR = PROJECT_ROOT / "docs" / "comparison-visualizations"

def load_data():
    """Load comparison data from JSON file"""
    if not DATA_FILE.exists():
        print(f"Error: Data file not found: {DATA_FILE}")
        print("Please run scripts/calculate-loc.sh first to generate data.")
        sys.exit(1)
    
    with open(DATA_FILE, 'r') as f:
        return json.load(f)

def create_bar_chart(data):
    """Create bar chart comparing LOC"""
    if not HAS_MATPLOTLIB:
        return
    
    scenarios = [s['name'] for s in data['scenarios']]
    uad_locs = [s['uad_loc'] for s in data['scenarios']]
    other_locs = [s['other_loc'] for s in data['scenarios']]
    
    x = np.arange(len(scenarios))
    width = 0.35
    
    fig, ax = plt.subplots(figsize=(12, 6))
    
    bars1 = ax.bar(x - width/2, uad_locs, width, label='UAD', color='#4A90E2')
    bars2 = ax.bar(x + width/2, other_locs, width, label='Other Language', color='#6C757D')
    
    ax.set_xlabel('Scenario', fontsize=12)
    ax.set_ylabel('Lines of Code (LOC)', fontsize=12)
    ax.set_title('UAD vs Other Languages: Lines of Code Comparison', fontsize=14, fontweight='bold')
    ax.set_xticks(x)
    ax.set_xticklabels(scenarios, rotation=45, ha='right')
    ax.legend()
    ax.grid(axis='y', alpha=0.3)
    
    # Add value labels on bars
    for bars in [bars1, bars2]:
        for bar in bars:
            height = bar.get_height()
            ax.text(bar.get_x() + bar.get_width()/2., height,
                   f'{int(height)}',
                   ha='center', va='bottom', fontsize=9)
    
    plt.tight_layout()
    
    output_file = OUTPUT_DIR / 'loc-comparison-bar.png'
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    plt.savefig(output_file, dpi=300, bbox_inches='tight')
    print(f"Created: {output_file}")
    plt.close()

def create_reduction_chart(data):
    """Create chart showing code reduction percentage"""
    if not HAS_MATPLOTLIB:
        return
    
    scenarios = [s['name'] for s in data['scenarios']]
    reductions = [s['reduction_percent'] for s in data['scenarios']]
    
    fig, ax = plt.subplots(figsize=(10, 6))
    
    colors = ['#28A745' if r > 30 else '#FFC107' if r > 20 else '#DC3545' for r in reductions]
    bars = ax.barh(scenarios, reductions, color=colors)
    
    ax.set_xlabel('Code Reduction (%)', fontsize=12)
    ax.set_title('UAD Code Reduction vs Other Languages', fontsize=14, fontweight='bold')
    ax.grid(axis='x', alpha=0.3)
    
    # Add value labels
    for i, (bar, reduction) in enumerate(zip(bars, reductions)):
        width = bar.get_width()
        ax.text(width, bar.get_y() + bar.get_height()/2.,
               f'{reduction:.1f}%',
               ha='left' if width > 0 else 'right', va='center', fontsize=10, fontweight='bold')
    
    plt.tight_layout()
    
    output_file = OUTPUT_DIR / 'reduction-percentage.png'
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    plt.savefig(output_file, dpi=300, bbox_inches='tight')
    print(f"Created: {output_file}")
    plt.close()

def create_radar_chart(data):
    """Create radar chart for concept mapping intuitiveness"""
    if not HAS_MATPLOTLIB:
        return
    
    # Concept mapping scores (subjective, from comparison report)
    scenarios = [s['name'] for s in data['scenarios']]
    
    # Scores: UAD vs Other (from comparison analysis)
    uad_scores = [5.0, 5.0, 5.0]  # UAD always 5/5
    other_scores = [3.0, 2.0, 3.0]  # Go: 3, 2; Python: 3
    
    # Categories for radar chart
    categories = scenarios
    num_vars = len(categories)
    
    # Compute angle for each category
    angles = [n / float(num_vars) * 2 * np.pi for n in range(num_vars)]
    angles += angles[:1]  # Complete the circle
    
    fig, ax = plt.subplots(figsize=(10, 10), subplot_kw=dict(projection='polar'))
    
    # Add scores
    uad_scores += uad_scores[:1]
    other_scores += other_scores[:1]
    
    # Plot
    ax.plot(angles, uad_scores, 'o-', linewidth=2, label='UAD', color='#4A90E2')
    ax.fill(angles, uad_scores, alpha=0.25, color='#4A90E2')
    
    ax.plot(angles, other_scores, 'o-', linewidth=2, label='Other Languages', color='#6C757D')
    ax.fill(angles, other_scores, alpha=0.25, color='#6C757D')
    
    # Customize
    ax.set_xticks(angles[:-1])
    ax.set_xticklabels(categories, fontsize=10)
    ax.set_ylim(0, 5)
    ax.set_yticks([1, 2, 3, 4, 5])
    ax.set_yticklabels(['1', '2', '3', '4', '5'], fontsize=9)
    ax.set_title('Concept Mapping Intuitiveness (1-5 scale)\nUAD vs Other Languages', 
                 fontsize=14, fontweight='bold', pad=20)
    ax.legend(loc='upper right', bbox_to_anchor=(1.3, 1.1))
    ax.grid(True)
    
    output_file = OUTPUT_DIR / 'concept-mapping-radar.png'
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    plt.savefig(output_file, dpi=300, bbox_inches='tight')
    print(f"Created: {output_file}")
    plt.close()

def generate_markdown_report(data):
    """Generate markdown report with embedded charts"""
    output_file = OUTPUT_DIR / 'comparison-report.md'
    
    with open(output_file, 'w') as f:
        f.write("# Language Expressiveness Comparison Report\n\n")
        f.write(f"**Generated:** {data.get('comparison_date', 'N/A')}\n\n")
        f.write("## Summary Statistics\n\n")
        f.write("| Scenario | UAD LOC | Other LOC | Reduction |\n")
        f.write("|----------|---------|-----------|-----------|\n")
        
        for scenario in data['scenarios']:
            f.write(f"| {scenario['name']} | {scenario['uad_loc']} | "
                   f"{scenario['other_loc']} ({scenario['other_language']}) | "
                   f"{scenario['reduction_percent']:.1f}% |\n")
        
        f.write("\n## Visualizations\n\n")
        f.write("### Lines of Code Comparison\n\n")
        f.write("![LOC Comparison](loc-comparison-bar.png)\n\n")
        f.write("### Code Reduction Percentage\n\n")
        f.write("![Reduction Percentage](reduction-percentage.png)\n\n")
        f.write("### Concept Mapping Intuitiveness\n\n")
        f.write("![Concept Mapping](concept-mapping-radar.png)\n\n")
    
    print(f"Created: {output_file}")

def main():
    """Main function"""
    print("Generating comparison visualizations...")
    print("=" * 60)
    
    data = load_data()
    
    if HAS_MATPLOTLIB:
        create_bar_chart(data)
        create_reduction_chart(data)
        create_radar_chart(data)
    else:
        print("Skipping visualization (matplotlib not available)")
    
    generate_markdown_report(data)
    
    print("=" * 60)
    print("Visualization complete!")
    print(f"Output directory: {OUTPUT_DIR}")

if __name__ == "__main__":
    main()


