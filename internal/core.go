package internal

import (
	"fmt"
	"runtime"
	"strings"
	"sync"

	"path/filepath"

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
	urlString := ""
	urls := []string{}
	if len(args) > 0 {
		for _, arg := range args {
			if strings.HasPrefix(arg, "https") {
				urls = append(urls, arg)
			}
		}
		if len(urls) > 0 {
			urlString = urls[0]
		}
		fmt.Println(urls)
	}

	if urlString != "" && !initValidations(urlString) {
		return
	}

	clientId := soundcloud.GetClientId(urlString)

	if clientId == "" {
		fmt.Println("Something went wrong while getting the Client Id!")
		return
	}

	if search {
		keyword := getUserSearch()
		apiUrl := soundcloud.GetSeachAPIUrl(keyword, SearchLimit, offset, clientId)
		searchResult := soundcloud.SearchTracksByKeyword(apiUrl, keyword, offset, clientId)

		// select one to download
		soundData = selectSearchUrl(searchResult)
	} else {
		for _, url := range urls {
			apiUrl := soundcloud.GetTrackInfoAPIUrl(url, clientId)
			soundData := soundcloud.GetSoundMetaData(apiUrl, url, clientId)
			if soundData == nil {
				fmt.Printf("%s URL : %s \n", theme.Red("[+]"), theme.Magenta(url))
				fmt.Println(theme.Yellow("URL doesn't return a valid track. Is the track publicly accessible?"))
				continue
			}
			fmt.Printf("%s %s found. Title : %s - Duration : %s\n", theme.Green("[+]"), strings.Title(soundData.Kind), theme.Magenta(soundData.Title), theme.Magenta(theme.FormatTime(soundData.Duration)))
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
						//soundcloud.AddMetadata(t, fp)

					}(dlT)
				}
				wg.Wait()

				fmt.Printf("\n%s Playlist saved to : %s\n", theme.Green("[-]"), theme.Magenta(downloadPath))
				return
			}
			downloadTracks := soundcloud.GetFormattedDL(soundData, clientId)
			os := runtime.GOOS
			filePath := ""
			track := getTrack(downloadTracks, bestQuality)
			if os == "windows" {
				filePath = soundcloud.Download(track, filepath.FromSlash(downloadPath))
			} else if os == "linux" {
				filePath = soundcloud.Download(track, downloadPath)
			}

			if filePath == "" {
				fmt.Printf("\n%s Track was already saved to : %s\n", theme.Green("[-]"), theme.Magenta(downloadPath))
				return
			}

			// err := soundcloud.AddMetadata(track, filePath)
			// if err != nil {
			// 	fmt.Println("\n" + theme.Red("An error occurred while adding tags to the track : "+"\n"+theme.Red(err)))
			// }
			fmt.Printf("\n%s Track saved to : %s\n", theme.Green("[-]"), theme.Magenta(filepath.FromSlash(filePath)))
		}
	}
}
