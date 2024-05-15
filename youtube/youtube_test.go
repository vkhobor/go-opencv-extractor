package youtube

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func urlParse(url string) (string, error) {
	id, err := NewYoutubeIDFromUrl(url)
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func TestUrlParse_ValidURL(t *testing.T) {
	url := "https://www.youtube.com/watch?v=abcdefghijk"
	expected := "abcdefghijk"

	result, err := urlParse(url)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestUrlParse_ValidURLLong(t *testing.T) {
	url := "https://www.youtube.com/watch?v=hstme7yg7gQ&pp=ygUOc2VhIG9mIHRoaWV2ZXM%3D"
	expected := "hstme7yg7gQ"

	result, err := urlParse(url)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestUrlParse_InvalidURL(t *testing.T) {
	url := "https://www.google.com"
	expected := ""

	result, err := urlParse(url)

	assert.EqualError(t, err, "cannot parse url")
	assert.Equal(t, expected, result)
}
