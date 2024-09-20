package service1

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Interface interface {
	UserGet(id int) (*User, error)
}

var (
	pathUserGet = "/user/%d"
)

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type Service1 struct {
	client *http.Client
	domain string
}

func New(domain string, client *http.Client) Service1 {
	if nil == client {
		client = &http.Client{}
	}

	return Service1{
		client: client,
		domain: domain,
	}
}

func (s Service1) UserGet(id int) (*User, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(fmt.Sprintf("%s%s", s.domain, pathUserGet), id), nil)
	if nil != err {
		return nil, err
	}

	res, err := s.client.Do(req)
	if nil != err {
		return nil, err
	}

	if http.StatusOK != res.StatusCode {
		return nil, fmt.Errorf("recieved status code %d", res.StatusCode)
	}

	defer res.Body.Close()

	var user User
	err = json.NewDecoder(res.Body).Decode(&user)
	if nil != err {
		return nil, fmt.Errorf("cannot decode response: %w", err)
	}

	return &user, nil
}