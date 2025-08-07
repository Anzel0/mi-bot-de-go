package bot

import (
	"log"
	"os"
	"time"

	"github.com/telego/telego"
	tu "github.com/telego/telego/telegoutil"
)

const (
	MaxVideoSizeMB = 4000
	DownloadDir    = "downloads"
)

func init() {
	os.MkdirAll(DownloadDir, os.ModePerm)
}

// startHandler maneja el comando /start
func (b *Bot) startHandler(update telego.Update) {
	chatID := update.Message.Chat.ID
	b.cleanup(chatID)
	b.telegoBot.SendMessage(tu.Message(tu.ID(chatID), "¡Hola! 👋 Soy tu bot para procesar videos.\n\nEnvíame un video para empezar."))
}

// videoHandler maneja los mensajes que contienen un video.
func (b *Bot) videoHandler(update telego.Update) {
	chatID := update.Message.Chat.ID
	if b.userHasActiveProcess(chatID) {
		b.telegoBot.SendMessage(tu.Message(tu.ID(chatID), "⚠️ Un proceso anterior se ha cancelado para iniciar uno nuevo."))
		b.cleanup(chatID)
	}

	if update.Message.Video.FileSize > MaxVideoSizeMB*1024*1024 {
		b.telegoBot.SendMessage(tu.Message(tu.ID(chatID), "❌ El video supera el límite de 4000 MB."))
		return
	}

	// Guardar estado inicial del usuario
	b.userStates.Lock()
	b.userStates.data[chatID] = &UserData{
		State:             "awaiting_action",
		OriginalMessageID: update.Message.ID,
		VideoFileName:     update.Message.Video.FileName,
		LastUpdateTime:    time.Now().Unix(),
	}
	b.userStates.Unlock()

	// Crear teclado de opciones
	keyboard := tu.InlineKeyboardMarkup(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🗜️ Comprimir Video").WithCallbackData("action_compress"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⚙️ Solo Enviar/Convertir").WithCallbackData("action_convert_only"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("❌ Cancelar").WithCallbackData("cancel"),
		),
	)

	b.telegoBot.SendMessage(tu.Message(tu.ID(chatID), "Video recibido. ¿Qué quieres hacer?").WithReplyMarkup(keyboard))
}

// callbackHandler maneja los callbacks de los botones en línea
func (b *Bot) callbackHandler(update telego.Update) {
	callback := update.CallbackQuery
	chatID := callback.Message.Chat.ID
	data := callback.Data

	b.userStates.RLock()
	userData, ok := b.userStates.data[chatID]
	b.userStates.RUnlock()

	if !ok {
		b.telegoBot.AnswerCallbackQuery(tu.AnswerCallbackQuery(callback.ID).WithText("Esta operación ha expirado."))
		return
	}
	
	// Lógica para manejar las diferentes acciones de los botones...
	switch data {
	case "cancel":
		b.userStates.Lock()
		userData.State = "cancelled"
		b.userStates.Unlock()
		b.telegoBot.EditMessageText(tu.EditMessageText(tu.ID(chatID), callback.Message.ID, "Operación cancelada."))
		b.cleanup(chatID)
	}

	b.telegoBot.AnswerCallbackQuery(tu.AnswerCallbackQuery(callback.ID))
}

// textHandler maneja los mensajes de texto del usuario
func (b *Bot) textHandler(update telego.Update) {
	// Lógica para manejar el comando /start o un nuevo nombre de archivo
	if update.Message.Text == "/start" {
		b.startHandler(update)
	}
}

// thumbnailHandler maneja las fotos enviadas como miniaturas
func (b *Bot) thumbnailHandler(update telego.Update) {
	// Lógica para manejar la miniatura
}

// userHasActiveProcess verifica si un usuario ya tiene un proceso activo
func (b *Bot) userHasActiveProcess(chatID int64) bool {
	b.userStates.RLock()
	defer b.userStates.RUnlock()
	_, ok := b.userStates.data[chatID]
	return ok
}

// cleanup elimina los archivos temporales y el estado del usuario
func (b *Bot) cleanup(chatID int64) {
	b.userStates.Lock()
	userData, ok := b.userStates.data[chatID]
	if !ok {
		b.userStates.Unlock()
		return
	}

	// Elimina archivos si existen
	if userData.DownloadPath != "" {
		os.Remove(userData.DownloadPath)
	}
	if userData.FinalPath != "" {
		os.Remove(userData.FinalPath)
	}
	if userData.ThumbnailPath != "" {
		os.Remove(userData.ThumbnailPath)
	}

	delete(b.userStates.data, chatID)
	b.userStates.Unlock()
	log.Printf("Datos del usuario %d limpiados.", chatID)
}


