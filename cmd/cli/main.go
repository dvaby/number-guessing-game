package main

import (
	"number-guessing-game/internal/delivery/cli"
	"number-guessing-game/internal/repository/memory"
	"number-guessing-game/internal/usecase"
)

func main() {
	// 1. Inisialisasi Layer Repository (Penyimpanan Data)
	scoreRepo := memory.NewScoreMemoryRepository()

	// 2. Inisialisasi Layer Usecase (Business Logic) dan suntikkan Repository
	gameUsecase := usecase.NewGameUsecase(scoreRepo)

	// 3. Inisialisasi Layer Delivery (CLI Handler) dan suntikkan Usecase
	cliHandler := cli.NewCLIHandler(gameUsecase)

	// 4. Jalankan aplikasi
	cliHandler.Run()
}