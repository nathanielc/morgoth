package validator

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TrueValidator struct{}

func (self TrueValidator) Validate() error {
	return nil
}

type FalseValidator struct{}

func (self FalseValidator) Validate() error {
	return errors.New("false")
}

func TestValidateAllShouldHandlePtr(t *testing.T) {
	assert := assert.New(t)

	type S struct {
		A TrueValidator
	}

	s := S{}
	assert.Nil(s.A.Validate())

	err := ValidateAll(&s)
	assert.Nil(err)
}

func TestValidateAllShouldPassAllTrue(t *testing.T) {
	assert := assert.New(t)

	type S struct {
		A TrueValidator
		B TrueValidator
		C TrueValidator
	}

	s := S{}
	assert.Nil(s.A.Validate())
	assert.Nil(s.B.Validate())
	assert.Nil(s.C.Validate())

	err := ValidateAll(s)
	assert.Nil(err)
}

func TestValidateAllShouldNotPassSomeTrue(t *testing.T) {
	assert := assert.New(t)

	type S struct {
		A TrueValidator
		B FalseValidator
		C TrueValidator
	}

	s := S{}
	assert.Nil(s.A.Validate())
	assert.NotNil(s.B.Validate())
	assert.Nil(s.C.Validate())

	err := ValidateAll(s)
	assert.NotNil(err)
}

func TestValidateAllShouldNotPassAllFalse(t *testing.T) {
	assert := assert.New(t)

	type S struct {
		A FalseValidator
		B FalseValidator
		C FalseValidator
	}

	s := S{}
	assert.NotNil(s.A.Validate())
	assert.NotNil(s.B.Validate())
	assert.NotNil(s.C.Validate())

	err := ValidateAll(s)
	assert.NotNil(err)
}

func TestValidateOneShouldPassAllTrue(t *testing.T) {
	assert := assert.New(t)

	type S struct {
		A TrueValidator
		B TrueValidator
		C TrueValidator
	}

	s := S{}
	assert.Nil(s.A.Validate())
	assert.Nil(s.B.Validate())
	assert.Nil(s.C.Validate())

	err := ValidateOne(s)
	assert.Nil(err)
}

func TestValidateOneShouldPassSomeTrue(t *testing.T) {
	assert := assert.New(t)

	type S struct {
		A TrueValidator
		B FalseValidator
		C TrueValidator
	}

	s := S{}
	assert.Nil(s.A.Validate())
	assert.NotNil(s.B.Validate())
	assert.Nil(s.C.Validate())

	err := ValidateOne(s)
	assert.Nil(err)
}

func TestValidateOneShouldNotPassAllFalse(t *testing.T) {
	assert := assert.New(t)

	type S struct {
		A FalseValidator
		B FalseValidator
		C FalseValidator
	}

	s := S{}
	assert.NotNil(s.A.Validate())
	assert.NotNil(s.B.Validate())
	assert.NotNil(s.C.Validate())

	err := ValidateOne(s)
	assert.NotNil(err)
}
