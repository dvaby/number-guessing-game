package entity

import "time"

// Difficulty merepresentasikan tingkat kesulitan permainan
type Difficulty int

const (
	DifficultyEasy   Difficulty = 1
	DifficultyMedium Difficulty = 2
	DifficultyHard   Difficulty = 3
)

// String mengembalikan nama tingkat kesulitan yang mudah dibaca orang
func (d Difficulty) String() string {
	switch d {
	case DifficultyEasy:
		return "Easy"
	case DifficultyMedium:
		return "Medium"
	case DifficultyHard:
		return "Hard"
	default:
		return "Unknown"
	}
}

// MaxChances menentukan jumlah maksimal kesempatan berdasarkan tingkat kesulitan
func (d Difficulty) MaxChances() int {
	switch d {
	case DifficultyEasy:
		return 10
	case DifficultyMedium:
		return 5
	case DifficultyHard:
		return 3
	default:
		return 5
	}
}

// IsValid memvalidasi apakah angka pilihan difficulty termasuk yang tersedia
func (d Difficulty) IsValid() bool {
	return d == DifficultyEasy || d == DifficultyMedium || d == DifficultyHard
}

// GameSession menyimpan seluruh keadaan (state) dalam satu ronde permainan
type GameSession struct {
	TargetNumber int
	Difficulty   Difficulty
	MaxChances   int
	AttemptsUsed int
	HintsUsed    int
	MaxHints     int
	StartTime    time.Time
	IsWon        bool
	IsOver       bool
}

// RemainingChances menghitung sisa kesempatan menebak
func (s *GameSession) RemainingChances() int {
	return s.MaxChances - s.AttemptsUsed
}

// HighScore buat nyimpen rekor terbaik pada kesulitan tertentu
type HighScore struct {
	Difficulty Difficulty
	Attempts   int
	Duration   time.Duration
}