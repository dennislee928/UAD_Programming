#!/bin/bash

# calculate-loc.sh
# Calculate Lines of Code (LOC) for comparison files
# Excludes empty lines and comments

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"

echo -e "${GREEN}UAD Language Expressiveness Comparison - LOC Statistics${NC}"
echo "=================================================================="
echo ""

# Function to count LOC (excluding empty lines and comments)
count_loc() {
    local file=$1
    local lang=$2
    
    if [ ! -f "$file" ]; then
        echo "0"
        return
    fi
    
    case "$lang" in
        "uad")
            # For UAD: exclude empty lines and comments (// and /* */)
            grep -v '^\s*$' "$file" | \
            grep -v '^\s*//' | \
            sed 's|/\*.*\*/||g' | \
            sed '/\/\*/,/\*\//d' | \
            grep -v '^\s*$' | \
            wc -l | tr -d ' '
            ;;
        "go")
            # For Go: exclude empty lines and comments
            grep -v '^\s*$' "$file" | \
            grep -v '^\s*//' | \
            sed 's|/\*.*\*/||g' | \
            sed '/\/\*/,/\*\//d' | \
            grep -v '^\s*$' | \
            wc -l | tr -d ' '
            ;;
        "python")
            # For Python: exclude empty lines and comments
            grep -v '^\s*$' "$file" | \
            grep -v '^\s*#' | \
            grep -v '^\s*""".*"""$' | \
            sed '/"""/,/"""/d' | \
            grep -v '^\s*$' | \
            wc -l | tr -d ' '
            ;;
        "rust")
            # For Rust: exclude empty lines and comments
            grep -v '^\s*$' "$file" | \
            grep -v '^\s*//' | \
            sed 's|/\*.*\*/||g' | \
            sed '/\/\*/,/\*\//d' | \
            grep -v '^\s*$' | \
            wc -l | tr -d ' '
            ;;
        *)
            # Default: just exclude empty lines
            grep -v '^\s*$' "$file" | wc -l | tr -d ' '
            ;;
    esac
}

# Count LOC for each scenario
echo -e "${YELLOW}Scenario 1: Ransomware Kill-Chain${NC}"
echo "----------------------------------------"
uad_file="$PROJECT_ROOT/examples/showcase/ransomware_killchain.uad"
go_file="$PROJECT_ROOT/examples/showcase/comparison/go/ransomware_killchain.go"

if [ -f "$uad_file" ]; then
    uad_loc=$(count_loc "$uad_file" "uad")
    echo "  UAD:   $uad_loc lines"
else
    echo "  UAD:   File not found"
    uad_loc=0
fi

if [ -f "$go_file" ]; then
    go_loc=$(count_loc "$go_file" "go")
    echo "  Go:    $go_loc lines"
else
    echo "  Go:    File not found"
    go_loc=0
fi

if [ "$uad_loc" -gt 0 ] && [ "$go_loc" -gt 0 ]; then
    reduction=$(echo "scale=1; (1 - $uad_loc / $go_loc) * 100" | bc)
    echo "  Reduction: ${reduction}%"
fi
echo ""

echo -e "${YELLOW}Scenario 2: SOC / JudgePipeline ERH Model${NC}"
echo "----------------------------------------"
uad_file="$PROJECT_ROOT/examples/showcase/psychohistory_scenario.uad"
go_file="$PROJECT_ROOT/examples/showcase/comparison/go/soc_erh_model.go"

if [ -f "$uad_file" ]; then
    # Estimate: ERH model portion is roughly 40% of psychohistory file
    total_loc=$(count_loc "$uad_file" "uad")
    uad_loc=$(echo "scale=0; $total_loc * 0.4" | bc | cut -d. -f1)
    echo "  UAD:   ~$uad_loc lines (estimated, ERH portion)"
else
    echo "  UAD:   File not found"
    uad_loc=0
fi

if [ -f "$go_file" ]; then
    go_loc=$(count_loc "$go_file" "go")
    echo "  Go:    $go_loc lines"
else
    echo "  Go:    File not found"
    go_loc=0
fi

if [ "$uad_loc" -gt 0 ] && [ "$go_loc" -gt 0 ]; then
    reduction=$(echo "scale=1; (1 - $uad_loc / $go_loc) * 100" | bc)
    echo "  Reduction: ${reduction}%"
fi
echo ""

echo -e "${YELLOW}Scenario 3: Psychohistory Macro Dynamics${NC}"
echo "----------------------------------------"
uad_file="$PROJECT_ROOT/examples/showcase/psychohistory_scenario.uad"
python_file="$PROJECT_ROOT/examples/showcase/comparison/python/psychohistory_scenario.py"

if [ -f "$uad_file" ]; then
    uad_loc=$(count_loc "$uad_file" "uad")
    echo "  UAD:   $uad_loc lines"
else
    echo "  UAD:   File not found"
    uad_loc=0
fi

if [ -f "$python_file" ]; then
    python_loc=$(count_loc "$python_file" "python")
    echo "  Python: $python_loc lines"
else
    echo "  Python: File not found"
    python_loc=0
fi

if [ "$uad_loc" -gt 0 ] && [ "$python_loc" -gt 0 ]; then
    reduction=$(echo "scale=1; (1 - $uad_loc / $python_loc) * 100" | bc)
    echo "  Reduction: ${reduction}%"
fi
echo ""

# Summary table
echo "=================================================================="
echo -e "${GREEN}Summary${NC}"
echo "=================================================================="
printf "%-40s %10s %10s %10s\n" "Scenario" "UAD" "Other" "Reduction"
echo "------------------------------------------------------------------"

# Scenario 1
uad_file="$PROJECT_ROOT/examples/showcase/ransomware_killchain.uad"
go_file="$PROJECT_ROOT/examples/showcase/comparison/go/ransomware_killchain.go"
if [ -f "$uad_file" ] && [ -f "$go_file" ]; then
    uad_loc=$(count_loc "$uad_file" "uad")
    go_loc=$(count_loc "$go_file" "go")
    reduction=$(echo "scale=1; (1 - $uad_loc / $go_loc) * 100" | bc)
    printf "%-40s %10d %10d %9.1f%%\n" "Ransomware Kill-Chain" "$uad_loc" "$go_loc" "$reduction"
fi

# Scenario 2
uad_file="$PROJECT_ROOT/examples/showcase/psychohistory_scenario.uad"
go_file="$PROJECT_ROOT/examples/showcase/comparison/go/soc_erh_model.go"
if [ -f "$uad_file" ] && [ -f "$go_file" ]; then
    total_loc=$(count_loc "$uad_file" "uad")
    uad_loc=$(echo "scale=0; $total_loc * 0.4" | bc | cut -d. -f1)
    go_loc=$(count_loc "$go_file" "go")
    reduction=$(echo "scale=1; (1 - $uad_loc / $go_loc) * 100" | bc)
    printf "%-40s %10d %10d %9.1f%%\n" "SOC ERH Model" "$uad_loc" "$go_loc" "$reduction"
fi

# Scenario 3
uad_file="$PROJECT_ROOT/examples/showcase/psychohistory_scenario.uad"
python_file="$PROJECT_ROOT/examples/showcase/comparison/python/psychohistory_scenario.py"
if [ -f "$uad_file" ] && [ -f "$python_file" ]; then
    uad_loc=$(count_loc "$uad_file" "uad")
    python_loc=$(count_loc "$python_file" "python")
    reduction=$(echo "scale=1; (1 - $uad_loc / $python_loc) * 100" | bc)
    printf "%-40s %10d %10d %9.1f%%\n" "Psychohistory" "$uad_loc" "$python_loc" "$reduction"
fi

echo ""

# Save to JSON for visualization
JSON_OUTPUT="$PROJECT_ROOT/docs/comparison-loc-data.json"
echo "Saving data to $JSON_OUTPUT..."
cat > "$JSON_OUTPUT" << EOF
{
  "comparison_date": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "scenarios": [
EOF

# Add scenario data
scenarios_data=""
first=true
for scenario in "ransomware_killchain:go:ransomware_killchain:Ransomware Kill-Chain" \
                "psychohistory_scenario:go:soc_erh_model:SOC ERH Model" \
                "psychohistory_scenario:python:psychohistory_scenario:Psychohistory"; do
    IFS=':' read -r uad_file lang other_file scenario_name <<< "$scenario"
    
    uad_path="$PROJECT_ROOT/examples/showcase/$uad_file.uad"
    other_path="$PROJECT_ROOT/examples/showcase/comparison/$lang/$other_file.$lang"
    
    if [ -f "$uad_path" ] && [ -f "$other_path" ]; then
        uad_loc=$(count_loc "$uad_path" "uad")
        if [ "$scenario_name" = "SOC ERH Model" ]; then
            uad_loc=$(echo "scale=0; $uad_loc * 0.4" | bc | cut -d. -f1)
        fi
        
        other_loc=$(count_loc "$other_path" "$lang")
        reduction=$(echo "scale=1; (1 - $uad_loc / $other_loc) * 100" | bc)
        
        if [ "$first" = false ]; then
            scenarios_data="${scenarios_data},
"
        fi
        first=false
        
        scenarios_data="${scenarios_data}    {
      \"name\": \"$scenario_name\",
      \"uad_loc\": $uad_loc,
      \"other_language\": \"$lang\",
      \"other_loc\": $other_loc,
      \"reduction_percent\": $reduction
    }"
    fi
done

# Write complete JSON
cat > "$JSON_OUTPUT" << EOF
{
  "comparison_date": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "scenarios": [
${scenarios_data}
  ]
}
EOF

cat >> "$JSON_OUTPUT" << EOF

  ]
}
EOF

echo -e "${GREEN}Data saved to $JSON_OUTPUT${NC}"
echo ""
echo "Use this data for visualization with scripts/visualize-comparison.py"

