if exists('g:loaded_syntastic_yanshi_yanshi_checker')
  finish
endif
let g:loaded_syntastic_yanshi_yanshi_checker = 1

let s:save_cpo = &cpo
set cpo&vim

fu! SyntaxCheckers_yanshi_yanshi_GetLocList() dict
  let makeprg = self.makeprgBuild({ 'args': '-d0 -c' })

  let errorformat =
        \ '%C  %.%#,'.
        \ '%E%f %l:%c-%*\d error %m,'.
        \ '%E%f %l-%*\d:%c-%*\d error %m,'.
        \ '%W%f %l:%c-%*\d warning %m,'.
        \ '%W%f %l-%*\d:%c-%*\d warning %m'

  return SyntasticMake({
        \ 'makeprg': makeprg,
        \ 'errorformat': errorformat })
endf

call g:SyntasticRegistry.CreateAndRegisterChecker({
      \ 'filetype': 'yanshi',
      \ 'name': 'yanshi'})

let &cpo = s:save_cpo
unlet s:save_cpo
