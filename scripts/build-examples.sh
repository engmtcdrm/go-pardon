#!/bin/bash

workspaceFolder=${WORKSPACE_FOLDER:-$(pwd)}

echo "Size before build:"
ls -la "$workspaceFolder/examples" |grep examples
ls -lh "$workspaceFolder/examples" |grep examples

echo ""
echo "Size after build:"
cd "$workspaceFolder/examples"
go build --ldflags "-s -w" -o examples
ls -la |grep examples
ls -lh |grep examples
