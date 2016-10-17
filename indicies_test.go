package sapt

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"testing"
)

func TestCreatePackageIndex(t *testing.T) {
	want := `Package: foo
Version: 1.2.3
Filename: foo_1.2.3_amd64.deb
MD5sum: 78455b8f4d492dba445d663b78ad50e2
SHA1: 02efeeec948c37999f37caddc1f64fd348fce963
SHA256: 507184bd48c8e7def1faf0eb96b3849d12110551e8bbc31720f423e8a1ea258c
Size: 7972868
Installed-Size: 57807
Architecture: amd64
Depends: bar (>= 1.4.0), banana (>= 1.4.0)
Recommends: x, y, z
Conflicts: oof
License: MIT
Vendor: unknown
Maintainer: Bob <bob@foo.bar>
Section: misc
Priority: optional
Homepage: https://foo.io
Description: Foo utility for barring some bananas

`

	packages := []packageMetadata{
		packageMetadata{
			Package:            "foo",
			Version:            "1.2.3",
			Filename:           "foo_1.2.3_amd64.deb",
			MD5sum:             "78455b8f4d492dba445d663b78ad50e2",
			SHA1:               "02efeeec948c37999f37caddc1f64fd348fce963",
			SHA256:             "507184bd48c8e7def1faf0eb96b3849d12110551e8bbc31720f423e8a1ea258c",
			Size:               "7972868",
			InstalledSize:      "57807",
			Architecture:       "amd64",
			Depends:            "bar (>= 1.4.0), banana (>= 1.4.0)",
			Recommends:         "x, y, z",
			Conflicts:          "oof",
			Maintainer:         "Bob <bob@foo.bar>",
			OriginalMaintainer: "Alice <alice@foo.bar>",
			Section:            "misc",
			Priority:           "optional",
			Homepage:           "https://foo.io",
			Bugs:               "https://bugs.foo.bar",
			Origin:             "https://github.com/bar/foo",
			License:            "MIT",
			Vendor:             "unknown",
			Description:        "Foo utility for barring some bananas",
		},
	}

	gzipData := createPackageIndex(packages)

	buf := bytes.NewBuffer(gzipData)
	r, err := gzip.NewReader(buf)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	s, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	if string(s) != want {
		t.Fatalf("not equal\nwant: %s\ngot:  %s", want, string(s))
	}
}
