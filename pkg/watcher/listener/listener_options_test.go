package listener

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOptions_IsAllowedExtension(t *testing.T) {

	o := options{
		ExtensionWhitelist: []string{"mkv", "mp4", "avi"},
	}

	ans := o.IsAllowedExtension("Avengers Endgame (2019) WEBDL-2160p.mkv.partial~")
	assert.False(t, ans)

	ans = o.IsAllowedExtension("Avengers Endgame (2019) WEBDL-2160p.mov")
	assert.False(t, ans)

	ans = o.IsAllowedExtension("Avengers Endgame (2019) WEBDL-2160p.mkv")
	assert.True(t, ans)

}
