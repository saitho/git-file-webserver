package git_test

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"

	"github.com/saitho/static-git-file-server/config"
	"github.com/saitho/static-git-file-server/git"
)

type MockedClient struct {
	mock.Mock
	git.ClientInterface
}

func (m *MockedClient) GetCurrentRepo() *config.RepoConfig {
	return m.Called().Get(0).(*config.RepoConfig)
}
func (m *MockedClient) GetTags(*config.RepoConfig) []git.GitTag {
	return m.Called().Get(0).([]git.GitTag)
}

func TestResolveVirtualTag(t *testing.T) {
	Convey("test resolving an existing virtual tag", t, func() {
		client := new(MockedClient)
		repo := &config.RepoConfig{}
		client.On("GetCurrentRepo").Return(repo)
		tags := []git.GitTag{
			{Tag: "v2.1.0", Date: time.Unix(21000, 0)},
			{Tag: "v2.0.0", Date: time.Unix(20000, 0)},
			{Tag: "v1.1.1", Date: time.Unix(12000, 0)},
			{Tag: "v1.1.0", Date: time.Unix(11000, 0)},
			{Tag: "v1.0.0", Date: time.Unix(10000, 0)},
		}
		client.On("GetTags", mock.Anything).Return(tags)

		resolvedTag, err := git.ResolveVirtualTag(client, "v1")

		So(err, ShouldBeNil)
		So(resolvedTag.Tag, ShouldEqual, tags[2].Tag)
	})

	Convey("test resolving an unknown virtual tag", t, func() {
		client := new(MockedClient)
		tag := git.GitTag{
			Tag: "v1.0.0",
		}
		repo := &config.RepoConfig{}
		client.On("GetCurrentRepo").Return(repo)
		client.On("GetTags", mock.Anything).Return([]git.GitTag{tag})

		resolvedTag, err := git.ResolveVirtualTag(client, "v2")

		So(err, ShouldBeError)
		So(resolvedTag, ShouldNotEqual, tag)
	})
}

func tagIncluded(tags []git.GitTag, tag git.GitTag) bool {
	for _, gitTag := range tags {
		if gitTag == tag {
			return true
		}
	}
	return false
}

func TestInsertVirtualTags(t *testing.T) {
	Convey("test creating virtual tags", t, func() {
		tags := []git.GitTag{
			{Tag: "v2.1.0", Date: time.Unix(21000, 0)},
			{Tag: "v2.0.0", Date: time.Unix(20000, 0)},
			{Tag: "v1.1.1", Date: time.Unix(12000, 0)},
			{Tag: "v1.1.0", Date: time.Unix(11000, 0)},
			{Tag: "v1.0.0", Date: time.Unix(10000, 0)},
		}

		Convey("with versions in descending order", func() {
			newTags := git.InsertVirtualTags(tags)
			So(tagIncluded(newTags, git.GitTag{Tag: "v1", Date: time.Unix(12000, 0)}), ShouldBeTrue)
			So(tagIncluded(newTags, git.GitTag{Tag: "v2", Date: time.Unix(21000, 0)}), ShouldBeTrue)
		})
	})
}
