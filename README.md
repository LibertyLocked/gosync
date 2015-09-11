# gosync
Gosync is a tool to quickly sync files inside a folder with those on a remote server


# Features
- Syncs all the files in client's folder with the files in server's folder (Think of it as a one-way Google Drive)
- Supports AES encryption for file transfers
- Uses SHA-1 checksum to ensure the integrity of files


## How to use: server
#### Basic server: 
gosync -s [port]

Example: gosync -s 9999

#### With AES encryption:
gosync -s [port] -key:[AESKey]

Example: gosync -s 9999 -key:myAwesomeKey


## How to use: client
#### Basic client: 
gosync -c [address:port]

Example: gosync -c localhost:9999

#### With AES encryption:
gosync -s [address:port] -key:[AESKey]

Example: gosync -c localhost:9999 -key:myAwesomeKey

#### Removing out-of-sync local files:
gosync -s [address:port] -rm

Example: gosync -c localhost:9999 -rm
