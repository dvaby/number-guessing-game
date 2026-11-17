package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"number-guessing-game/internal/domain/entity"
	"number-guessing-game/internal/usecase"
)

// CLIHandler bertanggung jawab menangani input dan output di terminal.
type CLIHandler struct {
	gameUsecase usecase.GameUsecase
	scanner     *bufio.Scanner
}

// NewCLIHandler membuat instance baru CLIHandler.
func NewCLIHandler(gameUsecase usecase.GameUsecase) *CLIHandler {
	return &CLIHandler{
		gameUsecase: gameUsecase,
		scanner:     bufio.NewScanner(os.Stdin),
	}
}

// Run menjalankan siklus utama aplikasi CLI (mendukung bermain berulang kali).
func (h *CLIHandler) Run() {
	h.printWelcomeBanner()

	for {
		difficulty := h.promptDifficulty()
		session, err := h.gameUsecase.StartNewGame(difficulty)
		if err != nil {
			fmt.Printf("Terjadi kesalahan saat memulai game: %v\n", err)
			continue
		}

		fmt.Printf("\nGreat! You have selected the %s difficulty level.\n", difficulty.String())
		fmt.Printf("You have %d chances to guess the correct number.\n", session.MaxChances)
		fmt.Println("Tip: Ketik 'hint' kapan saja jika butuh petunjuk (maksimal 2x per ronde).")
		fmt.Println("Let's start the game!")

		h.playRound(session)
		h.printHighScores()

		if !h.promptPlayAgain() {
			fmt.Println("\nTerima kasih sudah bermain! Sampai jumpa lagi.")
			break
		}
		fmt.Println("\n--------------------------------------------------")
	}
}

func (h *CLIHandler) printWelcomeBanner() {
	fmt.Println("==================================================")
	fmt.Println("       Welcome to the Number Guessing Game!       ")
	fmt.Println("==================================================")
	fmt.Println("Rules:")
	fmt.Println("- I'm thinking of a number between 1 and 100.")
	fmt.Println("- Choose a difficulty level to set your chances.")
	fmt.Println("- Guess the number before your chances run out!")
	fmt.Println()
}

func (h *CLIHandler) promptDifficulty() entity.Difficulty {
	for {
		fmt.Println("Please select the difficulty level:")
		fmt.Println("1. Easy (10 chances)")
		fmt.Println("2. Medium (5 chances)")
		fmt.Println("3. Hard (3 chances)")
		fmt.Print("\nEnter your choice: ")

		input := h.readLine()
		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("[!] Input tidak valid. Masukkan angka 1, 2, atau 3.\n")
			continue
		}

		difficulty := entity.Difficulty(choice)
		if !difficulty.IsValid() {
			fmt.Println("[!] Pilihan tidak tersedia. Silakan pilih 1, 2, atau 3.\n")
			continue
		}

		return difficulty
	}
}

func (h *CLIHandler) playRound(session *entity.GameSession) {
	for !session.IsOver {
		fmt.Print("\nEnter your guess: ")
		input := h.readLine()

		// Cek apakah pengguna meminta hint
		if strings.EqualFold(input, "hint") {
			clue, remainingHints, err := h.gameUsecase.GetHint(session)
			if err != nil {
				fmt.Printf("[!] %v\n", err)
			} else {
				fmt.Printf("[HINT] %s (Sisa hint: %d)\n", clue, remainingHints)
			}
			continue
		}

		guess, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("[!] Masukkan angka valid antara 1-100, atau ketik 'hint'.")
			continue
		}

		result, err := h.gameUsecase.MakeGuess(session, guess)
		if err != nil {
			fmt.Printf("[!] %v\n", err)
			continue
		}

		if result.IsWon {
			fmt.Printf("\nCongratulations! You guessed the correct number in %d attempts.\n", result.AttemptsUsed)
			fmt.Printf("Time taken: %.2f seconds.\n", result.ElapsedDuration.Seconds())
			if result.IsNewHighScore {
				fmt.Printf("*** NEW HIGH SCORE for %s difficulty! ***\n", session.Difficulty.String())
			}
			return
		}

		if result.Direction == usecase.DirectionLess {
			fmt.Printf("Incorrect! The number is less than %d.\n", guess)
		} else {
			fmt.Printf("Incorrect! The number is greater than %d.\n", guess)
		}

		if result.IsGameOver {
			fmt.Printf("\nGame Over! You've run out of chances. The correct number was %d.\n", result.TargetNumber)
			fmt.Printf("Time taken: %.2f seconds.\n", result.ElapsedDuration.Seconds())
		} else {
			fmt.Printf("Remaining chances: %d\n", result.RemainingChances)
		}
	}
}

func (h *CLIHandler) printHighScores() {
	scores := h.gameUsecase.GetHighScores()
	if len(scores) == 0 {
		return
	}

	fmt.Println("\n--- HIGH SCORES (Session Best) ---")
	difficulties := []entity.Difficulty{
		entity.DifficultyEasy,
		entity.DifficultyMedium,
		entity.DifficultyHard,
	}

	for _, diff := range difficulties {
		if score, exists := scores[diff]; exists {
			fmt.Printf("- %-6s : %d attempts (%.2fs)\n",
				diff.String(),
				score.Attempts,
				score.Duration.Seconds(),
			)
		}
	}
}

func (h *CLIHandler) promptPlayAgain() bool {
	for {
		fmt.Print("\nDo you want to play again? (y/n): ")
		answer := strings.ToLower(h.readLine())

		switch answer {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			fmt.Println("[!] Ketik 'y' untuk lanjut bermain atau 'n' untuk keluar.")
		}
	}
}

func (h *CLIHandler) readLine() string {
	if h.scanner.Scan() {
		return strings.TrimSpace(h.scanner.Text())
	}
	return ""
}