package main

import (
	"os"

	"github.com/nextrevision/sapt"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("sapt", "S3 apt repo utility that manages all packages remotely without a local mirror")

	bootstrap       = app.Command("bootstrap", "Bootstraps a new bucket")
	bootstrapPublic = bootstrap.Flag("public", "Make uploaded packages public").Bool()
	bootstrapBucket = bootstrap.Arg("bucket", "Name of bucket to use").Required().String()
	bootstrapRegion = bootstrap.Arg("region", "Region to use (defaults to AWS_REGION then us-east-1").String()

	upload       = app.Command("upload", "Uploads deb packages to S3")
	uploadPublic = upload.Flag("public", "Make uploaded packages public").Bool()
	uploadRm     = upload.Flag("rm", "Remove local packages after upload").Bool()
	uploadRoot   = upload.Arg("package_root", "Root path to packages/directory structure for upload").Required().File()
	uploadBucket = upload.Arg("bucket", "Name of bucket to use").Required().String()
	uploadRegion = upload.Arg("region", "Region to use (defaults to AWS_REGION then us-east-1").String()

	rescan       = app.Command("rescan", "Rescan the bucket and generate new indices")
	rescanPublic = rescan.Flag("public", "Make uploaded packages public").Bool()
	rescanBucket = rescan.Arg("bucket", "Name of bucket to use").Required().String()
	rescanRegion = rescan.Arg("region", "Region to use (defaults to AWS_REGION then us-east-1").String()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	// Bootstrap new repo
	case bootstrap.FullCommand():
		s3Conn := sapt.NewS3Connection(*bootstrapBucket, *bootstrapRegion, *bootstrapPublic)
		s3Conn.CreateBucket()
		sapt.ScanBucketPackages(s3Conn)

	// Upload packages
	case upload.FullCommand():
		s3Conn := sapt.NewS3Connection(*uploadBucket, *uploadRegion, *uploadPublic)
		sapt.UploadPackages(*uploadRoot, *uploadRm, s3Conn)
		sapt.ScanBucketPackages(s3Conn)

	// Rescan s3 and upload new apt data
	case rescan.FullCommand():
		s3Conn := sapt.NewS3Connection(*rescanBucket, *rescanRegion, *rescanPublic)
		sapt.ScanBucketPackages(s3Conn)
	}
}
