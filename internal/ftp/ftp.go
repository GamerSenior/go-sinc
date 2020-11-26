package ftp

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/jlaffaye/ftp"
	"github.com/spf13/viper"
)

//SendToFTP recebe dois parâmetros
//path: caminho do arquivo a ser enviado
//verbose: output detalhado  para depuração
//realiza o upload do arquivo para o servidor FTP configurado
func SendToFTP(path string, verbose bool) error {

	url := viper.GetString("FTP_URL")
	if len(url) == 0 {
		url = "ftp.supportweb.com.br"
	}
	port := viper.GetString("FTP_PORT")
	if len(port) > 0 {
		url += (":" + port)
	} else {
		url += ":21"
	}

	if verbose {
		log.Printf("Conectando-se ao servidor ftp %s", url)
	}

	client, err := ftp.Dial(url)
	if err != nil {
		return err
	}

	user := viper.GetString("FTP_USER")
	pass := viper.GetString("FTP_PASS")
	dir := viper.GetString("FTP_DIR")

	if len(user) == 0 || len(pass) == 0 {
		return errors.New("Variáveis de ambiente do FTP não parametrizadas")
	}

	if err := client.Login(user, pass); err != nil {
		return err
	}

	if verbose {
		log.Printf("Login realizado com sucesso")
	}

	f, err := os.Open(path)
	if err != nil {
		log.Printf("[FTP] Open: erro ao abrir arquivo")
		return err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		log.Printf("[FTP] Stat: erro ao abrir arquivo")
		return err
	}

	if verbose {
		log.Printf("File name: " + stat.Name())
		log.Printf("Size: " + strconv.FormatInt(stat.Size(), 10))
		log.Printf("Modification: " + stat.ModTime().String())
		log.Printf("File mode: %s %04o \n", stat.Mode(), stat.Mode().Perm())
	}

	if len(dir) > 0 {
		if verbose {
			log.Printf("Acessando diretório %s", dir)
		}
		if err := client.ChangeDir(dir); err != nil {
			return err
		}
	}

	err = client.Stor(stat.Name(), f)
	if err != nil {
		log.Printf("[FTP] Stor: erro ao salvar arquivo")
		return err
	}

	if err := client.Quit(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Upload realizado com sucesso")
	return nil
}
