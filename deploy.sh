#!/bin/bash

# Dogs API Deployment Script

echo "🐕 Dogs API Deployment Script"
echo "=============================="

# Check if Docker is installed
if ! command -v docker &> /dev/null
then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

echo "✅ Docker found"

# Build the Docker image
echo "🔨 Building Docker image..."
docker build -t dogs-api .

if [ $? -eq 0 ]; then
    echo "✅ Docker image built successfully"
else
    echo "❌ Docker image build failed"
    exit 1
fi

# Run the container
echo "🚀 Starting container..."
docker run -d -p 8080:8080 -v $(pwd)/data:/data --name dogs-api dogs-api

if [ $? -eq 0 ]; then
    echo "✅ Container started successfully"
    echo "🌐 Access the application at: http://localhost:8080"
    echo "📊 View logs with: docker logs -f dogs-api"
    echo "🛑 Stop the container with: docker stop dogs-api"
else
    echo "❌ Failed to start container"
    exit 1
fi
