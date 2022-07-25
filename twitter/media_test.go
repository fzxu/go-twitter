package twitter

import (
	"log"
	"net/http"
	"testing"

	"github.com/dghubble/oauth1"
)

func TestUploadImage(t *testing.T) {

	config := oauth1.NewConfig("", "")
	token := oauth1.NewToken("", "")
	httpClient := config.Client(oauth1.NoContext, token)

	client := NewClient(httpClient)

	resp, err := http.Get("https://lh3.googleusercontent.com/q0J1L80F6cRij9c9Fl4y4mVBCJEfiVbLGHSN2ZcVu0r699BigIfa638CBo_nfvJCilb6DWm1RXdyJDsO6Tg7omdBFL_V6zjuC_YS")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	m, resp, err := client.Media.Upload(resp.Body, &MediaUploadParams{MediaCategory: "tweet_image"})
	if err != nil {
		t.Fatal(err)
	}

	log.Println(resp.StatusCode)
	log.Println(m)

	tweet, resp2, err := client.Statuses.Update("test https://opensea.io", &StatusUpdateParams{MediaIds: []int64{m.MediaId}})
	if err != nil {
		t.Fatal(err)
	}

	log.Println(resp2.StatusCode)
	log.Println(tweet)
}
