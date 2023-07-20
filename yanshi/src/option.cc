#include "common.hh"
#include "option.hh"
#include <stdio.h>

bool opt_bytes, opt_check, opt_dump_action, opt_dump_assoc, opt_dump_automaton, opt_dump_embed, opt_dump_module, opt_dump_tree, opt_gen_c, opt_gen_extern_c, opt_keep_inaccessible, opt_standalone, opt_substring_grammar;

long AB = MAX_CODEPOINT+1, opt_max_return_stack = 100;
long debug_level = 3;
FILE* debug_file;
const char* opt_output_filename = "-";
const char* opt_output_header_filename;
Mode opt_mode = Mode::cxx;
vector<string> opt_include_paths;
