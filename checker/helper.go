package checker

// Normalized operations
const (
	PLUS   = "PLUS"
	EQUAL  = "EQ"
	LT     = "LT"
	GT     = "GT"
	MINUS  = "MINUS"
	TIMES  = "TIMES"
	DIVIDE = "DIVIDE"
	AND    = "AND"
	OR     = "OR"
	PRINT  = "PRINT"
)

// Variable types
const (
	INT_TYPE     = "Int"
	STRING_TYPE  = "String"
	BOOL_TYPE    = "Bool"
	NOTHING_TYPE = "Nothing"
)

// Signature of return values
type Signature struct {
	Return string
	// List of types
	Params []string
}

/*TypeTable : Methods defining various types */
var TypeTable = map[string]map[string]Signature{
	INT_TYPE: {
		PLUS:  {INT_TYPE, []string{INT_TYPE}},
		MINUS: {INT_TYPE, []string{INT_TYPE}},
		TIMES: {INT_TYPE, []string{INT_TYPE}},
		LT:    {BOOL_TYPE, []string{INT_TYPE}},
		GT:    {BOOL_TYPE, []string{INT_TYPE}},
		EQUAL: {BOOL_TYPE, []string{INT_TYPE}},
		PRINT: {NOTHING_TYPE, []string{}},
	},
	STRING_TYPE: {
		PLUS:  {STRING_TYPE, []string{STRING_TYPE}},
		PRINT: {NOTHING_TYPE, []string{}}},

	BOOL_TYPE: {
		AND:   {BOOL_TYPE, []string{BOOL_TYPE}},
		OR:    {BOOL_TYPE, []string{BOOL_TYPE}},
		PRINT: {NOTHING_TYPE, []string{}}},
}

/*Environment structure */
type Environment struct {
	Vals  map[string]string    // map identifier to type
	Funcs map[string]Signature // map function name to return type
	Types map[string]bool      // track valid types
}

var env Environment

/*IsBuiltin checks for a built-in function (command) */
func IsBuiltin(name string) bool {
	return name == "PRINT"
}

func NewEnvironment() Environment {
	return Environment{Vals: map[string]string{}, Funcs: map[string]Signature{}, Types: map[string]bool{}}
}

// Check if the method in the typetable exists
func MethodExist(kind, method string) bool {
	methods, ok := TypeTable[kind]
	if !ok {
		return false
	}

	_, ok = methods[method]
	return ok
}

func GetMethod(kind, method string) (Signature, bool) {
	methods, ok := TypeTable[kind]
	if !ok {
		return Signature{}, false
	}

	sig, ok := methods[method]
	return sig, ok
}

func (e *Environment) Set(name, kind string) {
	e.Vals[name] = kind
}

func SetFunctionSignature(name string, sig Signature) {
	env.Funcs[name] = sig
}

func GetFunctionSignature(name string) (Signature, bool) {
	kind, ok := env.Funcs[name]
	return kind, ok
}

func (e *Environment) Get(name string) (string, bool) {
	kind, ok := e.Vals[name]
	return kind, ok
}

func (e *Environment) IdentExist(kind string) bool {
	_, ok := e.Vals[kind]
	return ok
}

func GetIdentType(name string) (string, bool) {
	kind, ok := env.Vals[name]
	return kind, ok
}

func (e *Environment) TypeExist(kind string) bool {
	_, ok := e.Types[kind]
	return ok
}
