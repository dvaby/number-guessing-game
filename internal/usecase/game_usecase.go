package usecase

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"time"

	"number-guessing-game/internal/domain/entity"
	"number-guessing-game/internal/domain/repository"
)

// Sentinel Errors: daftar pesan kesalahan baku yang bisa terjadi dalam aturan permainan.
var (
	ErrInvalidDifficulty = errors.New("pilihan tingkat kesulitan tidak valid")
	ErrGuessOutOfRange   = errors.New("tebakan harus berada di antara angka 1 sampai 100")
	ErrGameAlreadyOver   = errors.New("permainan pada ronde ini sudah berakhir")
	ErrNoHintsRemaining  = errors.New("jatah hint pada ronde ini sudah habis")
)

// GuessDirection menandakan posisi angka rahasia terhadap tebakan pengguna.
type GuessDirection int

const (
	DirectionCorrect GuessDirection = iota
	DirectionGreater
	DirectionLess
)

// GuessResult adalah hasil evaluasi dari satu kali tebakan pengguna.
type GuessResult struct {
	Direction        GuessDirection
	AttemptsUsed     int
	RemainingChances int
	IsGameOver       bool
	IsWon            bool
	TargetNumber     int
	ElapsedDuration  time.Duration
	IsNewHighScore   bool
}

// GameUsecase mendefinisikan layanan bisnis utama untuk permainan tebak angka.
type GameUsecase interface {
	StartNewGame(difficulty entity.Difficulty) (*entity.GameSession, error)
	MakeGuess(session *entity.GameSession, guess int) (*GuessResult, error)
	GetHint(session *entity.GameSession) (string, int, error)
	GetHighScores() map[entity.Difficulty]entity.HighScore
}

type gameUsecase struct {
	scoreRepo repository.ScoreRepository
}

// NewGameUsecase membuat instance GameUsecase dengan menyuntikkan dependency repository.
func NewGameUsecase(scoreRepo repository.ScoreRepository) GameUsecase {
	return &gameUsecase{
		scoreRepo: scoreRepo,
	}
}

func (u *gameUsecase) StartNewGame(difficulty entity.Difficulty) (*entity.GameSession, error) {
	if !difficulty.IsValid() {
		return nil, ErrInvalidDifficulty
	}

	// rand.IntN(100) menghasilkan angka acak 0 sampai 99, ditambah 1 menjadi 1 sampai 100.
	secretNumber := rand.IntN(100) + 1

	session := &entity.GameSession{
		TargetNumber: secretNumber,
		Difficulty:   difficulty,
		MaxChances:   difficulty.MaxChances(),
		AttemptsUsed: 0,
		HintsUsed:    0,
		MaxHints:     2,
		StartTime:    time.Now(),
		IsWon:        false,
		IsOver:       false,
	}

	return session, nil
}

func (u *gameUsecase) MakeGuess(session *entity.GameSession, guess int) (*GuessResult, error) {
	if session.IsOver {
		return nil, ErrGameAlreadyOver
	}

	if guess < 1 || guess > 100 {
		return nil, ErrGuessOutOfRange
	}

	session.AttemptsUsed++
	elapsed := time.Since(session.StartTime)

	// Kondisi 1: Tebakan Benar
	if guess == session.TargetNumber {
		session.IsWon = true
		session.IsOver = true

		isNewRecord := u.checkAndSaveHighScore(session.Difficulty, session.AttemptsUsed, elapsed)

		return &GuessResult{
			Direction:        DirectionCorrect,
			AttemptsUsed:     session.AttemptsUsed,
			RemainingChances: session.RemainingChances(),
			IsGameOver:       true,
			IsWon:            true,
			TargetNumber:     session.TargetNumber,
			ElapsedDuration:  elapsed,
			IsNewHighScore:   isNewRecord,
		}, nil
	}

	// Kondisi 2: Tebakan Salah
	direction := DirectionGreater
	if session.TargetNumber < guess {
		direction = DirectionLess
	}

	if session.RemainingChances() <= 0 {
		session.IsOver = true
	}

	return &GuessResult{
		Direction:        direction,
		AttemptsUsed:     session.AttemptsUsed,
		RemainingChances: session.RemainingChances(),
		IsGameOver:       session.IsOver,
		IsWon:            false,
		TargetNumber:     session.TargetNumber,
		ElapsedDuration:  elapsed,
		IsNewHighScore:   false,
	}, nil
}

func (u *gameUsecase) GetHint(session *entity.GameSession) (string, int, error) {
	if session.IsOver {
		return "", 0, ErrGameAlreadyOver
	}

	if session.HintsUsed >= session.MaxHints {
		return "", 0, ErrNoHintsRemaining
	}

	session.HintsUsed++
	remainingHints := session.MaxHints - session.HintsUsed

	var clue string
	switch session.HintsUsed {
	case 1:
		// Hint 1: Menjelaskan apakah angka rahasia ganjil atau genap
		if session.TargetNumber%2 == 0 {
			clue = "Angka rahasia tersebut adalah bilangan GENAP."
		} else {
			clue = "Angka rahasia tersebut adalah bilangan GANJIL."
		}
	case 2:
		// Hint 2: Memberikan kisaran rentang terdekat
		low := max(1, session.TargetNumber-10)
		high := min(100, session.TargetNumber+10)
		clue = fmt.Sprintf("Angka rahasia berada di kisaran antara %d sampai %d.", low, high)
	}

	return clue, remainingHints, nil
}

func (u *gameUsecase) GetHighScores() map[entity.Difficulty]entity.HighScore {
	return u.scoreRepo.GetAll()
}

func (u *gameUsecase) checkAndSaveHighScore(difficulty entity.Difficulty, attempts int, duration time.Duration) bool {
	existing, found := u.scoreRepo.GetByDifficulty(difficulty)

	// Rekor terpecahkan jika belum ada skor sebelumnya,
	// atau jumlah percobaan lebih sedikit,
	// atau percobaan sama namun waktu penyelesaian lebih cepat.
	isBetter := !found ||
		attempts < existing.Attempts ||
		(attempts == existing.Attempts && duration < existing.Duration)

	if isBetter {
		_ = u.scoreRepo.Save(entity.HighScore{
			Difficulty: difficulty,
			Attempts:   attempts,
			Duration:   duration,
		})
		return true
	}

	return false
}