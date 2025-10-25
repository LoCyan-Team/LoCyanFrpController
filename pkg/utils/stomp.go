package utils

import "strings"

func ExtractHeader(frame, key string) string {
	lines := strings.Split(frame, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, key+":") {
			return strings.TrimSpace(strings.TrimPrefix(line, key+":"))
		}
	}
	return ""
}

func ExtractBody(frame string) string {
	parts := strings.SplitN(frame, "\n\n", 2)
	if len(parts) > 1 {
		return strings.TrimSuffix(parts[1], "\x00")
	}
	return ""
}
