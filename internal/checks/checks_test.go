package checks

import (
	"testing"

	"cultpedia/internal/models"
)

const (
	qtypeSingleChoice = "single_choice"
	qtypeTrueFalse    = "true_false"
)

func TestIsValidSlug(t *testing.T) {
	tests := []struct {
		name     string
		slug     string
		expected bool
	}{
		{"valid simple slug", "history", true},
		{"valid slug with hyphens", "world-war-2", true},
		{"valid slug with numbers", "science-101", true},
		{"empty slug", "", false},
		{"uppercase letters", "History", false},
		{"spaces", "world war", false},
		{"underscores", "world_war", false},
		{"starts with hyphen", "-history", false},
		{"ends with hyphen", "history-", false},
		{"special characters", "history!", false},
		{"valid complex slug", "science-physics-einstein-relativity-1905", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidSlug(tt.slug)
			if result != tt.expected {
				t.Errorf("isValidSlug(%q) = %v, expected %v", tt.slug, result, tt.expected)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		item     string
		expected bool
	}{
		{"item exists", []string{"a", "b", "c"}, "b", true},
		{"item not exists", []string{"a", "b", "c"}, "d", false},
		{"empty slice", []string{}, "a", false},
		{"first item", []string{"a", "b", "c"}, "a", true},
		{"last item", []string{"a", "b", "c"}, "c", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.slice, tt.item)
			if result != tt.expected {
				t.Errorf("contains(%v, %q) = %v, expected %v", tt.slice, tt.item, result, tt.expected)
			}
		})
	}
}

func createValidQuestion() models.Question {
	return models.Question{
		Kind:             "question",
		Slug:             "test-question-slug",
		Theme:            models.Theme{Slug: "history"},
		Qtype:            qtypeSingleChoice,
		Difficulty:       "beginner",
		EstimatedSeconds: 20,
		Sources:          []string{"https://example.com"},
		I18n: map[string]models.I18n{
			"fr": {Title: "Titre", Stem: "Question en francais ?", Explanation: "Explication detaillee ici."},
			"en": {Title: "Title", Stem: "Question in English?", Explanation: "Detailed explanation here."},
			"es": {Title: "Titulo", Stem: "Pregunta en espanol?", Explanation: "Explicacion detallada aqui."},
		},
		Answers: []models.Answer{
			{Slug: "a", IsCorrect: true, I18n: map[string]models.Label{"fr": {Label: "A"}, "en": {Label: "A"}, "es": {Label: "A"}}},
			{Slug: "b", IsCorrect: false, I18n: map[string]models.Label{"fr": {Label: "B"}, "en": {Label: "B"}, "es": {Label: "B"}}},
			{Slug: "c", IsCorrect: false, I18n: map[string]models.Label{"fr": {Label: "C"}, "en": {Label: "C"}, "es": {Label: "C"}}},
			{Slug: "d", IsCorrect: false, I18n: map[string]models.Label{"fr": {Label: "D"}, "en": {Label: "D"}, "es": {Label: "D"}}},
		},
	}
}

func TestValidateQuestion(t *testing.T) {
	t.Run("valid question", func(t *testing.T) {
		q := createValidQuestion()
		err := validateQuestion(q)
		if err != nil {
			t.Errorf("validateQuestion() returned unexpected error: %v", err)
		}
	})

	t.Run("invalid kind", func(t *testing.T) {
		q := createValidQuestion()
		q.Kind = "invalid"
		err := validateQuestion(q)
		if err == nil {
			t.Error("validateQuestion() should return error for invalid kind")
		}
	})

	t.Run("empty slug", func(t *testing.T) {
		q := createValidQuestion()
		q.Slug = ""
		err := validateQuestion(q)
		if err == nil {
			t.Error("validateQuestion() should return error for empty slug")
		}
	})

	t.Run("invalid slug format", func(t *testing.T) {
		q := createValidQuestion()
		q.Slug = "Invalid_Slug"
		err := validateQuestion(q)
		if err == nil {
			t.Error("validateQuestion() should return error for invalid slug format")
		}
	})

	t.Run("invalid qtype", func(t *testing.T) {
		q := createValidQuestion()
		q.Qtype = "invalid_type"
		err := validateQuestion(q)
		if err == nil {
			t.Error("validateQuestion() should return error for invalid qtype")
		}
	})

	t.Run("invalid difficulty", func(t *testing.T) {
		q := createValidQuestion()
		q.Difficulty = "expert"
		err := validateQuestion(q)
		if err == nil {
			t.Error("validateQuestion() should return error for invalid difficulty")
		}
	})

	t.Run("estimated seconds too low", func(t *testing.T) {
		q := createValidQuestion()
		q.EstimatedSeconds = 2
		err := validateQuestion(q)
		if err == nil {
			t.Error("validateQuestion() should return error for estimated_seconds below 5")
		}
	})

	t.Run("estimated seconds too high", func(t *testing.T) {
		q := createValidQuestion()
		q.EstimatedSeconds = 31
		err := validateQuestion(q)
		if err == nil {
			t.Error("validateQuestion() should return error for estimated_seconds above 30")
		}
	})

	t.Run("no sources", func(t *testing.T) {
		q := createValidQuestion()
		q.Sources = []string{}
		err := validateQuestion(q)
		if err == nil {
			t.Error("validateQuestion() should return error when no sources provided")
		}
	})

	t.Run("wrong answer count for single_choice", func(t *testing.T) {
		q := createValidQuestion()
		q.Answers = q.Answers[:2]
		err := validateQuestion(q)
		if err == nil {
			t.Error("validateQuestion() should return error for single_choice with != 4 answers")
		}
	})

	t.Run("no correct answer", func(t *testing.T) {
		q := createValidQuestion()
		for i := range q.Answers {
			q.Answers[i].IsCorrect = false
		}
		err := validateQuestion(q)
		if err == nil {
			t.Error("validateQuestion() should return error when no correct answer")
		}
	})

	t.Run("multiple correct answers", func(t *testing.T) {
		q := createValidQuestion()
		q.Answers[0].IsCorrect = true
		q.Answers[1].IsCorrect = true
		err := validateQuestion(q)
		if err == nil {
			t.Error("validateQuestion() should return error for multiple correct answers")
		}
	})

	t.Run("missing translation", func(t *testing.T) {
		q := createValidQuestion()
		delete(q.I18n, "es")
		err := validateQuestion(q)
		if err == nil {
			t.Error("validateQuestion() should return error for missing translation")
		}
	})
}

func TestValidateTrueFalseQuestion(t *testing.T) {
	t.Run("valid true_false question", func(t *testing.T) {
		q := createValidQuestion()
		q.Qtype = qtypeTrueFalse
		q.Answers = []models.Answer{
			{Slug: "true", IsCorrect: true, I18n: map[string]models.Label{"fr": {Label: "Vrai"}, "en": {Label: "True"}, "es": {Label: "Verdadero"}}},
			{Slug: "false", IsCorrect: false, I18n: map[string]models.Label{"fr": {Label: "Faux"}, "en": {Label: "False"}, "es": {Label: "Falso"}}},
		}
		err := validateQuestion(q)
		if err != nil {
			t.Errorf("validateQuestion() returned unexpected error for true_false: %v", err)
		}
	})

	t.Run("true_false with wrong slugs", func(t *testing.T) {
		q := createValidQuestion()
		q.Qtype = qtypeTrueFalse
		q.Answers = []models.Answer{
			{Slug: "yes", IsCorrect: true, I18n: map[string]models.Label{"fr": {Label: "Oui"}, "en": {Label: "Yes"}, "es": {Label: "Si"}}},
			{Slug: "no", IsCorrect: false, I18n: map[string]models.Label{"fr": {Label: "Non"}, "en": {Label: "No"}, "es": {Label: "No"}}},
		}
		err := validateQuestion(q)
		if err == nil {
			t.Error("validateQuestion() should return error for true_false with wrong answer slugs")
		}
	})

	t.Run("true_false with wrong answer count", func(t *testing.T) {
		q := createValidQuestion()
		q.Qtype = qtypeTrueFalse
		err := validateQuestion(q)
		if err == nil {
			t.Error("validateQuestion() should return error for true_false with != 2 answers")
		}
	})
}
