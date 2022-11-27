package main

import (
	"errors"
	"log"
	"time"

	"github.com/goguardian/pusher-ws-go"
)

const (
	bitstampAppKey = "5f1151a6d8354d05b649"
)

type Customs struct {
	Source string `json:"source"`
	From   string `json:"from"`
	Status string `json:"status"`
}

type Data struct {
	Title    string  `json:"title"`
	Body     string  `json:"body"`
	Priority string  `json:"priority"`
	ID       string  `json:"id"`
	Customs  Customs `json:"customs"`
}

var ErrPusherConnecting error = errors.New("error connect to pusher server")
var ErrPusherSubscribing error = errors.New("error subscribe to pusher channel")

func main() {

	// create error channel and start printing errors
	errChan := make(chan error)
	successChan := make(chan struct{})
	closeChan := make(chan chan struct{})

	retryDuration := 3 * time.Second
	var tm *time.Timer

	go func() {
		for {
			select {
			case err := <-errChan:
				log.Println("Error: ", err)

				if err == ErrPusherConnecting || err == ErrPusherSubscribing {
					tm = time.NewTimer(retryDuration)

					select {
					case <-tm.C:
						tm.Stop()
						log.Printf("trying to reconnect \n")
						go run(successChan, closeChan, errChan)
						break
					}

				} else {

					log.Println("connection lost")
					stop := make(chan struct{})
					closeChan <- stop
					<-stop

					log.Println("done connection lost")
					go run(successChan, closeChan, errChan)
					break
				}

			}

		}
	}()

	go run(successChan, closeChan, errChan)
	select {}

}

func run(successChan chan<- struct{}, closedChan <-chan chan struct{}, errChan chan error) {
	// instantiate Pusher client with the error channel;
	// commented options are the defaults
	pusherClient := &pusher.Client{
		// Insecure: false,
		Cluster: "ap1",
		Errors:  errChan,
	}

	// connect to the Bitstamp Pusher app
	err := pusherClient.Connect(bitstampAppKey)
	if err != nil {

		errChan <- ErrPusherConnecting
		log.Print("err 1 ")
		return
	}

	// subscribe to the BTC/USD order book
	usdOrderBook, err := pusherClient.Subscribe("darurat_channel")
	if err != nil {

		errChan <- ErrPusherSubscribing
		log.Print("err 2 ")
		return
	}

	// bind to data events on each order channel
	usdOrderData := usdOrderBook.Bind("darurat_event")

	// CHK:
	for {
		select {
		case usdOrder := <-usdOrderData:

			dd := &Data{}
			err := pusher.UnmarshalDataString(usdOrder, dd)
			if err != nil {
				log.Printf("Error unmarshaling : %v\n", err)
			}

			log.Println(dd)
		case stop := <-closedChan:
			pusherClient.Disconnect()
			close(stop)
			log.Println("done")
			return
		}
	}
}
