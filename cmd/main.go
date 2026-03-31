package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"cultpedia/internal/actions"
	"cultpedia/internal/checks"
	"cultpedia/internal/models"
	"cultpedia/internal/utils"
)

func main() {
	if len(os.Args) > 1 {
		handleCommand(os.Args[1], os.Args[2:])
		return
	}
	utils.PrintHelp()
}

func handleCommand(cmd string, args []string) {
	switch cmd {
	case "help", "--help", "-h":
		utils.PrintHelp()
		os.Exit(0)
	case "validate":
		err := checks.ValidateQuestions()
		if err != nil {
			fmt.Println("✗ Validation Failed:")
			fmt.Println()
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Println("✔ Validation Successful - All questions are valid!")
		}
	case "check-duplicates":
		result := checks.CheckDuplicates()
		fmt.Println(result)
		if strings.Contains(result, "detected") {
			os.Exit(1)
		}
	case "check-translations":
		result := checks.CheckTranslations()
		fmt.Println(result)
		if strings.Contains(result, "missing") {
			os.Exit(1)
		}
	case "add":
		handleAdd(args)
	case "preview":
		handlePreview(args)
	case "version":
		handleVersion()
	case "sync-themes":
		result := actions.SyncThemes()
		fmt.Println(result)
		if strings.Contains(result, "error") {
			os.Exit(1)
		}
	case "bump-version":
		version, err := actions.BumpVersion()
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(version)
	case "validate-geography":
		err := checks.ValidateGeography()
		if err != nil {
			fmt.Println("✗ Geography Validation Failed:")
			fmt.Println()
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Println("✔ Geography Validation Successful - All data is valid!")
		}
	case "check-geography-duplicates":
		result := checks.CheckGeographyDuplicates()
		fmt.Println(result)
		if strings.HasPrefix(result, "✗") {
			os.Exit(1)
		}
	case "check-geography-translations":
		result := checks.CheckGeographyTranslations()
		fmt.Println(result)
		if strings.HasPrefix(result, "✗") {
			os.Exit(1)
		}
	case "bump-geography-version":
		version, err := actions.BumpGeographyVersion()
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(version)
	case "init":
		defaultDir := "new-cultpedia-dataset"
		datasetName := "new-cultpedia-dataset"

		if len(args) > 0 {
			defaultDir = args[0]
			datasetName = args[0]
		}

		message, err := actions.InitCultpediaDataset(defaultDir, datasetName)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✔ " + message)
		actions.ShowStruct(datasetName)
	case "api":
		if len(args) > 0 {
			actions.RunAPIServer(args[0])
		} else {
			actions.RunAPIServer("8080")
		}

	// AI Generation commands
	case "generate":
		handleGenerate(args)

	case "pending":
		handlePending()

	case "approve":
		handleApprove(args)

	case "reject":
		handleReject(args)

	default:
		fmt.Printf("unknown command: %s\n", cmd)
		fmt.Println("use 'cultpedia help' to see available commands")
		os.Exit(1)
	}
}

func handleAdd(args []string) {
	qtype := ""
	for i := 0; i < len(args); i++ {
		if args[i] == "--type" && i+1 < len(args) {
			qtype = args[i+1]
			i++
		}
	}

	var question models.Question
	var err error
	if qtype != "" {
		question, err = actions.ValidateNewQuestionWithType(qtype)
	} else {
		question, err = actions.ValidateNewQuestion()
	}
	if err != nil {
		fmt.Println("✗ Cannot add question:")
		fmt.Println()
		fmt.Println(err)
		os.Exit(1)
	}

	message := actions.AddValidatedQuestion(question)
	fmt.Println("✔ " + message)

	if !strings.Contains(message, "error") {
		templateType := qtype
		if templateType == "" {
			if question.Qtype == "true_false" {
				templateType = "true_false"
			} else {
				templateType = "single_choice"
			}
		}
		if err := actions.ResetTemplate(templateType); err == nil {
			fmt.Println("\n✔ Template file has been reset for your next question.")
		}
	}

	if strings.Contains(message, "error") {
		os.Exit(1)
	}
}

func handlePreview(args []string) {
	qtype := ""
	for i := 0; i < len(args); i++ {
		if args[i] == "--type" && i+1 < len(args) {
			qtype = args[i+1]
			i++
		}
	}

	var question models.Question
	var err error
	if qtype != "" {
		question, err = actions.ValidateNewQuestionWithType(qtype)
	} else {
		question, err = actions.ValidateNewQuestion()
	}
	if err != nil {
		fmt.Println("✗ Cannot preview question:")
		fmt.Println()
		fmt.Println(err)
		os.Exit(1)
	}

	printQuestionPreview(question)
}

func printQuestionPreview(q models.Question) {
	fmt.Println()
	fmt.Println("── Question Preview ──")
	fmt.Println()
	fmt.Printf("  Slug:       %s\n", q.Slug)
	fmt.Printf("  Theme:      %s\n", q.Theme.Slug)

	if len(q.Subthemes) > 0 {
		subs := make([]string, len(q.Subthemes))
		for i, s := range q.Subthemes {
			subs[i] = s.Slug
		}
		fmt.Printf("  Subthemes:  %s\n", strings.Join(subs, ", "))
	}

	if len(q.Tags) > 0 {
		tags := make([]string, len(q.Tags))
		for i, t := range q.Tags {
			tags[i] = t.Slug
		}
		fmt.Printf("  Tags:       %s\n", strings.Join(tags, ", "))
	}

	fmt.Printf("  Difficulty: %s\n", q.Difficulty)
	fmt.Printf("  Type:       %s\n", q.Qtype)
	fmt.Println()

	langs := []string{"fr", "en", "es"}
	for _, lang := range langs {
		content, ok := q.I18n[lang]
		if !ok {
			continue
		}
		fmt.Printf("  [%s] Title:       %s\n", strings.ToUpper(lang), content.Title)
		fmt.Printf("  [%s] Stem:        %s\n", strings.ToUpper(lang), content.Stem)
		fmt.Printf("  [%s] Explanation: %s\n", strings.ToUpper(lang), content.Explanation)
		fmt.Println()
	}

	fmt.Println("  Answers:")
	for _, a := range q.Answers {
		mark := " "
		if a.IsCorrect {
			mark = "✓"
		}
		label := a.Slug
		if l, ok := a.I18n["en"]; ok {
			label = l.Label
		}
		fmt.Printf("    [%s] %s (%s)\n", mark, label, a.Slug)
	}
	fmt.Println()

	if len(q.Sources) > 0 {
		fmt.Println("  Sources:")
		for _, src := range q.Sources {
			fmt.Printf("    - %s\n", src)
		}
		fmt.Println()
	}
}

func handleVersion() {
	fmt.Printf("Cultpedia Version (API/CLI): %s\n\n", utils.Version)

	version := "unknown"
	if data, err := os.ReadFile(utils.ManifestFile); err == nil {
		var manifest models.Manifest
		if json.Unmarshal(data, &manifest) == nil {
			version = manifest.Version
		}
	}

	fmt.Printf("Local dataset version:  %s\n", version)

	remoteVersion, err := actions.GetRemoteVersion()
	if err != nil {
		fmt.Printf("Remote version: unavailable (%v)\n", err)
		return
	}

	fmt.Printf("Remote version: %s\n", remoteVersion)

	if version == remoteVersion {
		fmt.Println("\n✔ Your dataset is up to date.")
	} else {
		fmt.Println("\n⚠ Your local dataset is outdated. Run 'git pull' to update.")
	}
}

func handleGenerate(args []string) {
	if len(args) == 0 {
		fmt.Println("✗ Usage: cultpedia generate <topic> [options]")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  --theme <theme>       Theme slug (default: science)")
		fmt.Println("  --count <n>           Number of questions (default: 3, max: 10)")
		fmt.Println("  --difficulty <level>  beginner, intermediate, advanced (default: intermediate)")
		fmt.Println("  --type <qtype>        single_choice or true_false (default: single_choice)")
		fmt.Println()
		fmt.Println("Example:")
		fmt.Println("  cultpedia generate \"French Revolution\" --theme history --count 5")
		os.Exit(1)
	}

	config := actions.GenerateConfig{
		Topic:      args[0],
		Theme:      "science",
		Difficulty: "intermediate",
		Qtype:      "single_choice",
		Count:      3,
	}

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--theme":
			if i+1 < len(args) {
				config.Theme = args[i+1]
				i++
			}
		case "--count":
			if i+1 < len(args) {
				if n, err := strconv.Atoi(args[i+1]); err == nil {
					config.Count = n
				}
				i++
			}
		case "--difficulty":
			if i+1 < len(args) {
				config.Difficulty = args[i+1]
				i++
			}
		case "--type":
			if i+1 < len(args) {
				config.Qtype = args[i+1]
				i++
			}
		}
	}

	fmt.Printf("🤖 Generating %d questions about \"%s\"...\n", config.Count, config.Topic)
	fmt.Printf("   Theme: %s | Difficulty: %s | Type: %s\n\n", config.Theme, config.Difficulty, config.Qtype)

	questions, err := actions.GenerateQuestions(config)
	if err != nil {
		fmt.Printf("✗ Generation failed: %v\n", err)
		os.Exit(1)
	}

	if len(questions) == 0 {
		fmt.Println("✗ No questions were generated")
		os.Exit(1)
	}

	if err := actions.SavePendingQuestions(questions); err != nil {
		fmt.Printf("✗ Error saving questions: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✔ Generated %d questions and saved to pending\n\n", len(questions))
	for i, q := range questions {
		fmt.Println(actions.FormatPendingQuestionPreview(q, i))
	}
	fmt.Println("Next steps:")
	fmt.Println("  cultpedia pending              - View all pending questions")
	fmt.Println("  cultpedia approve <index>      - Approve a question")
	fmt.Println("  cultpedia reject <index>       - Reject a question")
}

func handlePending() {
	pending, err := actions.LoadPendingQuestions()
	if err != nil {
		fmt.Printf("✗ Error loading pending questions: %v\n", err)
		os.Exit(1)
	}

	if len(pending) == 0 {
		fmt.Println("📭 No pending questions")
		fmt.Println()
		fmt.Println("Generate some with:")
		fmt.Println("  cultpedia generate \"topic\" --theme science --count 3")
		return
	}

	fmt.Printf("📋 %d pending question(s):\n\n", len(pending))
	for i, q := range pending {
		fmt.Println(actions.FormatPendingQuestionPreview(q, i))
	}
	fmt.Println("Commands:")
	fmt.Println("  cultpedia approve <index>  - Approve and add to dataset")
	fmt.Println("  cultpedia reject <index>   - Reject and remove")
}

func handleApprove(args []string) {
	if len(args) == 0 {
		fmt.Println("✗ Usage: cultpedia approve <index>")
		fmt.Println("  Use 'cultpedia pending' to see available indices")
		os.Exit(1)
	}

	index, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("✗ Invalid index: %s\n", args[0])
		os.Exit(1)
	}

	result, err := actions.ApprovePendingQuestion(index)
	if err != nil {
		fmt.Printf("✗ %v\n", err)
		os.Exit(1)
	}
	fmt.Println(result)
}

func handleReject(args []string) {
	if len(args) == 0 {
		fmt.Println("✗ Usage: cultpedia reject <index>")
		fmt.Println("  Use 'cultpedia pending' to see available indices")
		os.Exit(1)
	}

	index, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("✗ Invalid index: %s\n", args[0])
		os.Exit(1)
	}

	result, err := actions.RejectPendingQuestion(index)
	if err != nil {
		fmt.Printf("✗ %v\n", err)
		os.Exit(1)
	}
	fmt.Println(result)
}
