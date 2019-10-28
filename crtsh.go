package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var newEntry bool
var subdomains []string
var latestsubdomains []string

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// TEST 
func SliceElementExist(a []string, x string) bool {
	for _, sub := range a {
		if x == sub {
			return true
		}
	}
	return false
}

func fetchSubDomains(domain string) {
	resp, err := http.Get("https://crt.sh/?q=%25." + domain)
	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	bodyText := string(body)
	r := regexp.MustCompile("<TD>.*" + domain)
	matched := r.FindAllString(bodyText, -1)
	for _, subdomain := range matched {
		subdomain = strings.Replace(subdomain, "<TD>", "", -1)
		if !SliceElementExist(subdomains, subdomain) {
			subdomains = append(subdomains, subdomain)
		}
	}

	if _, err := os.Stat(domain + ".subs"); err == nil {
		newEntry = false
	} else if os.IsNotExist(err) {
		newEntry = true
		os.Create(domain + ".subs")
	}
	latestsubdomains = subdomains
	f, err := os.OpenFile(domain+".subs", os.O_APPEND|os.O_WRONLY, 0644)
	check(err)
	if newEntry {
		for _, subdomain := range subdomains {
			f.WriteString(subdomain + "\n")
		}
		f.Close()
	} else {
		monitor(domain)
	}
	subdomains = subdomains[:0]
}

func monitor(domain string) {
	bot_api_key := "REPLACE"
	channel_name := "REPLACE"
	content, err := ioutil.ReadFile(domain + ".subs")
	check(err)
	lines := strings.Split(string(content), "\n")
	f, err := os.OpenFile(domain+".subs", os.O_APPEND|os.O_WRONLY, 0644)
	check(err)
	for _, x := range latestsubdomains {
		if !SliceElementExist(lines, x) {
			if x != "" {
				fmt.Println("New subdomain found : " + x + "\n")
				resp, err := http.Get("https://api.telegram.org/bot" + bot_api_key + "/sendMessage?chat_id=" + channel_name + "&text=New subdomain found : " + x)
				check(err)
				_ = resp
				f.WriteString(x + "\n")
			}
		}
	}
	f.Close()
}

func main() {
	target := os.Args[1]
	content, err := ioutil.ReadFile(target)
	check(err)
	targets := strings.Split(string(content), "\n")
	for _, domain := range targets {
		if domain != "" {
			fmt.Println("\n" + domain)
			fetchSubDomains(domain)
		}
	}
}
