#!/bin/bash

workspaceFolder=${WORKSPACE_FOLDER:-$(pwd)}
d=$(dirname "$0")
echo "Running $d in workspace folder: $workspaceFolder"

goDirs=$(find . -path "./examples" -prune -o -name '*.go' -exec dirname {} \; | sort -u)

readarray -t dirs <<< "$goDirs"
for dir in "${dirs[@]}"; do
    relDir=$(realpath --relative-to="$workspaceFolder" "$dir")

    # if [[ $relDir == "." ]]; then
    #     relDir=$(basename "$workspaceFolder")
    # fi

    echo "Generating documentation for '$relDir'..."
    gomarkdoc -v "$dir" --config "$workspaceFolder/.gomarkdoc.yml" -o "$dir/README.md"
done

# echo "Generating documentation for the project at $workspaceFolder..."
# gomarkdoc -v "$workspaceFolder" --config "$workspaceFolder/.gomarkdoc.yml" -o "$workspaceFolder/README.md"

# echo ""
# echo "Generating documentation for the 'grapheme' package..."
# gomarkdoc -v "$workspaceFolder/grapheme" --config "$workspaceFolder/.gomarkdoc.yml" -o "$workspaceFolder/grapheme/README.md"

# echo ""
# echo "Generating documentation for the 'internal/testutils' package..."
# gomarkdoc -v "$workspaceFolder/internal/testutils" --config "$workspaceFolder/.gomarkdoc.yml" -o "$workspaceFolder/internal/testutils/README.md"
