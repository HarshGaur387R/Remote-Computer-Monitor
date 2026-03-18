#include "commands.h"
#include <windows.h>
#include <string>	
#include <winuser.h>

void DisplayMessage(std::wstring title, std::wstring message) {
  MessageBoxW(NULL, message.c_str(), title.c_str(), MB_OK);
}

int LockScreen() { return LockWorkStation(); }
