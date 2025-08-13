package parser

import (
	"path/filepath"
	"regexp"
	"strings"
)

var (
	numericPattern   = regexp.MustCompile(`^(.+?)[-_.\s]*(\d{2,})$`)
	separatorPattern = regexp.MustCompile(`[-_.\s]+$`)
)

func ParseGroup(filename, prefix string) (string, bool) {
	stem := strings.TrimSuffix(filename, filepath.Ext(filename))

	if prefix != "" {
		if group, ok := tryPrefixForm(stem, prefix); ok {
			return group, true
		}
	}

	return tryNumericForm(stem)
}

func tryNumericForm(stem string) (string, bool) {
	matches := numericPattern.FindStringSubmatch(stem)
	if len(matches) != 3 {
		return "", false
	}

	group := matches[1]
	group = separatorPattern.ReplaceAllString(group, "")

	if group == "" {
		return "", false
	}

	return group, true
}

func tryPrefixForm(stem, prefix string) (string, bool) {
	escapedPrefix := regexp.QuoteMeta(prefix)

	pattern := regexp.MustCompile(`^(.+?)([-_.\s]?)` + escapedPrefix + `(\d{2,})`)
	matches := pattern.FindStringSubmatch(stem)
	if len(matches) != 4 {
		return "", false
	}

	groupPart := matches[1]
	separator := matches[2]

	if separator == "" {
		if len(groupPart) == 0 {
			return "", false
		}
		lastChar := groupPart[len(groupPart)-1]
		if lastChar != '-' && lastChar != '_' && lastChar != '.' && lastChar != ' ' {
			return "", false
		}
	}

	group := separatorPattern.ReplaceAllString(groupPart, "")
	if group == "" {
		return "", false
	}

	return group, true
}
