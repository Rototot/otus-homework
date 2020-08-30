package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %s", err)
	}

	return countDomains(u, domain)
}

type users []User

func getUsers(r io.Reader) (users, error) {
	var result users
	content := bufio.NewReader(r)
	for {
		line, _, err := content.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		var user User
		if err := json.Unmarshal(line, &user); err != nil {
			return nil, err
		}

		result = append(result, user)
	}

	return result, nil
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	domain = strings.ToLower(domain)

	for _, user := range u {
		email := strings.ToLower(user.Email)
		index := strings.Index(email, "@")
		if index >= 0 && strings.Contains(email, "."+domain) {
			emailSuffix := email[index+1:]
			result[emailSuffix]++
		}
	}

	return result, nil
}
