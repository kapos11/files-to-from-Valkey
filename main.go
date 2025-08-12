package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/valkey-io/valkey-go"
)

func uploadFiles(dirPath string) {
	numOfSuccFiles := 0
	context := context.Background()
	// connect valkey
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{"localhost:6379"},
	})
	if err != nil {
		fmt.Println("Failed to connect Valkey : ", err)
		return
	}
	defer client.Close()

	//read dir
	files, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}
	//loop to upload files in valkey
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filePath := filepath.Join(dirPath, file.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", file.Name(), err)
			continue
		}
		//key dirName/fileName
		key := fmt.Sprintf("%s/%s", filepath.Base(dirPath), file.Name())

		//SET in valkey
		err = client.Do(context, client.B().Set().Key(key).Value(string(content)).Build()).Error()
		if err != nil {
			fmt.Printf("Error uploading file %s: %v\n", file.Name(), err)
			continue
		} else {
			numOfSuccFiles++
			fmt.Printf("Uploaded: %s\n", file.Name())
		}
	}
	fmt.Printf(" Number of successfully uploaded : %d files", numOfSuccFiles)
}

// downloadFiles from valkey

func downloadFiles(sourceDir, targetDir string) {
	numOfSuccFiles := 0
	context := context.Background()
	//connect valkey
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{"localhost:6379"},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	//if directory not found
	_, err = os.Stat(targetDir)

	if os.IsNotExist(err) {
		err = os.Mkdir(targetDir, 0755)
		if err != nil {
			fmt.Println("Error creating target directory:", err)
			return
		}
	}
	//get keys
	getKeys := client.Do(context, client.B().Keys().Pattern(sourceDir+"/*").Build())

	keys, err := getKeys.AsStrSlice()
	if err != nil {
		fmt.Println("Error in getting keys from Valkey:", err)
		return
	}
	if len(keys) == 0 {
		fmt.Println("No files on this directory.")
		return
	}

	//download each file
	for _, key := range keys {
		keyValue := client.Do(context, client.B().Get().Key(key).Build())
		content, err := keyValue.ToString()
		if err != nil {
			fmt.Printf("Error getting content for key %s: %v\n", key, err)
			continue
		}

		fileName := filepath.Base(key)
		targetPath := filepath.Join(targetDir, fileName)

		//push content to the file
		err = os.WriteFile(targetPath, []byte(content), 0644)
		if err != nil {
			fmt.Println("Error writing file ", err)
		} else {
			numOfSuccFiles++
			fmt.Printf("Downloaded: %s\n", fileName)
		}
	}
	fmt.Printf(" Number of successfully downloaded : %d files\n", numOfSuccFiles)
}

func main() {

	//select run command
	if len(os.Args) > 0 || len(os.Args) > 4 {
		flag := os.Args[1]
		if flag == "-u" && len(os.Args) == 3 {
			uploadFiles(os.Args[2])
			return
		} else if flag == "-d" && len(os.Args) == 4 {
			downloadFiles(os.Args[2], os.Args[3])
			return
		} else {
			fmt.Println("no arguments in this case")
			fmt.Println("upload-folder: ./copy-dir -u <directory>")
			fmt.Println("download-from-valkey : ./copy-dir -d <valkey_dir> <target_dir>")
		}

	}
}
