package storage

import "time"

type Event struct {
	ID           string    // Уникальный идентификатор события (UUID)
	Title        string    // Заголовок события
	StartTime    time.Time // Дата и время начала события
	EndTime      time.Time // Дата и время окончания события
	Description  string    // Описание события (опционально)
	UserID       string    // ID пользователя, владельца события
	NotifyBefore string    // Время до уведомления (INTERVAL как строка)
	CreatedAt    time.Time // Время создания записи
}
