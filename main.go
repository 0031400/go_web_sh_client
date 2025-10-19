package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/gorilla/websocket"
	"golang.org/x/term"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	server := flag.String("s", "ws://127.0.0.1:3000/ws", "the server websocket path")
	auth := flag.String("a", "", "the authorization header")
	flag.Parse()
	var err error
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Panicln(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	c, _, err := websocket.DefaultDialer.Dial(*server, map[string][]string{"authorization": {*auth}})
	if err != nil {
		log.Panicln(err)
	}
	defer c.Close()
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil && err != io.EOF {
				log.Println(err)
				return
			}
			c.WriteMessage(websocket.BinaryMessage, buf[:n])
		}
	}()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if mt == websocket.BinaryMessage {
			os.Stdout.Write(message)
		}
	}
}
