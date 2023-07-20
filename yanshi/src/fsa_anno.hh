#pragma once
#include "fsa.hh"
#include "syntax.hh"

enum class ExprTag {
  start = 1,
  inner = 2,
  final = 4,
};

extern inline bool has_start(ExprTag x) { return long(x) & long(ExprTag::start); }
extern inline bool has_inner(ExprTag x) { return long(x) & long(ExprTag::inner); }
extern inline bool has_final(ExprTag x) { return long(x) & long(ExprTag::final); }

bool operator<(ExprTag x, ExprTag y);
bool assoc_has_expr(vector<pair<Expr*, ExprTag>>& as, const Expr* x);

struct FsaAnno {
  bool deterministic;
  Fsa fsa;
  vector<vector<pair<Expr*, ExprTag>>> assoc;
  void accessible(const vector<long>* starts, vector<long>& mapping);
  void add_assoc(Expr& expr);
  void complement(ComplementExpr* expr);
  void co_accessible(const vector<bool>* final, vector<long>& mapping);
  void concat(FsaAnno& rhs, ConcatExpr* expr);
  void determinize(const vector<long>* starts, vector<vector<long>>* mapping);
  void difference(FsaAnno& rhs, DifferenceExpr* expr);
  void intersect(FsaAnno& rhs, IntersectExpr* expr);
  void minimize(vector<vector<long>>* mapping);
  void plus(PlusExpr* expr);
  void question(QuestionExpr* expr);
  void repeat(RepeatExpr& expr);
  void star(StarExpr* expr);
  void substring_grammar();
  void union_(FsaAnno& rhs, UnionExpr* expr);
  static FsaAnno bracket(BracketExpr& expr);
  static FsaAnno call(CallExpr& expr);
  static FsaAnno collapse(CollapseExpr& expr);
  static FsaAnno dot(DotExpr* expr);
  static FsaAnno embed(EmbedExpr& expr);
  static FsaAnno epsilon_fsa(EpsilonExpr* expr);
  static FsaAnno literal(LiteralExpr& expr);
};
