package memory

import (
	"sync"

	"number-guessing-game/internal/domain/entity"
	"number-guessing-game/internal/domain/repository"
)

// scoreMemoryRepository adalah implementasi ScoreRepository yang menyimpan data di RAM.
type scoreMemoryRepository struct {
	mu     sync.RWMutex
	scores map[entity.Difficulty]entity.HighScore
}

// NewScoreMemoryRepository membuat instance baru penyimpanan High Score di memori.
func NewScoreMemoryRepository() repository.ScoreRepository {
	return &scoreMemoryRepository{
		scores: make(map[entity.Difficulty]entity.HighScore),
	}
}

func (r *scoreMemoryRepository) GetByDifficulty(difficulty entity.Difficulty) (*entity.HighScore, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	score, exists := r.scores[difficulty]
	if !exists {
		return nil, false
	}
	return &score, true
}

func (r *scoreMemoryRepository) Save(score entity.HighScore) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.scores[score.Difficulty] = score
	return nil
}

func (r *scoreMemoryRepository) GetAll() map[entity.Difficulty]entity.HighScore {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Mengembalikan salinan (copy) map agar data asli tidak termodifikasi secara tidak sengaja
	result := make(map[entity.Difficulty]entity.HighScore, len(r.scores))
	for k, v := range r.scores {
		result[k] = v
	}
	return result
}