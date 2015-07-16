package fn

/*
import (
	"fmt"
	"math/rand"
	"bufio"
//	"strconv"
)
*/
type Vector struct {
	A []float64 // data
	L int       // length
}

func NewVector(length int) (v *Vector) {
	v = new(Vector)
	v.L = length
	v.A = make([]float64, length)
	return v
}

type Matrix struct {
	R int
	C int
	A []float64
}

func NewMatrix(rows, cols int) (m *Matrix) {
	m = new(Matrix)
	m.R = rows
	m.C = cols
	m.A = make([]float64, rows*cols)
	return m
}

func (m Matrix) Get(i int, j int) float64 {
	return m.A[i*m.C+j]
}

func (m Matrix) Set(i int, j int, x float64) {
	//	m.A[i+j*m.C] = x
	m.A[i*m.C+j] = x
}

func (v Vector) Set(i int, x float64) {
	v.A[i] = x
}

func (v Vector) Get(i int) float64 {
	return v.A[i]
}

/*

func (p Vector) Swap(i int, j int) {
	x := p.A[i]
	p.A[i] = p.A[j]
	p.A[j] = x
}

func (v Vector) Len() int {
	return v.L
}

func (v Vector) Copy(w Vector) {
	for i := 0; i < len(v); i++ {
		v[i] = w[i]
	}
}

func (v Vector) Print() {
	for i := 0; i < len(v); i++ {
		fmt.Printf("%d ", v[i])
	}
	fmt.Print("\n")
}

func (m *Matrix) Print() {
	var i, j int
	for i = 0; i < m.R; i++ {
		for j = 0; j < m.C; j++ {
			fmt.Printf("%d ", m.Get(i, j))
		}
		fmt.Print("\n")
	}
}

func Perm(p Vector) {
	n := int(len(p))
	var i int
	for i = 0; i < n; i++ {
		p[i] = float64(i)
	}
	for i = 0; i < n; i++ {
		p.Swap(i, i+rand.Intn(n-i))
	}
}

func skip(rd *bufio.Reader) {
	var b byte = ' '
	var err error
	for b == ' ' || b == '\t' || b == '\n' {
		b, err = rd.ReadByte()
		if err != nil {
			return
		}
	}
	rd.UnreadByte()
}

func wskip(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] != ' ' && s[i] != '\t' {
			return s[i:]
		}
	}
	return ""
}

func end(s string) (i int) {
	for i = 0; i < int(len(s)); i++ {
		if s[i] == ' ' || s[i] == '\t' || s[i] == '\n'{
			return i
		}
	}
	return 0
}
func readUint(s string) (int, int){
	i := end(s)
	x, _ := strconv.ParseInt(s[:i], 10, 64)
	return int(x), i
}

func readMatrix(rd *bufio.Reader, n int) *Matrix {
	M := NewMatrix(n)
	var i, j int
	for i = 0; i < n; i++ {
		skip(rd)
		line, _ := rd.ReadString('\n')
		for j = 0; j < n; j++ {
			line = wskip(line)
			x, p := readUint(line)
			M.Set(j, i, x)
			if p == 0 {
				panic("bad int")
			}
			line = line[p:]
		}
	}
	return M
}
*/
