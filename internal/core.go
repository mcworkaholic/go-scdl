package internal

import (
	"fmt"
	"strings"
	"sync"

	"github.com/mcworkaholic/go-scdl/pkg/soundcloud"
	"github.com/mcworkaholic/go-scdl/pkg/theme"
)

var (
	defaultQuality = "low"
	soundData      = &soundcloud.SoundData{}
	SearchLimit    = 6 // FIXME: hard-coded search limit
	offset         = 0
)

func Sc(args []string, downloadPath string, bestQuality bool, search bool) {

	url := ""
	if len(args) > 0 {
		url = args[0]
	}

	if url != "" && !initValidations(url) {
		return
	}

	clientId := soundcloud.GetClientId(url)

	if clientId == "" {
		fmt.Println("Something went wrong while getting the Client Id!")
		return
	}
	// --search-and-download
	if search {
		keyword := getUserSearch()
		apiUrl := soundcloud.GetSeachAPIUrl(keyword, SearchLimit, offset, clientId)
		searchResult := soundcloud.SearchTracksByKeyword(apiUrl, keyword, offset, clientId)

		// select one to download
		soundData = selectSearchUrl(searchResult)
	} else {

		apiUrl := soundcloud.GetTrackInfoAPIUrl(url, clientId)
		soundData = soundcloud.GetSoundMetaData(apiUrl, url, clientId)
		if soundData == nil {
			fmt.Printf("%s URL : %s \n", theme.Red("[+]"), theme.Magenta(url))
			fmt.Println(theme.Yellow("URL doesn't return a valid track. Is the track publicly accessible?"))
			return
		}

		fmt.Printf("%s %s found. Title : %s - Duration : %s\n", theme.Green("[+]"), strings.Title(soundData.Kind), theme.Magenta(soundData.Title), theme.Magenta(theme.FormatTime(soundData.Duration)))
	}

	// check if the url is a playlist
	if soundData.Kind == "playlist" {
		var wg sync.WaitGroup
		plDownloadTracks := getPlaylistDownloadTracks(soundData, clientId)

		for _, dlT := range plDownloadTracks {

			wg.Add(1)

			go func(dlT []soundcloud.DownloadTrack) {
				defer wg.Done()
				// bestQuality is true to avoid prompting the user for quality choosing each time and speed up
				// TODO: get a single progress bar, this will require the use of "https://github.com/cheggaaa/pb" since the current pb doesn't support download pool (I think)
				t := getTrack(dlT, true)
				fp := soundcloud.Download(t, downloadPath)

				// silent indication of already existing files
				if fp == "" {
					return
				}
				soundcloud.AddMetadata(t, fp)

			}(dlT)
		}
		wg.Wait()

		fmt.Printf("\n%s Playlist saved to : %s\n", theme.Green("[-]"), theme.Magenta(downloadPath))
		return
	}

	downloadTracks := soundcloud.GetFormattedDL(soundData, clientId)

	track := getTrack(downloadTracks, bestQuality)
	filePath := soundcloud.Download(track, downloadPath)

	// add tags
	if filePath == "" {
		fmt.Printf("\n%s Track was already saved to : %s\n", theme.Green("[-]"), theme.Magenta(downloadPath))
		return
	}
	defer soundcloud.AddMetadata(track, strings.Replace(filePath, "/", "\\", -1))
	// if err != nil {
	// 	fmt.Println()
	// 	fmt.Println(theme.Red("An error occurred while adding tags to the track : "))
	// 	fmt.Println(theme.Red(err))
	// }
	fmt.Printf("\n%s Track saved to : %s\n", theme.Green("[-]"), theme.Magenta(filePath))
}
