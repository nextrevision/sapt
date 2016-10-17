package sapt

import (
	"reflect"
	"strings"
	"testing"
)

func TestMetadataFromDeb(t *testing.T) {
	want := &packageMetadata{
		Filename:      "./fixtures/less_458-3_amd64.deb",
		Package:       "less",
		Version:       "458-3",
		Size:          "124466",
		Architecture:  "amd64",
		Maintainer:    "Anibal Monsalve Salazar <anibal@debian.org>",
		InstalledSize: "253",
		Depends:       "libc6 (>= 2.14), libtinfo5, debianutils (>= 1.8)",
		Section:       "text",
		Priority:      "important",
		Homepage:      "http://www.greenwoodsoftware.com/less/",
		Description: `pager program similar to more
 This package provides "less", a file pager (that is, a memory-efficient
 utility for displaying text one screenful at a time). Less has many
 more features than the basic pager "more". As part of the GNU project,
 it is widely regarded as the standard pager on UNIX-derived systems.
 .
 Also provided are "lessecho", a simple utility for ensuring arguments
 with spaces are correctly quoted; "lesskey", a tool for modifying the
 standard (vi-like) keybindings; and "lesspipe", a filter for specific
 types of input, such as .doc or .txt.gz files.
`,
	}

	metadata := metadataFromDeb("./fixtures/less_458-3_amd64.deb")
	if !reflect.DeepEqual(want, metadata) {
		t.Fatalf("not equal\nwant: %+v\ngot: %+v", want, metadata)
	}
}

func TestMetadataFromControl(t *testing.T) {
	want := &packageMetadata{
		Package:       "less",
		Version:       "458-3",
		Architecture:  "amd64",
		Maintainer:    "Anibal Monsalve Salazar <anibal@debian.org>",
		InstalledSize: "253",
		Depends:       "libc6 (>= 2.14), libtinfo5, debianutils (>= 1.8)",
		Section:       "text",
		Priority:      "important",
		Homepage:      "http://www.greenwoodsoftware.com/less/",
		Description: `pager program similar to more
 This package provides "less", a file pager (that is, a memory-efficient
 utility for displaying text one screenful at a time). Less has many
 more features than the basic pager "more". As part of the GNU project,
 it is widely regarded as the standard pager on UNIX-derived systems.
 .
 Also provided are "lessecho", a simple utility for ensuring arguments
 with spaces are correctly quoted; "lesskey", a tool for modifying the
 standard (vi-like) keybindings; and "lesspipe", a filter for specific
 types of input, such as .doc or .txt.gz files.
`,
	}

	controlOutput := `Package: less
Version: 458-3
Architecture: amd64
Maintainer: Anibal Monsalve Salazar <anibal@debian.org>
Installed-Size: 253
Depends: libc6 (>= 2.14), libtinfo5, debianutils (>= 1.8)
Section: text
Priority: important
Multi-Arch: foreign
Homepage: http://www.greenwoodsoftware.com/less/
Description: pager program similar to more
 This package provides "less", a file pager (that is, a memory-efficient
 utility for displaying text one screenful at a time). Less has many
 more features than the basic pager "more". As part of the GNU project,
 it is widely regarded as the standard pager on UNIX-derived systems.
 .
 Also provided are "lessecho", a simple utility for ensuring arguments
 with spaces are correctly quoted; "lesskey", a tool for modifying the
 standard (vi-like) keybindings; and "lesspipe", a filter for specific
 types of input, such as .doc or .txt.gz files.
`
	metadata := metadataFromControl(strings.Split(controlOutput, "\n"))
	if !reflect.DeepEqual(want, metadata) {
		t.Fatalf("not equal\nwant: %+v\ngot: %+v", want, metadata)
	}
}

func TestMetadataFromHeaders(t *testing.T) {
	want := &packageMetadata{
		Package:       "less",
		Version:       "458-3",
		Architecture:  "amd64",
		Maintainer:    "Anibal Monsalve Salazar <anibal@debian.org>",
		InstalledSize: "253",
		Depends:       "libc6 (>= 2.14), libtinfo5, debianutils (>= 1.8)",
		Section:       "text",
		Priority:      "important",
		Homepage:      "http://www.greenwoodsoftware.com/less/",
	}

	headers := map[string][]string{
		"X-Amz-Meta-Package":        []string{"less"},
		"X-Amz-Meta-Version":        []string{"458-3"},
		"X-Amz-Meta-Architecture":   []string{"amd64"},
		"X-Amz-Meta-Maintainer":     []string{"Anibal Monsalve Salazar <anibal@debian.org>"},
		"X-Amz-Meta-Installed-Size": []string{"253"},
		"X-Amz-Meta-Depends":        []string{"libc6 (>= 2.14), libtinfo5, debianutils (>= 1.8)"},
		"X-Amz-Meta-Section":        []string{"text"},
		"X-Amz-Meta-Priority":       []string{"important"},
		"X-Amz-Meta-Homepage":       []string{"http://www.greenwoodsoftware.com/less/"},
	}

	metadata := metadataFromHeaders(headers)
	if !reflect.DeepEqual(want, metadata) {
		t.Fatalf("not equal\nwant: %+v\ngot: %+v", want, metadata)
	}
}
