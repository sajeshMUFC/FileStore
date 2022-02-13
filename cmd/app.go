package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	// we call ioutil.TempFile which returns either a file
	// or an error.
	// we specify the directory we want to create these temp files in
	// for this example we'll use `car-images`, and we'll define
	// a pattern which will be used for naming our car images
	// in this case car-*.png
	file, err := ioutil.TempFile("./car-images/", "car-*.png")
	if err != nil {
		fmt.Println(err)
	}
	// We can choose to have these files deleted on program close
	defer os.Remove(file.Name())
	// We can then have a look and see the name
	// of the image that has been generated for us
	fmt.Println(file.Name())
}
