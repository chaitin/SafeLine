#include "compiler.hh"
#include "fsa_anno.hh"
#include "loader.hh"
#include "option.hh"

#include <algorithm>
#include <ctype.h>
#include <limits.h>
#include <map>
#include <sstream>
#include <stack>
#include <unordered_map>
using namespace std;

unordered_map<DefineStmt*, FsaAnno> compiled;
static unordered_map<DefineStmt*, vector<pair<long, long>>> stmt2call_addr;
static unordered_map<DefineStmt*, vector<bool>> stmt2final;

void print_assoc(const FsaAnno& anno)
{
  magenta(); printf("=== Associated Expr of each state\n"); sgr0();
  REP(i, anno.fsa.n()) {
    printf("%ld:", i);
    for (auto aa: anno.assoc[i]) {
      auto a = aa.first;
      printf(" %s%s%s%s(%ld-%ld", a->name().c_str(),
             has_start(aa.second) ? "^" : "",
             has_inner(aa.second) ? "." : "",
             has_final(aa.second) ? "$" : "",
             a->loc.start, a->loc.end);
      if (a->entering.size())
        printf(",>%zd", a->entering.size());
      if (a->leaving.size())
        printf(",%%%zd", a->leaving.size());
      if (a->finishing.size())
        printf(",@%zd", a->finishing.size());
      if (a->transiting.size())
        printf(",$%zd", a->transiting.size());
      printf(")");
    }
    puts("");
  }
  puts("");
}

void print_automaton(const Fsa& fsa)
{
  magenta(); printf("=== Automaton\n"); sgr0();
  green(); printf("start: %ld\n", fsa.start);
  red(); printf("finals:");
  for (long i: fsa.finals)
    printf(" %ld", i);
  puts("");
  sgr0(); puts("edges:");
  REP(i, fsa.n()) {
    printf("%ld:", i);
    for (auto it = fsa.adj[i].begin(); it != fsa.adj[i].end(); ) {
      long from = it->first.first, to = it->first.second, v = it->second;
      while (++it != fsa.adj[i].end() && to == it->first.first && it->second == v)
        to = it->first.second;
      if (from == to-1)
        printf(" (%ld,%ld)", from, v);
      else
        printf(" (%ld-%ld,%ld)", from, to-1, v);
    }
    puts("");
  }
  puts("");
}

Expr* find_lca(Expr* u, Expr* v)
{
  if (u->depth > v->depth)
    swap(u, v);
  if (u->depth < v->depth)
    for (long k = 63-__builtin_clzl(v->depth-u->depth); k >= 0; k--)
      if (u->depth <= v->depth-(1L<<k))
        v = v->anc[k];
  if (u == v)
    return u;
  if (v->depth)
    for (long k = 63-__builtin_clzl(v->depth); k >= 0; k--)
      if (k < u->anc.size() && u->anc[k] != v->anc[k])
        u = u->anc[k], v = v->anc[k];
  return u->anc[0]; // NULL if two trees
}

struct Compiler : Visitor<Expr> {
  stack<FsaAnno> st;
  stack<Expr*> path;
  long tick = 0;

  void pre_expr(Expr& expr) {
    expr.pre = tick++;
    expr.depth = path.size();
    if (path.size()) {
      expr.anc.assign(1, path.top());
      for (long k = 1; 1L << k <= expr.depth; k++)
        expr.anc.push_back(expr.anc[k-1]->anc[k-1]);
    } else
      expr.anc.assign(1, nullptr);
    path.push(&expr);
    DP(5, "%s(%ld-%ld)", expr.name().c_str(), expr.loc.start, expr.loc.end);
  }
  void post_expr(Expr& expr) {
    path.pop();
    expr.post = tick;
#ifdef DEBUG
    st.top().fsa.check();
#endif
  }

  void visit(Expr& expr) override {
    pre_expr(expr);
    expr.accept(*this);
    post_expr(expr);
  }
  void visit(BracketExpr& expr) override {
    st.push(FsaAnno::bracket(expr));
  }
  void visit(CallExpr& expr) override {
    st.push(FsaAnno::call(expr));
  }
  void visit(CollapseExpr& expr) override {
    st.push(FsaAnno::collapse(expr));
  }
  void visit(ComplementExpr& expr) override {
    visit(*expr.inner);
    st.top().complement(&expr);
  }
  void visit(ConcatExpr& expr) override {
    visit(*expr.rhs);
    FsaAnno rhs = move(st.top());
    visit(*expr.lhs);
    st.top().concat(rhs, &expr);
  }
  void visit(DifferenceExpr& expr) override {
    visit(*expr.rhs);
    FsaAnno rhs = move(st.top());
    visit(*expr.lhs);
    st.top().difference(rhs, &expr);
  }
  void visit(DotExpr& expr) override {
    st.push(FsaAnno::dot(&expr));
  }
  void visit(EmbedExpr& expr) override {
    st.push(FsaAnno::embed(expr));
  }
  void visit(EpsilonExpr& expr) override {
    st.push(FsaAnno::epsilon_fsa(&expr));
  }
  void visit(IntersectExpr& expr) override {
    visit(*expr.rhs);
    FsaAnno rhs = move(st.top());
    visit(*expr.lhs);
    st.top().intersect(rhs, &expr);
  }
  void visit(LiteralExpr& expr) override {
    st.push(FsaAnno::literal(expr));
  }
  void visit(PlusExpr& expr) override {
    visit(*expr.inner);
    st.top().plus(&expr);
  }
  void visit(QuestionExpr& expr) override {
    visit(*expr.inner);
    st.top().question(&expr);
  }
  void visit(RepeatExpr& expr) override {
    visit(*expr.inner);
    st.top().repeat(expr);
  }
  void visit(StarExpr& expr) override {
    visit(*expr.inner);
    st.top().star(&expr);
  }
  void visit(UnionExpr& expr) override {
    visit(*expr.rhs);
    FsaAnno rhs = move(st.top());
    visit(*expr.lhs);
    st.top().union_(rhs, &expr);
  }
};

void compile(DefineStmt* stmt)
{
  if (compiled.count(stmt))
    return;
  FsaAnno& anno = compiled[stmt];
  Compiler comp;
  comp.visit(*stmt->rhs);
  anno = move(comp.st.top());
  anno.determinize(NULL, NULL);
  anno.minimize(NULL);
  DP(4, "size(%s::%s) = %ld", stmt->module->filename.c_str(), stmt->lhs.c_str(), anno.fsa.n());
}

void generate_transitions(DefineStmt* stmt)
{
  FsaAnno& anno = compiled[stmt];
  auto& call_addr = stmt2call_addr[stmt];
  auto& sub_final = stmt2final[stmt];
  auto find_within = [&](long u) {
    vector<pair<Expr*, ExprTag>> within;
    Expr* last = NULL;
    sort(ALL(anno.assoc[u]), [](const pair<Expr*, ExprTag>& x, const pair<Expr*, ExprTag>& y) {
      if (x.first->pre != y.first->pre)
        return x.first->pre < y.first->pre;
      return x.second < y.second;
    });
    for (auto aa: anno.assoc[u]) {
      Expr* stop = last ? find_lca(last, aa.first) : NULL;
      last = aa.first;
      for (Expr* x = aa.first; x != stop; x = x->anc[0])
        within.emplace_back(x, aa.second);
    }
    sort(ALL(within));
    auto j = within.begin();
    for (auto i = within.begin(); i != within.end(); ) {
      Expr* x = i->first;
      long t = long(i->second);
      while (++i != within.end() && x == i->first)
        t |= long(i->second);
      *j++ = {x, ExprTag(t)};
    }
    within.erase(j, within.end());
    return within;
  };
  decltype(anno.assoc) withins(anno.fsa.n());
  REP(i, anno.fsa.n())
    withins[i] = move(find_within(i));

  auto get_code = [](Action* action) {
    if (auto t = dynamic_cast<InlineAction*>(action))
      return t->code;
    else if (auto t = dynamic_cast<RefAction*>(action))
      return t->define_stmt->code;
    else
      assert(0);
    return string();
  };

#define D(S) if (opt_dump_action) { \
               if (auto t = dynamic_cast<InlineAction*>(action.first)) { \
                 if (from == to-1) \
                   printf(S " %ld %ld %ld %s\n", u, from, v, t->code.c_str()); \
                 else \
                   printf(S " %ld %ld-%ld %ld %s\n", u, from, to-1, v, t->code.c_str()); \
               } else if (auto t = dynamic_cast<RefAction*>(action.first)) { \
                 if (from == to-1) \
                   printf(S " %ld %ld %ld %s\n", u, from, v, t->define_stmt->code.c_str()); \
                 else \
                   printf(S " %ld %ld-%ld %ld %s\n", u, from, to-1, v, t->define_stmt->code.c_str()); \
               } \
             }

  if (output_header) {
    if (opt_gen_c) {
      if (opt_gen_extern_c)
        fputs("extern \"C\" ", output_header);
      fprintf(output_header, "long yanshi_%s_transit(long* ret_stack, long* ret_stack_len, long u, long c", stmt->lhs.c_str());
    }
    else
      fprintf(output_header, "long yanshi_%s_transit(vector<long>& ret_stack, long u, long c", stmt->lhs.c_str());
    if (stmt->export_params.size())
      fprintf(output_header, ", %s", stmt->export_params.c_str());
    fprintf(output_header, ");\n");
  }
  if (opt_gen_c) {
    if (opt_gen_extern_c)
      fputs("extern \"C\" ", output);
    fprintf(output, "long yanshi_%s_transit(long* ret_stack, long* ret_stack_len, long u, long c", stmt->lhs.c_str());
  }
  else
    fprintf(output, "long yanshi_%s_transit(vector<long>& ret_stack, long u, long c", stmt->lhs.c_str());
  if (stmt->export_params.size())
    fprintf(output, ", %s", stmt->export_params.c_str());
fprintf(output,
")\n"
"{\n"
"  long v = -1;\n"
"again:\n"
"  switch (u) {\n");
  REP(u, anno.fsa.n()) {
    if (call_addr[u].first >= 0) { // no other transitions
      fprintf(output,
"  case %ld:\n"
"    u = %ld;\n"
, u, call_addr[u].first);
      if (opt_gen_c)
        fprintf(output,
"    if (*ret_stack_len >= %ld) return -1;\n"
"    ret_stack[(*ret_stack_len)++] = %ld;\n"
, opt_max_return_stack, call_addr[u].second);
      else
        fprintf(output,
"    ret_stack.push_back(%ld);\n"
, call_addr[u].second);
      fprintf(output,
"    goto again;\n");
      continue;
    }
    if (anno.fsa.adj[u].empty() && ! sub_final[u])
      continue;
    indent(output, 1);
    fprintf(output, "case %ld:\n", u);
    indent(output, 2);
    fprintf(output, "switch (c) {\n");

    unordered_map<long, pair<vector<pair<long, long>>, vector<pair<Action*, long>>>> v2case;
    for (auto it = anno.fsa.adj[u].begin(); it != anno.fsa.adj[u].end(); ) {
      long from = it->first.first, to = it->first.second, v = it->second;
      while (++it != anno.fsa.adj[u].end() && to == it->first.first && it->second == v)
        to = it->first.second;
      v2case[v].first.emplace_back(from, to);
      auto& body = v2case[v].second;

      auto ie = withins[u].end(), je = withins[v].end();

      // leaving = Expr(u) - Expr(v)
      for (auto i = withins[u].begin(), j = withins[v].begin(); i != ie; ++i) {
        while (j != je && i->first > j->first)
          ++j;
        if (j == je || i->first != j->first)
          for (auto action: i->first->leaving) {
            D("%%");
            body.push_back(action);
          }
      }

      // entering = Expr(v) - Expr(u)
      for (auto i = withins[u].begin(), j = withins[v].begin(); j != je; ++j) {
        while (i != ie && i->first < j->first)
          ++i;
        if (i == ie || i->first != j->first)
          for (auto action: j->first->entering) {
            D(">");
            body.push_back(action);
          }
      }

      // transiting = intersect(Expr(u), Expr(v))
      for (auto i = withins[u].begin(), j = withins[v].begin(); j != je; ++j) {
        while (i != ie && i->first < j->first)
          ++i;
        if (i != ie && i->first == j->first)
          for (auto action: j->first->transiting) {
            D("$");
            body.push_back(action);
          }
      }

      // finishing = intersect(Expr(u), Expr(v)) & Expr(v).has_final(v)
      for (auto i = withins[u].begin(), j = withins[v].begin(); j != je; ++j) {
        while (i != ie && i->first < j->first)
          ++i;
        if (i != ie && i->first == j->first && has_final(j->second))
          for (auto action: j->first->finishing) {
            D("@");
            body.push_back(action);
          }
      }
    }

    for (auto& x: v2case) {
      for (auto& y: x.second.first) {
        indent(output, 2);
        if (y.first == y.second-1)
          fprintf(output, "case %ld:\n", y.first);
        else
          fprintf(output, "case %ld ... %ld:\n", y.first, y.second-1);
      }
      indent(output, 3);
      fprintf(output, "v = %ld;\n", x.first);

      // actions
      sort(ALL(x.second.second), [](const pair<Action*, long>& a0, const pair<Action*, long>& a1) {
        return a0.second != a1.second ? a0.second < a1.second : a0.first < a1.first;
      });
      x.second.second.erase(unique(ALL(x.second.second)), x.second.second.end());
      for (auto a: x.second.second)
        fprintf(output, "{%s}\n", get_code(a.first).c_str());
      indent(output, 3);
      fprintf(output, "break;\n");
    }
    // return from finals of DefineStmt called by CallExpr
    if (sub_final[u]) {
      indent(output, 2);
      fprintf(output, "default:\n");
      indent(output, 3);
      fprintf(output, opt_gen_c ?
"if (*ret_stack_len) { u = ret_stack[--*ret_stack_len]; goto again; }\n"
:
"if (ret_stack.size()) { u = ret_stack.back(); ret_stack.pop_back(); goto again; }\n");
      indent(output, 3);
      fprintf(output, "break;\n");
    }

    indent(output, 2);
    fprintf(output, "}\n");
    indent(output, 2);
    fprintf(output, "break;\n");
  }
  indent(output, 1);
  fprintf(output, "}\n");
  indent(output, 1);
  fprintf(output, "return v;\n");
  fprintf(output, "}\n\n");
}

bool compile_export(DefineStmt* stmt)
{
  DP(2, "Exporting %s", stmt->lhs.c_str());
  FsaAnno& anno = compiled[stmt];

  DP(3, "Construct automaton with all DefineStmt associated to referenced CallExpr/CollapseExpr");
  vector<vector<Edge>> adj;
  decltype(anno.assoc) assoc;
  vector<vector<DefineStmt*>> cllps;
  long allo = 0;
  unordered_map<DefineStmt*, long> stmt2offset;
  unordered_map<DefineStmt*, long> stmt2start;
  unordered_map<long, DefineStmt*> start2stmt;
  vector<long> starts;
  vector<bool> sub_final;
  function<void(DefineStmt*)> allocate = [&](DefineStmt* stmt) {
    if (stmt2offset.count(stmt))
      return;
    DP(4, "Allocate %ld to %s", allo, stmt->lhs.c_str());
    FsaAnno& anno = compiled[stmt];
    long base = stmt2offset[stmt] = allo;
    allo += anno.fsa.n();
    sub_final.resize(allo);
    if (used_as_call.count(stmt)) {
      stmt2start[stmt] = base+anno.fsa.start;
      start2stmt[base+anno.fsa.start] = stmt;
      starts.push_back(base+anno.fsa.start);
      for (long f: anno.fsa.finals)
        sub_final[base+f] = true;
    }
    adj.insert(adj.end(), ALL(anno.fsa.adj));
    REP(i, anno.fsa.n())
      for (auto& e: adj[base+i])
        e.second += base;
    assoc.insert(assoc.end(), ALL(anno.assoc));
    FOR(i, base, base+anno.fsa.n()) {
      for (auto aa: assoc[i])
        if (has_start(aa.second)) {
          if (auto* e = dynamic_cast<CallExpr*>(aa.first)) {
            DefineStmt* v = e->define_stmt;
            allocate(v);
          } else if (auto* e = dynamic_cast<CollapseExpr*>(aa.first)) {
            DefineStmt* v = e->define_stmt;
            allocate(v);
            // (i@{CollapseExpr,...}, special, _) -> ({CollapseExpr,...}, epsilon, CollapseExpr.define_stmt.start)
            sorted_emplace(adj[i], epsilon, stmt2offset[v]+compiled[v].fsa.start);
          }
        }
      long j = adj[i].size();
      while (j && collapse_label_base < adj[i][j-1].first.second) {
        long v = adj[i][j-1].second;
        if (adj[i][j-1].first.first < collapse_label_base)
          adj[i][j-1].first.second = collapse_label_base;
        else
          j--;
        CollapseExpr* e;
        for (auto aa: assoc[v])
          if (has_final(aa.second) && (e = dynamic_cast<CollapseExpr*>(aa.first))) {
            DefineStmt* w = e->define_stmt;
            allocate(w);
            // (_, special, v@{CollapseExpr,...}) -> (CollapseExpr.define_stmt.final, epsilon, v)
            for (long f: compiled[w].fsa.finals) {
              long g = stmt2offset[w]+f;
              sorted_emplace(adj[g], epsilon, v);
              if (g == i)
                j++;
            }
          }
      }
      // remove (i, special, _)
      adj[i].resize(j);
    }
  };
  allocate(stmt);
  anno.fsa.adj = move(adj);
  anno.assoc = move(assoc);
  anno.deterministic = false;
  DP(3, "# of states: %ld", anno.fsa.n());

  // substring grammar & this nonterminal is not marked as intact
  if (opt_substring_grammar && ! stmt->intact) {
    DP(3, "Constructing substring grammar");
    anno.substring_grammar();
    DP(3, "# of states: %ld", anno.fsa.n());
  }

  vector<vector<long>> map0;
  DP(3, "Determinize");
  anno.determinize(&starts, &map0);
  vector<bool> sub_final2(anno.fsa.n());
  REP(i, anno.fsa.n())
    for (long u: map0[i]) {
      if (sub_final[u])
        sub_final2[i] = true;
      if (start2stmt.count(u)) {
        DefineStmt* stmt = start2stmt[u];
        if (stmt2start[stmt] < 0) {
          stmt->module->locfile.error_context(stmt->loc, "the start has been included in multiple DFA states");
          return false;
        }
        stmt2start[stmt] = ~ i;
      }
    }
  sub_final = move(sub_final2);
  start2stmt.clear();
  for (auto& it: stmt2start) {
    it.second = ~ it.second;
    start2stmt[it.second] = it.first;
  }
  DP(3, "# of states: %ld", anno.fsa.n());

  DP(3, "Minimize");
  map0.clear();
  anno.minimize(&map0);
  sub_final2.assign(anno.fsa.n(), false);
  REP(i, anno.fsa.n())
    for (long u: map0[i]) {
      if (sub_final[u])
        sub_final2[i] = true;
      if (start2stmt.count(u)) {
        DefineStmt* stmt = start2stmt[u];
        stmt2start[stmt] = i;
      }
    }
  sub_final = move(sub_final2);
  start2stmt.clear();
  for (auto& it: stmt2start)
    start2stmt[it.second] = it.first;
  DP(3, "# of states: %ld", anno.fsa.n());

  if (! opt_keep_inaccessible) {
    DP(3, "Keep accessible states");
    // roots: start, starts of DefineStmt associated to CallExpr
    starts.clear();
    for (auto& it: stmt2start)
      starts.push_back(it.second);
    vector<long> map1;
    anno.accessible(&starts, map1);
    sub_final2.assign(anno.fsa.n(), false);
    REP(i, anno.fsa.n()) {
      long u = map1[i];
      sub_final2[i] = sub_final[u];
      if (start2stmt.count(u))
        stmt2start[start2stmt[u]] = i;
    }
    sub_final = move(sub_final2);
    start2stmt.clear();
    for (auto& it: stmt2start)
      start2stmt[it.second] = it.first;
    DP(3, "# of states: %ld", anno.fsa.n());

    DP(3, "Keep co-accessible states");
    // roots: finals, finals of DefineStmt associated to CallExpr
    map1.clear();
    anno.co_accessible(&sub_final, map1);
    sub_final2.assign(anno.fsa.n(), false);
    REP(i, anno.fsa.n()) {
      long u = map1[i];
      sub_final2[i] = sub_final[u];
      if (start2stmt.count(u))
        stmt2start[start2stmt[u]] = i;
    }
    sub_final = move(sub_final2);
    start2stmt.clear();
    for (auto& it: stmt2start)
      start2stmt[it.second] = it.first;
    DP(3, "# of states: %ld", anno.fsa.n());
  }

  stmt2final[stmt] = sub_final;
  auto& call_addr = stmt2call_addr[stmt];
  call_addr.assign(anno.fsa.n(), make_pair(-1L, -1L));
  DP(3, "CallExpr");
  REP(i, anno.fsa.n())
    if (anno.fsa.has_call(i)) {
      if (anno.fsa.adj[i].size() != 1 || anno.fsa.adj[i][0].first.second-anno.fsa.adj[i][0].first.first > 1) {
        stmt->module->locfile.error_context(stmt->loc, "state %ld: CallExpr cannot coexist with other transitions", i);
        for (auto it = anno.fsa.adj[i].begin(); it != anno.fsa.adj[i].end(); ) {
          long from = it->first.first, to = it->first.second, v = it->second;
          while (++it != anno.fsa.adj[i].end() && to == it->first.first && it->second == v)
            to = it->first.second;
          fprintf(stderr, "  (%ld,%ld)\n", from, to-1);
        }
        return false;
      }
      for (auto aa: anno.assoc[i])
        if (has_start(aa.second))
          if (auto* e = dynamic_cast<CallExpr*>(aa.first)) // unique
            call_addr[i] = {stmt2start[e->define_stmt], anno.fsa.adj[i][0].second};
    }

  DP(3, "Removing action/CallExpr labels");
  REP(i, anno.fsa.n()) {
    long j = anno.fsa.adj[i].size();
    while (j && action_label_base < anno.fsa.adj[i][j-1].first.second)
      if (anno.fsa.adj[i][j-1].first.first < action_label_base)
        anno.fsa.adj[i][j-1].first.second = action_label_base;
      else
        j--;
    anno.fsa.adj[i].resize(j);
  }

  return true;
}

//// Graphviz dot renderer

void generate_graphviz(Module* mo)
{
  fprintf(output, "// Generated by 偃师, %s\n", mo->filename.c_str());
  for (Stmt* x = mo->toplevel; x; x = x->next)
    if (auto stmt = dynamic_cast<DefineStmt*>(x)) {
      if (stmt->export_) {
        FsaAnno& anno = compiled[stmt];

        fprintf(output, "digraph \"%s\" {\n", mo->filename.c_str());
        bool start_is_final = false;

        // finals
        indent(output, 1);
        fprintf(output, "node[shape=doublecircle,color=olivedrab1,style=filled,fontname=Monospace];");
        for (long f: anno.fsa.finals)
          if (f == anno.fsa.start)
            start_is_final = true;
          else
            fprintf(output, " %ld", f);
        fprintf(output, "\n");

        // start
        indent(output, 1);
        if (start_is_final)
          fprintf(output, "node[shape=doublecircle,color=orchid];");
        else
          fprintf(output, "node[shape=circle,color=orchid];");
        fprintf(output, " %ld\n", anno.fsa.start);

        // other states
        indent(output, 1);
        fprintf(output, "node[shape=circle,color=black,style=\"\"]\n");

        // edges
        REP(u, anno.fsa.n()) {
          unordered_map<long, stringstream> labels;
          bool first = true;
          auto it = anno.fsa.adj[u].begin();
          for (; it != anno.fsa.adj[u].end(); ++it) {
            stringstream& lb = labels[it->second];
            if (! lb.str().empty())
              lb << ',';
            if (it->first.first == it->first.second-1)
              lb << it->first.first;
            else
              lb << it->first.first << '-' << it->first.second-1;
          }
          for (auto& lb: labels) {
            indent(output, 1);
            fprintf(output, "%ld -> %ld[label=\"%s\"]\n", u, lb.first, lb.second.str().c_str());
          }
        }
      }
    }
  fprintf(output, "}\n");
}

//// C++ renderer

static void generate_final(const char* name, const vector<bool>& final)
{
  // comment
  fprintf(output, "  //static const long %sfinals[] = {", name);
  bool first = true;
  REP(i, final.size())
    if (final[i]) {
      if (first) first = false;
      else fprintf(output, ",");
      fprintf(output, "%ld", i);
    }
  fprintf(output, "};\n");

  first = true;
  fprintf(output, "  static const unsigned long %sfinal[] = {", name);
  for (long j = 0, i = 0; i < final.size(); i += CHAR_BIT*sizeof(long)) {
      ulong mask = 0;
      for (; j < final.size() && j < i+CHAR_BIT*sizeof(long); j++)
        if (final[j]) {
          mask |= 1uL << (j-i);
        }
      if (i) fprintf(output, ",");
      fprintf(output, "%#lx", mask);
  }
  fprintf(output, "};\n");
}

static void generate_cxx_export(DefineStmt* stmt)
{
  FsaAnno& anno = compiled[stmt];

  // yanshi_%s_init
  if (output_header)
    fprintf(output_header, "extern long yanshi_%s_start;\n", stmt->lhs.c_str());
  fprintf(output, "long yanshi_%s_start = %ld;\n\n", stmt->lhs.c_str(), anno.fsa.start);

  // yanshi_%s_is_final
  if (output_header) {
    if (opt_gen_extern_c) fputs("extern \"C\" ", output_header);
    fprintf(output_header, opt_gen_c ?
"bool yanshi_%s_is_final(const long* ret_stack, long ret_stack_len, long u);\n"
:
"bool yanshi_%s_is_final(const vector<long>& ret_stack, long u);\n"
, stmt->lhs.c_str());
  }
  if (opt_gen_extern_c) fputs("extern \"C\" ", output);
  fprintf(output, opt_gen_c ?
"bool yanshi_%s_is_final(const long* ret_stack, long ret_stack_len, long u)\n"
:
"bool yanshi_%s_is_final(const vector<long>& ret_stack, long u)\n"
          , stmt->lhs.c_str());
  fprintf(output, "{\n");
  vector<bool> final(anno.fsa.n());
  for (long f: anno.fsa.finals)
    final[f] = true;
  generate_final("", final);
  generate_final("sub_", stmt2final[stmt]);
  fprintf(output, opt_gen_c ?
"  for (long i = ret_stack_len; i; u = ret_stack[--i])\n"
:
"  for (auto i = ret_stack.size(); i; u = ret_stack[--i])\n"
);
  fprintf(output,
"    if (! (0 <= u && u < %ld && sub_final[u/(CHAR_BIT*sizeof(long))] >> (u%%(CHAR_BIT*sizeof(long))) & 1))\n"
"      return false;\n"
"  return 0 <= u && u < %ld && final[u/(CHAR_BIT*sizeof(long))] >> (u%%(CHAR_BIT*sizeof(long))) & 1;\n"
"};\n\n"
, anno.fsa.n() , anno.fsa.n()
);
  generate_transitions(stmt);
}

void generate_cxx(Module* mo)
{
  fprintf(output, "// Generated by 偃师, %s\n", mo->filename.c_str());
  fprintf(output, "#include <limits.h>\n");
  if (! opt_gen_c) {
    fprintf(output, "#include <vector>\n");
    fprintf(output, "using namespace std;\n");
  } else {
    fprintf(output, "#include <stdbool.h>\n");
  }
  if (opt_standalone) {
    fputs(
"#include <algorithm>\n"
"#include <cinttypes>\n"
"#include <clocale>\n"
"#include <codecvt>\n"
"#include <cstdint>\n"
"#include <cstdio>\n"
"#include <cstring>\n"
"#include <cwctype>\n"
"#include <iostream>\n"
"#include <locale>\n"
"#include <string>\n"
"using namespace std;\n"
, output);
  }
  if (output_header) {
    fputs("#pragma once\n", output_header);
    if (! opt_gen_c) {
      fprintf(output_header, "#include <vector>\n");
      fprintf(output_header, "using std::vector;\n");
    } else {
      fprintf(output_header, "#include <stdbool.h>\n");
    }
  }
  fprintf(output, "\n");
  DefineStmt* main_export = NULL;
  for (Stmt* x = mo->toplevel; x; x = x->next)
    if (auto xx = dynamic_cast<DefineStmt*>(x)) {
      if (xx->export_) {
        if (! main_export)
          main_export = xx;
        generate_cxx_export(xx);
      }
    } else if (auto xx = dynamic_cast<CppStmt*>(x))
      fprintf(output, "%s", xx->code.c_str());
  if (opt_standalone && main_export) {
    fprintf(output,
"\n"
"int main(int argc, char* argv[])\n"
"{\n"
"  setlocale(LC_ALL, \"\");\n"
"  string utf8;\n"
"  const char* p;\n"
"  long c, u = yanshi_%s_start, pref = 0;\n"
, main_export->lhs.c_str());
    if (opt_gen_c)
      fprintf(output, "  long ret_stack[%ld], ret_stack_len = 0;\n", opt_max_return_stack);
    else
      fprintf(output, "  vector<long> ret_stack;\n");
    fprintf(output,
"  if (argc == 2)\n"
"    utf8 = argv[1];\n"
"  else {\n"
"    FILE* f = argc == 1 ? stdin : fopen(argv[1], \"r\");\n"
"    while ((c = fgetc(f)) != EOF)\n"
"      utf8 += c;\n"
"    fclose(f);\n"
"  }\n"
"  u32string utf32 = wstring_convert<codecvt_utf8<char32_t>, char32_t>{}.from_bytes(utf8);\n"
);
    fprintf(output, opt_gen_c ?
"  printf(\"\\033[%%s33m%%ld \\033[m\", yanshi_%s_is_final(ret_stack, ret_stack_len, u) ? \"1;\" : \"\", u);\n"
:
"  printf(\"\\033[%%s33m%%ld \\033[m\", yanshi_%s_is_final(ret_stack, u) ? \"1;\" : \"\", u);\n"
, main_export->lhs.c_str()
);
    fprintf(output,
"  for (char32_t c: utf32) {\n");
    fprintf(output, opt_gen_c ?
"    u = yanshi_%s_transit(ret_stack, &ret_stack_len, u, c);\n"
:
"    u = yanshi_%s_transit(ret_stack, u, c);\n"
, main_export->lhs.c_str());
    fprintf(output,
"    if (c > WCHAR_MAX || iswcntrl(c)) printf(\"%%\" PRIuLEAST32 \" \", c);\n"
"    else cout << wstring_convert<codecvt_utf8<char32_t>, char32_t>{}.to_bytes(c) << ' ';\n");
    fprintf(output, opt_gen_c ?
"    printf(\"\\033[%%s33m%%ld \\033[m\", yanshi_%s_is_final(ret_stack, ret_stack_len, u) ? \"1;\" : \"\", u);\n"
:
"    printf(\"\\033[%%s33m%%ld \\033[m\", yanshi_%s_is_final(ret_stack, u) ? \"1;\" : \"\", u);\n"
, main_export->lhs.c_str());
    fprintf(output,
"    if (u < 0) break;\n"
"    pref++;\n"
"  }\n");
    fprintf(output, opt_gen_c ?
"  printf(\"\\nlen: %%zd\\npref: %%ld\\nstate: %%ld\\nfinal: %%s\\n\", utf32.size(), pref, u, yanshi_%s_is_final(ret_stack, ret_stack_len, u) ? \"true\" : \"false\");\n"
"}\n"
:
"  printf(\"\\nlen: %%zd\\npref: %%ld\\nstate: %%ld\\nfinal: %%s\\n\", utf32.size(), pref, u, yanshi_%s_is_final(ret_stack, u) ? \"true\" : \"false\");\n"
"}\n"
, main_export->lhs.c_str());
  }
}
