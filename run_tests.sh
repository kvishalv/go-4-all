#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Running Go Backend Tests...${NC}"
echo "=================================="

# Run tests with verbose output
go test -v ./...

# Check if tests passed
if [ $? -eq 0 ]; then
    echo -e "\n${GREEN}✅ All tests passed!${NC}"
else
    echo -e "\n${RED}❌ Some tests failed!${NC}"
    exit 1
fi

echo -e "\n${YELLOW}Running tests with coverage...${NC}"
echo "=================================="

# Run tests with coverage
go test -v -cover ./...

echo -e "\n${YELLOW}Generating coverage report...${NC}"
echo "=================================="

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

echo -e "${GREEN}Coverage report generated: coverage.html${NC}"
echo -e "${GREEN}Open coverage.html in your browser to view detailed coverage${NC}"
