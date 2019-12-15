package main

import (
	"image/png"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"

	"github.com/fogleman/gg"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
)

var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func (p *program) run() {
	gin.SetMode(gin.ReleaseMode)
	f, _ := os.Create("./gin.log")
	gin.DefaultWriter = io.MultiWriter(f)
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	router := gin.Default()
	//new template engine
	router.HTMLRender = ginview.Default()

	router.GET("/", indexHandler)
	router.POST("/create", formHandler)

	router.Run(":80")
}

func main() {
	svcConfig := &service.Config{
		Name:        "QR-code-generator",
		DisplayName: "QR code generator",
		Description: "",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}

func indexHandler(ctx *gin.Context) {
	//var err error

	ctx.HTML(http.StatusOK, "index", gin.H{
		"Title": "QR Code Generator"},
	)

}

func formHandler(ctx *gin.Context) {

	dataString := ctx.PostForm("dataString")

	qrCode, _ := qr.Encode(dataString, qr.L, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 1024, 1024)

	//png.Encode(ctx.Writer, qrCode)

	im := qrCode

	dc := gg.NewContext(1024, 1050)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace("./Arial.ttf", 16); err != nil {
		panic(err)
	}
	dc.DrawStringAnchored(dataString, 512, 1045, 0.5, 0)

	dc.DrawRoundedRectangle(0, 0, 1024, 1050, 0)
	dc.DrawImage(im, 0, 0)
	dc.DrawStringAnchored(dataString, 512, 1045, 0.5, 0)
	dc.Clip()
	//dc.SavePNG("./out.png")

	png.Encode(ctx.Writer, dc.Image())

	ctx.String(http.StatusOK, "Done")
}
