package format

import (
	"fmt"
)

func TruncateText(text string, size int, formatted bool) string {

	if len(text) > size {
		text = fmt.Sprintf("%s...", text[:size])
	}

	if len(text) > size && formatted {
		return fmt.Sprintf("%s)", text) // add closing parenthesis
	}

	return text
}
