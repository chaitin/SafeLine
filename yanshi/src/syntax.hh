#pragma once
#include "common.hh"
#include "location.hh"

#include <cxxabi.h>
#include <memory>
#include <string.h>
#include <string>
#include <typeinfo>
#include <vector>
using std::move;
using std::pair;
using std::string;
using std::vector;

//// Visitor

template<class T>
struct Visitor;

template<class T>
struct VisitableBase {
  virtual void accept(Visitor<T>& visitor) = 0;
};

template<class Base, class Derived>
struct Visitable : Base {
  void accept(Visitor<Base>& visitor) override {
    visitor.visit(static_cast<Derived&>(*this));
  }
};

struct Action;
struct InlineAction;
struct RefAction;
template<>
struct Visitor<Action> {
  virtual void visit(Action& action) = 0;
  virtual void visit(InlineAction&) = 0;
  virtual void visit(RefAction&) = 0;
};

struct Expr;
struct BracketExpr;
struct CallExpr;
struct CollapseExpr;
struct ComplementExpr;
struct ConcatExpr;
struct DifferenceExpr;
struct DotExpr;
struct EmbedExpr;
struct EpsilonExpr;
struct IntersectExpr;
struct LiteralExpr;
struct PlusExpr;
struct RepeatExpr;
struct QuestionExpr;
struct StarExpr;
struct UnionExpr;
template<>
struct Visitor<Expr> {
  virtual void visit(Expr&) = 0;
  virtual void visit(BracketExpr&) = 0;
  virtual void visit(CallExpr&) = 0;
  virtual void visit(CollapseExpr&) = 0;
  virtual void visit(ComplementExpr&) = 0;
  virtual void visit(ConcatExpr&) = 0;
  virtual void visit(DifferenceExpr&) = 0;
  virtual void visit(DotExpr&) = 0;
  virtual void visit(EmbedExpr&) = 0;
  virtual void visit(EpsilonExpr&) = 0;
  virtual void visit(IntersectExpr&) = 0;
  virtual void visit(LiteralExpr&) = 0;
  virtual void visit(PlusExpr&) = 0;
  virtual void visit(RepeatExpr&) = 0;
  virtual void visit(QuestionExpr&) = 0;
  virtual void visit(StarExpr&) = 0;
  virtual void visit(UnionExpr&) = 0;
};

struct Stmt;
struct ActionStmt;
struct CppStmt;
struct DefineStmt;
struct EmptyStmt;
struct ImportStmt;
struct PreprocessDefineStmt;
template<>
struct Visitor<Stmt> {
  virtual void visit(Stmt&) = 0;
  virtual void visit(ActionStmt&) = 0;
  virtual void visit(CppStmt&) = 0;
  virtual void visit(DefineStmt&) = 0;
  virtual void visit(EmptyStmt&) = 0;
  virtual void visit(ImportStmt&) = 0;
  virtual void visit(PreprocessDefineStmt&) = 0;
};

//// Action

struct Action : VisitableBase<Action> {
  Location loc;
  virtual ~Action() = default;
};

struct InlineAction : Visitable<Action, InlineAction> {
  string code;
  InlineAction(string& code) : code(move(code)) {}
};

struct Module;
struct RefAction : Visitable<Action, RefAction> {
  string qualified, ident;
  ActionStmt* define_stmt; // set by ModuleUse
  RefAction(string& qualified, string& ident) : qualified(move(qualified)), ident(move(ident)) {}
};

//// Expr

struct Expr : VisitableBase<Expr> {
  Location loc;
  long pre, post, depth; // set by Compiler
  vector<Expr*> anc; // set by Compiler
  vector<pair<Action*, long>> entering, finishing, leaving, transiting;
  DefineStmt* stmt = NULL; // set by ModuleImportDef
  virtual ~Expr() {
    for (auto a: entering)
      delete a.first;
    for (auto a: finishing)
      delete a.first;
    for (auto a: leaving)
      delete a.first;
    for (auto a: transiting)
      delete a.first;
  }
  string name() const {
    int status;
    std::unique_ptr<char, void(*)(void*)> r{
      abi::__cxa_demangle(typeid(*this).name(), NULL, NULL, &status),
      free
    };
    std::string t = r.get();
    t = t.substr(0, t.size()-4); // suffix 'Expr'
    return t;
  }
  bool no_action() const {
    return entering.empty() && finishing.empty() && leaving.empty() && transiting.empty();
  }
};

struct BracketExpr : Visitable<Expr, BracketExpr> {
  DisjointIntervals intervals;
  BracketExpr(DisjointIntervals* intervals) : intervals(std::move(*intervals)) { delete intervals; }
};

struct CallExpr : Visitable<Expr, CallExpr> {
  string qualified, ident;
  DefineStmt* define_stmt = NULL; // set by ModuleUse
  CallExpr(string& qualified, string& ident) : qualified(move(qualified)), ident(move(ident)) {}
};

struct CollapseExpr : Visitable<Expr, CollapseExpr> {
  string qualified, ident;
  DefineStmt* define_stmt = NULL; // set by ModuleUse
  CollapseExpr(string& qualified, string& ident) : qualified(move(qualified)), ident(move(ident)) {}
};

struct ComplementExpr : Visitable<Expr, ComplementExpr> {
  Expr* inner;
  ComplementExpr(Expr* inner) : inner(inner) {}
  ~ComplementExpr() {
    delete inner;
  }
};

struct ConcatExpr : Visitable<Expr, ConcatExpr> {
  Expr *lhs, *rhs;
  ConcatExpr(Expr* lhs, Expr* rhs) : lhs(lhs), rhs(rhs) {}
  ~ConcatExpr() {
    delete lhs;
    delete rhs;
  }
};

struct DifferenceExpr : Visitable<Expr, DifferenceExpr> {
  Expr *lhs, *rhs;
  DifferenceExpr(Expr* lhs, Expr* rhs) : lhs(lhs), rhs(rhs) {}
  ~DifferenceExpr() {
    delete lhs;
    delete rhs;
  }
};

struct DotExpr : Visitable<Expr, DotExpr> {};

struct EmbedExpr : Visitable<Expr, EmbedExpr> {
  string qualified, ident;
  DefineStmt* define_stmt = NULL; // set by ModuleUse
  long macro_value; // set by ModuleUse
  EmbedExpr(string& qualified, string& ident) : qualified(move(qualified)), ident(move(ident)) {}
};

struct EpsilonExpr : Visitable<Expr, EpsilonExpr> {};

struct IntersectExpr : Visitable<Expr, IntersectExpr> {
  Expr *lhs, *rhs;
  IntersectExpr(Expr* lhs, Expr* rhs) : lhs(lhs), rhs(rhs) {}
  ~IntersectExpr() {
    delete lhs;
    delete rhs;
  }
};

struct LiteralExpr : Visitable<Expr, LiteralExpr> {
  string literal;
  LiteralExpr(string& literal) : literal(move(literal)) {}
};

struct PlusExpr : Visitable<Expr, PlusExpr> {
  Expr* inner;
  PlusExpr(Expr* inner) : inner(inner) {}
  ~PlusExpr() {
    delete inner;
  }
};

struct RepeatExpr : Visitable<Expr, RepeatExpr> {
  Expr* inner;
  long low, high;
  RepeatExpr(Expr* inner, long low, long high) : inner(inner), low(low), high(high) {}
  ~RepeatExpr() {
    delete inner;
  }
};

struct QuestionExpr : Visitable<Expr, QuestionExpr> {
  Expr* inner;
  QuestionExpr(Expr* inner) : inner(inner) {}
  ~QuestionExpr() {
    delete inner;
  }
};

struct StarExpr : Visitable<Expr, StarExpr> {
  Expr* inner;
  StarExpr(Expr* inner) : inner(inner) {}
  ~StarExpr() {
    delete inner;
  }
};

struct UnionExpr : Visitable<Expr, UnionExpr> {
  Expr *lhs, *rhs;
  UnionExpr(Expr* lhs, Expr* rhs) : lhs(lhs), rhs(rhs) {}
  ~UnionExpr() {
    delete lhs;
    delete rhs;
  }
};

//// Stmt

struct Stmt {
  Location loc;
  Stmt *prev = NULL, *next = NULL;
  virtual ~Stmt() = default;
  virtual void accept(Visitor<Stmt>& visitor) = 0;
};

struct EmptyStmt : Visitable<Stmt, EmptyStmt> {};

struct ActionStmt : Visitable<Stmt, ActionStmt> {
  string ident, code;
  ActionStmt(string& ident, string& code) : ident(move(ident)), code(move(code)) {}
};

struct CppStmt : Visitable<Stmt, CppStmt> {
  string code;
  CppStmt(string& code) : code(move(code)) {}
};

struct DefineStmt : Visitable<Stmt, DefineStmt> {
  bool export_ = false, intact = false;
  string export_params, lhs;
  Expr* rhs;
  Module* module; // used in topological sort
  DefineStmt(string& lhs, Expr* rhs) : lhs(move(lhs)), rhs(rhs) {}
  ~DefineStmt() {
    delete rhs;
  }
};

struct ImportStmt : Visitable<Stmt, ImportStmt> {
  string filename, qualified;
  ImportStmt(string& filename, string& qualified) : filename(move(filename)), qualified(move(qualified)) {}
};

struct PreprocessDefineStmt : Visitable<Stmt, PreprocessDefineStmt> {
  string ident;
  long value;
  PreprocessDefineStmt(string& ident, long value) : ident(move(ident)), value(value) {}
};

void stmt_free(Stmt* stmt);

//// Visitor implementations

struct StmtPrinter : Visitor<Action>, Visitor<Expr>, Visitor<Stmt> {
  int depth = 0;

  void visit(Action& action) override {
    action.accept(*this);
  }
  void visit(InlineAction& action) override {
    printf("%*s%s\n", 2*depth, "", "InlineAction");
    printf("%*s%s\n", 2*(depth+1), "", action.code.c_str());
  }
  void visit(RefAction& action) override {
    printf("%*s%s\n", 2*depth, "", "RefAction");
    printf("%*s%s\n", 2*(depth+1), "", action.ident.c_str());
  }

  void visit(Expr& expr) override {
    if (expr.entering.size()) {
      printf("%*s%s\n", 2*depth, "", "@entering");
      depth++;
      for (auto a: expr.entering) {
        indent(stdout, depth);
        printf("%ld\n", a.second);
        a.first->accept(*this);
      }
      depth--;
    }
    if (expr.finishing.size()) {
      printf("%*s%s\n", 2*depth, "", "@finishing");
      depth++;
      for (auto a: expr.finishing) {
        indent(stdout, depth);
        printf("%ld\n", a.second);
        a.first->accept(*this);
      }
      depth--;
    }
    if (expr.leaving.size()) {
      printf("%*s%s\n", 2*depth, "", "@entering");
      depth++;
      for (auto a: expr.leaving) {
        indent(stdout, depth);
        printf("%ld\n", a.second);
        a.first->accept(*this);
      }
      depth--;
    }
    if (expr.transiting.size()) {
      printf("%*s%s\n", 2*depth, "", "@transiting");
      depth++;
      for (auto a: expr.transiting) {
        indent(stdout, depth);
        printf("%ld\n", a.second);
        a.first->accept(*this);
      }
      depth--;
    }
    expr.accept(*this);
  }
  void visit(BracketExpr& expr) override {
    printf("%*s%s\n", 2*depth, "", "BracketExpr");
    printf("%*s", 2*(depth+1), "");
    for (auto& x: expr.intervals.to)
      printf("(%ld,%ld) ", x.first, x.second);
    puts("");
  }
  void visit(CallExpr& expr) override {
    printf("%*s%s\n", 2*depth, "", "CallExpr");
    printf("%*s", 2*(depth+1), "");
    if (expr.qualified.size())
      printf("%s::%s\n", expr.qualified.c_str(), expr.ident.c_str());
    else
      printf("%s\n", expr.ident.c_str());
  }
  void visit(CollapseExpr& expr) override {
    printf("%*s%s\n", 2*depth, "", "CollapseExpr");
    printf("%*s", 2*(depth+1), "");
    if (expr.qualified.size())
      printf("%s::%s\n", expr.qualified.c_str(), expr.ident.c_str());
    else
      printf("%s\n", expr.ident.c_str());
  }
  void visit(ComplementExpr& expr) override {
    printf("%*s%s\n", 2*depth, "", "ComplementExpr");
    depth++;
    visit(*expr.inner);
    depth--;
  }
  void visit(ConcatExpr& expr) override {
    printf("%*s%s\n", 2*depth, "", "ConcatExpr");
    depth++;
    visit(*expr.lhs);
    visit(*expr.rhs);
    depth--;
  }
  void visit(DifferenceExpr& expr) override {
    printf("%*s%s\n", 2*depth, "", "DifferenceExpr");
    depth++;
    visit(*expr.lhs);
    visit(*expr.rhs);
    depth--;
  }
  void visit(DotExpr& expr) override {
    printf("%*s%s\n", 2*depth, "", "DotExpr");
  }
  void visit(EmbedExpr& expr) override {
    printf("%*s%s\n", 2*depth, "", "EmbedExpr");
    printf("%*s", 2*(depth+1), "");
    if (expr.qualified.size())
      printf("%s::%s\n", expr.qualified.c_str(), expr.ident.c_str());
    else
      printf("%s\n", expr.ident.c_str());
  }
  void visit(EpsilonExpr& expr) override {
    printf("%*s%s\n", 2*depth, "", "EpsilonExpr");
  }
  void visit(IntersectExpr& expr) override {
    printf("%*s%s\n", 2*depth, "", "IntersectExpr");
    depth++;
    visit(*expr.lhs);
    visit(*expr.rhs);
    depth--;
  }
  void visit(LiteralExpr& expr) override {
    printf("%*s%s\n", 2*depth, "", "LiteralExpr");
    printf("%*s%s\n", 2*(depth+1), "", expr.literal.c_str());
  }
  void visit(PlusExpr& expr) override {
    printf("%*s%s\n", 2*depth, "", "PlusExpr");
    depth++;
    visit(*expr.inner);
    depth--;
  }
  void visit(RepeatExpr& expr) override {
    printf("%*s%s\n", 2*depth, "", "RepeatExpr");
    printf("%*s%ld,%ld\n", 2*(depth+1), "", expr.low, expr.high);
    depth++;
    visit(*expr.inner);
    depth--;
  }
  void visit(QuestionExpr& expr) override {
    printf("%*s%s\n", 2*depth, "", "QuestionExpr");
    depth++;
    visit(*expr.inner);
    depth--;
  }
  void visit(StarExpr& expr) override {
    printf("%*s%s\n", 2*depth, "", "StarExpr");
    depth++;
    visit(*expr.inner);
    depth--;
  }
  void visit(UnionExpr& expr) override {
    printf("%*s%s\n", 2*depth, "", "UnionExpr");
    depth++;
    visit(*expr.lhs);
    visit(*expr.rhs);
    depth--;
  }

  void visit(Stmt& stmt) override {
    stmt.accept(*this);
  }
  void visit(ActionStmt& stmt) override {
    printf("%*s%s\n", 2*depth, "", "ActionStmt");
    printf("%*s%s\n", 2*(depth+1), "", stmt.ident.c_str());
    printf("%*s%s\n", 2*(depth+1), "", stmt.code.c_str());
  }
  void visit(CppStmt& stmt) override {
    printf("%*s%s\n", 2*depth, "", "CppStmt");
    printf("%*s%s\n", 2*(depth+1), "", stmt.code.c_str());
  }
  void visit(DefineStmt& stmt) override {
    printf("%*s%s%s\n", 2*depth, "", "DefineStmt", stmt.export_ ? " export" : "");
    depth++;
    indent(stdout, depth);
    if (stmt.export_params.size())
      printf("(%s) ", stmt.export_params.c_str());
    printf("%s\n", stmt.lhs.c_str());
    visit(*stmt.rhs);
    depth--;
  }
  void visit(EmptyStmt& stmt) override {
    printf("%*s%s\n", 2*depth, "", "EmptyStmt");
  }
  void visit(ImportStmt& stmt) override {
    printf("%*s%s\n", 2*depth, "", "ImportStmt");
    printf("%*s%s\n", 2*(depth+1), "", stmt.filename.c_str());
    if (stmt.qualified.size())
      printf("%*sas %s\n", 2*(depth+1), "", stmt.qualified.c_str());
  }
  void visit(PreprocessDefineStmt& stmt) override {
    printf("%*s%s\n", 2*depth, "", "PreprocessDefineStmt");
    printf("%*s%s %ld\n", 2*(depth+1), "", stmt.ident.c_str(), stmt.value);
  }
};

//// Visitor implementation

struct PreorderStmtVisitor : Visitor<Stmt> {
  void visit(Stmt& stmt) override { stmt.accept(*this); }
  void visit(ActionStmt& stmt) override {}
  void visit(CppStmt& stmt) override {}
  void visit(DefineStmt& stmt) override {}
  void visit(EmptyStmt& stmt) override {}
  void visit(ImportStmt& stmt) override {}
  void visit(PreprocessDefineStmt&) override {}
};

struct PrePostActionExprStmtVisitor : Visitor<Action>, Visitor<Expr>, Visitor<Stmt> {
  virtual void pre_action(Action&) {}
  virtual void post_action(Action&) {}
  virtual void pre_expr(Expr&) {}
  virtual void post_expr(Expr&) {}
  virtual void pre_stmt(Stmt&) {}
  virtual void post_stmt(Stmt&) {}

  void visit(Action& action) override {
    pre_action(action);
    action.accept(*this);
    post_action(action);
  }
  void visit(InlineAction&) override {}
  void visit(RefAction&) override {}

  void visit(Expr& expr) override {
    pre_expr(expr);
    expr.accept(*this);
    post_expr(expr);
  }
  void visit(BracketExpr& expr) override {}
  void visit(CallExpr& expr) override {}
  void visit(CollapseExpr& expr) override {}
  void visit(ComplementExpr& expr) override { visit(*expr.inner); }
  void visit(ConcatExpr& expr) override {
    visit(*expr.lhs);
    visit(*expr.rhs);
  }
  void visit(DifferenceExpr& expr) override {
    visit(*expr.lhs);
    visit(*expr.rhs);
  }
  void visit(DotExpr& expr) override {}
  void visit(EmbedExpr& expr) override {}
  void visit(EpsilonExpr& expr) override {}
  void visit(IntersectExpr& expr) override {
    visit(*expr.lhs);
    visit(*expr.rhs);
  }
  void visit(LiteralExpr& expr) override {}
  void visit(PlusExpr& expr) override { visit(*expr.inner); }
  void visit(RepeatExpr& expr) override { visit(*expr.inner); }
  void visit(QuestionExpr& expr) override { visit(*expr.inner); }
  void visit(StarExpr& expr) override { visit(*expr.inner); }
  void visit(UnionExpr& expr) override {
    visit(*expr.lhs);
    visit(*expr.rhs);
  }

  void visit(Stmt& stmt) override {
    pre_stmt(stmt);
    stmt.accept(*this);
    post_stmt(stmt);
  }
  void visit(ActionStmt& stmt) override {}
  void visit(CppStmt& stmt) override {}
  void visit(DefineStmt& stmt) override { stmt.rhs->accept(*this); }
  void visit(EmptyStmt& stmt) override {}
  void visit(ImportStmt& stmt) override {}
  void visit(PreprocessDefineStmt&) override {}
};
