package task

import "github.com/kyokomi/emoji"

type CompletedState struct {
	Completed bool
}

func (c CompletedState) String() string {
	if c.Completed {
		return emoji.Sprint(":white_check_mark:")
	}
	return emoji.Sprint(":black_square_button:")
}
