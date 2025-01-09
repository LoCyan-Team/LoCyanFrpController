package main

import (
	log2 "log"
	"strings"
)

func main1() {
	log := "2025-01-09T14:30:08+08:00   ERROR   ruleset match error     {\"name\": \"Block Socks\", \"id\": 1877241426203385856, \"src\": \"95.214.52.191:40366\", \"dst\": \"36.50.226.25:12000\", \"error\": \"interface conversion: interface {} is nil, not bool (1:27)\\n | socks != nil && socks.yes && (port.dst >= 10000) && (port.dst <= 65535)\\n | ..........................^\"}"
	log2.Printf(SplitLog(log)[2])
}

func SplitLog(log string) []string {
	var parts []string
	var currentPart strings.Builder
	var inQuotes bool  // Flag to indicate if we are inside quotes
	var spaceCount int // Count of consecutive spaces

	// Iterate over the string
	for i := 0; i < len(log); i++ {
		ch := log[i]

		// Check for quotes to toggle the inQuotes flag
		if ch == '"' {
			inQuotes = !inQuotes
			currentPart.WriteByte(ch)
			spaceCount = 0 // Reset space count after a quote
			continue
		}

		// If we are inside quotes, just append the character
		if inQuotes {
			currentPart.WriteByte(ch)
			continue
		}

		// If we encounter a space, increment the space count
		if ch == ' ' {
			spaceCount++
			// If this is the second consecutive space, it's a separator
			if spaceCount == 2 {
				parts = append(parts, currentPart.String())
				currentPart.Reset()
				spaceCount = 0 // Reset space count after adding a part
			}
			continue
		}

		// If the current character is not a space, append it to the current part
		if spaceCount > 0 {
			// If there was at least one space before this character, add a single space to the current part
			currentPart.WriteByte(' ')
			spaceCount = 0 // Reset space count
		}
		currentPart.WriteByte(ch)
	}

	// Add the last part if there is any content
	if currentPart.Len() > 0 {
		parts = append(parts, currentPart.String())
	}
	var rs []string
	for index, ch := range parts {
		if strings.ReplaceAll(string(ch), " ", "") != "" {
			rs = append(rs, parts[index])
		}
	}
	return rs
}
