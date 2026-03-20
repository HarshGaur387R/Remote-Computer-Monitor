#ifndef COMMANDS_H
#define COMMANDS_H
#include <string>

int LockScreen();
int SendNotification(const std::wstring &title, const std::wstring &message);
int DisplayMessage(const std::wstring &title, const std::wstring &message);

#endif
