package importing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYoutubeParser_ValidURL(t *testing.T) {
	url := "https://www.youtube.com/watch?v=abcdefghijk"
	expected := "abcdefghijk"

	result, err := YoutubeParser(url)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestYoutubeParser_ValidURLLong(t *testing.T) {
	url := "https://www.youtube.com/watch?v=hstme7yg7gQ&pp=ygUOc2VhIG9mIHRoaWV2ZXM%3D"
	expected := "hstme7yg7gQ"

	result, err := YoutubeParser(url)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestYoutubeParser_InvalidURL(t *testing.T) {
	url := "https://www.google.com"
	expected := ""

	result, err := YoutubeParser(url)

	assert.EqualError(t, err, "cannot parse url")
	assert.Equal(t, expected, result)
}
