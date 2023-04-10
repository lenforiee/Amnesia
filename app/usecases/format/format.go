package format

import (
	"fmt"
)

func TruncateText(text string, size int) string {
	if len(text) > size {
		return fmt.Sprintf("%s...", text[:size])
	}

	return text
}
