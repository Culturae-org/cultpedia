package checks

import (
	"fmt"
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
	if q.Theme.Slug == "" {
		return fmt.Errorf("theme.slug is required")
	}
	if len(q.Answers) != 4 {
		return fmt.Errorf("must have exactly 4 answers")
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

	validDifficulties := []string{"beginner", "intermediate", "expert"}
	difficultyValid := false
	for _, d := range validDifficulties {
		if q.Difficulty == d {
			difficultyValid = true
			break
		}
	}
	if !difficultyValid {
		errors = append(errors, fmt.Sprintf("✗ Invalid difficulty '%s' (allowed: beginner, intermediate, expert)", q.Difficulty))
	}

	if q.Points < 0.5 || q.Points > 5.0 {
		errors = append(errors, fmt.Sprintf("✗ Points must be between 0.5 and 5.0 (got %.1f)", q.Points))
	}

	if len(q.Sources) == 0 {
		errors = append(errors, "✗ At least one source URL is required")
	}

	validQtypes := []string{"single_choice", "multiple_choice"}
	qtypeValid := false
	for _, qt := range validQtypes {
		if q.Qtype == qt {
			qtypeValid = true
			break
		}
	}
	if !qtypeValid {
		errors = append(errors, fmt.Sprintf("✗ Invalid qtype '%s' (allowed: single_choice, multiple_choice)", q.Qtype))
	}

	if q.EstimatedSeconds < 5 || q.EstimatedSeconds > 300 {
		errors = append(errors, fmt.Sprintf("✗ Estimated seconds must be between 5 and 300 (got %d)", q.EstimatedSeconds))
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
