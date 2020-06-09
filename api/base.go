package api

import (
	"encoding/json"
	"fmt"
	"log"
  "net/http"
  "os"
  "strconv"
  "strings"
  "regexp"
  "time"

	"github.com/PuerkitoBio/goquery"
)

type match struct {
	date, team1, team2, result, csMap string
	matchId string
}

type team struct {
	Name string    `json:"name"`
	KdDiff string  `json:"kdDiff"`
	Maps int       `json:"maps"`
	Kd float64     `json:"kd"`
	Rating float64 `json:"rating"`
}

var (
	teams []*team
	matches []match
	client = &http.Client{}
	userAgent = os.Getenv("USER_AGENT")
)

func GetTeams() []byte {
	url := os.Getenv("TEAM_URL")
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

  doc.Find(".stats-table.player-ratings-table tbody tr").Each(func(i int, s *goquery.Selection) {
  	maps, _ := strconv.Atoi(strings.TrimSpace(s.Find("td.statsDetail").Text()))
  	kd, _ := strconv.ParseFloat(strings.TrimSpace(s.Find("td.statsDetail").Text()), 64)
  	rating, _ := strconv.ParseFloat(strings.TrimSpace(s.Find("td.ratingCol").Text()), 64)

  	team := &team{
    	Name: strings.TrimSpace(s.Find("td.teamCol-teams-overview a").Text()),
    	KdDiff: strings.TrimSpace(s.Find("td.kdDiffCol.won").Text()),
    	Maps: maps,
    	Kd: kd,
    	Rating: rating,
  	}

    teams = append(teams, team)
  })

  json, err := json.Marshal(teams)
  if err != nil {
  	log.Fatal(err)	
  }

  return json
}

func GetMatch(matchId int, confrontationName string) []byte {
	url := os.Getenv("MATCHES_URL") + "\\" + string(matchId) + "\\" + confrontationName
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

  // Get total records count
  paginationStr := doc.Find(".pagination-component.pagination-top span.pagination-data").Text()
	re := regexp.MustCompile(`of\s(\d+)`)
	match := re.FindSubmatch([]byte(paginationStr))
  totalCount, _ := strconv.Atoi(string(match[len(match)-1]))

  for i := 0; i < totalCount; i += 50 {
  	go getMatchesByPage(url, i)
  }
  fmt.Println(matches)
}

func getMatchesByPage(url string, offset int) {
	if offset != 0 {
		url = url + "?offset=" + strconv.Itoa(offset)
	}
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
  time.Sleep(300 * time.Millisecond)

  // fmt.Println(len(matches))
}
// func (&m match) LoadMatch(url string) {
// 
// }