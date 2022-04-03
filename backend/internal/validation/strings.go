package validation

import "fmt"

func StringMinLen(value string, min int) error {
	l := len([]rune(value))
	if l >= min {
		return nil
	}

	return fmt.Errorf("expected mininal length: %d, received: %d", min, l)
}

func StringMaxLen(value string, max int) error {
	l := len([]rune(value))
	if l < max {
		return nil
	}

	return fmt.Errorf("expected mininal length: %d, received: %d", max, l)
}
