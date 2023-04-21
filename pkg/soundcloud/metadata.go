// Package soundcloud adding tags to the track after downloading it.
package soundcloud

import (
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/bogem/id3v2"
	"github.com/mcworkaholic/go-scdl/pkg/client"
	"github.com/mcworkaholic/go-scdl/pkg/theme"
)

func AddMetadata(track DownloadTrack, filePath string) error {
	t500 := "t500x500" // for getting a higher res img
	imgBytes := make([]byte, 0)
	os := runtime.GOOS

	// check for artist thing
	if track.SoundData.ArtworkUrl != "" {
		url := strings.Replace(track.SoundData.ArtworkUrl, "large", t500, 1)

		// fetching the data
		statusCode, data, err := client.Get(url)
		if err != nil || statusCode != http.StatusOK {
			return err
		}
		imgBytes = data
	}

	var tag *id3v2.Tag
	var err error
	switch os {
	case "windows":
		windowsPath := filepath.FromSlash(filePath)
		tag, err = id3v2.Open(windowsPath, id3v2.Options{Parse: true})
	default:
		tag, err = id3v2.Open(filePath, id3v2.Options{Parse: true})
	}
	if err != nil {
		return err
	}
	defer func(tag *id3v2.Tag) {
		err := tag.Close()
		if err != nil {
			log.Fatalln("\n" + "id3v2 close error: " + theme.Red(err))
		}
	}(tag)

	// setting metadata
	tag.SetTitle(track.SoundData.Title)
	// TODO:
	//tag.SetArtist()
	//tag.SetAlbum()
	tag.SetGenre(track.SoundData.Genre)
	tag.SetYear(track.SoundData.CreatedAt)

	// extracting the usr
	artistName := strings.Split(track.SoundData.PermalinkUrl, "/")
	tag.SetArtist(artistName[3])

	if imgBytes != nil {
		tag.AddAttachedPicture(
			id3v2.PictureFrame{
				Encoding:    id3v2.EncodingUTF8,
				MimeType:    "image/jpeg",
				Picture:     imgBytes,
				Description: track.SoundData.Description, // well, coz why not :D
			},
		)
	}
	if err = tag.Save(); err != nil {
		log.Fatalln("\n" + theme.Red("id3v2 save error: ") + theme.Red(err))
	}
	return nil
}
