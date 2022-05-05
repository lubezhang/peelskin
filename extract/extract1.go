package extract

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Step2 struct {
	url   string
	time  string
	ua    string
	cip   string
	vkey  string
	fvkey string
}

const (
	// CONST_BASE_UA   = "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1"
	CONST_BASE_UA   = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36"
	CONST_BASE_URL  = "https://jx.xmflv.com"
	CONST_BASE_URL1 = CONST_BASE_URL + "/?url="       // https://jx.m3u8.pw/?url=
	CONST_BASE_URL2 = CONST_BASE_URL + "/player.php"  // https://jx.xmflv.com/player.php?time=1650878119&url=https://v.qq.com/x/cover/mzc00200imi2b3v/l00344c9o6b.html
	CONST_BASE_URL3 = CONST_BASE_URL + "/xmflv-1.SVG" // sdf
)

func ExtractM3u8(videoPageUrl string) error {
	LoggerDebug("视频源：" + videoPageUrl)
	doc1, err := HttpGetStep1(videoPageUrl)
	if err != nil {
		log.Fatal(err)
	}

	step1Time := ""
	doc1.Find("script").Each(func(i int, s *goquery.Selection) {
		if i == 1 {
			step1Time = extractVal(s.Text(), "time")
			fmt.Println()
		}
	})

	if len(step1Time) != 0 {
		stepParam := Step2{
			url: videoPageUrl,
		}
		doc2, err := HttpGetStep2(videoPageUrl, step1Time)
		if err != nil {
			log.Fatal(err)
		}

		doc2.Find("script").Each(func(i int, s *goquery.Selection) {
			if i == 6 {
				stepParam.time = extractVal(s.Text(), "time")
				stepParam.ua = extractVal(s.Text(), "ua")
				stepParam.cip = extractVal(s.Text(), "cip")
				stepParam.vkey = extractVal(s.Text(), "vkey")
				stepParam.fvkey = extractVal(s.Text(), "fvkey")
			}
		})
		// fmt.Println(stepParam)

		time.Sleep(time.Second * 2)

		HttpGetData("https://jx.xmflv.com/css/DPlayer.min.css")

		m3u8Url, _ := HttpGetStep3(stepParam)
		fmt.Println("m3u8Url:", m3u8Url)
	}
	return nil
}

func extractVal(str string, key string) string {
	arrHls := strings.Split(str, "\n")
	for _, v := range arrHls {
		val := strings.TrimSpace(v)
		if len(val) != 0 && strings.Contains(val, "var "+key) {
			val = strings.Replace(val, "var "+key+" = '", "", -1)
			val = strings.Replace(val, "';", "", -1)
			return val
		}
	}
	return ""
}

func HttpGetStep1(url string) (*goquery.Document, error) {
	res, err := http.Get(CONST_BASE_URL1 + url)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range res.Cookies() {
		fmt.Printf("Cookies:%+v\n", v)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc, nil
}

func HttpGetStep2(url string, time string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", CONST_BASE_URL2, nil)
	req.Header.Set("User-Agent", CONST_BASE_UA)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := req.URL.Query()
	q.Add("url", url)
	q.Add("time", time)
	req.URL.RawQuery = q.Encode()

	// fmt.Println(req.URL.String())
	resp, err := http.DefaultClient.Do(req)
	// fmt.Println("resp.Header::12:::", resp.Cookies())
	for _, v := range resp.Cookies() {
		fmt.Printf("Cookies:%+v\n", v)
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc, nil
}

func HttpGetStep3(stepParam Step2) (string, error) {
	data := url.Values{
		// "start":  {"0"},
		// "offset": {"xxxx"},
		"url":   {stepParam.url},
		"time":  {stepParam.time},
		"ua":    {stepParam.ua},
		"cip":   {stepParam.cip},
		"vkey":  {stepParam.vkey},
		"fvkey": {stepParam.fvkey},
	}
	body := strings.NewReader(data.Encode())
	req, err := http.NewRequest("POST", CONST_BASE_URL3, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	// req.Header.Set(":authority", "jx.xmflv.com")
	// req.Header.Set(":method", "POST")
	// req.Header.Set(":path", "/xmflv-1.SVG")
	// req.Header.Set(":scheme", "https")
	req.Header.Set("Origin", "https://jx.xmflv.com")
	req.Header.Set("User-Agent", CONST_BASE_UA)
	req.Header.Set("Cookie", "__51vcke__JKT8F7tQPL0PV9l2=962827e3-ed8c-5caf-828d-2134ddc99785; __51vuft__JKT8F7tQPL0PV9l2=1650869305543; __51uvsct__JKT8F7tQPL0PV9l2=2; __vtins__JKT8F7tQPL0PV9l2=%7B%22sid%22%3A%20%227f7819fb-f945-5e92-a6d8-f8e3b0cfbedb%22%2C%20%22vd%22%3A%203%2C%20%22stt%22%3A%2073802%2C%20%22dr%22%3A%2051115%2C%20%22expires%22%3A%201650879855316%2C%20%22ct%22%3A%201650878055316%7D")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	// fmt.Println(req.Header)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	for _, v := range resp.Cookies() {
		fmt.Printf("Cookies:%+v\n", v)
	}

	buf, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		fmt.Println(err2)
	}
	fmt.Println("================", string(buf))

	return "", nil
}
