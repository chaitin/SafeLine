#pragma once

char* aprintf(const char* fmt, ...)
  __attribute__((format(printf, 1, 2)));
