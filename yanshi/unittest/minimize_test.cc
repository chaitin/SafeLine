#include "fsa.hh"
#include "unittest/unittest_helper.hh"

#include <algorithm>
#include <iostream>
#include <type_traits>
#include <unistd.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>
using namespace std;

const char test[] =
"4 4 1\n"
"3  \n"
"0 0 1\n"
"0 1 2\n"
"1 0 3\n"
"2 0 3\n"
;

int main(int argc, char *argv[])
{
  if (argc == 1) {
    char filename[] = "/tmp/XXXXXX";
    int fd = mkstemp(filename);
    write(fd, test, sizeof test-1);
    close(fd);
    freopen(filename, "r", stdin);
    unlink(filename);
  }

  auto relate = [](const vector<long>&) {};
  Fsa fsa = read_dfa().minimize(relate);
  print_fsa(fsa);

  if (argc == 1)
    return fsa.n() == 3 ? 0 : 1;
}
