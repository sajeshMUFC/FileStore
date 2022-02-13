package filestore

import (
	"bufio"
	"errors"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type Service interface {
	AddFile() (string, error)
}

type FileStore struct {
	FileVolumne string
}
type Response struct {
	status  string
	message string
}

// Init the file mount
func NewFileStoreService() *FileStore {
	fsObj := FileStore{FileVolumne: "./saved/"}
	return &fsObj
}

func (fst *FileStore) AddFiles(c *gin.Context) (string, error) {
	// Multipart form
	var message string
	var err error
	form, err := c.MultipartForm()
	if err != nil {
		log.Println("error while reading form", err)
		return "bad request", err
	}
	if files, found := form.File["add_file"]; found {
		for _, file := range files {
			if _, err := os.Stat(fst.FileVolumne + file.Filename); err == nil {
				message = message + file.Filename + " already exists \n"
				continue
			}
			dst, err := os.Create(fst.FileVolumne + file.Filename)
			defer dst.Close()
			fileContent, err := file.Open()
			if err != nil {
				log.Println("ERROR: " + err.Error())
				message = message + file.Filename + "issue while uploading"
				continue
			}
			data, err := ioutil.ReadAll(fileContent)
			fileContent.Close()
			if err != nil {
				log.Println("ERROR: " + err.Error())
				message = message + file.Filename + " issue while uploading\n"
				continue
			}
			// Copy the uploaded file to the created file on the filesystem
			err = ioutil.WriteFile(fst.FileVolumne+file.Filename, data, 0644)
			if err != nil {
				message = message + file.Filename + " issue while saving \n"
				continue
			}
			message = message + file.Filename + "uploaded successfully\n"
		}
		return message, nil

	}
	log.Println("error while reading formfile")
	return "No files uploaded", errors.New("Invalid file")

}

func (fst *FileStore) ListFiles(c *gin.Context) (string, error) {
	// Multipart form
	var fileNames string
	files, err := ioutil.ReadDir(fst.FileVolumne)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fileNames = fileNames + " \n " + file.Name()
	}
	return fileNames, nil

}

func (fst *FileStore) DeleteFile(c *gin.Context, fname string) (string, error) {
	// Multipart form
	e := os.Remove(fst.FileVolumne + fname)
	if e != nil {
		return fname + " No such file found", errors.New("Failed to delete")
	}
	return fname + " successfully deleted the file", nil
}

func (fst *FileStore) UpdateFile(c *gin.Context) (string, error) {
	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		log.Println("error while reading form", err)
		return "bad request", err
	}
	if files, found := form.File["update_file"]; found {
		for _, file := range files {
			fileContent, err := file.Open()
			if err != nil {
				return file.Filename + "failed to open", errors.New("failed to open ")
			}
			data, err := ioutil.ReadAll(fileContent)
			if err != nil {
				return file.Filename + "failed to read", errors.New("failed to read ")
			}
			fileContent.Close()
			//if file exists
			if _, err := os.Stat(fst.FileVolumne + file.Filename); err == nil {
				//read the file
				freader, err := os.OpenFile(fst.FileVolumne+file.Filename, os.O_RDWR, 0644)
				if err != nil {
					log.Fatalf("failed opening file: %s", err)
					return file.Filename + "failed to update", errors.New("failed to update")
				}
				defer freader.Close()
				//write to the file
				_, err = freader.Write(data)
				if err != nil {
					log.Fatalf("failed writing to file: %s", err)
					return "failed writing to the file " + file.Filename, errors.New("failed to update")
				}
				return file.Filename + "Updated successfully ", nil

			}
			err = ioutil.WriteFile(fst.FileVolumne+file.Filename, data, 0644)
			if err != nil {
				return file.Filename + " update failed \n", errors.New("update failed")

			}
			return file.Filename + "updated successfully\n", nil

		}
	}
	return "bad request", err
}

func (fst *FileStore) WordCountInFiles(c *gin.Context, word string) (string, error) {
	files, err := ioutil.ReadDir(fst.FileVolumne)
	if err != nil {
		log.Fatal(err)
		return "No files found", errors.New("No file")
	}
	totalCount := 0
	fileWordCountCh := make(chan int)
	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go fst.getCountByWord(word, file, fileWordCountCh, &wg)
	}
	// close the channel in the background
	go func() {
		wg.Wait()
		close(fileWordCountCh)
	}()
	// read from channel as they come in until its closed
	for countRes := range fileWordCountCh {
		totalCount = totalCount + countRes
	}
	log.Println("totalCount: ", totalCount)
	return strconv.Itoa(totalCount), nil

}

func (fst *FileStore) getCountByWord(word string, file fs.FileInfo, ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	f, err := os.Open(fst.FileVolumne + file.Name())
	if err != nil {
		log.Println("err: ", err)
		ch <- 0
	}
	defer f.Close()
	count := 0
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		wordFromFile := scanner.Text()
		if strings.ToLower(wordFromFile) == strings.ToLower(word) {
			count++
		}
	}
	ch <- count
}
