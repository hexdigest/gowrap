package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMethod_Declaration(t *testing.T) {
	m := Method{Name: "method"}
	assert.Equal(t, "method() ()", m.Declaration())
}

func TestMethod_Signature(t *testing.T) {
	m := Method{
		Name:    "method",
		Params:  []Param{{Name: "args", Type: "...string"}},
		Results: []Param{{Name: "err", Type: "error"}},
	}
	assert.Equal(t, "(args ...string) (err error)", m.Signature())
}

func TestMethod_ReturnStruct(t *testing.T) {
	t.Run("with results", func(t *testing.T) {
		m := Method{
			Name:    "method",
			Results: []Param{{Name: "err"}},
		}
		assert.Equal(t, "return s.err", m.ReturnStruct("s"))
	})

	t.Run("no results", func(t *testing.T) {
		m := Method{
			Name: "method",
		}
		assert.Equal(t, "return", m.ReturnStruct("s"))
	})
}

func TestMethod_HasResults(t *testing.T) {
	t.Run("with results", func(t *testing.T) {
		m := Method{
			Name:    "method",
			Results: []Param{{Name: "err"}},
		}
		assert.True(t, m.HasResults())
	})

	t.Run("no results", func(t *testing.T) {
		m := Method{
			Name: "method",
		}
		assert.False(t, m.HasResults())
	})
}

func TestMethod_HasParams(t *testing.T) {
	t.Run("with params", func(t *testing.T) {
		m := Method{
			Name:   "method",
			Params: []Param{{}},
		}
		assert.True(t, m.HasParams())
	})

	t.Run("no params", func(t *testing.T) {
		m := Method{
			Name: "method",
		}
		assert.False(t, m.HasParams())
	})
}

func TestMethod_ResultsStruct(t *testing.T) {
	m := Method{
		Name:    "method",
		Results: []Param{{Name: "s", Type: "string"}},
	}
	assert.Equal(t, "struct{\ns string}", m.ResultsStruct())
}

func TestMethod_ResultsNames(t *testing.T) {
	m := Method{
		Name:    "method",
		Results: []Param{{Name: "s"}, {Name: "t"}},
	}
	assert.Equal(t, "s, t", m.ResultsNames())
}

func TestMethod_Pass(t *testing.T) {
	t.Run("no results", func(t *testing.T) {
		m := Method{
			Name:   "method",
			Params: []Param{{Name: "s"}, {Name: "t"}},
		}
		assert.Equal(t, "d.method(s, t)\nreturn", m.Pass("d."))
	})

	t.Run("with results", func(t *testing.T) {
		m := Method{
			Name:    "method",
			Params:  []Param{{Name: "s"}, {Name: "t"}},
			Results: []Param{{Name: "err"}, {Name: "error"}},
		}
		assert.Equal(t, "return d.method(s, t)", m.Pass("d."))
	})
}

func TestMethod_Call(t *testing.T) {
	m := Method{
		Name:   "method",
		Params: []Param{{Name: "s"}, {Name: "t"}},
	}
	assert.Equal(t, "method(s, t)", m.Call())
}
