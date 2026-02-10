package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"cultpedia/internal/models"
)

const (
	ManifestFile             = "datasets/general-knowledge/manifest.json"
	QuestionsFile            = "datasets/general-knowledge/questions.ndjson"
	ThemesFile               = "datasets/general-knowledge/themes.ndjson"
	SubthemesFile            = "datasets/general-knowledge/subthemes.ndjson"
	TagsFile                 = "datasets/general-knowledge/tags.ndjson"
	NewQuestionFile          = "datasets/new-question.json"
	NewQuestionTrueFalseFile = "datasets/new-question-true-false.json"

	GeographyManifestFile  = "datasets/geography/manifest.json"
	CountriesFile          = "datasets/geography/countries.ndjson"
	ContinentsFile         = "datasets/geography/continents.ndjson"
	RegionsFile            = "datasets/geography/regions.ndjson"
	FlagsSVGDir            = "datasets/geography/assets/flags/svg"

	PendingQuestionsFile = "datasets/pending-questions.ndjson"
	LLMSGuideFile        = "LLMS.md"
)

func LoadQuestions() ([]models.Question, error) {
	data, err := os.ReadFile(QuestionsFile)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	var questions []models.Question
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var q models.Question
		if err := json.Unmarshal([]byte(line), &q); err != nil {
			return nil, fmt.Errorf("json parsing error at line %d: %v", len(questions)+1, err)
		}
		questions = append(questions, q)
	}
	return questions, nil
}

func SaveQuestion(q models.Question) error {
	minified, err := json.Marshal(q)
	if err != nil {
		return fmt.Errorf("minification error: %v", err)
	}
	ndjsonLine := string(minified) + "\n"

	if _, err := os.Stat(QuestionsFile); err == nil {
		f, err := os.Open(QuestionsFile)
		if err != nil {
			return fmt.Errorf("error opening file for check: %v", err)
		}
		defer func() { _ = f.Close() }()
		stat, err := f.Stat()
		if err != nil {
			return fmt.Errorf("error getting file stat: %v", err)
		}
		size := stat.Size()
		if size > 0 {
			buf := make([]byte, 1)
			_, err := f.ReadAt(buf, size-1)
			if err != nil {
				return fmt.Errorf("error reading file end: %v", err)
			}
			if buf[0] != '\n' {
				ndjsonLine = "\n" + ndjsonLine
			}
		}
	}

	f, err := os.OpenFile(QuestionsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer func() {
		_ = f.Close()
	}()
	if _, err := f.WriteString(ndjsonLine); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	return nil
}

func SlugExists(slug string) bool {
	questions, err := LoadQuestions()
	if err != nil {
		return false
	}
	for _, q := range questions {
		if q.Slug == slug {
			return true
		}
	}
	return false
}

func DetectModifiedTemplateFile() (filePath string, questionType string) {
	if isTemplateModified(NewQuestionFile, "default-question-slug") {
		return NewQuestionFile, "single_choice"
	}
	if isTemplateModified(NewQuestionTrueFalseFile, "default-true-false-question-slug") {
		return NewQuestionTrueFalseFile, "true_false"
	}
	return "", ""
}

func isTemplateModified(filePath, defaultSlug string) bool {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}
	var q models.Question
	if err := json.Unmarshal(data, &q); err != nil {
		return false
	}
	return q.Slug != defaultSlug
}

func PrintHelp() {
	helpText := `
Cultpedia - Question Dataset Management Tool

USAGE:
  cultpedia [command] [options]

COMMANDS:
  help                          Show this help message
  version                       Show local and remote dataset versions

  Questions Dataset:
  validate                      Validate the questions dataset for consistency and correctness
  check-duplicates              Check for duplicate questions in the dataset
  check-translations            Check for missing translations in the dataset
  preview [--type <qtype>]      Preview the current question template
  add [--type <qtype>]          Validate and add a new question to the dataset
  sync-themes                   Synchronize themes and subthemes with the questions dataset
  bump-version                  Increment version and update manifest (automated in CI)

  AI Generation (requires PPLX_API_KEY):
  generate <topic>              Generate questions using Perplexity AI
    --theme <theme>             Theme slug (default: science)
    --count <n>                 Number of questions (default: 3, max: 10)
    --difficulty <level>        beginner, intermediate, advanced
    --type <qtype>              single_choice or true_false
  pending                       List pending questions awaiting review
  approve <index>               Approve a pending question and add to dataset
  reject <index>                Reject and remove a pending question

  Geography Dataset:
  validate-geography            Validate the geography dataset (countries, continents, regions)
  check-geography-duplicates    Check for duplicate entries in geography dataset
  check-geography-translations  Check for missing translations in geography dataset
  bump-geography-version        Increment geography version and update checksums (automated in CI)

  General:
  init [dataset-name]           Initialize a new Cultpedia dataset structure
  api [port]                    Start the REST API server (default: 8080)

CONTRIBUTION GUIDE:
  For questions: Fork → Edit template → cultpedia preview → cultpedia add → Create PR
  For code: Fork → Edit code → Run tests → Create PR
  For AI generation: Set PPLX_API_KEY → generate → pending → approve/reject → Create PR

For more info, see CONTRIBUTING.md in the docs/ folder.
Or visit:
  https://docs.culturae.me/cultpedia/

Thank you for contributing to Cultpedia!
`
	fmt.Println(helpText)
}

func LoadThemeTranslations(filePath string) (map[string]models.ThemeTranslation, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]models.ThemeTranslation), nil
		}
		return nil, err
	}
	result := make(map[string]models.ThemeTranslation)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var t models.ThemeTranslation
		if err := json.Unmarshal([]byte(line), &t); err != nil {
			return nil, fmt.Errorf("json parsing error: %v", err)
		}
		result[t.Slug] = t
	}
	return result, nil
}

func LoadCountries() ([]models.Country, error) {
	data, err := os.ReadFile(CountriesFile)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	var countries []models.Country
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var c models.Country
		if err := json.Unmarshal([]byte(line), &c); err != nil {
			return nil, fmt.Errorf("json parsing error at line %d: %v", len(countries)+1, err)
		}
		countries = append(countries, c)
	}
	return countries, nil
}

func LoadContinents() ([]models.Continent, error) {
	data, err := os.ReadFile(ContinentsFile)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	var continents []models.Continent
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var c models.Continent
		if err := json.Unmarshal([]byte(line), &c); err != nil {
			return nil, fmt.Errorf("json parsing error at line %d: %v", len(continents)+1, err)
		}
		continents = append(continents, c)
	}
	return continents, nil
}

func LoadRegions() ([]models.Region, error) {
	data, err := os.ReadFile(RegionsFile)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	var regions []models.Region
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var r models.Region
		if err := json.Unmarshal([]byte(line), &r); err != nil {
			return nil, fmt.Errorf("json parsing error at line %d: %v", len(regions)+1, err)
		}
		regions = append(regions, r)
	}
	return regions, nil
}
