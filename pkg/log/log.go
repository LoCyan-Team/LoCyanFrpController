package log

import "strings"

func SplitLog(log string) []string {
	var parts []string
	var currentPart strings.Builder
	var inQuotes bool
	var spaceCount int

	for i := 0; i < len(log); i++ {
		ch := log[i]

		if ch == '"' {
			inQuotes = !inQuotes
			currentPart.WriteByte(ch)
			spaceCount = 0
			continue
		}

		if inQuotes {
			currentPart.WriteByte(ch)
			continue
		}

		if ch == ' ' {
			spaceCount++
			if spaceCount == 2 {
				parts = append(parts, currentPart.String())
				currentPart.Reset()
				spaceCount = 0
			}
			continue
		}

		if spaceCount > 0 {
			currentPart.WriteByte(' ')
			spaceCount = 0
		}
		currentPart.WriteByte(ch)
	}

	if currentPart.Len() > 0 {
		parts = append(parts, currentPart.String())
	}
	var rs []string
	for index, ch := range parts {
		if strings.ReplaceAll(ch, " ", "") != "" {
			rs = append(rs, parts[index])
		}
	}
	return rs
}
