package checkurl_test

import (
	"context"
	"testing"

	checkurl "github.com/blindlobstar/go-interview-problems/13-check-urls"
	"github.com/stretchr/testify/assert"
)

func Test_Check_IncorrectUrl(t *testing.T) {
	urls := []string{
		"https://go.dev/",
		"https://go.dev/learn/",
		"not_url",
	}

	got, _ := checkurl.CheckUrls(context.Background(), 2, urls)
	assert.ElementsMatch(t, urls[:2], got, "remove wrong url")
}
