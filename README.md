J
===

J (short for 'jump') is a simple but extremely handy tool for jumping to
frequently visited directories on the command-line. It records how many times
each directory has been visited and based on that finds the best match for a
given string and `cd`s the user into that directory.

J does this by trying to find a match for the given command-line parameter in
the list of recently visited directories. Those are ranked by your visiting
frequency and thus the first match will lead to you being `cd`d into that
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

Exact matches are always ranked above all others. So if you have
`/home/makkes/Doc` and `/home/makkes/Documents` in J's database, entering `j
Doc` will always lead you to `/home/makkes/Doc`, even if you jumped to the
`Documents` directory more often.

Requirements
===

* Bash
* Linux or macOS

Installation
===

```sh
curl -o- https://raw.githubusercontent.com/makkes/j/v1.0.7/install.sh | bash
```

or

```sh
wget -qO- https://raw.githubusercontent.com/makkes/j/v1.0.7/install.sh | bash
```

Usage
===

* Just use `j` whenever you used to invoke `cd`. One difference you have to be
  aware of is that pressing TAB twice after typing `j ` will not give you
  completions for directories in the current working dir as `cd` does.

"Why the heck are you doing this? There's autojump, z, ..."
===

I like to craft my own tools because I know how to handle them, how to modify
them and I get to learn new stuff by creating things I can actually use.
