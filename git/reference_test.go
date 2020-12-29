package git_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/saitho/static-git-file-server/config"
	"github.com/saitho/static-git-file-server/git"
)

func TestGetShowRef(t *testing.T) {
	Convey("test getting show-ref of tag reference", t, func() {
		client := new(git.Client)
		client.CurrentRepo = &config.RepoConfig{
			Title: "saithos Profile",
			Slug:  "saitho",
			Url:   "http://github.com/saitho/saitho.git",
		}

		subject := &git.Reference{
			Client: client,
			Type:   "tag",
			Name:   "v1.2.3",
		}

		Convey("without workdir", func() {
			client.CurrentRepo.WorkDir = ""
			So(subject.GetShowRef("foo/bar.ext"), ShouldEqual, "refs/tags/v1.2.3:foo/bar.ext")
		})

		Convey("with workdir", func() {
			client.CurrentRepo.WorkDir = "public"
			So(subject.GetShowRef("foo/bar.ext"), ShouldEqual, "refs/tags/v1.2.3:public/foo/bar.ext")
		})

		Convey("with empty file path", func() {
			client.CurrentRepo.WorkDir = ""
			So(subject.GetShowRef(""), ShouldEqual, "refs/tags/v1.2.3:./")
		})
	})

	Convey("test getting show-ref of branch reference", t, func() {
		client := new(git.Client)
		client.CurrentRepo = &config.RepoConfig{
			Title: "saithos Profile",
			Slug:  "saitho",
			Url:   "http://github.com/saitho/saitho.git",
		}

		subject := &git.Reference{
			Client: client,
			Type:   "branch",
			Name:   "feature/some-epic-feature",
		}

		Convey("without workdir", func() {
			client.CurrentRepo.WorkDir = ""
			So(subject.GetShowRef("foo/bar.ext"), ShouldEqual, "origin/feature/some-epic-feature:foo/bar.ext")
		})

		Convey("with workdir", func() {
			client.CurrentRepo.WorkDir = "public"
			So(subject.GetShowRef("foo/bar.ext"), ShouldEqual, "origin/feature/some-epic-feature:public/foo/bar.ext")
		})

		Convey("with empty file path", func() {
			client.CurrentRepo.WorkDir = ""
			So(subject.GetShowRef(""), ShouldEqual, "origin/feature/some-epic-feature:./")
		})
	})
}
