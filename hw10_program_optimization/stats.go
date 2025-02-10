package hw10programoptimization

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// User представляет пользователя.
// Обратите внимание: имена полей (с учётом регистра) соответствуют входящему JSON.
type User struct {
	ID       int    `json:"Id"`
	Name     string `json:"Name"`
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Phone    string `json:"Phone"`
	Password string `json:"Password"`
	Address  string `json:"Address"`
}

// DomainStat – статистика по доменам: ключ – доменная часть email, значение – количество.
type DomainStat map[string]int

// GetDomainStat читает данные из r (каждая строка – отдельный JSON-объект) и подсчитывает,
// сколько раз встречается доменная часть email, если она заканчивается на "."+domain (без учёта регистра).
//
// В этой реализации мы используем json.Decoder для последовательного декодирования объектов,
// что позволяет обрабатывать данные «на лету» без накопления всех пользователей в памяти.
func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	stat := make(DomainStat)
	dec := json.NewDecoder(r)
	targetSuffix := "." + strings.ToLower(domain)

	for {
		var user User
		err := dec.Decode(&user)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("get users error: %w", err)
		}
		// Обрабатываем email: находим символ '@'
		email := user.Email
		at := strings.LastIndex(email, "@")
		if at < 0 || at >= len(email)-1 {
			continue
		}
		domPart := email[at+1:]
		// Если доменная часть (в нижнем регистре) оканчивается на "."+domain, увеличиваем счётчик.
		if strings.HasSuffix(strings.ToLower(domPart), targetSuffix) {
			key := strings.ToLower(domPart)
			stat[key]++
		}
	}

	return stat, nil
}
