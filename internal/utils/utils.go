package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"cultpedia/internal/models"
)

const (
	ManifestFile  = "datasets/general-knowledge/manifest.json"
	QuestionsFile = "datasets/general-knowledge/questions.ndjson"
	ThemesFile    = "datasets/general-knowledge/themes.ndjson"
	SubthemesFile = "datasets/general-knowledge/subthemes.ndjson"
	TagsFile      = "datasets/general-knowledge/tags.ndjson"
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
	f, err := os.OpenFile(QuestionsFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer func() {
		_ = f.Close()
	}()
	_, err = f.Seek(0, 2)
	if err != nil {
		return err
	}
	stat, err := f.Stat()
	if err != nil {
		return err
	}
	if stat.Size() > 0 {
		lastByte := make([]byte, 1)
		_, err = f.ReadAt(lastByte, stat.Size()-1)
		if err != nil {
			return err
		}
		if lastByte[0] != '\n' {
			_, err = f.WriteString("\n")
			if err != nil {
				return err
			}
		}
	}
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

func PrintHelp() {
	helpText := `
Cultpedia - Question Dataset Management Tool

USAGE FOR CONTRIBUTORS:
  ./cultpedia                  Launch interactive UI (recommended for adding questions)

USAGE FOR MAINTAINERS:
  ./cultpedia [command]

COMMANDS:
  help                  Show this help message
  validate              Validate the questions dataset for consistency and correctness
  check-duplicates      Check for duplicate questions in the dataset
  check-translations    Check for missing translations in the dataset
  add                   Add a new question to the dataset via interactive prompts
  sync-themes           Synchronize themes and subthemes with the questions dataset
  bump-version          Increment version and update manifest (automated in CI)

CONTRIBUTION GUIDE:
  For questions: Fork → Edit questions.ndjson → Create PR
  For code: Fork → Edit code → Run tests → Create PR
  See docs/CONTRIBUTING.md for detailed instructions
`
	fmt.Println(helpText)
}
