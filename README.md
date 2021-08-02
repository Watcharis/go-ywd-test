# การสร้าง Project Go By Echo

    go version ควร install ตั้งเเต่ Go v1.13 ขึ้นไป

    step 1 สร้าง package module ของ project go
        - go mod init <module name>

    step 2 ทำการ install echo
        - go get github.com/labstack/echo/v4
        - refferen (https://echo.labstack.com/guide/)
    
    step 3 สร้าง file main.go
        - เขียน echo
        ------------------------------------------------------------
            package main

            import (
                "net/http"
                
                "github.com/labstack/echo/v4"
            )

            func main() {
                e := echo.New()
                e.GET("/", func(c echo.Context) error {
                    return c.String(http.StatusOK, "Hello, World!")
                })
                e.Logger.Fatal(e.Start(":1323"))
            }
        -------------------------------------------------------------

    step 4 run file main
        - go run main.go

------------------------------------------------------------------------------------------
# การเขียน file ใน golang มีด้วยกันหลายวิธี

    อันนี้ คือ วิธ๊การเขียน file โดยใช้ "io/ioutil" เเละ "path/filepath"

    หากเรามี file อยู่เเล้ว สามารถทำตาม code ชุดนี้ได้เลย

    // os.Getwd() คือ การเรียก root path ของ dirctory 
    path, err := os.Getwd()
	if err != nil {
	 	log.Println(err)
	}

	inputFile := filepath.Join(path, "img", <filename>)

    // ioutil.ReadFile() คือ การอ่าน file_input เเละ return ออกมาเป็น []byte
	dataFile, err := ioutil.ReadFile(inputFile)
	if err != nil {
	 	logrus.Errorln("Error: Readfile ->", err)
	}

	outPutFile := filepath.Join(path, "img", <new filename>)

    // ioutil.WriteFile() คือ การเขียน file
    if err := ioutil.WriteFile(newPathFile, dataFile, 0644); err != nil {
	    logrus.Errorln("Error: WriteFile ->", err)
	}

--------------------------------------------------------------------------------------------
# golang หาก deploy เป็น container เเล้ว