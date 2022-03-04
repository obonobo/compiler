package constructs

import "github.com/obonobo/esac/core/token"

type FuncDef struct {
	Id         token.Kind
	FParamList token.ASTNodeInterface
	Type       token.Kind
	StatBlock  token.ASTNodeInterface
}
