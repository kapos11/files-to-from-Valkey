package main

import (
	"fmt"
	"log"
	"myproject/controller"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	// Load env variables
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	valkeyHost := os.Getenv("VALKEY_HOST")

	//select run command
	if len(os.Args) > 0 || len(os.Args) > 4 {
		flag := os.Args[1]
		if flag == "-u" && len(os.Args) == 3 {
			controller.UploadFiles(os.Args[2], valkeyHost)
			return
		} else if flag == "-d" && len(os.Args) == 4 {
			controller.DownloadFiles(os.Args[2], os.Args[3], valkeyHost)
			return
		} else {
			fmt.Println("no arguments in this case")
			fmt.Println("upload-folder: ./copy-dir -u <directory>")
			fmt.Println("download-from-valkey : ./copy-dir -d <valkey_dir> <target_dir>")
		}

	}
}
