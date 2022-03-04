# Workspace

## _Files:_

- `grammar.grm`: the unmangled (non-LL(1)) grammar

- `grammar-mangled.grm`: the LL(1) version of the grammar, mangled by
  factorization.

- `grammar-mangled-attribute.grm`: the mangled grammar with attributes (semantic
  actions) added to it.

## _Format:_

Semantic actions that are embedded in the attribute grammar are written like so:

    <nonterminal1> ::= 'terminal1' (ACTION1) <nonterminal2>

That is, the semantic action is identified by a symbol written in all caps
surrounded with parentheses `()`.

Then, at the bottom of the file, after all grammar rules are written, the
semantic actions are described:

    (ACTION1)
    ---
    Make a node using the previous token and push it on the semantic stack
