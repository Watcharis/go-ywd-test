package service

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"watcharis/ywd-test/model"
	"watcharis/ywd-test/src/product"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type FileHeader struct {
	Filename string
	Header   textproto.MIMEHeader
	Size     int64

	content []byte
	tmpfile string
}

type ProductRepository struct {
	repository product.ProductRepository
}

func NewProductService(repository product.ProductRepository) *ProductRepository {
	return &ProductRepository{
		repository: repository,
	}
}

func (r ProductRepository) AddProduct(ctx echo.Context) error {
	bodyProduct := model.ProductRequest{}
	if reqbody := ctx.Bind(&bodyProduct); reqbody != nil {
		// fmt.Println("reqbody :", reqbody)
		return reqbody
	}
	if err := ctx.Validate(&bodyProduct); err != nil {
		// fmt.Println("reqbody :", reqbody)
		return err
	}
	if bodyProduct.ProductName == "" {
		return ctx.JSON(http.StatusOK, model.JsonResponse{Message: "invalide format ProductName", Status: "fail", Data: ""})
	}
	proDuctPoint, err := strconv.Atoi(bodyProduct.ProductPoint)
	if proDuctPoint <= 0 {
		return ctx.JSON(http.StatusOK, model.JsonResponse{Message: "ProductPonit must more than zero", Status: "fail", Data: ""})
	}
	if err != nil {
		return ctx.JSON(http.StatusOK, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
	}
	if bodyProduct.ProductName != "" && proDuctPoint > 0 {
		dataProduct := model.ProductRequest{
			ProductName:  bodyProduct.ProductName,
			ProductPoint: bodyProduct.ProductPoint,
		}
		validateProduct := r.repository.FindProductByName(dataProduct)
		if len(validateProduct) != 0 {
			return ctx.JSON(http.StatusOK, model.JsonResponse{Message: "product is exists", Status: "fail", Data: ""})
		} else {
			insertItemProduct := r.repository.InsertProduct(dataProduct)
			if insertItemProduct == "insert product success" {
				return ctx.JSON(http.StatusOK, model.JsonResponse{Message: insertItemProduct, Status: "success", Data: ""})
			} else {
				return ctx.JSON(http.StatusOK, model.JsonResponse{Message: "add product fail", Status: "fail", Data: ""})
			}
		}
	} else {
		return ctx.JSON(http.StatusOK, model.JsonResponse{Message: "product_point must more than zero and invalide format product_name", Status: "fail", Data: ""})
	}
}

func (r ProductRepository) GetProduct(ctx echo.Context) error {
	getProduct := r.repository.FindAllProduct()
	return ctx.JSON(http.StatusOK, model.JsonResponse{Message: "get product success", Status: "success", Data: getProduct})
}

func (r ProductRepository) FileProduct(ctx echo.Context) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return err
	}

	fileType := strings.Split(file.Header["Content-Type"][0], "/")
	pathCwd, _ := os.Getwd()
	// Open() คือ การอ่าน file จาก content ของ ctx.FormFile("file")
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	var fileName string
	if fileType[1] == "jpeg" || fileType[1] == "jpg" || fileType[1] == "png" {
		fileName = fmt.Sprintf("%s.%s", strconv.Itoa(int(time.Now().UnixNano())), fileType[1])
		// Destination
		pathFile := filepath.Join(pathCwd, "img", fileName)
		// fmt.Println("pathFile :", pathFile)

		// os.Create() คือ การสร้าง destination ของ file
		dst, err := os.Create(pathFile)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		// os.Copy() คือ การสร้าง file โดยบอก destination ของ file เเละ เเนบ data มา
		if _, err = io.Copy(dst, src); err != nil {
			return ctx.JSON(http.StatusBadRequest, model.JsonResponse{Message: err.Error(), Status: "fail", Data: ""})
		}

		// os.Stat() จะทำการ check path file ของ file ที่เรา input
		// os.IsNotExist() จะรับ error หาก ไม่มี file ที่เราค้นหาจะ return true เเต่ถ้ามีจะ return false
		if _, err := os.Stat(pathFile); os.IsNotExist(err) {
			// กรณี ไม่มี file ที่เราค้นหาจะ return true
			logrus.Errorln("File does not exist ->", os.IsNotExist(err))

			if err := os.Remove(pathFile); err != nil {
				logrus.Errorln("Remove file fail ->", err)
			}
			return ctx.JSON(http.StatusBadRequest, model.JsonResponse{Message: "Create File fail -> File does not exist", Status: "fail", Data: ""})
		}
		//มี file
		dataFile, err := ioutil.ReadFile(pathFile)
		if err != nil {
			logrus.Errorln("Error: Readfile ->", err)
		}
		// http.DetectContentType() จะทำการ return type file ที่อ่านมา
		mimeType := http.DetectContentType(dataFile)

		var edcodeByteToBase64 string
		edcodeByteToBase64 = base64.StdEncoding.EncodeToString(dataFile)

		var imageBase64 string
		imageBase64 = HeaderMineTypeBase64(mimeType) + edcodeByteToBase64

		var responseJson = map[string]string{
			"image_base64": imageBase64,
		}
		// return ctx.File(pathFile)
		return ctx.JSON(http.StatusBadRequest, model.JsonResponse{Message: "encode image to base64 success", Status: "success", Data: responseJson})
	}
	return ctx.JSON(http.StatusBadRequest, model.JsonResponse{Message: fmt.Sprintf("invalid type file -> %s", fileType[1]), Status: "fail", Data: ""})
}

func HeaderMineTypeBase64(mimeType string) string {
	var base64Encoding string
	switch base64Encoding = mimeType; base64Encoding {
	case "image/jpeg":
		base64Encoding = "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding = "data:image/png;base64,"
	}
	return base64Encoding
}
