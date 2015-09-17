package sapt

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Package struct {
	Path     string
	Content  []byte
	Metadata *PackageMetadata
}

// TODO: refactor...
func UploadPackages(file *os.File, rm bool, conn *S3) {
	defer file.Close()
	f, _ := file.Stat()
	if f.IsDir() {
		filepath.Walk(file.Name(), func(p string, f os.FileInfo, err error) error {
			if !f.IsDir() && filepath.Ext(p) == ".deb" {
				fn, _ := os.Open(p)
				defer fn.Close()
				conn.UploadPackage(NewPackage(fn))
				if rm {
					os.Remove(p)
				}
			}
			return nil
		})
	} else if filepath.Ext(file.Name()) == ".deb" {
		conn.UploadPackage(NewPackage(file))
		if rm {
			os.Remove(file.Name())
		}
	}
}

func NewPackage(file *os.File) *Package {
	path := file.Name()

	// read in file contents
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, file)

	// retrieve metadata
	metadata := MetadataFromFile(file)
	metadata.MD5sum = Hash(buf.Bytes(), "md5")
	metadata.SHA1 = Hash(buf.Bytes(), "sha1")
	metadata.SHA256 = Hash(buf.Bytes(), "sha256")

	return &Package{
		Path:     path,
		Content:  buf.Bytes(),
		Metadata: metadata,
	}
}

func Hash(content []byte, crypto string) string {
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
