#include "server.h"
#include "../include/httplib.h"

void startServer(int port) {

  using namespace httplib;

  Server svr;

  svr.Get("/hi", [](const Request &req, Response &res) {
	res.set_content("Hello World!", "text/plain");
  });
  
  svr.listen("0.0.0.0", port <= 0? port : 8080);
}
