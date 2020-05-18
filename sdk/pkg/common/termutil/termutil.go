package termutil

import (
	"sort"
	"strings"

	"github.com/getcouragenow/core-bs/sdk/pkg/common/colorutil"
)

const defaultIndent = 20

type Contents map[string][]string

func (c Contents) MergeContents(c2 Contents) Contents {
	for k, v := range c2 {
		c[k] = v
	}
	return c
}

func (c Contents) String(title string) string {
	var s strings.Builder
	keys, indent := c.getKeys()
	if title != "" {
		s.WriteString(colorutil.ColorMagenta(strings.Repeat("―", defaultIndent*4)))
		s.WriteRune('\n')
		s.WriteString(colorutil.ColorYellow(strings.ToUpper(title)))
		s.WriteRune('\n')
		s.WriteString(colorutil.ColorMagenta(strings.Repeat("―", defaultIndent*4)))
		s.WriteRune('\n')
	}
	for _, k := range keys {
		vals := c[k]
		s.WriteString(colorutil.ColorGreen(k))
		s.WriteRune(':')
		newIndent := indent - len(k) - 1
		if len(vals) < 1 {
			s.WriteRune('\n')
		}
		for _, v := range vals {
			s.WriteString(strings.Repeat(" ", newIndent))
			s.WriteString(v)
			s.WriteRune('\n')
			newIndent = indent
		}
	}
	s.WriteString(colorutil.ColorMagenta(strings.Repeat("―", defaultIndent*4)))
	s.WriteRune('\n')
	return s.String()
}

func (c Contents) getKeys() ([]string, int) {
	var keys []string
	indent := defaultIndent
	for k, v := range c {
		if len(v) > 0 {
			if len(k)+2 > indent {
				indent = len(k) + 2
			}
			keys = append(keys, k)
		}
	}
	c.sortKeys(keys)
	return keys, indent
}

func (c Contents) sortKeys(keys []string) []string {
	sort.Slice(keys, func(i, j int) bool {
		return strings.ToLower(keys[i]) < strings.ToLower(keys[j])
	})
	return keys
}
