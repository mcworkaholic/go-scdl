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

func SaveResponse(filePath string, apiUrl string, i int) *SoundData {
	resp, err := http.Get(apiUrl)
	if err != nil {
		return nil
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

	// Convert the map to an interface{}
	resultIfc := interface{}(result)

	// Check if the result is a map
	if resultMap, ok := resultIfc.(map[string]interface{}); ok {
		// Check if the "kind" key is set to "playlist"
		if kind, ok := resultMap["kind"].(string); ok && kind == "playlist" {
			// Check if the "tracks" key exists
			if tracks, ok := resultMap["tracks"].([]interface{}); ok {
				// "tracks" is an array of tracks

				// Loop through each track and add the extra fields
				for i := 0; i < len(tracks); i++ {
					if track, ok := tracks[i].(map[string]interface{}); ok {
						if title, ok := track["title"].(string); ok {
							// Set the file path, name, artist attrs of the JSON file
							filepath := filepath.FromSlash(path.Join(filePath, title+".ogg"))
							filename := title + ".ogg"
							// Add the extra fields to the track object

							t500 := "t500x500" // for getting a higher res img
							if track["artwork_url"] != "" {
								track["artwork_url"] = strings.Replace(track["artwork_url"].(string), "large", t500, 1)
							}

							track["file_path"] = filepath
							track["file_name"] = filename
						}
					}
				}
			}
		} else {
			// Check if the "tracks" key does not exist and "kind" is not set to "playlist"
			if _, tracksExist := resultMap["tracks"]; !tracksExist && resultMap["kind"].(string) != "playlist" {
				// Set the file path, name, artist attrs of the JSON file
				filepath := filepath.FromSlash(path.Join(filePath, resultMap["title"].(string)+".ogg"))
				filename := resultMap["title"].(string) + ".ogg"
				// Add the extra fields to the result object

				t500 := "t500x500" // for getting a higher res img
				if resultMap["artwork_url"] != "" {
					resultMap["artwork_url"] = strings.Replace(resultMap["artwork_url"].(string), "large", t500, 1)
				}
				resultMap["file_path"] = filepath
				resultMap["file_name"] = filename
			}
		}
		// Format the JSON response for writing to file
		formattedJson, err := json.MarshalIndent(&result, "", "    ")
		if err != nil {
			panic(err)
		}
		WriteJSON(formattedJson, i)
	}
	// Unmarshal the JSON response into a SoundData struct
	var soundData SoundData
	err = json.Unmarshal(body, &soundData)
	if err != nil {
		panic(err)
	}
	t500 := "t500x500" // for getting a higher res img
	if soundData.ArtworkUrl != "" {
		soundData.ArtworkUrl = strings.Replace(soundData.ArtworkUrl, "large", t500, 1)
	}
	return &soundData
}
func getExecutableDir() string {
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exeDir := filepath.Dir(exe)
	return exeDir
}

func getJSONFilePath() string {
	exeDir := getExecutableDir()
	filePath := filepath.Join(exeDir, "json", "download-cache.json")
	return filePath
}

func WriteJSON(resp []byte, i int) {
	// Delete the file if it exists
	filePath := getJSONFilePath()
	err := os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	// Create or open the file for appending
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close() // ensure that file is closed after writing

	// Determine whether to write opening bracket or comma separator
	var prefix []byte
	switch fi, _ := file.Stat(); {
	case fi.Size() == 0:
		prefix = []byte("[")
	case i > 0:
		prefix = []byte(",")
	}

	// Write prefix and response to file
	_, err = file.Write(prefix)
	if err != nil {
		panic(err)
	}
	_, err = file.Write(resp)
	if err != nil {
		panic(err)
	}
}

func CloseJSON() {
	// Create or open the file for appending
	filePath := getJSONFilePath()
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	_, err = file.Write([]byte("]"))
	if err != nil {
		panic(err)
	}
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
