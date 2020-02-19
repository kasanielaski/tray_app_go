package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/robfig/cron"

	"github.com/getlantern/systray"
)

type state struct {
	Stars string
	Cron  *cron.Cron
}

type Result struct {
	Id     int    `json:"_id"`
	Status string `json: "status"`
	Rubric string `json: "rubric"`
}

type Metadata struct {
	Resultset Resultset `json: "resultset"`
}

type Resultset struct {
	Limit  int `json: "limit"`
	Offset int `json: "offset"`
	Count  int `json: "count"`
}

type APIResponse struct {
	Result   []Result `json: "result"`
	Metadata Metadata `json: "metadata"`
}

func main() {
	s := &state{}
	systray.Run(s.onReady, s.onExit)
}

func getCount(body []byte) (*APIResponse, error) {
	var c = new(APIResponse)
	err := json.Unmarshal(body, &c)
	if err != nil {
		log.Fatal(err)
	}
	return c, err
}

func (s *state) onReady() {
	s.updateFlats()
	s.Cron = cron.New()
	s.Cron.AddFunc("@daily", s.updateFlats)
	s.Cron.Start()
}
func (s *state) onExit() {
	s.Cron.Stop()
}

func (s *state) updateFlats() {
	url := "https://api.n1.ru/api/v1/offers/?limit=25&offset=0&sort=-billing_weight%2C-order_date%2C-creation_date&query%5B0%5D%5Bdeal_type%5D=sell&query%5B0%5D%5Brubric%5D=flats&filter%5Bcity_id%5D=89026%2C89030%2C89027%2C2355612%2C89946%2C89966%2C89958%2C89978%2C89963%2C89996%2C90006%2C89955%2C89967%2C89994%2C90027%2C90008%2C89999%2C89982%2C90004%2C89989%2C89028%2C90014%2C89979%2C89973%2C90015%2C89952%2C89975%2C89984%2C90019%2C89956%2C89912%2C89992%2C89961%2C89985%2C90007%2C90013%2C89968%2C89969%2C90024%2C89948%2C89997%2C89976%2C89962%2C89551%2C89540%2C89951%2C89988%2C89945%2C89974%2C89953%2C90011%2C90012%2C89892%2C89936%2C89548%2C89981%2C89957%2C89954%2C89980%2C89585%2C89983%2C89516%2C89959%2C90020%2C90003%2C89584%2C89919%2C89572%2C89986%2C2355854%2C2355857%2C89313%2C89995&filter%5Bregion_id%5D=1054&filter_or%5Baddresses%5D%5B0%5D%5Bstreet_id%5D=865620&filter_or%5Baddresses%5D%5B0%5D%5Bhouse_number%5D=11&status=published"
	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic(res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	c, err := getCount(body)

	fmt.Println(c)
}
