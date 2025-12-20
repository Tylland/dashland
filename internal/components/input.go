package components

type InputComponent struct {
	RightKeyPressed bool
	LeftKeyPressed  bool
	DownKeyPressed  bool
	UpKeyPressed    bool
}

func NewInputComponent() *InputComponent {
	return &InputComponent{}
}
