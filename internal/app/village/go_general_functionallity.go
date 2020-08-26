package village

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

// This method iterates over the Attributes of an Struct, and print all the information about the structure
func shortDescriptionInterface(in interface{}) {
	var date int64

	if reflect.ValueOf(in).Kind() == reflect.Struct {
		fmt.Println("Is a struct")
		nameFunction := reflect.TypeOf(in).Name()
		fmt.Println("Struct name: " + nameFunction)
		for i := 0; i < reflect.ValueOf(in).NumField(); i++ {
			nameAtribute := reflect.TypeOf(in).Field(i).Name
			fmt.Println("Field name: " + nameAtribute)
			attribute := reflect.ValueOf(in).Field(i).String()
			fmt.Println("Value: " + attribute)

			if reflect.TypeOf(in).Field(i).Type == reflect.TypeOf(date) {
				// Is an int64
				valueNewAttribute, _ := reflect.ValueOf(in).Field(i).Interface().(int64)
				fmt.Print("Value: ")
				fmt.Println(valueNewAttribute)
			}
		}
	}
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func unzipDirectory(src string, dest string) ([]string, error) {
	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil

}
func zipDirectory(directoryToZIP string, dest string, nameZip string, directoryToCut string) {
	files, err := listFiles(directoryToZIP)
	if err != nil {
		panic(err)
	}
	zipMe(files, dest+nameZip+".zip", directoryToCut)
	for _, f := range files {
		fmt.Println(f)
	}
	fmt.Println("Done!")
}

// List all files inside the directory given
func listFiles(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil

}

// Zips all files given to the directory target given.
// target SHOULD EN IN .ZIP and have the full path
// diractoryToCut is the last directory it should remain, otherwise will create the whole path structure of folder from the system
func zipMe(filepaths []string, target string, directoryToCut string) error {

	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(target, flags, 0644)

	if err != nil {
		return fmt.Errorf("Failed to open zip for writing: %s", err)
	}
	defer file.Close()

	zipw := zip.NewWriter(file)
	defer zipw.Close()

	for _, filepath := range filepaths {
		last := strings.LastIndex(filepath, directoryToCut)
		filename := filepath[last:len(filepath)]
		fmt.Print("The filename is " + filename)
		if err := addFileToZip(filepath, filename, zipw); err != nil {
			return fmt.Errorf("Failed to add file %s to zip: %s", filepath, err)
		}
	}
	return nil

}

// Adds a new file to the ZIP
func addFileToZip(filepath string, filename string, zipw *zip.Writer) error {
	file, err := os.Open(filepath)

	if err != nil {
		return fmt.Errorf("Error opening file %s: %s", filepath, err)
	}
	defer file.Close()

	wr, err := zipw.Create(filename)
	if err != nil {
		return fmt.Errorf("Error adding file; '%s' to zip : %s", filename, err)
	}

	if _, err := io.Copy(wr, file); err != nil {
		return fmt.Errorf("Error writing %s to zip: %s", filename, err)
	}

	return nil
}

// copyFile copies a single file from src to dst, taking care of creating the destiny directory if does not exist
// The src and the dst should be the FULL ROUTE TO THE FILE, INCLUDING THE EXTENSION
func copyFile(src, dst string) error {
	// First create the directory just in case is not there
	last := strings.LastIndex(dst, "/")
	os.MkdirAll(dst[0:last+1], os.ModePerm)

	fmt.Println("Copy " + src)
	fmt.Println("Into " + dst[0:last+1])

	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	fmt.Println("Copy file finished executing")

	return os.Chmod(dst, srcinfo.Mode())
}

// copyDirectory copies a whole directory recursively, taking care of creating the destiny directory if does not exist
func copyDirectory(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = copyDirectory(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = copyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

// This method orders in alphabetic order the array of Items
func bubbleSort(animals []Item) {
	for i := 0; i < len(animals); i++ {
		sweep(animals)
	}
}

// This method orders in alphabetic order the array of Items
func sweep(animals []Item) {
	N := len(animals)
	firstIndex := 0
	secondIndex := 1

	for secondIndex < N {
		firstVal := animals[firstIndex]
		secondVal := animals[secondIndex]

		var a string
		a = firstVal.getName()
		b := secondVal.getName()

		if greater(a, b) {
			animals[firstIndex] = secondVal
			animals[secondIndex] = firstVal
		}

		firstIndex++
		secondIndex++
	}
}

// Compares the two strings given and returns true if a is greater than b, is INSENSITIVE to capital letters
func greater(a, b string) bool {
	if strings.ToLower(a) > strings.ToLower(b) {
		return true
	}
	return false
}

// Returns the name of the Item to be able to order alphabetical
func (i Item) getName() string {
	return i.Name
}

// Returns the ID of the Item to be able to order alphabetical
func (i Item) getID() bson.ObjectId {
	return i.ID
}

// Reverts the order of the array SALES
// Go has no generics, so can not convert array X into array interface.
func revertSales(input []Sale) []Sale {
	if len(input) == 0 {
		return input
	}
	return append(revertSales(input[1:]), input[0])
}
