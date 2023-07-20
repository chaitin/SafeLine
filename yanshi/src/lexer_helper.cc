#include "common.hh"

#include <cstdarg>
#include <cstdio>
using namespace std;

char* aprintf(const char* fmt, ...)
{
  va_list va;
  va_start(va, fmt);
  char* r = NULL;
  vasprintf(&r, fmt, va);
  va_end(va);
  return r;
}
