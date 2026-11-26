package components

type Input struct {
	RightKeyPressed bool
	LeftKeyPressed  bool
	DownKeyPressed  bool
	UpKeyPressed    bool
}

func NewInput() *Input {
	return &Input{}
}
