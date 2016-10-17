package sapt

import (
	"bytes"
	"compress/gzip"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"sync"
	"text/template"

	"github.com/goamz/goamz/s3"
)

const packagesTemplate = `{{range .}}Package: {{.Package}}
Version: {{.Version}}
Filename: {{.Filename}}
MD5sum: {{.MD5sum}}
SHA1: {{.SHA1}}
SHA256: {{.SHA256}}
Size: {{.Size}}
Installed-Size: {{.InstalledSize}}
{{ if .Architecture }}{{ printf "Architecture: %s\n" .Architecture }}{{ end -}}
{{ if .Depends }}{{ printf "Depends: %s\n" .Depends }}{{ end -}}
{{ if .Recommends }}{{ printf "Recommends: %s\n" .Recommends }}{{ end -}}
{{ if .Conflicts }}{{ printf "Conflicts: %s\n" .Conflicts }}{{ end -}}
{{ if .License }}{{ printf "License: %s\n" .License}}{{ end -}}
{{ if .Vendor }}{{ printf "Vendor: %s\n" .Vendor }}{{ end -}}
{{ if .Maintainer }}{{ printf "Maintainer: %s\n" .Maintainer }}{{ end -}}
{{ if .Section }}{{ printf "Section: %s\n" .Section }}{{ end -}}
{{ if .Priority }}{{ printf "Priority: %s\n" .Priority }}{{ end -}}
{{ if .Homepage }}{{ printf "Homepage: %s\n" .Homepage }}{{ end -}}
{{ if .Description }}{{ printf "Description: %s\n" .Description }}{{ end -}}
{{ printf "\n" }}{{end}}`

type Index struct {
	Path    string
	Content []byte
}

func ScanBucketPackages(conn *S3) {
	packages := []packageMetadata{}
	contents := conn.getBucketContents()
	packageList := getBucketPackages(contents)

	var wg sync.WaitGroup
	wg.Add(len(packageList))
	headerChan := make(chan http.Header, 10)
	for _, pkg := range packageList {
		go func(pkg string) {
			headerChan <- conn.getObjectHeaders(pkg)
		}(pkg)
	}
	for i := 0; i < len(packageList); i++ {
		go func() {
			for {
				headers, ok := <-headerChan
				if !ok {
					break
				}
				m := metadataFromHeaders(headers)
				packages = append(packages, *m)
				wg.Done()
			}
		}()
	}
	wg.Wait()

	packageIndex := createPackageIndex(packages)

	indices := getIndexPaths(contents)
	if len(indices) == 0 {
		indices = append(indices, Index{Path: "repo"})
	}
	for _, index := range indices {
		index.Content = packageIndex
		conn.uploadPackageIndex(&index)
	}
}

func getIndexPaths(contents *map[string]s3.Key) []Index {
	indices := []Index{}
	pathRe := regexp.MustCompile(`^(.*)/Packages.gz$`)

	for key := range *contents {
		result := pathRe.FindStringSubmatch(key)
		if result != nil {
			index := Index{Path: result[1]}
			indices = append(indices, index)
		}
	}
	return indices
}

func getBucketPackages(contents *map[string]s3.Key) []string {
	packages := []string{}

	for key := range *contents {
		if filepath.Ext(key) == ".deb" {
			packages = append(packages, key)
		}
	}
	return packages
}

func createPackageIndex(packages []packageMetadata) []byte {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	t, err := template.New("Package Template").Parse(packagesTemplate)
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(writer, packages)

	writer.Close()
	return buf.Bytes()
}
