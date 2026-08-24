package cloud

import (
	"strings"
	"unicode"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/domain"
)

var serviceSearchAliases = map[string]string{
	"auth": "authentication",
	"db":   "database",
}

var serviceSearchStopWords = map[string]struct{}{
	"a":       {},
	"an":      {},
	"cloud":   {},
	"service": {},
	"the":     {},
}

func cloudServiceMatchesQuery(service domain.CloudService, query string) bool {
	queryTerms := meaningfulSearchTerms(normalizeSearchTerms(query))
	candidateTerms := normalizeSearchTerms(service.ID + " " + service.Name)

	for _, queryTerm := range queryTerms {
		matched := false
		for _, candidateTerm := range candidateTerms {
			if searchTermsMatch(queryTerm, candidateTerm) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	return len(queryTerms) > 0
}

func normalizeSearchTerms(value string) []string {
	terms := strings.FieldsFunc(strings.ToLower(strings.TrimSpace(value)), func(character rune) bool {
		return !unicode.IsLetter(character) && !unicode.IsNumber(character)
	})

	for index, term := range terms {
		if alias, exists := serviceSearchAliases[term]; exists {
			terms[index] = alias
		}
	}

	return terms
}

func meaningfulSearchTerms(terms []string) []string {
	meaningful := make([]string, 0, len(terms))
	for _, term := range terms {
		if _, isStopWord := serviceSearchStopWords[term]; !isStopWord {
			meaningful = append(meaningful, term)
		}
	}

	if len(meaningful) == 0 {
		return terms
	}
	return meaningful
}

func searchTermsMatch(queryTerm, candidateTerm string) bool {
	if queryTerm == candidateTerm {
		return true
	}
	if len([]rune(queryTerm)) >= 4 && strings.Contains(candidateTerm, queryTerm) {
		return true
	}

	maximumDistance := allowedSearchDistance(len([]rune(queryTerm)))
	return maximumDistance > 0 && levenshteinDistance(queryTerm, candidateTerm) <= maximumDistance
}

func allowedSearchDistance(length int) int {
	switch {
	case length >= 11:
		return 3
	case length >= 7:
		return 2
	case length >= 4:
		return 1
	default:
		return 0
	}
}

func levenshteinDistance(left, right string) int {
	leftRunes := []rune(left)
	rightRunes := []rune(right)
	previous := make([]int, len(rightRunes)+1)
	current := make([]int, len(rightRunes)+1)

	for index := range previous {
		previous[index] = index
	}

	for leftIndex, leftRune := range leftRunes {
		current[0] = leftIndex + 1
		for rightIndex, rightRune := range rightRunes {
			cost := 0
			if leftRune != rightRune {
				cost = 1
			}
			current[rightIndex+1] = min(
				current[rightIndex]+1,
				previous[rightIndex+1]+1,
				previous[rightIndex]+cost,
			)
		}
		previous, current = current, previous
	}

	return previous[len(rightRunes)]
}
