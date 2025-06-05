#!/bin/bash

# Configuration
COMPONENTS_DIR="./config/components"
MAX_RETRIES=3
BASE_DELAY=30  # Base delay in seconds
MAX_DELAY=300  # Maximum delay in seconds

# Set command prefix based on environment
if [[ -n "$GITHUB_ACTIONS" ]]; then
    COMMAND_PREFIX="./ofc component compute -a -c"
else
    COMMAND_PREFIX="go run ./cmd/root.go component compute -a -c"
fi

# Tracking files
PROCESSED_FILE="processed_components.txt"
FAILED_FILE="failed_components.txt"
SKIPPED_FILE="skipped_components.txt"

# Check if single component argument is provided
SINGLE_COMPONENT="$1"

# Function to check GitHub API rate limit
check_rate_limit() {
    local remaining=$(curl -s \
        -H "Authorization: token $GITHUB_TOKEN" \
        -H "Accept: application/vnd.github.v3+json" \
        "https://api.github.com/rate_limit" | \
        jq -r '.rate.remaining // 0')

    echo "$remaining"
}

# Function to wait for rate limit reset
wait_for_rate_limit_reset() {
    local reset_time=$(curl -s \
        -H "Authorization: token $GITHUB_TOKEN" \
        -H "Accept: application/vnd.github.v3+json" \
        "https://api.github.com/rate_limit" | \
        jq -r '.rate.reset // 0')

    if [[ "$reset_time" != "0" && "$reset_time" != "null" ]]; then
        local current_time=$(date +%s)
        local wait_time=$((reset_time - current_time + 60))  # Add 1 minute buffer

        if [[ $wait_time -gt 0 && $wait_time -lt 3600 ]]; then  # Don't wait more than 1 hour
            echo "Rate limit exceeded. Waiting $wait_time seconds for reset..."
            sleep "$wait_time"
            return 0
        fi
    fi

    # Fallback: wait 5 minutes
    echo "Rate limit exceeded. Waiting 5 minutes..."
    sleep 300
}

# Function to calculate exponential backoff delay
calculate_delay() {
    local attempt=$1
    local delay=$((BASE_DELAY * (2 ** (attempt - 1))))

    if [[ $delay -gt $MAX_DELAY ]]; then
        delay=$MAX_DELAY
    fi

    echo "$delay"
}

# Enhanced function to process a single component with retry logic
process_component() {
    local component="$1"
    local component_file="$COMPONENTS_DIR/component-$component.yaml"
    local attempt=1

    # Verify the file still exists
    if [[ ! -f "$component_file" ]]; then
        echo "Warning: Component file not found: $component_file"
        echo "$component" >> "$SKIPPED_FILE"
        return 1
    fi

    while [[ $attempt -le $MAX_RETRIES ]]; do
        echo "Processing component: $component (attempt $attempt/$MAX_RETRIES)"
        echo "  File: $component_file"
        echo "  Command: $COMMAND_PREFIX $component"

        # Check rate limit before processing
        local remaining=$(check_rate_limit)
        if [[ "$remaining" != "null" && "$remaining" != "0" && $remaining -lt 100 ]]; then
            echo "  Rate limit low ($remaining remaining). Waiting for reset..."
            wait_for_rate_limit_reset
        fi

        # Execute the command and capture both stdout and stderr
        local output
        local exit_code

        output=$($COMMAND_PREFIX "$component" 2>&1)
        exit_code=$?

        if [[ $exit_code -eq 0 ]]; then
            echo "  ✓ Successfully processed component: $component"
            echo "$component" >> "$PROCESSED_FILE"

            # Adaptive delay based on success - shorter for single component
            if [[ -z "$SINGLE_COMPONENT" ]]; then
                echo "  Waiting 10 seconds before next component..."
                sleep 10
            fi
            return 0
        else
            echo "  ✗ Failed to process component: $component (attempt $attempt)"
            echo "  Error output: $output"

            # Check if it's a rate limit error
            if echo "$output" | grep -qi "rate.limit\|403\|too.many.requests"; then
                echo "  Rate limit detected. Waiting for reset..."
                wait_for_rate_limit_reset
                ((attempt++))
                continue
            fi

            # Check if it's a temporary error that should be retried
            if echo "$output" | grep -qi "timeout\|connection.refused\|service.unavailable\|502\|503\|504"; then
                if [[ $attempt -lt $MAX_RETRIES ]]; then
                    local delay=$(calculate_delay $attempt)
                    echo "  Temporary error detected. Retrying in $delay seconds..."
                    sleep "$delay"
                    ((attempt++))
                    continue
                fi
            fi

            # For other errors or max retries reached, fail immediately
            echo "$component" >> "$FAILED_FILE"
            return 1
        fi
    done

    echo "  ✗ Max retries exceeded for component: $component"
    echo "$component" >> "$FAILED_FILE"
    return 1
}

# Function to extract component names from YAML files
get_component_names() {
    find "$COMPONENTS_DIR" -name "component-*.yaml" -type f | \
    sed 's|.*/component-\(.*\)\.yaml|\1|' | \
    sort
}

# Function to validate if component exists
component_exists() {
    local component="$1"
    local component_file="$COMPONENTS_DIR/component-$component.yaml"
    [[ -f "$component_file" ]]
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
    if [[ -n "$SINGLE_COMPONENT" ]]; then
        echo "========== Single Component Statistics =========="
        echo "Component: $SINGLE_COMPONENT"
        if [[ $total_new_processed -eq 1 ]]; then
            echo "Status: Successfully processed"
        elif [[ $total_new_failed -eq 1 ]]; then
            echo "Status: Failed to process"
        elif [[ $total_prev_processed -eq 1 ]]; then
            echo "Status: Previously processed (skipped)"
        elif [[ $total_prev_failed -eq 1 ]]; then
            echo "Status: Previously failed (skipped)"
        elif [[ $total_skipped -eq 1 ]]; then
            echo "Status: Component file not found"
        fi
    else
        echo "========== Component Apply Statistics =========="
        echo "Total components found: $total_components"
        echo "New components processed successfully: $total_new_processed"
        echo "Previously processed components (skipped): $total_prev_processed"
        echo "New components failed: $total_new_failed"
        echo "Previously failed components (skipped): $total_prev_failed"
        echo "Components skipped (file not found): $total_skipped"
    fi
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

# Handle single component processing
if [[ -n "$SINGLE_COMPONENT" ]]; then
    echo "Single component mode: processing '$SINGLE_COMPONENT'"

    # Validate component exists
    if ! component_exists "$SINGLE_COMPONENT"; then
        echo "Error: Component '$SINGLE_COMPONENT' not found"
        echo "  Expected file: $COMPONENTS_DIR/component-$SINGLE_COMPONENT.yaml"
        exit 1
    fi

    total_components=1

    # Check if component was already processed
    if is_component_processed "$SINGLE_COMPONENT"; then
        echo "Component '$SINGLE_COMPONENT' was already processed"
        if grep -Fxq "$SINGLE_COMPONENT" "$PROCESSED_FILE"; then
            echo "  Status: Previously processed successfully"
            total_prev_processed=1
        elif grep -Fxq "$SINGLE_COMPONENT" "$FAILED_FILE"; then
            echo "  Status: Previously failed"
            total_prev_failed=1
        fi
        cleanup
    fi

    # Process the single component
    echo ""
    if process_component "$SINGLE_COMPONENT"; then
        total_new_processed=1
    else
        total_new_failed=1
    fi

    cleanup
fi

# Batch processing mode
echo "Batch processing mode: processing all components"

# Get all component names
component_names=$(get_component_names)

if [[ -z "$component_names" ]]; then
    echo "No component-*.yaml files found in $COMPONENTS_DIR"
    exit 1
fi

echo "Found component files:"
echo "$component_names" | sed 's/^/  - component-/' | sed 's/$/.yaml/'
echo ""

# Check initial rate limit status
initial_remaining=$(check_rate_limit)
echo "Initial GitHub API rate limit remaining: $initial_remaining"
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

echo "Starting processing of ${#components_to_process[@]} components with rate limit handling..."
echo ""

# Process components one by one with enhanced error handling
for component in "${components_to_process[@]}"; do
    if process_component "$component"; then
        ((total_new_processed++))
    else
        ((total_new_failed++))
    fi
    echo "" # Add blank line between components for readability
done

echo "All components processed."

cleanup