#pragma once
#ifndef _GNU_SOURCE
# define _GNU_SOURCE
#endif
#include <assert.h>
#include <map>
#include <stdarg.h>
#include <stdint.h>
#include <stdio.h>
#include <type_traits>
#include <vector>
using std::map;
using std::vector;

typedef int8_t i8;
typedef int16_t i16;
typedef int32_t i32;
typedef int64_t i64;
typedef uint8_t u8;
typedef uint16_t u16;
typedef uint32_t u32;
typedef uint64_t u64;
typedef unsigned long ulong;

#ifdef __APPLE__
#include <crt_externs.h>
extern char*** _NSGetArgv(void);
#define program_invocation_name (((char **)*_NSGetArgv())[0])
#define program_invocation_short_name (((char **)*_NSGetArgv())[0])
#endif

#define LEN(x) (sizeof(x)/sizeof(*x))
#define ALL(x) (x).begin(), (x).end()
#define REP(i, n) FOR(i, 0, n)
#define FOR(i, a, b) for (typename std::remove_cv<typename std::remove_reference<decltype(b)>::type>::type i = (a); i < (b); i++)
#define ROF(i, a, b) for (typename std::remove_cv<typename std::remove_reference<decltype(b)>::type>::type i = (b); --i >= (a); )

#define SGR0 "\x1b[m"
#define RED "\x1b[1;31m"
#define GREEN "\x1b[1;32m"
#define YELLOW "\x1b[1;33m"
#define BLUE "\x1b[1;34m"
#define MAGENTA "\x1b[1;35m"
#define CYAN "\x1b[1;36m"
#define NORMAL_YELLOW "\x1b[33m"
const long MAX_CODEPOINT = 0x10ffff;
extern long action_label_base, action_label, call_label_base, call_label, collapse_label_base, collapse_label;

void bold(long fd = 1);
void blue(long fd = 1);
void cyan(long fd = 1);
void green(long fd = 1);
void magenta(long fd = 1);
void red(long fd = 1);
void sgr0(long fd = 1);
void yellow(long fd = 1);
void normal_yellow(long fd = 1);
void indent(FILE* f, int d);

const size_t BUF_SIZE = 512;

void output_error(bool use_err, const char *format, va_list ap);
void err_msg(const char *format, ...);
void err_exit(int exitno, const char *format, ...);

long get_long(const char *arg);

void log_generic(const char *prefix, const char *format, va_list ap);
void log_event(const char *format, ...);
void log_action(const char *format, ...);
void log_status(const char *format, ...);

extern long debug_level;
extern FILE* debug_file;
#define DP(level, ...)  do {           \
    if (level <= debug_level) {        \
      fprintf(debug_file, "%s:%d ", __FILE__, __LINE__); \
      fprintf(debug_file, __VA_ARGS__);\
      fprintf(debug_file, "\n");       \
      fflush(debug_file);              \
    }                                  \
  } while (0)

template<class T, class... Args>
void emplace_front(vector<T>& a, Args&&... args)
{
  a.emplace(a.begin(), args...);
}

template<class T>
void sorted_insert(vector<T>& a, const T& x)
{
  a.emplace_back();
  auto it = a.end();
  while (a.begin() != --it && x < it[-1])
    *it = it[-1];
  *it = x;
}

template<class T, class... Args>
void sorted_emplace(vector<T>& a, Args&&... args)
{
  T x{args...};
  a.emplace_back();
  auto it = a.end();
  while (a.begin() != --it && x < it[-1])
    *it = it[-1];
  *it = x;
}

struct DisjointIntervals
{
  typedef std::pair<long, long> value_type;
  std::map<long, long> to;
  template<class... Args>
  void emplace(Args&&... args) {
    value_type x{args...};
    auto it = to.lower_bound(x.first);
    if (it != to.begin() && x.first <= prev(it)->second)
      x.first = (--it)->first;
    auto it2 = to.upper_bound(x.second);
    if (it2 != to.begin() && prev(it2)->first <= x.second && x.second < prev(it2)->second)
      x.second = prev(it2)->second;
    while (it != it2)
      it = to.erase(it);
    to.emplace(x);
  }
  void flip();
  void print(); // XXX
};
