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

# Function to check and remove existing binary
check_existing() {
    INSTALL_PATH="/usr/local/bin/rick"
    if [ -f "$INSTALL_PATH" ]; then
        echo -e "${YELLOW}Found existing installation at: ${GREEN}$INSTALL_PATH${NC}"
        echo -e "${YELLOW}Removing previous version...${NC}"
        sudo rm "$INSTALL_PATH"
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}Previous version removed successfully${NC}"
        else
            echo -e "${RED}Failed to remove previous version${NC}"
            exit 1
        fi
    fi
}

# Function to setup binary location
setup_binary() {
    INSTALL_PATH="/usr/local/bin/rick"
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    SOURCE_BINARY="$SCRIPT_DIR/rick"
    
    # Check if binary exists in the same directory as the script
    if [ ! -f "$SOURCE_BINARY" ]; then
        echo -e "${RED}Error: rick binary not found in script directory${NC}"
        echo -e "${YELLOW}Expected location: ${GREEN}$SOURCE_BINARY${NC}"
        exit 1
    fi
    
    # Check if /usr/local/bin exists and is in PATH
    if [ ! -d "/usr/local/bin" ]; then
        echo -e "${YELLOW}Creating /usr/local/bin directory${NC}"
        sudo mkdir -p /usr/local/bin
    fi

    # Copy binary to /usr/local/bin
    echo -e "${GREEN}Installing rick...${NC}"
    sudo cp "$SOURCE_BINARY" "$INSTALL_PATH"
    sudo chmod +x "$INSTALL_PATH"

    # Verify installation
    if [ ! -x "$INSTALL_PATH" ]; then
        echo -e "${RED}Failed to install rick binary${NC}"
        exit 1
    fi
    
    echo -e "${YELLOW}Rick installed to: ${GREEN}$INSTALL_PATH${NC}"
}

# Function to add or update alias
add_alias() {
    local shell_config="$1"
    
    # Remove existing rick alias if it exists
    sed -i '' '/alias rick=/d' "$shell_config"
    
    # Add the new alias
    echo "alias rick='rick'" >> "$shell_config"
    echo -e "${YELLOW}Rick is getting ready >> $shell_config${NC}"
}

# Print ASCII art
print_rick_ascii

# Check and remove existing binary
check_existing

# Setup binary
setup_binary

# Determine shell configuration file
if [ -n "$ZSH_VERSION" ] || [ "$SHELL" = "/bin/zsh" ] || [ "$SHELL" = "/usr/bin/zsh" ]; then
    SHELL_CONFIG="$HOME/.zshrc"
    SHELL_NAME="zsh"
elif [ -n "$BASH_VERSION" ] || [ "$SHELL" = "/bin/bash" ] || [ "$SHELL" = "/usr/bin/bash" ]; then
    SHELL_CONFIG="$HOME/.bashrc"
    SHELL_NAME="bash"
else
    echo "Unsupported shell. Please add the alias manually."
    exit 1
fi

# Add or update the alias
add_alias "$SHELL_CONFIG"

# Print final instructions with installation details
echo -e "\n${YELLOW}╔════════════════════════════════════════════╗${NC}"
echo -e "${YELLOW}║  ${GREEN}Installation Complete!                    ${YELLOW}║${NC}"
echo -e "${YELLOW}║  ${GREEN}To activate rick, run:                    ${YELLOW}║${NC}"
echo -e "${YELLOW}║  ${PURPLE}source $SHELL_CONFIG                  ${YELLOW}║${NC}"
echo -e "${YELLOW}╚════════════════════════════════════════════╝${NC}\n"