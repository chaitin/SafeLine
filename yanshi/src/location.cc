#include "common.hh"
#include "location.hh"

#include <algorithm>
#include <cstdarg>
using namespace std;

LocationFile::LocationFile(const string& filename, const string& data) : filename(filename), data(data)
{
  // data ends with '\n'
  long nlines = count(data.begin(), data.end(), '\n');
  linemap.assign(nlines+1, 0);
  long line = 1;
  for (long i = 0; i < data.size(); i++)
    if (data[i] == '\n')
      linemap[line++] = i+1;
}

void LocationFile::context(const Location& loc) const
{
  long line1, col1, line2, col2;
  locate(loc, line1, col1, line2, col2);
  if (line1 == line2) {
    fputs("  ", stderr);
    FOR(i, linemap[line1], line1+1 < linemap.size() ? linemap[line1+1] : data.size()) {
      if (i == loc.start)
        magenta();
      fputc(data[i], stderr);
      if (i+1 == loc.end)
        sgr0();
    }
  } else {
    bool first = true;
    fputs("  ", stderr);
    FOR(i, linemap[line1], linemap[line1+1]) {
      if (i == loc.start)
        magenta();
      fputc(data[i], stderr);
    }
    if (line2-line1 < 8) {
      FOR(i, linemap[line1+1], linemap[line2]) {
        if (first) { first = false; fputs("  ", stderr); }
        fputc(data[i], stderr);
        if (data[i] == '\n') first = true;
      }
    } else {
      FOR(i, linemap[line1+1], linemap[line1+4]) {
        if (first) { first = false; fputs("  ", stderr); }
        fputc(data[i], stderr);
        if (data[i] == '\n') first = true;
      }
      fputs("  ........\n", stderr);
      FOR(i, linemap[line2-3], linemap[line2]) {
        if (first) { first = false; fputs("  ", stderr); }
        fputc(data[i], stderr);
        if (data[i] == '\n') first = true;
      }
    }
    FOR(i, linemap[line2], line2+1 < linemap.size() ? linemap[line2+1] : data.size()) {
      if (first) { first = false; fputs("  ", stderr); }
      fputc(data[i], stderr);
      if (i+1 == loc.end)
        sgr0();
    }
  }
}

void LocationFile::locate(const Location& loc, long& line1, long& col1, long& line2, long& col2) const
{
  line1 = upper_bound(ALL(linemap), loc.start) - linemap.begin() - 1;
  line2 = upper_bound(ALL(linemap), max(loc.end-1, 0L)) - linemap.begin() - 1;
  col1 = loc.start - linemap[line1];
  col2 = loc.end - linemap[line2];
}

void LocationFile::report_location(const Location& loc) const
{
  long line1, col1, line2, col2;
  locate(loc, line1, col1, line2, col2);
  yellow(2);
  fprintf(stderr, "%s ", filename.c_str());
  cyan(2);
  if (line1 == line2)
    fprintf(stderr, "%ld:%ld-%ld ", line1+1, col1+1, col2);
  else
    fprintf(stderr, "%ld-%ld:%ld-%ld ", line1+1, line2+1, col1+1, col2);
}

void LocationFile::error(const Location& loc, const char* fmt, ...) const
{
  report_location(loc);
  red(2);
  fprintf(stderr, "error ");
  va_list va;
  va_start(va, fmt);
  vfprintf(stderr, fmt, va);
  va_end(va);
  fputs("\n", stderr);
  sgr0(2);
}

void LocationFile::warning(const Location& loc, const char* fmt, ...) const
{
  report_location(loc);
  yellow(2);
  fprintf(stderr, "warning ");
  va_list va;
  va_start(va, fmt);
  vfprintf(stderr, fmt, va);
  va_end(va);
  fputs("\n", stderr);
  sgr0(2);
}
