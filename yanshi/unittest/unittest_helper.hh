#pragma once
#include "common.hh"
#include "fsa.hh"

#include <algorithm>
#include <sysexits.h>
#include <err.h>
#include <assert.h>
#include <iostream>
using namespace std;

static Fsa read_nfa()
{
  long n, m, k, u, a, v;
  cin >> n >> m >> k;
  Fsa r;
  r.start = 0;
  r.adj.resize(n);
  while (k--) {
    cin >> u;
    r.finals.push_back(u);
  }
  sort(ALL(r.finals));
  while (m--) {
    cin >> u >> a >> v;
    if (u < 0 || u >= n)
      errx(EX_DATAERR, "%ld: 0 <= u < n", u);
    if (a < -1 || a >= 256)
      errx(EX_DATAERR, "%ld: -1 <= c < 256", a);
    if (v < 0 || v >= n)
      errx(EX_DATAERR, "%ld: 0 <= v < n", v);
    r.adj[u].emplace_back(a, v);
  }
  assert(cin.good());
  REP(i, n)
    sort(ALL(r.adj[i]));
  return r;
}

static Fsa read_dfa()
{
  Fsa r = read_nfa();
  REP(i, r.n())
    if (r.adj[i].size()) {
      if (r.adj[i][0].first < 0)
        errx(EX_DATAERR, "epsilon edge found for %ld", i);
      REP(j, r.adj[i].size()-1)
        if (r.adj[i][j].first == r.adj[i][j+1].first)
          errx(EX_DATAERR, "duplicate labels %ld found for %ld", r.adj[i][j].first, i);
    }
  assert(cin.good());
  return r;
}

static void print_fsa(const Fsa& fsa)
{
  printf("finals:");
  for (long i: fsa.finals)
    printf(" %ld", i);
  puts("");
  puts("edges:");
  REP(i, fsa.n()) {
    printf("%ld:", i);
    for (auto& x: fsa.adj[i])
      printf(" (%ld,%ld)", x.first, x.second);
    puts("");
  }
}
