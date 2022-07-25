package twitter

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/dghubble/sling"
)

type MediaService struct {
	sling *sling.Sling
}

type MediaUploadParams struct {
	MediaCategory string `json:"media_category"`
}

type Media struct {
	MediaId          int64      `json:"media_id"`
	MediaIdString    string     `json:"media_id_string"`
	MediaKey         string     `json:"media_key"`
	Size             int        `json:"size"`
	ExpiresAfterSecs int        `json:"expires_after_secs"`
	Image            MediaImage `json:"image"`
}

type MediaImage struct {
	ImageType string `json:"image_type"`
	Width     int    `json:"w"`
	Height    int    `json:"h"`
}

func newMediaService(sling *sling.Sling) *MediaService {
	return &MediaService{
		sling: sling.Path("media/"),
	}
}

type multiFormBodyProvider struct {
	b    []byte
	form *multipart.Writer
}

func (p multiFormBodyProvider) ContentType() string {
	return p.form.FormDataContentType()
}

func (p multiFormBodyProvider) Body() (io.Reader, error) {
	return bytes.NewReader(p.b), nil
}

func newMultiFormProvider(in io.Reader) (*multiFormBodyProvider, error) {
	b := &bytes.Buffer{}
	form := multipart.NewWriter(b)

	// create media paramater
	fw, err := form.CreateFormFile("media", "file.png")
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(fw, in)
	if err != nil {
		return nil, err
	}

	// close form
	form.Close()

	return &multiFormBodyProvider{b: b.Bytes(), form: form}, nil
}

func (s *MediaService) Upload(file io.Reader, params *MediaUploadParams) (*Media, *http.Response, error) {
	multiFormProvider, err := newMultiFormProvider(file)
	if err != nil {
		return nil, nil, err
	}

	urlPath := fmt.Sprintf("upload.json?media_category=%s", params.MediaCategory)
	media := new(Media)
	apiError := new(APIError)
	resp, err := s.sling.New().Post(urlPath).BodyProvider(multiFormProvider).Receive(media, apiError)
	return media, resp, relevantError(err, *apiError)
}
