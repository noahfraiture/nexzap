#!/bin/bash

# Function to zip the tutorials directory and send it to the server
server_url="http://localhost:8080"
zip_file="tutorials.zip"

# Check if the tutorials directory exists
if [ ! -d "tutorials" ]; then
	echo "Error: 'tutorials' directory not found."
	return 1
fi

# Zip the tutorials directory
(cd ./tutorials && zip -r "$zip_file" ./* && mv "$zip_file" ..)
if [ $? -ne 0 ]; then
	echo "Error: Failed to create zip file."
	return 1
fi

# Send the zip file to the server
curl -X POST -F "tutorial_zip=@$zip_file" "$server_url/import"
if [ $? -ne 0 ]; then
	echo "Error: Failed to upload zip file to the server."
	rm "$zip_file" # Clean up the zip file on failure
	return 1
fi

# Clean up the zip file after successful upload
rm "$zip_file"
echo "Tutorials uploaded successfully."

# Example usage:
# upload_tutorials "http://localhost:8080"
