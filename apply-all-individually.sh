#!/bin/bash

# Configuration
COMPONENTS_DIR="./config/components"
COMMAND_PREFIX="go run ./cmd/root.go component apply -l ./config/components -c"

# Tracking files
PROCESSED_FILE="processed_components.txt"
FAILED_FILE="failed_components.txt"
SKIPPED_FILE="skipped_components.txt"

# Maximum number of parallel jobs (adjust based on your system)
MAX_JOBS=10

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

# Function to process a single component (runs in background)
process_component() {
    local component="$1"
    local component_file="$COMPONENTS_DIR/component-$component.yaml"
    local temp_dir="./temp_results"

    # Create temp directory if it doesn't exist
    mkdir -p "$temp_dir"

    # Verify the file still exists
    if [[ ! -f "$component_file" ]]; then
        echo "Warning: Component file not found: $component_file"
        echo "$component" > "$temp_dir/skipped_$component"
        return 1
    fi

    echo "Processing component: $component"
    echo "  File: $component_file"
    echo "  Command: $COMMAND_PREFIX $component"

    # Execute the command
    if $COMMAND_PREFIX "$component"; then
        echo "  ✓ Successfully processed component: $component"
        echo "$component" > "$temp_dir/success_$component"
        return 0
    else
        echo "  ✗ Failed to process component: $component"
        echo "$component" > "$temp_dir/failed_$component"
        return 1
    fi
}

# Function to collect results from temp files
collect_results() {
    local temp_dir="./temp_results"

    if [[ -d "$temp_dir" ]]; then
        # Collect successful components
        for file in "$temp_dir"/success_*; do
            if [[ -f "$file" ]]; then
                cat "$file" >> "$PROCESSED_FILE"
            fi
        done

        # Collect failed components
        for file in "$temp_dir"/failed_*; do
            if [[ -f "$file" ]]; then
                cat "$file" >> "$FAILED_FILE"
            fi
        done

        # Collect skipped components
        for file in "$temp_dir"/skipped_*; do
            if [[ -f "$file" ]]; then
                cat "$file" >> "$SKIPPED_FILE"
            fi
        done

        # Clean up temp directory
        rm -rf "$temp_dir"
    fi
}

# Variables for statistics
total_components=0
total_new_processed=0
total_new_failed=0
total_prev_processed=0
total_prev_failed=0
total_skipped=0

# Array to store background process PIDs
declare -a pids=()

# Cleanup function to print statistics
cleanup() {
    echo ""
    echo "Terminating running processes..."

    # Kill all background processes
    for pid in "${pids[@]}"; do
        if kill -0 "$pid" 2>/dev/null; then
            kill "$pid" 2>/dev/null
        fi
    done

    # Wait a moment for processes to terminate gracefully
    sleep 2

    # Force kill any remaining processes
    for pid in "${pids[@]}"; do
        if kill -0 "$pid" 2>/dev/null; then
            kill -9 "$pid" 2>/dev/null
        fi
    done

    # Collect any remaining results
    collect_results

    echo ""
    echo "========== Component Apply Statistics =========="
    echo "Total components found: $total_components"
    echo "New components processed successfully: $total_new_processed"
    echo "Previously processed components (skipped): $total_prev_processed"
    echo "New components failed: $total_new_failed"
    echo "Previously failed components (skipped): $total_prev_failed"
    echo "Components skipped (file not found): $total_skipped"
    echo "=============================================="
    echo "Script terminated"
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

echo "Starting parallel processing of ${#components_to_process[@]} components..."
echo "Maximum concurrent jobs: $MAX_JOBS"
echo ""

# Process components in parallel with job control
job_count=0
for component in "${components_to_process[@]}"; do
    # Wait if we've reached the maximum number of jobs
    if [[ $job_count -ge $MAX_JOBS ]]; then
        # Wait for any job to finish
        wait -n
        ((job_count--))
    fi

    # Start the component processing
    process_component "$component" &
    pids+=($!)
    ((job_count++))
done

echo "All component processes started (with job limit of $MAX_JOBS)."
echo "Waiting for all processes to complete..."
echo ""

# Wait for all remaining background processes to complete
for pid in "${pids[@]}"; do
    if wait "$pid"; then
        ((total_new_processed++))
    else
        ((total_new_failed++))
    fi
done

echo ""
echo "All parallel processes completed."

# Collect results from temporary files
collect_results

cleanup