package git

import (
	"regexp"
	"strings"
)

func (c *Client) GetBranches() []string {
	output, _ := c.runGitCommand("branch", "-l", "-r", "--no-color")
	var branches []string
	for _, v := range strings.Split(output, "\n") {
		v = strings.TrimSpace(v)
		v = strings.TrimPrefix(v, "origin/")
		branches = append(branches, v)
	}
	return c.filterBranches(branches)
}

func (c *Client) filterBranches(references []string) []string {
	filters := c.Cfg.Display.Branches.Filter
	if len(filters) == 0 {
		return references
	}
	var output []string
	for _, v := range references {
		valid := false
		for _, f := range filters {
			if strings.HasPrefix(f, "/") && strings.HasSuffix(f, "/") {
				// RegEx
				expression := strings.Trim(f, "/")
				valid = regexp.MustCompile(expression).MatchString(v)
			} else {
				valid = f == v
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
