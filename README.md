J
===

J (short for 'jump') is a simple but extremely handy tool for jumping to
frequently visited directories on the command-line. It records how many times
each directory has been visited and based on that finds the best match for a
given string and `cd`s the user into that directory.

J does this by trying to find a match for the given command-line parameter in
the list of recently visited directories. Those are ranked by your visiting
frequency and thus the first match will lead to you being `cd` into that
directory. For getting into the secondly ranked directory you have to be so
specific that the first one does not match anymore.

This tool comes with a tab completion script so that you can always do `j <TAB>`
to see the list of directories that match your input.

J only works in Bash.

Examples
===

Say I go into the directory `/home/makkes/Documents/articles` very often then I
only need to enter `j` to get to that directory.

If I go to the directories `/home/makkes/Documents` and `/home/makkes/Music`
very often then I enter `j D` for the first one and `j M` for the latter one.
Patterns such as `j oc`, `j Doc`, `j usic` would also work.

Requirements
===

* Bash
* Python 3

Installation
===

1. Copy the files 'jump' and 'jump.py' to somewhere in your `PATH` like
   `/usr/local/bin`.
1. Add this line to your `.bashrc` (replace '/PATH/TO/' with the real path to
   `j.sh`):
   source /PATH/TO/j.sh
1. If you want tab completion, copy the file `j_completion` to
   `/etc/bash_completion.d/` or wherever your distribution searches for the
   completion scripts.

"Why the heck are you doing this? There's autojump, z, ..."
===

I like to craft my own tools because I know how to handle them, how to modify
them and I get to learn new stuff by creating things I can actually use.
