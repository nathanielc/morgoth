package kstest

import (
	"github.com/nathanielc/morgoth"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	defer glog.Flush()
	os.Exit(m.Run())
}

func TestCalcDShouldBe0(t *testing.T) {
	assert := assert.New(t)

	data := make([]float64, 10)
	for i := range data {
		data[i] = float64(i+1) / float64(len(data))
	}

	expectedD := 0.0

	d := calcD(data, data)
	assert.InDelta(expectedD, d, 1e-5)
}

func TestCalcDShouldBeSmall(t *testing.T) {
	assert := assert.New(t)

	data1 := make([]float64, 10)
	for i := range data1 {
		data1[i] = float64(i+1) / float64(len(data1))
	}

	data2 := make([]float64, 10)
	for i := range data2 {
		data2[i] = float64(i) / float64(len(data2))
	}

	expectedD := 0.1

	d := calcD(data1, data2)
	assert.InDelta(expectedD, d, 1e-5)

	d = calcD(data2, data1)
	assert.InDelta(expectedD, d, 1e-5)

}

func TestCalcDShouldBe1(t *testing.T) {
	assert := assert.New(t)

	data1 := make([]float64, 10)
	for i := range data1 {
		data1[i] = 0.0
	}

	data2 := make([]float64, 10)
	for i := range data2 {
		data2[i] = float64(i+1) / float64(len(data2))
	}

	expectedD := 1.0

	d := calcD(data1, data2)
	assert.InDelta(expectedD, d, 1e-5)

	d = calcD(data2, data1)
	assert.InDelta(expectedD, d, 1e-5)

}

func BenchmarkCalcD(b *testing.B) {

	data := make([]float64, 100)
	for i := range data {
		data[i] = float64(i+1) / float64(len(data))
	}

	for i := 0; i < b.N; i++ {
		calcD(data, data)
	}
}

func BenchmarkIsMatch(b *testing.B) {

	data1 := make([]float64, 100)
	for i := range data1 {
		data1[i] = float64(-i) / float64(len(data1))
	}

	data2 := make([]float64, 100)
	for i := range data2 {
		data2[i] = float64(i+1) / float64(len(data2))
	}

	ks := KSTest{
		confidence: 4,
	}

	f1 := ks.Fingerprint(morgoth.Window{
		Data: data1,
	})

	f2 := ks.Fingerprint(morgoth.Window{
		Data: data2,
	})

	for i := 0; i < b.N; i++ {
		f1.IsMatch(f2)
	}
}

func BenchmarkFingerprint(b *testing.B) {

	data := make([]float64, 100)
	for i := range data {
		data[i] = float64(i) / float64(len(data))
	}

	ks := KSTest{
		confidence: 4,
	}

	w := morgoth.Window{
		Data: data,
	}

	for i := 0; i < b.N; i++ {
		ks.Fingerprint(w)
	}
}
