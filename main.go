package main

import (
	"io/ioutil"

	"github.com/ilovelili/LandingPageGenerator/config"
	"github.com/ilovelili/LandingPageGenerator/ftp"
)

var (
	cfg *config.Config
)

func init() {
	confg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	cfg = confg
}

func main() {
	ftp := &ftp.FTP{
		IP:       "188.166.244.244",
		Port:     "21",
		UserName: "wechat",
		Password: "Aa7059970599",
	}

	bufchan := make(chan []byte)
	go func() {
		ftp.Download(bufchan, "testftp.txt")
	}()

	// todo: I have buf now, unmarshall to csv and QR generator
	ioutil.WriteFile("testftp.txt", <-bufchan, 0640)

}

/*


// GenerateQRCodeFromURLString generate QR from url string
func GenerateQRCodeFromURLString(urlstring, outputfilename string) error {
	_, err := url.ParseRequestURI(urlstring)
	if err != nil {
		return err
	}

	return qrcode.WriteFile(urlstring, qrcode.Highest, 256, outputfilename)
}


*/
