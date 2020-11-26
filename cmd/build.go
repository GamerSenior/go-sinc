/*
Copyright © 2020 Giovani Garcia Abel <abelgiovani@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/GamerSenior/go-sinc/internal/ftp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var buildCore bool
var upload bool

func runMavenCleanInstall() error {
	mvnCmd := exec.Command("mvn", "clean", "install")
	var stdoutBuf, stderrBuf bytes.Buffer
	mvnCmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	mvnCmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err := mvnCmd.Run()
	if err != nil {
		return errors.New("Erro durante compilação do projeto")
	}

	outStr := string(stdoutBuf.Bytes())
	fmt.Printf("\nout:\n%s", outStr)
	return nil
}

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the project especified",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("É necessário um módulo para realizar a build")
		}
		module := args[0]
		sincDir := viper.GetString("JAVA_SINC_DIR")
		if _, err := os.Stat(sincDir + "/" + module); os.IsNotExist(err) {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		sincDir := viper.GetString("JAVA_SINC_DIR")

		if buildCore {
			err := os.Chdir(sincDir + "/sinc-core")
			if err != nil {
				fmt.Println("Não foi possível acessar JAVA_SINC_DIR")
				return
			}

			if err := runMavenCleanInstall(); err != nil {
				log.Println(err)
				return
			}
		}

		modulePath := sincDir + "/" + args[0]
		err := os.Chdir(modulePath)
		if err != nil {
			fmt.Println("Erro ao acessar diretório " + modulePath)
			return
		}
		if err := runMavenCleanInstall(); err != nil {
			log.Println(err)
			return
		}

		if upload {
			filePath := modulePath + "/war/target/" + args[0] + ".war"
			if err := ftp.SendToFTP(filePath, Verbose); err != nil {
				log.Printf("Error ao fazer upload do war: %s", err)
				return
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	buildCmd.Flags().BoolVarP(&buildCore, "core", "c", false, "Compila o módulo Core da aplicação")
	buildCmd.Flags().BoolVarP(&upload, "upload", "u", false, "Realiza upload do módulo compilado no FTP configurado")
}
