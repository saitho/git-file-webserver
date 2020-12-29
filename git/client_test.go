package git_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/saitho/static-git-file-server/config"
	"github.com/saitho/static-git-file-server/git"
)

func TestGetRepositoryBySlug(t *testing.T) {
	Convey("test getting existing repository by slug", t, func() {
		subject := new(git.Client)
		subject.Cfg = new(config.Config)
		repo := &config.RepoConfig{
			Title: "saithos Profile",
			Slug:  "saitho",
			Url:   "http://github.com/saitho/saitho.git",
		}
		subject.Cfg.Git.Repositories = []*config.RepoConfig{repo}

		So(subject.GetRepositoryBySlug("saitho"), ShouldEqual, repo)
	})

	Convey("test getting unknown repository by slug", t, func() {
		subject := new(git.Client)
		subject.Cfg = new(config.Config)

		So(subject.GetRepositoryBySlug("unknown"), ShouldBeNil)
	})
}

func TestSelectRepository(t *testing.T) {
	Convey("test selecting an existing repository", t, func() {
		subject := new(git.Client)
		subject.Cfg = new(config.Config)
		repo := &config.RepoConfig{
			Title: "saithos Profile",
			Slug:  "saitho",
			Url:   "http://github.com/saitho/saitho.git",
		}
		subject.Cfg.Git.Repositories = []*config.RepoConfig{repo}

		So(subject.CurrentRepo, ShouldBeNil)
		err := subject.SelectRepository("saitho")
		So(subject.CurrentRepo, ShouldEqual, repo)
		So(err, ShouldBeNil)
	})

	Convey("test selecting an unknown repository", t, func() {
		subject := new(git.Client)
		subject.Cfg = new(config.Config)

		So(subject.CurrentRepo, ShouldBeNil)
		err := subject.SelectRepository("unknown")
		So(err, ShouldBeError)
		So(subject.CurrentRepo, ShouldBeNil)
	})
}
