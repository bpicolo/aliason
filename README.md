# aliason
aliason is a shell tool for managing project-specific shell aliases.

Installation:
```bash
go get github.com/bpicolo/aliason
```

I created aliason primarily out of the desire to not remember specific test-running
commands for different repos. Task runners like Make have left
me shorthanded because passing arguments to them can be a huge pain, whereas shell aliases
are useful, but not portable between projects.


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
proably a valid strategy using PROMPT_COMMAND, and probably any number of other alternatives as well.
Anything that runs `eval $(aliason env)` when desired should do the trick.

## Declaring aliases
aliason will look for an .aliasonrc file in your current directory. An .aliasonrc file is a
simple mapping of alias names to commands in yaml syntax.

```yaml
ping: echo pong
test: tox
```

## todo
1. Support preserving global aliases when moving between directories (or also support not-overwriting?)
2. Quoting is probably not an entirely sufficient security measure, though could be worse. Hopefully people
tend to cd into relatively trusted directories on their boxes... When can this fail horribly? (Other than people
being jerks in a shared repo and overwriting a builtin with 'sudo rm -rf --no-preserve-root /'). Perhaps this is part over
the "not overwriting" above.
3. Tests
