# Static Git File Server

[![semantic-release](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg)](https://github.com/go-semantic-release/semantic-release)
[![CI](https://github.com/saitho/git-file-webserver/workflows/CI/badge.svg?branch=master)](https://github.com/saitho/git-file-webserver/actions?query=workflow%3ACI+branch%3Amaster)
[![Go Report Card](https://goreportcard.com/badge/github.com/saitho/git-file-webserver)](https://goreportcard.com/report/github.com/saitho/git-file-webserver)

## Features

This is binary allows serving selective content of a Git repository.
The user can access those files per branch or tag.

This was originally built for publishing up-to-date JSONSchema files per version.

* Set up a **static file mirror** of your Git Repository
* Limit display to **specific files**
* Make files **accessible per branch or tag**

## Usage

You may use the binaries from our [release section](https://github.com/saitho/git-file-webserver/releases), build one yourself or use our [pre-made Docker container](https://hub.docker.com/r/saitho/git-file-webserver).

### Docker

Webserver is served at port 8080.

```
docker run -p 8080:80 -v "`pwd`/config.yml.dist":/config/config.yml saitho/git-file-webserver:latest
```

### Binary

Requires Git to be installed! Webserver is served at port 8080.

```
./git-file-webserver -p 8080 -c config.yml.dist
```

## Configuration

The `config.yml` file is used to configure the repository that should be displayed.
It can also be used to limit the displayed files or set a work directory.

### git

Inside the `git` section you have to set the path to your repository in `url` setting.

Additionally you may set a `work_dir`, which means that only the files in this directory are considered to be served.
**Note:** As of right now this is a global option. Keep that in mind if the folder name changes in releases or branches.

The `cache_time` specifies when the repository content is invalidated and updated.
This will be eventually replaced or extended by Webhook functionality, so it is automatically updated when changes are pushed.

```yaml
---
git:
  url: https://github.com/getstackhead/stackhead.git
  work_dir: schemas
  cache_time: 3600 # default: 60 minutes
```

### files

The `files` section is an array of fileglobs you can use to specify which files should be displayed.

In the example below, only JSON files are served. So files that do not end with `.json` and folders not containing `.json` files are not displayed.

```yaml
---
files:
  - "**/*.json"
```

### display

In the `display` section you can define how the frontend should look like.

You may change the `order` of tags or hide the tag date in `tags` subsection.

You can also toggle the display of branches or tags for the index page in the `index` subsection.

```yaml
---
display:
  tags:
    order: desc
    show_date: true
  index:
    show_branches: true
    show_tags: true
```
