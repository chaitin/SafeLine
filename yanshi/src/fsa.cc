#include "common.hh"
#include "fsa.hh"
#include "option.hh"

#include <algorithm>
#include <assert.h>
#include <limits.h>
#include <queue>
#include <set>
#include <stack>
#include <tuple>
#include <unordered_map>
#include <unordered_set>
#include <utility>
#include <vector>
using namespace std;

namespace std
{
  template<typename T>
  struct hash<vector<T>> {
    size_t operator()(const vector<T>& v) const {
      hash<T> h;
      size_t r = 0;
      for (auto x: v)
        r = r*17+h(x);
      return r;
    }
  };
}

void Fsa::check() const
{
  REP(i, n())
    FOR(j, 1, adj[i].size())
      assert(adj[i][j-1].first.second == 0 &&
             adj[i][j].first.second == 0 ||
             adj[i][j-1].first.second <= adj[i][j].first.first);
}

bool Fsa::has(long u, long c) const
{
  auto it = upper_bound(ALL(adj[u]), make_pair(make_pair(c, LONG_MAX), LONG_MAX));
  return it != adj[u].begin() && c < (--it)->first.second;
}

bool Fsa::has_call(long u) const
{
  auto it = upper_bound(ALL(adj[u]), make_pair(make_pair(call_label_base, LONG_MAX), LONG_MAX));
  return (it != adj[u].end() && it->first.first < call_label) || (it != adj[u].begin() && call_label_base < (--it)->first.second);
}

bool Fsa::has_call_or_collapse(long u) const
{
  auto it = upper_bound(ALL(adj[u]), make_pair(make_pair(call_label_base, LONG_MAX), LONG_MAX));
  return it != adj[u].end() || (it != adj[u].begin() && call_label_base < (--it)->first.second);
}

long Fsa::transit(long u, long c) const
{
  auto it = upper_bound(ALL(adj[u]), make_pair(make_pair(c, LONG_MAX), LONG_MAX));
  return it != adj[u].begin() && c < (--it)->first.second ? it->second : -1;
}

bool Fsa::is_final(long x) const
{
  return binary_search(ALL(finals), x);
}

void Fsa::epsilon_closure(vector<long>& src) const
{
  static vector<bool> vis;
  if (n() > vis.size())
    vis.resize(n());
  for (long i: src)
    vis[i] = true;
  REP(i, src.size()) {
    long u = src[i];
    for (auto& e: adj[u]) {
      if (-1 < e.first.first) break;
      if (! vis[e.second]) {
        vis[e.second] = true;
        src.push_back(e.second);
      }
    }
  }
  for (long i: src)
    vis[i] = false;
  sort(ALL(src));
}

Fsa Fsa::operator~() const
{
  long accept = n();
  Fsa r;
  r.start = start;
  r.adj.resize(accept+1);
  REP(i, accept) {
    long j = 0;
    for (auto& e: adj[i]) {
      if (j < e.first.first)
        r.adj[i].emplace_back(make_pair(j, e.first.first), accept);
      r.adj[i].emplace_back(e.first, e.second);
      j = e.first.second;
    }
    if (j < AB)
      r.adj[i].emplace_back(make_pair(j, AB), accept);
  }
  r.adj[accept].emplace_back(make_pair(0, AB), accept);
  vector<long> new_finals;
  auto j = finals.begin();
  REP(i, accept+1) {
    while (j != finals.end() && *j < i)
      ++j;
    if (j == finals.end() || *j != i)
      new_finals.push_back(i);
  }
  r.finals = move(new_finals);
  return r;
}

void Fsa::accessible(const vector<long>* starts, function<void(long)> relate)
{
  vector<long> q{start}, id(n(), 0);
  id[start] = 1;
  if (starts)
    for (long u: *starts)
      if (! id[u]) {
        id[u] = 1;
        q.push_back(u);
      }
  REP(i, q.size()) {
    long u = q[i];
    for (auto& e: adj[u]) {
      //if (e.first.first >= AB) break;
      if (! id[e.second]) {
        id[e.second] = 1;
        q.push_back(e.second);
      }
    }
  }

  long j = 0;
  REP(i, n())
    id[i] = id[i] ? j++ : -1;

  auto it = finals.begin(), it2 = it;
  REP(i, n())
    if (id[i] >= 0) {
      relate(i);
      if (start == i)
        start = id[i];
      while (it != finals.end() && *it < i)
        ++it;
      if (it != finals.end() && *it == i)
        *it2++ = id[i];
      long k = 0;
      for (auto& e: adj[i])
        if (id[e.second] >= 0)
          adj[i][k++] = {e.first, id[e.second]}; // unordered unless deterministic
      adj[i].resize(k);
      if (id[i] != i)
        adj[id[i]] = move(adj[i]);
    }
  finals.erase(it2, finals.end());
  adj.resize(j);
}

void Fsa::co_accessible(const vector<bool>*final, function<void(long)> relate)
{
  vector<vector<long>> radj(n());
  REP(i, n())
    for (auto& e: adj[i]) {
      //if (e.first.first >= AB) break;
      radj[e.second].push_back(i);
    }
  REP(i, n())
    sort(ALL(radj[i]));
  vector<long> q = finals, id(n(), 0);
  for (long f: finals)
    id[f] = 1;
  if (final)
    REP(i, n())
      if ((*final)[i] && ! id[i]) {
        id[i] = 1;
        q.push_back(i);
      }
  REP(i, q.size()) {
    long u = q[i];
    for (auto& v: radj[u])
      if (! id[v]) {
        id[v] = 1;
        q.push_back(v);
      }
  }
  if (! id[start]) {
    start = 0;
    finals.clear();
    adj.assign(1, {});
    return;
  }

  long j = 0;
  REP(i, n())
    id[i] = id[i] ? j++ : -1;

  auto it = finals.begin(), it2 = it;
  REP(i, n())
    if (id[i] >= 0) {
      relate(i);
      if (start == i)
        start = id[i];
      while (it != finals.end() && *it < i)
        ++it;
      if (it != finals.end() && *it == i)
        *it2++ = id[i];
      long k = 0;
      for (auto& e: adj[i])
        if (id[e.second] >= 0)
          adj[i][k++] = {e.first, id[e.second]}; // unordered unless deterministic
      adj[i].resize(k);
      if (id[i] != i)
        adj[id[i]] = move(adj[i]);
    }
  finals.erase(it2, finals.end());
  adj.resize(j);
}

Fsa Fsa::difference(const Fsa& rhs, function<void(long)> relate) const
{
  Fsa r;
  vector<pair<long, long>> q;
  unordered_map<long, long> m;
  q.emplace_back(start, rhs.start);
  m[(rhs.n()+1) * start + rhs.start] = 0;
  r.start = 0;
  REP(i, q.size()) {
    long u0, u1, v0, v1;
    tie(u0, u1) = q[i];
    if (is_final(u0) && ! rhs.is_final(u1))
      r.finals.push_back(i);
    r.adj.emplace_back();
    relate(u0);
    vector<Edge>::const_iterator it0 = adj[u0].begin(), it1, it1e;
    if (u1 == rhs.n())
      it1 = it1e = rhs.adj[0].end();
    else {
      it1 = rhs.adj[u1].begin();
      it1e = rhs.adj[u1].end();
    }
    long last = LONG_MIN;
    while (it0 != adj[u0].end()) {
      long from = max(last, it0->first.first), to = it0->first.second;
      while (it1 != it1e && it1->first.second <= from)
        ++it1;
      if (it1 != it1e)
        to = min(to, from < it1->first.first ? it1->first.first : it1->first.second);
      last = to;
      long v1 = it1 != it1e && it1->first.first <= from ? it1->second : rhs.n(),
           t = (rhs.n()+1) * it0->second + v1;
      auto mit = m.find(t);
      if (mit == m.end()) {
        mit = m.emplace(t, m.size()).first;
        q.emplace_back(it0->second, v1);
      }
      r.adj[i].emplace_back(make_pair(from, to), mit->second);
      if (to == it0->first.second)
        ++it0;
    }
  }
  return r;
}

Fsa Fsa::intersect(const Fsa& rhs, function<void(long, long)> relate) const
{
  Fsa r;
  vector<pair<long, long>> q;
  long u0, u1, v0, v1;
  unordered_map<long, long> m;
  q.emplace_back(start, rhs.start);
  m[rhs.n() * start + rhs.start] = 0;
  r.start = 0;
  REP(i, q.size()) {
    tie(u0, u1) = q[i];
    if (is_final(u0) && rhs.is_final(u1))
      r.finals.push_back(i);
    r.adj.emplace_back();
    relate(u0, u1);
    auto it0 = adj[u0].begin(), it1 = rhs.adj[u1].begin();
    while (it0 != adj[u0].end() && it1 != rhs.adj[u1].end()) {
      if (it0->first.second <= it1->first.first)
        ++it0;
      else if (it1->first.second <= it0->first.first)
        ++it1;
      else {
        long t = rhs.n() * it0->second + it1->second;
        auto mit = m.find(t);
        if (mit == m.end()) {
          mit = m.emplace(t, m.size()).first;
          q.emplace_back(it0->second, it1->second);
        }
        r.adj[i].emplace_back(make_pair(max(it0->first.first, it1->first.first), min(it0->first.second, it1->first.second)), mit->second);
        if (it0->first.second < it1->first.second)
          ++it0;
        else if (it0->first.second > it1->first.second)
          ++it1;
        else
          ++it0, ++it1;
      }
    }
  }
  return r;
}

Fsa Fsa::determinize(const vector<long>* starts, function<void(long, const vector<long>&)> relate) const
{
  Fsa r;
  r.start = 0;
  unordered_map<vector<long>, long> m;
  vector<vector<Edge>::const_iterator> its(n());
  vector<long> vs{start};
  vector<pair<long, long>> events;
  stack<vector<long>> st;
  epsilon_closure(vs);
  m[vs] = 0;
  st.push(move(vs));
  if (starts)
    for (long u: *starts) {
      vs.assign(1, u);
      epsilon_closure(vs);
      if (! m.count(vs)) {
        m.emplace(vs, m.size());
        st.push(move(vs));
      }
    }
  while (st.size()) {
    vector<long> x = move(st.top());
    st.pop();
    long id = m[x];
    if (id+1 > r.adj.size())
      r.adj.resize(id+1);
    relate(id, x);
    bool final = false;
    events.clear();
    for (long u: x) {
      if (is_final(u))
        final = true;
      for (auto& e: adj[u]) {
        events.emplace_back(e.first.first, e.second);
        events.emplace_back(e.first.second, ~ e.second);
      }
    }
    if (final)
      r.finals.push_back(id);
    long last = 0;
    multiset<long> live;
    sort(ALL(events));
    for (auto& ev: events) {
      if (last < ev.first) {
        if (live.size()) {
          vs.assign(ALL(live));
          vs.erase(unique(ALL(vs)), vs.end());
          epsilon_closure(vs);
          auto mit = m.find(vs);
          if (mit == m.end()) {
            mit = m.emplace(vs, m.size()).first;
            st.push(vs);
          }
          if (r.adj[id].size() && r.adj[id].back().first.second == last && r.adj[id].back().second == mit->second) // coalesce two edges
            r.adj[id].back().first.second = ev.first;
          else
            r.adj[id].emplace_back(make_pair(last, ev.first), mit->second);
        }
        last = ev.first;
      }
      if (ev.second >= 0)
        live.insert(ev.second);
      else
        live.erase(live.find(~ ev.second));
    }
  }
  sort(ALL(r.finals));
  return r;
}

Fsa Fsa::distinguish(function<void(vector<long>&)> relate) const
{
  vector<long> scale;
  REP(i, n())
    for (auto& e: adj[i]) {
      scale.push_back(e.first.first);
      scale.push_back(e.first.second);
    }
  sort(ALL(scale));
  scale.erase(unique(ALL(scale)), scale.end());

  vector<vector<pair<long, long>>> radj(n());
  REP(i, n())
    for (auto& e: adj[i]) {
      long from = lower_bound(ALL(scale), e.first.first) - scale.begin(),
           to = lower_bound(ALL(scale), e.first.second) - scale.begin();
      FOR(j, from, to)
        radj[e.second].emplace_back(j, i);
    }
  REP(i, n())
    sort(ALL(radj[i]));
  vector<long> L(n()), R(n()), B(n()), C(n(), 0), CC(n(), 0);
  vector<bool> mark(n(), false);

  // distinguish finals & non-finals
  long fx = -1, x = -1, fy = -1, y = -1, j = 0;
  REP(i, n())
    if (j < finals.size() && finals[j] == i) {
      j++;
      if (y < 0)
        fy = i;
      else
        R[y] = i;
      C[B[i] = fy]++;
      L[i] = y;
      y = i;
    } else {
      if (x < 0)
        fx = i;
      else
        R[x] = i;
      C[B[i] = fx]++;
      L[i] = x;
      x = i;
    }
  if (x >= 0)
    L[fx] = x, R[x] = fx;
  if (y >= 0)
    L[fy] = y, R[y] = fy;

  set<pair<long, long>> refines;
  auto labels = [&](long fx) {
    vector<long> lb;
    for (long x = fx; ; ) {
      for (auto& e: radj[x])
        lb.push_back(e.first);
      if ((x = R[x]) == fx) break;
    }
    sort(ALL(lb));
    lb.erase(unique(ALL(lb)), lb.end());
    return lb;
  };

  if (fx >= 0)
    for (long a: labels(fx))
      refines.emplace(a, fx);
  if (fy >= 0)
    for (long a: labels(fy))
      refines.emplace(a, fy);
  while (refines.size()) {
    long a;
    tie(a, fx) = *refines.begin();
    refines.erase(refines.begin());
    // count
    vector<long> bs;
    for (x = fx; ; ) {
      auto it = lower_bound(ALL(radj[x]), make_pair(a, 0L)),
           ite = upper_bound(ALL(radj[x]), make_pair(a, n()));
      for (; it != ite; ++it) {
        y = it->second;
        if (! CC[B[y]]++)
          bs.push_back(B[y]);
        mark[y] = true;
      }
      if ((x = R[x]) == fx) break;
    }
    // for each refinable set
    for (long fy: bs) {
      if (CC[fy] < C[fy]) {
        long fu = -1, u = -1, cu = 0,
             fv = -1, v = -1, cv = 0;
        vector<long> lb = labels(fy);
        for (long i = fy; ; ) {
          if (mark[i]) {
            mark[i] = false;
            if (u < 0)
              C[fu = i] = 0;
            else
              R[u] = i;
            C[fu]++;
            B[i] = fu;
            L[i] = u;
            u = i;
          } else {
            if (v < 0)
              C[fv = i] = 0;
            else
              R[v] = i;
            C[fv]++;
            B[i] = fv;
            L[i] = v;
            v = i;
          }
          if ((i = R[i]) == fy) break;
        }
        L[fu] = u, R[u] = fu;
        L[fv] = v, R[v] = fv;
        //REP(a, AB+1)
        for (long a: lb)
          if (refines.count({a, fy}))
            refines.emplace(a, fu != fy ? fu : fv);
          else
            refines.emplace(a, C[fu] < C[fv] ? fu : fv);
      } else
        for (long i = fy; ; ) {
          mark[i] = false;
          if ((i = R[i]) == fy) break;
        }
      CC[fy] = 0;
    }
    // clear marks
    for (x = fx; ; ) {
      auto it = lower_bound(ALL(radj[x]), make_pair(a, 0L)),
           ite = upper_bound(ALL(radj[x]), make_pair(a, n()));
      for (; it != ite; ++it) {
        y = it->second;
        CC[B[y]] = 0;
        mark[y] = false;
      }
      if ((x = R[x]) == fx) break;
    }
  }

  Fsa r;
  long nn = 0;
  vector<long> vs;
  REP(i, n())
    if (B[i] == i) {
      vs.clear();
      for (long j = i; ; ) {
        B[j] = nn;
        vs.push_back(j);
        if ((j = R[j]) == i) break;
      }
      relate(vs);
      if (binary_search(ALL(finals), i))
        r.finals.push_back(nn);
      nn++;
    }
  r.start = B[start];
  r.adj.resize(nn);
  REP(i, n())
    for (auto& e: adj[i])
      r.adj[B[i]].emplace_back(e.first, B[e.second]);
  REP(i, nn) {
    // merge edges with the same destination
    sort(ALL(r.adj[i]), [](const Edge& x, const Edge& y) {
      return x.second != y.second ? x.second < y.second : x.first < y.first;
    });
    auto it2 = r.adj[i].begin();
    for (auto it = r.adj[i].begin(); it != r.adj[i].end(); ) {
      long v = it->second, from = it->first.first, to = it->first.second;
      while (++it != r.adj[i].end() && it->second == v)
        if (it->first.first <= to)
          to = max(to, it->first.second);
        else {
          *it2++ = make_pair(make_pair(from, to), v);
          tie(from, to) = it->first;
        }
      *it2++ = make_pair(make_pair(from, to), v);
    }
    r.adj[i].erase(it2, r.adj[i].end());
    sort(ALL(r.adj[i]));
  }
  return r;
}
