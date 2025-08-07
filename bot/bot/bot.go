package bot

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/telego/telego"
	"github.com/telego/telego/telegoutil"
)

// Bot representa la estructura de tu bot de Telegram
type Bot struct {
	telegoBot  *telego.Bot
	updates    <-chan telego.Update
	userStates *userState
	done       chan struct{}
}

// userState almacena datos de los usuarios de forma segura
type userState struct {
	sync.RWMutex
	data map[int64]*UserData
}

// UserData almacena la información de cada usuario en el bot
type UserData struct {
	State              string
	OriginalMessageID  int64
	VideoFileName      string
	DownloadPath       string
	FinalPath          string
	ThumbnailPath      string
	NewName            string
	CompressionOptions map[string]string
	SendAsFile         bool
	LastUpdateTime     int64
	StatusMessageID    int64
}

// New crea una nueva instancia del bot
func New(apiID, apiHash, botToken string) (*Bot, error) {
	bot, err := telego.NewBot(botToken)
	if err != nil {
		return nil, fmt.Errorf("error al crear el bot de Telego: %w", err)
	}

	updates, err := bot.UpdatesViaWebhook(telego.With.WebhookSet(&telego.SetWebhookParams{
		URL: os.Getenv("WEBHOOK_URL"),
	}))
	if err != nil {
		return nil, fmt.Errorf("error al obtener actualizaciones vía webhook: %w", err)
	}

	return &Bot{
		telegoBot: bot,
		updates:   updates,
		userStates: &userState{
			data: make(map[int64]*UserData),
		},
		done: make(chan struct{}),
	}, nil
}

// Start inicia el bot y su bucle de actualización
func (b *Bot) Start() {
	log.Println("Bot en línea...")
	for update := range b.updates {
		go b.handleUpdate(update)
	}
}

// Stop detiene el bot
func (b *Bot) Stop() {
	close(b.done)
	b.telegoBot.StopWebhook()
}

// handleUpdate dirige cada actualización al manejador adecuado
func (b *Bot) handleUpdate(update telego.Update) {
	if update.Message != nil {
		if update.Message.Video != nil {
			b.videoHandler(update)
		} else if update.Message.Text != "" {
			b.textHandler(update)
		} else if update.Message.Photo != nil {
			b.thumbnailHandler(update)
		}
	} else if update.CallbackQuery != nil {
		b.callbackHandler(update)
	}
}

