package main

import (
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"os"
	"path"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/ilovelili/LandingPageGenerator/config"
	"github.com/ilovelili/LandingPageGenerator/ftp"
	qrcode "github.com/skip2/go-qrcode"
)

var (
	cfg              *config.Config
	filename         = flag.String("filename", "", "filename to download")
	landingpageitems = make([]*LandingPageItem, 0)
)

func init() {
	confg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	cfg = confg
}

// LandingPageItem landing page entity
type LandingPageItem struct {
	Name     string `csv:"name"`
	URL      string `csv:"url"`
	Password string `csv:"password"`
	Created  string `csv:"time"`
}

func main() {
	flag.Parse()
	if *filename == "" {
		*filename = time.Now().Format("20060102") + ".csv"
	}

	ftp := &ftp.FTP{
		IP:       cfg.IP,
		Port:     cfg.Port,
		UserName: cfg.UserName,
		Password: cfg.Password,
	}

	bufchan := make(chan []byte)
	go func() {
		ftp.Download(bufchan, *filename)
	}()

	for {
		select {
		case <-time.After(time.Second * 5):
			{
				fmt.Println("Nothing to download. Bye")
				os.Exit(0)
			}

		case buf := <-bufchan:
			{
				if err := gocsv.UnmarshalBytesToCallback(buf, generateLandingPageItem); err != nil {
					panic(err)
				}

				if _, err := os.Stat("./output"); !os.IsNotExist(err) {
					os.Mkdir("./output", os.ModePerm)
				}

				outputhtml := path.Join("./output", "index.html")
				// remove if exists
				if _, err := os.Stat(outputhtml); !os.IsNotExist(err) {
					os.Remove(outputhtml)
				}

				if err := generateHTML(outputhtml); err != nil {
					panic(err)
				}

				deleted := make(chan bool)
				go ftp.Delete(deleted, *filename)

				if <-deleted {
					fmt.Println("All Done")
					os.Exit(0)
				} else {
					panic(errors.New("file not deleted"))
				}
			}
		}
	}
}

func generateLandingPageItem(lp *LandingPageItem) {
	// resolve qrcode
	qrcode, err := qrcode.Encode(lp.URL, qrcode.Highest, 256)
	if err != nil {
		// skip this item
		return
	}
	base64img := base64.StdEncoding.EncodeToString(qrcode)

	landingpageitems = append(landingpageitems, &LandingPageItem{
		Name:     lp.Name,
		URL:      base64img,
		Password: lp.Password,
		Created:  lp.Created,
	})
}

func generateHTML(output string) error {
	fmap := template.FuncMap{
		"passthrough": passThrough,
	}
	t := template.Must(template.New("index.templ").Funcs(fmap).ParseFiles("./template/index.templ"))

	file, err := os.OpenFile(output, os.O_CREATE|os.O_RDWR, os.ModePerm)
	defer file.Close()
	if err != nil {
		return err
	}

	return t.Execute(file, landingpageitems)
}

func passThrough(s string) template.URL {
	return template.URL(s)
}
