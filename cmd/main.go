package main

import (
	"fmt"
	"os"
	"strings"

	"cultpedia/internal/actions"
	"cultpedia/internal/checks"
	"cultpedia/internal/ui"
	"cultpedia/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) > 1 {
		handleCommand(os.Args[1])
		return
	}
	p := tea.NewProgram(ui.InitialMainModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	clearTerminal()
}

func clearTerminal() {
	fmt.Print("\033[H\033[2J")
}

func handleCommand(cmd string) {
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
		question, err := actions.ValidateNewQuestion()
		if err != nil {
			fmt.Println("✗ Cannot add question:")
			fmt.Println()
			fmt.Println(err)
			os.Exit(1)
		}
		message := actions.AddValidatedQuestion(question)
		fmt.Println("✔ " + message)
		if strings.Contains(message, "error") {
			os.Exit(1)
		}
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
	default:
		fmt.Printf("unknown command: %s\n", cmd)
		fmt.Printf("use 'cultpedia help' to see available commands or if you are a contributor, please use the interactive UI with ./cultpedia\n")
		os.Exit(1)
	}
}
