package main

import (
	"encoding/json"
	"os"

	"github.com/charmbracelet/bubbles/list"
)

type connection struct {
	Name string `json:"name"`
	Data struct {
		User     string `json:"user"`
		Password string `json:"password"`
		IP       string `json:"IP"`
	} `json:"data"`
}

func generateList(filePath string) (list.Model, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return list.Model{}, err
	}

	var connections []connection

	if err := json.Unmarshal(bytes, &connections); err != nil {
		return list.Model{}, err
	}

	items := []list.Item{}

	for _, conn := range connections {
		items = append(items, Item{
			name:     conn.Name,
			ip:       conn.Data.IP,
			user:     conn.Data.User,
			password: conn.Data.Password,
		})
	}

	newList := list.New(items, list.NewDefaultDelegate(), 0, 0)

	return newList, nil

}
