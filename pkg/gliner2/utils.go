package gliner2

import (
	"regexp"
	"strings"
)

// parseSection extracts content between <TAG> and </TAG>
func parseSection(text, tag string) string {
	re := regexp.MustCompile(regexp.QuoteMeta(`<` + tag + `>\s*([\s\S]*?)\s*</` + tag + `>`))
	match := re.FindStringSubmatch(text)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

// parseTSV parses simple TSV string into slice of records
func parseTSV(tsv string) [][]string {
	lines := strings.Split(tsv, "\n")
	var records [][]string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		// clean quotes?
		for i := range parts {
			parts[i] = strings.Trim(parts[i], "\"")
		}
		records = append(records, parts)
	}
	return records
}
