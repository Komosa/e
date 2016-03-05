### Lets start with git commit msg editor
Use cases:
- allow convenient interactive **rebase**
- allow possibility to most basic text editing (**commit** msg)

Assumptions:
- no need to display parts of _too long_ lines (just drop rest)
- no need for any scrolling (both vertical and horizontal)
- no need for any code optimization

Features:
- move by whole words `ctrl+left, ctrl+right`
- delete whole words (rebase!) `alt+d`
- delete/undelete whole line (like in nano, but only one line at time) (rebase!) `ctrl+k, ctrl+u`
- clear everything `ctrl+e`
- quit and save `ctrl+q`
- discard changes `ctrl+]`

At this point I can even say, that this is not fully abstract project - I will use it! (nano doesn't have emacs' _M-d_, and I don't want to use emacs/vim).

I will probably add features in future. Configuration is done be recompilation.
