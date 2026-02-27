package constant

import "regexp"

const (
	SourceIndeed = "Indeed"
	SourceITViec = "ITViec"
	SourceTopDev = "TopDev"
)

var KeywordRegex = regexp.MustCompile(`(?i)\b(backend|back-end|back\s+end|golang|go\s+developer|go\s+backend|\bGo\b|java\b|rust\b|software\s+engineer|developer)\b`)

const MaxSeenJobs = 5000
