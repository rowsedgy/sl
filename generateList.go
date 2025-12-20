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
		IP       string `json:"ip"`
		WebIP    string `json:"webip"`
	} `json:"data"`
}

func (c *cfg) generateList() (list.Model, error) {
	bytes, err := os.ReadFile(c.filepath)
	if err != nil {
		return list.Model{}, err
	}

	var connections []connection

	if len(bytes) == 0 {
		return list.New(nil, list.NewDefaultDelegate(), 0, 0), nil
	}

	if err := json.Unmarshal(bytes, &connections); err != nil {
		return list.Model{}, err
	}

	items := []list.Item{}

	for _, conn := range connections {
		items = append(items, Item{
			name:     conn.Name,
			ip:       conn.Data.IP,
			webip:    conn.Data.WebIP,
			user:     conn.Data.User,
			password: conn.Data.Password,
		})
	}

	newList := list.New(items, list.NewDefaultDelegate(), 0, 0)

	return newList, nil

}
