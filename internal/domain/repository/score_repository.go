package repository

import "number-guessing-game/internal/domain/entity"

// ScoreRepository adalah kontrak (interface) untuk menyimpan dan mengambil High Score.
// Layer Usecase hanya bergantung pada kontrak ini, bukan pada jenis database-nya.
type ScoreRepository interface {
	GetByDifficulty(difficulty entity.Difficulty) (*entity.HighScore, bool)
	Save(score entity.HighScore) error
	GetAll() map[entity.Difficulty]entity.HighScore
}