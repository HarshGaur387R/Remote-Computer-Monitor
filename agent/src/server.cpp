#include "server.h"
#include "../include/httplib.h"
#include "commands.h"
#include <iostream>

// Helper: convert UTF-8 std::string to UTF-16 std::wstring
std::wstring utf8_to_wstring(const std::string &str) {
  if (str.empty())
    return std::wstring();

  int size_needed =
      MultiByteToWideChar(CP_UTF8, 0, str.c_str(), (int)str.size(), nullptr, 0);

  std::wstring wstr(size_needed, 0);
  MultiByteToWideChar(CP_UTF8, 0, str.c_str(), (int)str.size(), &wstr[0],
                      size_needed);
  return wstr;
}

void startServer(int port) {
  using namespace std;
  using namespace httplib;

  Server svr;

  svr.Get("/lock", [](const Request &req, Response &res) {
    if (LockScreen() > 0) {
      cout << "Screen Locked successfully." << endl;
      res.set_content("Screen Locked", "text/plain");
    } else {
      cerr << "Failed to lock the screen." << endl;
      res.set_content("Failed to lock Screen", "text/plain");
    }
  });

  svr.Get("/notification", [](const Request &req, Response &res) {
    string body = req.get_param_value("body");
    string title = req.get_param_value("title");

    // Convert safely from UTF-8 to UTF-16
    wstring wbody = utf8_to_wstring(body);
    wstring wtitle = utf8_to_wstring(title);

    SendNotification(wtitle, wbody);
    res.set_content("Message delivered", "text/plain");
  });

  svr.Get("/message", [](const Request &req, Response &res) {
    std::string body = req.get_param_value("body");
    std::string title = req.get_param_value("title");

    std::wstring wbody = utf8_to_wstring(body);
    std::wstring wtitle = utf8_to_wstring(title);

    // Launch DisplayMessage in its own thread
    std::thread([wtitle, wbody]() {
      DisplayMessage(wtitle, wbody); // blocking, but only in this thread
    }).detach();

    // Respond immediately
    res.set_content("Message delivered", "text/plain");
  });

  int corrected_port = port > 0 ? port : 8080;

  cout << "Server will listen at " << corrected_port << endl;
  if (!svr.listen("0.0.0.0", corrected_port)) {
    cout << "PORT " << corrected_port << " is busy, trying different port...\n";

    int s_a_p = svr.bind_to_any_port("0.0.0.0"); // System Allocated Port

    if (s_a_p > 0) {
      cout << "Hooray! Server is listening at: " << s_a_p << endl;
      svr.listen_after_bind();
    } else {
      cerr << "Failed to run local agent server." << endl;
    }
  }
}

/* curl "http://localhost:4567/message?body=Hello world&title=hello" ` -Method
 * GET ` -Headers @{ "Content-Type" = "application/x-www-form-urlencoded" } */
