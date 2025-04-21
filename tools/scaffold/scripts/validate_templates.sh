#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Directory containing templates
TEMPLATE_DIR="../templates"

# Counter for validation results
VALID_COUNT=0
INVALID_COUNT=0
TOTAL_COUNT=0

echo -e "${YELLOW}Starting template validation...${NC}\n"

# Function to validate Go template syntax
validate_go_template() {
    local file=$1
    if [[ $file == *.go.tmpl ]]; then
        # Remove template syntax temporarily
        temp_file=$(mktemp)
        sed 's/{{.*}}/dummy/g' "$file" > "$temp_file"
        
        # Check Go syntax
        if ! go fmt "$temp_file" > /dev/null 2>&1; then
            echo -e "${RED}✗ Invalid Go syntax in $file${NC}"
            return 1
        fi
        rm "$temp_file"
    fi
    return 0
}

# Function to validate YAML template syntax
validate_yaml_template() {
    local file=$1
    if [[ $file == *.yml.tmpl || $file == *.yaml.tmpl ]]; then
        # Remove template syntax temporarily
        temp_file=$(mktemp)
        sed 's/{{.*}}/dummy/g' "$file" > "$temp_file"
        
        # Check YAML syntax (requires yq)
        if ! yq eval "$temp_file" > /dev/null 2>&1; then
            echo -e "${RED}✗ Invalid YAML syntax in $file${NC}"
            return 1
        fi
        rm "$temp_file"
    fi
    return 0
}

# Validate each template file
find "$TEMPLATE_DIR" -type f -name "*.tmpl" | while read -r template; do
    TOTAL_COUNT=$((TOTAL_COUNT + 1))
    filename=$(basename "$template")
    
    # Check if file exists and is readable
    if [[ ! -r "$template" ]]; then
        echo -e "${RED}✗ Cannot read template: $filename${NC}"
        INVALID_COUNT=$((INVALID_COUNT + 1))
        continue
    fi
    
    # Check for empty files
    if [[ ! -s "$template" ]]; then
        echo -e "${RED}✗ Empty template: $filename${NC}"
        INVALID_COUNT=$((INVALID_COUNT + 1))
        continue
    }
    
    # Validate template syntax
    if ! grep -q "{{" "$template" && ! grep -q "}}" "$template"; then
        echo -e "${YELLOW}⚠ Warning: No template syntax found in $filename${NC}"
    fi
    
    # Validate based on file type
    validate_go_template "$template"
    validate_yaml_template "$template"
    
    if [[ $? -eq 0 ]]; then
        echo -e "${GREEN}✓ Valid template: $filename${NC}"
        VALID_COUNT=$((VALID_COUNT + 1))
    else
        INVALID_COUNT=$((INVALID_COUNT + 1))
    fi
done

# Print summary
echo -e "\n${YELLOW}Validation Summary:${NC}"
echo -e "Total templates: $TOTAL_COUNT"
echo -e "${GREEN}Valid templates: $VALID_COUNT${NC}"
echo -e "${RED}Invalid templates: $INVALID_COUNT${NC}"

# Exit with error if any templates are invalid
[[ $INVALID_COUNT -eq 0 ]] 