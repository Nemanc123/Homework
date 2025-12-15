package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"
)

//go:generate easyjson -all
type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomains(r, domain)
}

//type users [100_000]User

func countDomains(r io.Reader, domain string) (result DomainStat, err error) {
	scanner := bufio.NewScanner(r)
	var user User
	result = make(DomainStat)
	for scanner.Scan() {
		line := scanner.Bytes()
		if err = user.UnmarshalJSON(line); err != nil {
			return
		}
		matched := strings.Contains(user.Email, "."+domain)
		if matched {
			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}
	return
}
