#include "common.hh"
#include "compiler.hh"
#include "fsa_anno.hh"
#include "loader.hh"
#include "option.hh"

#include <algorithm>
#include <limits.h>
#include <map>
#include <unicode/utf8.h>
#include <utility>
using namespace std;

bool operator<(ExprTag x, ExprTag y)
{
  return long(x) < long(y);
}

bool assoc_has_expr(vector<pair<Expr*, ExprTag>>& as, Expr* x)
{
  auto it = lower_bound(ALL(as), make_pair(x, ExprTag(0)));
  return it != as.end() && it->first == x;
}

void sort_assoc(vector<pair<Expr*, ExprTag>>& as)
{
  sort(ALL(as));
  auto i = as.begin(), j = i, k = i;
  for (; i != as.end(); i = j) {
    while (++j != as.end() && i->first == j->first)
      i->second = ExprTag(long(i->second) | long(j->second));
    *k++ = *i;
  }
  as.erase(k, as.end());
}

void FsaAnno::add_assoc(Expr& expr)
{
  // has actions: actions need tags to differentiate 'entering', 'leaving', ...
  // 'intact': states with the 'inner' tag cannot be connected to start/final in substring grammar
  // 'CallExpr' 'CollapseExpr': differentiate states representing 'CallExpr' 'CollapseExpr' (u, special, v)
  // 'opt_mode': displaying possible positions for given strings in interactive mode
  if (expr.no_action() && ! expr.stmt->intact && ! dynamic_cast<CallExpr*>(&expr) && ! dynamic_cast<CollapseExpr*>(&expr) && opt_mode != Mode::interactive)
    return;
  auto j = fsa.finals.begin();
  REP(i, fsa.n()) {
    ExprTag tag = ExprTag(0);
    if (i == fsa.start)
      tag = ExprTag::start;
    while (j != fsa.finals.end() && *j < i)
      ++j;
    if (j != fsa.finals.end() && *j == i)
      tag = ExprTag(long(tag) | long(ExprTag::final));
    if (tag == ExprTag(0))
      tag = ExprTag::inner;
    sorted_insert(assoc[i], make_pair(&expr, tag));
  }
  // Add pseudo transitions with labels [ACTION_LABEL_BASE, COLLAPSE_LABEL_BASE) to prevent its merge with other states
  if (expr.leaving.size() || expr.entering.size() || expr.transiting.size())
    for (auto action: expr.transiting)
      REP(i, fsa.n()) {
        fsa.adj[i].emplace_back(make_pair(action_label, action_label+1), i);
        action_label++;
      }
  else if (expr.finishing.size())
    for (long f: fsa.finals) {
      fsa.adj[f].emplace_back(make_pair(action_label, action_label+1), f);
      action_label++;
    }
}

void FsaAnno::accessible(const vector<long>* starts, vector<long>& mapping) {
  long allo = 0;
  auto relate = [&](long x) {
    if (allo != x)
      assoc[allo] = move(assoc[x]);
    allo++;
    mapping.push_back(x);
  };
  fsa.accessible(starts, relate);
  assoc.resize(allo);
}

void FsaAnno::co_accessible(const vector<bool>* final, vector<long>& mapping) {
  long allo = 0;
  auto relate = [&](long x) {
    if (allo != x)
      assoc[allo] = move(assoc[x]);
    allo++;
    mapping.push_back(x);
  };
  fsa.co_accessible(final, relate);
  if (fsa.finals.empty()) { // 'start' does not produce acceptable strings
    assoc.assign(1, {});
    mapping.assign(1, 0);
    deterministic = true;
    return;
  }
  if (! deterministic)
    REP(i, fsa.n())
      sort(ALL(fsa.adj[i]));
  assoc.resize(allo);
}

void FsaAnno::complement(ComplementExpr* expr) {
  if (! deterministic)
    fsa = fsa.determinize(NULL, [&](long, const vector<long>&){});
  fsa = ~ fsa;
  assoc.assign(fsa.n(), {});
  deterministic = true;
}

void FsaAnno::concat(FsaAnno& rhs, ConcatExpr* expr) {
  long ln = fsa.n(), rn = rhs.fsa.n();
  for (long f: fsa.finals)
    emplace_front(fsa.adj[f], epsilon, ln+rhs.fsa.start);
  for (auto& es: rhs.fsa.adj) {
    for (auto& e: es)
      e.second += ln;
    fsa.adj.emplace_back(move(es));
  }
  fsa.finals = move(rhs.fsa.finals);
  for (long& f: fsa.finals)
    f += ln;
  assoc.resize(fsa.n());
  REP(i, rhs.fsa.n())
    assoc[ln+i] = move(rhs.assoc[i]);
  if (expr)
    add_assoc(*expr);
  deterministic = false;
}

void FsaAnno::determinize(const vector<long>* starts, vector<vector<long>>* mapping) {
  if (deterministic)
    return;
  decltype(assoc) new_assoc;
  auto relate = [&](long id, const vector<long>& xs) {
    if (id+1 > new_assoc.size()) {
      new_assoc.resize(id+1);
      if (mapping)
        mapping->resize(id+1);
    }
    auto& as = new_assoc[id];
    for (long x: xs)
      as.insert(as.end(), ALL(assoc[x]));
    sort_assoc(as);
    if (mapping)
      (*mapping)[id] = xs;
  };
  fsa = fsa.determinize(starts, relate);
  assoc = move(new_assoc);
  deterministic = true;
}

void FsaAnno::difference(FsaAnno& rhs, DifferenceExpr* expr) {
  vector<vector<long>> rel0;
  decltype(rhs.assoc) new_assoc;
  auto relate0 = [&](long id, const vector<long>& xs) {
    if (id+1 > rel0.size())
      rel0.resize(id+1);
    rel0[id] = xs;
  };
  auto relate = [&](long x) {
    if (rel0.empty())
      new_assoc.emplace_back(assoc[x]);
    else {
      new_assoc.emplace_back();
      auto& as = new_assoc.back();
      for (long u: rel0[x])
        as.insert(as.end(), ALL(assoc[u]));
      sort_assoc(as);
    }
  };
  if (! deterministic)
    fsa = fsa.determinize(NULL, relate0);
  if (! rhs.deterministic)
    rhs.fsa = rhs.fsa.determinize(NULL, [](long, const vector<long>&) {});
  fsa = fsa.difference(rhs.fsa, relate);
  assoc = move(new_assoc);
  if (expr)
    add_assoc(*expr);
  deterministic = true;
}

FsaAnno FsaAnno::epsilon_fsa(EpsilonExpr* expr) {
  FsaAnno r;
  r.fsa.start = 0;
  r.fsa.finals.push_back(0);
  r.fsa.adj.resize(1);
  r.assoc.resize(1);
  if (expr)
    r.add_assoc(*expr);
  r.deterministic = true;
  return r;
}

void FsaAnno::intersect(FsaAnno& rhs, IntersectExpr* expr) {
  decltype(rhs.assoc) new_assoc;
  vector<vector<long>> rel0, rel1;
  auto relate0 = [&](long id, const vector<long>& xs) {
    if (id+1 > rel0.size())
      rel0.resize(id+1);
    rel0[id] = xs;
  };
  auto relate1 = [&](long id, const vector<long>& xs) {
    if (id+1 > rel1.size())
      rel1.resize(id+1);
    rel1[id] = xs;
  };
  auto relate = [&](long x, long y) {
    new_assoc.emplace_back();
    auto& as = new_assoc.back();
    if (rel0.empty())
      as.insert(as.end(), ALL(assoc[x]));
    else
      for (long u: rel0[x])
        as.insert(as.end(), ALL(assoc[u]));
    if (rel1.empty())
      as.insert(as.end(), ALL(rhs.assoc[y]));
    else
      for (long v: rel1[y])
        as.insert(as.end(), ALL(rhs.assoc[v]));
    sort_assoc(as);
  };
  if (! deterministic)
    fsa = fsa.determinize(NULL, relate0);
  if (! rhs.deterministic)
    rhs.fsa = rhs.fsa.determinize(NULL, relate1);
  fsa = fsa.intersect(rhs.fsa, relate);
  assoc = move(new_assoc);
  if (expr)
    add_assoc(*expr);
  deterministic = true;
}

void FsaAnno::minimize(vector<vector<long>>* mapping) {
  assert(deterministic);
  decltype(assoc) new_assoc;
  auto relate = [&](vector<long>& xs) {
    new_assoc.emplace_back();
    auto& as = new_assoc.back();
    for (long x: xs)
      as.insert(as.end(), ALL(assoc[x]));
    sort_assoc(as);
    if (mapping)
      mapping->push_back(xs);
  };
  fsa = fsa.distinguish(relate);
  assoc = move(new_assoc);
}

void FsaAnno::union_(FsaAnno& rhs, UnionExpr* expr) {
  long ln = fsa.n(), rn = rhs.fsa.n(), src = ln+rn,
       old_lsrc = fsa.start;
  fsa.start = src;
  for (long f: rhs.fsa.finals)
    fsa.finals.push_back(ln+f);
  for (auto& es: rhs.fsa.adj) {
    for (auto& e: es)
      e.second += ln;
    fsa.adj.emplace_back(move(es));
  }
  fsa.adj.emplace_back();
  fsa.adj[src].emplace_back(epsilon, old_lsrc);
  fsa.adj[src].emplace_back(epsilon, ln+rhs.fsa.start);
  assoc.resize(fsa.n());
  REP(i, rhs.fsa.n())
    assoc[ln+i] = move(rhs.assoc[i]);
  if (expr)
    add_assoc(*expr);
  deterministic = false;
}

void FsaAnno::plus(PlusExpr* expr) {
  for (long f: fsa.finals)
    emplace_front(fsa.adj[f], epsilon, fsa.start);
  if (expr)
    add_assoc(*expr);
  deterministic = false;
}

void FsaAnno::question(QuestionExpr* expr) {
  long src = fsa.n(), sink = src+1, old_src = fsa.start;
  fsa.start = src;
  fsa.adj.emplace_back();
  fsa.adj.emplace_back();
  fsa.adj[src].emplace_back(epsilon, old_src);
  fsa.adj[src].emplace_back(epsilon, sink);
  fsa.finals.push_back(sink);
  assoc.resize(fsa.n());
  if (expr)
    add_assoc(*expr);
  deterministic = false;
}

void FsaAnno::repeat(RepeatExpr& expr) {
  FsaAnno r = epsilon_fsa(NULL);
  REP(i, expr.low) {
    FsaAnno t = *this;
    r.concat(t, NULL);
  }
  if (expr.high == LONG_MAX) {
    star(NULL);
    r.concat(*this, NULL);
  } else if (expr.low < expr.high) {
    FsaAnno rhs = epsilon_fsa(NULL), x = *this;
    ROF(i, 0, expr.high-expr.low) {
      FsaAnno t = x;
      rhs.union_(t, NULL);
      if (i) {
        t = *this;
        x.concat(t, NULL);
      }
    }
    r.concat(rhs, NULL);
  }
  r.deterministic = false;
  *this = move(r);
}

void FsaAnno::star(StarExpr* expr) {
  long src = fsa.n(), sink = src+1, old_src = fsa.start;
  fsa.start = src;
  fsa.adj.emplace_back();
  fsa.adj.emplace_back();
  fsa.adj[src].emplace_back(epsilon, old_src);
  fsa.adj[src].emplace_back(epsilon, sink);
  for (long f: fsa.finals) {
    sorted_emplace(fsa.adj[f], epsilon, old_src);
    sorted_emplace(fsa.adj[f], epsilon, sink);
  }
  fsa.finals.assign(1, sink);
  assoc.resize(fsa.n());
  if (expr)
    add_assoc(*expr);
  deterministic = false;
}

FsaAnno FsaAnno::bracket(BracketExpr& expr) {
  FsaAnno r;
  r.fsa.start = 0;
  r.fsa.finals = {1};
  r.fsa.adj.resize(2);
  for (auto& x: expr.intervals.to)
    r.fsa.adj[0].emplace_back(x, 1);
  r.assoc.resize(2);
  r.add_assoc(expr);
  r.deterministic = true;
  return r;
}

FsaAnno FsaAnno::call(CallExpr& expr) {
  // represented by (0, special, 1)
  FsaAnno r;
  r.fsa.start = 0;
  r.fsa.finals = {1};
  r.fsa.adj.resize(2);
  r.fsa.adj[0].emplace_back(make_pair(call_label, call_label+1), 1);
  call_label++;
  r.assoc.resize(2);
  r.add_assoc(expr);
  r.deterministic = true;
  return r;
}

FsaAnno FsaAnno::collapse(CollapseExpr& expr) {
  // represented by (0, special, 1)
  FsaAnno r;
  r.fsa.start = 0;
  r.fsa.finals = {1};
  r.fsa.adj.resize(2);
  r.fsa.adj[0].emplace_back(make_pair(collapse_label, collapse_label+1), 1);
  collapse_label++;
  r.assoc.resize(2);
  r.add_assoc(expr);
  r.deterministic = true;
  return r;
}

FsaAnno FsaAnno::dot(DotExpr* expr) {
  FsaAnno r;
  r.fsa.start = 0;
  r.fsa.finals = {1};
  r.fsa.adj.resize(2);
  r.fsa.adj[0].emplace_back(make_pair(0L, AB), 1);
  r.assoc.resize(2);
  if (expr)
    r.add_assoc(*expr);
  r.deterministic = true;
  return r;
}

FsaAnno FsaAnno::embed(EmbedExpr& expr) {
  if (expr.define_stmt) {
    FsaAnno r = compiled[expr.define_stmt];
    // change the labels to differentiate instances of CallExpr
    REP(i, r.fsa.n()) {
      auto it = upper_bound(ALL(r.fsa.adj[i]), make_pair(make_pair(call_label_base, LONG_MAX), LONG_MAX));
      if (it != r.fsa.adj[i].begin() && call_label_base < (it-1)->first.second)
        --it;
      for (; it != r.fsa.adj[i].end() && it->first.first < call_label; ++it) {
        assert(call_label_base <= it->first.first);
        long t = it->first.second-it->first.first;
        it->first.first = call_label;
        call_label += t;
        it->first.second = call_label;
        assert(it->first.second <= call_label);
      }
    }
    r.add_assoc(expr);
    return r;
  } else { // macro
    FsaAnno r;
    r.fsa.start = 0;
    r.fsa.finals = {1};
    r.fsa.adj.resize(2);
    r.fsa.adj[0].emplace_back(make_pair(expr.macro_value, expr.macro_value+1), 1);
    r.assoc.resize(2);
    r.add_assoc(expr);
    r.deterministic = true;
    return r;
  }
}

FsaAnno FsaAnno::literal(LiteralExpr& expr) {
  FsaAnno r;
  r.fsa.start = 0;
  long len = 0;
  if (opt_bytes) {
    len = expr.literal.size();
    r.fsa.adj.resize(len+1);
    REP(i, expr.literal.size()) {
      long c = (u8)expr.literal[i];
      r.fsa.adj[i].emplace_back(make_pair(c, c+1), i+1);
    }
  } else {
    for (i32 c, i = 0; i < expr.literal.size(); len++) {
      U8_NEXT_OR_FFFD(expr.literal.c_str(), i, expr.literal.size(), c);
      r.fsa.adj.emplace_back();
      r.fsa.adj[len].emplace_back(make_pair(c, c+1), len+1);
    }
    r.fsa.adj.emplace_back();
  }
  r.fsa.finals.push_back(len);
  r.assoc.resize(len+1);
  r.add_assoc(expr);
  r.deterministic = true;
  return r;
}

void FsaAnno::substring_grammar() {
  long src = fsa.n(), sink = src+1, old_src = fsa.start;
  fsa.start = src;
  fsa.adj.emplace_back();
  fsa.adj.emplace_back();
  REP(i, src) {
    bool ok = true;
    for (auto aa: assoc[i])
      if (auto e = dynamic_cast<CollapseExpr*>(aa.first)) {
        if (e->define_stmt->intact && has_inner(aa.second)) {
          ok = false;
          break;
        }
      } else if (aa.first->stmt->intact && has_inner(aa.second)) {
        ok = false;
        break;
      }
    if (ok || i == old_src)
      fsa.adj[src].emplace_back(epsilon, i);
    if (ok || fsa.is_final(i))
      emplace_front(fsa.adj[i], epsilon, sink);
  }
  fsa.finals.assign(1, sink);
  assoc.resize(fsa.n());
  deterministic = false;
}
