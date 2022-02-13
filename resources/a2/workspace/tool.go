///usr/bin/env go run "$0" "$@"; exit "$?"

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
		}
	}()

	config := struct {
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
	}

	if config.table {
		rules, _, _, firsts, follows := parseFile(config.fileHandle)
		config.fileHandle.Close()
		compileTable(os.Stdout, rules, firsts, follows)
		os.Exit(0)
	}

	rules, terms, nonterms, firsts, follows := parseFile(config.fileHandle)
	config.fileHandle.Close()
	printAll(os.Stdout, rules, terms, nonterms, firsts, follows)
	os.Exit(0)
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
) {

	type key struct{ a, t string }
	entries := make(map[key][]string, 512)

	fprintentry := func(a, t, rhs string) {
		rhs = fmt.Sprintf("token.Rule{%v, []Kind{%v}", a, rhs)
		entry := fmt.Sprintf("{%v, %v}: %v}", a, t, rhs)
		entries[key{a, t}] = append(entries[key{a, t}], entry)
		fmt.Fprintf(fh, "	%v,\n", entry)
	}

	// Returns the first set of an entire RHS
	first := func(rhs []string) map[string]struct{} {
		firstSet := make(map[string]struct{}, len(rhs)*2)
		add := func(symbol string) { firstSet[symbol] = struct{}{} }
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

	fmt.Fprintln(fh, "map[Key][]token.Rule{")

	for _, rs := range rules {
		for _, r := range rs {
			alpha := strings.Join(r.RHS, ", ")
			a := r.LHS
			firstalpha := first(r.RHS)

			// 2. For all terminals in FIRST(t), add r.RHS to TT[a, t]
			for t := range firstalpha {
				if t != EPSILON && t != DOLLAR_SIGN {
					fprintentry(a, t, alpha)
				}
			}

			// 3. If EPSILON in firstalpha, for all t in FOLLOW(a), add r.RHS TT[a, t]
			if setContains(firstalpha, EPSILON) {
				followa, ok := follows[a]
				if !ok {
					// All nonterminals should be in the FOLLOW set. If you hit
					// this panic, then there is an error in the grammar
					panic(fmt.Errorf("no FOLLOW entry for '%v'", a))
				}
				for t := range followa {
					if t != EPSILON && t != DOLLAR_SIGN {
						fprintentry(a, t, alpha)
					}
				}
			}
		}
	}
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
	compileTable(fh, rules, firsts, follows)
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
) {
	data, err := io.ReadAll(inputFile)
	if err != nil {
		panic(err)
	}

	contents := string(data)
	rules = scanRules(contents)
	terminals, nonterminals = scanTerminalsAndNonterminals(contents)
	firsts = firstSets(nonterminals, terminals, rules)
	follows = followSets(nonterminals, terminals, rules, firsts)

	return rules, terminals, nonterminals, firsts, follows
}

// Computes all FIRST sets for the provided terminals and nonterminals
func firstSets(
	nonterminals StringSet,
	terminals StringSet,
	rules Rules,
) (set map[string]StringSet) {
	set = make(map[string]StringSet, len(nonterminals)+len(terminals))
	for nonterminal := range nonterminals {
		set[nonterminal] = firstSet(nonterminal, rules, terminals)
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
	set = make(map[string]StringSet, len(nonterminals))
	for nonterminal := range nonterminals {
		followSet(
			set, make(map[string]struct{}, 64),
			nonterminal, rules, terminals, firstSets)
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
				r.RHS = append(r.RHS, word)
			}
		}
		return r
	}

	rules := make(map[string][]Rule, 100)
	scnr := bufio.NewScanner(bytes.NewBufferString(input))
	for scnr.Scan() {
		line := scnr.Text()
		if line == "" {
			continue
		}
		rule := scanRule(line)
		rules[rule.LHS] = append(rules[rule.LHS], rule)
	}

	return rules
}

func scanTerminalsAndNonterminals(input string) (terminals, nonterminals StringSet) {
	nonterminals = make(StringSet, 100)
	terminals = make(StringSet, 100)
	anythingExceptSpace := `[(){}\[\]\w\d,\.:;\"\\\/|\*\&\^\%\$\#\@\!\~\+\=\-\_]`
	terminalRegex := regexp.MustCompile(fmt.Sprintf(`('%s*'|EPSILON)`, anythingExceptSpace))
	nonterminalRegex := regexp.MustCompile(fmt.Sprintf(`<%s*>`, anythingExceptSpace))
	scnr := bufio.NewScanner(bytes.NewBufferString(input))

	for scnr.Scan() {
		line := scnr.Text()
		for _, terminal := range terminalRegex.FindAllString(line, -1) {
			terminals[terminal] = struct{}{}
		}
		for _, nonterminal := range nonterminalRegex.FindAllString(line, -1) {
			nonterminals[nonterminal] = struct{}{}
		}
	}

	return terminals, nonterminals
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

// UNUSED: some unused extra functions that may be handy
var EXTRA_FUNCTIONS = struct {
	compilerFirstOrFollowSets func(*os.File, map[string]StringSet)
	compileRules              func(*os.File, Rules)
}{
	compilerFirstOrFollowSets: func(fh *os.File, firstOrFollowSet map[string]StringSet) {
		variableIfy := func(s string) (ret string) {
			convert := map[string]string{
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
				"NEQ":      "NOTEQ",
				"INTLIT":   "INTNUM",
				"FLOATLIT": "FLOATNUM",
			}

			if strings.Contains(s, "'") && !regexp.MustCompile(`\w`).MatchString(s) {
				if c, ok := convert[s]; ok {
					ret = c
				} else {
					ret = s
				}
			} else {
				ret = strings.ToUpper(strings.ReplaceAll(strings.Trim(s, "<>'"), "-", ""))
			}
			if c, ok := convert[ret]; ok {
				ret = c
			}
			return ret
		}
		for k, v := range firstOrFollowSet {
			if k == "'$'" {
				continue
			}
			lhs := variableIfy(k)
			rhs := make([]string, 0, len(v))
			for s := range v {
				if s != "'$'" {
					rhs = append(rhs, variableIfy(s)+": {}")
				}
			}
			fmt.Fprintf(fh, "%v: {%v},\n", lhs, strings.Join(rhs, ", "))
		}
	},

	// Prints rules as a Go map[string]Rule declaration
	compileRules: func(fh *os.File, rules Rules) {
		fmt.Fprintln(fh, "Rules{")
		variableIfy := func(s string) string {
			return strings.ToUpper(strings.ReplaceAll(strings.Trim(s, "<>'"), "-", ""))
		}
		for _, rs := range rules {
			if len(rs) < 1 {
				continue
			}
			rulesToVar := make([]string, 0, len(rs))
			for _, r := range rs {
				varl := variableIfy(r.LHS)
				varr := make([]string, 0, len(r.RHS))
				for _, s := range r.RHS {
					ss := variableIfy(s)
					varr = append(varr, ss)
				}
				varrs := strings.Join(varr, ", ")
				print := fmt.Sprintf("{%v, []Kind{%v}}", varl, varrs)
				rulesToVar = append(rulesToVar, print)
			}
			rulesToVars := strings.Join(rulesToVar, ", ")
			lhs := variableIfy(rs[0].LHS)
			printout := fmt.Sprintf(`	%v: []Rule{%v},`, lhs, rulesToVars)
			fmt.Fprintln(fh, printout)
		}
		fmt.Fprintln(fh, "}")
	},
}
