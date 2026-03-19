#include "commandPrompts.h"
#include <iostream>
#include <windows.h>
using namespace std;

int main(int argc, char *argv[]) {

  cout << "\033[1;34mThis is RCM Agent\033[0m\n"<< endl;
  HANDLE hMutex = CreateMutexA(NULL, TRUE, "Global\\RCMA_Server_Mutex");

  if (hMutex == NULL) {
    cerr << "Failed to create mutex.\n";
    return 1;
  }

  if (GetLastError() == ERROR_ALREADY_EXISTS) {
    cout << "RCMA server is already running.\n";
    CloseHandle(hMutex);
    return 1;
  }

  handle_command_prompts(argc, argv);

  CloseHandle(hMutex);
  return 0;
};
