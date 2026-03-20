#include "commands.h"
#include "wintoastlib.h"
#include <string>
#include <windows.h>
#include <winuser.h>

using namespace WinToastLib;

class CustomHandler : public IWinToastHandler {
public:
  void toastActivated() const {
    std::wcout << L"Toast activated!" << std::endl;
  }

  void toastActivated(int actionIndex) const {
    std::wcout << L"Toast activated with" << actionIndex << std::endl;
  }

  void toastActivated(std::wstring response) const {
    std::wcout << L"Toast response: " << response << std::endl;
  }

  void toastDismissed(WinToastDismissalReason state) const {
    switch (state) {
    case UserCanceled:
      std::wcout << L"The user dismissed this toast" << std::endl;
      break;
    case TimedOut:
      std::wcout << L"The toast has timed out" << std::endl;
      break;
    case ApplicationHidden:
      std::wcout << L"The application hid the toast using ToastNotifier.hide()"
                 << std::endl;
      break;
    default:
      std::wcout << L"Toast not activated" << std::endl;
      break;
    }
  }

  void toastFailed() const {
    std::wcout << L"Error showing current toast" << std::endl;
  }
};

int DisplayMessage(const std::wstring &title, const std::wstring &message) {
  return MessageBoxW(NULL, message.c_str(), title.c_str(), MB_OK);
}

int SendNotification(const std::wstring &title, const std::wstring &message) {

  WinToast::instance()->setAppName(L"RCMA");
  WinToast::instance()->setAppUserModelId(L"RCMAid");
  WinToast::instance()->createShortcut();

  if (!WinToast::instance()->initialize()) {
    std::wcerr << L"Error initializing WinToast" << std::endl;
    return 0;
  }

  WinToastTemplate templ(WinToastTemplate::Text02);
  templ.setTextField(title, WinToastTemplate::FirstLine);
  templ.setTextField(message, WinToastTemplate::SecondLine);

  CustomHandler *handler = new CustomHandler();
  WinToast::instance()->showToast(templ, handler);
  return 1;
}

int LockScreen() { return LockWorkStation(); }
