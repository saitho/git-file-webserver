package git

import (
	"regexp"
	"strings"
	"time"
)

type GitTag struct {
	Tag  string
	Date time.Time
}

func (c *Client) GetTags() []GitTag {
	output, _ := c.runGitCommand("for-each-ref", "--sort=-creatordate", "--format=%(refname)---%(creatordate)", "refs/tags")
	var tags []GitTag
	for _, v := range strings.Split(output, "\n") {
		if v == "" {
			continue
		}
		split := strings.Split(v, "---")
		intDate, _ := time.Parse("Mon Jan 2 15:04:05 2006 -0700", split[1])
		tags = append(tags, GitTag{
			Tag:  strings.TrimPrefix(split[0], "refs/tags/"),
			Date: intDate,
		})
	}

	return c.filterTags(tags)
}

func (c *Client) filterTags(references []GitTag) []GitTag {
	filters := c.Cfg.Display.Tags.Filter
	if len(filters) == 0 {
		return references
	}
	var output []GitTag
	for _, v := range references {
		valid := false
		for _, f := range filters {
			if strings.HasPrefix(f, "/") && strings.HasSuffix(f, "/") {
				// RegEx
				expression := strings.Trim(f, "/")
				valid = regexp.MustCompile(expression).MatchString(v.Tag)
			} else {
				valid = f == v.Tag
			}
			if valid {
				break
			}
		}
		if !valid {
			continue
		}
		output = append(output, v)
	}
	return output
}
