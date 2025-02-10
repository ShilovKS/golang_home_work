package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.Config{
	EscapeHTML:             true,
	SortMapKeys:            true,
	ValidateJsonRawMessage: true,
	CaseSensitive:          false,
}.Froze()

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	targetSuffix := "." + strings.ToLower(domain)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var user User
		if err := json.Unmarshal(line, &user); err != nil {
			return nil, fmt.Errorf("get users error: %w", err)
		}

		email := user.Email
		at := strings.LastIndex(email, "@")
		if at < 0 || at >= len(email)-1 {
			continue
		}
		domPart := email[at+1:]
		if strings.HasSuffix(strings.ToLower(domPart), targetSuffix) {
			key := strings.ToLower(domPart)
			result[key]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
