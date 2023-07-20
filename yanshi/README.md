# 偃师 (yanshi)

yanshi is a finite state automaton generator like Ragel. Use inline operators to embed C++ code in the recognition of a language. yanshi is enhanced with features to approximate context-free grammar:

- Approximation of substring grammar
- Approximation of recursive automaton to match expressions with recursion.

The motivation to create yanshi is that Ragel does not provide a mechanism to serialize its representation of finite state automata, making it difficult to post-process generated automata and obtain the substring grammar recognizer.

Later on, I found a simplified SQL grammar might contain more than 10000 states. It was not only slow to generate the automaton, making it hard to do trial-and-error experiments, but wasted memory to store the automaton. I introduced `CollapseExpr` to allow circular references.

`CallExpr` takes one step further, maintains a return address stack to imitate function calls. It can be seen as an augmented `CollapseExpr`, removing a lot of false positive cases.

## Name

From <https://en.wikipedia.org/wiki/Automaton>:

<blockquote>
In ancient China, a curious account of automata is found in the Lie Zi text (列子), written in the 3rd century BC. Within it there is a description of a much earlier encounter between King Mu of Zhou (周穆王, 1023-957 BC) and a mechanical engineer known as Yan Shi (偃师), an 'artificer'.
</blockquote>

## Build

* Debug: `make`
* Release: `make build=release`

## Getting Started

* Create file `a.ys`:
  ```
  export foo = 'hello'
  ```

  Run `yanshi -S a.ys -o /tmp/a.cc` to generate a C++ file from the yanshi source file `a.ys`.

  + `yanshi_foo_start`: the start state is 0. States are represented by natural numbers.
  + `yanshi_foo_is_final`: leave aside `ret_stack` and look at the last line, it checks whether `u` is one of the final states.
  + `yanshi_foo_transit`: leave aside `ret_stack`, `u` is the current state and `c` is the next input codepoint or label.

  With the `-S` option, yanshi will generate a standalone C++ file.
  ```
  % make -C /tmp a
  make: Entering directory '/tmp'
  g++     a.cc   -o a
  make: Leaving directory '/tmp'
  % /tmp/a hello
  0 h 1 e 2 l 3 l 4 o 5
  len: 5
  pref: 5
  state: 5
  final: true
  % /tmp/a
  hello<press C-d>0 h 1 e 2 l 3 l 4 o 5
  len: 5
  pref: 5
  state: 5
  final: true
  ```

  States are yellow and interleaved with transition labels. Final states are bold yellow.
  + `len`: the length of input codepoints or labels
  + `pref`: the length of the longest prefix that does not enter the dead state
  + `state`: the state entered after consuming the input
  + `final`: whether the state is one of final states

* Interactive mode

  The `-i` option enables interactive mode.
  ```
  % yanshi -i a.ys
  Testing foo
  foo :: DefineStmt
  .integer mode
  Commands available from the prompt:
    .automaton    dump automaton
    .assoc        dump associated AST Expr for each state
    .help         display this help
    .integer      input is a list of non-negative integers, macros(#define) or ''  quoted strings
    .macro        display defined macros
    .string       input is a string
    .stmt <ident> change target DefineStmt to <ident>
    .quit         exit interactive mode
  λ 104 101 108 108 111
  0 104 1 101 2 108 3 108 4 111 5
  export foo = 'hello':
  λ .string
  .string mode
  λ hello
  0 h 1 e 2 l 3 l 4 o 5
  export foo = 'hello':
  λ
  ```

* Regex-like syntax
  ```
  export hello = [gh] 'e' l{2} 'o'
  l = 'l'
  ```

  `[gh]` is a bracket expression and `l{2}` denotes to matches `l` at least twice. This grammar matches `hello`, `gello`, `helllo`, ...

  + Union: `c = a | b`
  + Intersection: `c = a && b`
  + Difference: `c = a - b`
  + Concatenation: `c = a b`
  + Complement: `c = ~ a`

* Actions (embedded C++ code)
  ```
  c++ {
  #include <stdio.h>
  }
  export hello = '喵' @ { puts("meow"); } {2}
  ```
  I have not thought clearly on the implementation. The executing point may be counter-intuitive.

* Modules
  ```
  # a.ys
  import 'b.ys' as B # B::bar
  import 'b.ys' # qux

  export foo = B::bar | qux
  bar = '4'

  # b.ys
  bar = '3'
  qux = '5'
  ```

* Substring grammar
  Specify the `--substring-grammar` option to generate code for substring grammar. That is, the generated code matches every substring of the grammar. The implementation creates a new start state and a new final state, connects the start state to the old start state, and old final states to the new final state.

* `EmbedExpr`, reference a nonterminal without modifiers
   ```
   foo = bar
   bar = [0-9]
   ```
   The complete automaton of bar will be duplicated in each reference site. If the referenced automaton is large, `EmbedExpr` will significantly increase the number of states. `EmbedExpr` defines dependencies among states and no cyclic dependency is allowed.

* `CollapseExpr`, reference a nonterminal with the `!` modifier
  ```
  export foo = 'pre' !bar 'post'
  bar = [\u0300-\u034E]
  quz = 'meow' !bar 'meow'
  ```
  The final state of `'pre'` and the start state of `'post'` will be connected by a special directed arc. When exporting, an epsilon transition will be added from the tail of the arc to the start state of `bar`, others will be added from the final states of `bar` to the head of the arc. `CollapseExpr` behaves like function calls, however, the return address is not preserved (hence the name `CollapseExpr`) and the state may go to other call sites. In this example, the state may go to either `foo` or `quz` after traveling through `bar`, causing false positives.

* `CallExpr`, reference a nonterminal with the `&` modifier
  ```
  export foo = 'pre' &bar 'post'
  bar = '4'
  ```
  This is a refinement of `CollapseExpr`. Suppose state `&B` is contained in `A`'s definition (`A` calls `B`). `&B` will be represented as an pseudo arc (`u -> v`), where `u` is the state before `&B` and `v` is the state after `&B`. If arcs of `u` do not collide with arcs of `B`, the transition function will push `v` to the return stack if current state set contains `u` and there is no other transition. Note `B` is disconnected from `A`, which is different from the `CollapseExpr` case. The machine will transit on automaton `B` greedily. If there is no transition, it will pop a return address(`v` in this case) and jumps to it.

## Contrib

### Vim

Syntax highlighting, and a syntax checking plugin for Synaptics

```zsh
ln -sr contrib/vim/compiler/yanshi.vim ~/.vim/compiler/
ln -sr contrib/vim/ftdetect/yanshi.vim ~/.vim/ftdetect/
ln -sr contrib/vim/ftplugin/yanshi.vim ~/.vim/ftplugin/
ln -sr contrib/vim/syntax/yanshi.vim ~/.vim/syntax/
ln -sr contrib/vim/syntax_checker/yanshi ~/.vim/syntax_checker/
```

### Zsh

Command line completion

```
# ~/.zshrc
fpath=(~/.zsh $fpath)

# ln -sr contrib/zsh/_yanshi ~/.zsh/
```

## Internals

```
src
  common.{cc,hh}
  main.{cc,hh}
  syntax.{cc,hh}
  loader.{cc,hh}
  fsa.{cc,hh}
  fsa_anno.{cc,hh}
  compiler.{cc,hh}
  parser.y
  lexer.l
  location.cc
```

* Lex `lexer.l`
* Parse and generate a syntax tree `parser.y`
* `loader.cc`
  + Get a list of definitions
  + Recursively load for each `import`
  + Resolve references and associate uses to definitions
  + Build a dependency graph from `EmbedExpr`
  + Compile automaton for each nonterminal in topological order. `CollapseExpr` and `CallExpr` are represented by special directed arcs.
  + Generate code for `export` nonterminals, resolving `CollapseExpr` and `CallExpr`

### Finite state automaton

Each node of the syntax tree corresponds to an automaton. The parent builds an automaton from its children according to the semantics. The automaton of the parent may contain states from the automaton of one of the children, or it is a state introduced by the parent.

`assoc[i]` records the associative nodes in the automaton tree (which part of the syntax tree has associations with this state) and positions (start state, final state or inner state) for state `i`. It serves three purposes:

* Check which action should be triggered
* Look for inner states (neither start nor final) in the implementation of substring grammar
* Check whether it is associated to a `CallExpr` or `CollapseExpr`

### `CollapseExpr`
