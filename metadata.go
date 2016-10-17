package sapt

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
)

// packageMetadata represents the control information of a package
type packageMetadata struct {
	Package            string `json:"X-Amz-Meta-Package"`
	Priority           string `json:"X-Amz-Meta-Priority"`
	Section            string `json:"X-Amz-Meta-Section"`
	InstalledSize      string `json:"X-Amz-Meta-InstalledSize"`
	Maintainer         string `json:"X-Amz-Meta-Maintainer"`
	OriginalMaintainer string `json:"X-Amz-Meta-OriginalMaintainer"`
	Architecture       string `json:"X-Amz-Meta-Architecture"`
	Version            string `json:"X-Amz-Meta-Version"`
	Depends            string `json:"X-Amz-Meta-Depends"`
	Recommends         string `json:"X-Amz-Meta-Recommends"`
	Conflicts          string `json:"X-Amz-Meta-Conflicts"`
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

// metadataFromDeb extracts the control information from a deb package
// returning a packageMetadata representation
func metadataFromDeb(path string) *packageMetadata {
	stat, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	dpkgOut, err := exec.Command("dpkg-deb", "-f", path).Output()
	if err != nil {
		log.Fatal(err)
	}

	controlOutput := strings.Split(string(dpkgOut), "\n")
	metadata := metadataFromControl(controlOutput)

	metadata.Filename = path
	metadata.Size = strconv.FormatInt(stat.Size(), 10)

	return metadata
}

// metadataFromHeaders extracts the metadata from S3 object headers
func metadataFromHeaders(headers http.Header) *packageMetadata {
	var metadata packageMetadata

	for key, value := range headers {
		switch key {
		case "X-Amz-Meta-Package":
			metadata.Package = value[0]
		case "X-Amz-Meta-Priority":
			metadata.Priority = value[0]
		case "X-Amz-Meta-Section":
			metadata.Section = value[0]
		case "X-Amz-Meta-Installed-Size":
			metadata.InstalledSize = value[0]
		case "X-Amz-Meta-Maintainer":
			metadata.Maintainer = value[0]
		case "X-Amz-Meta-Original-Maintainer":
			metadata.OriginalMaintainer = value[0]
		case "X-Amz-Meta-Architecture":
			metadata.Architecture = value[0]
		case "X-Amz-Meta-Version":
			metadata.Version = value[0]
		case "X-Amz-Meta-Depends":
			metadata.Depends = value[0]
		case "X-Amz-Meta-Recommends":
			metadata.Recommends = value[0]
		case "X-Amz-Meta-Conflicts":
			metadata.Conflicts = value[0]
		case "X-Amz-Meta-Filename":
			metadata.Filename = value[0]
		case "X-Amz-Meta-Description":
			metadata.Description = value[0]
		case "X-Amz-Meta-Description-md5":
			metadata.DescriptionMd5 = value[0]
		case "X-Amz-Meta-Homepage":
			metadata.Homepage = value[0]
		case "X-Amz-Meta-Bugs":
			metadata.Bugs = value[0]
		case "X-Amz-Meta-Origin":
			metadata.Origin = value[0]
		case "X-Amz-Meta-License":
			metadata.License = value[0]
		case "X-Amz-Meta-Vendor":
			metadata.Vendor = value[0]
		}
	}
	return &metadata
}

// metadataFromControl converts a control file to a packageMetadata object
func metadataFromControl(control []string) *packageMetadata {
	metadata := packageMetadata{}
	inDescription := false

	for _, line := range control {
		if line == "" {
			continue
		}

		if inDescription && strings.HasPrefix(line, " ") {
			metadata.Description += fmt.Sprintf("%s\n", line)
			continue
		} else if inDescription {
			inDescription = false
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
		case "Recommends":
			metadata.Recommends = value
		case "Conflicts":
			metadata.Conflicts = value
		case "Filename":
			metadata.Filename = value
		case "Description":
			inDescription = true
			metadata.Description = fmt.Sprintf("%s\n", value)
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

	return &metadata
}

// metadataToMap converts packageMetadata to a string map
func metadataToMap(pm packageMetadata) map[string][]string {
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
