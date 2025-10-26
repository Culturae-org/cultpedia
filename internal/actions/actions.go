package actions

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"cultpedia/internal/checks"
	"cultpedia/internal/models"
	"cultpedia/internal/utils"
)

func ValidateNewQuestion() (models.Question, error) {
	jsonFilePath := "datasets/new-question.json"
	if _, err := os.Stat(jsonFilePath); os.IsNotExist(err) {
		return models.Question{}, fmt.Errorf("file new-question.json not found")
	}
	data, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return models.Question{}, fmt.Errorf("read error: %v", err)
	}

	var question models.Question
	if err := json.Unmarshal(data, &question); err != nil {
		return models.Question{}, fmt.Errorf("json parsing error: %v", err)
	}

	if question.Kind != "question" {
		return models.Question{}, fmt.Errorf("kind must be 'question'")
	}
	if question.Slug == "" {
		return models.Question{}, fmt.Errorf("slug is required")
	}
	if question.Theme.Slug == "" {
		return models.Question{}, fmt.Errorf("theme.slug is required")
	}
	if len(question.Answers) != 4 {
		return models.Question{}, fmt.Errorf("must have exactly 4 answers")
	}
	correctCount := 0
	for _, a := range question.Answers {
		if a.IsCorrect {
			correctCount++
		}
	}
	if correctCount != 1 {
		return models.Question{}, fmt.Errorf("must have exactly one correct answer")
	}

	if question.Slug == "default-question-slug" {
		return models.Question{}, fmt.Errorf("slug cannot be the default template value 'default-question-slug' \nedit datasets/new-question.json to set a unique slug")
	}
	if question.Theme.Slug == "default-theme" {
		return models.Question{}, fmt.Errorf("theme slug cannot be the default template value 'default-theme' \nedit datasets/new-question.json to set a unique theme slug")
	}
	for _, sub := range question.Subthemes {
		if sub.Slug == "default-subtheme1" || sub.Slug == "default-subtheme2" {
			return models.Question{}, fmt.Errorf("subtheme slug cannot be the default template value '%s' \nedit datasets/new-question.json to set a unique subtheme slug", sub.Slug)
		}
	}
	for _, tag := range question.Tags {
		if tag.Slug == "default-tag1" || tag.Slug == "default-tag2" {
			return models.Question{}, fmt.Errorf("tag slug cannot be the default template value '%s' \nedit datasets/new-question.json to set a unique tag slug", tag.Slug)
		}
	}

	existingQuestions, err := utils.LoadQuestions()
	if err != nil {
		return models.Question{}, fmt.Errorf("error reading questions file: %v", err)
	}
	for i, q := range existingQuestions {
		if q.Slug == question.Slug {
			return models.Question{}, fmt.Errorf("✗ duplicate slug detected\n\n  Slug '%s' already exists at line %d\n  Please use a unique slug", question.Slug, i+1)
		}
	}

	langs := []string{"fr", "en", "es"}
	for _, lang := range langs {
		if _, ok := question.I18n[lang]; !ok {
			return models.Question{}, fmt.Errorf("missing %s translation in question (%s)", lang, question.Slug)
		}
		for j, a := range question.Answers {
			if _, ok := a.I18n[lang]; !ok {
				return models.Question{}, fmt.Errorf("missing %s translation in answer %d of question %s", lang, j+1, question.Slug)
			}
		}
	}

	if err := checks.ValidateQuestionStrict(question); err != nil {
		return models.Question{}, fmt.Errorf("validation failed:\n%v", err)
	}

	return question, nil
}

func AddValidatedQuestion(question models.Question) string {
	existingQuestions, err := utils.LoadQuestions()
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	lineNumber := len(existingQuestions) + 1

	if err := utils.SaveQuestion(question); err != nil {
		return fmt.Sprintf("error: %v", err)
	}

	message := fmt.Sprintf("✔ Question '%s' added successfully at line %d\n\n", question.Slug, lineNumber)
	message += "Next steps:\n"
	message += "  1. git add datasets/general-knowledge/questions.ndjson\n"
	message += "  2. git commit -m \"feat: add " + question.Slug + "\"\n"
	message += "  3. git push \n"
	message += "  4. Create a Pull Request in Github\n\n"
	message += "Our CI will automatically:\n"
	message += "  • Validate your question\n"
	message += "  • Sync themes and tags\n"
	message += "  • Bump the version"

	return message
}

func SyncThemes() string {
	questions, err := utils.LoadQuestions()
	if err != nil {
		return fmt.Sprintf("✗ error reading questions: %v", err)
	}
	themeSlugs := make(map[string]bool)
	subthemeSlugs := make(map[string]bool)
	tagSlugs := make(map[string]bool)

	for _, q := range questions {
		themeSlugs[q.Theme.Slug] = true
		for _, sub := range q.Subthemes {
			subthemeSlugs[sub.Slug] = true
		}
		for _, tag := range q.Tags {
			tagSlugs[tag.Slug] = true
		}
	}

	if err := writeSlugFile(utils.ThemesFile, themeSlugs); err != nil {
		return fmt.Sprintf("✗ error writing themes: %v", err)
	}
	if err := writeSlugFile(utils.SubthemesFile, subthemeSlugs); err != nil {
		return fmt.Sprintf("✗ error writing subthemes: %v", err)
	}
	if err := writeSlugFile(utils.TagsFile, tagSlugs); err != nil {
		return fmt.Sprintf("✗ error writing tags: %v", err)
	}

	if err := updateManifest(len(questions), len(themeSlugs), len(subthemeSlugs), len(tagSlugs)); err != nil {
		return fmt.Sprintf("✗ error updating manifest: %v", err)
	}
	return fmt.Sprintf("✔ Themes synced successfully\n  - %d questions\n  - %d themes\n  - %d subthemes\n  - %d tags", len(questions), len(themeSlugs), len(subthemeSlugs), len(tagSlugs))
}

func BumpVersion() (string, error) {
	data, err := os.ReadFile(utils.ManifestFile)
	if err != nil {
		return "", fmt.Errorf("error reading manifest: %v", err)
	}

	var manifest models.Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return "", fmt.Errorf("error parsing manifest: %v", err)
	}

	parts := strings.Split(manifest.Version, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("version format invalid: expected major.minor.patch, got %s", manifest.Version)
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return "", fmt.Errorf("invalid major version: %v", err)
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("invalid minor version: %v", err)
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", fmt.Errorf("invalid patch version: %v", err)
	}

	patch++
	newVersion := fmt.Sprintf("%d.%d.%d", major, minor, patch)

	manifest.Version = newVersion
	manifest.UpdatedAt = time.Now()

	checksums, err := calculateChecksums()
	if err != nil {
		return "", fmt.Errorf("error calculating checksums: %v", err)
	}
	manifest.Checksums = checksums

	updatedData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshaling manifest: %v", err)
	}

	if err := os.WriteFile(utils.ManifestFile, updatedData, 0644); err != nil {
		return "", fmt.Errorf("error writing manifest: %v", err)
	}

	return fmt.Sprintf("✔ Version bumped: %s → %s\n✔ Checksums calculated and updated", strings.Join(parts, "."), newVersion), nil
}

func writeSlugFile(filePath string, slugs map[string]bool) error {
	var lines []string
	for slug := range slugs {
		lines = append(lines, fmt.Sprintf(`{"slug": "%s"}`, slug))
	}
	data := strings.Join(lines, "\n") + "\n"
	return os.WriteFile(filePath, []byte(data), 0644)
}

func updateManifest(questionCount, themeCount, subthemeCount, tagCount int) error {
	data, err := os.ReadFile(utils.ManifestFile)
	if err != nil {
		return err
	}

	var manifest models.Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return err
	}

	manifest.Counts["questions"] = questionCount
	manifest.Counts["themes"] = themeCount
	manifest.Counts["subthemes"] = subthemeCount
	manifest.Counts["tags"] = tagCount

	updatedData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(utils.ManifestFile, updatedData, 0644)
}

func calculateChecksums() (map[string]string, error) {
	checksums := make(map[string]string)
	files := []string{
		utils.QuestionsFile,
		utils.ThemesFile,
		utils.SubthemesFile,
		utils.TagsFile,
	}

	for _, filePath := range files {
		hash, err := calculateSHA256(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				hash = calculateEmptySHA256()
			} else {
				return nil, fmt.Errorf("error calculating hash for %s: %v", filePath, err)
			}
		}
		fileName := strings.TrimPrefix(filePath, "datasets/general-knowledge/")
		checksums[fileName] = hash
	}

	return checksums, nil
}

func calculateSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashBytes := hash.Sum(nil)
	return "sha256-" + hex.EncodeToString(hashBytes), nil
}

func calculateEmptySHA256() string {
	return "sha256-e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
}
