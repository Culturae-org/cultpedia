package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"cultpedia/internal/models"
	"cultpedia/internal/utils"

	perplexity "github.com/sgaunet/perplexity-go/v2"
)

type GenerateConfig struct {
	Topic      string
	Theme      string
	Difficulty string
	Qtype      string
	Count      int
}

func GenerateQuestions(config GenerateConfig) ([]models.Question, error) {
	apiKey := os.Getenv("PPLX_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("PPLX_API_KEY environment variable is not set\n\nTo set it:\n  export PPLX_API_KEY=\"your-api-key\"\n\nGet your API key at: https://www.perplexity.ai/settings/api")
	}

	client := perplexity.NewClient(apiKey)

	prompt := buildPrompt(config)

	messages := []perplexity.Message{
		{
			Role:    "system",
			Content: getSystemPrompt(),
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	req := perplexity.NewCompletionRequest(
		perplexity.WithMessages(messages),
	)

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("request validation error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	res, err := client.SendCompletionRequestWithContext(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("API error: %v", err)
	}

	content := res.GetLastContent()
	questions, err := parseAIResponse(content)
	if err != nil {
		return nil, fmt.Errorf("parsing error: %v\n\nRaw response:\n%s", err, content)
	}

	return questions, nil
}

func getSystemPrompt() string {
	llmsContent, err := os.ReadFile(utils.LLMSGuideFile)
	if err != nil {
		return `You are an educational quiz question generator for Cultpedia.
Each question must:
- Be factual and verifiable
- Have Wikipedia or academic sources
- Be translated into EN, FR, and ES
- Follow the exact Cultpedia JSON format`
	}
	return string(llmsContent)
}

func buildPrompt(config GenerateConfig) string {
	qtype := config.Qtype
	if qtype == "" {
		qtype = "single_choice"
	}
	difficulty := config.Difficulty
	if difficulty == "" {
		difficulty = "intermediate"
	}
	theme := config.Theme
	if theme == "" {
		theme = "science"
	}
	count := config.Count
	if count <= 0 || count > 10 {
		count = 3
	}

	var answersInfo string
	if qtype == "true_false" {
		answersInfo = `- Type: true_false (exactly 2 answers with slugs "true" and "false")`
	} else {
		answersInfo = `- Type: single_choice (exactly 4 answers, only one correct)`
	}

	prompt := fmt.Sprintf(`Generate %d quiz questions about the following topic: "%s"

Configuration:
- Theme: %s
- Difficulty: %s
%s

IMPORTANT RULES:
1. Each question MUST have verifiable sources (Wikipedia, official websites)
2. Information MUST be factually correct
3. Translate EVERYTHING into French (fr), English (en), and Spanish (es)
4. The slug must follow the format: {theme}-{subtheme}-{element}-{detail}
5. Wrong answers must be plausible but clearly incorrect

Return ONLY a JSON array with the questions, no additional text.
Expected format: [{"kind": "question", ...}, ...]

Example structure for a question:
{
  "kind": "question",
  "version": "1.0",
  "slug": "theme-subtheme-element-detail",
  "theme": { "slug": "theme" },
  "subthemes": [{ "slug": "subtheme" }],
  "tags": [{ "slug": "tag1" }],
  "qtype": "%s",
  "difficulty": "%s",
  "estimated_seconds": 15,
  "points": 1.0,
  "shuffle_answers": true,
  "i18n": {
    "fr": { "title": "...", "stem": "...?", "explanation": "..." },
    "en": { "title": "...", "stem": "...?", "explanation": "..." },
    "es": { "title": "...", "stem": "...?", "explanation": "..." }
  },
  "answers": [...],
  "sources": ["https://..."]
}`, count, config.Topic, theme, difficulty, answersInfo, qtype, difficulty)

	return prompt
}

func parseAIResponse(content string) ([]models.Question, error) {
	rawResponseFile := "datasets/pending_raw_response.json"
	_ = os.WriteFile(rawResponseFile, []byte(content), 0644)

	content = strings.TrimSpace(content)

	re := regexp.MustCompile("```(?:json)?\\s*([\\s\\S]*?)\\s*```")
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		content = strings.TrimSpace(matches[1])
	}

	var questions []models.Question
	if err := json.Unmarshal([]byte(content), &questions); err != nil {
		startIdx := strings.Index(content, "[")
		if startIdx != -1 {
			arrayContent := content[startIdx:]
			
			if err := json.Unmarshal([]byte(arrayContent), &questions); err != nil {
				questions = extractValidQuestions(arrayContent)
				if len(questions) == 0 {
					return nil, fmt.Errorf("failed to parse JSON array: %v\n\nRaw response saved to: %s", err, rawResponseFile)
				}
				fmt.Printf("Response was truncated. Recovered %d valid question(s).\n", len(questions))
				fmt.Printf("   Raw response saved to: %s\n\n", rawResponseFile)
			}
		} else {
			var singleQuestion models.Question
			if err := json.Unmarshal([]byte(content), &singleQuestion); err != nil {
				return nil, fmt.Errorf("failed to parse response as JSON: %v\n\nRaw response saved to: %s", err, rawResponseFile)
			}
			questions = append(questions, singleQuestion)
		}
	}

	for i := range questions {
		if questions[i].Kind == "" {
			questions[i].Kind = "question"
		}
		if questions[i].Version == "" {
			questions[i].Version = "1.0"
		}
	}

	return questions, nil
}

func extractValidQuestions(content string) []models.Question {
	var questions []models.Question
	
	depth := 0
	start := -1
	inString := false
	escaped := false
	
	for i, char := range content {
		if escaped {
			escaped = false
			continue
		}
		
		if char == '\\' && inString {
			escaped = true
			continue
		}
		
		if char == '"' && !escaped {
			inString = !inString
			continue
		}
		
		if inString {
			continue
		}
		
		switch char {
		case '{':
			if depth == 0 {
				start = i
			}
			depth++
		case '}':
			depth--
			if depth == 0 && start != -1 {
				objStr := content[start : i+1]
				var q models.Question
				if err := json.Unmarshal([]byte(objStr), &q); err == nil {
					if q.Slug != "" && len(q.Answers) > 0 {
						questions = append(questions, q)
					}
				}
				start = -1
			}
		}
	}
	
	return questions
}

func SavePendingQuestion(q models.Question) error {
	minified, err := json.Marshal(q)
	if err != nil {
		return fmt.Errorf("minification error: %v", err)
	}
	ndjsonLine := string(minified) + "\n"

	dir := "datasets"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}

	if _, err := os.Stat(utils.PendingQuestionsFile); err == nil {
		f, err := os.Open(utils.PendingQuestionsFile)
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

	f, err := os.OpenFile(utils.PendingQuestionsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer func() { _ = f.Close() }()

	if _, err := f.WriteString(ndjsonLine); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	return nil
}

func SavePendingQuestions(questions []models.Question) error {
	for _, q := range questions {
		if err := SavePendingQuestion(q); err != nil {
			return err
		}
	}
	return nil
}

func LoadPendingQuestions() ([]models.Question, error) {
	if _, err := os.Stat(utils.PendingQuestionsFile); os.IsNotExist(err) {
		return []models.Question{}, nil
	}

	data, err := os.ReadFile(utils.PendingQuestionsFile)
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
			return nil, fmt.Errorf("json parsing error: %v", err)
		}
		questions = append(questions, q)
	}
	return questions, nil
}

func ApprovePendingQuestion(index int) (string, error) {
	pending, err := LoadPendingQuestions()
	if err != nil {
		return "", fmt.Errorf("error loading pending questions: %v", err)
	}

	if index < 0 || index >= len(pending) {
		return "", fmt.Errorf("invalid index: %d (have %d pending questions)", index, len(pending))
	}

	question := pending[index]

	if utils.SlugExists(question.Slug) {
		return "", fmt.Errorf("slug '%s' already exists in the dataset", question.Slug)
	}

	if err := utils.SaveQuestion(question); err != nil {
		return "", fmt.Errorf("error saving question: %v", err)
	}

	pending = append(pending[:index], pending[index+1:]...)
	if err := rewritePendingQuestions(pending); err != nil {
		return "", fmt.Errorf("error updating pending file: %v", err)
	}

	return fmt.Sprintf("✔ Question '%s' approved and added to dataset", question.Slug), nil
}

func RejectPendingQuestion(index int) (string, error) {
	pending, err := LoadPendingQuestions()
	if err != nil {
		return "", fmt.Errorf("error loading pending questions: %v", err)
	}

	if index < 0 || index >= len(pending) {
		return "", fmt.Errorf("invalid index: %d (have %d pending questions)", index, len(pending))
	}

	slug := pending[index].Slug

	pending = append(pending[:index], pending[index+1:]...)
	if err := rewritePendingQuestions(pending); err != nil {
		return "", fmt.Errorf("error updating pending file: %v", err)
	}

	return fmt.Sprintf("✔ Question '%s' rejected and removed", slug), nil
}

func ClearPendingQuestions() error {
	return os.Remove(utils.PendingQuestionsFile)
}

func rewritePendingQuestions(questions []models.Question) error {
	if len(questions) == 0 {
		if _, err := os.Stat(utils.PendingQuestionsFile); err == nil {
			return os.Remove(utils.PendingQuestionsFile)
		}
		return nil
	}

	var lines []string
	for _, q := range questions {
		minified, err := json.Marshal(q)
		if err != nil {
			return fmt.Errorf("minification error: %v", err)
		}
		lines = append(lines, string(minified))
	}

	content := strings.Join(lines, "\n") + "\n"
	return os.WriteFile(utils.PendingQuestionsFile, []byte(content), 0644)
}

func FormatPendingQuestionPreview(q models.Question, index int) string {
	lang := "fr"
	i18n, ok := q.I18n[lang]
	if !ok {
		if i18n, ok = q.I18n["en"]; !ok {
			return fmt.Sprintf("[%d] %s (no translation available)", index, q.Slug)
		}
	}

	correctAnswer := ""
	for _, a := range q.Answers {
		if a.IsCorrect {
			if label, ok := a.I18n[lang]; ok {
				correctAnswer = label.Label
			}
			break
		}
	}

	sources := "none"
	if len(q.Sources) > 0 {
		sources = fmt.Sprintf("%d source(s)", len(q.Sources))
	}

	return fmt.Sprintf(`[%d] %s
    Theme: %s | Difficulty: %s | Type: %s
    %s
    %s
    %s
`, index, q.Slug, q.Theme.Slug, q.Difficulty, q.Qtype, i18n.Stem, correctAnswer, sources)
}
