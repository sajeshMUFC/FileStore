package handler

import (
	filestore "http-filestore/platform/file-ops"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FileUploadHandler(fs *filestore.FileStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := fs.AddFiles(c)
		if err != nil {
			c.String(http.StatusInternalServerError, res)
		} else {
			c.String(http.StatusOK, res)
		}
	}
}

func FileListHandler(fs *filestore.FileStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := fs.ListFiles(c)
		if err != nil {
			c.String(http.StatusInternalServerError, res)
		} else {
			c.String(http.StatusOK, res)
		}
	}
}

func FileDeleteHandler(fs *filestore.FileStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		fname := c.Param("filename")
		res, err := fs.DeleteFile(c, fname)
		if err != nil {
			c.String(http.StatusInternalServerError, res)
		} else {
			c.String(http.StatusOK, res)
		}
	}
}

func FileUpdateHandler(fs *filestore.FileStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := fs.UpdateFile(c)
		if err != nil {
			c.String(http.StatusInternalServerError, res)
		} else {
			c.String(http.StatusOK, res)
		}
	}
}

func FileWordCountHandler(fs *filestore.FileStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		word := c.Param("search")
		res, err := fs.WordCountInFiles(c, word)
		if err != nil {
			c.String(http.StatusInternalServerError, res)
		} else {
			c.String(http.StatusOK, res)
		}
	}
}
