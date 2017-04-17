- [aliason](#aliason)
  - [Using aliason](#using-aliason)
  - [Declaring aliases](#declaring-aliases)

# aliason
aliason is a shell tool for managing project-specific shell aliases.

Installation:
```bash
go get github.com/bpicolo/aliason
```

I created aliason primarily out of the desire to not remember specific test-running
commands for different repos. Task runners like Make have left
me shorthanded because passing arguments can be such a pain, whereas shell aliases
have been useful, but not portable between projects.


## Using aliason
Getting started with aliason is straightforward, simply:

```bash
aliason install >> ~/.bash_profile  # or bashrc/other shell source file of choice
```

and then source that file.

This does two things:
1. adds a `cd` function that, as an alternative to cd, will cd followed by `eval $(aliason env)`
2. adds an `eval $(aliason env)` directly, so creating new shells will also source your env.

These options seemed the most straightforward, but there a couple alternatives. Manually
running `eval $(aliason env)` should pick up the aliases in your current env. There's also
a valid path through PROMPT_COMMAND, and probably any number of other alternatives as well.
Anything that runs `eval $(aliason env)` when desired should do the trick.

## Declaring aliases
aliason will look for and .aliasrc file in your current directory. An .aliasrc file is a
simple mapping of alias names to commands in yaml-syntax.

```yaml
ping: pong
test: tox
```

## todo
1. Support preserving global aliases when moving between directories (or also support not-overwriting?)
2. Tests
