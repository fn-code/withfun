package main

import (
	"fmt"

	pusher "github.com/pusher/pusher-http-go"
)

type Messeging struct {
	Title    string      `json:"title"`
	Body     string      `json:"body"`
	Channel  string      `json:"channel"`
	Priority string      `json:"priority"`
	Id       string      `json:"id"`
	Customs  interface{} `json:"customs"`
}

func main() {

	cstm := map[string]string{
		"source":         "ESP002",
		"from":           "2",
		"status":         "2",
		"emergency_type": "2",
	}

	notifMsg := &Messeging{
		Title:    "Darurat Kesehatan",
		Body:     fmt.Sprintf("Segera tindak lanjuti laporan darurat Kesehatan dari %s", "Ludin Nento"),
		Channel:  "darurat",
		Priority: "high",
		Id:       "id",
		Customs:  cstm,
	}

	pusherClient := pusher.Client{
		AppID:   "APPID",
		Key:     "KEY",
		Secret:  "SECRET",
		Cluster: "ap1",
		Secure:  false,
	}

	pusherClient.Trigger("darurat_channel", "darurat_event", notifMsg)
}
