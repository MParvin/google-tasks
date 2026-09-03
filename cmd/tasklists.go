package cmd

import (
	"github.com/BRO3886/gtasks/api"
	"github.com/BRO3886/gtasks/internal/utils"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"google.golang.org/api/tasks/v1"
)

var tasklistsCmd = &cobra.Command{
	Use:   "tasklists",
	Short: "Manage tasklists",
	Long: `List and manage tasklists for the currently signed-in account.

Running this command with no subcommand lists all tasklists.

  gtasks tasklists
  gtasks tasklists add "Work"
  gtasks tasklists add --title "Work"
  gtasks tasklists update Work -t "New name"
  gtasks tasklists rm Work`,
	Example: `  gtasks tasklists
  gtasks tasklists add "Work"
  gtasks tasklists update Work -t "Personal"
  gtasks tasklists rm Work`,
	Run: runTasklistsView,
}

var showlistsCmd = &cobra.Command{
	Use:        "view",
	Short:      "List tasklists",
	Long:       `List tasklists for the account currently signed in.`,
	Deprecated: "use \"gtasks tasklists\" instead",
	Run:        runTasklistsView,
}

var addListcmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Add a tasklist",
	Long: `Add a tasklist for the currently signed-in account.

The title can be given as a positional argument or with --title / -t.`,
	Example: `  gtasks tasklists add "Work"
  gtasks tasklists add --title "Work"
  gtasks tasklists add -t "Work"`,
	Args: cobra.MaximumNArgs(1),
	Run:  runTasklistsAdd,
}

var removeListCmd = &cobra.Command{
	Use:   "rm [tasklist]",
	Short: "Delete a tasklist",
	Long: `Delete a tasklist for the currently signed-in account.

If no tasklist name is given, you will be prompted to select one.`,
	Example: `  gtasks tasklists rm
  gtasks tasklists rm Work`,
	Args: cobra.MaximumNArgs(1),
	Run:  runTasklistsRm,
}

var updateTitleCmd = &cobra.Command{
	Use:   "update [tasklist]",
	Short: "Update a tasklist",
	Long: `Update a tasklist title for the currently signed-in account.

The tasklist to update can be given as a positional argument.
The new title is set with --title / -t. If no tasklist is given,
you will be prompted to select one.`,
	Example: `  gtasks tasklists update Work -t "Personal"
  gtasks tasklists update -t "New name"`,
	Args: cobra.MaximumNArgs(1),
	Run:  runTasklistsUpdate,
}

func runTasklistsView(cmd *cobra.Command, args []string) {
	srv, err := api.GetService()
	if err != nil {
		utils.ErrorP("Failed to get service: %v\n", err)
		return
	}
	list, err := api.GetTaskLists(srv)
	if err != nil {
		utils.ErrorP("Error: %v\n", err)
	}

	for index, i := range list {
		utils.Print("[%d] %s\n", index+1, i.Title)
	}
}

func runTasklistsAdd(cmd *cobra.Command, args []string) {
	listTitle := resolveAddTitle(args, addListFlags.title)
	if listTitle == "" {
		utils.ErrorP("Error: title is required\nUsage: gtasks tasklists add <title>\n")
		return
	}

	srv, err := api.GetService()
	if err != nil {
		utils.ErrorP("Failed to get service: %v\n", err)
		return
	}
	t := &tasks.TaskList{Title: listTitle}
	r, err := srv.Tasklists.Insert(t).Do()
	if err != nil {
		utils.ErrorP("Unable to create task list. %v", err)
	}
	utils.Info("task list created: %s", r.Title)
}

func runTasklistsRm(cmd *cobra.Command, args []string) {
	srv, err := api.GetService()
	if err != nil {
		utils.ErrorP("Failed to get service: %v\n", err)
		return
	}
	list, err := api.GetTaskLists(srv)
	if err != nil {
		utils.ErrorP("Error: %v\n", err)
	}

	var selected tasks.TaskList
	if len(args) == 1 {
		resolved, _, err := resolveTasklist(list, args[0])
		if err != nil {
			utils.ErrorP("Error: %s\n", err.Error())
		}
		selected = resolved
	} else {
		utils.Print("Choose a Tasklist: ")
		var l []string
		for _, i := range list {
			l = append(l, i.Title)
		}

		prompt := promptui.Select{
			Label: "Select Tasklist",
			Items: l,
		}
		option, _, err := prompt.Run()
		if err != nil {
			color.Red("Error: " + err.Error())
			return
		}
		selected = list[option]
	}

	utils.Print("%s: %s\n", utils.WarnStyle.Sprint("Deleting list..."), selected.Title)

	err = api.DeleteTaskList(srv, selected.Id)
	if err != nil {
		utils.ErrorP("Error deleting tasklist: %s", err.Error())
		return
	}
	utils.Info("Tasklist deleted")
}

func runTasklistsUpdate(cmd *cobra.Command, args []string) {
	newTitle := updateListFlags.title
	if newTitle == "" {
		utils.ErrorP("Error: new title is required\nUsage: gtasks tasklists update [tasklist] --title <new-title>\n")
		return
	}

	srv, err := api.GetService()
	if err != nil {
		utils.ErrorP("Failed to get service: %v\n", err)
		return
	}

	list, err := api.GetTaskLists(srv)
	if err != nil {
		utils.ErrorP("Error: %v\n", err)
	}

	var selected tasks.TaskList
	if len(args) == 1 {
		resolved, _, err := resolveTasklist(list, args[0])
		if err != nil {
			utils.ErrorP("Error: %s\n", err.Error())
		}
		selected = resolved
	} else {
		utils.Print("Choose a Tasklist:")
		var l []string
		for _, i := range list {
			l = append(l, i.Title)
		}

		prompt := promptui.Select{
			Label: "Select Tasklist",
			Items: l,
		}
		option, _, err := prompt.Run()
		if err != nil {
			utils.ErrorP("Error: %s", err.Error())
		}
		selected = list[option]
	}

	selected.Title = newTitle
	_, err = api.UpdateTaskList(srv, &selected)
	if err != nil {
		utils.ErrorP("Error updating tasklist: %s", err.Error())
	}
	utils.Info("Tasklist title updated")
}

var (
	addListFlags struct {
		title string
	}
	updateListFlags struct {
		title string
	}
)

func init() {
	addListcmd.Flags().StringVarP(&addListFlags.title, "title", "t", "", "title of the tasklist")
	updateTitleCmd.Flags().StringVarP(&updateListFlags.title, "title", "t", "", "new title for the tasklist")
	tasklistsCmd.AddCommand(showlistsCmd, addListcmd, removeListCmd, updateTitleCmd)
	rootCmd.AddCommand(tasklistsCmd)
}
