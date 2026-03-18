cmake_minimum_required(VERSION 3.10)

# Project name and language
project(agent LANGUAGES CXX)

# Set C++ standard
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Define source files
set(SOURCES
    src/main.cpp 
    src/server.cpp
    src/commands.cpp
)

# Define include directories
include_directories(include)

# Create executable
add_executable(agent ${SOURCES})

# Output directory for build
set_target_properties(agent PROPERTIES
    RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/bin
)

