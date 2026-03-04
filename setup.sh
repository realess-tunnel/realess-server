#!/bin/bash

# Check if the user has permission to write to /usr/local/bin
if [ "$(id -u)" != "0" ]; then
   echo "Error: This command must be run as root." 1>&2
   exit 1
fi

echo "🚀 Starting installation of Realess Server..."
# 2. Copy binary file to system directory
# Note: ./rlss refers to the file in the extracted temporary directory
cp -f ./rlss /usr/local/bin/rlss

# 3. Set execute permissions
chmod 755 /usr/local/bin/rlss

echo "✅ Installation completed successfully!"
echo "👉 Now use 'rlss' at anywhere."
