package main 

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"strings"
)

type Song struct {
	ID        string `json:"id"`
	Placement string `json:"placement"`
	Title     string `json:"name"`
	Artist    string `json:"artist"`
}

type Album struct {
	ID        string `json:"id"`
	Placement string `json:"placement"`
	Title     string `json:"name"`
	Artist    string `json:"artist"`
}

var (
	count     = 0
	songList  = ""
	albumList = ""
	Songs     []Song
	Albums    []Album
)

func main() {
    var records []Song
	j := MusicTopTen()

	if err := json.Unmarshal([]byte(j), &records); err != nil {
		panic(err)
	}

    fmt.Println("------------------------------------")
	fmt.Println("Here are the top 10 songs this week:\n")
	for _, record := range records {
        fmt.Printf("%3s %1s %-2s by %2s\n", record.Placement, "-",  record.Title, record.Artist)
	}
    fmt.Println("")
}

func MusicTopTen() string {
	ScrapeUrl := "https://musicchartsarchive.com/"

	c := colly.NewCollector(colly.AllowedDomains("musicchartsarchive.com"))

	c.OnRequest(func(r *colly.Request) {
	})

	c.OnHTML("tr.odd td", func(h *colly.HTMLElement) {
		if count < 30 {
			songList += h.Text + "\n"
		} else {
			albumList += h.Text + "\n"
		}
		count++
	})

	c.OnScraped(func(r *colly.Response) {
	})

	c.Visit(ScrapeUrl)
	songsJson := ParseSongs(songList)

	return songsJson
}

func ParseSongs(body string) string {
	scanner := bufio.NewScanner(strings.NewReader(body))
	objCount := 0
	propCount := 0
	currentObj := 0
	stempObj := Song{} //update temp obj and whenever objCount iterates, save and set to a new object. //update temp obj and whenever objCount iterates, save and set to a new object.
	for scanner.Scan() {
		//keep track of what object we are creating: objCount
		//keep track of what properties we are assigning values to: propCount
		//How?
		//everytime propCount == 6, objCount++
		//everytime propCount == 6, propCount = 0

		if propCount == 3 {
			objCount++
			propCount = 1
		} else {
			propCount++
		}

		//if current object isn't
		if currentObj != objCount {
			//move to a different object to populate with data
			//do something then set to object count
			if objCount <= 10 {
				Songs = append(Songs, stempObj)
			}
			currentObj = objCount
		}

		switch propCount {
		case 1:
			//assign the id to the id.id

			stempObj.ID = scanner.Text()
			stempObj.Placement = scanner.Text()
		case 2:
			stempObj.Title = scanner.Text()
		case 3:
			stempObj.Artist = scanner.Text()
		default:
		}

	}

	//fmt.Println("-------------------------------------S    I    T     E    S----------------------------------")
	Songs = append(Songs, stempObj)
	j, _ := json.MarshalIndent(Songs, "", "  ")
	//log.Println(string(j))
	//fmt.Println(string(j))

	if err := scanner.Err(); err != nil {
		fmt.Printf("error occurred: %v\n", err)
	}
	return string(j)
}

func ParseAlbums(body string) string {
	scanner := bufio.NewScanner(strings.NewReader(body))

	objCount := 0
	propCount := 0
	currentObj := 0
	atempObj := Album{} //update temp obj and whenever objCount iterates, save and set to a new object. //update temp obj and whenever objCount iterates, save and set to a new object.
	for scanner.Scan() {
		//keep track of what object we are creating: objCount
		//keep track of what properties we are assigning values to: propCount
		//How?
		//everytime propCount == 6, objCount++
		//everytime propCount == 6, propCount = 0

		if propCount == 3 {
			objCount++
			propCount = 1
		} else {
			propCount++
		}

		//if current object isn't equal to the objCount, which increments every 3 lines,
		if currentObj != objCount {
			//move to a different object to populate with data
			//do something then set to object count
			if objCount < 11 {
				Albums = append(Albums, atempObj)
			}
			currentObj = objCount
		}

		switch propCount {
		case 1:
			//assign the id to the id.id
			atempObj.ID = scanner.Text()
			atempObj.Placement = scanner.Text()
		case 2:
			atempObj.Title = scanner.Text()
		case 3:
			atempObj.Artist = scanner.Text()
		default:

		}

	}

	//fmt.Println("-------------------------------------S    I    T     E    S----------------------------------")
	Albums = append(Albums, atempObj)
	j, _ := json.MarshalIndent(Albums, "", "  ")
	//log.Println(string(j))
	//fmt.Println(string(j))

	if err := scanner.Err(); err != nil {
		fmt.Printf("error occurred: %v\n", err)
	}
	return string(j)
}
