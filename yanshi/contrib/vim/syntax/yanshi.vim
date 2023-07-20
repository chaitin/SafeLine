if exists('b:current_syntax')
  finish
endif

syn cluster yanshiCommentGroup contains=yanshiTodo
syn include @yanshiCcode syntax/cpp.vim
syn keyword yanshiAction action
syn keyword yanshiMacro semicolon nosemicolon
syn keyword yanshiStorageClass export intact
syn keyword yanshiTodo contained TODO FIXME XXX
syn match yanshiCpp 'c++'
syn match yanshiActionOperator '[>$@%]'
syn match yanshiCall '\^\w\+\(::\w\+\)\?'
syn match yanshiCollapse '!\w\+\(::\w\+\)\?'
syn match yanshiHighOp '[+\*?]'
syn match yanshiIdent '\w\+\(::\w\+\)\?'
syn match yanshiCpp display "^c++\s*" skipwhite nextgroup=yanshiBrace
syn match yanshiImport display "^import\s*" contains=yanshiImported
syn match yanshiLowOp '[-&|]'
syn match yanshiSpecial display contained "\\\(x\x\x\|.\|$\)"
syn region yanshiBrace matchgroup=Delimiter start='{' end='}' fold contains=@yanshiCcode
syn region yanshiBracket start='\[' skip=+\\\\\|\\]+ end=']'
syn region yanshiComment start='/\*' end='\*/' keepend contains=@yanshiCommentGroup,@Spell
syn region yanshiImported display contained start="+" skip=+\\\\\|\\"+ end=+"+
syn region yanshiLineComment start='#\|//' skip='\\$' end='$' keepend contains=@yanshiCommentGroup,@Spell
syn region yanshiPreprocess start="#define" skip="\\$" end="$" keepend
syn region yanshiQQString start=+"+ skip=+\\.+ end=+"+ contains=yanshiSpecial
syn region yanshiQString start=+'+ skip=+\\.+ end=+'+

syn region yanshiDefineStmt start='^\w\+\s*[=:]' end='$' skipnl contains=@yanshiExpr,yanshiComment,yanshiLineComment,yanshiParen0

syn cluster yanshiExpr contains=yanshiActionOperator,yanshiBrace,yanshiBracket,yanshiCall,yanshiCollapse,yanshiIdent,yanshiHighOp,yanshiLowOp,yanshiQString,yanshiQQString,
sy region yanshiParen0 matchgroup=yanshiParen0 start='(' end=')' contains=@yanshiExpr,yanshiParen1
sy region yanshiParen1 matchgroup=yanshiParen1 start='(' end=')' contains=@yanshiExpr,yanshiParen2 contained
sy region yanshiParen2 matchgroup=yanshiParen2 start='(' end=')' contains=@yanshiExpr,yanshiParen3 contained
sy region yanshiParen3 matchgroup=yanshiParen3 start='(' end=')' contains=@yanshiExpr,yanshiParen4 contained
sy region yanshiParen4 matchgroup=yanshiParen4 start='(' end=')' contains=@yanshiExpr,yanshiParen5 contained
sy region yanshiParen5 matchgroup=yanshiParen5 start='(' end=')' contains=@yanshiExpr,yanshiParen0 contained
hi yanshiParen0 ctermfg=brown guifg=#3bb9ff
hi yanshiParen1 ctermfg=DarkBlue guifg=#f88017
hi yanshiParen2 ctermfg=darkgray guifg=#5efb6e
hi yanshiParen3 ctermfg=darkgreen guifg=#f62817
hi yanshiParen4 ctermfg=darkcyan guifg=#fdd017
hi yanshiParen5 ctermfg=darkmagenta guifg=#faafba

hi link yanshiIdent          Identifier
"TODO color mismatch of {}
"hi link yanshiBrace          Statement
"hi link yanshiDefineStmt     Statement
hi def link yanshiCall           Constant
hi def link yanshiCollapse       Constant
hi def link yanshiAction         Structure
hi def link yanshiActionOperator Type
hi def link yanshiBracket        Function
hi def link yanshiCpp            Structure
hi def link yanshiComment        Comment
hi def link yanshiHighOp         Operator
hi def link yanshiImport         Include
hi def link yanshiImported       String
hi def link yanshiLineComment    Comment
hi def link yanshiLowOp          Conditional
hi def link yanshiMacro          Macro
hi def link yanshiPreprocess     Macro
hi def link yanshiQQString       String
hi def link yanshiQString        String
hi def link yanshiSpecial        SpecialChar
hi def link yanshiStorageClass   StorageClass
hi def link yanshiTodo           Todo

let b:current_syntax = 'yanshi'
