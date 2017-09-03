package urlstorage

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
)

var (
	wd               string
	urls             []string
	maxReturnedItems = 10
)

func SetMaxItems(n int) {
	maxReturnedItems = n
}

func Initialize(dir string) {
	wd = dir

	readUserURLs()
	readListURLs()
}

func readListURLs() {
	filename := path.Join(wd, "list_urls.txt")

	f, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}
}

var userURLs struct {
	PinnedTabs []string
	Bookmarks  []string
	History    map[string]int
}

func readUserURLs() {
	filename := path.Join(wd, "user_urls.toml")

	_, err := toml.DecodeFile(filename, &userURLs)
	if err != nil {
		log.Println(err)
		return
	}
}

func saveUserURLs() {
	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	err := enc.Encode(userURLs)
	if err != nil {
		log.Println(err)
		return
	}

	filename := path.Join(wd, "user_urls.toml")
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = f.Write(buf.Bytes())
	if err != nil {
		log.Println(err)
		return
	}
}

func GetURLs(substr string) []string {
	if len(urls) == 0 {
		return nil
	}

	var chosen []string
	for _, u := range urls {
		if strings.HasPrefix(u, substr) {
			chosen = append(chosen, u)
			if len(chosen) >= maxReturnedItems {
				return chosen
			}
		}
	}

	// for _, addr := range addrs {
	// 	if strings.Contains(addr, substr) && !ChosenContains(chosen, addr) {
	// 		chosen = append(chosen, addr)
	// 		if len(chosen) >= maxReturnedItems {
	// 			return chosen
	// 		}
	// 	}
	// }

	// substr = strings.Replace(substr, ".", "", -1)
	// for _, addr := range addrs {
	// 	if strings.ContainsAny(addr, substr) && !ChosenContains(chosen, addr) {
	// 		chosen = append(chosen, addr)
	// 		if len(chosen) >= maxReturnedItems {
	// 			return chosen
	// 		}
	// 	}
	// }

	return chosen
}

func ChosenContains(chosen []string, u string) bool {
	for _, c := range chosen {
		if c == u {
			return true
		}
	}
	return false
}

func GetPinnedTabs() []string {
	return userURLs.PinnedTabs
}

func AddPinnedTab(u string) (n int) {
	n = len(userURLs.PinnedTabs)

	userURLs.PinnedTabs = append(userURLs.PinnedTabs, u)

	saveUserURLs()
	return
}

func DelPinnedTab(n int) {
	userURLs.PinnedTabs = append(userURLs.PinnedTabs[:n], userURLs.PinnedTabs[n+1:]...)
	saveUserURLs()
}

func MovePinnedTab(n0, n1 int) {
	userURLs.PinnedTabs[n0], userURLs.PinnedTabs[n1] = userURLs.PinnedTabs[n1], userURLs.PinnedTabs[n0]
	saveUserURLs()
}
