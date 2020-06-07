package api

import (
	"fmt"
	"log"
  "net/http"
  "os"
  "strings"

	"github.com/PuerkitoBio/goquery"
)

type match struct {
	date, team1, team2, result, csMap string
	matchId string
}

var (
	teams []string
	matches []match
	client = &http.Client{}
	userAgent = os.Getenv("USER_AGENT")
)

func LoadTeams() {
	url := os.Getenv("TEAM_URL")
	req, _ := http.NewRequest(http.MethodGet, url + "?startDate=2020-01-01&endDate=2020-12-31", nil)
	req.Header.Set("User-Agent", userAgent)

	res, err := client.Do(req)
  if err != nil {
    log.Fatal(err)
  }
  defer res.Body.Close()
  if res.StatusCode != 200 {
  	fmt.Println(res)
    log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
  }

  doc, err := goquery.NewDocumentFromReader(res.Body)
  if err != nil {
    log.Fatal(err)
  }

  doc.Find(".stats-table.player-ratings-table tbody tr").Each(func(i int, s *goquery.Selection) {
    name := strings.TrimSpace(s.Find("td.teamCol-teams-overview a").Text())   
    teams = append(teams, name)
  })
}

func GetTeams() []string {
	fmt.Println(teams)
	return teams;
}

func LoadMatches() {
	url := os.Getenv("MATCHES_URL")
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", userAgent)

  res, err := client.Do(req)
  if err != nil {
    log.Fatal(err)
  }
  defer res.Body.Close()
  if res.StatusCode != 200 {
  	fmt.Println(res)
    log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
  }

  doc, err := goquery.NewDocumentFromReader(res.Body)
  if err != nil {
    log.Fatal(err)
  }

  doc.Find(".stats-table.matches-table tbody tr").Each(func(i int, s *goquery.Selection) {
  	matchId, _ := s.Find("td.date-col a").Attr("href")
    match := match{
    	date: strings.TrimSpace(s.Find("td.date-col").Text()),
    	team1: strings.TrimSpace(s.Find("td.team-col").First().Find("a").Text()),
    	team2: strings.TrimSpace(s.Find("td.team-col").Last().Find("a").Text()),
    	result: strings.TrimSpace(s.Find("td.team-col .score").First().Text()) + "-" + strings.TrimSpace(s.Find("td.team-col .score").Last().Text()),
    	csMap: strings.TrimSpace(s.Find("td.statsDetail .dynamic-map-name-full").Text()),
    	matchId: matchId,
  	}

  	// fmt.Println(match)
  	matches = append(matches, match)
  })

  // fmt.Println(matches)
}

// func (&m match) LoadMatch(url string) {
// 
// }