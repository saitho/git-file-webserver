package git

import (
	"fmt"
	"regexp"

	"github.com/saitho/static-git-file-server/utils"
)

var re = regexp.MustCompile(`^(v?\d+)\.\d+\.\d+$`)

func ResolveVirtualTag(client ClientInterface, virtualTag string) (GitTag, error) {
	for _, tag := range client.GetTags(client.GetCurrentRepo()) {
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
