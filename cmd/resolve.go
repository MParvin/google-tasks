package cmd

import (
	"fmt"
	"strconv"

	"github.com/BRO3886/gtasks/internal/config"
	"google.golang.org/api/tasks/v1"
)

// effectiveTasklist returns the tasklist name to use: -l/--tasklist wins, then config/env default.
func effectiveTasklist(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	return config.GetDefaultTaskList()
}

// resolveTasklist picks a tasklist from the available lists.
// If specified is non-empty, it must match a list title exactly.
// If specified is empty and exactly one list exists, that list is returned.
// If specified is empty and multiple lists exist, needsPrompt is true.
func resolveTasklist(list []tasks.TaskList, specified string) (tasks.TaskList, bool, error) {
	if specified != "" {
		for _, tl := range list {
			if tl.Title == specified {
				return tl, false, nil
			}
		}
		return tasks.TaskList{}, false, fmt.Errorf("tasklist %q not found", specified)
	}

	if len(list) == 0 {
		return tasks.TaskList{}, false, fmt.Errorf("no tasklists found")
	}

	if len(list) == 1 {
		return list[0], false, nil
	}

	return tasks.TaskList{}, true, nil
}

// parseTaskNumber converts a 1-based task number into a 0-based slice index.
func parseTaskNumber(arg string, count int) (int, error) {
	n, err := strconv.Atoi(arg)
	if err != nil || n < 1 || n > count {
		return -1, fmt.Errorf("task %s not found", arg)
	}
	return n - 1, nil
}

// resolveAddTitle returns the task title from a positional argument, falling back to --title.
func resolveAddTitle(args []string, flagTitle string) string {
	if len(args) > 0 && args[0] != "" {
		return args[0]
	}
	return flagTitle
}
