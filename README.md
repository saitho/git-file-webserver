# Static Git File Server

[![semantic-release](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg)](https://github.com/go-semantic-release/semantic-release)
[![CI](https://github.com/saitho/git-file-webserver/workflows/CI/badge.svg?branch=master)](https://github.com/saitho/git-file-webserver/actions?query=workflow%3ACI+branch%3Amaster)
[![Go Report Card](https://goreportcard.com/badge/github.com/saitho/git-file-webserver)](https://goreportcard.com/report/github.com/saitho/git-file-webserver)

## Features

This binary allows serving selective content of a Git repository.
The user can access those files per branch or tag.

This was originally built for publishing up-to-date JSONSchema files per version.

* Set up a **static file mirror** of your Git Repository
* Limit display to **specific files**
* Make files **accessible per branch or tag**

## Live Demo

Find a live demo here: https://schema.stackhead.io

This demo serves all JSON files of the "schema" folder of the [StackHead repository](https://github.com/getstackhead/stackhead).

Configuration:
```yaml
---
git:
  url: https://github.com/getstackhead/stackhead.git
  work_dir: schemas
  update:
    mode: webhook_github
    webhook:
      github:
        secret: 'REDACTED'
files:
  - "**/*.json"
display:
  branches:
    filter:
      - master
      - next
  tags:
    show_date: false
    virtual_tags:
      enable_semver_major: true
```

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

### Git repository settings

Inside the `git` section you have to set the path to your repository in `url` setting.

Additionally you may set a `work_dir`, which means that only the files in this directory are considered to be served.
**Note:** As of right now this is a global option. Keep that in mind if the folder name changes in releases or branches.

Updating the repository can be done two ways: time-based (cache-like) or webhook-based (mirror-like).

Setting the `update.mode` to "cache" will refresh the repository every hour (per default).
You may change the update time by setting `update.cache.time` (in seconds).

Setting the `update.mode` to "webhook_github" will refresh the repository on new commits or tags to the repository.
The repository needs to be setup manually for that (see below).

```yaml
---
git:
  repositories:
    - title: StackHead
      slug: stackhead
      url: https://github.com/getstackhead/stackhead.git
      work_dir: schemas
      update:
        mode: cache # either "cache" (default) or "webhook_github"
        cache:
           time: 3600 # default: 60 minutes
        webhook:
          github:
            secret: foobar # secret to be used with GitHub webhook
```

### File settings

The `files` section is an array of fileglobs you can use to specify which files should be displayed.

In the example below, only JSON files are served. So files that do not end with `.json` and folders not containing `.json` files are not displayed.

```yaml
---
files:
  - "**/*.json"
```

### Display settings

In the `display` section you can define how the frontend should look like.

```yaml
---
display:
  branches:
     filter:
        - master
        - develop
        - /feature/.*/ # regex: any feature branches
  tags:
    filter: []
    order: desc
    show_date: true
    virtual_tags:
      enable_semver_major: true
  index:
    show_branches: true
    show_tags: true
```

You may change the `order` of tags or hide the tag date in `tags` subsection.

You can also toggle the display of branches or tags for the index page in the `index` subsection.

#### Virtual tags

Virtual major tags will always point to the latest version inside a major release.
They can be enabled by setting `display.tags.virtual_tags.enable_semver_major` to true.

This only will consider semantic versions.

_Example:_ Given the tags v1.0.0, v1.1.0, v2.0.0, two virtual tags "v1" and "v2" will be displayed.
"v1" links to "v1.1.0" and "v2" links to "v2.0.0".
Pushing a release "v2.0.1" will make "v2" automatically refer to "v2.0.1".

This also works when you don't use the "v" prefix in your tags. The virtual tag will then also not have a "v" prefix.

#### Filters

Using the `filter` settings in _branches_ and _tags_ subsection you can filter which branches or tags should be displayed.
Wrap the string in / to evaluate the inside as Golang regular expression.

Note that the filters apply before the virtual tags.
Given the repostory has two tags v1.0.0 and v1.1.0.
For some reason, you exclude v1.1.0 via a filter.
The virtual tag will then point to v1.0.0, not v1.1.0!

### Logging

A log file is created at `./logs.txt`.
The setting `loc_level` can be used to set the logging level.

Allowed values: debug, info, warning, error, panic
Default: warning

## Mirroring with Webhooks

If you want to update the repository whenever something is pushed or tagged in your repository, you can use GitHub webhooks.

Set the `mode` in `git.update` subsection to "webhook_github" and define a `secret`.

```yaml
---
git:
  update:
    mode: webhook_github
    webhook:
      github:
        secret: your-secret-here
```

Then, create your webhook on GitHub as follows:

1. Go to your Repository Settings
2. Select the "Webhooks" option in the left navigation menu
3. Click the button "Add webhook" at the top right
4. Add set your server URL as "Payload URL" (ending with `/webhook/github`), e.g. `https://schema.stackhead.io/webhook/github`
5. Select "application/json" as "Content type"
6. Set your secret from configuration as "Secret"
7. Choose the option "Let me select individual events." and enable the folowing events:
   * Branch or tag creation
   * Branch or tag deletion
   * Pushes

You're ready to go. New changes to your repository should be mirrored automatically to your webserver.
