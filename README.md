# FileStore
Provide a CLI client to perform actions on the HTTP Server on file objects . 

## CLI Details
### Name : FileStore
- Description : File Store CLI - Let's you add , delete ,update ,get a file and search by keyword !
- USAGE: store command [arguments...]
- ### COMMANDS: 
            - ls            List all the files
            - rm            Removes the file
            - add           Upload Multiple Files
            - update        Update the content of given file
            - wc            Counts the occurances of given word in all the files
            - freq-words    [--n <limit> order <desc>]Displays the most frequent words/less frequent used words based on the limit



## PREREQUISITE 
- Docker


## STEPS TO BUILD 
```
git clone https://github.com/sajeshMUFC/FileStore
cd FileStore/server
docker build -t file_store_http_server .
docker run -p 8000:8000 file_store_http_server

//run the client
cd FileStore/cmd
./store <commands>
```


## REFERENCES

Go CLI -  https://www.rapid7.com/blog/post/2016/08/04/build-a-simple-cli-tool-with-golang/
Go GIN Framework - https://github.com/gin-gonic/gin