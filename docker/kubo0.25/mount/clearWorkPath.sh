#!/bin/bash

set -ex
# Define the directory
DIR="/home/Cherry/data/p2pWork"

# Define the files to keep
KEEP_FILES=("1080left.mp4" "4k30left.mp4" "left.mp4" "1080right.mp4" "4k30right.mp4" "right.mp4")

# Change to the target directory
cd "$DIR"

# Check if the directory change was successful
if [ $? -ne 0 ]; then
    echo "Error: Could not change to directory $DIR. Script aborted."
    exit 1
fi

# Loop through all files in the directory
for file in *; do
    # Check if the file is not in the keep list
    if [[ ! " ${KEEP_FILES[@]} " =~ " $file " ]]; then
        # Delete the file
        rm -f "$file"
    fi
done

echo "Cleanup complete."
