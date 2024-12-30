#!/bin/bash

# Function to extract version number from git tags
get_latest_version() {
    git tag -l "v*" | sort -V | tail -n1 | sed 's/v//'
}

# Function to increment version
increment_version() {
    local version=$1
    local position=$2
    local IFS=.
    local array=($version)
    
    case $position in
        major)
            ((array[0]++))
            array[1]=0
            array[2]=0
            ;;
        minor)
            ((array[1]++))
            array[2]=0
            ;;
        patch)
            ((array[2]++))
            ;;
    esac
    
    echo "${array[0]}.${array[1]}.${array[2]}"
}

# Get the current version
current_version=$(get_latest_version)
if [ -z "$current_version" ]; then
    current_version="0.1.0"
fi

# Determine which version number to increment
echo "Current version: v$current_version"
echo "Select version increment type:"
echo "1) major (x.0.0)"
echo "2) minor (0.x.0)"
echo "3) patch (0.0.x)"
read -p "Enter choice [1-3]: " choice

case $choice in
    1) new_version=$(increment_version $current_version "major");;
    2) new_version=$(increment_version $current_version "minor");;
    3) new_version=$(increment_version $current_version "patch");;
    *) echo "Invalid choice"; exit 1;;
esac

# Confirm the new version
echo "New version will be: v$new_version"
read -p "Continue? [y/N] " confirm
if [[ $confirm != [yY] ]]; then
    echo "Aborted"
    exit 1
fi

# Make sure all changes are committed
if [ -n "$(git status --porcelain)" ]; then
    echo "There are uncommitted changes. Please commit them first."
    exit 1
fi

# Create and push the new tag
git tag "v$new_version"
git push origin "v$new_version"

echo "Successfully created and pushed tag v$new_version" 