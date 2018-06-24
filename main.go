package main

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/rs/xid"
)

const cookieName = "_concater"

func main() {
	e := echo.New()
	e.Debug = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/static", "assets")
	e.File("/", "index.html")

	e.GET("/cookie", func(c echo.Context) error {
		cookie, err := c.Cookie(cookieName)
		if err != nil {
			return c.String(http.StatusOK, "")
		}
		return c.String(http.StatusOK, cookie.Value)
	})

	e.DELETE("/cookie", func(c echo.Context) error {
		c.SetCookie(&http.Cookie{
			Name:    cookieName,
			Value:   "",
			Expires: time.Now(),
		})
		return c.String(http.StatusOK, "")
	})

	e.POST("/upload", func(c echo.Context) error {
		uid := xid.New().String()
		c.SetCookie(&http.Cookie{
			Name:    cookieName,
			Value:   uid,
			Expires: time.Now().Add(time.Minute),
		})

		form, err := c.MultipartForm()
		if err != nil {
			panic(err)
		}

		dir := "assets/video/" + uid
		if err := os.MkdirAll(dir, 0777); err != nil {
			panic(err)
		}

		inputText := uid + ".txt"
		files := []string{}
		for _, v := range form.File {
			file := v[0]
			src, err := file.Open()
			if err != nil {
				panic(err)
			}
			defer src.Close()

			path := dir + "/" + file.Filename
			dst, err := os.Create(path)
			if err != nil {
				panic(err)
			}
			defer dst.Close()

			if _, err = io.Copy(dst, src); err != nil {
				panic(err)
			}
			files = append(files, "file '"+path+"'")
		}

		txt, err := os.Create(inputText)
		if err != nil {
			panic(err)
		}
		defer txt.Close()
		txt.Write(([]byte)(strings.Join(files, "\n")))

		err = exec.Command("ffmpeg", "-f", "concat", "-i", inputText, "-c", "copy", dir+".mp4").Run()
		if err != nil {
			panic(err)
		}

		if err := os.Remove(inputText); err != nil {
			panic(err)
		}
		if err := os.RemoveAll(dir); err != nil {
			panic(err)
		}

		return c.String(http.StatusOK, uid)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
