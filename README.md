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

