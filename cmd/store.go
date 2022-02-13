package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

const (
	put    = "PUT"
	post   = "POST"
	get    = "GET"
	delete = "DELETE"
)

var serverAddr string

func main() {
	// Subcommands

	serverAddr = os.Getenv("SERVER_URL")
	if serverAddr == "" {
		serverAddr = "http://localhost:8000/"

	}
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	listCmd := flag.NewFlagSet("ls", flag.ExitOnError)
	rmCmd := flag.NewFlagSet("remove", flag.ExitOnError)
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	searchCmd := flag.NewFlagSet("wc", flag.ExitOnError)
	freqCount := flag.NewFlagSet("freq-words", flag.ExitOnError)
	//wordCountSubCmd := freqCount.String("n", "", "order", "")
	limitPtr := freqCount.String("n", "query limit", "limit")
	SortPtr := freqCount.String("order", "sort order", "limit")

	// Verify that a subcommand has been provided
	// os.Arg[0] is the main command
	if len(os.Args) < 2 {
		fmt.Println("Either of add/ls/remove/update command is required")
		os.Exit(0)
	}

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		uploadFiles(os.Args[2:])
	case "ls":
		listCmd.Parse(os.Args[2:])
		getAllFiles()
	case "rm":
		rmCmd.Parse(os.Args[2:])
		deleteFile(os.Args[2:])
	case "update":
		updateCmd.Parse(os.Args[2:])
		updateFile(os.Args[2:])
	case "wc":
		searchCmd.Parse(os.Args[2:])
		getWordCount(os.Args[2:])
	case "freq-words":
		freqCount.Parse(os.Args[2:])

	default:
		fmt.Println("No command found")
		os.Exit(0)
	}

	if freqCount.Parsed() {
		if *limitPtr == "" {
			fmt.Println("specify the limits, --n")
			os.Exit(0)
		}
		getFreqWordCount(*limitPtr, *SortPtr)
	}
}
func getFreqWordCount(limit string, sort string) {
	q := url.Values{}
	q.Add("limit", limit)
	q.Add("sort", sort)
	res, err := commonHttpRequest(get, serverAddr+"v1/file/freqword", q, nil)
	if err != nil {
		fmt.Println("unable to connect to file store server", err)
	}
	fmt.Println(string(res))

}

func uploadFiles(filenames []string) {
	var paths []string
	if len(filenames) <= 0 {
		fmt.Println("Atleast one file is required.")
	} else {
		for _, file := range filenames {
			path, _ := os.Getwd()
			path = path + "/" + file
			paths = append(paths, path)
		}
		res, err := multipleFileUploadRequest(post, serverAddr+"v1/file/", "add_file", paths)
		if err != nil {
			fmt.Println(" connection error")
		}
		fmt.Println(string(res))
	}
}

func updateFile(filenames []string) {
	var paths []string
	if len(filenames) <= 0 {
		fmt.Println("Atleast one file is required.")
	} else if len(filenames) > 1 {
		fmt.Println("multiple updates is not supported.")
	} else {
		for _, file := range filenames {
			path, _ := os.Getwd()
			path = path + "/" + file
			paths = append(paths, path)
		}
		res, err := multipleFileUploadRequest(put, serverAddr+"v1/file/", "update_file", paths)
		if err != nil {
			fmt.Println(" connection error")
		}
		fmt.Println(string(res))
	}
}

// Multiple file upload
func multipleFileUploadRequest(mtype string, uri string, paramName string, paths []string) ([]byte, error) {
	var err error
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			fmt.Print("unable to open the file, provide the right path")
			os.Exit(0)
		}
		defer file.Close()
		part, err := writer.CreateFormFile(paramName, filepath.Base(path))
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(part, file)
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(mtype, uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return responseBody, err
}

func deleteFile(filenames []string) {
	if len(filenames) <= 0 {
		fmt.Println("file name to delete is missing.")
	} else if len(filenames) > 1 {
		fmt.Println("multiple delete is not supported.")
	} else {
		res, err := commonHttpRequest(delete, serverAddr+"v1/file/"+filenames[0], nil, nil)
		if err != nil {
			fmt.Println("unable to connect to file store server")
		}
		fmt.Println(string(res))
	}
}

func getAllFiles() {
	res, err := commonHttpRequest(get, serverAddr+"v1/file", nil, nil)
	if err != nil {
		fmt.Println("unable to connect to file store server")
	}
	fmt.Println(string(res))
}

func getWordCount(words []string) {
	if len(words) <= 0 {
		fmt.Println("word is required for getting the count")
	} else if len(words) > 1 {
		fmt.Println("multiple word count is not supported.")
	} else {
		res, err := commonHttpRequest(get, serverAddr+"v1/file/"+words[0], nil, nil)
		if err != nil {
			fmt.Println("unable to connect to file store server")
		}
		fmt.Println(string(res))
	}
}

func commonHttpRequest(method string, url string, q url.Values, body []byte) ([]byte, error) {
	var req *http.Request
	var err error
	if method == put {
		req, err = http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	} else if method == post {
		req, err = http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
	} else if method == get {
		req, err = http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Content-Type", "application/json")
	} else if method == delete {
		req, err = http.NewRequest(http.MethodDelete, url, nil)
	}
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//Add the result in goroutine channel
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return responseBody, err
}
