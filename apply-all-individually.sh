#!/bin/bash

# Configuration
COMPONENTS_DIR="./config/components"
COMMAND_PREFIX="go run ./cmd/root.go component apply -l ./config/components -c"

# Tracking files
PROCESSED_FILE="processed_components.txt"
FAILED_FILE="failed_components.txt"
SKIPPED_FILE="skipped_components.txt"

# Function to extract component names from YAML files
get_component_names() {
    find "$COMPONENTS_DIR" -name "component-*.yaml" -type f | \
    sed 's|.*/component-\(.*\)\.yaml|\1|' | \
    sort
}

# Ensure tracking files exist
touch "$PROCESSED_FILE"
touch "$FAILED_FILE"
touch "$SKIPPED_FILE"

# Function to check if component was already processed
is_component_processed() {
    local component="$1"
    grep -Fxq "$component" "$PROCESSED_FILE" || grep -Fxq "$component" "$FAILED_FILE"
}

# Function to process a single component (runs serially)
process_component() {
    local component="$1"
    local component_file="$COMPONENTS_DIR/component-$component.yaml"

    # Verify the file still exists
    if [[ ! -f "$component_file" ]]; then
        echo "Warning: Component file not found: $component_file"
        echo "$component" >> "$SKIPPED_FILE"
        return 1
    fi

    echo "Processing component: $component"
    echo "  File: $component_file"
    echo "  Command: $COMMAND_PREFIX $component"

    # Execute the command
    if $COMMAND_PREFIX "$component"; then
        echo "  ✓ Successfully processed component: $component"
        echo "$component" >> "$PROCESSED_FILE"
        return 0
    else
        echo "  ✗ Failed to process component: $component"
        echo "$component" >> "$FAILED_FILE"
        return 1
    fi
}

# Variables for statistics
total_components=0
total_new_processed=0
total_new_failed=0
total_prev_processed=0
total_prev_failed=0
total_skipped=0

# Cleanup function to print statistics
cleanup() {
    echo ""
    echo "========== Component Apply Statistics =========="
    echo "Total components found: $total_components"
    echo "New components processed successfully: $total_new_processed"
    echo "Previously processed components (skipped): $total_prev_processed"
    echo "New components failed: $total_new_failed"
    echo "Previously failed components (skipped): $total_prev_failed"
    echo "Components skipped (file not found): $total_skipped"
    echo "=============================================="
    echo "Script completed"
    exit 0
}

# Trap signals for graceful shutdown
trap cleanup SIGTERM SIGINT

# Validate components directory exists
if [[ ! -d "$COMPONENTS_DIR" ]]; then
    echo "Error: Components directory '$COMPONENTS_DIR' does not exist"
    exit 1
fi

# Get all component names
component_names=$(get_component_names)

if [[ -z "$component_names" ]]; then
    echo "No component-*.yaml files found in $COMPONENTS_DIR"
    exit 1
fi

echo "Found component files:"
echo "$component_names" | sed 's/^/  - component-/' | sed 's/$/.yaml/'
echo ""

# First pass: identify components to process and skip already processed ones
components_to_process=()
for component in $component_names; do
    ((total_components++))

    component_file="$COMPONENTS_DIR/component-$component.yaml"

    # Verify the file exists
    if [[ ! -f "$component_file" ]]; then
        echo "Warning: Component file not found: $component_file"
        echo "$component" >> "$SKIPPED_FILE"
        ((total_skipped++))
        continue
    fi

    # Check if component was already processed
    if is_component_processed "$component"; then
        echo "Skipping component: $component (already processed)"
        if grep -Fxq "$component" "$PROCESSED_FILE"; then
            ((total_prev_processed++))
        elif grep -Fxq "$component" "$FAILED_FILE"; then
            ((total_prev_failed++))
        fi
        continue
    fi

    components_to_process+=("$component")
done

# If no components to process, exit
if [[ ${#components_to_process[@]} -eq 0 ]]; then
    echo "No new components to process."
    cleanup
fi

echo "Starting serial processing of ${#components_to_process[@]} components..."
echo ""

# Process components one by one (serially)
for component in "${components_to_process[@]}"; do
    if process_component "$component"; then
        ((total_new_processed++))
    else
        ((total_new_failed++))
    fi
    echo "" # Add blank line between components for readability
done

echo "All components processed serially."

cleanup