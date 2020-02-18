package main

import (
	"net/http"

	"github.com/robfig/cron"

	"github.com/PuerkitoBio/goquery"
	"github.com/getlantern/systray"
)

type state struct {
	Stars string
	Cron  *cron.Cron
}

func main() {
	s := &state{}
	systray.Run(s.onReady, s.onExit)
}

func (s *state) onReady() {
	s.updateStars()
	s.Cron = cron.New()
	s.Cron.AddFunc("@every 10m", s.updateStars)
	s.Cron.Start()
}
func (s *state) onExit() {
	s.Cron.Stop()
}

func (s *state) updateStars() {
	url := "https://github.com/vuejs/vue/stargazers"
	res, err := http.Get(url)

	if err != nil {
		return
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return
	}

	stars := doc.Find("a[href='/vuejs/vue/stargazers']>.Counter").Text()

	systray.SetTitle("Vue stars " + stars)

}
