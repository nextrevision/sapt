package sapt

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
)

type S3 struct {
	Auth       aws.Auth
	Region     aws.Region
	ACL        s3.ACL
	Bucket     *s3.Bucket
	Connection *s3.S3
}

func ConnectS3(bucket string, region string, public bool) *S3 {
	var acl s3.ACL
	var auth aws.Auth
	var err error

	// set auth
	auth, err = aws.EnvAuth()
	if err != nil {
		auth, err = aws.SharedAuth()
		if err != nil {
			log.Fatal(err)
		}
	}

	// set region
	if region == "" {
		region = os.Getenv("AWS_REGION")
		if region == "" {
			region = "us-east-1"
		}
	}
	awsRegion := aws.Regions[region]

	acl = s3.ACL("private")
	if public {
		acl = s3.ACL("public-read")
	}

	// establish connection
	conn := s3.New(auth, awsRegion)

	// set bucket
	bkt := conn.Bucket(bucket)

	return &S3{
		Auth:       auth,
		Region:     awsRegion,
		Bucket:     bkt,
		Connection: conn,
		ACL:        acl,
	}
}

func (s *S3) CreateBucket() {
	if err := s.Bucket.PutBucket(s.ACL); err != nil {
		log.Fatal(err)
	}
}

func (s *S3) uploadPackage(pkg *Package) {
	fileType := http.DetectContentType(pkg.Content)

	opts := s3.Options{
		Meta: MetadataToMap(*pkg.Metadata),
	}

	if err := s.Bucket.Put(pkg.Path, pkg.Content, fileType, s.ACL, opts); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Uploaded %s\n", pkg.Path)
}

func (s *S3) uploadPackageIndex(index *Index) {
	path := fmt.Sprintf("%s/Packages.gz", index.Path)
	fileType := http.DetectContentType(index.Content)
	opts := s3.Options{}
	if err := s.Bucket.Put(path, index.Content, fileType, s.ACL, opts); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Uploaded Package Index %s/Packages.gz\n", index.Path)
}

func (s *S3) getBucketContents() *map[string]s3.Key {
	contents, err := s.Bucket.GetBucketContents()
	if err != nil {
		log.Fatal(err)
	}
	return contents
}

func (s *S3) getObjectHeaders(object string) http.Header {
	headers := map[string][]string{}
	response, err := s.Bucket.Head(object, headers)
	if err != nil {
		log.Fatal(err)
	}
	return response.Header
}
