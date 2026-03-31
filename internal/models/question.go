package models

type Question struct {
	Kind             string          `json:"kind"`
	Version          string          `json:"version,omitempty"`
	Slug             string          `json:"slug"`
	Theme            Theme           `json:"theme"`
	Subthemes        []Theme         `json:"subthemes,omitempty"`
	Tags             []Theme         `json:"tags,omitempty"`
	Qtype            string          `json:"qtype"`
	Difficulty       string          `json:"difficulty"`
	EstimatedSeconds int             `json:"estimated_seconds"`
	ShuffleAnswers   bool            `json:"shuffle_answers"`
	I18n             map[string]I18n `json:"i18n"`
	Answers          []Answer        `json:"answers"`
	Sources          []string        `json:"sources,omitempty"`
}

type Theme struct {
	Slug string `json:"slug"`
}

type I18n struct {
	Title       string `json:"title"`
	Stem        string `json:"stem"`
	Explanation string `json:"explanation"`
}

type Answer struct {
	Slug      string           `json:"slug"`
	IsCorrect bool             `json:"is_correct"`
	I18n      map[string]Label `json:"i18n"`
}

type Label struct {
	Label string `json:"label"`
}

type ThemeTranslation struct {
	Slug string           `json:"slug"`
	I18n map[string]Label `json:"i18n,omitempty"`
}
