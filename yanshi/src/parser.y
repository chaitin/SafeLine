%code requires {
#include "common.hh"
#include "location.hh"
#include "option.hh"
#include "syntax.hh"

#include <limits.h>
#include <unicode/utf8.h>

#define YYINITDEPTH 1000
#define YYLTYPE Location
#define YYLLOC_DEFAULT(Loc, Rhs, N)             \
  do {                                          \
    if (N) {                                    \
      (Loc).start = YYRHSLOC(Rhs, 1).start;     \
      (Loc).end = YYRHSLOC(Rhs, N).end;         \
    } else {                                    \
      (Loc).start = YYRHSLOC(Rhs, 0).start;     \
      (Loc).end = YYRHSLOC(Rhs, 0).end;         \
    }                                           \
  } while (0)

int parse(const LocationFile& locfile, Stmt*& res);
}

%locations
%error-verbose
%define api.pure

%parse-param {Stmt*& res}
%parse-param {long& errors}
%parse-param {const LocationFile& locfile}
%parse-param {void** lexer}
%lex-param {Stmt*& res}
%lex-param {long& errors}
%lex-param {const LocationFile& locfile}
%lex-param {void** lexer}

%union {
  long integer;
  string* str;
  DisjointIntervals* intervals;
  Action* action;
  Expr* expr;
  Stmt* stmt;
  char* errmsg;
}
%destructor { delete $$; } <str>
%destructor { delete $$; } <action>
%destructor { delete $$; } <expr>
%destructor { delete $$; } <intervals>
%destructor { delete $$; } <stmt>

%token ACTION AMPERAMPER AS COLONCOLON CPP DOTDOT EPSILON EXPORT IMPORT INTACT INVALID_CHARACTER PREPROCESS_DEFINE
%token <integer> CHAR INTEGER
%token <str> IDENT
%token <str> BRACED_CODE
%token <str> STRING_LITERAL

%type <action> action
%type <expr> concat_expr difference_expr factor repeat intersect_expr union_expr union_expr2 unop_expr
%type <intervals> bracket bracket_items
%type <stmt> define_stmt preprocess stmt stmt_list

%{
#include "lexer.hh"

#define FAIL(loc, errmsg)                                          \
  do {                                                             \
    Location l = loc;                                              \
    yyerror(&l, res, errors, locfile, lexer, errmsg);              \
  } while (0)

void yyerror(YYLTYPE* loc, Stmt*& res, long& errors, const LocationFile& locfile, yyscan_t* lexer, const char *errmsg)
{
  errors++;
  locfile.error_context(*loc, "%s", errmsg);
}

int yylex(YYSTYPE* yylval, YYLTYPE* loc, Stmt*& res, long& errors, const LocationFile& locfile, yyscan_t* lexer)
{
  int token = raw_yylex(yylval, loc, *lexer);
  if (token == INVALID_CHARACTER) {
    FAIL(*loc, yylval->errmsg ? yylval->errmsg : "invalid character");
    free(yylval->errmsg);
  }
  return token;
}

#define gen_repeat(x, inner, low, high) \
  if (low < 0) {                     \
    FAIL(yyloc, "negative"); \
  } \
  if (low > high) { \
    FAIL(yyloc, "low > high"); \
  } \
  x = new RepeatExpr(inner, low, high)
%}

%%

toplevel:
  stmt_list { res = $1; }

stmt_list:
    %empty { $$ = new EmptyStmt; }
  | '\n' stmt_list { $$ = $2; }
  | stmt stmt_list { $1->next = $2; $2->prev = $1; $$ = $1; }
  | error stmt_list { $$ = $2; }

stmt:
    define_stmt { $$ = $1; }
  | preprocess '\n' { $$ = $1; }
  | IMPORT STRING_LITERAL AS IDENT '\n' { $$ = new ImportStmt(*$2, *$4); delete $2; delete $4; $$->loc = yyloc; }
  | IMPORT STRING_LITERAL '\n' { string t; $$ = new ImportStmt(*$2, t); delete $2; $$->loc = yyloc; }
  | ACTION IDENT BRACED_CODE '\n' { $$ = new ActionStmt(*$2, *$3); delete $2; delete $3; $$->loc = yyloc; }
  | CPP BRACED_CODE '\n' { $$ = new CppStmt(*$2); delete $2; $$->loc = yyloc; }

preprocess:
    PREPROCESS_DEFINE IDENT INTEGER { $$ = new PreprocessDefineStmt(*$2, $3); delete $2; $$->loc = yyloc; }

eq:
    '='
  | ':'

define_stmt:
    IDENT eq union_expr '\n' { $$ = new DefineStmt(*$1, $3); delete $1; $$->loc = yyloc; }
  | IDENT eq '|' union_expr '\n' { $$ = new DefineStmt(*$1, $4); delete $1; $$->loc = yyloc; }
  | IDENT eq '\n' union_expr2 '\n' { $$ = new DefineStmt(*$1, $4); delete $1; $$->loc = yyloc; }
  | IDENT eq '\n' '|' union_expr2 '\n' { $$ = new DefineStmt(*$1, $5); delete $1; $$->loc = yyloc; }
  | EXPORT define_stmt { $$ = $2; ((DefineStmt*)$$)->export_ = true; $$->loc = yyloc; }
  | EXPORT BRACED_CODE define_stmt { $$ = $3; ((DefineStmt*)$$)->export_ = true; ((DefineStmt*)$$)->export_params = *$2; delete $2; $$->loc = yyloc; }
  | INTACT define_stmt { $$ = $2; ((DefineStmt*)$$)->intact = true; $$->loc = yyloc; }

union_expr:
    intersect_expr { $$ = $1; }
  | union_expr '|' intersect_expr { $$ = new UnionExpr($1, $3); $$->loc = yyloc; }

union_expr2:
    intersect_expr { $$ = $1; }
  | union_expr2 '|' intersect_expr { $$ = new UnionExpr($1, $3); $$->loc = yyloc; }
  | union_expr2 '\n' '|' intersect_expr { $$ = new UnionExpr($1, $4); $$->loc = yyloc; }

intersect_expr:
    difference_expr { $$ = $1; }
  | intersect_expr AMPERAMPER difference_expr { $$ = new IntersectExpr($1, $3); $$->loc = yyloc; }

difference_expr:
    concat_expr { $$ = $1; }
  | difference_expr '-' concat_expr { $$ = new DifferenceExpr($1, $3); $$->loc = yyloc; }

concat_expr:
    unop_expr { $$ = $1; }
  | concat_expr unop_expr { $$ = new ConcatExpr($1, $2); $$->loc = yyloc; }

unop_expr:
    factor { $$ = $1; }
  | '~' unop_expr { $$ = new ComplementExpr($2); $$->loc = yyloc; }

factor:
    EPSILON { $$ = new EpsilonExpr; $$->loc = yyloc; }
  | IDENT { string t; $$ = new EmbedExpr(t, *$1); delete $1; $$->loc = yyloc; }
  | IDENT COLONCOLON IDENT { $$ = new EmbedExpr(*$1, *$3); delete $1; delete $3; $$->loc = yyloc; }
  | '!' IDENT { string t; $$ = new CollapseExpr(t, *$2); delete $2; $$->loc = yyloc; }
  | '!' IDENT COLONCOLON IDENT { $$ = new CollapseExpr(*$2, *$4); delete $2; delete $4; $$->loc = yyloc; }
  | '&' IDENT { string t; $$ = new CallExpr(t, *$2); delete $2; $$->loc = yyloc; }
  | '&' IDENT COLONCOLON IDENT { $$ = new CallExpr(*$2, *$4); delete $2; delete $4; $$->loc = yyloc; }
  | STRING_LITERAL { $$ = new LiteralExpr(*$1); delete $1; $$->loc = yyloc; }
  | '.' { $$ = new DotExpr(); $$->loc = yyloc; }
  | INTEGER {
      if (opt_bytes && 256 <= $1) {
        FAIL(yyloc, "literal integers should be less than 256 in bytes mode");
        $$ = new DotExpr;
      } else {
        auto t = new DisjointIntervals;
        t->emplace($1, $1+1);
        $$ = new BracketExpr(t);
        $$->loc = yyloc;
      }
    }
  | bracket { $$ = new BracketExpr($1); $$->loc = yyloc; }
  | STRING_LITERAL DOTDOT STRING_LITERAL {
      i32 c0, c1, i = 0, j = 0;
      if (opt_bytes) {
        c0 = u8((*$1)[0]);
        c1 = u8((*$3)[0]);
        i = j = 1;
      } else {
        U8_NEXT($1->c_str(), i, $1->size(), c0);
        U8_NEXT($3->c_str(), j, $3->size(), c1);
      }
      delete $1;
      delete $3;
      if (i != $1->size() || j != $3->size()) {
        FAIL(yyloc, "endpoints of Unicode range should be of length 1");
        $$ = new DotExpr;
      } else if (c0 > c1) {
        FAIL(yyloc, "negative Unicode range");
        $$ = new DotExpr;
      } else {
        auto t = new DisjointIntervals;
        t->emplace(c0, c1+1);
        $$ = new BracketExpr(t);
        $$->loc = yyloc;
      }
    }
  | '(' union_expr ')' { $$ = $2; }
  | '(' error ')' { $$ = new DotExpr; }
  | repeat { $$ = $1; }
  | factor '>' action { $$ = $1; $$->entering.emplace_back($3, 0L); }
  | factor '>' INTEGER action { $$ = $1; $$->entering.emplace_back($4, $3); }
  | factor '@' action { $$ = $1; $$->finishing.emplace_back($3, 0L); }
  | factor '@' INTEGER action { $$ = $1; $$->finishing.emplace_back($4, $3); }
  | factor '%' action { $$ = $1; $$->leaving.emplace_back($3, 0L); }
  | factor '%' INTEGER action { $$ = $1; $$->leaving.emplace_back($4, $3); }
  | factor '$' action { $$ = $1; $$->transiting.emplace_back($3, 0L); }
  | factor '$' INTEGER action { $$ = $1; $$->transiting.emplace_back($4, $3); }
  | factor '+' { $$ = new PlusExpr($1); $$->loc = yyloc; }
  | factor '?' { $$ = new QuestionExpr($1); $$->loc = yyloc; }
  | factor '*' { $$ = new StarExpr($1); $$->loc = yyloc; }

repeat:
    factor '{' INTEGER ',' INTEGER '}' { gen_repeat($$, $1, $3, $5); $$->loc = yyloc; }
  | factor '{' INTEGER ',' '}' { gen_repeat($$, $1, $3, LONG_MAX); $$->loc = yyloc; }
  | factor '{' INTEGER '}' { gen_repeat($$, $1, $3, $3); $$->loc = yyloc; }
  | factor '{' ',' INTEGER '}' { gen_repeat($$, $1, 0, $4); $$->loc = yyloc; }

action:
    IDENT { string t; $$ = new RefAction(t, *$1); delete $1; $$->loc = yyloc; }
  | IDENT COLONCOLON IDENT { $$ = new RefAction(*$1, *$3); delete $1; delete $3; $$->loc = yyloc; }
  | BRACED_CODE { $$ = new InlineAction(*$1); delete $1; $$->loc = yyloc; }

bracket:
    '[' bracket_items ']' { $$ = $2; }
  | '[' '^' bracket_items ']' {
      $$ = $3;
      $$->flip();
    }

bracket_items:
    bracket_items CHAR '-' CHAR {
      $$ = $1;
      if ($2 > $4)
        FAIL(yyloc, "negative range in character class");
      else
        $$->emplace($2, $4+1);
    }
  | bracket_items CHAR {
      $$ = $1;
      $$->emplace($2, $2+1);
    }
  | %empty { $$ = new DisjointIntervals; }

%%

int parse(const LocationFile& locfile, Stmt*& res)
{
  yyscan_t lexer;
  raw_yylex_init_extra(0, &lexer);
  YY_BUFFER_STATE buf = raw_yy_scan_bytes(locfile.data.c_str(), locfile.data.size(), lexer);
  long errors = 0;
  yyparse(res, errors, locfile, &lexer);
  raw_yy_delete_buffer(buf, lexer);
  raw_yylex_destroy(lexer);
  if (errors > 0) {
    stmt_free(res);
    res = NULL;
  }
  return errors;
}
