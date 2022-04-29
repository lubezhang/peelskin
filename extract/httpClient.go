package extract

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func HttpGet(url string) (io.Reader, error) {
	// LoggerDebug("HttpGetFile: " + url)
	resp, err := http.Get(url)
	// fmt.Println("resp.Header:::::", resp.Cookies())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp.Body, nil
}

func HttpGetData(url string) ([]byte, error) {
	body, err1 := HttpGet(url)
	if err1 != nil {
		return nil, err1
	}

	buf, err2 := ioutil.ReadAll(body)
	if err2 != nil {
		return nil, err2
	}
	return buf, nil
}

func HttpPostForm(apiUrl string, data string) (string, error) {
	formData := url.Values{
		"url": {data},
	}
	resp, err := http.PostForm(apiUrl, formData)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	buf, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return "", err2
	}

	return string(buf), nil
}

func HttpGetDocument(url string) (*goquery.Document, error) {
	res, err1 := http.Get(url)
	if err1 != nil {
		return nil, err1
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc, nil
}
