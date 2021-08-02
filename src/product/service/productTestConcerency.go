package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"watcharis/ywd-test/model"
	"watcharis/ywd-test/src/product"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func SendRequestMethodPOST(url string, c chan model.JsonResponse, e chan error) {
	// fmt.Println("url :", url)
	values := map[string]string{"name": "John Doe", "occupation": "gardener"}
	jsonData, err := json.Marshal(values)
	if err != nil {
		logrus.Errorln("Error json.Marshal() ->", err)
	}
	// fmt.Println("jsonData :", jsonData)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logrus.Errorln("Err request post ->", err)
	}
	fmt.Println("resp :", resp.Body)
	defer resp.Body.Close()

	bodyPost, errPost := ioutil.ReadAll(resp.Body)
	if errPost != nil {
		log.Fatal(errPost)
	}
	fmt.Println("bodyPost :", string(bodyPost))

	var res model.JsonResponse
	if err := json.Unmarshal(bodyPost, &res); err != nil {
		logrus.Errorln("Error json.Unmarshal ->", err)
		e <- err
	}
	c <- res
}

func (r ProductRepository) TestConcerency(ctx echo.Context) error {

	c := make(chan model.JsonResponse)
	e := make(chan error)

	FlaskUrlConOne := os.Getenv("FLASK_URL_CONONE")
	FlaskUrlConTwo := os.Getenv("FLASK_URL_CONTWO")
	FlaskUrlConThree := os.Getenv("FLASK_URL_CONTHREE")
	FlaskUrlConFour := os.Getenv("FLASK_URL_CONFOUR")
	// fmt.Println("product.UrlPython :", product.UrlPython)
	// fmt.Println("url test:", FlaskUrlConOne)
	urlFlask := []string{
		FlaskUrlConOne,
		FlaskUrlConTwo,
		FlaskUrlConThree,
		FlaskUrlConFour,
	}
	for _, v := range urlFlask {
		go SendRequestMethodPOST(v, c, e)
	}
	// fmt.Println("result :", result)

	// var result []interface{}
	// go func() {
	// 	defer close(c)
	// 	for i := range c {
	// 		fmt.Println("i :", i)
	// 	}
	// }()

	var result []model.JsonResponse
	for i := 0; i < len(product.UrlPython); i++ {
		select {
		case data := <-c:
			// fmt.Println("data :", data)
			result = append(result, data)
		case data := <-e:
			logrus.Errorln("Error: channel err ->", data)
		default:
			data := <-c
			result = append(result, data)
		}
	}

	finalResult := []map[string]int{}
	for _, values := range result {
		convertData, _ := json.Marshal(values.Data)
		// fmt.Println("convertData :", convertData)
		var data product.Data
		if err := json.Unmarshal(convertData, &data); err != nil {
			logrus.Errorln("Error json.Unmarshal() ->", err)
		}
		one := data.One
		two := data.Two
		response := map[string]int{
			"res_one": one,
			"res_two": two,
		}
		finalResult = append(finalResult, response)
	}
	return ctx.JSON(http.StatusOK, model.JsonResponse{Message: "data", Status: "success", Data: finalResult})
}
