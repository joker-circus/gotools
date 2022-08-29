package hashset

import (
	"testing"

	"github.com/joker-circus/gotools/internal/template"
)

func TestSetTypesGen(t *testing.T) {
	template.TypesGen("hashset.template", "%sset.go")
}

func TestSafeSetTypesGen(t *testing.T) {
	template.TypesGen("safehashset.template", "safe%sset.go")
}
