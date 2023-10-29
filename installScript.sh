#!/bin/bash

# Variables needed for Script
COVERT_PI_REPO_URL="https://github.com/ArmanSandhu/Covert-Pi-Server.git"
DESKTOP_PATH="$HOME/Desktop"
SERVER_DETAILS_DIR_PATH="$DESKTOP_PATH/CovertPiServerDetails"
SERVER_CONF_FILE_PATH="$DESKTOP_PATH/CovertPiServerDetails/covertpi.conf"
SERVER_CAPTURE_DIR_PATH="$DESKTOP_PATH/CovertPiCaptures"

# Move to the Desktop
cd "$DESKTOP_PATH"

# Get the server files
git clone "$COVERT_PI_REPO_URL"

echo "Repository Downloaded!"

# Build the server
cd "$DESKTOP_PATH/Covert-Pi-Server"
go build -o covert-pi-server

echo "CovertPi Server has been built! Please use the openssl utility to generate the required server and client keys (It is recommended to place these files within the CovertPiServerDetails directory!)"

# Create the Server Details and Capture Directories
mkdir "$SERVER_DETAILS_DIR_PATH"
mkdir "$SERVER_CAPTURE_DIR_PATH"

# Create conf file for modification
echo "hostIP=SERVER_HOST_IP" >> "$SERVER_CONF_FILE_PATH"
echo "hostPort=SERVER_HOST_PORT" >> "$SERVER_CONF_FILE_PATH"
echo "captureDir=$SERVER_CAPTURE_DIR_PATH" >> "$SERVER_CONF_FILE_PATH"
echo "serverKeyFilePath=COVERTPI_SERVER_KEY_PATH" >> "$SERVER_CONF_FILE_PATH"
echo "serverCertFilePath=COVERTPI_SERVER_CERT_PATH" >> "$SERVER_CONF_FILE_PATH"

echo "Please modify the conf file located in CovertPiServerDetails directory!"