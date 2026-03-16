#include "server.h"
#include "../include/httplib.h"
#include <iostream>

using namespace std;

void startServer(int port) {

  using namespace httplib;

  Server svr;

  svr.Get("/hi", [](const Request &req, Response &res) {
    res.set_content("Hello World!", "text/plain");
  });

  int corrected_port = port > 0 ? port : 8080;

  cout << "Spining up server .." << endl;
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
