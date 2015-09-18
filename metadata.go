package sapt

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
)

type PackageMetadata struct {
	Package            string `json:"X-Amz-Meta-Package"`
	Priority           string `json:"X-Amz-Meta-Priority"`
	Section            string `json:"X-Amz-Meta-Section"`
	InstalledSize      string `json:"X-Amz-Meta-InstalledSize"`
	Maintainer         string `json:"X-Amz-Meta-Maintainer"`
	OriginalMaintainer string `json:"X-Amz-Meta-OriginalMaintainer"`
	Architecture       string `json:"X-Amz-Meta-Architecture"`
	Version            string `json:"X-Amz-Meta-Version"`
	Depends            string `json:"X-Amz-Meta-Depends"`
	Filename           string `json:"X-Amz-Meta-Filename"`
	Size               string `json:"X-Amz-Meta-Size"`
	MD5sum             string `json:"X-Amz-Meta-Md5sum"`
	SHA1               string `json:"X-Amz-Meta-Sha1"`
	SHA256             string `json:"X-Amz-Meta-Sha256"`
	Description        string `json:"X-Amz-Meta-Description"`
	DescriptionMd5     string `json:"X-Amz-Meta-DescriptionMd5"`
	Homepage           string `json:"X-Amz-Meta-Homepage"`
	Bugs               string `json:"X-Amz-Meta-Bugs"`
	Origin             string `json:"X-Amz-Meta-Origin"`
	License            string `json:"X-Amz-Meta-License"`
	Vendor             string `json:"X-Amz-Meta-Vendor"`
}

func MetadataFromFile(path string) *PackageMetadata {
	metadata := PackageMetadata{}

	stat, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	dpkgOut, err := exec.Command("dpkg-deb", "-f", path).Output()
	if err != nil {
		log.Fatal(err)
	}
	for _, line := range strings.Split(string(dpkgOut), "\n") {
		if line == "" {
			continue
		}
		data := strings.SplitAfterN(line, ":", 2)
		if len(data) != 2 {
			continue
		}
		key := strings.TrimSuffix(data[0], ":")
		value := strings.TrimSpace(data[1])

		switch key {
		case "Package":
			metadata.Package = value
		case "Priority":
			metadata.Priority = value
		case "Section":
			metadata.Section = value
		case "Installed-Size":
			metadata.InstalledSize = value
		case "Maintainer":
			metadata.Maintainer = value
		case "Original-Maintainer":
			metadata.OriginalMaintainer = value
		case "Architecture":
			metadata.Architecture = value
		case "Version":
			metadata.Version = value
		case "Depends":
			metadata.Depends = value
		case "Filename":
			metadata.Filename = value
		case "Description":
			metadata.Description = value
		case "Description-md5":
			metadata.DescriptionMd5 = value
		case "Homepage":
			metadata.Homepage = value
		case "Bugs":
			metadata.Bugs = value
		case "Origin":
			metadata.Origin = value
		case "License":
			metadata.License = value
		case "Vendor":
			metadata.Vendor = value
		}
	}
	metadata.Filename = path
	metadata.Size = strconv.FormatInt(stat.Size(), 10)
	return &metadata
}

func MetadataFromHeaders(headers http.Header) *PackageMetadata {
	var metadata PackageMetadata
	mapping := map[string]string{}
	// silly song and dance to convert
	// from map[string][]string to map[string]string
	for key, value := range headers {
		mapping[key] = value[0]
	}
	jsonMapping, err := json.Marshal(mapping)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(jsonMapping, &metadata)
	return &metadata
}

func metadataToMap(pm PackageMetadata) map[string][]string {
	mapping := map[string][]string{}
	v := reflect.ValueOf(pm)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		key := t.Field(i).Name
		value := v.Field(i).String()
		if key != "mapping" {
			mapping[key] = []string{value}
		}
	}

	return mapping
}
