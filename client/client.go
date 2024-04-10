package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

func sender(c *websocket.Conn, username string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		if errInput := scanner.Err(); errInput != nil {
			log.Fatal("read:", errInput)
		}

		if strings.Compare(msg, "!exit") == 0 {
			break
		}

		err := c.WriteMessage(
			websocket.TextMessage,
			[]byte(username+": "+string(msg)),
		)
		if err != nil {
			log.Fatal("write:", err)
		}
	}
}

func receiver(c *websocket.Conn, username string) {
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Fatal("read:", err)
		}
		if !strings.Contains(string(msg), username) {
			fmt.Println(msg)
		}
	}
}

func main() {
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	var username string
	fmt.Print("Enter your user: ")
	_, errr := fmt.Scanln(&username)
	if errr != nil {
		log.Fatal("read:", err)
	}

	go receiver(c, username)
	go sender(c, username)
	go func() {
		for {
			_, msg, err := c.ReadMessage()
			println(string(msg))
			if err != nil {
				log.Fatal("read:", err)
			}
			if !strings.Contains(string(msg), username) {
				fmt.Println(msg)
			}
		}
	}()
}
