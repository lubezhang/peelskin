package extract

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractO8tvGetM3u8(t *testing.T) {
	assetObj := assert.New(t)
	assetObj.Equal("e10adc3949ba59abbe56e057f20f883e", "e10adc3949ba59abbe56e057f20f883e")

	// pageUrl := "https://www.o8tv.com/vodplay/344227-1-1.html"
	pageUrl := "https://www.o8tv.com/vodplay/96177-1-21.html"
	extractO8tv := NewExtractO8tv(pageUrl)
	extractO8tv.ExtractM3u8Url()
	// extractO8tv.getVideoPlayPageUrl()

	// fmt.Println("m3u8Url:", m3u8Url)
}
