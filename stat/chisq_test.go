package stat_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/nvcook42/morgoth/stat"
)

func TestSpecificChiSqInc(t *testing.T) {

	assert := assert.New(t)

	ppf := stat.Xsquare_InvCDF(60)
	assert.Equal(79.66881012367774, ppf(0.95))
}
