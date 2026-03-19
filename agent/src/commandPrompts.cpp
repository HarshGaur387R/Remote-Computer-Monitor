#include "commandPrompts.h"
#include "config.h"
#include "server.h"
// #include <iostream>
// #include <string>
// using namespace std;

// void print_help();

void handle_command_prompts(int argc, char *argv[]) {

  const int server_port = AGENT_PORT;
  const int beacon_port = BEACON_PORT;

  startServer(server_port);

  /*
  if (argc <= 1) {
    cout << "No command provided. Use --help for usage.\n";
    return;
  }
  string command = argv[1];
  if (command == "--help") {
    print_help();
  } else if (command == "start") {
    cout << "Starting agent server...\n";
    startServer(server_port);
  } else if (command == "stop") {
    cout << "Stopping agent server...\n";
  } else {
    cout << "Unknown command: " << command << "\n";
    cout << "Use --help to see available commands.\n";
  }
*/
}

/*
void print_help() {
  cout << "Agent Server - Command Line Interface\n";
  cout << "Usage:\n";
  cout << "  agent [command] \n\n";
  cout << "Available commands:\n";
  cout << "  --help           Show this help message\n";
  cout << "  start            Start the agent server in terminal\n";
  cout << "  stop             Stop the agent server\n\n";
  cout << "Examples:\n";
  cout << "  agent --help\n";
  cout << "  agent start\n";
  cout << "  agent start -b\n";
  cout << "  agent stop\n";
}
*/
