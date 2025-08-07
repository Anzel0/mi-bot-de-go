package bot

import (
	"fmt"
	"time"
)

// formatSize convierte un tamaño en bytes a un formato legible (KB, MB, GB).
func formatSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d Bytes", size)
	}
	if size < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(size)/1024)
	}
	if size < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(size)/1024/1024)
	}
	return fmt.Sprintf("%.2f GB", float64(size)/1024/1024/1024)
}

// humanReadableTime convierte segundos a un formato de tiempo H:M:S.
func humanReadableTime(seconds int) string {
	if seconds <= 0 {
		return "00:00:00"
	}
	duration := time.Duration(seconds) * time.Second
	h := duration / time.Hour
	m := (duration % time.Hour) / time.Minute
	s := (duration % time.Minute) / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

// getProgressBar genera una barra de progreso en texto.
func getProgressBar(percentage float64) string {
	completedBlocks := int(percentage / 10)
	progress := ""
	for i := 0; i < completedBlocks; i++ {
		progress += "█"
	}
	for i := 0; i < 10-completedBlocks; i++ {
		progress += " "
	}
	return progress
}

