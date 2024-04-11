package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

func sender(ch chan string, username string) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		msg := scanner.Text()
		if errInput := scanner.Err(); errInput != nil {
			log.Fatal("read:", errInput)
		}

		if strings.Compare(msg, "!exit") == 0 {
			os.Exit(0)
		}

		fullMsg := fmt.Sprintf("%s: %s", username, msg)

		ch <- fullMsg
	}

	close(ch)
}

func receiver(ch chan string, c *websocket.Conn) {
	for {
		_, msg, err := c.ReadMessage()
		// _, reader, err := c.NextReader()
		if err != nil {
			log.Fatal("read:", err)
		}

		// msg, err := io.ReadAll(reader)
		if err != nil {
			log.Fatal("read:", err)
		}

		log.Println(string(msg))

		ch <- string(msg)

		// close(ch)
	}
}

func main() {
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	var username string
	fmt.Print("Enter your user: ")
	_, errr := fmt.Scanln(&username)
	if errr != nil {
		log.Fatal("read:", err)
	}

	sendMsg := make(chan string)
	receiveMsg := make(chan string)

	go sender(sendMsg, username)
	go receiver(receiveMsg, c)

	for {
		select {
		case msg := <-sendMsg:
			if msg != "" {
				err := c.WriteMessage(
					websocket.TextMessage,
					[]byte(msg),
				)

				if err != nil {
					log.Fatal("write:", err)
				}
			}
		case msg := <-receiveMsg:
			if msg != "" {
				fmt.Println(msg)
			}
		}
	}
}
