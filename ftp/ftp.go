package ftp

import (
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/jlaffaye/ftp"
)

var (
	client *ftp.ServerConn
	once   sync.Once
)

// FTP ftp type
type FTP struct {
	IP       string
	Port     string
	UserName string
	Password string
}

// client init singleton client
func (f *FTP) client() (*ftp.ServerConn, error) {
	var client *ftp.ServerConn
	var err error
	once.Do(func() {
		client, err = ftp.Dial(f.connectionstring())
	})

	return client, err
}

// connectionstring resolve connection string
func (f *FTP) connectionstring() string {
	return fmt.Sprintf("%v:%v", f.IP, f.Port)
}

// Download get file from ftp
func (f *FTP) Download(bufchan chan<- ([]byte), filenames ...string) (err error) {
	client, err := f.client()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = client.Login(f.UserName, f.Password); err != nil {
		return
	}

	for _, entry := range filenames {
		reader, err := client.RetrFrom(entry, 0)
		if err != nil {
			return err
		}

		defer reader.Close()

		buf, err := ioutil.ReadAll(reader)
		bufchan <- buf

		return err

	}

	return
}
