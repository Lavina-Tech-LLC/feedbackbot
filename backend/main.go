package main

import (
	"log"

	"github.com/Lavina-Tech-LLC/feedbackbot/internal/config"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/db"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/db/models"
	"github.com/Lavina-Tech-LLC/feedbackbot/internal/tgbot"
	webServer "github.com/Lavina-Tech-LLC/feedbackbot/internal/webserver"
	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
)

func main() {
	config.Init()
	db.Init()

	go webServer.Listen()

	// Start polling for all verified bots
	var bots []models.Bot
	models.DB.Where("verified = ?", true).Find(&bots)
	for _, bot := range bots {
		log.Printf("[main] Starting bot polling for @%s", bot.BotUsername)
		go tgbot.StartPolling(bot)
	}

	lvn.WaitExitSignal()
	log.Println("[main] Shutting down bot polling...")
	tgbot.StopAll()
}
