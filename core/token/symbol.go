package token

type StringSet = map[string]struct{}
type KindSet = map[Kind]struct{}

const (
	EPSILON             Kind = "EPSILON"
	START               Kind = "<START>"
	ASSIGNSTAT          Kind = "<assignStat>"
	IMPLDEF             Kind = "<implDef>"
	RETURNTYPE          Kind = "<returnType>"
	TERM                Kind = "<term>"
	REPTVARIABLE0       Kind = "<rept-variable0>"
	IDNEST              Kind = "<idnest>"
	MEMBERDECL          Kind = "<memberDecl>"
	REPTIMPLDEF3        Kind = "<rept-implDef3>"
	TYPE                Kind = "<type>"
	REPTVARIABLE2       Kind = "<rept-variable2>"
	ARRAYSIZE           Kind = "<arraySize>"
	VARIABLE            Kind = "<variable>"
	SIGN                Kind = "<sign>"
	ASSIGNOP            Kind = "<assignOp>"
	RIGHTRECARITHEXPR   Kind = "<rightrec-arithExpr>"
	ADDOP               Kind = "<addOp>"
	REPTIDNEST1         Kind = "<rept-idnest1>"
	FUNCDEF             Kind = "<funcDef>"
	FUNCBODY            Kind = "<funcBody>"
	REPTFPARAMSTAIL4    Kind = "<rept-fParamsTail4>"
	REPTOPTSTRUCTDECL22 Kind = "<rept-opt-structDecl22>"
	FUNCHEAD            Kind = "<funcHead>"
	STATBLOCK           Kind = "<statBlock>"
	ARITHEXPR           Kind = "<arithExpr>"
	VARDECL             Kind = "<varDecl>"
	FUNCTIONCALL        Kind = "<functionCall>"
	REPTFPARAMS3        Kind = "<rept-fParams3>"
	APARAMSTAIL         Kind = "<aParamsTail>"
	VARDECLORSTAT       Kind = "<varDeclOrStat>"
	MULTOP              Kind = "<multOp>"
	RELOP               Kind = "<relOp>"
	OPTSTRUCTDECL2      Kind = "<opt-structDecl2>"
	REPTSTRUCTDECL4     Kind = "<rept-structDecl4>"
	FPARAMS             Kind = "<fParams>"
	INDICE              Kind = "<indice>"
	REPTFPARAMS4        Kind = "<rept-fParams4>"
	REPTAPARAMS1        Kind = "<rept-aParams1>"
	STRUCTDECL          Kind = "<structDecl>"
	VISIBILITY          Kind = "<visibility>"
	RELEXPR             Kind = "<relExpr>"
	REPTFUNCTIONCALL0   Kind = "<rept-functionCall0>"
	APARAMS             Kind = "<aParams>"
	EXPR                Kind = "<expr>"
	RIGHTRECTERM        Kind = "<rightrec-term>"
	PROG                Kind = "<prog>"
	REPTPROG0           Kind = "<rept-prog0>"
	REPTFUNCBODY1       Kind = "<rept-funcBody1>"
	FUNCDECL            Kind = "<funcDecl>"
	REPTVARDECL4        Kind = "<rept-varDecl4>"
	STATEMENT           Kind = "<statement>"
	FACTOR              Kind = "<factor>"
	FPARAMSTAIL         Kind = "<fParamsTail>"
	STRUCTORIMPLORFUNC  Kind = "<structOrImplOrFunc>"
	REPTSTATBLOCK1      Kind = "<rept-statBlock1>"
)

var terminals = TERMINALS()
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
	}
}

func IsTerminal(symbol Kind) bool {
	_, ok := terminals[symbol]
	return ok
}

var nonterminals = NONTERMINALS()
var NONTERMINALS = func() KindSet {
	return KindSet{
		START:               {},
		ASSIGNSTAT:          {},
		IMPLDEF:             {},
		RETURNTYPE:          {},
		TERM:                {},
		REPTVARIABLE0:       {},
		IDNEST:              {},
		MEMBERDECL:          {},
		REPTIMPLDEF3:        {},
		TYPE:                {},
		REPTVARIABLE2:       {},
		ARRAYSIZE:           {},
		VARIABLE:            {},
		SIGN:                {},
		ASSIGNOP:            {},
		RIGHTRECARITHEXPR:   {},
		ADDOP:               {},
		REPTIDNEST1:         {},
		FUNCDEF:             {},
		FUNCBODY:            {},
		REPTFPARAMSTAIL4:    {},
		REPTOPTSTRUCTDECL22: {},
		FUNCHEAD:            {},
		STATBLOCK:           {},
		ARITHEXPR:           {},
		VARDECL:             {},
		FUNCTIONCALL:        {},
		REPTFPARAMS3:        {},
		APARAMSTAIL:         {},
		VARDECLORSTAT:       {},
		MULTOP:              {},
		RELOP:               {},
		OPTSTRUCTDECL2:      {},
		REPTSTRUCTDECL4:     {},
		FPARAMS:             {},
		INDICE:              {},
		REPTFPARAMS4:        {},
		REPTAPARAMS1:        {},
		STRUCTDECL:          {},
		VISIBILITY:          {},
		RELEXPR:             {},
		REPTFUNCTIONCALL0:   {},
		APARAMS:             {},
		EXPR:                {},
		RIGHTRECTERM:        {},
		PROG:                {},
		REPTPROG0:           {},
		REPTFUNCBODY1:       {},
		FUNCDECL:            {},
		REPTVARDECL4:        {},
		STATEMENT:           {},
		FACTOR:              {},
		FPARAMSTAIL:         {},
		STRUCTORIMPLORFUNC:  {},
		REPTSTATBLOCK1:      {},
	}
}

func IsNonterminal(symbol Kind) bool {
	_, ok := nonterminals[symbol]
	return ok
}

var FIRSTS = func() map[Kind]KindSet {
	return map[Kind]KindSet{
		FUNCTIONCALL:        {ID: {}, EPSILON: {}},
		REPTFUNCBODY1:       {ID: {}, IF: {}, READ: {}, RETURN: {}, LET: {}, WHILE: {}, WRITE: {}, EPSILON: {}},
		FUNC:                {FUNC: {}},
		CLOSECUBR:           {CLOSECUBR: {}},
		SEMI:                {SEMI: {}},
		OR:                  {OR: {}},
		DIV:                 {DIV: {}},
		FUNCDEF:             {FUNC: {}},
		RIGHTRECARITHEXPR:   {PLUS: {}, MINUS: {}, OR: {}, EPSILON: {}},
		REPTVARDECL4:        {EPSILON: {}, OPENSQBR: {}},
		REPTSTATBLOCK1:      {IF: {}, READ: {}, WRITE: {}, ID: {}, EPSILON: {}, WHILE: {}, RETURN: {}},
		ARITHEXPR:           {ID: {}, EPSILON: {}, INTNUM: {}, FLOATNUM: {}, OPENPAR: {}, NOT: {}, PLUS: {}, MINUS: {}},
		REPTVARIABLE2:       {OPENSQBR: {}, EPSILON: {}},
		RELEXPR:             {FLOATNUM: {}, OPENPAR: {}, PLUS: {}, ID: {}, EPSILON: {}, INTNUM: {}, NOT: {}, MINUS: {}},
		STRUCT:              {STRUCT: {}},
		INHERITS:            {INHERITS: {}},
		STATBLOCK:           {WHILE: {}, READ: {}, WRITE: {}, RETURN: {}, ID: {}, EPSILON: {}, IF: {}, OPENCUBR: {}},
		LET:                 {LET: {}},
		AND:                 {AND: {}},
		FLOAT:               {FLOAT: {}},
		VOID:                {VOID: {}},
		RELOP:               {LT: {}, GT: {}, LEQ: {}, GEQ: {}, EQ: {}, NOTEQ: {}},
		MULTOP:              {MULT: {}, DIV: {}, AND: {}},
		REPTFPARAMS3:        {OPENSQBR: {}, EPSILON: {}},
		VARDECL:             {LET: {}},
		FPARAMS:             {ID: {}, EPSILON: {}},
		TYPE:                {INTEGER: {}, FLOAT: {}, ID: {}},
		EPSILON:             {EPSILON: {}},
		REPTOPTSTRUCTDECL22: {COMMA: {}, EPSILON: {}},
		FUNCHEAD:            {FUNC: {}},
		FUNCDECL:            {FUNC: {}},
		ID:                  {ID: {}},
		GEQ:                 {GEQ: {}},
		LEQ:                 {LEQ: {}},
		INTEGER:             {INTEGER: {}},
		STRUCTDECL:          {STRUCT: {}},
		APARAMSTAIL:         {COMMA: {}},
		ELSE:                {ELSE: {}},
		FLOATNUM:            {FLOATNUM: {}},
		OPENSQBR:            {OPENSQBR: {}},
		WRITE:               {WRITE: {}},
		MINUS:               {MINUS: {}},
		NOT:                 {NOT: {}},
		ASSIGN:              {ASSIGN: {}},
		STATEMENT:           {WRITE: {}, ID: {}, IF: {}, READ: {}, EPSILON: {}, WHILE: {}, RETURN: {}},
		REPTFPARAMSTAIL4:    {OPENSQBR: {}, EPSILON: {}},
		REPTSTRUCTDECL4:     {PUBLIC: {}, PRIVATE: {}, EPSILON: {}},
		REPTAPARAMS1:        {EPSILON: {}, COMMA: {}},
		START:               {EPSILON: {}, STRUCT: {}, IMPL: {}, FUNC: {}},
		VISIBILITY:          {PUBLIC: {}, PRIVATE: {}},
		INTNUM:              {INTNUM: {}},
		IDNEST:              {ID: {}},
		REPTIMPLDEF3:        {EPSILON: {}, FUNC: {}},
		FUNCBODY:            {OPENCUBR: {}},
		WHILE:               {WHILE: {}},
		EQ:                  {EQ: {}},
		OPENPAR:             {OPENPAR: {}},
		REPTVARIABLE0:       {ID: {}, EPSILON: {}},
		REPTPROG0:           {STRUCT: {}, IMPL: {}, FUNC: {}, EPSILON: {}},
		APARAMS:             {OPENPAR: {}, NOT: {}, PLUS: {}, MINUS: {}, FLOATNUM: {}, EPSILON: {}, INTNUM: {}, ID: {}},
		CLOSESQBR:           {CLOSESQBR: {}},
		GT:                  {GT: {}},
		FPARAMSTAIL:         {COMMA: {}},
		ASSIGNSTAT:          {ID: {}, EPSILON: {}},
		COLON:               {COLON: {}},
		RIGHTRECTERM:        {MULT: {}, DIV: {}, AND: {}, EPSILON: {}},
		VARDECLORSTAT:       {EPSILON: {}, WRITE: {}, RETURN: {}, LET: {}, ID: {}, IF: {}, WHILE: {}, READ: {}},
		REPTFPARAMS4:        {COMMA: {}, EPSILON: {}},
		IMPL:                {IMPL: {}},
		DOT:                 {DOT: {}},
		NOTEQ:               {NOTEQ: {}},
		PRIVATE:             {PRIVATE: {}},
		PROG:                {IMPL: {}, FUNC: {}, EPSILON: {}, STRUCT: {}},
		ASSIGNOP:            {ASSIGN: {}},
		REPTFUNCTIONCALL0:   {ID: {}, EPSILON: {}},
		REPTIDNEST1:         {OPENSQBR: {}, EPSILON: {}},
		OPENCUBR:            {OPENCUBR: {}},
		IF:                  {IF: {}},
		MULT:                {MULT: {}},
		ADDOP:               {OR: {}, PLUS: {}, MINUS: {}},
		INDICE:              {OPENSQBR: {}},
		ARRAYSIZE:           {OPENSQBR: {}},
		FACTOR:              {EPSILON: {}, INTNUM: {}, OPENPAR: {}, PLUS: {}, ID: {}, FLOATNUM: {}, NOT: {}, MINUS: {}},
		THEN:                {THEN: {}},
		PUBLIC:              {PUBLIC: {}},
		LT:                  {LT: {}},
		READ:                {READ: {}},
		RETURNTYPE:          {FLOAT: {}, ID: {}, VOID: {}, INTEGER: {}},
		IMPLDEF:             {IMPL: {}},
		VARIABLE:            {ID: {}, EPSILON: {}},
		EXPR:                {ID: {}, EPSILON: {}, INTNUM: {}, NOT: {}, FLOATNUM: {}, OPENPAR: {}, PLUS: {}, MINUS: {}},
		TERM:                {EPSILON: {}, FLOATNUM: {}, OPENPAR: {}, ID: {}, INTNUM: {}, NOT: {}, PLUS: {}, MINUS: {}},
		COMMA:               {COMMA: {}},
		RETURN:              {RETURN: {}},
		CLOSEPAR:            {CLOSEPAR: {}},
		STRUCTORIMPLORFUNC:  {STRUCT: {}, IMPL: {}, FUNC: {}},
		SIGN:                {PLUS: {}, MINUS: {}},
		MEMBERDECL:          {FUNC: {}, LET: {}},
		OPTSTRUCTDECL2:      {INHERITS: {}, EPSILON: {}},
		PLUS:                {PLUS: {}},
	}
}

var FOLLOWS = func() map[Kind]KindSet {
	return map[Kind]KindSet{
		EXPR:                {COMMA: {}, SEMI: {}, CLOSEPAR: {}},
		RELEXPR:             {COMMA: {}, SEMI: {}, CLOSEPAR: {}},
		FPARAMS:             {CLOSEPAR: {}},
		STRUCTDECL:          {FUNC: {}, STRUCT: {}, IMPL: {}},
		ASSIGNOP:            {ID: {}, MINUS: {}, NOT: {}, INTNUM: {}, FLOATNUM: {}, PLUS: {}, OPENPAR: {}},
		VISIBILITY:          {LET: {}, FUNC: {}},
		ARRAYSIZE:           {SEMI: {}, OPENSQBR: {}, CLOSEPAR: {}, COMMA: {}},
		STATEMENT:           {RETURN: {}, LET: {}, SEMI: {}, IF: {}, WRITE: {}, WHILE: {}, ID: {}, READ: {}, CLOSECUBR: {}},
		MULTOP:              {MINUS: {}, ID: {}, NOT: {}, INTNUM: {}, PLUS: {}, OPENPAR: {}, FLOATNUM: {}},
		IDNEST:              {ID: {}},
		REPTIMPLDEF3:        {CLOSECUBR: {}},
		REPTFPARAMS4:        {CLOSEPAR: {}},
		REPTFPARAMSTAIL4:    {COMMA: {}, CLOSEPAR: {}},
		VARDECL:             {PRIVATE: {}, IF: {}, RETURN: {}, LET: {}, CLOSECUBR: {}, READ: {}, PUBLIC: {}, WRITE: {}, ID: {}, WHILE: {}},
		APARAMSTAIL:         {COMMA: {}, CLOSEPAR: {}},
		REPTFUNCBODY1:       {CLOSECUBR: {}},
		FUNCDEF:             {STRUCT: {}, CLOSECUBR: {}, FUNC: {}, IMPL: {}},
		FUNCTIONCALL:        {AND: {}, PLUS: {}, OR: {}, MINUS: {}, COMMA: {}, SEMI: {}, CLOSESQBR: {}, CLOSEPAR: {}, DIV: {}, MULT: {}},
		REPTFUNCTIONCALL0:   {ID: {}},
		SIGN:                {NOT: {}, INTNUM: {}, ID: {}, PLUS: {}, FLOATNUM: {}, OPENPAR: {}, MINUS: {}},
		REPTFPARAMS3:        {CLOSEPAR: {}, COMMA: {}},
		REPTVARIABLE2:       {MINUS: {}, OR: {}, MULT: {}, CLOSESQBR: {}, DIV: {}, CLOSEPAR: {}, ASSIGN: {}, COMMA: {}, PLUS: {}, SEMI: {}, AND: {}},
		PROG:                {},
		FUNCHEAD:            {OPENCUBR: {}, SEMI: {}},
		FUNCDECL:            {PUBLIC: {}, CLOSECUBR: {}, PRIVATE: {}},
		ASSIGNSTAT:          {SEMI: {}},
		REPTOPTSTRUCTDECL22: {OPENCUBR: {}},
		ARITHEXPR:           {SEMI: {}, COMMA: {}, CLOSESQBR: {}, CLOSEPAR: {}},
		START:               {},
		INDICE:              {PLUS: {}, DOT: {}, SEMI: {}, ASSIGN: {}, MULT: {}, OPENSQBR: {}, MINUS: {}, CLOSEPAR: {}, COMMA: {}, DIV: {}, OR: {}, AND: {}, CLOSESQBR: {}},
		REPTSTRUCTDECL4:     {CLOSECUBR: {}},
		RETURNTYPE:          {OPENCUBR: {}, SEMI: {}},
		APARAMS:             {CLOSEPAR: {}},
		IMPLDEF:             {IMPL: {}, STRUCT: {}, FUNC: {}},
		REPTVARDECL4:        {SEMI: {}},
		REPTIDNEST1:         {DOT: {}},
		REPTAPARAMS1:        {CLOSEPAR: {}},
		RELOP:               {FLOATNUM: {}, NOT: {}, OPENPAR: {}, MINUS: {}, ID: {}, INTNUM: {}, PLUS: {}},
		VARIABLE:            {COMMA: {}, PLUS: {}, ASSIGN: {}, OR: {}, DIV: {}, MINUS: {}, AND: {}, SEMI: {}, MULT: {}, CLOSEPAR: {}, CLOSESQBR: {}},
		FUNCBODY:            {IMPL: {}, STRUCT: {}, FUNC: {}, CLOSECUBR: {}},
		REPTVARIABLE0:       {ID: {}},
		FACTOR:              {DIV: {}, CLOSESQBR: {}, MULT: {}, CLOSEPAR: {}, COMMA: {}, PLUS: {}, OR: {}, SEMI: {}, AND: {}, MINUS: {}},
		TYPE:                {OPENCUBR: {}, PUBLIC: {}, CLOSEPAR: {}, WHILE: {}, OPENSQBR: {}, RETURN: {}, CLOSECUBR: {}, PRIVATE: {}, SEMI: {}, ID: {}, IF: {}, LET: {}, READ: {}, WRITE: {}, COMMA: {}},
		STATBLOCK:           {SEMI: {}},
		REPTPROG0:           {},
		STRUCTORIMPLORFUNC:  {STRUCT: {}, IMPL: {}, FUNC: {}},
		FPARAMSTAIL:         {CLOSEPAR: {}, COMMA: {}},
		TERM:                {SEMI: {}, MINUS: {}, CLOSEPAR: {}, COMMA: {}, OR: {}, CLOSESQBR: {}, PLUS: {}},
		RIGHTRECARITHEXPR:   {CLOSESQBR: {}, CLOSEPAR: {}, SEMI: {}, COMMA: {}},
		OPTSTRUCTDECL2:      {OPENCUBR: {}},
		ADDOP:               {INTNUM: {}, ID: {}, PLUS: {}, NOT: {}, MINUS: {}, FLOATNUM: {}, OPENPAR: {}},
		REPTSTATBLOCK1:      {CLOSECUBR: {}},
		VARDECLORSTAT:       {ID: {}, CLOSECUBR: {}, RETURN: {}, WRITE: {}, IF: {}, READ: {}, WHILE: {}, LET: {}},
		RIGHTRECTERM:        {PLUS: {}, CLOSEPAR: {}, OR: {}, SEMI: {}, COMMA: {}, MINUS: {}, CLOSESQBR: {}},
		MEMBERDECL:          {PUBLIC: {}, CLOSECUBR: {}, PRIVATE: {}},
	}
}
