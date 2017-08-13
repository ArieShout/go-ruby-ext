package rb

type PBreak struct {
	Label string
}

func Label(label string, action func()) {
	if label == "" {
		panic("Empty label")
	}
	defer RecoverBreak(label)
	action()
}

func RecoverBreak(label string) {
	if err := recover(); err != nil {
		if brk, ok := err.(PBreak); ok {
			if brk.Label == label {
				return
			}
		}
		panic(err)
	}
}

func Break() {
	panic(PBreak{""})
}

func BreakLabel(label string) {
	if label == "" {
		panic("Empty label")
	}
	panic(PBreak{label})
}
