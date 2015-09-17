package sapt

import (
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"regexp"
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
{{if .Architecture}}Architecture: {{.Architecture}}{{printf "\n"}}{{end}}{{if .License}}License: {{.License}}{{printf "\n"}}{{end}}{{if .Vendor}}Vendor: {{.Vendor}}{{printf "\n"}}{{end}}{{if .Maintainer}}Maintainer: {{.Maintainer}}{{printf "\n"}}{{end}}{{if .Section}}Section: {{.Section}}{{printf "\n"}}{{end}}{{if .Priority}}Priority: {{.Priority}}{{printf "\n"}}{{end}}{{if .Homepage}}Homepage: {{.Homepage}}{{printf "\n"}}{{end}}{{if .Description}}Description: {{.Description}}{{printf "\n"}}{{end}}
{{end}}`

type Index struct {
	Path    string
	Content []byte
}

func RescanBucket(s *S3) {
	packages := []PackageMetadata{}
	contents := s.GetBucketContents()
	packageList := getBucketPackages(contents)
	for _, pkg := range packageList {
		m := MetadataFromHeaders(s.GetObjectHeaders(pkg))
		packages = append(packages, *m)
	}
	packageIndex := createPackageIndex(packages)

	indicies := getIndexPaths(contents)
	if len(indicies) == 0 {
		indicies = append(indicies, Index{Path: "repo"})
	}
	for _, index := range indicies {
		index.Content = packageIndex
		s.UploadPackageIndex(&index)
	}
}

func getIndexPaths(contents *map[string]s3.Key) []Index {
	indicies := []Index{}
	pathRe := regexp.MustCompile(`^(.*)/Packages.gz$`)

	// TODO: err handle
	for key := range *contents {
		result := pathRe.FindStringSubmatch(key)
		if result != nil {
			index := Index{Path: result[1]}
			indicies = append(indicies, index)
		}
	}
	return indicies
}

func getBucketPackages(contents *map[string]s3.Key) []string {
	packages := []string{}

	// TODO: err handle
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

	// TODO: err handles
	t, _ := template.New("Package Template").Parse(packagesTemplate)
	t.Execute(os.Stdout, packages)
	t.Execute(writer, packages)

	writer.Close()
	return buf.Bytes()
}
