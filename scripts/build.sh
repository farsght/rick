#!/bin/zsh


# Define colors using ANSI escape codes
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
PURPLE='\033[0;35m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to print ASCII art
print_rick_ascii() {
    echo -e "${CYAN}
    ██████╗ ██╗ ██████╗██╗  ██╗
    ██╔══██╗██║██╔════╝██║ ██╔╝
    ██████╔╝██║██║     █████╔╝ 
    ██╔══██╗██║██║     ██╔═██╗ 
    ██║  ██║██║╚██████╗██║  ██╗
    ╚═╝  ╚═╝╚═╝ ╚═════╝╚═╝  ╚═╝${NC}
    "
}

# Function to increment build number
increment_build() {
    VERSION_FILE="cmd/version.go"
    if [ ! -f "$VERSION_FILE" ]; then
        echo -e "${YELLOW}Creating version.go file in cmd directory...${NC}"
        mkdir -p cmd
        echo 'package cmd

var (
    Version string = "1.0.0"
    BuildNum string = "0"
)' > "$VERSION_FILE"
    fi
    
    # Read current build number
    current_build=$(grep 'BuildNum.*=.*"[0-9]*"' "$VERSION_FILE" | grep -o '[0-9]*')
    if [ -z "$current_build" ]; then
        current_build=0
    fi
    
    # Increment build number
    new_build=$((current_build + 1))
    
    # Update the file with new build number
    sed -i '' "s/BuildNum.*=.*\"[0-9]*\"/BuildNum string = \"$new_build\"/" "$VERSION_FILE"
    
    echo -e "${GREEN}Build number incremented to: ${YELLOW}$new_build${NC}"
    export CURRENT_BUILD=$new_build
}

# Build the binary with embedded variables and version
build_binary() {
    echo -e "${GREEN}Building rick with embedded configuration...${NC}"
    
    # Increment build number before building
    increment_build
    
    VERSION="1.0.0"
    
    echo -e "${YELLOW}Building version ${VERSION} (Build ${CURRENT_BUILD})${NC}"
    
    go build -ldflags "
        -X 'github.com/farsght/rick/cmd.OpenAIAPIKey=$API_KEY' 
        -X 'github.com/farsght/rick/cmd.OpenAIModel=$MODEL' 
        -X 'github.com/farsght/rick/cmd.Version=${VERSION}' 
        -X 'github.com/farsght/rick/cmd.BuildNum=${CURRENT_BUILD}'
    " -o rick
    
    if [ $? -ne 0 ]; then
        echo -e "${RED}Failed to build rick${NC}"
        exit 1
    fi
    echo -e "${GREEN}Binary built successfully (Version ${VERSION}, Build #${YELLOW}${CURRENT_BUILD}${GREEN})${NC}"
}

# Function to setup binary location
setup_binary() {
    # Check if /usr/local/bin exists and is in PATH
    if [ ! -d "/usr/local/bin" ]; then
        echo -e "${YELLOW}Creating /usr/local/bin directory${NC}"
        sudo mkdir -p /usr/local/bin
    fi

    # Copy binary to /usr/local/bin
    echo -e "${GREEN}Installing rick binary...${NC}"
    sudo cp rick /usr/local/bin/
    sudo chmod +x /usr/local/bin/rick

    # Verify installation
    if [ ! -x "/usr/local/bin/rick" ]; then
        echo -e "${RED}Failed to install rick binary${NC}"
        exit 1
    fi
    
    echo -e "${YELLOW}Rick installed to: ${GREEN}/usr/local/bin/rick${NC}"
}

# Print ASCII art
print_rick_ascii

# Build the binary
build_binary

# Setup binary
setup_binary

# Print completion message
echo -e "\n${GREEN}Installation complete! ${PURPLE}rick${GREEN} is ready to use.${NC}\n"