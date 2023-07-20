#include "common.hh"
#include "compiler.hh"
#include "loader.hh"
#include "option.hh"
#include "parser.hh"
#include "repl.hh"

#include <algorithm>
#include <errno.h>
#include <functional>
#include <stdio.h>
#include <stack>
#include <string.h>
#include <sys/stat.h>
#include <sysexits.h>
#include <unordered_map>
#include <unordered_set>
#include <utility>
using namespace std;

static map<pair<dev_t, ino_t>, Module> inode2module;
static unordered_map<DefineStmt*, vector<DefineStmt*>> depended_by; // key ranges over all DefineStmt
map<DefineStmt*, vector<Expr*>> used_as_call, used_as_collapse, used_as_embed;
static DefineStmt* main_export;
Module* main_module;
FILE *output, *output_header;

void print_module_info(Module& mo)
{
  yellow(); printf("filename: %s\n", mo.filename.c_str());
  cyan(); puts("qualified imports:"); sgr0();
  for (auto& x: mo.qualified_import)
    printf("  %s as %s\n", x.second->filename.c_str(), x.first.c_str());
  cyan(); puts("unqualified imports:"); sgr0();
  for (auto& x: mo.unqualified_import)
    printf("  %s\n", x->filename.c_str());
  cyan(); puts("defined actions:"); sgr0();
  for (auto& x: mo.defined_action)
    printf("  %s\n", x.first.c_str());
  cyan(); puts("defined:"); sgr0();
  for (auto& x: mo.defined)
    printf("  %s\n", x.first.c_str());
}

Stmt* resolve(Module& mo, const string qualified, const string& ident)
{
  if (qualified.size()) {
    if (! mo.qualified_import.count(qualified))
      return NULL;
    auto it = mo.qualified_import[qualified]->defined.find(ident);
    if (it == mo.qualified_import[qualified]->defined.end())
      return NULL;
    return it->second;
  } else {
    Stmt* r = NULL;
    if (mo.macro.count(ident))
      r = mo.macro[ident];
    if (mo.defined.count(ident)) {
      if (r) return (Stmt*)1;
      r = mo.defined[ident];
    }
    for (auto* import: mo.unqualified_import) {
      if (import->macro.count(ident)) {
        if (r) return (Stmt*)1;
        r = import->macro[ident];
      }
      if (import->defined.count(ident)) {
        if (r) return (Stmt*)1;
        r = import->defined[ident];
      }
    }
    return r;
  }
}

ActionStmt* resolve_action(Module& mo, const string qualified, const string& ident)
{
  if (qualified.size()) {
    if (! mo.qualified_import.count(qualified))
      return NULL;
    auto it = mo.qualified_import[qualified]->defined_action.find(ident);
    if (it == mo.qualified_import[qualified]->defined_action.end())
      return NULL;
    return it->second;
  } else {
    ActionStmt* r = NULL;
    if (mo.defined_action.count(ident))
      r = mo.defined_action[ident];
    for (auto* import: mo.unqualified_import)
      if (import->defined_action.count(ident)) {
        if (r) return (ActionStmt*)1;
        r = import->defined_action[ident];
      }
    return r;
  }
}

struct ModuleImportDef : PreorderStmtVisitor {
  Module& mo;
  long& n_errors;
  ModuleImportDef(Module& mo, long& n_errors) : mo(mo), n_errors(n_errors) {}

  void visit(ActionStmt& stmt) override {
    if (mo.defined_action.count(stmt.ident)) {
      n_errors++;
      mo.locfile.error(stmt.loc, "redefined '%s'", stmt.ident.c_str());
    } else
      mo.defined_action[stmt.ident] = &stmt;
  }
  // TODO report error: import 'aa.hs' (#define d 3) ; #define d 4
  void visit(DefineStmt& stmt) override {
    if (mo.defined.count(stmt.lhs) || mo.macro.count(stmt.lhs)) {
      n_errors++;
      mo.locfile.error(stmt.loc, "redefined '%s'", stmt.lhs.c_str());
    } else {
      mo.defined.emplace(stmt.lhs, &stmt);
      stmt.module = &mo;
      depended_by[&stmt]; // empty
    }
  }
  void visit(ImportStmt& stmt) override {
    Module* m = load_module(n_errors, stmt.filename);
    if (! m) {
      n_errors++;
      mo.locfile.error(stmt.loc, "'%s': %s", stmt.filename.c_str(), errno ? strerror(errno) : "parse error");
      return;
    }
    if (stmt.qualified.size())
      mo.qualified_import[stmt.qualified] = m;
    else if (count(ALL(mo.unqualified_import), m) == 0)
      mo.unqualified_import.push_back(m);
  }
  void visit(PreprocessDefineStmt& stmt) override {
    if (mo.defined.count(stmt.ident) || mo.macro.count(stmt.ident)) {
      n_errors++;
      mo.locfile.error(stmt.loc, "redefined '%s'", stmt.ident.c_str());
    } else
      mo.macro[stmt.ident] = &stmt;
  }
};

struct ModuleUse : PrePostActionExprStmtVisitor {
  Module& mo;
  long& n_errors;
  DefineStmt* stmt = NULL;
  ModuleUse(Module& mo, long& n_errors) : mo(mo), n_errors(n_errors) {}

  void pre_expr(Expr& expr) override {
    expr.stmt = stmt;
  }

  void post_expr(Expr& expr) override {
    for (auto a: expr.entering)
      PrePostActionExprStmtVisitor::visit(*a.first);
    for (auto a: expr.finishing)
      PrePostActionExprStmtVisitor::visit(*a.first);
    for (auto a: expr.leaving)
      PrePostActionExprStmtVisitor::visit(*a.first);
    for (auto a: expr.transiting)
      PrePostActionExprStmtVisitor::visit(*a.first);
  }

  void visit(RefAction& action) override {
    ActionStmt* r = resolve_action(mo, action.qualified, action.ident);
    if (! r) {
      n_errors++;
      if (action.qualified.size())
        mo.locfile.error(action.loc, "'%s::%s' undefined", action.qualified.c_str(), action.ident.c_str());
      else
        mo.locfile.error(action.loc, "'%s' undefined", action.ident.c_str());
    } else if (r == (Stmt*)1) {
      n_errors++;
      mo.locfile.error(action.loc, "'%s' redefined", action.ident.c_str());
    } else
      action.define_stmt = r;
  }

  void visit(BracketExpr& expr) override {
    for (auto& x: expr.intervals.to)
      AB = max(AB, x.second);
  }
  void visit(CallExpr& expr) override {
    Stmt* r = resolve(mo, expr.qualified, expr.ident);
    if (! r)
      error_undefined(expr.loc, expr.qualified, expr.ident);
    else if (r == (Stmt*)1)
      error_ambiguous(expr.loc, expr.ident);
    else if (auto d = dynamic_cast<PreprocessDefineStmt*>(r))
      error_misuse_macro("CallExpr", expr.loc, expr.qualified, expr.ident);
    else if (auto d = dynamic_cast<DefineStmt*>(r)) {
      used_as_call[d].push_back(&expr);
      expr.define_stmt = d;
    } else
      assert(0);
  }
  void visit(CollapseExpr& expr) override {
    Stmt* r = resolve(mo, expr.qualified, expr.ident);
    if (! r)
      error_undefined(expr.loc, expr.qualified, expr.ident);
    else if (r == (Stmt*)1)
      error_ambiguous(expr.loc, expr.ident);
    else if (auto d = dynamic_cast<PreprocessDefineStmt*>(r))
      error_misuse_macro("CollapseExpr", expr.loc, expr.qualified, expr.ident);
    else if (auto d = dynamic_cast<DefineStmt*>(r)) {
      used_as_collapse[d].push_back(&expr);
      expr.define_stmt = d;
    } else
      assert(0);
  }
  void visit(DefineStmt& stmt) override {
    this->stmt = &stmt;
    PrePostActionExprStmtVisitor::visit(*stmt.rhs);
    this->stmt = NULL;
  }
  void visit(EmbedExpr& expr) override {
    // introduce dependency
    Stmt* r = resolve(mo, expr.qualified, expr.ident);
    if (! r)
      error_undefined(expr.loc, expr.qualified, expr.ident);
    else if (r == (Stmt*)1)
      error_ambiguous(expr.loc, expr.ident);
    else if (auto d = dynamic_cast<PreprocessDefineStmt*>(r)) {
      // enlarge alphabet
      expr.define_stmt = NULL;
      expr.macro_value = d->value;
      AB = max(AB, d->value+1);
    } else if (auto d = dynamic_cast<DefineStmt*>(r)) {
      depended_by[d].push_back(stmt);
      used_as_embed[d].push_back(&expr);
      expr.define_stmt = d;
    } else
      assert(0);
  }
private:
  void error_undefined(const Location& loc, const string& qualified, const string& ident) {
    n_errors++;
    if (qualified.size())
      mo.locfile.error(loc, "'%s::%s' undefined", qualified.c_str(), ident.c_str());
    else
      mo.locfile.error(loc, "'%s' undefined", ident.c_str());
  }
  void error_ambiguous(const Location& loc, const string& ident) {
    n_errors++;
    mo.locfile.error(loc, "ambiguous '%s'", ident.c_str());
  }
  void error_misuse_macro(const char* name, const Location& loc, const string& qualified, const string& ident) {
    n_errors++;
    if (qualified.size())
      mo.locfile.error(loc, "macro '%s::%s' used as %s", qualified.c_str(), ident.c_str(), name);
    else
      mo.locfile.error(loc, "macro '%s' used as %s", ident.c_str(), name);
  }
};

Module* load_module(long& n_errors, const string& filename)
{
  FILE* file = stdin;
  if (filename != "-") {
    file = fopen(filename.c_str(), "r");
    for (string& include: opt_include_paths) {
      if (file) break;
      file = fopen((include+'/'+filename).c_str(), "r");
    }
  }
  if (! file) {
    n_errors++;
    return NULL;
  }

  pair<dev_t, ino_t> inode{0, 0}; // stdin -> {0, 0}
  if (file != stdin) {
    struct stat sb;
    if (fstat(fileno(file), &sb) < 0)
      err_exit(EX_OSFILE, "fstat '%s'", filename.c_str());
    inode = {sb.st_dev, sb.st_ino};
  }
  if (inode2module.count(inode)) {
    fclose(file);
    return &inode2module[inode];
  }
  Module& mo = inode2module[inode];

  string module{file != stdin ? filename : "main"};
  string::size_type t = module.find('.');
  if (t != string::npos)
    module.erase(t, module.size()-t);

  long r;
  char buf[BUF_SIZE];
  string data;
  while ((r = fread(buf, 1, sizeof buf, file)) > 0) {
    data += string(buf, buf+r);
    if (r < sizeof buf) break;
  }
  fclose(file);
  if (data.empty() || data.back() != '\n')
    data.push_back('\n');
  LocationFile locfile(filename, data);

  Stmt* toplevel = NULL;
  mo.locfile = locfile;
  mo.filename = filename;
  long errors = parse(locfile, toplevel);
  if (! toplevel) {
    n_errors += errors;
    mo.status = BAD;
    mo.toplevel = NULL;
    return &mo;
  }
  mo.toplevel = toplevel;
  return &mo;
}

static vector<DefineStmt*> topo_define_stmts(long& n_errors)
{
  vector<DefineStmt*> topo;
  vector<DefineStmt*> st;
  unordered_map<DefineStmt*, i8> vis; // 0: unvisited; 1: in stack; 2: visited; 3: in a cycle
  unordered_map<DefineStmt*, long> cnt;
  function<bool(DefineStmt*)> dfs = [&](DefineStmt* u) {
    if (vis[u] == 2)
      return false;
    if (vis[u] == 3)
      return true;
    if (vis[u] == 1) {
      u->module->locfile.error_context(u->loc, "'%s': circular embedding", u->lhs.c_str());
      long i = st.size();
      while (st[i-1] != u)
        i--;
      st.push_back(st[i-1]);
      for (; i < st.size(); i++) {
        vis[st[i]] = 3;
        fputs("  ", stderr);
        st[i]->module->locfile.error_context(st[i]->loc, "required by %s", st[i]->lhs.c_str());
      }
      fputs("\n", stderr);
      return true;
    }
    cnt[u] = u->export_ ? 1 : 0;
    vis[u] = 1;
    st.push_back(u);
    bool cycle = false;
    for (auto v: depended_by[u])
      if (dfs(v))
        cycle = true;
      else
        cnt[u] += cnt[v];
    st.pop_back();
    vis[u] = 2;
    topo.push_back(u);
    return cycle;
  };
  for (auto& d: depended_by)
    if (! vis[d.first] && dfs(d.first)) // detected cycle
      n_errors++;
  reverse(ALL(topo));
  if (opt_dump_embed) {
    magenta(); printf("=== Embed\n"); sgr0();
    for (auto stmt: topo)
      if (cnt[stmt] > 0)
        printf("count(%s::%s) = %ld\n", stmt->module->filename.c_str(), stmt->lhs.c_str(), cnt[stmt]);
  }
  return topo;
}

long load(const string& filename)
{
  long n_errors = 0;
  Module* mo = load_module(n_errors, filename);
  if (! mo) {
    err_exit(EX_OSFILE, "fopen", filename.c_str());
    return n_errors;
  }
  if (mo->status == BAD)
    return n_errors;
  main_module = mo;

  DP(1, "Processing import & def");
  for(;;) {
    bool done = true;
    for (auto& it: inode2module)
      if (it.second.status == UNPROCESSED) {
        done = false;
        Module& mo = it.second;
        mo.status = GOOD;
        long old = n_errors;
        ModuleImportDef p{mo, n_errors};
        for (Stmt* x = mo.toplevel; x; x = x->next)
          x->accept(p);
        mo.status = old == n_errors ? GOOD : BAD;
      }
    if (done) break;
  }
  if (n_errors)
    return n_errors;

  DP(1, "Processing use");
  for (auto& it: inode2module)
    if (it.second.status == GOOD) {
      Module& mo = it.second;
      ModuleUse p{mo, n_errors};
      for (Stmt* x = mo.toplevel; x; x = x->next)
        x->accept(p);
    }
  if (n_errors)
    return n_errors;

  // warning: not used solely as CallExpr, CollapseExpr or EmbedExpr
  {
    auto it0 = used_as_call.begin(), it0e = used_as_call.end(),
         it1 = used_as_collapse.begin(), it1e = used_as_collapse.end(),
         it2 = used_as_embed.begin(), it2e = used_as_embed.end();
    while (it0 != it0e || it1 != it1e || it2 != it2e) {
      long k = 0;
      long c = 0;
      DefineStmt* x = NULL;
      if (it0 != it0e && (! x || it0->first < x)) x = it0->first;
      if (it1 != it1e && (! x || it1->first < x)) x = it1->first;
      if (it2 != it2e && (! x || it2->first < x)) x = it2->first;
      if (it0 != it0e && it0->first == x) c++;
      if (it1 != it1e && it1->first == x) c++;
      if (it2 != it2e && it2->first == x) c++;
      if (c > 1) {
        x->module->locfile.warning(x->loc, "'%s' is not used solely as CallExpr, CollapseExpr or EmbedExpr", x->lhs.c_str());
        if (it0 != it0e && it0->first == x)
          for (auto* y: it0->second) {
            fputs("  ", stderr);
            y->stmt->module->locfile.warning_context(y->loc, "required by %s", y->stmt->lhs.c_str());
          }
        if (it1 != it1e && it1->first == x)
          for (auto* y: it1->second) {
            fputs("  ", stderr);
            y->stmt->module->locfile.warning_context(y->loc, "required by %s", y->stmt->lhs.c_str());
          }
        if (it2 != it2e && it2->first == x)
          for (auto* y: it2->second) {
            fputs("  ", stderr);
            y->stmt->module->locfile.warning_context(y->loc, "required by %s", y->stmt->lhs.c_str());
          }
      }
      if (it0 != it0e && it0->first == x) ++it0, c++;
      if (it1 != it1e && it1->first == x) ++it1, c++;
      if (it2 != it2e && it2->first == x) ++it2, c++;
    }
  }

  if (opt_dump_module) {
    magenta(); printf("=== Module\n"); sgr0();
    for (auto& it: inode2module)
      if (it.second.status == GOOD) {
        Module& mo = it.second;
        print_module_info(mo);
      }
    puts("");
  }

  if (opt_dump_tree) {
    magenta(); printf("=== Tree\n"); sgr0();
    StmtPrinter p;
    for (auto& it: inode2module)
      if (it.second.status == GOOD) {
        Module& mo = it.second;
        yellow(); printf("filename: %s\n", mo.filename.c_str()); sgr0();
        for (Stmt* x = mo.toplevel; x; x = x->next)
          x->accept(p);
      }
    puts("");
  }

  DP(1, "Topological sorting");
  vector<DefineStmt*> topo = topo_define_stmts(n_errors);
  if (n_errors)
    return n_errors;

  if (opt_check)
    return 0;

  // AB has been updated by ModuleUse
  action_label_base = action_label = AB;
  call_label_base = call_label = action_label+1000000;
  collapse_label_base = collapse_label = call_label+1000000;

  DP(1, "Compiling DefineStmt");
  for (auto stmt: topo)
    compile(stmt);

  output = strcmp(opt_output_filename, "-") ? fopen(opt_output_filename, "w") : stdout;
  if (! output) {
    n_errors++;
    err_exit(EX_OSFILE, "fopen", opt_output_filename);
    return n_errors;
  }

  unordered_map<DefineStmt*, vector<pair<long, long>>> stmt2call_addr;
  DP(1, "Compiling exporting DefineStmt (coalescing referenced CallExpr/CollapseExpr)");
  for (Stmt* x = main_module->toplevel; x; x = x->next)
    if (auto xx = dynamic_cast<DefineStmt*>(x))
      if (xx->export_ && ! compile_export(xx))
        n_errors++;
  if (n_errors)
    return n_errors;

  for (Stmt* x = main_module->toplevel; x; x = x->next)
    if (auto xx = dynamic_cast<DefineStmt*>(x))
      if (xx->export_) {
        FsaAnno& anno = compiled[xx];
        if (opt_dump_automaton)
          print_automaton(anno.fsa);
        if (opt_dump_assoc)
          print_assoc(anno);
      }

  if (opt_mode == Mode::cxx) {
    if (opt_output_header_filename) {
      output_header = fopen(opt_output_header_filename, "w");
      if (! output_header) {
        n_errors++;
        err_exit(EX_OSFILE, "fopen", opt_output_header_filename);
        return n_errors;
      }
    }
    DP(1, "Generating C++");
    generate_cxx(mo);
    if (output_header)
      fclose(output_header);
  } else if (opt_mode == Mode::graphviz) {
    DP(1, "Generating Graphviz dot");
    generate_graphviz(mo);
  } else if (opt_mode == Mode::interactive) {
    DP(1, "Testing given string");
    DefineStmt* main_export = NULL;
    for (Stmt* x = main_module->toplevel; x; x = x->next)
      if (auto xx = dynamic_cast<DefineStmt*>(x))
        if (xx->export_) {
          main_export = xx;
          break;
        }
    if (! main_export)
      puts("no exporting DefineStmt");
    else {
      printf("Testing %s\n", main_export->lhs.c_str());
      repl(main_export);
    }
  }

  fclose(output);
  return n_errors;
}

void unload_all()
{
  for (auto& it: inode2module) {
    Module& mo = it.second;
    stmt_free(mo.toplevel);
  }
}
