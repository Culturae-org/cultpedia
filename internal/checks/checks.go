package checks

import (
	"fmt"
	"net/url"
	"strings"

	"cultpedia/internal/models"
	"cultpedia/internal/utils"
)

func ValidateQuestions() error {
	questions, err := utils.LoadQuestions()
	if err != nil {
		return err
	}
	slugs := make(map[string]bool)
	var errors []string

	for i, q := range questions {
		if err := validateQuestion(q); err != nil {
			errors = append(errors, fmt.Sprintf("line %d (slug: %s): %v", i+1, q.Slug, err))
			continue
		}
		if slugs[q.Slug] {
			errors = append(errors, fmt.Sprintf("duplicate detected for slug '%s' at line %d", q.Slug, i+1))
		} else {
			slugs[q.Slug] = true
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors:\n%s", strings.Join(errors, "\n"))
	}

	return nil
}

func validateQuestion(q models.Question) error {
	if q.Kind != "question" {
		return fmt.Errorf("kind must be 'question'")
	}
	if q.Slug == "" {
		return fmt.Errorf("slug is required")
	}
	if !isValidSlug(q.Slug) {
		return fmt.Errorf("slug must be lowercase with hyphens only (got '%s')", q.Slug)
	}
	if q.Theme.Slug == "" {
		return fmt.Errorf("theme.slug is required")
	}

	validQtypes := []string{"single_choice", "true_false"}
	if !contains(validQtypes, q.Qtype) {
		return fmt.Errorf("qtype must be one of: %s (got '%s')", strings.Join(validQtypes, ", "), q.Qtype)
	}

	validDifficulties := []string{"beginner", "intermediate", "advanced", "pro"}
	if !contains(validDifficulties, q.Difficulty) {
		return fmt.Errorf("difficulty must be one of: %s (got '%s')", strings.Join(validDifficulties, ", "), q.Difficulty)
	}

	if q.Points < 0.5 || q.Points > 5.0 {
		return fmt.Errorf("points must be between 0.5 and 5.0 (got %.1f)", q.Points)
	}

	if q.EstimatedSeconds < 5 || q.EstimatedSeconds > 30 {
		return fmt.Errorf("estimated_seconds must be between 5 and 30 (got %d)", q.EstimatedSeconds)
	}

	if len(q.Sources) == 0 {
		return fmt.Errorf("at least one source URL is required")
	}

	for i, source := range q.Sources {
		if err := validateURL(source); err != nil {
			return fmt.Errorf("invalid source URL #%d (%s): %v", i+1, source, err)
		}
	}

	if q.Qtype == "true_false" {
		if len(q.Answers) != 2 {
			return fmt.Errorf("true_false questions must have exactly 2 answers (got %d)", len(q.Answers))
		}
		hasTrue, hasFalse := false, false
		for _, a := range q.Answers {
			if a.Slug == "true" {
				hasTrue = true
			}
			if a.Slug == "false" {
				hasFalse = true
			}
		}
		if !hasTrue || !hasFalse {
			return fmt.Errorf("true_false questions must have answers with slugs 'true' and 'false'")
		}
	} else {
		if len(q.Answers) != 4 {
			return fmt.Errorf("must have exactly 4 answers (got %d)", len(q.Answers))
		}
	}

	correctCount := 0
	for _, a := range q.Answers {
		if a.IsCorrect {
			correctCount++
		}
		if a.Slug == "" {
			return fmt.Errorf("answer slug is required")
		}
	}
	if correctCount != 1 {
		return fmt.Errorf("must have exactly one correct answer")
	}
	requiredLangs := []string{"fr", "en", "es"}
	for _, lang := range requiredLangs {
		if _, ok := q.I18n[lang]; !ok {
			return fmt.Errorf("missing %s translation in question", lang)
		}
		for _, a := range q.Answers {
			if _, ok := a.I18n[lang]; !ok {
				return fmt.Errorf("missing %s translation in answer %s", lang, a.Slug)
			}
		}
	}
	return nil
}
func isValidSlug(slug string) bool {
	if slug == "" {
		return false
	}
	for _, c := range slug {
		isLower := c >= 'a' && c <= 'z'
		isDigit := c >= '0' && c <= '9'
		isHyphen := c == '-'
		if !isLower && !isDigit && !isHyphen {
			return false
		}
	}
	if slug[0] == '-' || slug[len(slug)-1] == '-' {
		return false
	}
	return true
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func validateURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format")
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("URL must use http or https scheme")
	}

	if parsed.Host == "" {
		return fmt.Errorf("URL must have a host")
	}

	return nil
}

func ValidateQuestionStrict(q models.Question) error {
	var errors []string

	if err := validateQuestion(q); err != nil {
		errors = append(errors, fmt.Sprintf("✗ %v", err))
	}

	requiredLangs := []string{"fr", "en", "es"}
	if len(q.I18n) != len(requiredLangs) {
		errors = append(errors, fmt.Sprintf("✗ Exactly 3 languages required (fr, en, es), got %d", len(q.I18n)))
	}
	for lang := range q.I18n {
		found := false
		for _, req := range requiredLangs {
			if lang == req {
				found = true
				break
			}
		}
		if !found {
			errors = append(errors, fmt.Sprintf("✗ Invalid language '%s' (only fr, en, es allowed)", lang))
		}
	}

	minStemLength := 10
	minExplanationLength := 20
	for lang, content := range q.I18n {
		if len(strings.TrimSpace(content.Stem)) < minStemLength {
			errors = append(errors, fmt.Sprintf("✗ %s stem too short (min %d chars, got %d)", lang, minStemLength, len(content.Stem)))
		}
		if len(strings.TrimSpace(content.Explanation)) < minExplanationLength {
			errors = append(errors, fmt.Sprintf("✗ %s explanation too short (min %d chars, got %d)", lang, minExplanationLength, len(content.Explanation)))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("%s", strings.Join(errors, "\n"))
	}

	return nil
}

func CheckDuplicates() string {
	questions, err := utils.LoadQuestions()
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	slugs := make(map[string]int)
	var duplicates []string
	for i, q := range questions {
		if firstLine, exists := slugs[q.Slug]; exists {
			duplicates = append(duplicates, fmt.Sprintf("slug '%s' duplicated: first occurrence line %d, occurrence line %d", q.Slug, firstLine+1, i+1))
		} else {
			slugs[q.Slug] = i
		}
	}
	if len(duplicates) > 0 {
		return fmt.Sprintf("duplicates detected:\n%s", strings.Join(duplicates, "\n"))
	} else {
		return "No duplicates."
	}
}

func CheckTranslations() string {
	questions, err := utils.LoadQuestions()
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	valid := true
	var missing []string
	for i, q := range questions {
		langs := []string{"fr", "en", "es"}
		for _, lang := range langs {
			if _, ok := q.I18n[lang]; !ok {
				valid = false
				missing = append(missing, fmt.Sprintf("question line %d (slug: %s): missing %s translation in title/question/explanation", i+1, q.Slug, lang))
			}
			for j, a := range q.Answers {
				if _, ok := a.I18n[lang]; !ok {
					valid = false
					missing = append(missing, fmt.Sprintf("answer %d of question line %d (slug: %s): missing %s translation", j+1, i+1, q.Slug, lang))
				}
			}
		}
	}
	if valid {
		return "All translations present."
	} else {
		return fmt.Sprintf("missing translations:\n%s", strings.Join(missing, "\n"))
	}
}
