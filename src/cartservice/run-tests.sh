#!/bin/bash

# Script to run cartservice tests with JUnit output using Docker

echo "Running cartservice tests with JUnit output..."

# Build and run tests in Docker
docker run --rm \
  -v "$(pwd)":/workspace \
  -w /workspace \
  mcr.microsoft.com/dotnet/sdk:9.0 \
  bash -c "
    echo 'Restoring packages...'
    dotnet restore
    
    echo 'Running tests with JUnit logger...'
    dotnet test --logger:junit --results-directory ./TestResults
    
    echo 'Tests completed. Results are in ./TestResults/TestResults.xml'
    ls -la ./TestResults/
  "

echo "JUnit test results available in ./TestResults/TestResults.xml"
