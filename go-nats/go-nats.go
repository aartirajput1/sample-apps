package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	nats "github.com/nats-io/go-nats"
)

func main() {
	http.HandleFunc("/echo", func(w http.ResponseWriter, req *http.Request) {
		nc, err := nats.Connect(os.Getenv("NATS_URI"))
		if err != nil {
			http.Error(w,
				fmt.Sprintf("Error connecting to NATS at %q: %s\n", os.Getenv("NATS_URI"), err),
				http.StatusInternalServerError)
			return
		}
		defer nc.Close()

		msg := req.FormValue("msg")

		sub, err := nc.Subscribe("test", func(m *nats.Msg) {
			nc.Publish(m.Reply, []byte(msg))
		})
		if err != nil {
			http.Error(w,
				fmt.Sprintf("Can't subscribe to 'test' subject: %s\n", err),
				http.StatusInternalServerError)
			return
		}
		defer sub.Unsubscribe()

		reply, err := nc.Request("test", []byte(msg), time.Second)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("Request failed: %s\n", err),
				http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, string(reply.Data)+"\n")
	})

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("No PORT defined")
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Error starting HTTP server: ", err)
	}
}
