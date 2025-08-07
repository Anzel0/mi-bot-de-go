package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"mi-bot-de-go/bot"
)

func main() {
	apiID := os.Getenv("API_ID")
	apiHash := os.Getenv("API_HASH")
	botToken := os.Getenv("BOT_TOKEN")

	if apiID == "" || apiHash == "" || botToken == "" {
		log.Fatal("Variables de entorno API_ID, API_HASH o BOT_TOKEN no configuradas.")
	}

	myBot, err := bot.New(apiID, apiHash, botToken)
	if err != nil {
		log.Fatalf("Error al crear el bot: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go myBot.Start()

	<-quit
	log.Println("Deteniendo el bot...")
	myBot.Stop()
	log.Println("Bot detenido.")
}

