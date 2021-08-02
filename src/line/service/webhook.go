package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"watcharis/ywd-test/model"
	"watcharis/ywd-test/src/line"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type service struct{}

func NewLineService() service {
	return service{}
}

func replyMessage(message line.ReplyMessage) (string, error) {
	jsonData, err := json.Marshal(message)
	if err != nil {
		logrus.Errorln("Error json.Marshal() ->", err)
	}

	url := "https://api.line.me/v2/bot/message/reply"
	ChannelToken := os.Getenv("CHANNELTOKEN")
	var jsonStr = []byte(jsonData)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+ChannelToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()

	bodyPost, errPost := ioutil.ReadAll(resp.Body)
	if errPost != nil {
		log.Fatal(errPost)
		return "", errPost
	}
	if string(bodyPost) == "{}" {
		return "success", nil
	}
	return "", nil
}

func validateSignature(channelSecret, signature string, body []byte) bool {
	decoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		logrus.Errorln("Error decoded signature from header ->", err)
		return false
	}
	hash := hmac.New(sha256.New, []byte(channelSecret))

	_, err = hash.Write(body)
	if err != nil {
		return false
	}

	// Compare decoded signature and `hash.Sum(nil)` by using `hmac.Equal`
	//hmac.Equal(decoded, hash.Sum(nil))
	return hmac.Equal(decoded, hash.Sum(nil))
}

func (l service) Webhook(ctx echo.Context) error {

	Line := line.LineMessage{}
	// header := ctx.Request().Header
	// fmt.Println("header :", header)

	if err := ctx.Bind(&Line); err != nil {
		logrus.Errorln("Error bind data webhook ->", err)
		return ctx.JSON(http.StatusOK, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
	}
	if err := ctx.Validate(Line); err != nil {
		return err
	}

	// fmt.Println("ctx.Request().Body ->", ctx.Request().Body)
	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		logrus.Errorln("Error : read body ->", err)
	}

	// ติด ***bug*** ยังไม่สามารถ หาวิธี validate signature ที่ถูกต้องได้
	channelSecret := os.Getenv("CHANNELSECRET")
	signature := ctx.Request().Header.Get("X-Line-Signature")
	if !validateSignature(channelSecret, signature, body) {
		return ctx.JSON(http.StatusOK, model.JsonResponse{Message: "false", Status: "fail", Data: ""})
	}

	message := line.ReplyMessage{
		ReplyToken: Line.Events[0].ReplyToken,
		Messages: []line.Text{
			{Type: "text", Text: "สวัสดีจ้า"},
			{Type: "text", Text: "ได้คุยกันสักที"},
		},
		NotificationDisabled: false,
	}
	statusMessage, err := replyMessage(message)
	if err != nil {
		return err
	}
	fmt.Println("statusMessage :", statusMessage)

	return ctx.JSON(http.StatusOK, model.JsonResponse{Message: "success ok", Status: "success", Data: ""})
}
