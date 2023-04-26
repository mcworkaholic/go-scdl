package soundcloud

type Transcode struct {
	ApiUrl  string `json:"url"`
	Quality string `json:"quality"`
	Format  Format `json:"format"`
}

type Format struct {
	Protocol string `json:"protocol"`
	MimeType string `json:"mime_type"`
}

type SoundData struct {
	Id           int64      `json:"id"`
	Filepath     string     `json:"file_path"`
	Filename     string     `json:"file_name"`
	CreatedAt    string     `json:"created_at"`
	Title        string     `json:"title"`
	Username     string     `json:"username"`
	Genre        string     `json:"genre"`
	Duration     int64      `json:"duration"`
	Kind         string     `json:"kind"`
	TrackFormat  string     `json:"track_format,omitempty"`
	PermalinkUrl string     `json:"permalink_url"`
	UserId       int64      `json:"user_id"`
	ArtworkUrl   string     `json:"artwork_url"`
	Transcodes   Transcodes `json:"media"`
	LikesCount   int        `json:"likes_count"`
	Downloadable bool       `json:"downloadable"`
	Description  string     `json:"description,omitempty"`
}

type Transcodes struct {
	Transcodings []Transcode `json:"transcodings"`
}

type Media struct {
	Url string `json:"url"`
}

type DownloadTrack struct {
	Url       string
	Size      int
	Data      []byte
	Quality   string
	Ext       string
	SoundData *SoundData
}

type SearchResult struct {
	Sounds []SoundData `json:"collection"`
	Next   string      `json:"next_href"`
}

type Track struct {
	ArtworkUrl       string      `json:"artwork_url"`
	Filepath         string      `json:"file_path"`
	Filename         string      `json:"file_name"`
	Caption          interface{} `json:"caption"`
	CommentCount     int         `json:"comment_count"`
	Commentable      bool        `json:"commentable"`
	CreatedAt        string      `json:"created_at"`
	Description      string      `json:"description"`
	DisplayDate      string      `json:"display_date"`
	DownloadCount    int         `json:"download_count"`
	Downloadable     bool        `json:"downloadable"`
	Duration         int         `json:"duration"`
	EmbeddableBy     string      `json:"embeddable_by"`
	FullDuration     int         `json:"full_duration"`
	Genre            string      `json:"genre"`
	HasDownloadsLeft bool        `json:"has_downloads_left"`
	ID               int         `json:"id"`
	Kind             string      `json:"kind"`
	LabelName        interface{} `json:"label_name"`
	LastModified     string      `json:"last_modified"`
	License          string      `json:"license"`
	LikesCount       int         `json:"likes_count"`
	Media            struct {
		Transcodings []struct {
			Duration int `json:"duration"`
			Format   struct {
				MimeType string `json:"mime_type"`
				Protocol string `json:"protocol"`
			} `json:"format"`
			Preset  string `json:"preset"`
			Quality string `json:"quality"`
			Snipped bool   `json:"snipped"`
			URL     string `json:"url"`
		} `json:"transcodings"`
	} `json:"media"`
	MonetizationModel string `json:"monetization_model"`
	Permalink         string `json:"permalink"`
	PermalinkUrl      string `json:"permalink_url"`
	PlaybackCount     int    `json:"playback_count"`
	Policy            string `json:"policy"`
	Public            bool   `json:"public"`
	PublisherMetadata struct {
		ContainsMusic bool   `json:"contains_music"`
		ID            int    `json:"id"`
		Urn           string `json:"urn"`
	} `json:"publisher_metadata"`
	PurchaseTitle    string      `json:"purchase_title"`
	PurchaseURL      interface{} `json:"purchase_url"`
	ReleaseDate      interface{} `json:"release_date"`
	RepostsCount     int         `json:"reposts_count"`
	SecretToken      interface{} `json:"secret_token"`
	Sharing          string      `json:"sharing"`
	State            string      `json:"state"`
	StationPermalink string      `json:"station_permalink"`
	StationURN       string      `json:"station_urn"`
	Streamable       bool        `json:"streamable"`
	TagList          string      `json:"tag_list"`
	Title            string      `json:"title"`
	User             struct {
		AvatarUrl        string     `json:"avatar_url"`
		Badges           UserBadges `json:"badges"`
		City             string     `json:"city"`
		CountryCode      string     `json:"country_code"`
		FirstName        string     `json:"first_name"`
		FollowersCount   int        `json:"followers_count"`
		FullName         string     `json:"full_name"`
		Id               int        `json:"id"`
		Kind             string     `json:"kind"`
		LastModified     string     `json:"last_modified"`
		LastName         string     `json:"last_name"`
		Permalink        string     `json:"permalink"`
		PermalinkUrl     string     `json:"permalink_url"`
		StationPermalink string     `json:"station_permalink"`
		StationUrn       string     `json:"station_urn"`
		Uri              string     `json:"uri"`
		Urn              string     `json:"urn"`
		Username         string     `json:"username"`
		Verified         bool       `json:"verified"`
	}
}

type UserBadges struct {
	Pro          bool `json:"pro"`
	ProUnlimited bool `json:"pro_unlimited"`
	Verified     bool `json:"verified"`
}
