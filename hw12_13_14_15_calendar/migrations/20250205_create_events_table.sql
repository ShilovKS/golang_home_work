-- +goose Up
CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Уникальный идентификатор события
    title VARCHAR(255) NOT NULL,                   -- Заголовок события (короткий текст)
    start_time TIMESTAMP NOT NULL,                 -- Дата и время начала события
    end_time TIMESTAMP NOT NULL,                   -- Дата и время окончания события
    description TEXT,                              -- Описание события (длинный текст, опционально)
    user_id UUID NOT NULL,                         -- ID пользователя, владельца события
    notify_before INTERVAL,                        -- Время до уведомления (опционально)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Время создания записи
);

-- +goose Down
DROP TABLE IF EXISTS events;