package main

import (
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func convert(c echo.Context, inputExt string, outputExt string) error {
	file, err := c.FormFile("file")

	newDirName := uuid.New().String()

	os.MkdirAll("files/"+newDirName, os.ModePerm)

	if err != nil {
		log.Fatal(err)
		return err
	}

	src, err := file.Open()

	if err != nil {
		log.Fatal(err)
		return err
	}

	
	dst, err := os.Create(filepath.Join("files/"+newDirName, filepath.Base(file.Filename)))
	
	if err != nil {
		log.Fatal(err)
		return err
	}

	
	if _, err = io.Copy(dst, src); err != nil {
		log.Fatal(err)
		return err
	}
	
	fileNameNoExt := strings.TrimSuffix(file.Filename, inputExt)
	
	cmd := exec.Command("ffmpeg", "-i", "files/"+newDirName+"/"+file.Filename, "files/"+newDirName+"/"+fileNameNoExt+outputExt)
	
	err = cmd.Run()
	
	if err != nil {
		log.Fatal(err)
		return err
	}
	
	outputFile, err := os.ReadFile("files/"+newDirName+"/"+fileNameNoExt+outputExt)
	
	if err != nil {
		log.Fatal(err)
		return err
	}
	
	src.Close()
	dst.Close()
	if err := os.RemoveAll("files/"+newDirName); err != nil {
		log.Fatal(err)
		return err
	}

	return c.Blob(http.StatusOK, mime.TypeByExtension(outputExt), outputFile)
}

func main() {
	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to file convertion server")
	})

	e.POST("/convert/avitomp4", func(c echo.Context) error {
		return convert(c, ".avi", ".mp4")
	})

	e.POST("/convert/mp4toavi", func(c echo.Context) error {
		return convert(c, ".mp4", ".avi")
	})

	e.Logger.Fatal(e.Start(":8080"))
}