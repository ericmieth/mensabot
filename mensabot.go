package main

import (
	"encoding/json"
	"github.com/42wim/matterbridge/matterhook"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"os"
	"time"
)

// contains the configuration items
type config struct {
	WebhookPath        string
	MensaID            string
	MattermostUsername string
	MattermostChannel  string
	MattermostIconURL  string
}

// returns a formatted date (e.g. like "2006-01-02")
func formatDateLink(date time.Time) string {
	return date.Format("2006") +
		"-" +
		date.Format("01") +
		"-" +
		date.Format("02")
}

// returns a formatted date with german months (e.g. like "02. Januar 2006")
func formatDate(date time.Time) string {
	monthMap := map[string]string{
		"January":   "Januar",
		"February":  "Februar",
		"March":     "März",
		"April":     "April",
		"May":       "Mai",
		"June":      "Juni",
		"July:":     "Juli",
		"August":    "August",
		"September": "September",
		"October":   "Oktober",
		"November":  "November",
		"December":  "Dezember",
	}
	return date.Format("02.") +
		" " +
		monthMap[date.Format("January")] +
		" " +
		date.Format("2006")
}

// check, if the next TextToken is a location (e.g. like "Mensa am Park")
func locationNext(t html.Token, mensaID string) bool {
	for _, a := range t.Attr {
		if a.Key == "value" {
			if a.Val == mensaID {
				return true
			}
		}
	}
	return false
}

// check, if the next TextToken is a section (e.g. like "Pizza")
func sectionNext(t html.Token) bool {
	for _, a := range t.Attr {
		if a.Key == "class" {
			if a.Val == "menu_title" {
				return true
			}
		}
	}
	return false
}

// check, if the next TextToken is a dish (e.g. like "Pizza California")
func dishNext(t html.Token) bool {
	for _, a := range t.Attr {
		if a.Key == "class" {
			if a.Val == "menu_name1" || a.Val == "menu_name2" {
				return true
			}
		}
	}
	return false
}

// check, if the next TextToken is a dish (e.g. like "1.50€…")
func priceNext(t html.Token) bool {
	for _, a := range t.Attr {
		if a.Key == "class" {
			if a.Val == "menu_price" {
				return true
			}
		}
	}
	return false
}

// parse the website and prepare the string with the menu in markdown
func parse(url string) string {

	resp, err := http.Get(url)
	if err != nil {
		log.Print(err)
		return ""
	}
	b := resp.Body
	defer b.Close()

	z := html.NewTokenizer(b)

	var sectionIsNext bool
	var dishIsNext bool
	var priceIsNext bool
	var newListItem bool
	var menuData string

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return menuData

		case tt == html.StartTagToken:
			t := z.Token()

			isDiv := t.Data == "div"
			if !isDiv {
				continue
			}
			dishIsNext = dishNext(t)
			priceIsNext = priceNext(t)
			sectionIsNext = sectionNext(t)

		case tt == html.TextToken:
			if sectionIsNext {
				t := z.Token()
				newListItem = true
				menuData = menuData + "\n**" + t.Data + "**\n"
			}
			if dishIsNext {
				t := z.Token()
				if newListItem {
					menuData = menuData + "* "
					newListItem = false
				}
				menuData = menuData + t.Data + ", "
			}
			if priceIsNext {
				t := z.Token()
				newListItem = true
				menuData = menuData + "**" + t.Data[:4] + "€ **\n"
			}
		}
	}
}

// returns the name of a canteen
func getCanteenName(url string, mensaID string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Print(err)
		return ""
	}
	b := resp.Body
	defer b.Close()

	z := html.NewTokenizer(b)

	var locationIsNext bool
	var mensaName string

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return mensaName

		case tt == html.StartTagToken:
			t := z.Token()

			isOption := t.Data == "option"
			if !isOption {
				continue
			}
			locationIsNext = locationNext(t, mensaID)

		case tt == html.TextToken:
			if locationIsNext {
				t := z.Token()
				return t.Data
			}
		}
	}
}

func main() {

	// read configuration file
	file, err := os.Open("config.json")
	decoder := json.NewDecoder(file)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	configuration := config{}
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatal(err)
	}

	// prepare url to parse
	now := time.Now()
	todayLink := formatDateLink(now)
	mensaURL := "https://www.studentenwerk-leipzig.de/mensa/menu?date=" + todayLink + "&location=" + configuration.MensaID
	menu := parse(mensaURL)

	mensaName := getCanteenName(mensaURL, configuration.MensaID)

	// prepare mattermost message
	today := formatDate(now)
	m := matterhook.New(configuration.WebhookPath,
		matterhook.Config{DisableServer: true})
	msg := matterhook.OMessage{}
	msg.UserName = configuration.MattermostUsername
	msg.Channel = configuration.MattermostChannel
	msg.IconURL = configuration.MattermostIconURL
	msg.Text = "Speiseplan am " + today + " (" + mensaName + ")\n" + menu
	if menu != "" {
		m.Send(msg)
	}
}
