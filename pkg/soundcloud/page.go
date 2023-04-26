package soundcloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/mcworkaholic/go-scdl/pkg/client"
	"github.com/mcworkaholic/go-scdl/pkg/theme"
)

var Sound *SoundData

func SaveResponse(url string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}

	// Format the JSON response for writing to file
	formattedJson, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}

	// Create or open the file for appending
	file, err := os.OpenFile(path.Join(".\\json", "response.json"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	} else {
		// If file is empty, add opening bracket for JSON array
		fi, _ := file.Stat()
		if fi.Size() == 0 {
			_, err = file.Write([]byte("["))
			if err != nil {
				panic(err)
			}

			// Write the formatted JSON response to file
			_, err = file.Write([]byte(formattedJson))
			if err != nil {
				panic(err)
			}

		} else { // If file is not empty, add comma separator
			_, err = file.Write([]byte(","))
			if err != nil {
				panic(err)
			}

			// Write the formatted JSON response to file
			_, err = file.Write([]byte(formattedJson))
			if err != nil {
				panic(err)
			}
		}
	}
}

func CloseJSON() {
	// Create or open the file for appending
	file, err := os.OpenFile(path.Join(".\\json", "response.json"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	_, err = file.Write([]byte("]"))
	if err != nil {
		panic(err)
	}
}

// extract some meta data under : window.__sc_hydration
// write to JSON file
// check if the track exists and open to public
func GetSoundMetaData(filePath string, apiUrl string, url string, clientId string) *SoundData {
	statusCode, body, err := client.Get(apiUrl)
	if err != nil || statusCode != http.StatusOK {
		return nil
	}

	// Unmarshal the JSON response into a SoundData struct
	var soundData SoundData
	err = json.Unmarshal(body, &soundData)
	if err != nil {
		panic(err)
	}

	// Set the file path, name, artist attrs of the JSON file
	soundData.Filepath = filepath.FromSlash(path.Join(filePath, soundData.Title+".ogg"))
	soundData.Filename = soundData.Title + ".ogg"
	t500 := "t500x500" // for getting a higher res img
	if soundData.ArtworkUrl != "" {
		soundData.ArtworkUrl = strings.Replace(soundData.ArtworkUrl, "large", t500, 1)
	}
	permaUrl := soundData.PermalinkUrl
	// Find the index of the last "/" and the second to last "/"
	lastSlashIndex := strings.LastIndex(permaUrl, "/")
	secondToLastSlashIndex := strings.LastIndex(permaUrl[:lastSlashIndex], "/")

	// Extract the text between the slashes
	textBetweenSlashes := permaUrl[secondToLastSlashIndex+1 : lastSlashIndex]

	// Replace any "_" or "-" characters with spaces
	cleanedText := strings.ReplaceAll(textBetweenSlashes, "_", " ")
	cleanedText = strings.ReplaceAll(cleanedText, "-", " ")

	// Split the cleaned text into words and capitalize each one
	words := strings.Split(cleanedText, " ")
	for i, word := range words {
		words[i] = strings.Title(word)
	}

	// Join the words back into a single string
	formattedUsername := strings.Join(words, " ")
	soundData.Username = formattedUsername

	return &soundData
}
func GetClientId(url string) string {

	if url == "" {
		// the best url ever, if you find this then you're so cool :D I love you :DDD
		url = "https://soundcloud.com/ahmed-yehia0"
	}

	statusCode, bodyData, err := client.Get(url)
	if statusCode != 200 {
		fmt.Println(theme.Red("Bad URL for Client ID. "), theme.Red(url))
		os.Exit(1)
	} else if err != nil {
		fmt.Printf("An Error : %s occurred while requesting : %s", err, url)
		os.Exit(1)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bodyData))
	if err != nil {
		fmt.Printf("failed to parse HTML: %s", err)
		os.Exit(1)
	}

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
	bodyData, err = io.ReadAll(resp.Body)
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
				fmt.Println(theme.Red("Error decoding json: "), theme.Red(err))
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
