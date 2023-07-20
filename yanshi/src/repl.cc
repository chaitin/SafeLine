#include "compiler.hh"
#include "fsa_anno.hh"
#include "loader.hh"
#include "parser.hh"
#include "lexer.hh" // after parser.hh
#include "syntax.hh"

#include <algorithm>
#include <assert.h>
#include <functional>
#include <inttypes.h>
#include <sstream>
#include <stdlib.h>
#include <type_traits>
#include <unicode/utf8.h>
#include <unordered_map>
#include <wctype.h>
#ifdef HAVE_READLINE
# include <readline/readline.h>
# include <readline/history.h>
#endif
using namespace std;

enum class ReplMode {string, integer};
static ReplMode mode = ReplMode::string;
static const FsaAnno* anno;
static bool quit;

struct Command
{
  const char* name;
  function<void(const char*)> fn;
} commands[] = {
  {".automaton", [](const char*) {print_automaton(anno->fsa); }},
  {".assoc", [](const char*) {print_assoc(*anno); }},
  {".help",
    [](const char*) {
      fputs("Commands available from the prompt:\n"
             "  .automaton    dump automaton\n"
             "  .assoc        dump associated AST Expr for each state\n"
             "  .help         display this help\n"
             "  .integer      input is a list of non-negative integers, macros(#define) or '' "" quoted strings\n"
             "  .macro        display defined macros\n"
             "  .string       input is a string\n"
             "  .stmt <ident> change target DefineStmt to <ident>\n"
             "  .quit         exit interactive mode\n"
             , stdout);
    }},
  {".integer",
    [](const char*) {
      mode = ReplMode::integer; puts(".integer mode");
    }},
  {".macro",
    [](const char*) {
      for (auto& it: main_module->macro)
        printf("%s\t%ld\n", it.first.c_str(), it.second->value);
      for (auto* import: main_module->unqualified_import)
        for (auto& it: import->macro)
          printf("%s\t%ld\n", it.first.c_str(), it.second->value);
    }},
  {".stmt",
    [](const char* arg) {
      Stmt* r = resolve(*main_module, "", arg);
      if (! r)
        printf("'%s' undefined\n", arg);
      else if (r == (Stmt*)1)
        printf("ambiguous '%s'\n", arg);
      else if (auto d = dynamic_cast<PreprocessDefineStmt*>(r))
        printf("'%s' is a macro\n", arg);
      else if (auto d = dynamic_cast<DefineStmt*>(r)) {
        anno = &compiled[d];
        printf("%s :: DefineStmt\n", d->lhs.c_str());
      } else
        assert(0);
    }},
  {".string",
    [](const char*) {
      mode = ReplMode::string; puts(".string mode");
    }
  },
  {".quit",
    [](const char*) {
      puts("Leaving interactive mode");
      quit = true;
    }
  },
};

#ifdef HAVE_READLINE
static char* command_completer(const char* text, int state)
{
  static long i = 0;
  if (! state)
    i = 0;
  while (i < LEN(commands)) {
    Command* x = &commands[i++];
    if (! strncmp(x->name, text, strlen(text)))
      return strdup(x->name);
  }
  return NULL;
}

static char* macro_completer(const char* text, int state)
{
  static Stmt* x;
  if (! state)
    x = main_module->toplevel;
  while (x) {
    auto xx = dynamic_cast<PreprocessDefineStmt*>(x);
    x = x->next;
    if (xx && ! strncmp(xx->ident.c_str(), text, strlen(text)))
      return strdup(xx->ident.c_str());
  }
  return NULL;
}

static char* stmt_completer(const char* text, int state)
{
  static Stmt* x;
  if (! state)
    x = main_module->toplevel;
  while (x) {
    auto xx = dynamic_cast<DefineStmt*>(x);
    x = x->next;
    if (xx && ! strncmp(xx->lhs.c_str(), text, strlen(text)))
      return strdup(xx->lhs.c_str());
  }
  return NULL;
}

static char** on_complete(const char* text, int start, int end)
{
  rl_attempted_completion_over = 1;
  if (! start)
    return rl_completion_matches(text, command_completer);
  if (6 <= start && ! strncmp(rl_line_buffer, ".stmt ", 6))
    return rl_completion_matches(text, stmt_completer);
  if (mode == ReplMode::integer)
    return rl_completion_matches(text, macro_completer);
  return NULL;
}
#else
char* readline(const char* prompt)
{
  char* r = NULL;
  size_t s = 0;
  ssize_t n;
  fputs(prompt, stdout);
  if ((n = getline(&r, &s, stdin)) > 0)
    r[n-1] = '\0';
  else {
    free(r);
    r = NULL;
  }
  return r;
}
#endif

static void run_command(char* line)
{
  size_t p = 1;
  while (line[p] && isalnum(line[p]))
    p++;
  Command* com = NULL;
  REP(i, LEN(commands))
    if (! strncmp(commands[i].name, line, p)) {
      if (com) com = (Command*)1;
      else com = &commands[i];
    }
  if (! com)
    printf("Unknown command '%s'\n", line);
  else if (com == (Command*)1)
    printf("Ambiguous command '%s'\n", line);
  else {
    while (line[p] && isspace(line[p]))
      p++;
    size_t len = strlen(line);
    while (len && isspace(line[len-1]))
      line[--len] = '\0';
    com->fn(line+p);
  }
}

void repl(DefineStmt* stmt)
{
#ifdef HAVE_READLINE
  rl_attempted_completion_function = on_complete;
#endif
  char buf[BUF_SIZE];
  snprintf(buf, sizeof buf, ".stmt %s", stmt->lhs.c_str());
  run_command(buf);
  strcpy(buf, ".integer");
  run_command(buf);
  strcpy(buf, ".help");
  run_command(buf);
  if (! anno) return;
  char* line;
  stringstream ss;
  while ((line = readline("λ ")) != NULL) {
#ifdef HAVE_READLINE
    if (line[0])
      add_history(line);
#endif
    if (line[0] == '.') {
      run_command(line);
      free(line);
      if (quit) break;
      continue;
    }

    long u = anno->fsa.start;

    if (mode == ReplMode::string) {
      i32 i = 0, len;
      long c;
      len = strlen(line);
      if (anno->fsa.is_final(u)) yellow(1);
      else normal_yellow(1);
      printf("%ld ", u); sgr0();
      while (i < len) {
        U8_NEXT_OR_FFFD(line, i, len, c);
        if (iswcntrl(c)) printf("%ld ", c);
        else printf("%lc ", wint_t(c));
        u = anno->fsa.transit(u, c);
        if (anno->fsa.is_final(u)) yellow();
        else normal_yellow();
        printf("%ld ", u); sgr0();
        if (u < 0) break;
      }
    } else {
      vector<long> input;
      int token;
      yyscan_t lexer;
      raw_yylex_init_extra(0, &lexer);
      YY_BUFFER_STATE buf = raw_yy_scan_bytes(line, strlen(line), lexer);
      YYSTYPE yylval;
      YYLTYPE yylloc;
      while (u >= 0 && (token = raw_yylex(&yylval, &yylloc, lexer)) != 0) {
        switch (token) { // all tokens with a destructor should be listed
        case IDENT: {
          Stmt* stmt = resolve(*main_module, "", yylval.str->c_str());
          if (! stmt) {
            printf("'%s' undefined", yylval.str->c_str());
            u = -1;
          } else if (stmt == (Stmt*)1) {
            printf("ambiguous '%s'", yylval.str->c_str());
            u = -1;
          } else if (auto d = dynamic_cast<PreprocessDefineStmt*>(stmt))
            input.push_back(d->value);
          else if (auto d = dynamic_cast<DefineStmt*>(stmt)) {
            printf("'%s' is not a macro", yylval.str->c_str());
            u = -1;
          } else
            assert(0);
          delete yylval.str;
          break;
        }
        case INTEGER:
          input.push_back(yylval.integer);
          break;
        case STRING_LITERAL:
          if (opt_bytes)
            for (unsigned char c: *yylval.str)
              input.push_back(c);
          else
            for (i32 c, i = 0; i < yylval.str->size(); ) {
              U8_NEXT_OR_FFFD(yylval.str->c_str(), i, yylval.str->size(), c);
              input.push_back(c);
            }
          delete yylval.str;
          break;
        case BRACED_CODE:
          delete yylval.str;
          // fall through
        default:
          printf("invalid token at column %ld-%ld\n", yylloc.start+1, yylloc.end);
          u = -1;
          break;
        }
      }
      raw_yy_delete_buffer(buf, lexer);
      raw_yylex_destroy(lexer);

      if (u >= 0) {
        if (anno->fsa.is_final(u)) yellow(1);
        else normal_yellow(1);
        printf("%ld ", u); sgr0();
        for (long c: input) {
          printf("%ld ", c);
          u = anno->fsa.transit(u, c);
          if (anno->fsa.is_final(u)) yellow();
          else normal_yellow();
          printf("%ld ", u); sgr0();
          if (u < 0) break;
        }
      }
    }
    free(line);
    puts("");
    if (u >= 0) {
      unordered_map<DefineStmt*, vector<long>> start_finals;
      unordered_map<DefineStmt*, vector<pair<long, long>>> inners;
      for (auto aa: anno->assoc[u]) {
        if (has_start(aa.second))
          start_finals[aa.first->stmt].push_back(aa.first->loc.start);
        if (has_inner(aa.second)) {
          inners[aa.first->stmt].emplace_back(aa.first->loc.start, 1);
          inners[aa.first->stmt].emplace_back(aa.first->loc.end, -1);
        }
        if (has_final(aa.second))
          start_finals[aa.first->stmt].push_back(aa.first->loc.end);
      }
      vector<DefineStmt*> stmts;
      for (auto& it: start_finals)
        stmts.push_back(it.first);
      for (auto& it: inners)
        stmts.push_back(it.first);
      // sort DefineStmt by location
      sort(ALL(stmts), [](const DefineStmt* x, const DefineStmt* y) {
        if (x->module != y->module)
          return x->module < y->module;
        if (x->loc.start != y->loc.start)
          return x->loc.start < y->loc.start;
        return x->loc.end < y->loc.end;
      });
      stmts.erase(unique(ALL(stmts)), stmts.end());
      for (auto* stmt: stmts) {
        auto& start_final = start_finals[stmt];
        auto& inner = inners[stmt];
        sort(ALL(start_final));
        sort(ALL(inner));
        auto it0 = start_final.begin();
        auto it1 = inner.begin();
        long nest = 0;
        FOR(i, stmt->loc.start, stmt->loc.end) {
          for (; it0 != start_final.end() && *it0 < i; ++it0);
          if (it0 != start_final.end() && *it0 == i) {
            cyan();
            putchar(':');
            sgr0();
          }
          for (; it1 != inner.end() && it1->first <= i; ++it1)
            nest += it1->second;
          if (nest)
            yellow();
          putchar(stmt->module->locfile.data[i]);
          if (nest)
            sgr0();
        }
      }
    }
  }
}
