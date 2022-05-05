package extract

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

func GetMD5(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	return hex.EncodeToString(md5Ctx.Sum(nil))
}

func NullUnPadding(in []byte) []byte {
	return bytes.TrimRight(in, string([]byte{0}))
}

func Logger() *zerolog.Logger {
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	// consoleWriter.FormatLevel = func(i interface{}) string {
	// 	return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	// }

	consoleWriter.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}
	// multi := zerolog.MultiLevelWriter(consoleWriter, os.Stdout)
	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()
	return &logger
}

func LoggerDebug(msg string) {
	Logger().Debug().Msg(msg)
}
func LoggerInfo(msg string) {
	Logger().Info().Msg(msg)
}
func LoggerError(msg string) {
	Logger().Error().Msg(msg)
}
