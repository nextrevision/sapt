sapt
====

[![Build Status](https://travis-ci.org/nextrevision/sapt.svg?branch=master)](https://travis-ci.org/nextrevision/sapt)

`sapt` is a s3 apt repo management utility.

Traditionally, one would need to store all packages locally on disk and run
`dpkg-scanpackages` or similar tool to generate indices for apt. If s3 storage
is desired, a sync operation would then need to be performed, uploading new
packages as well as the repo indices. `sapt` removes the need for keeping a
local mirror on disk by storing package metadata along with the package in s3.
Repo indices are then generated and uploaded to s3 by quering each package for
their metadata attributes.

It is similar to [deb-s3](https://github.com/krobertson/deb-s3), however sapt stores the package index data as s3 metadata.

## Requirements

* dpkg-deb

## Installation

### Binary

* OSX: ([64-bit](https://github.com/nextrevision/sapt/releases/download/0.1.0/sapt_darwin_amd64.zip) | [32-bit](https://github.com/nextrevision/sapt/releases/download/0.1.0/sapt_darwin_386.zip))
* Linux: ([64-bit](https://github.com/nextrevision/sapt/releases/download/0.1.0/sapt_linux_amd64.zip) | [32-bit](https://github.com/nextrevision/sapt/releases/download/0.1.0/sapt_linux_386.zip) | [Arm](https://github.com/nextrevision/sapt/releases/download/0.1.0/sapt_linux_arm.zip))

### Go

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
        Rescan the bucket and generate new indices

Create a new private bucket:

    sapt bootstrap my-s3-apt-repo

Upload packages to the bucket:

    sapt upload ./packages/ my-s3-apt-repo

Rescan the bucket and generate new indices:

    sapt rescan my-s3-apt-repo

### Authentication

`sapt` first looks for AWS credentials passed via environment variables
(`AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`). If neither are specified,
`~/.aws/credentials` is consulted. If neither method produces valid credentials,
an error will be thrown.

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

`sapt` is designed to work with [apt-transport-s3](https://github.com/kyleshank/apt-transport-s3) in mind. Simply use `sapt`
to upload your packages and then use apt-transport-s3 to access them on your client systems.

#### Creating a public bucket with the apt-transport-s3

This assumes a Ubuntu OS.

Download dependencies:

    apt-get update && apt-get install -y build-essential dpkg-dev python-apt-dev libcurl4-openssl-dev debhelper wget unzip

Create a public repo with sapt:

     sapt bootstrap --public my-s3-public-repo

Download and create a package for apt-transport-s3 (or just download from here https://launchpad.net/ubuntu/+source/apt-transport-s3):

    wget https://github.com/kyleshank/apt-transport-s3/archive/master.zip
    unzip master.zip
    cd apt-transport-s3-master
    make
    dpkg-buildpackage -us -uc && dpkg-deb -b debian/apt-transport-s3

Upload package with sapt:

    sapt upload --public debian my-s3-public-repo

Now configure your apt clients with the following entry:

    deb http://s3.amazonaws.com/my-s3-public-repo repo/
