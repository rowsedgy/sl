package main

import (
	"encoding/json"
	"os"

	"github.com/charmbracelet/bubbles/list"
)

// type connection struct {
// 	Name string `json:"name"`
// 	Data struct {
// 		User       string `json:"user"`
// 		Password   string `json:"password"`
// 		Pubauth    bool   `json:"pubauth"`
// 		Key        string `json:"key"`
// 		IP         string `json:"ip"`
// 		WebIP      string `json:"webip"`
// 		Tunnel     bool   `json:"tunnel"`
// 		TunnelHost string `json:"tunnelhost"`
// 	} `json:"data"`
// }

type connections struct {
	TunnelHosts map[string]TunnelHost `json:"tunnelhosts"`
	Hosts       map[string]Host       `json:"hosts"`
}

type TunnelHost struct {
	User     string `json:"user"`
	Password string `json:"password"`
	IP       string `json:"ip"`
}

type Host struct {
	User       string `json:"user"`
	Password   string `json:"password"`
	Pubauth    bool   `json:"pubauth"`
	Key        string `json:"key"`
	IP         string `json:"ip"`
	WebIP      string `json:"webip"`
	Tunnel     bool   `json:"tunnel"`
	TunnelHost string `json:"tunnelhost"`
}

func (c *cfg) generateList() (list.Model, error) {
	bytes, err := os.ReadFile(c.filepath)
	if err != nil {
		return list.Model{}, err
	}

	var connections connections

	if len(bytes) == 0 {
		return list.New(nil, list.NewDefaultDelegate(), 0, 0), nil
	}

	if err := json.Unmarshal(bytes, &connections); err != nil {
		return list.Model{}, err
	}

	c.connections = connections

	items := []list.Item{}

	for name, data := range connections.Hosts {
		items = append(items, Item{
			name:       name,
			ip:         data.IP,
			webip:      data.WebIP,
			pubauth:    data.Pubauth,
			user:       data.User,
			password:   data.Password,
			key:        data.Key,
			tunnel:     data.Tunnel,
			tunnelHost: data.TunnelHost,
		})
	}

	newList := list.New(items, list.NewDefaultDelegate(), 0, 0)

	return newList, nil
}
