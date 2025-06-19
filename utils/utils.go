package utils

import (
	tsize "github.com/kopoli/go-terminal-size"
)

/*
Returns the input int if it's between min and max, otherwise
returns min or max
*/
func ClampInt(input int, minimum int, maximum int) int {
	return max(minimum, min(input, maximum))
}

/*
Returns the input int between the two boundaries and wraps
it if it's out of the given bounds
*/
func WrapInt(value int, min int, max int) int {
	diff := max - min
	if value < min {
		value += diff * ((min-value)/diff + 1)
	}
	return min + (value-min)%diff
}

func GetTerminalSize() (width int, height int, err error) {
	s, err := tsize.GetSize()
	if err != nil {
		return 0, 0, err
	}
	return s.Width, s.Height, nil
}
