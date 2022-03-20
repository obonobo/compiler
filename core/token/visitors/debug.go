package visitors

import (
	"fmt"

	"github.com/obonobo/esac/core/token"
	"github.com/obonobo/esac/util"
)

// A placeholder function for testing tree traversal
func debugNode(node *token.ASTNode) {
	fmt.Println(node)
}

func debugAll() map[token.Kind]token.Visit {
	m := make(map[token.Kind]token.Visit, len(token.DEFAULT_VISITOR_METHODS))
	for k := range token.DEFAULT_VISITOR_METHODS {
		m[k] = debugNode
	}
	return m
}

func debugOrOverride(overrides ...map[token.Kind]token.Visit) map[token.Kind]token.Visit {
	return mergeMaps(append([]map[token.Kind]token.Visit{debugAll()}, overrides...)...)
}

//
// Just screwing around with generics below
//

// Returns a new map containing all entries from m1, or if an entry with the
// same key exists in m2, then that entry overrides the entry from m1
func mergeMaps[K comparable, V any](maps ...map[K]V) map[K]V {
	mapSize := util.Max(transform(func(m map[K]V) int { return len(m) }, maps...)...)
	merged := make(map[K]V, mapSize)
	for _, mp := range maps {
		for k, v := range mp {
			merged[k] = v
		}
	}
	return merged
}

// A generic `map` function
func transform[T, V any](mapper func(T) V, values ...T) []V {
	ret := make([]V, 0, len(values))
	for _, v := range values {
		ret = append(ret, mapper(v))
	}
	return ret
}
