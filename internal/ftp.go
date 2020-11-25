package ftp

import (
	"log"
	"os"

	"github.com/jlaffaye/ftp"
	"github.com/spf13/viper"
)

func sendToFTP(path string) error {
	client, err := ftp.Dial("ftp.supportweb.com.br")
	if err != nil {
		return err
	}

	user := viper.GetString("FTP_USER")
	pass := viper.GetString("FTP_PASS")
	dir := viper.GetString("FTP_DIR")

	if err := client.Login(user, pass); err != nil {
		return err
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	if len(dir) > 0 {
		if err := client.ChangeDir(dir); err != nil {
			return err
		}
	}

	err = client.Stor(f.Name(), f)
	if err != nil {
		return err
	}

	if err := client.Quit(); err != nil {
		log.Fatal(err)
	}
	return nil
}
