package filestore

import (
	"errors"
	"io/ioutil"
	"log"
	"os"

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

func (fs *FileStore) AddFiles(c *gin.Context) (string, error) {
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
			if _, err := os.Stat(fs.FileVolumne + file.Filename); err == nil {
				message = message + file.Filename + " already exists \n"
				continue
			}
			dst, err := os.Create(fs.FileVolumne + file.Filename)
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
			err = ioutil.WriteFile(fs.FileVolumne+file.Filename, data, 0644)
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

func (fs *FileStore) ListFiles(c *gin.Context) (string, error) {
	// Multipart form
	var fileNames string
	files, err := ioutil.ReadDir(fs.FileVolumne)
	if err != nil {
		log.Fatal("====", err)
	}
	for _, file := range files {
		fileNames = fileNames + " \n " + file.Name()
	}
	return fileNames, nil

}

func (fs *FileStore) DeleteFile(c *gin.Context, fname string) (string, error) {
	// Multipart form
	e := os.Remove(fs.FileVolumne + fname)
	if e != nil {
		return fname + " No such file found", errors.New("Failed to delete")
	}
	return fname + " successfully deleted the file", nil
}

func (fs *FileStore) UpdateFile(c *gin.Context) (string, error) {
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
			if _, err := os.Stat(fs.FileVolumne + file.Filename); err == nil {
				//read the file
				freader, err := os.OpenFile(fs.FileVolumne+file.Filename, os.O_RDWR, 0644)
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
			err = ioutil.WriteFile(fs.FileVolumne+file.Filename, data, 0644)
			if err != nil {
				return file.Filename + " update failed \n", errors.New("update failed")

			}
			return file.Filename + "updated successfully\n", nil

		}
	}
	return "bad request", err
}
