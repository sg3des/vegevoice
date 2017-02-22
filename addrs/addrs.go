package addrs

import (
	"bufio"
	"log"
	"os"
	"path"
	"strings"
)

var addrs []string
var maxReturnedItems int

func ReadUrls(dir string) {
	addrsfile := path.Join(dir, "addrs.list")

	f, err := os.Open(addrsfile)
	if err != nil {
		log.Println(err)
		return
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		addrs = append(addrs, scanner.Text())
	}
}

func GetAddrs(substr string) []string {
	if len(addrs) == 0 {
		return nil
	}

	var chosen []string
	for _, addr := range addrs {
		if strings.HasPrefix(addr, substr) {
			chosen = append(chosen, addr)
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

func ChosenContains(chosen []string, addr string) bool {
	for _, c := range chosen {
		if c == addr {
			return true
		}
	}
	return false
}

func SetMaxItems(n int) {
	maxReturnedItems = n
}
