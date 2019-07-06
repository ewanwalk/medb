package handbrake

import (
	"context"
	"encoder-backend/pkg/models"
	"encoding/json"
	"testing"
)

func TestCommand_scan(t *testing.T) {

	cmd := Command{
		binary: "HandBrakeCLI",
		profile: models.QualityProfile{
			AudioContainer: "aac",
			AudioTracks:    2,
		},
	}

	t.Logf("%#v", cmd.profile.AudioCodecMap())

	file := "/mnt/s/coding/golang/encoder-backend/library/movies/Anchorman 2 The Legend Continues (2013)/Anchorman 2 The Legend Continues (2013) Bluray-1080p.mkv"
	//file = "/mnt/freenas/plex-main/media/movies/10 Things I Hate About You (1999)/10 Things I Hate About You.mkv"
	//file = "/mnt/freenas/plex-main/media/movies/A Billion Lives (2016)/A Billion Lives (2016) WEBDL-720p.mkv"
	//file = "/mnt/freenas/plex-main/media/movies/Everybody Wants Some!! (2016)/Everybody Wants Some!! (2016) Bluray-1080p.mkv"
	//file = "/mnt/freenas/plex-main/media/movies/Blade Runner (1982)/Blade Runner (1982) Bluray-1080p.mkv"
	//file = "/mnt/freenas/plex-main/media/movies/The Greatest Showman (2017)/The Greatest Showman (2017) WEBDL-1080p.mkv"
	//file = "/mnt/freenas/plex-main/media/movies/Bunohan (2011)/Bunohan (2011).avi"
	//file = "/mnt/freenas/plex-main/media/movies/The Dark Tower (2017)/The Dark Tower (2017) Remux-2160p.mkv"
	//file = "/mnt/freenas/plex-main/media/anime_movies/Castle in the Sky (1986).mkv"
	//file = "/mnt/freenas/plex-main/media/anime_movies/Harmony (2015).mkv"

	//_ = file
	title, err := cmd.scan(context.Background(), file)
	if err != nil {
		t.Error(err)
		return
	}

	out, err := json.MarshalIndent(title.AudioList, "", "  ")

	t.Logf("%s", string(out))

	err = cmd.get(context.Background(), file)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%#v", cmd.args)
}
