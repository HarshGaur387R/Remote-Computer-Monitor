#include "commands.h"
#include <windows.h>
#include <string>	
#include <winuser.h>

void DisplayMessage(const std::wstring &title, const std::wstring &message) {
  MessageBoxW(NULL, message.c_str(), title.c_str(), MB_OK);
}

int LockScreen() { return LockWorkStation(); }
