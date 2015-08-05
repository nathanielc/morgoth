package defaults

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"testing"
)

//////////////////////////////////////////////////////////
// SetDefault Tests
//////////////////////////////////////////////////////////

func TestDefaultShouldDefaultErrorOnNonPtr(t *testing.T) {
	assert := assert.New(t)

	type S struct {
		A string
	}

	s := S{}

	_, err := SetDefault(s, "A")
	assert.NotNil(err)

}
func TestDefaultShouldDefaultErrorOnMissingTag(t *testing.T) {
	assert := assert.New(t)

	type S struct {
		A string
	}

	s := S{}

	_, err := SetDefault(&s, "A")
	assert.NotNil(err)

}

func TestDefaultShouldDefaultErrorOnMissingField(t *testing.T) {
	assert := assert.New(t)

	type S struct {
		A string `default:"a"`
	}

	s := S{}

	_, err := SetDefault(&s, "B")
	assert.NotNil(err)

}

func TestDefaultShouldDefaultBool(t *testing.T) {
	assert := assert.New(t)

	type S struct {
		A bool `default:"true"`
	}

	s := S{}

	assert.Equal(s.A, false)
	v, err := SetDefault(&s, "A")
	assert.Nil(err)
	assert.Equal(s.A, true)
	assert.Equal(s.A, v)

}

func TestDefaultShouldNotDefaultInvalidBool(t *testing.T) {
	assert := assert.New(t)

	type S struct {
		A bool `default:"invalid"`
	}

	s := S{}

	assert.Equal(s.A, false)
	_, err := SetDefault(&s, "A")
	assert.NotNil(err)

}
func TestDefaultShouldDefaultInt(t *testing.T) {
	assert := assert.New(t)

	type S struct {
		A int `default:"-1"`
	}

	s := S{}

	assert.Equal(s.A, 0)
	v, err := SetDefault(&s, "A")
	assert.Nil(err)
	assert.Equal(s.A, -1)
	assert.Equal(s.A, v)

}

func TestDefaultShouldNotDefaultInvalidInt(t *testing.T) {
	assert := assert.New(t)

	type S struct {
		A int `default:"1.5"`
	}

	s := S{}

	assert.Equal(s.A, 0)
	_, err := SetDefault(&s, "A")
	assert.NotNil(err)

}

func TestDefaultShouldDefaultUint(t *testing.T) {
	assert := assert.New(t)

	type S struct {
		A uint `default:"1"`
	}

	s := S{}

	assert.Equal(s.A, uint(0))
	v, err := SetDefault(&s, "A")
	assert.Nil(err)
	assert.Equal(s.A, uint(1))
	assert.Equal(s.A, v)

}

func TestDefaultShouldNotDefaultInvalidUint(t *testing.T) {
	assert := assert.New(t)

	type S struct {
		A uint `default:"-1"`
	}

	s := S{}

	assert.Equal(s.A, uint(0))
	_, err := SetDefault(&s, "A")
	assert.NotNil(err)

}

func TestDefaultShouldDefaultFloat(t *testing.T) {
	assert := assert.New(t)

	type S struct {
		A float64 `default:"1.0"`
	}

	s := S{}

	assert.Equal(s.A, 0.0)
	v, err := SetDefault(&s, "A")
	assert.Nil(err)
	assert.Equal(s.A, 1.0)
	assert.Equal(s.A, v)

}

func TestDefaultShouldNotDefaultInvalidFloat(t *testing.T) {
	assert := assert.New(t)

	type S struct {
		A float64 `default:"1.0a"`
	}

	s := S{}

	assert.Equal(s.A, 0.0)
	_, err := SetDefault(&s, "A")
	assert.NotNil(err)

}

func TestDefaultShouldDefaultString(t *testing.T) {
	assert := assert.New(t)

	type S struct {
		A string `default:"1.0"`
	}

	s := S{}

	assert.Equal(s.A, "")
	v, err := SetDefault(&s, "A")
	assert.Nil(err)
	assert.Equal(s.A, "1.0")
	assert.Equal(s.A, v)
}

func TestDefaultShouldDefaultMultipleFields(t *testing.T) {
	assert := assert.New(t)

	type S struct {
		A string `default:"i am a string"`
		B string
		C int `default:"1"`
	}

	s := S{}

	assert.Equal(s.A, "")
	assert.Equal(s.B, "")
	assert.Equal(s.C, 0)

	v, err := SetDefault(&s, "A")
	assert.Nil(err)
	assert.Equal(s.A, "i am a string")
	assert.Equal(s.A, v)

	_, err = SetDefault(&s, "B")
	assert.NotNil(err)
	assert.Equal(s.B, "")

	v, err = SetDefault(&s, "C")
	assert.Nil(err)
	assert.Equal(s.C, 1)
	assert.Equal(s.C, v)
}

//////////////////////////////////////////////////////////
// HasDefault Tests
//////////////////////////////////////////////////////////

func TestHasDefaultShouldWorkWithPtr(t *testing.T) {

	assert := assert.New(t)

	type S struct {
		A int `default:"1"`
	}

	s := S{}
	assert.Equal(s.A, 0)

	def, err := HasDefault(&s, "A")
	assert.Nil(err)
	assert.True(def)
}

func TestHasDefaultShouldWorkWithoutPtr(t *testing.T) {

	assert := assert.New(t)

	type S struct {
		A int `default:"1"`
	}

	s := S{}
	assert.Equal(s.A, 0)

	def, err := HasDefault(s, "A")
	assert.Nil(err)
	assert.True(def)
}

//////////////////////////////////////////////////////////
// SetAllDefaults Tests
//////////////////////////////////////////////////////////

func TestSetAllDefaultsShouldErrorOnNonPtr(t *testing.T) {
	assert := assert.New(t)

	type S struct {
		A string `default:"i am a string"`
	}

	s := S{}
	assert.Equal(s.A, "")

	err := SetAllDefaults(s)
	assert.NotNil(err)

}

func TestSetAllDefaultsShouldDefaultAll(t *testing.T) {
	assert := assert.New(t)

	type S struct {
		A string `default:"i am a string"`
		B string
		C int `default:"1"`
	}

	s := S{}
	assert.Equal(s.A, "")
	assert.Equal(s.B, "")
	assert.Equal(s.C, 0)

	err := SetAllDefaults(&s)
	assert.Nil(err)

	assert.Equal(s.A, "i am a string")
	assert.Equal(s.B, "")
	assert.Equal(s.C, 1)

}

func TestSetAllDefaultsShouldIgnoreNotDefaults(t *testing.T) {
	assert := assert.New(t)

	type S struct {
		A string `default:"i am a string"`
		B string
		C int `default:"1"`
	}

	s := S{
		A: "original",
		B: "non default",
		C: 42,
	}
	assert.Equal(s.A, "original")
	assert.Equal(s.B, "non default")
	assert.Equal(s.C, 42)

	err := SetAllDefaults(&s)
	assert.Nil(err)

	assert.Equal(s.A, "i am a string")
	assert.Equal(s.B, "non default")
	assert.Equal(s.C, 1)

}
