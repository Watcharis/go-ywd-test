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

func (l service) Webhook(ctx echo.Context) error {

	Line := line.LineMessage{}
	header := ctx.Request().Header
	fmt.Println("header :", header)

	if err := ctx.Bind(&Line); err != nil {
		logrus.Errorln("Error bind data webhook ->", err)
		return ctx.JSON(http.StatusOK, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
	}
	if err := ctx.Validate(Line); err != nil {
		return err
	}

	// fmt.Println(ctx.Request().Body)
	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		logrus.Errorln("Error : ")
	}
	fmt.Println("body :", body)

	signature := ctx.Request().Header.Get("x-line-signature")
	fmt.Println("signature :", signature)
	decoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		logrus.Errorln("Error decoded signature from header ->", err)
	}

	hash := hmac.New(sha256.New, []byte("83d8f32b73ea4c16285000eb25c823e3"))
	hash.Write(body)

	//// Compare decoded signature and `hash.Sum(nil)` by using `hmac.Equal`
	compare := hmac.Equal(decoded, hash.Sum(nil))
	fmt.Println("compare :", compare)

	fmt.Println("Line :", Line.Events[0].ReplyToken)

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
