#include "config.h"
#include "server.h"
#include <iostream>
using namespace std;

int main() {
  const int server_port = AGENT_PORT;
  const int beacon_port = BEACON_PORT;

  cout << "This is Agent server!" << endl;

  startServer(server_port);

  return 0;
};
