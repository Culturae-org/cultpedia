package actions

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"cultpedia/internal/checks"
	"cultpedia/internal/models"
	"cultpedia/internal/utils"
)

func ValidateNewQuestion() (models.Question, error) {
	return ValidateNewQuestionWithType("")
}

func ValidateNewQuestionWithType(forceType string) (models.Question, error) {
	var jsonFilePath string
	var questionType string

	if forceType == "true_false" {
		jsonFilePath = utils.NewQuestionTrueFalseFile
		questionType = "true_false"
	} else if forceType == "single_choice" {
		jsonFilePath = utils.NewQuestionFile
		questionType = "single_choice"
	} else {
		jsonFilePath, questionType = utils.DetectModifiedTemplateFile()
		if jsonFilePath == "" {
			jsonFilePath = utils.NewQuestionFile
			questionType = "single_choice"
		}
	}

	if _, err := os.Stat(jsonFilePath); os.IsNotExist(err) {
		return models.Question{}, fmt.Errorf("file %s not found", jsonFilePath)
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

	if questionType == "true_false" || question.Qtype == "true_false" {
		if len(question.Answers) != 2 {
			return models.Question{}, fmt.Errorf("true_false questions must have exactly 2 answers")
		}
		hasTrue, hasFalse := false, false
		for _, a := range question.Answers {
			if a.Slug == "true" {
				hasTrue = true
			}
			if a.Slug == "false" {
				hasFalse = true
			}
		}
		if !hasTrue || !hasFalse {
			return models.Question{}, fmt.Errorf("true_false questions must have answers with slugs 'true' and 'false'")
		}
	} else {
		if len(question.Answers) != 4 {
			return models.Question{}, fmt.Errorf("must have exactly 4 answers")
		}
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

	existingThemes, err := utils.LoadThemeTranslations(utils.ThemesFile)
	if err != nil {
		return fmt.Sprintf("✗ error reading existing themes: %v", err)
	}
	existingSubthemes, err := utils.LoadThemeTranslations(utils.SubthemesFile)
	if err != nil {
		return fmt.Sprintf("✗ error reading existing subthemes: %v", err)
	}
	existingTags, err := utils.LoadThemeTranslations(utils.TagsFile)
	if err != nil {
		return fmt.Sprintf("✗ error reading existing tags: %v", err)
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

	mergedThemes := mergeThemeTranslations(themeSlugs, existingThemes)
	mergedSubthemes := mergeThemeTranslations(subthemeSlugs, existingSubthemes)
	mergedTags := mergeThemeTranslations(tagSlugs, existingTags)

	if err := writeThemeFile(utils.ThemesFile, mergedThemes); err != nil {
		return fmt.Sprintf("✗ error writing themes: %v", err)
	}
	if err := writeThemeFile(utils.SubthemesFile, mergedSubthemes); err != nil {
		return fmt.Sprintf("✗ error writing subthemes: %v", err)
	}
	if err := writeThemeFile(utils.TagsFile, mergedTags); err != nil {
		return fmt.Sprintf("✗ error writing tags: %v", err)
	}

	if err := updateManifest(len(questions), len(mergedThemes), len(mergedSubthemes), len(mergedTags)); err != nil {
		return fmt.Sprintf("✗ error updating manifest: %v", err)
	}
	return fmt.Sprintf("✔ Themes synced successfully\n  - %d questions\n  - %d themes\n  - %d subthemes\n  - %d tags", len(questions), len(mergedThemes), len(mergedSubthemes), len(mergedTags))
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

	questions := manifest.Counts["questions"]
	if questions > 0 {
		if questions%1000 == 0 {
			minor++
			patch = 0
		} else {
			patch++
		}
	} else {
		patch++
	}
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

func BumpGeographyVersion() (string, error) {
	data, err := os.ReadFile(utils.GeographyManifestFile)
	if err != nil {
		return "", fmt.Errorf("error reading geography manifest: %v", err)
	}

	var manifest models.GeographyManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return "", fmt.Errorf("error parsing geography manifest: %v", err)
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

	// Bump patch version
	patch++
	newVersion := fmt.Sprintf("%d.%d.%d", major, minor, patch)

	manifest.Version = newVersion
	manifest.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	checksums, err := calculateGeographyChecksums()
	if err != nil {
		return "", fmt.Errorf("error calculating checksums: %v", err)
	}
	manifest.Checksums = checksums

	countries, _ := utils.LoadCountries()
	continents, _ := utils.LoadContinents()
	regions, _ := utils.LoadRegions()

	manifest.Counts["countries"] = len(countries)
	manifest.Counts["continents"] = len(continents)
	manifest.Counts["regions"] = len(regions)

	flagCount := 0
	if entries, err := os.ReadDir(utils.FlagsSVGDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".svg") {
				flagCount++
			}
		}
	}
	manifest.Counts["flags"] = flagCount

	updatedData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshaling manifest: %v", err)
	}

	if err := os.WriteFile(utils.GeographyManifestFile, updatedData, 0644); err != nil {
		return "", fmt.Errorf("error writing manifest: %v", err)
	}

	return fmt.Sprintf("✔ Geography version bumped: %s → %s\n✔ Checksums calculated and updated\n✔ Counts updated: %d countries, %d continents, %d regions, %d flags",
		strings.Join(parts, "."), newVersion, len(countries), len(continents), len(regions), flagCount), nil
}

func InitCultpediaDataset(targetDir, datasetName string) (string, error) {
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	filesToCreate := []string{
		"questions.ndjson",
		"themes.ndjson",
		"subthemes.ndjson",
		"tags.ndjson",
	}

	for _, filename := range filesToCreate {
		if err := os.WriteFile(fmt.Sprintf("%s/%s", targetDir, filename), []byte{}, 0644); err != nil {
			return "", fmt.Errorf("failed to create %s: %v", filename, err)
		}
	}

	manifest := models.NewQuestionManifest(datasetName)
	manifest.Counts = map[string]int{
		"questions": 0,
		"themes":    0,
		"subthemes": 0,
		"tags":      0,
	}
	manifest.Checksums = map[string]string{
		"questions.ndjson": calculateEmptySHA256(),
		"themes.ndjson":    calculateEmptySHA256(),
		"subthemes.ndjson": calculateEmptySHA256(),
		"tags.ndjson":      calculateEmptySHA256(),
	}

	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal manifest: %v", err)
	}

	manifestPath := fmt.Sprintf("%s/manifest.json", targetDir)
	if err := os.WriteFile(manifestPath, manifestData, 0644); err != nil {
		return "", fmt.Errorf("failed to write manifest: %v", err)
	}

	return fmt.Sprintf("✔ Dataset '%s' initialized successfully", manifest.Dataset), nil
}

// --------------------------------
// Helper functions
// --------------------------------
func mergeThemeTranslations(slugs map[string]bool, existing map[string]models.ThemeTranslation) map[string]models.ThemeTranslation {
	merged := make(map[string]models.ThemeTranslation)
	for slug := range slugs {
		if t, ok := existing[slug]; ok {
			merged[slug] = t
		} else {
			merged[slug] = models.ThemeTranslation{
				Slug: slug,
				I18n: map[string]models.Label{
					"fr": {Label: ""},
					"en": {Label: ""},
					"es": {Label: ""},
				},
			}
		}
	}
	return merged
}

func writeThemeFile(filePath string, themes map[string]models.ThemeTranslation) error {
	slugs := make([]string, 0, len(themes))
	for slug := range themes {
		slugs = append(slugs, slug)
	}
	sort.Strings(slugs)

	var lines []string
	for _, slug := range slugs {
		data, err := json.Marshal(themes[slug])
		if err != nil {
			return fmt.Errorf("error marshaling theme %s: %v", slug, err)
		}
		lines = append(lines, string(data))
	}
	content := strings.Join(lines, "\n") + "\n"
	return os.WriteFile(filePath, []byte(content), 0644)
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

func calculateGeographyChecksums() (map[string]string, error) {
	checksums := make(map[string]string)
	files := []string{
		utils.CountriesFile,
		utils.ContinentsFile,
		utils.RegionsFile,
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
		fileName := strings.TrimPrefix(filePath, "datasets/geography/")
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

func GetRemoteVersion() (string, error) {
	resp, err := http.Get("https://raw.githubusercontent.com/Culturae-org/cultpedia/main/datasets/general-knowledge/manifest.json")
	if err != nil {
		return "", fmt.Errorf("failed to fetch remote manifest: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("remote manifest not found (status %d)", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	var manifest models.Manifest
	if err := json.Unmarshal(body, &manifest); err != nil {
		return "", fmt.Errorf("failed to parse remote manifest: %v", err)
	}

	return manifest.Version, nil
}

func ShowStruct(datasetName string) {
	helpText := fmt.Sprintf(`
. %s
├── questions.ndjson
├── themes.ndjson
├── subthemes.ndjson
├── tags.ndjson
└── manifest.json
	`, datasetName)
	fmt.Println(helpText)
}

func calculateEmptySHA256() string {
	return "sha256-e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
}

func ResetTemplate(questionType string) error {
	if questionType == "true_false" {
		return resetTrueFalseTemplate()
	}
	return resetSingleChoiceTemplate()
}

func resetSingleChoiceTemplate() error {
	template := `{
  "kind": "question",
  "version": "1.0",
  "slug": "default-question-slug",
  "theme": { "slug": "default-theme" },
  "subthemes": [ { "slug": "default-subtheme1" }, { "slug": "default-subtheme2" } ],
  "tags": [ { "slug": "default-tag1" }, { "slug": "default-tag2" } ],
  "qtype": "single_choice",
  "difficulty": "beginner",
  "estimated_seconds": 15,
  "shuffle_answers": true,
  "i18n": {
    "fr": { "title": "Titre par défaut", "stem": "Question par défaut ?", "explanation": "Explication par défaut." },
    "en": { "title": "Default Title", "stem": "Default question?", "explanation": "Default explanation." },
    "es": { "title": "Título por defecto", "stem": "¿Pregunta por defecto?", "explanation": "Explicación por defecto." }
  },
  "answers": [
    { "slug": "default-answer1", "is_correct": true,  "i18n": { "fr": { "label": "Réponse 1" }, "en": { "label": "Answer 1" }, "es": { "label": "Respuesta 1" } } },
    { "slug": "default-answer2", "is_correct": false, "i18n": { "fr": { "label": "Réponse 2" }, "en": { "label": "Answer 2" }, "es": { "label": "Respuesta 2" } } },
    { "slug": "default-answer3", "is_correct": false, "i18n": { "fr": { "label": "Réponse 3" }, "en": { "label": "Answer 3" }, "es": { "label": "Respuesta 3" } } },
    { "slug": "default-answer4", "is_correct": false, "i18n": { "fr": { "label": "Réponse 4" }, "en": { "label": "Answer 4" }, "es": { "label": "Respuesta 4" } } }
  ],
  "sources": [
    "https://example.com/default-source"
  ]
}
`
	return os.WriteFile(utils.NewQuestionFile, []byte(template), 0644)
}

func resetTrueFalseTemplate() error {
	template := `{
  "kind": "question",
  "version": "1.0",
  "slug": "default-true-false-question-slug",
  "theme": { "slug": "default-theme" },
  "subthemes": [ { "slug": "default-subtheme1" }, { "slug": "default-subtheme2" } ],
  "tags": [ { "slug": "default-tag1" }, { "slug": "default-tag2" } ],
  "qtype": "true_false",
  "difficulty": "beginner",
  "estimated_seconds": 10,
  "shuffle_answers": false,
  "i18n": {
    "fr": { "title": "Titre vrai/faux", "stem": "Cette affirmation est-elle vraie ou fausse?", "explanation": "Explication par défaut à remplacer par votre explication détaillée." },
    "en": { "title": "True/False Title", "stem": "Is this statement true or false?", "explanation": "Default explanation to replace with your detailed explanation." },
    "es": { "title": "Titulo verdadero/falso", "stem": "Es verdadera o falsa esta afirmacion?", "explanation": "Explicación por defecto a reemplazar con su explicación detallada." }
  },
  "answers": [
    { "slug": "true", "is_correct": true, "i18n": { "fr": { "label": "Vrai" }, "en": { "label": "True" }, "es": { "label": "Verdadero" } } },
    { "slug": "false", "is_correct": false, "i18n": { "fr": { "label": "Faux" }, "en": { "label": "False" }, "es": { "label": "Falso" } } }
  ],
  "sources": [
    "https://example.com/default-source"
  ]
}
`
	return os.WriteFile(utils.NewQuestionTrueFalseFile, []byte(template), 0644)
}

func GetAvailableThemes() ([]string, error) {
	questions, err := utils.LoadQuestions()
	if err != nil {
		return nil, err
	}

	themeSlugs := make(map[string]bool)
	for _, q := range questions {
		themeSlugs[q.Theme.Slug] = true
	}

	themes := make([]string, 0, len(themeSlugs))
	for slug := range themeSlugs {
		themes = append(themes, slug)
	}
	return themes, nil
}
