sapt
====

sapt is a s3 apt repo utility that manages all packages remotely without a local mirror.

Traditionally, one would need to store all packages locally on disk and run
`dpkg-scanpackages` or similar tool to generate indicies for apt. If s3 storage
is desired, a sync operation would then need to be performed, uploading new
packages as well as the repo indicies. `sapt` removes the need for keeping a
local mirror on disk by storing package metadata along with the package in s3.
Repo indicies are then generated and uploaded to s3 by quering each package for
their metadata attributes.

## Requirements

* dpkg-deb

## Installation

Directly with go:

    go get github.com/nextrevision/sapt

## Usage

    usage: sapt [<flags>] <command> [<args> ...]

    S3 apt repo utility that manages all packages remotely without a local mirror

    Flags:
      --help  Show help (also see --help-long and --help-man).

    Commands:
      help [<command>...]
        Show help.

      bootstrap [<flags>] <bucket> [<region>]
        Bootstraps a new bucket

      upload [<flags>] <package_root> <bucket> [<region>]
        Uploads deb packages to S3

      rescan [<flags>] <bucket> [<region>]
        Rescan the bucket and generate new indicies

Create a new private bucket:

    sapt bootstrap my-s3-apt-repo

Upload packages to the bucket:

    sapt upload ./packages my-s3-apt-repo

Rescan the bucket and generate new indicies:

    sapt rescan my-s3-apt-repo

### Public and Private Sources

Buckets and their contents can either be public or private. Apt client
configurations vary depending on this.

Public repo sample `sources.list` entry:

    # with a flat repo structure (sapt default)
    deb http://s3.amazonaws.com/my-s3-public-repo repo/
    # with a dist repo structure (ex. trusty)
    deb http://s3.amazonaws.com/my-s3-public-repo trusty main

Private repo sample `sources.list` entry (requires apt-transport-s3):

    # with a flat repo structure (sapt default)
    deb s3://<access_key>:[<secret_key>]@s3.amazonaws.com/my-s3-public-repo repo/
    # with a dist repo structure (ex. trusty)
    deb s3://<access_key>:[<secret_key>]@s3.amazonaws.com/my-s3-public-repo trusty main

Note the secret key is surrounded with brackets intentionally.

### Using with [apt-transport-s3](https://github.com/kyleshank/apt-transport-s3)

`sapt` is designed to work with [apt-transport-s3(https://github.com/kyleshank/apt-transport-s3)] in mind. Simply use `sapt`
to upload your packages and then use apt-transport-s3 to access them on your client systems.

#### Creating a public bucket with the apt-transport-s3

This assumes a Ubuntu OS.

Download dependencies:

    apt-get update && apt-get install -y build-essential dpkg-deb wget

Create a public repo with sapt:

     boossapttrap --public my-s3-public-repo

Download and create a package for apt-transport-s3:

    wget https://github.com/kyleshank/apt-transport-s3/archive/master.zip
    unzip apt-transport-s3-master.zip
    cd apt-transport-s3-master
    make
    dpkg-buildpackage -us -uc && dpkg-deb -b debian/apt-transport-s3

Upload package with sapt:

    sapt upload --public debian my-s3-public-repo

Now configure your apt clients with the following entry:

    deb http://s3.amazonaws.com/my-s3-public-repo repo/
