package mkutil

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

/*
The package will scan makefiles and its target recipes.
*/

const (
	ruleIndent                      = 20 // indent for makefile recipes to os.Stdout
	builtinTargetPhony              = ".PHONY"
	builtinTargetSuffixes           = ".SUFFIXES"
	builtinTargetDefault            = ".DEFAULT"
	builtinTargetIntermediate       = ".INTERMEDIATE"
	builtinTargetSecondary          = ".SECONDARY"
	builtinTargetIgnore             = ".IGNORE"
	builtinTargetSilent             = ".SILENT"
	builtinTargetExportAllVariables = ".EXPORT_ALL_VARIABLES"
)

var (
	ruleRx           = regexp.MustCompile(`^([^\s%]+)\s*:`)
	isBuiltinTargets = map[string]bool{
		builtinTargetDefault:            true,
		builtinTargetPhony:              true,
		builtinTargetSuffixes:           true,
		builtinTargetIntermediate:       true,
		builtinTargetSecondary:          true,
		builtinTargetIgnore:             true,
		builtinTargetSilent:             true,
		builtinTargetExportAllVariables: true,
	}
)

type MakeRules map[string][]string

func ScanMakefiles(filepath string) (MakeRules, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to open file: %w", err))
	}
	m := MakeRules{}
	var buf []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()

		if strings.HasPrefix(line, "## ") {
			buf = append(buf, line[3:])
			continue
		}

		if matches := ruleRx.FindStringSubmatch(line); len(matches) > 1 {
			target := matches[1]
			if isBuiltinTargets[target] {
				continue
			}
			m[target] = buf
		}
		if len(buf) > 0 {
			buf = []string{}
		}
	}
	if err = sc.Err(); err != nil {
		return nil, errors.New(fmt.Sprintf("scan failed: %w", err))
	}
	return m, nil
}

func (m MakeRules) mergeRules(rules MakeRules) MakeRules {
	for k, v := range rules {
		m[k] = v
	}
	return m
}

func (m MakeRules) string() string {
	var s strings.Builder
	tasks, indent := m.getRules()
	for _, task := range tasks {
		taskHelps := m[task]
		s.WriteString(task)
		s.WriteString(":")
		newIndent := indent - len(task) - 1
		if len(taskHelps) < 1 {
			s.WriteRune('\n')
		}
		for _, th := range taskHelps {
			s.WriteString(strings.Repeat(" ", newIndent))
			s.WriteString(th)
			s.WriteRune('\n')
			newIndent = indent
		}
	}
	return s.String()
}

func (m MakeRules) getRules() ([]string, int) {
	var tasks []string // makefile tasks / recipes
	var indent = 0
	for k, v := range m {
		if len(v) > 0 {
			if len(k)+2 > ruleIndent {
				indent = len(k) + 2
			}
			tasks = append(tasks, k)
		}
	}
	m.sortTasks(tasks)
	return tasks, indent
}

func (m MakeRules) sortTasks(tasks []string) []string {
	sort.Slice(tasks, func(i, j int) bool {
		return strings.ToLower(tasks[i]) < strings.ToLower(tasks[j])
	})
	return tasks
}
