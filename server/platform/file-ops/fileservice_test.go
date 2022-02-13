package filestore

import (
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestDeleteFile(t *testing.T) {
	//create a temp file
	tempfile, e := os.Create("#prefix.txt")
	if e != nil {
		log.Fatal(e)
	}
	tempfile.Close()
	//delete it after function exit
	defer os.Remove(tempfile.Name())

	fsObj := NewFileStoreService()
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	fsObj.FileVolumne = ""
	actual, _ := fsObj.DeleteFile(c, "#prefix.txt")
	expected := "#prefix.txt successfully deleted the file"
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}
func TestListFile(t *testing.T) {
	//delete it after function exit
	fsObj := NewFileStoreService()
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	fsObj.FileVolumne = "./"

	actual, _ := fsObj.ListFiles(c)
	expected := "fileservice.gofileservice_test.go"
	if strings.Contains(actual, expected) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}
