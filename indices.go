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
{{if .Architecture}}Architecture: {{.Architecture}}{{printf "\n"}}{{end}}{{if .Depends}}Depends: {{.Depends}}{{printf "\n"}}{{end}}{{if .Recommends}}Recommends: {{.Recommends}}{{printf "\n"}}{{end}}{{if .Conflicts}}Conflicts: {{.Conflicts}}{{printf "\n"}}{{end}}{{if .License}}License: {{.License}}{{printf "\n"}}{{end}}{{if .Vendor}}Vendor: {{.Vendor}}{{printf "\n"}}{{end}}{{if .Maintainer}}Maintainer: {{.Maintainer}}{{printf "\n"}}{{end}}{{if .Section}}Section: {{.Section}}{{printf "\n"}}{{end}}{{if .Priority}}Priority: {{.Priority}}{{printf "\n"}}{{end}}{{if .Homepage}}Homepage: {{.Homepage}}{{printf "\n"}}{{end}}{{if .Description}}Description: {{.Description}}{{printf "\n"}}{{end}}
{{end}}`

type Index struct {
	Path    string
	Content []byte
}

func ScanBucketPackages(conn *S3) {
	packages := []PackageMetadata{}
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
				m := MetadataFromHeaders(headers)
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

func createPackageIndex(packages []PackageMetadata) []byte {
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
