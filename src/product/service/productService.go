package service

import (
	"fmt"
	"io"
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
)

type FileHeader struct {
	Filename string
	Header   textproto.MIMEHeader
	Size     int64

	content []byte
	tmpfile string
}

type productRepository struct {
	repository product.ProductRepository
}

func NewProductService(repository product.ProductRepository) *productRepository {
	return &productRepository{
		repository: repository,
	}
}

func (r productRepository) AddProduct(ctx echo.Context) error {
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

func (r productRepository) GetProduct(ctx echo.Context) error {
	getProduct := r.repository.FindAllProduct()
	return ctx.JSON(http.StatusOK, model.JsonResponse{Message: "get product success", Status: "success", Data: getProduct})
}

func (r productRepository) FileProduct(ctx echo.Context) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return err
	}
	// fmt.Println("file :", file)
	fmt.Println("file_header:", file.Header)
	fmt.Println("file_name :", file.Filename)

	fileType := strings.Split(file.Header["Content-Type"][0], "/")

	pathCwd, _ := os.Getwd()
	// Open() คือ การอ่าน file จาก content ของ ctx.FormFile("file")
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	var fileName string
	if fileType[1] == "jpeg" || fileType[1] == "jpg" {
		fileName = fmt.Sprintf("%s.%s", strconv.Itoa(int(time.Now().UnixNano())), fileType[1])
		// Destination
		pathFile := filepath.Join(pathCwd, "img", fileName)

		// os.Create() คือ การสร้าง destination ของ file
		dst, err := os.Create(pathFile)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		// os.Copy() คือ การสร้าง file โดยบอก destination ของ file เเละ เเนบ data มา
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}
		return ctx.File(pathFile)
	}
	return ctx.JSON(http.StatusBadRequest, model.JsonResponse{Message: fmt.Sprintf("invalid type file -> %s", fileType[1]), Status: "fail", Data: ""})
}
