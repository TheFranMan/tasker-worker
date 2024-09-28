package service3

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Interface interface {
	UserGet(email string) (*User, error)
	DeleteUser(email string) error
	UpdateUser(oldEmail, newEmail string) error
}

const (
	pathUser = "/user"
)

type Service3 struct {
	client *http.Client
	domain string
}

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func New(domain string, client *http.Client) Service3 {
	if nil == client {
		client = &http.Client{}
	}

	return Service3{
		domain: domain,
		client: client,
	}
}

func (s Service3) UserGet(email string) (*User, error) {
	req, err := http.NewRequest(http.MethodGet, s.domain+pathUser, strings.NewReader(fmt.Sprintf(`{"email": "%s"}`, email)))
	if nil != err {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if nil != err {
		return nil, err
	}

	if http.StatusOK != resp.StatusCode {
		return nil, fmt.Errorf("recieved status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	var user User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if nil != err {
		return nil, err
	}

	return &user, nil
}

func (s Service3) DeleteUser(email string) error {
	req, err := http.NewRequest(http.MethodDelete, s.domain+pathUser, strings.NewReader(fmt.Sprintf(`{"email": "%s"}`, email)))
	if nil != err {
		return err
	}

	resp, err := s.client.Do(req)
	if nil != err {
		return err
	}

	if http.StatusOK != resp.StatusCode {
		return fmt.Errorf("recieved status code %d", resp.StatusCode)
	}

	return nil
}

func (s Service3) UpdateUser(oldEmail, newEmail string) error {
	req, err := http.NewRequest(http.MethodPost, s.domain+pathUser, strings.NewReader(fmt.Sprintf(`{"old_email": "%s", "new_email": "%s"}`, oldEmail, newEmail)))
	if nil != err {
		return err
	}

	resp, err := s.client.Do(req)
	if nil != err {
		return err
	}

	if http.StatusOK != resp.StatusCode {
		return fmt.Errorf("recieved status code %d", resp.StatusCode)
	}

	return nil
}
