package extract

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func NewExtractO8tv(pageUrl string) ExtractO8tv {
	return ExtractO8tv{
		pageUrl: pageUrl,
	}
}

type ExtractO8tv struct {
	pageUrl string // 视频播放页面url
}

func (obj *ExtractO8tv) ExtractM3u8Url() (string, error) {
	videoPlayPageUrl, _ := obj.getVideoPlayPageUrl()
	fmt.Println("videoPlayPageUrl:", videoPlayPageUrl)

	cryptoM3u8Url, err2 := obj.getCryptoM3u8Url(videoPlayPageUrl)
	if err2 != nil {
		fmt.Println("", err2)
	}
	fmt.Println("cryptoM3u8Url:", cryptoM3u8Url)

	m3u8Url, _ := obj.decryptM3u8Url(cryptoM3u8Url)
	fmt.Println("m3u8Url:", m3u8Url)
	return "", nil
}

// 获取加密的m3u8的url
func (obj *ExtractO8tv) getCryptoM3u8Url(url string) (result string, err error) {
	baseApiUrl := "https://zyz.021huaying.com/duoduo/api.php"
	res, err := HttpPostForm(baseApiUrl, url)
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(res), &dat); err == nil {
		result = dat["url"].(string)
	}
	return
}

// 获取视频播放组件的url
func (obj *ExtractO8tv) getVideoPlayPageUrl() (result string, err error) {
	doc, err := HttpGetDocument(obj.pageUrl)
	if err != nil {
		return
	}
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		htmlStr := s.Text()
		if strings.Contains(htmlStr, "var player_aaaa=") {
			htmlStr = strings.Replace(htmlStr, "var player_aaaa=", "", -1)
			var dat map[string]interface{}
			if err := json.Unmarshal([]byte(htmlStr), &dat); err == nil {
				originUrl := dat["url"]
				decodeData, _ := base64.StdEncoding.DecodeString(originUrl.(string))
				decodeData1 := string(decodeData)
				decodeData2, _ := url.QueryUnescape(decodeData1)
				result = decodeData2
			} else {
				fmt.Println("==============json str 转map 失败=======================", err)
			}
		}
	})
	return
}

func (obj *ExtractO8tv) decryptM3u8Url(cryptoM3u8Url string) (result string, err error) {
	key := GetMD5("rXjWvXl6")
	iv := []byte("NXbHoWJbpsEOin8b")

	decodeData, _ := base64.StdEncoding.DecodeString(cryptoM3u8Url)
	block, err := aes.NewCipher([]byte(key)) //选择加密算法
	if err != nil {
		return "", fmt.Errorf("key 长度必须 16/24/32长度: %s", err)
	}
	blockModel := cipher.NewCBCDecrypter(block, iv[:block.BlockSize()])
	plantText := make([]byte, len(decodeData))
	blockModel.CryptBlocks(plantText, decodeData)
	plantText = NullUnPadding(plantText)
	result = string(plantText)
	return
}
