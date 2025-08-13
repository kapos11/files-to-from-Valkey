package controller

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/valkey-io/valkey-go"
)

func UploadFiles(dirPath string, valkeyHost string) {
	numOfSuccFiles := 0
	context := context.Background()
	// connect valkey
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{valkeyHost},
	})
	if err != nil {
		fmt.Println("Failed to connect Valkey : ", err)
		return
	}
	defer client.Close()

	//filepath.Walk to process all files and directories
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return err //Continue
		}
		if info.IsDir() {
			return nil // Skip directories foucs files only
		}
		content, err := os.ReadFile(path)

		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", path, err)
			return nil // Continue
		}
		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			fmt.Printf("Error getting relative path for %s: %v\n", path, err)
			return nil // Continue
		}

		//key dirName/fileName
		key := fmt.Sprintf("%s/%s", filepath.Base(dirPath), relPath)
		//SET in valkey
		err = client.Do(context, client.B().Set().Key(key).Value(string(content)).Build()).Error()
		if err != nil {
			fmt.Printf("Error uploading file %s: %v\n", path, err)
			return nil // Continue
		} else {
			numOfSuccFiles++
			fmt.Printf("Uploaded: %s\n", path)
		}
		return nil // Continue

	})
	if err != nil {
		fmt.Printf("Error walking the path %s: %v\n", dirPath, err)
		return
	}
	fmt.Println("Number of successfully uploaded files: ", numOfSuccFiles)
}

// downloadFiles from valkey

func DownloadFiles(sourceDir, targetDir string, valkeyHost string) {
	numOfSuccFiles := 0
	context := context.Background()
	//connect valkey
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{valkeyHost},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	//if directory not found
	_, err = os.Stat(targetDir)

	if os.IsNotExist(err) {
		err = os.MkdirAll(targetDir, 0755)
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

		//get REL path
		relPath, err := filepath.Rel(sourceDir, key)
		if err != nil {
			fmt.Printf("Error getting relative path for key %s: %v\n", key, err)
			continue
		}

		targetPath := filepath.Join(targetDir, relPath)

		//push content to the file and make directory if not exists
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			fmt.Printf("Error creating directory for %s: %v\n", targetPath, err)
			continue
		}
		//push content
		err = os.WriteFile(targetPath, []byte(content), 0644)
		if err != nil {
			fmt.Println("Error writing file ", err)
		} else {
			numOfSuccFiles++
			fmt.Println("Downloaded: ", targetPath)
		}
	}
	fmt.Printf(" Number of successfully downloaded : %d files\n", numOfSuccFiles)
}
