package morgoth

type Window struct {
	Data []float64
}

func (self *Window) Copy() *Window {
	data := make([]float64, len(self.Data))
	copy(data, self.Data)

	return &Window{
		Data: data,
	}
}
