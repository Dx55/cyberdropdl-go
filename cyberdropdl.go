package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

var resp *http.Response
var cyberlink string

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func retrieveLinks(cyberlink string) (links []string) {
	// Request to cyberdrop
	cyberreq, err := http.Get(cyberlink)
	if err != nil {
		print(err)
	}

	// Declaring body from request
	body, err := ioutil.ReadAll(cyberreq.Body)
	defer cyberreq.Body.Close()
	if err != nil {
		print(err)
	}

	// Regex for finding links (painful but it works)
	re := regexp.MustCompile(`href="(https:\/\/fs([\w\-_]+(?:(?:\.[\w\-_]+)+))([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?).([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])`)
	rawData := re.FindAllString(string(body), -1)

	data := fmt.Sprint(rawData)
	reLinks := regexp.MustCompile(`(https:\/\/fs([\w\-_]+(?:(?:\.[\w\-_]+)+))([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?).([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])[:>.jpg|.jpeg|.png|.gif|.webp|.mp4|.webm|.mov|.mkv]`)
	links = reLinks.FindAllString(data, -1)

	return links
}

func request(cyberlink string, links []string) {
	links = removeDuplicateStr(links)
	for iter, link := range links {
		fmt.Println("---------------------------------------------------------")
		// Seems painful here too idk
		if strings.Contains(link, "href") == true {
			reSpecial := regexp.MustCompile(`(https:\/\/fs([\w\-_]+(?:(?:\.[\w\-_]+)+))([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?)`)
			link := reSpecial.FindString(link)

			// Starting request for file
			fmt.Println(fmt.Sprintf("Starting request for %s(%d)", cyberlink[24:], iter))
			resp, err := http.Get(link)
			if err != nil {
				fmt.Println("Request error")
				fmt.Println(err)
				break
			}
			defer resp.Body.Close()

			download(link, iter, resp)
		} else {
			download(link, iter, resp)
		}
	}
	time.Sleep(3)
}

func download(link string, iter int, resp *http.Response) {
	formatString := fmt.Sprint(link)
	reFormat := regexp.MustCompile(`\.jpg|.jpeg|.png|.gif|.webp|.mp4|.webm|.mov|.mkv`)
	fileFormat := reFormat.FindString(formatString)

	fmt.Println("Creating file")
	// Creating file
	file, err := os.Create(fmt.Sprintf("%s/%s(%d)%s", cyberlink[24:], cyberlink[24:], iter, fileFormat))
	if err != nil {
		fmt.Println("Error with creating the image")
		fmt.Println(err)
	}
	defer file.Close()

	fmt.Println("Copying data into image (may take some time)")
	// Copy data from HTTP response to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Error when copying data from HTTP response to file")
		fmt.Println(err)
	}
	fmt.Println(fmt.Sprintf("Downloaded %s(%d)%s", cyberlink[24:], iter, fileFormat))
	fmt.Println("---------------------------------------------------------")
	time.Sleep(1)
}

func main() {
	// Taking input from args
	cyberlink = os.Args[1]
	// If directory doesn't exist, create one
	if _, err := os.Stat(fmt.Sprintf("%s", cyberlink[24:])); os.IsNotExist(err) {
		if err != nil {
			fmt.Println("Directory doesn't exist")
			fmt.Println("Creating a new one...")
			_ = os.Mkdir(fmt.Sprintf("%s", cyberlink[24:]), 0755)
			request(cyberlink, retrieveLinks(cyberlink))
		}
	} else {
		// Starting the functions
		request(cyberlink, retrieveLinks(cyberlink))
	}
}
