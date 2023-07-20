#pragma once
#include <string>
#include <vector>
using std::string;
using std::vector;

extern bool opt_bytes, opt_check, opt_dump_action, opt_dump_assoc, opt_dump_automaton, opt_dump_embed, opt_dump_module, opt_dump_tree, opt_gen_c, opt_gen_extern_c, opt_keep_inaccessible, opt_standalone, opt_substring_grammar;
extern long AB, opt_max_return_stack;
extern const char* opt_output_filename;
extern const char* opt_output_header_filename;
enum class Mode {cxx, graphviz, interactive};
extern Mode opt_mode;
extern vector<string> opt_include_paths;
