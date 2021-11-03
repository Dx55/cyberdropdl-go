package main

import (
	"bufio"
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
var tr = &http.Transport{
	MaxIdleConns:        30,
	MaxIdleConnsPerHost: 30,
	MaxConnsPerHost:     30,
}
var netClient = &http.Client{Transport: tr}

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
		// Seems painful here too idk
		if strings.Contains(link, "href") == true {
			reSpecial := regexp.MustCompile(`(https:\/\/fs([\w\-_]+(?:(?:\.[\w\-_]+)+))([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?)`)
			link := reSpecial.FindString(link)

			if iter <= len(links) {
				// Starting request for file
				fmt.Println(fmt.Sprintf("Starting request for %s(%d)", cyberlink[23:], iter))

				resp, err := netClient.Get(link)
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
	}
	time.Sleep(2)
}

func download(link string, iter int, resp *http.Response) {
	formatString := fmt.Sprint(link)
	reFormat := regexp.MustCompile(`\.jpg|.jpeg|.png|.gif|.webp|.mp4|.webm|.mov|.mkv`)
	fileFormat := reFormat.FindString(formatString)

	fmt.Println("Creating file")
	// Creating file
	file, err := os.Create(fmt.Sprintf("%s/%s(%d)%s", cyberlink[23:], cyberlink[23:], iter, fileFormat))
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
	fmt.Println(fmt.Sprintf("Downloaded %s(%d)%s", cyberlink[23:], iter, fileFormat))
	fmt.Println("---------------------------------------------------------")
}

func folderVerifier(cyberlink string) {
	if _, err := os.Stat(fmt.Sprintf("%s", cyberlink[23:])); os.IsNotExist(err) {
		if err != nil {
			fmt.Println("Directory doesn't exist")
			fmt.Println("Creating a new one...")
			_ = os.Mkdir(fmt.Sprintf("%s", cyberlink[23:]), 0755)
			request(cyberlink, retrieveLinks(cyberlink))
		}
	} else {
		// Starting the functions
		request(cyberlink, retrieveLinks(cyberlink))
	}
}

func main() {
	// Taking input from args
	if os.Args[1] == "-m" {
		file, err := os.Open(os.Args[2])
		if err != nil {
			fmt.Println("Unable to read file")
			fmt.Println(err)
		}
		defer file.Close()

		line_scan := bufio.NewScanner(file)

		for line_scan.Scan() {
			cyberlink = line_scan.Text()
			fmt.Println(cyberlink)
			folderVerifier(cyberlink)
		}
	} else {
		cyberlink = os.Args[1]

		folderVerifier(cyberlink)
	}
}
