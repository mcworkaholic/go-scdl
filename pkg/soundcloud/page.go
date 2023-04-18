package soundcloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
	"github.com/mcworkaholic/go-scdl/pkg/client"
	"github.com/mcworkaholic/go-scdl/pkg/theme"
)

var Sound *SoundData

// extract some meta data under : window.__sc_hydration
// check if the track exists and open to public
func GetSoundMetaData(apiUrl string, url string, clientId string) *SoundData {

	statusCode, body, err := client.Get(apiUrl)

	if err != nil || statusCode != http.StatusOK {
		return nil
	}

	json.Unmarshal(body, &Sound)

	return Sound
}

func GetClientId() string {
	url := ""
	err := godotenv.Load()
	if err != nil {
		fmt.Println(theme.Red("Error loading .env file"))
	} else {
		url = os.Getenv("URL")
	}

	statusCode, bodyData, err := client.Get(url)
	if statusCode != 200 {
		fmt.Println(theme.Red("Bad URL for Client ID. Check your .env file. -->"), theme.Red(url))
		os.Exit(1)
	} else if err != nil {
		fmt.Printf("An Error : %s occurred while requesting : %s", err, url)
		os.Exit(1)
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bodyData))

	// find the last src under the body
	apiurl, exists := doc.Find("body > script").Last().Attr("src")
	if !exists {
		return ""
	}

	// making a GET request to find the client_id
	resp, err := http.Get(apiurl)
	if err != nil {
		fmt.Printf("Something went wrong while requesting : %s , Error : %s", apiurl, err)
	}

	// reading the body
	bodyData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	// search for the client_id
	pattern := ",client_id:\"([^\"]*?.[^\"]*?)\""
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(string(bodyData), 1)

	return matches[0][1]
}

func GetFormattedDL(track *SoundData, clientId string) []DownloadTrack {

	ext := "mp3" // the default extension type
	tracks := make([]DownloadTrack, 0)
	data := track.Transcodes.Transcodings
	var wg sync.WaitGroup

	for _, tcode := range data {
		wg.Add(1)
		go func(tcode Transcode) {
			defer wg.Done()

			url := tcode.ApiUrl + "?client_id=" + clientId
			statusCode, body, err := client.Get(url)
			if err != nil && statusCode != http.StatusOK {
				return
			}
			q := mapQuality(tcode.ApiUrl, tcode.Format.MimeType)
			if q == "high" {
				ext = "ogg"
			}
			mediaUrl := Media{}
			dec := json.NewDecoder(bytes.NewReader(body))
			if err := dec.Decode(&mediaUrl); err != nil {
				log.Println("Error decoding json: ", err)
				return
			}
			tmpTrack := DownloadTrack{
				Url:       mediaUrl.Url,
				Quality:   q,
				SoundData: track,
				Ext:       ext,
			}
			tracks = append(tracks, tmpTrack)

		}(tcode)
	}
	wg.Wait()
	return tracks
}

// check if the trackUrl is mp3:progressive or ogg:hls
func mapQuality(url string, format string) string {
	tmp := strings.Split(url, "/")
	if tmp[len(tmp)-1] == "hls" && strings.HasPrefix(format, "audio/ogg") {
		return "high"
	} else if tmp[len(tmp)-1] == "hls" && strings.HasPrefix(format, "audio/mpeg") {
		return "medium"
	}
	return "low"
}

func SearchTracksByKeyword(apiUrl string, keyword string, offset int, clientId string) *SearchResult {

	statusCode, body, err := client.Get(apiUrl)

	if err != nil && statusCode != http.StatusOK {
		return nil
	}

	var result = SearchResult{}
	json.Unmarshal(body, &result)

	return &result
}
