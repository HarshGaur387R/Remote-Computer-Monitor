#include "server.h"
#include "../include/httplib.h"
#include "commands.h"
#include <iostream>

void startServer(int port) {
  using namespace std;
  using namespace httplib;

  Server svr;

  svr.Get("/lock", [](const Request &req, Response &res) {
    if (LockScreen() > 0) {
      cout << "Screen Locked succesfully." << endl;
      res.set_content("Screen Locked", "text/plain");
    } else {
      cerr << "Failed to lock the screen." << endl;
      res.set_content("Failed to lock Screen", "text/plain");
    }
  });

  svr.Get("/message", [](const Request &req, Response &res) {
    DisplayMessage(L"Message from admin", L"You have only 20 minutes left.");
    res.set_content("message diliverd", "text/plain");
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
