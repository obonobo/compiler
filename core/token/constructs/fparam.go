package constructs

import "github.com/obonobo/esac/core/token"

type FParam struct {
	Id      token.Kind
	Type    token.Kind
	DimList DimList
}
