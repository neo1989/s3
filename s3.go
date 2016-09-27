package main

import (
	"fmt"
	"github.com/gosexy/to"
	"github.com/gosexy/yaml"
	"github.com/kataras/iris"
	"github.com/satori/go.uuid"
	"io"
	"os"
	"strings"
	"time"
)

type Response struct {
	Status   bool   `json:"status"`
	Url      string `json:"url"`
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Path     string `json:"path"`
	Filename string `json:"filename"`
}

type ErrorResponst struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}

func loadConfig() *yaml.Yaml {

	configFile := "./s3.yaml"

	configs, err := yaml.Open(configFile)

	if err != nil {
		panic("s3.yaml error")
	}

	return configs
}

func main() {

	configs := loadConfig()
	FILESERVER := to.String(configs.Get("FILESERVER"))
	UPLOADSERVER := to.String(configs.Get("UPLOADSERVER"))

	fileServer := iris.New()
	fileServer.StaticServe("./uploads", "")

	fileServer.OnError(iris.StatusNotFound, func(ctx *iris.Context) {
		ctx.Write("CUSTOM 404 NOT FOUND ERROR PAGE")
	})

	go fileServer.Listen(FILESERVER)

	iris.New()
	iris.Get("/", view)
	iris.Post("/", uploadHandler)

	mainConfig := iris.ServerConfiguration{ListeningAddr: UPLOADSERVER, MaxRequestBodySize: 32 << 20}
	iris.ListenTo(mainConfig)
}

func uploadHandler(c *iris.Context) {

	configs := loadConfig()
	FILEACCESSHOST := to.String(configs.Get("FILEACCESSHOST"))
	PROTOCOL := to.String(configs.Get("PROTOCOL"))

	handler, err := c.FormFile("file")
	folder := c.FormValueString("type")
	if folder == "" {
		folder = "default"
	}
	if err == nil {

		source, _ := handler.Open()
		defer source.Close()

		date := time.Now().Format("20060102")
		uploadRoot := fmt.Sprintf("%s/%s/%s", "uploads", folder, date)
		if _, err := os.Stat(uploadRoot); os.IsNotExist(err) {
			os.MkdirAll(uploadRoot, 0777)
		}

		st := strings.Split(handler.Filename, ".")
		filename := fmt.Sprintf("%s.%s", uuid.NewV4(), st[len(st)-1])

		dst, err := os.OpenFile(fmt.Sprintf("%s/%s", uploadRoot, filename), os.O_WRONLY|os.O_CREATE, 0666)
		defer dst.Close()

		if err == nil {
			io.Copy(dst, source)
			path := fmt.Sprintf("%s/%s/%s", folder, date, filename)
			url := fmt.Sprintf("%s://%s/%s", PROTOCOL, FILEACCESSHOST, path)
			c.JSON(200, Response{Status: true, Url: url, Protocol: PROTOCOL, Host: FILEACCESSHOST, Path: path, Filename: filename})
		} else {
			c.JSON(200, ErrorResponst{Status: false, Error: err.Error()})
		}

	} else {

		c.JSON(200, ErrorResponst{Status: false, Error: "file error!"})
	}
}

func view(c *iris.Context) {
	c.HTML(200, `
		<html>
			<head>
				<title>Wizcloud s3</title>
			</head>
			<body>
				<form action="/" method="POST" enctype="multipart/form-data">
					<input type="file" name="file">
					<input type="submit" value="Submit" />
				</form>
			</body>
		</html>
	`)

}
