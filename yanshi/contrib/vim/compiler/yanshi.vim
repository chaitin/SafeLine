"if exists('current_compiler')
"  finish
"endif
let current_compiler = 'yanshi'

if exists(':CompilerSet') != 2
  command -nargs=* CompilerSet setlocal <args>
endif
CompilerSet errorformat=
      \%E%f\ %l:%c-%*\\d\ error\ %m,
      \%E%f\ %l-%*\\d:%c-%*\\d\ error\ %m,
      \%W%f\ %l:%c-%*\\d\ warning\ %m,
      \%W%f\ %l-%*\\d:%c-%*\\d\ warning\ %m,
      \%C%.%#
CompilerSet makeprg=yanshi\ -d0\ -c\ $*\ %
