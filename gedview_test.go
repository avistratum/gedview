package gedview

import (
	"strings"
	"testing"
)

func Test_CreateAST(t *testing.T) {
	t.Run("higher level", func(t *testing.T) {
		c := strings.NewReader(`
0 HEAD
1 GEDC
2 VERS 5.5.5
		`)

		ast, err := CreateAST(c)

		if err != nil {
			t.Error(err)
		}

		if ast.Children[0].Children[0].Children[0].Raw != "2 VERS 5.5.5" {
			t.Error("mapping didn't resolve correctly")
		}
	})

	t.Run("equal level", func(t *testing.T) {
		c := strings.NewReader(`
0 HEAD
1 GEDC
2 VERS 5.5.5
2 FORM LINEAGE-LINKED
			`)

		ast, err := CreateAST(c)

		if err != nil {
			t.Error(err)
		}

		if ast.Children[0].Children[0].Children[1].Raw != "2 FORM LINEAGE-LINKED" {
			t.Error("mapping didn't resolve correctly")
		}

	})

	t.Run("lower level", func(t *testing.T) {
		c := strings.NewReader(`
0 HEAD
1 GEDC
2 VERS 5.5.5
2 FORM LINEAGE-LINKED
3 VERS 5.5.5
1 CHAR UTF-8
			`)

		ast, err := CreateAST(c)

		if err != nil {
			t.Error(err)
		}

		if ast.Children[0].Children[1].Raw != "1 CHAR UTF-8" {
			t.Error("mapping didn't resolve correctly")
		}
	})
}
