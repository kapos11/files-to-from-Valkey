# files-to-from-Valkey

This project, written in **Go Lang**, allows you to **upload** files from a local directory to a **Valkey** in-memory database,  
and **download** files stored in Valkey back to your local machine.

##  Features
- **Upload** files from a local folder to Valkey.
- **Download** files from Valkey to a local folder.

##  Requirements
- **Go** (version 1.18 or higher).
- Running **Valkey Server** locally or remotely.
- Connection to Valkey at `localhost:6379` (modify the code if using a different host/port).
- Install Valkey Go client:
    go get github.com/valkey-io/valkey-go

## Build the project:

- go build -o copy-dir

## HOW TO USE
1️- Upload files from a local directory to Valkey

    ./copy-dir -u <directory_path>

     Example:

    ./copy-dir -u ./photos

2- Download files from Valkey to a local directory

    ./copy-dir -d <valkey_directory> <target_directory>

    Example:

    ./copy-dir -d photos backup_photos