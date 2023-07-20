#include "common.hh"
#include "fsa.hh"
#include "loader.hh"
#include "option.hh"

#include <errno.h>
#include <getopt.h>
#include <locale.h>
#include <stdarg.h>
#include <stdio.h>
#include <string.h>
#include <string>
#include <sysexits.h>
#include <unistd.h>
using namespace std;

void print_help(FILE *fh)
{
  fprintf(fh, "Usage: %s [OPTIONS] dir\n", program_invocation_short_name);
  fputs(
        "\n"
        "Options:\n"
        "  -b,--bytes                make labels range over [0,256), Unicode literals will be treated as UTF-8 bytes\n"
        "  -C                        generate C source code (default: C++)\n"
        "  --check                   check syntax & use/def\n"
        "  --debug                   debug level\n"
        "  --debug-output            filename for debug output\n"
        "  --dump-action             dump associated actions for each edge\n"
        "  --dump-assoc              dump associated AST Expr for each state\n"
        "  --dump-automaton          dump automata\n"
        "  --dump-embed              dump statistics of EmbedExpr\n"
        "  --dump-module             dump module use/def/...\n"
        "  --dump-tree               dump AST\n"
        "  --extern-c                generate extern \"C\" specifier\n"
        "  -G,--graph <dir>          output a Graphviz dot file\n"
        "  -I,--import <dir>         add <dir> to search path for 'import'\n"
        "  -i,--interactive          interactive mode\n"
        "  --max-return-stack        max length of return stack in C generator (default: 100)\n"
        "  -k,--keep-inaccessible    do not perform accessible/co-accessible\n"
        "  -S,--standalone           generate header and 'main()'\n"
        "  --substring-grammar       construct regular approximation of the substring grammar. Inner states of nonterminals labeled 'intact' are not connected to start/final\n"
        "  -o,--output <file>        .cc output filename\n"
        "  -O,--output-header <file> .hh output filename\n"
        "  -h, --help                display this help and exit\n"
        "\n"
        , fh);
  exit(fh == stdout ? 0 : EX_USAGE);
}

int main(int argc, char *argv[])
{
  setlocale(LC_ALL, "");
  int opt;
  static struct option long_options[] = {
    {"bytes",               no_argument,       0,   'b'},
    {"check",               required_argument, 0,   'c'},
    {"debug",               required_argument, 0,   'd'},
    {"debug-output",        required_argument, 0,   'l'},
    {"dump-action",         no_argument,       0,   1000},
    {"dump-assoc",          no_argument,       0,   1001},
    {"dump-automaton",      no_argument,       0,   1002},
    {"dump-embed",          no_argument,       0,   1003},
    {"dump-module",         no_argument,       0,   1004},
    {"dump-tree",           no_argument,       0,   1005},
    {"extern-c",            no_argument,       0,   1007},
    {"graph",               no_argument,       0,   'G'},
    {"import",              required_argument, 0,   'I'},
    {"interactive",         no_argument,       0,   'i'},
    {"max-return-stack",    required_argument, 0,   1006},
    {"keep-inaccessible",   no_argument,       0,   'k'},
    {"standalone",          no_argument,       0,   'S'},
    {"substring-grammar",   no_argument,       0,   's'},
    {"output",              required_argument, 0,   'o'},
    {"output-header",       required_argument, 0,   'O'},
    {"help",                no_argument,       0,   'h'},
    {0,                     0,                 0,   0},
  };

  while ((opt = getopt_long(argc, argv, "bCDcd:GhI:ikl:O:o:Ss", long_options, NULL)) != -1) {
    switch (opt) {
    case 'b':
      opt_bytes = true;
      AB = 256;
      break;
    case 'C':
      opt_gen_c = true;
      break;
    case 'D':
      break;
    case 'c':
      opt_check = true;
      break;
    case 'd':
      debug_level = get_long(optarg);
      break;
    case 'G':
      opt_mode = Mode::graphviz;
      break;
    case 'h':
      print_help(stdout);
      break;
    case 'I':
      opt_include_paths.push_back(string(optarg));
      break;
    case 'i':
      opt_mode = Mode::interactive;
      break;
    case 'k':
      opt_keep_inaccessible = true;
      break;
    case 'l':
      if (debug_file)
        err_exit(EX_USAGE, "multiple '-l'");
      debug_file = fopen(optarg, "w");
      if (! debug_file)
        err_exit(EX_OSFILE, "fopen");
      break;
    case 'O':
      opt_output_header_filename = optarg;
      break;
    case 'o':
      opt_output_filename = optarg;
      break;
    case 'S':
      opt_standalone = true;
      break;
    case 's':
      opt_substring_grammar = true;
      break;
    case 1000: opt_dump_action = true; break;
    case 1001: opt_dump_assoc = true; break;
    case 1002: opt_dump_automaton = true; break;
    case 1003: opt_dump_embed = true; break;
    case 1004: opt_dump_module = true; break;
    case 1005: opt_dump_tree = true; break;
    case 1006:
      opt_max_return_stack = get_long(optarg);
      break;
    case 1007: opt_gen_extern_c = true; break;
    case '?':
      print_help(stderr);
      break;
    }
  }
  if (! debug_file)
    debug_file = stderr;
  argc -= optind;
  argv += optind;

  long n_errors = load(argc ? argv[0] : "-");
  unload_all();
  fclose(debug_file);
  return n_errors ? 2 : 0;
}
