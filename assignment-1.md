# ASSIGNMENT 1 - LEXICAL ANALYZER (SCANNER)

## Regular Expressions

**\* Check out the [Regular Expression Rules](#1---regular-expression-rules)
section for the notation used here**

_id:_ `[aA-zZ]([aA-zZ]|[0-9]|_)*`

_integer:_ `(0|[1-9][0-9]*)`

_float:_ `(0|[1-9][0-9]*)\.(0|[0-9]*[1-9])`

## Appendix

### 1 - Regular Expression Rules

I wrote the regular expressions out using a subset of PCRE that is close to the
original lexical specifications.

Notation:

```plaintext
<exp1><exp2>  =  Concatentation

(<exp1>|<exp2>) = Union

<exp>*     =  Kleene's start, repeat <exp> 0 to infinite times

[<c1><c2>...]  = Matches any character within the set

\<c1>   =   Matches any character exactly

[a-z]   =   Matches any alphabetical character between a to z

[A-Z]   =   Same as above but capitalized

[aA-zZ]    =  ([a-z]|[A-Z])

(<exp>)    =  Brackets are a grouping statement

[0-9]      =  Any digit from 0 to 9, includes: 0,1,2,3,4,5,6,7,8,9

[1-9]      =  Any digit from 1 to 9, includes: 1,2,3,4,5,6,7,8,9

<exp>+     =  Repeat <exp> 1 to infinite times

<exp>?     =  Repeat <exp> 0 or 1 time

<exp>{x}   =  Repeat <exp> exactly x times

<exp>{x,y} =  Repeat <exp> between x and y times
```
