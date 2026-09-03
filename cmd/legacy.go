package cmd

import (
	"github.com/spf13/cobra"
)

// tasksCmd is a hidden compatibility namespace for the pre-refactor commands.
// Implementations are shared with the top-level commands; these entries only
// preserve the old gtasks tasks <subcommand> invocation path.
var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Task commands (deprecated: use ls, add, done, rm, ...)",
	Long: `Deprecated compatibility namespace for task commands.

Prefer the top-level commands:
  gtasks ls          (was: gtasks tasks view)
  gtasks add         (was: gtasks tasks add)
  gtasks done        (was: gtasks tasks done)
  gtasks undo        (was: gtasks tasks undo)
  gtasks update      (was: gtasks tasks update)
  gtasks info        (was: gtasks tasks info)
  gtasks rm          (was: gtasks tasks rm)
  gtasks clear       (was: gtasks tasks clear)`,
	Hidden: true,
}

func registerLegacyTasksCommand() {
	legacyView := &cobra.Command{
		Use:        "view",
		Short:      "List tasks",
		Deprecated: "use \"gtasks ls\" instead",
		Args:       cobra.NoArgs,
		Run:        runLs,
	}
	legacyAdd := &cobra.Command{
		Use:        "add",
		Short:      "Add a task",
		Deprecated: "use \"gtasks add\" instead",
		Args:       cobra.MaximumNArgs(1),
		Run:        runAdd,
	}
	legacyDone := &cobra.Command{
		Use:        "done",
		Short:      "Mark a task as done",
		Deprecated: "use \"gtasks done\" instead",
		Args:       cobra.MaximumNArgs(1),
		Run:        runDone,
	}
	legacyUndo := &cobra.Command{
		Use:        "undo",
		Short:      "Mark a completed task as incomplete",
		Deprecated: "use \"gtasks undo\" instead",
		Args:       cobra.MaximumNArgs(1),
		Run:        runUndo,
	}
	legacyUpdate := &cobra.Command{
		Use:        "update",
		Short:      "Update an existing task",
		Deprecated: "use \"gtasks update\" instead",
		Args:       cobra.MaximumNArgs(1),
		Run:        runUpdate,
	}
	legacyInfo := &cobra.Command{
		Use:        "info",
		Short:      "View detailed information about a task",
		Deprecated: "use \"gtasks info\" instead",
		Args:       cobra.MaximumNArgs(1),
		Run:        runInfo,
	}
	legacyRm := &cobra.Command{
		Use:        "rm",
		Short:      "Delete a task",
		Deprecated: "use \"gtasks rm\" instead",
		Args:       cobra.MaximumNArgs(1),
		Run:        runRm,
	}
	legacyClear := &cobra.Command{
		Use:        "clear",
		Short:      "Hide all completed tasks",
		Deprecated: "use \"gtasks clear\" instead",
		Args:       cobra.NoArgs,
		Run:        runClear,
	}

	registerViewFlags(legacyView)
	registerAddFlags(legacyAdd)
	registerClearFlags(legacyClear)
	registerUpdateFlags(legacyUpdate)
	registerInfoFlags(legacyInfo)

	tasksCmd.PersistentFlags().StringVarP(&taskListFlag, "tasklist", "l", "", "specify a tasklist")
	tasksCmd.AddCommand(legacyView, legacyAdd, legacyDone, legacyUndo, legacyUpdate, legacyInfo, legacyRm, legacyClear)
	rootCmd.AddCommand(tasksCmd)
}
