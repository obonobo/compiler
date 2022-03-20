///usr/bin/env -S TOOL="$0" go run "$0" "$@"; exit "$?"

// *****************************************************************************
// PARSER BOOTSTRAPPER SCRIPT
//
// This script consumes disambiguated grammar and produces the parsing table.
//
// Run like so: ./tool.go [flags] <grammar file>
//
// If no <grammar_file> is specified, reads from stdin.
//
// Flags:
// 	--all, -a
// 		Default action. Prints out all information collected.
//
// 	--table, -t
// 		Print only the compiled parser table.
//
// 	--compile, -c
// 		Compile everything.
//
// 	--out, -o
//		Output file for generated code
//
// The script will parse the source grammar file and output info about it
// including all the rules, terminals, and nonterminals parsed from the source
// grammar, as well as the FIRST, FOLLOW sets, and the final parsing table.
//
// It is not a full parser generator - you have to write the algorithm yourself,
// but it will do the hard work of generating the table after you disambiguate
// your grammar.
// *****************************************************************************

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
)

var USAGE = strings.TrimLeft(`
Usage: %v [flags] <grammar file>

%v parses a source grammar file and outputs info about it including all the
rules, terminals, and nonterminals parsed from the source grammar, as well as
the FIRST, FOLLOW sets, and the final parsing table.

If no <grammar file> is specified, reads from stdin.

Flags:
	--all, -a
		Default action. Prints out all information collected.

	--table, -t
		Print only the compiled parser table.

	--compile, -c
		Compile everything.

	--out, -o
		Output file for generated code
`, "\n\r\t ")

const (
	EPSILON     = "EPSILON"
	DOLLAR_SIGN = "'$'"
)

var hr = strings.Repeat("-", 40)

type Rule struct {
	LHS string
	RHS []string
}

type Rules = map[string][]Rule
type StringSet = map[string]struct{}

func main() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}()

	config := struct {
		out        string
		compile    bool
		all, table bool
		file, prog string
		fileHandle *os.File
	}{}

	config.prog = path.Base(os.Args[0])
	flag.Usage = func() {
		fmt.Printf(
			USAGE, config.prog,
			strings.ToUpper(string(config.prog[0]))+config.prog[1:])
	}

	flag.BoolVar(&config.all, "all", true, "")
	flag.BoolVar(&config.all, "a", true, "")
	flag.BoolVar(&config.table, "table", false, "")
	flag.BoolVar(&config.table, "t", false, "")
	flag.BoolVar(&config.compile, "compile", false, "")
	flag.BoolVar(&config.compile, "c", false, "")
	flag.StringVar(&config.out, "out", "", "")
	flag.StringVar(&config.out, "o", "", "")
	flag.Parse()

	config.file = flag.Arg(0)
	if config.file == "" {
		fmt.Fprintln(os.Stderr, "No file provided, reading from stdin...")
		config.fileHandle = os.Stdin
	} else {
		fh, err := os.Open(config.file)
		if err != nil {
			panic(fmt.Sprintf("Failed to open file: %v", err))
		}
		config.fileHandle = fh
		defer config.fileHandle.Close()
	}

	rules, terms, nonterms, firsts, follows, semActions := parseFile(config.fileHandle)
	config.fileHandle.Close()

	toolpath, ok := os.LookupEnv("TOOL")
	if !ok {
		toolpath = os.Args[0]
	}

	out := os.Stdout
	if config.out != "" && config.out != "-" {
		fh, err := os.Create(config.out)
		if err != nil {
			panic(fmt.Sprintf("Failed to open output file ('%v'): %v", config.out, err))
		}
		out = fh
	}

	switch {
	case config.compile:
		compileAll(
			out, rules, terms,
			nonterms, firsts, follows,
			toolpath, config.file, semActions)
	case config.table:
		compileTable(os.Stdout, rules, firsts, follows, false)
	default:
		printAll(os.Stdout, rules, terms, nonterms, firsts, follows)
	}

	os.Exit(0)
}

func compileAll(
	fh *os.File,
	rules map[string][]Rule,
	terminals, nonterminals StringSet,
	firsts, follows map[string]StringSet,
	toolpath, grammarfile string,
	semanticActions StringSet,
) {
	compileTypesKindsAndTerminals(fh, toolpath, grammarfile)
	fmt.Fprintln(fh)

	compileNonterminals(fh, nonterminals)
	fmt.Fprintln(fh)

	compileSemanticActions(fh, semanticActions)
	fmt.Fprintln(fh)

	compileTerminals(fh)
	fmt.Fprintln(fh)

	compileRules(fh, rules)
	fmt.Fprintln(fh)

	compileFirstOrFollowSet(fh, "FIRSTS", firsts)
	fmt.Fprintln(fh)

	compileFirstOrFollowSet(fh, "FOLLOWS", follows)
	fmt.Fprintln(fh)

	compileTable(fh, rules, firsts, follows, true)
}

func compileSemanticActions(fh *os.File, semanticActions StringSet) {
	entries := orderifyAndVariablify(semanticActions)

	fmt.Fprintln(fh, "// SEMANTIC ACTIONS")
	fmt.Fprintln(fh, "const (")
	for _, e := range entries {
		fmt.Fprintf(fh, `	%v Kind = "%v"`+"\n", e.variablified, e.original)
	}
	fmt.Fprintln(fh, ")")
	fmt.Fprintln(fh)

	fmt.Fprintln(fh, "var semActions = SEMANTIC_ACTIONS()")
	fmt.Fprintln(fh, "var SEMANTIC_ACTIONS = func() KindSet {")
	fmt.Fprintln(fh, "	return KindSet{")
	for _, e := range entries {
		fmt.Fprintf(fh, "		%v: {},\n", e.variablified)
	}
	fmt.Fprintln(fh, "	}")
	fmt.Fprintln(fh, "}")

	fmt.Fprintf(fh, `
// Returns true if the symbol is a semantic action, false otherwise
func IsSemAction(symbol Kind) bool {
	_, ok := semActions[symbol]
	return ok
}

// The default semantic action is to push a new node on the stack
func defaultSemAction(stack *[]*ASTNode, action Kind, tok Token) {
	pushNode(stack, action, tok)
}

// Invokes the defaultSemAction function, of the function from the override map
// if available
func defaultSemActionOrOverride(lookup Kind, tok Token, stack *[]*ASTNode) {
	if act, ok := semDisptachOverride[lookup]; ok {
		act(stack, tok)
		return
	}
	defaultSemAction(stack, lookup, tok)
}
`)

	// Generate stubs for semantic actions. The default functionality will be to
	// pop the stack and place a
	fmt.Fprintln(fh, "// Default action is to pop, change type, and repush")
	fmt.Fprintln(fh, "var SEM_DISPATCH = map[Kind]SemanticAction{")
	for i := 0; i < len(entries); i++ {
		e := entries[i]
		if i > 0 {
			fmt.Fprintln(fh)
		}
		fmt.Fprintf(fh, "	%v: func(stack *[]*ASTNode, tok Token) {", e.variablified)
		fmt.Fprintf(fh, `defaultSemActionOrOverride(%v, tok, stack)`+"\n", e.variablified)
		fmt.Fprintln(fh, "	},")
	}
	fmt.Fprintln(fh, "}")
}

// These terminals are restricted based on the tokens that I coded by hand when
// making the lexer, so its easiest to just hardcode these
func compileTerminals(fh *os.File) {
	fmt.Fprint(fh, `var terminals = TERMINALS()
var TERMINALS = func() KindSet {
	return KindSet{
		EPSILON:   {},
		OPENPAR:   {},
		CLOSEPAR:  {},
		OPENCUBR:  {},
		CLOSECUBR: {},
		OPENSQBR:  {},
		CLOSESQBR: {},
		AND:       {},
		FLOAT:     {},
		DIV:       {},
		PUBLIC:    {},
		ELSE:      {},
		INTEGER:   {},
		INTNUM:    {},
		WHILE:     {},
		ID:        {},
		EQ:        {},
		VOID:      {},
		COMMA:     {},
		LET:       {},
		MULT:      {},
		SEMI:      {},
		THEN:      {},
		STRUCT:    {},
		FLOATNUM:  {},
		WRITE:     {},
		GT:        {},
		PLUS:      {},
		IMPL:      {},
		MINUS:     {},
		ASSIGN:    {},
		LEQ:       {},
		OR:        {},
		PRIVATE:   {},
		IF:        {},
		COLON:     {},
		NOTEQ:     {},
		LT:        {},
		DOT:       {},
		GEQ:       {},
		READ:      {},
		RETURN:    {},
		NOT:       {},
		INHERITS:  {},
		FUNC:      {},
		ARROW:     {},
	}
}

func IsTerminal(symbol Kind) bool {
	_, ok := terminals[symbol]
	return ok
}
`)
}

func compileNonterminals(fh *os.File, nonterminals StringSet) {
	entries := orderifyAndVariablify(nonterminals)
	fmt.Fprintln(fh, "const (")
	for _, e := range entries {
		fmt.Fprintf(fh, `	%v Kind = "%v"`+"\n", e.variablified, e.original)
	}
	fmt.Fprintln(fh, ")")
	fmt.Fprintln(fh)

	fmt.Fprintln(fh, "var nonterminals = NONTERMINALS()")
	fmt.Fprintln(fh, "var NONTERMINALS = func() KindSet {")
	fmt.Fprintln(fh, "	return KindSet{")
	for _, e := range entries {
		fmt.Fprintf(fh, "		%v: {},\n", e.variablified)
	}
	fmt.Fprintln(fh, "	}")
	fmt.Fprintln(fh, "}")

	fmt.Fprintf(fh, `
func IsNonterminal(symbol Kind) bool {
	_, ok := nonterminals[symbol]
	return ok
}
`)
}

func sortedStringSet(set StringSet) []string {
	sorted := make([]string, 0, len(set))
	for k := range set {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)
	return sorted
}

func compileFirstOrFollowSet(fh *os.File, name string, firstOrFollowSet map[string]StringSet) {
	type entry struct {
		key   string
		value StringSet
	}
	ordered := make([]entry, 0, len(firstOrFollowSet))
	for k, v := range firstOrFollowSet {
		ordered = append(ordered, entry{k, v})
	}
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].key < ordered[j].key })

	fmt.Fprintf(fh, "var %v = func() map[Kind]KindSet {\n", name)
	fmt.Fprintln(fh, "	return map[Kind]KindSet{")
	for _, e := range ordered {
		k, vv := e.key, sortedStringSet(e.value)
		if k == "'$'" {
			continue
		}
		lhs := variablify(k)
		rhs := make([]string, 0, len(vv))
		for _, s := range vv {
			if s != "'$'" {
				rhs = append(rhs, variablify(s)+": {}")
			}
		}
		fmt.Fprintf(fh, "		%v: {%v},\n", lhs, strings.Join(rhs, ", "))
	}
	fmt.Fprintln(fh, "	}")
	fmt.Fprintln(fh, "}")
}

func compileRules(fh *os.File, rules Rules) {
	sortedRules := sortedRules(rules)

	fmt.Fprintf(fh, `// Returns a printout of the rules
func RulesToString(rules Rules) string {
	ret := ""
	for _, v := range rules {
		for _, r := range v {
			rhs := make([]string, 0, len(r.RHS))
			for _, k := range r.RHS {
				rhs = append(rhs, string(k))
			}
			ret += fmt.Sprintf("%v ::= %v\n", r.LHS, strings.Join(rhs, " "))
		}
	}
	return ret
}
`, "%v", "%v")

	fmt.Fprint(fh, `
var RULES = func() Rules {
`)
	fmt.Fprintln(fh, "	return Rules{")
	// for _, rs := range rules {
	for _, e := range sortedRules {
		rs := e.value
		if len(rs) < 1 {
			continue
		}
		rulesToVar := make([]string, 0, len(rs))
		for _, r := range rs {
			varl := variablify(r.LHS)
			varr := make([]string, 0, len(r.RHS))
			for _, s := range r.RHS {
				ss := variablify(s)
				varr = append(varr, ss)
			}
			varrs := strings.Join(varr, ", ")
			print := fmt.Sprintf("{%v, []Kind{%v}}", varl, varrs)
			rulesToVar = append(rulesToVar, print)
		}
		rulesToVars := strings.Join(rulesToVar, ", ")
		lhs := variablify(rs[0].LHS)
		printout := fmt.Sprintf(`		%v: []Rule{%v},`, lhs, rulesToVars)
		fmt.Fprintln(fh, printout)
	}
	fmt.Fprintln(fh, "	}")
	fmt.Fprintln(fh, "}")
}

// Converts grammar symbols to Go variables
func variablify(s string) (ret string) {
	translate := map[string]string{
		"'+'":      "PLUS",
		"'-'":      "MINUS",
		"'*'":      "MULT",
		"'/'":      "DIV",
		"'='":      "ASSIGN",
		"'['":      "OPENSQBR",
		"']'":      "CLOSESQBR",
		"'{'":      "OPENCUBR",
		"'}'":      "CLOSECUBR",
		"'('":      "OPENPAR",
		"')'":      "CLOSEPAR",
		"';'":      "SEMI",
		"'.'":      "DOT",
		"':'":      "COLON",
		"'::'":     "DOUBLECOLON",
		"','":      "COMMA",
		"'->'":     "ARROW",
		"NEQ":      "NOTEQ",
		"INTLIT":   "INTNUM",
		"FLOATLIT": "FLOATNUM",
	}

	if strings.Contains(s, "'") && !regexp.MustCompile(`\w`).MatchString(s) {
		if c, ok := translate[s]; ok {
			ret = c
		} else {
			ret = s
		}
	} else {
		ret = strings.ToUpper(strings.ReplaceAll(strings.Trim(s, "<>'()"), "-", "_"))
	}
	if c, ok := translate[ret]; ok {
		ret = c
	}
	return ret
}

type ruleEntry struct {
	key   string
	value []Rule
}

// Cache for storing certain function return values
var cache = struct{ sortedRules []ruleEntry }{}

func sortedRules(rules map[string][]Rule) []ruleEntry {
	if cache.sortedRules != nil {
		return cache.sortedRules
	}

	entries := make([]ruleEntry, 0, len(rules))
	for k, v := range rules {
		entries = append(entries, ruleEntry{k, v})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].key < entries[j].key
	})

	cache.sortedRules = entries
	return entries
}

// Constructs a parser table from the provided rules, first and follow sets, in
// the format of a Go map[Key][]token.Rule variable which can be copy and pasted
// into your program. Reports on any ambiguous table entries (duplicate map
// keys). If ambiguous entries are found, the map variable will not be legal due
// to the duplicate keys.
func compileTable(
	fh *os.File,
	rules map[string][]Rule,
	firsts, follows map[string]StringSet,
	variablifyEnabled bool,
) {
	// Sort the rules firstly
	sortedRules := sortedRules(rules)

	type key struct{ a, t string }
	entries := make(map[key][]string, 512)

	fprintentry := func(a, t string, rhs []string) {
		if variablifyEnabled {
			for i, r := range rhs {
				rhs[i] = variablify(r)
			}
			a = variablify(a)
			t = variablify(t)
		}
		rhsPrint := fmt.Sprintf("{%v, []Kind{%v}", a, strings.Join(rhs, ", "))
		entry := fmt.Sprintf("{%v, %v}: %v}", a, t, rhsPrint)
		entries[key{a, t}] = append(entries[key{a, t}], entry)
		fmt.Fprintf(fh, "		%v,\n", entry)
	}

	fmt.Fprint(fh, `type Key struct {
	Nonterminal Kind
	Terminal    Kind
}
`)

	fmt.Fprintln(fh)
	fmt.Fprintln(fh, "var TABLE = func() map[Key]Rule {")
	fmt.Fprintln(fh, "	return map[Key]Rule{")

	// for _, rs := range rules {
	// for _, rs := range rules {
	for _, e := range sortedRules {
		rs := e.value

		for _, r := range rs {
			a := r.LHS
			firstalpha := first(firsts, r.RHS)
			sortedFirstAlpha := sortedStringSet(firstalpha)

			// 2. For all terminals in FIRST(t), add r.RHS to TT[a, t]
			for _, t := range sortedFirstAlpha {
				if t != EPSILON && t != DOLLAR_SIGN {
					fprintentry(a, t, r.RHS)
				}
			}

			// 3. If EPSILON in firstalpha, for all t in FOLLOW(a), add r.RHS TT[a, t]
			if setContains(firstalpha, EPSILON) {
				followa, ok := follows[a]
				sortedFollowA := sortedStringSet(followa)

				if !ok {
					// All nonterminals should be in the FOLLOW set. If you hit
					// this panic, then there is an error in the grammar
					panic(fmt.Errorf("no FOLLOW entry for '%v'", a))
				}
				for _, t := range sortedFollowA {
					if t != EPSILON && t != DOLLAR_SIGN {
						fprintentry(a, t, r.RHS)
					}
				}
			}
		}
	}

	fmt.Fprintln(fh, "	}")
	fmt.Fprintln(fh, "}")

	// Report on the ambiguous entries found
	printout, printMe := "AMBIGUITIES ENCOUNTERED\n", false
	for _, entry := range entries {
		if len(entry) >= 2 {
			printMe = true
			for _, e := range entry {
				printout += fmt.Sprintf("%v\n", e)
			}
		}
	}
	if printMe {
		fmt.Fprintf(fh, "\n%v", printout)
	}
}

// This functions generates a header specifying which tool + grammar was used in
// the codegen, then it generates the type declarations we will need, followed
// by a series of hardcoded constants that correspond to the constants that I
// created by hand when I coded the scanner - these
func compileTypesKindsAndTerminals(fh *os.File, toolpath, grammarfile string) {
	fmt.Fprintf(fh, `package token

//
// CODEGEN - DO NOT MODIFY
//
// TOOL:    %v
// GRAMMAR: %v
//
// This file was generated by a tool, it should not be modified by hand. Instead,
// modify the grammar file listed above and rerun the codegen tool.
//

import (
	"fmt"
	"strings"
)

type SemanticAction func(stack *[]*ASTNode, tok Token)

type Kind string

type StringSet = map[string]struct{}

type KindSet = map[Kind]struct{}

type Rules = map[Kind][]Rule

type Rule struct {
	LHS Kind   // The Left Hand Side nonterminal symbol for this rule
	RHS []Kind // The RHS sentential form for this rule
}

func (r Rule) String() string {
	var rhs string
	for _, r := range r.RHS {
		rhs += string(r)
	}
	return fmt.Sprintf("%v ::= %v", r.LHS, rhs)
}

`, toolpath, grammarfile, "%v", "%v")

	fmt.Fprintln(fh)
	fmt.Fprint(fh, `const (
	EPSILON Kind = "EPSILON" // Empty string ''
	ASSIGN  Kind = "assign"  // Assignment operator '='
	ARROW   Kind = "arrow"   // Right-pointing arrow operator '->'

	EQ    Kind = "eq"    // Arithmetic operator: equality '=='
	PLUS  Kind = "plus"  // Arithmetic operator: addition '+'
	MINUS Kind = "minus" // Arithmetic operator: subtraction '-'
	MULT  Kind = "mult"  // Arithmetic operator: multiplication '*'
	DIV   Kind = "div"   // Arithmetic operator: division '/'

	LT    Kind = "lt"    // Comparison operator: less than '<'
	NOTEQ Kind = "noteq" // Comparison operator: not equal '<>'
	LEQ   Kind = "leq"   // Comparison operator: less than or equal '<='
	GT    Kind = "gt"    // Comparison operator: greater than '>'
	GEQ   Kind = "geq"   // Comparison operator: greater than or equal '>='

	OR  Kind = "or"  // Logical operator: OR '|'
	AND Kind = "and" // Logical operator: AND '&'
	NOT Kind = "not" // Logical operator: NOT '!'

	OPENPAR   Kind = "openpar"   // Bracket: opening parenthesis '('
	CLOSEPAR  Kind = "closepar"  // Bracket: closing parenthesis ')'
	OPENCUBR  Kind = "opencubr"  // Bracket: opening curly bracket '{'
	CLOSECUBR Kind = "closecubr" // Bracket: closing curly bracket '}'
	OPENSQBR  Kind = "opensqbr"  // Bracket: opening square bracket '['
	CLOSESQBR Kind = "closesqbr" // Bracket: closing square bracket ']'

	DOT        Kind = "dot"        // Period '.'
	COMMA      Kind = "comma"      // Comma ','
	SEMI       Kind = "semi"       // Semicolon ';'
	COLON      Kind = "colon"      // Colon ':'
	COLONCOLON Kind = "coloncolon" // Double colon '::'

	INLINECMT   Kind = "inlinecmt"   // Single-line comment '// ... \n'
	BLOCKCMT    Kind = "blockcmt"    // Multi-line comment '/* ... */'
	CLOSEINLINE Kind = "closeinline" // End of an inline comment '\n'
	CLOSEBLOCK  Kind = "closeblock"  // End of a block comment '*/'
	OPENINLINE  Kind = "openinline"  // Start of an inline comment '//'
	OPENBLOCK   Kind = "openblock"   // Start of a block comment '/*'

	ID       Kind = "id"       // Identifier 'exampleId_123'
	INTNUM   Kind = "intnum"   // Integer '123'
	FLOATNUM Kind = "floatnum" // Floating-point number '1.23'

	IF       Kind = "if"       // Reserved word 'if'
	THEN     Kind = "then"     // Reserved word 'then'
	ELSE     Kind = "else"     // Reserved word 'else'
	INTEGER  Kind = "integer"  // Reserved word 'integer'
	FLOAT    Kind = "float"    // Reserved word 'float'
	VOID     Kind = "void"     // Reserved word 'void'
	PUBLIC   Kind = "public"   // Reserved word 'public'
	PRIVATE  Kind = "private"  // Reserved word 'private'
	FUNC     Kind = "func"     // Reserved word 'func'
	VAR      Kind = "var"      // Reserved word 'var'
	STRUCT   Kind = "struct"   // Reserved word 'struct'
	WHILE    Kind = "while"    // Reserved word 'while'
	READ     Kind = "read"     // Reserved word 'read'
	WRITE    Kind = "write"    // Reserved word 'write'
	RETURN   Kind = "return"   // Reserved word 'return'
	SELF     Kind = "self"     // Reserved word 'self'
	INHERITS Kind = "inherits" // Reserved word 'inherits'
	LET      Kind = "let"      // Reserved word 'let'
	IMPL     Kind = "impl"     // Reserved word 'impl'

	INVALIDID           Kind = "invalidid"           // Error token
	INVALIDNUM          Kind = "invalidnum"          // Error token
	INVALIDCHAR         Kind = "invalidchar"         // Error token
	UNTERMINATEDCOMMENT Kind = "unterminatedcomment" // Error token
)

func Comments() []Kind {
	return []Kind{INLINECMT, BLOCKCMT, CLOSEBLOCK, CLOSEINLINE}
}
`)
}

func first(firsts map[string]StringSet, rhs []string) StringSet {
	firstSet := make(map[string]struct{}, len(rhs)*2)
	add := func(symbol string) { firstSet[symbol] = struct{}{} }

	rhs = filterSemanticActionsOneRhs(rhs)
	for i, s := range rhs {
		firstOfS, ok := firsts[s]
		if !ok {
			// All symbols should be present in the FIRST set, even
			// nonterminals. If you hit this panic, then there is an error
			// in the grammar
			panic(fmt.Errorf("FIRST(%v) not found", s))
		}

		var hasEpsilon bool
		for f := range firstOfS {
			if f == EPSILON {
				hasEpsilon = true
			}
			if f != EPSILON || i == len(rhs)-1 {
				// If we are at the last element, add epsilon
				add(f)
			}
		}

		if !hasEpsilon {
			break
		}
	}
	return firstSet
}

func printAll(
	fh *os.File,
	rules Rules,
	terminals, nonterminals StringSet,
	firsts, follows map[string]StringSet,
) {
	printRules(fh, "RULES", rules)
	fmt.Fprintln(fh)

	printStringSet(fh, "TERMINALS", terminals)
	fmt.Fprintln(fh)

	printStringSet(fh, "NONTERMINALS", nonterminals)
	fmt.Fprintln(fh)

	printFirstOrFollowSet(fh, "FIRST SETS", firsts, true)
	fmt.Fprintln(fh)

	printFirstOrFollowSet(fh, "FOLLOW SETS", follows, false)
	fmt.Fprintln(fh)

	printTable(fh, "TABLE", rules, firsts, follows)
}

func printTable(
	fh *os.File,
	header string,
	rules map[string][]Rule,
	firsts, follows map[string]map[string]struct{},
) {
	printHeader(fh, header)
	compileTable(fh, rules, firsts, follows, false)
}

func printFirstOrFollowSet(
	fh *os.File,
	header string,
	firstOrFollowSet map[string]StringSet,
	first bool,
) {
	printHeader(fh, header)
	for k, v := range firstOrFollowSet {
		stub := "FIRST"
		if !first {
			stub = "FOLLOW"
		}
		printout := fmt.Sprintf("%v(%v) = {", stub, k)
		var some bool
		for s := range v {
			some = true
			printout += fmt.Sprintf("%v, ", s)
		}
		if some {
			printout = printout[:len(printout)-2]
		}
		printout += "}"
		fmt.Fprintln(fh, printout)
	}
}

func iterateRules(rules Rules, consumer func(Rule)) {
	for _, rs := range rules {
		for _, rule := range rs {
			consumer(rule)
		}
	}
}

func printRules(fh *os.File, header string, rules Rules) {
	printHeader(fh, header)
	iterateRules(rules, func(rule Rule) {
		printout := fmt.Sprintf("%v ::= ", rule.LHS)
		for i, rhs := range rule.RHS {
			printout += rhs
			if i != len(rule.RHS)-1 {
				printout += " "
			}
		}
		fmt.Fprintln(fh, printout)
	})
}

func printStringSet(fh *os.File, header string, ss StringSet) {
	printHeader(fh, header)
	for s := range ss {
		fmt.Fprintln(fh, s)
	}
}

func printHeader(fh *os.File, header string) {
	fmt.Fprintln(fh, header)
	fmt.Fprintln(fh, hr)
}

// Parses the provided file producing the terminals and nonterminal sets, the
// FIRST and FOLLOW sets, as well as the set of production rules of the grammar
func parseFile(inputFile io.Reader) (
	rules Rules,
	terminals, nonterminals StringSet,
	firsts, follows map[string]StringSet,
	semanticActions StringSet,
) {
	data, err := io.ReadAll(inputFile)
	if err != nil {
		panic(err)
	}

	contents := string(data)
	rules = scanRules(contents)
	terminals, nonterminals, semanticActions = scanTerminalsAndNonterminals(contents)
	firsts = firstSets(nonterminals, terminals, rules)
	follows = followSets(nonterminals, terminals, rules, firsts)

	return rules, terminals, nonterminals, firsts, follows, semanticActions
}

// Removes the semantic actions from the RHS of all rules. Returns a new set of
// rules with all the semantic actions removed
func filterSemanticActions(rules Rules) Rules {
	out := make(Rules, len(rules))
	for k, v := range rules {
		rhs := make([]Rule, 0, len(v))
		for _, r := range v {
			rhs = append(rhs, Rule{r.LHS, filterSemanticActionsOneRhs(r.RHS)})
		}
		out[k] = rhs
	}
	return out
}

// Filters a single rhs
func filterSemanticActionsOneRhs(rhs []string) []string {
	out := make([]string, 0, len(rhs))
	for _, r := range rhs {
		if !(strings.HasPrefix(r, "(") && strings.HasSuffix(r, ")")) {
			out = append(out, r)
		}
	}
	return out
}

// Computes all FIRST sets for the provided terminals and nonterminals
func firstSets(
	nonterminals StringSet,
	terminals StringSet,
	rules Rules,
) (set map[string]StringSet) {
	newRules := filterSemanticActions(rules)

	set = make(map[string]StringSet, len(nonterminals)+len(terminals))
	for nonterminal := range nonterminals {
		set[nonterminal] = firstSet(nonterminal, newRules, terminals)
	}
	for terminal := range terminals {
		set[terminal] = StringSet{terminal: {}}
	}
	return set
}

// Computes all FOLLOW sets for the provided nonterminals
func followSets(
	nonterminals StringSet,
	terminals StringSet,
	rules Rules,
	firstSets map[string]StringSet,
) (set map[string]StringSet) {
	newRules := filterSemanticActions(rules)

	set = make(map[string]StringSet, len(nonterminals))
	for nonterminal := range nonterminals {
		followSet(
			set, make(map[string]struct{}, 64),
			nonterminal, newRules, terminals, firstSets)
	}

	for terminal := range terminals {
		followSet(
			set, make(map[string]struct{}, 64),
			terminal, newRules, terminals, firstSets)
	}
	return set
}

// This procedure adds the follow set to the cache that you provide. It also
// returns the same value that it adds to the cache. This is to support
// recursion.
//
// The path parameter should be an empty map. The procedure uses the path to
// keep track of its path as it traverses the grammar rules
//
// ALGO:
//
// 1. if A == S then FOLLOW(A) includes {$}
//
// 2. if there exists a rule B -> <alpha>A<beta>, then FOLLOW(A) includes
// FIRST(<beta>) - {<epsilon>}
//
// 3. if there exists a rule B -> <alpha>A<beta> and <beta> can derive
// <epsilon>, then FOLLOW(A) includes FOLLOW(B)
func followSet(
	cache map[string]StringSet,
	path map[string]struct{},
	nonterminal string,
	rules Rules,
	terminals StringSet,
	firstSets map[string]StringSet,
) (follow StringSet) {

	path[nonterminal] = struct{}{}
	defer func() {
		// Only add this result if we are at the root of our search, i.e. if
		// this function was original called with this nonterminal
		if len(path) == 1 {
			cache[nonterminal] = follow
		}
		delete(path, nonterminal) // Remove ourself from the path
	}()

	// If cache already has this nonterminal, return cached value
	if got, ok := cache[nonterminal]; ok {
		return got
	}

	follow = make(StringSet, 32)

	// Returns true if the symbol is <START>
	isStart := func(symbol string) bool {
		return symbol == "<START>"
	}

	// Searches for rules that match B -> <alpha>A<beta>
	existsRulesForProducing := func(symbol string) (Bs, betas []string, ok bool) {
		Bs = make([]string, 0, len(rules))
		betas = make([]string, 0, len(rules))
		for _, v := range rules {
			for _, rule := range v {
				index := -1
				var contains bool
				for i, x := range rule.RHS {
					if x == nonterminal {
						contains = true
						index = i
					}
				}
				if contains {
					// Then we have found a rule that matches B -> <alpha>A<beta>
					ok = true
					Bs = append(Bs, rule.LHS)
					if index < len(rule.RHS)-1 {
						betas = append(betas, rule.RHS[index+1])
					} else {
						betas = append(betas, EPSILON)
					}
				}
			}
		}
		return Bs, betas, ok
	}

	// Determines whether the given symbol can derive <epsilon>
	var canDeriveEpsilon func(string) bool
	canDeriveEpsilon = func(symbol string) bool {
		// Base case: we are epsilon
		if symbol == EPSILON {
			return true
		}

		// Resursive case: we can still produce epsilon if we have a rule where
		// ALL symbols of the RHS can produce epsilon,
		//
		// e.g.: B -> FT; F -> epsilon; T -> epsilon;
		for _, rule := range rules[symbol] {
			all := true
			for _, s := range rule.RHS {
				if !canDeriveEpsilon(s) {
					all = false
					break
				}
			}
			if all {
				return true
			}
		}

		return false // Otherwise, this symbol cannot produce epsilon
	}

	// 1. A == S
	if isStart(nonterminal) {
		// follow["'$'"] = struct{}{}
		return follow // Start is followed by $
	}

	bs, betas, ok := existsRulesForProducing(nonterminal)
	if !ok {
		return follow
	}

	for i, beta := range betas {

		// 2. Add beta's first set
		first := firstSets[beta]
		minusEpsilon := setDiff(first, StringSet{EPSILON: {}})
		for k := range minusEpsilon {
			follow[k] = struct{}{}
		}

		// 3. Add FOLLOW(B) if beta can derive epsilon
		if _, seen := path[bs[i]]; !seen && canDeriveEpsilon(beta) {
			followB := followSet(cache, path, bs[i], rules, terminals, firstSets)
			for k := range followB {
				follow[k] = struct{}{}
			}
		}
	}

	return follow
}

// LOOKAHEAD TREE GENERATOR
func firstSet(nonterminal string, rules Rules, terminals StringSet) StringSet {
	isTerminal := func(terminal string) (yes bool) {
		_, ok := terminals[terminal]
		return ok
	}

	var recurse func(string) []string
	recurse = func(first string) (set []string) {
		if isTerminal(first) { // Base case: we've reached a terminal
			return []string{first}
		}

		// Otherwise, grab the set of productions where this nonterminal is LHS
		rules, ok := rules[first]
		if !ok {
			// All nonterminals should have rules for how to derive a terminal
			// sentential form i.e. a sentence
			panic(fmt.Errorf("no rules found for nonterminal symbol '%v'", first))
		}

		// For each rule, expand the first symbol on RHS. If that symbol's first
		// set does not contain epsilon then we are done - that first set shall
		// be added to our first set. However, if that symbol's first set does
		// contain epsilon, then we must look at the subsequent symbols on the
		// RHS. Look at each symbol one-after-another, stop at the first symbol
		// whose first set does not contain epsilon, add the union of all first
		// sets encountered minus the set {epsilon}. If all symbol first sets on
		// RHS contain epsilon, our first set shall contain the union of all
		// first sets on RHS.
		for _, r := range rules {
			for i, rhs := range r.RHS {
				frhs := recurse(rhs)
				if i == len(r.RHS)-1 { // last element, add everything
					set = append(set, frhs...)
				} else {
					set = append(set, sliceExcept(frhs, EPSILON)...)
				}
				if !sliceContains(frhs, EPSILON) {
					break
				}
			}
		}

		return set
	}

	set := recurse(nonterminal)
	ret := make(StringSet, len(set))
	for _, term := range set {
		ret[term] = struct{}{}
	}
	return ret
}

func scanRules(input string) map[string][]Rule {
	scanRule := func(line string) (r Rule) {
		scnr := bufio.NewScanner(bytes.NewBufferString(line))
		scnr.Split(bufio.ScanWords)
		for first := true; scnr.Scan(); {
			word := scnr.Text()
			if first {
				first = !first
				r.LHS = word
			} else if word != "::=" {
				r.RHS = append(r.RHS, prefixSemanticAction(word))
			}
		}
		return r
	}

	rules := make(map[string][]Rule, 100)
	scnr := bufio.NewScanner(bytes.NewBufferString(input))
	for scnr.Scan() {
		line := scnr.Text()
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		rule := scanRule(line)
		rules[rule.LHS] = append(rules[rule.LHS], rule)
	}

	return rules
}

func prefixSemanticAction(actionSymbol string) string {
	if strings.HasPrefix(actionSymbol, "(") &&
		strings.HasSuffix(actionSymbol, ")") &&
		!strings.HasPrefix(actionSymbol, "(SEM-") {
		return "(SEM-" + actionSymbol[1:]
	}
	return actionSymbol
}

func scanTerminalsAndNonterminals(input string) (terminals, nonterminals, semanticActions StringSet) {
	nonterminals = make(StringSet, 100)
	terminals = make(StringSet, 100)
	semanticActions = make(StringSet, 100)
	anythingExceptSpace := `[(){}\[\]\w\d,\.:;\"\\\/|\*\&\^\%\$\#\@\!\~\+\=\-\_\>]`
	terminalRegex := regexp.MustCompile(fmt.Sprintf(`('%s*'|EPSILON)`, anythingExceptSpace))
	nonterminalRegex := regexp.MustCompile(fmt.Sprintf(`<%s*>`, anythingExceptSpace))
	semanticActionRegex := regexp.MustCompile(fmt.Sprintf(`\(%s*\)`, anythingExceptSpace))
	scnr := bufio.NewScanner(bytes.NewBufferString(input))

	for scnr.Scan() {
		line := scnr.Text()

		// Filter out comments
		if strings.HasPrefix(line, "//") {
			continue
		}

		for _, terminal := range terminalRegex.FindAllString(line, -1) {
			terminals[terminal] = struct{}{}
		}
		for _, nonterminal := range nonterminalRegex.FindAllString(line, -1) {
			nonterminals[nonterminal] = struct{}{}
		}
		for _, semAction := range semanticActionRegex.FindAllString(line, -1) {
			semanticActions[prefixSemanticAction(semAction)] = struct{}{}
		}
	}

	return terminals, nonterminals, semanticActions
}

// Checks if the needle is in your map haystack
func setContains(haystack map[string]struct{}, needle string) bool {
	_, ok := haystack[needle]
	return ok
}

// Linear search to check if needle is in haystack
func sliceContains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

// Returns a new slice containing all elements of s excludng those in the except
// slice
func sliceExcept(s []string, except ...string) []string {
	acc := make([]string, 0, len(s))
	exceptm := make(map[string]struct{}, len(except))
	for _, e := range except {
		exceptm[e] = struct{}{}
	}
	for _, ss := range s {
		if _, ok := exceptm[ss]; !ok {
			acc = append(acc, ss)
		}
	}
	return acc
}

// Returns a new set containing all elements of s1 that are not also in s2
func setDiff(s1, s2 StringSet) StringSet {
	ret := make(StringSet, len(s1))
	for k := range s1 {
		if _, ok := s2[k]; !ok {
			ret[k] = struct{}{}
		}
	}
	return ret
}

type entry struct{ original, variablified string }

func orderifyAndVariablify(set StringSet) []entry {
	entries := make([]entry, 0, len(set))
	for n := range set {
		entries = append(entries, entry{n, variablify(n)})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].variablified < entries[j].variablified
	})
	return entries
}
