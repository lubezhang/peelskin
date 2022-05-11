package extract

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	// CONST_BASE_UA   = "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1"
	CONST_BASE_XMJX_UA   = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36"
	CONST_BASE_XMJX_URL  = "https://jx.xmflv.com"
	CONST_BASE_XMJX_URL1 = CONST_BASE_XMJX_URL + "/?url="       // https://jx.m3u8.pw/?url=
	CONST_BASE_XMJX_URL2 = CONST_BASE_XMJX_URL + "/player.php"  // https://jx.xmflv.com/player.php?time=1650878119&url=https://v.qq.com/x/cover/mzc00200imi2b3v/l00344c9o6b.html
	CONST_BASE_XMJX_URL3 = CONST_BASE_XMJX_URL + "/xmflv-1.SVG" // sdf
)

func NewXmjx(url string) Xmjx {
	return Xmjx{
		videoPageUrl: url,
	}
}

type XmjxStep2Params struct {
	url   string
	time  string
	ua    string
	cip   string
	vkey  string
	fvkey string
}

type XmjxVideo struct {
	Code    int
	Success int
	Player  string
	Title   string
	Vtype   string
	Url     string
}

type Xmjx struct {
	videoPageUrl string // 视频播放页面链接
}

func (xmjx *Xmjx) ExtractVideo() (XmjxVideo, error) {
	// time1, err1 := xmjx.getStep1Time()
	// if err1 != nil {
	// 	LoggerError(err1.Error())
	// 	return XmjxVideo{}, err1
	// }
	params, err2 := xmjx.getStep2Params("1650878119")
	if err2 != nil {
		LoggerError(err2.Error())
		return XmjxVideo{}, err2
	}
	return xmjx.getStep3Video(params)
}

func (xmjx *Xmjx) getStep1Time() (result string, err error) {
	err = nil
	url := CONST_BASE_XMJX_URL1 + xmjx.videoPageUrl
	LoggerDebug("Step1 url: " + url)

	doc1, err := HttpGetDocument(url)
	if err != nil {
		return "", err
	}

	doc1.Find("script").Each(func(i int, s *goquery.Selection) {
		if i == 1 {
			result = extractVal(s.Text(), "time")
		}
	})
	LoggerDebug("Step1 time: " + result)
	if result == "" {
		err = errors.New("没有获取到时间")
	}
	return
}

func (xmjx *Xmjx) getStep2Params(time string) (result XmjxStep2Params, err error) {
	url := fmt.Sprintf("%s%s", CONST_BASE_XMJX_URL1, xmjx.videoPageUrl)
	LoggerDebug("Step2 url: " + url)
	result = XmjxStep2Params{}
	doc1, err := HttpGetDocument(url)
	if err != nil {
		return result, err
	}

	doc1.Find("script").Each(func(i int, s *goquery.Selection) {
		if i == 6 {
			pageText := s.Text()
			result.url = extractVal(pageText, "url")
			result.time = extractVal(pageText, "time")
			result.ua = extractVal(pageText, "ua")
			result.cip = extractVal(pageText, "cip")
			result.vkey = extractVal(pageText, "vkey")
			result.fvkey = extractVal(pageText, "fvkey")
		}
	})
	LoggerDebug("Step2 params fvkey: " + result.fvkey)
	return
}

func (xmjx *Xmjx) getStep3Video(params XmjxStep2Params) (result XmjxVideo, err error) {
	result = XmjxVideo{}
	formData := url.Values{
		"url":   {params.url},
		"time":  {params.time},
		"ua":    {params.ua},
		"cip":   {params.cip},
		"vkey":  {params.vkey},
		"fvkey": {xmjx.encryptFVkey(params.fvkey)},
	}

	LoggerDebug("请求参数：" + formData.Encode())
	body := strings.NewReader(formData.Encode())
	req, _ := http.NewRequest("POST", CONST_BASE_URL3, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", CONST_BASE_UA)
	// req.Header.Set("Origin", "https://jx.xmflv.com")
	resp, err2 := http.DefaultClient.Do(req)
	if err != nil {
		LoggerError(err2.Error())
		return result, err2
	}
	defer resp.Body.Close()
	buf, _ := ioutil.ReadAll(resp.Body)
	LoggerDebug("视频数据：" + string(buf))

	var dat map[string]interface{}
	if err := json.Unmarshal(buf, &dat); err == nil {
		result.Code = int(dat["code"].(float64))
		result.Success = int(dat["success"].(float64))
		result.Player = dat["player"].(string)
		result.Title = dat["title"].(string)
		result.Vtype = dat["type"].(string)
		result.Url = dat["url"].(string)
	}

	LoggerDebug("视频数据：" + result.Url)
	return
}

func (xmjx *Xmjx) encryptFVkey(fvkey string) (result string) {
	key := GetMD5(fvkey)
	iv := []byte("UVE1NTY4MDY2NQ==")

	block, err := aes.NewCipher([]byte(key)) //选择加密算法
	if err != nil {
		return ""
	}
	decodeData := []byte(fvkey)
	blockModel := cipher.NewCBCEncrypter(block, iv[:block.BlockSize()])
	plantText := make([]byte, len(decodeData))
	blockModel.CryptBlocks(plantText, decodeData)
	plantText = NullUnPadding(plantText)
	result = base64.StdEncoding.EncodeToString(plantText)
	LoggerDebug("fvkey：" + fvkey + ",  fvkey加密后：" + result)
	return
}
