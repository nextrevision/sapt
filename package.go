package sapt

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Package struct {
	Path     string
	Content  []byte
	Metadata *PackageMetadata
}

func UploadPackages(file *os.File, rm bool, conn *S3) {
	stat, _ := file.Stat()
	path := file.Name()
	fileList := []string{}
	if stat.IsDir() {
		filepath.Walk(path, func(p string, f os.FileInfo, err error) error {
			if !f.IsDir() && filepath.Ext(p) == ".deb" {
				fileList = append(fileList, p)
			}
			return nil
		})
	} else if filepath.Ext(path) == ".deb" {
		fileList = append(fileList, path)
	}
	file.Close()

	var wg sync.WaitGroup
	wg.Add(len(fileList))
	uploadChan := make(chan *Package, 10)
	for _, f := range fileList {
		go func(f string) {
			uploadChan <- NewPackage(f, path)
		}(f)
	}
	for i := 0; i < len(fileList); i++ {
		go func() {
			for {
				pkg, ok := <-uploadChan
				if !ok {
					break
				}
				conn.uploadPackage(pkg)
				time.Sleep(time.Millisecond * 100)
				wg.Done()
			}
		}()
	}
	wg.Wait()
}

func NewPackage(path string, basePath string) *Package {
	var name string
	if path == basePath {
		name = basename(path)
	} else {
		name = strings.Replace(path, basePath, "", 1)
	}

	// read in file contents
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	// retrieve metadata
	metadata := MetadataFromFile(path)
	metadata.Filename = name
	metadata.MD5sum = hash(content, "md5")
	metadata.SHA1 = hash(content, "sha1")
	metadata.SHA256 = hash(content, "sha256")

	return &Package{
		Path:     name,
		Content:  content,
		Metadata: metadata,
	}
}

func hash(content []byte, crypto string) string {
	var hash string

	if crypto == "md5" {
		hash = fmt.Sprintf("%x", md5.Sum(content))
	} else if crypto == "sha1" {
		hash = fmt.Sprintf("%x", sha1.Sum(content))
	} else if crypto == "sha256" {
		hash = fmt.Sprintf("%x", sha256.Sum256(content))
	}

	return hash
}

func basename(path string) string {
	i := strings.LastIndex(path, "/")
	return path[i:]
}
