package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type User struct {
	ID       int    `json:"-"`
	Name     string `json:"-"`
	Username string `json:"-"`
	Email    string
	Phone    string `json:"-"`
	Password string `json:"-"`
	Address  string `json:"-"`
}

type DomainStat map[string]int

//nolint
func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %s", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	reader := bufio.NewReader(r)
	i := 0
	var user User
	var lineBytes []byte
	ok := true
	for ok {
		lineBytes, err = reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				// чтобы прочитать последнюю строку
				ok = false
				err = nil
			} else {
				break
			}
		}

		err = user.UnmarshalJSON(lineBytes)
		if err == nil {
			result[i] = user
		}
		i++
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	var fullDomain string
	for _, user := range u {
		if user.Email == "" {
			// считаем, что данные кончились
			break
		}
		fullDomain = strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])

		if strings.HasSuffix(fullDomain, domain) {
			result[fullDomain]++
		}
	}
	return result, nil
}
