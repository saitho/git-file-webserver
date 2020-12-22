package git

import (
	"fmt"
	"github.com/saitho/static-git-file-server/utils"
	"regexp"
)

var re = regexp.MustCompile(`^(v?\d+)\.\d+\.\d+$`)

func ResolveVirtualTag(gitHandler *GitHandler, virtualTag string) (GitTag, error) {
	for _, tag := range gitHandler.GetTags() {
		majorTag := re.FindStringSubmatch(tag.Tag)
		if len(majorTag) < 2 {
			continue
		}
		if majorTag[1] == virtualTag {
			return tag, nil
		}
	}
	return GitTag{}, fmt.Errorf("cannot resolve virtual tag")
}

func InsertVirtualTags(tags []GitTag) []GitTag {
	var newTags []GitTag
	var processedMajorTags []string
	for _, tag := range tags {
		majorTag := re.FindStringSubmatch(tag.Tag)[1]
		if len(majorTag) > 0 && !utils.Contains(processedMajorTags, majorTag) {
			newTags = append(newTags, GitTag{
				Tag:  majorTag,
				Date: tag.Date,
			})
			processedMajorTags = append(processedMajorTags, majorTag)
		}
		newTags = append(newTags, tag)
	}
	return newTags
}
