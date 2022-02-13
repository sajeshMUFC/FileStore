# FileStore
Provide a CLI client to perform actions on the HTTP Server on file objects . 


NAME:
   File Store CLI - Let's you add , delete ,update ,get a file and search by keyword !

USAGE:
   store command [arguments...]

VERSION:
   0.0.0

COMMANDS:
     ls       List all the files
     rm       Removes the file 
     add      Upload Multiple Files
     update   Update the content of given file
     wc       Counts the occurances of given word in all the files
     


** prerequisite **
1 Docker


** STEPS TO BUILD **
```
git clone https://github.com/sajeshMUFC/FileStore
cd FileStore/server
docker build -t file_store_http_server .
docker run -p 8000:8000 file_store_http_server

//No run the client
cd FileStore/cmd
./store <commands>
```