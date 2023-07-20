#pragma once
#include "fsa_anno.hh"
#include "syntax.hh"

#include <unordered_map>
using std::unordered_map;

void print_assoc(const FsaAnno& anno);
void print_automaton(const Fsa& fsa);
void compile(DefineStmt*);
bool compile_export(DefineStmt* stmt);
void generate_cxx(Module* mo);
void generate_graphviz(Module* mo);
extern unordered_map<DefineStmt*, FsaAnno> compiled;
