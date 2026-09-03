package cmd

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/BRO3886/gtasks/api"
	"github.com/BRO3886/gtasks/internal/utils"
	"github.com/araddon/dateparse"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"google.golang.org/api/tasks/v1"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List tasks",
	Long: `List tasks in a tasklist.

Uses the default tasklist when -l/--tasklist is not set.
You can control output with --format: table (default), json, csv.`,
	Example: `  gtasks ls
  gtasks ls -l work
  gtasks ls --tasklist work --sort due
  gtasks ls -i --format json`,
	Args: cobra.NoArgs,
	Run:  runLs,
}

var addCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Add a task",
	Long: `Add a task to a tasklist.

The title can be given as a positional argument or with --title.
If neither is provided, you will be prompted interactively.

Supports recurring tasks with --repeat:
  gtasks add "Standup" -d "2025-02-10" --repeat daily --repeat-count 5
  gtasks add "Weekly sync" -d "2025-02-10" --repeat weekly --repeat-until "2025-03-10"`,
	Example: `  gtasks add "Buy milk"
  gtasks add "Deploy ParsOps" -l work
  gtasks add -t "Call dentist" -d tomorrow
  gtasks add "Standup" -d "2025-02-10" --repeat daily --repeat-count 5`,
	Args: cobra.MaximumNArgs(1),
	Run:  runAdd,
}

var doneCmd = &cobra.Command{
	Use:   "done [task-number]",
	Short: "Mark a task as done",
	Long: `Mark a task as completed.

Task numbers match the 1-based index shown by gtasks ls.
If no task number is given, you will be prompted to select a task.`,
	Example: `  gtasks done 1
  gtasks done 3 -l work`,
	Args: cobra.MaximumNArgs(1),
	Run:  runDone,
}

var undoCmd = &cobra.Command{
	Use:   "undo [task-number]",
	Short: "Mark a completed task as incomplete",
	Long: `Mark a completed task as incomplete.

Task numbers match the 1-based index of completed tasks.
If no task number is given, you will be prompted to select a task.`,
	Example: `  gtasks undo 1
  gtasks undo 1 -l work`,
	Args: cobra.MaximumNArgs(1),
	Run:  runUndo,
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Hide all completed tasks",
	Long: `Hide all completed tasks from a tasklist.

This marks completed tasks as hidden so they won't be returned
by the API. Primarily affects tasks completed via the CLI.`,
	Example: `  gtasks clear
  gtasks clear -l work
  gtasks clear --force`,
	Args: cobra.NoArgs,
	Run:  runClear,
}

var rmCmd = &cobra.Command{
	Use:   "rm [task-number]",
	Short: "Delete a task",
	Long: `Delete a task from a tasklist.

Task numbers match the 1-based index shown by gtasks ls.
If no task number is given, you will be prompted to select a task.`,
	Example: `  gtasks rm 1
  gtasks rm 2 -l work`,
	Args: cobra.MaximumNArgs(1),
	Run:  runRm,
}

var infoCmd = &cobra.Command{
	Use:   "info [task-number]",
	Short: "View detailed information about a task",
	Long: `View detailed information about a task, including links, notes, and other metadata.

By default, only pending tasks are considered (matching ls). Use -i to include completed tasks.`,
	Example: `  gtasks info 1
  gtasks info 3 -l work
  gtasks info 3 -l work -i`,
	Args: cobra.MaximumNArgs(1),
	Run:  runInfo,
}

var updateCmd = &cobra.Command{
	Use:   "update [task-number]",
	Short: "Update an existing task",
	Long: `Update an existing task in a tasklist.

Interactive mode (no flags): prompts for each field showing current values.
Press Enter to keep the current value, or type a new value.

Flag mode: only update fields that are explicitly provided.`,
	Example: `  gtasks update 1
  gtasks update 1 --title "New title"
  gtasks update 1 --note "Updated note" --due tomorrow`,
	Args: cobra.MaximumNArgs(1),
	Run:  runUpdate,
}

func runLs(cmd *cobra.Command, args []string) {
	srv, err := api.GetService()
	if err != nil {
		utils.ErrorP("Failed to get service: %v\n", err)
		return
	}
	tList := getTaskLists(srv)

	taskItems, err := api.GetTasks(srv, tList.Id, viewTasksFlags.includeCompleted || viewTasksFlags.onlyCompleted, viewTasksFlags.max)
	if err != nil {
		color.Red(err.Error())
		return
	}

	utils.Sort(taskItems, viewTasksFlags.sort)

	var filteredTasks []*tasks.Task
	for _, task := range taskItems {
		if viewTasksFlags.onlyCompleted && task.Status == "needsAction" {
			continue
		}
		filteredTasks = append(filteredTasks, task)
	}

	switch viewTasksFlags.format {
	case "json":
		outputJSON(filteredTasks)
	case "csv":
		outputCSV(filteredTasks)
	default:
		outputTable(filteredTasks, tList.Title)
	}
}

func runAdd(cmd *cobra.Command, args []string) {
	srv, err := api.GetService()
	if err != nil {
		utils.ErrorP("Failed to get service: %v\n", err)
		return
	}
	tList := getTaskLists(srv)
	utils.Warn("Creating task in %s\n", tList.Title)

	title := resolveAddTitle(args, addTaskFlags.title)
	var notes string
	var dateInput string

	if title == "" && (addTaskFlags.note != "" || addTaskFlags.due != "") {
		utils.ErrorP("Error: task title is required\n")
		return
	} else if title != "" {
		notes = addTaskFlags.note
		dateInput = addTaskFlags.due
	} else {
		reader := bufio.NewReader(os.Stdin)

		utils.Print("Title: ")
		title = getInput(reader)
		utils.Print("Note: ")
		notes = getInput(reader)
		utils.Print("Due Date: ")
		dateInput = getInput(reader)
	}

	repeatPattern, err := parseRepeatUnit(addTaskFlags.repeat)
	if err != nil {
		utils.ErrorP("%v\n", err)
		return
	}

	if repeatPattern != repeatNone && dateInput == "" {
		utils.ErrorP("Error: due date (--due) is required when using --repeat\n")
		return
	}

	var startDate time.Time
	if dateInput != "" {
		t, err := dateparse.ParseAny(dateInput)
		if err != nil {
			utils.ErrorP("Error: date format incorrect. Valid examples: https://github.com/araddon/dateparse#extended-example\n")
			return
		}
		startDate = t
	}

	var untilDate *time.Time
	if addTaskFlags.repeatUntil != "" {
		t, err := dateparse.ParseAny(addTaskFlags.repeatUntil)
		if err != nil {
			utils.ErrorP("Error: repeat-until date format incorrect. Valid examples: https://github.com/araddon/dateparse#extended-example\n")
			return
		}
		untilDate = &t
	}

	var dates []time.Time
	if repeatPattern != repeatNone {
		dates = expandRepeatSchedule(startDate, repeatPattern, addTaskFlags.repeatCount, untilDate)
	} else if dateInput != "" {
		dates = []time.Time{startDate}
	} else {
		dates = []time.Time{}
	}

	if len(dates) == 0 {
		task := &tasks.Task{Title: title, Notes: notes}
		_, err = api.CreateTask(srv, task, tList.Id)
		if err != nil {
			utils.ErrorStyle.Printf("Unable to create task: %v", err)
			return
		}
		utils.Info("Task created\n")
	} else if len(dates) == 1 {
		task := &tasks.Task{Title: title, Notes: notes, Due: dates[0].Format(time.RFC3339)}
		_, err = api.CreateTask(srv, task, tList.Id)
		if err != nil {
			utils.ErrorStyle.Printf("Unable to create task: %v", err)
			return
		}
		utils.Info("Task created\n")
	} else {
		utils.Info("Creating %d recurring tasks...\n", len(dates))
		for i, d := range dates {
			task := &tasks.Task{Title: title, Notes: notes, Due: d.Format(time.RFC3339)}
			_, err = api.CreateTask(srv, task, tList.Id)
			if err != nil {
				utils.ErrorStyle.Printf("Unable to create task %d: %v\n", i+1, err)
				return
			}
		}
		utils.Info("Created %d tasks\n", len(dates))
	}
}

func runDone(cmd *cobra.Command, args []string) {
	srv, err := api.GetService()
	if err != nil {
		utils.ErrorP("Failed to get service: %v\n", err)
		return
	}
	tList := getTaskLists(srv)
	tID := tList.Id

	taskItems, err := api.GetTasks(srv, tID, false, 0)
	if err != nil {
		color.Red(err.Error())
		return
	}

	ind := getTaskIndex(args, taskItems, tList.Title)
	t := taskItems[ind]
	t.Status = "completed"

	_, err = api.UpdateTask(srv, t, tID)
	if err != nil {
		color.Red("Unable to mark task as completed: %v", err)
		return
	}
	utils.Info("Marked as complete: %s\n", t.Title)
}

func runUndo(cmd *cobra.Command, args []string) {
	srv, err := api.GetService()
	if err != nil {
		utils.ErrorP("Failed to get service: %v\n", err)
		return
	}
	tList := getTaskLists(srv)
	tID := tList.Id

	taskItems, err := api.GetTasks(srv, tID, true, 0)
	if err != nil {
		color.Red(err.Error())
		return
	}

	var completedTasks []*tasks.Task
	for _, task := range taskItems {
		if task.Status == "completed" {
			completedTasks = append(completedTasks, task)
		}
	}

	if len(completedTasks) == 0 {
		utils.Info("No completed tasks to undo\n")
		return
	}

	ind := getTaskIndex(args, completedTasks, tList.Title)
	t := completedTasks[ind]
	t.Status = "needsAction"
	t.Completed = nil

	_, err = api.UpdateTask(srv, t, tID)
	if err != nil {
		color.Red("Unable to mark task as incomplete: %v", err)
		return
	}
	utils.Info("Marked as incomplete: %s\n", t.Title)
}

func runClear(cmd *cobra.Command, args []string) {
	srv, err := api.GetService()
	if err != nil {
		utils.ErrorP("Failed to get service: %v\n", err)
		return
	}
	tList := getTaskLists(srv)

	if !clearTasksFlags.force {
		prompt := promptui.Prompt{
			Label:     fmt.Sprintf("Clear all completed tasks from '%s'", tList.Title),
			IsConfirm: true,
		}
		_, err := prompt.Run()
		if err != nil {
			utils.Info("Cancelled\n")
			return
		}
	}

	err = api.ClearTasks(srv, tList.Id)
	if err != nil {
		color.Red("Unable to clear completed tasks: %v", err)
		return
	}
	utils.Info("Cleared completed tasks from %s\n", tList.Title)
}

func runRm(cmd *cobra.Command, args []string) {
	srv, err := api.GetService()
	if err != nil {
		utils.ErrorP("Failed to get service: %v\n", err)
		return
	}
	tList := getTaskLists(srv)
	tID := tList.Id

	taskItems, err := api.GetTasks(srv, tID, false, 0)
	if err != nil {
		color.Red(err.Error())
		return
	}

	ind := getTaskIndex(args, taskItems, tList.Title)
	t := taskItems[ind]

	err = api.DeleteTask(srv, t.Id, tID)
	if err != nil {
		color.Red("Unable to delete task: %v", err)
		return
	}
	utils.Info("Deleted: %s\n", t.Title)
}

func runInfo(cmd *cobra.Command, args []string) {
	srv, err := api.GetService()
	if err != nil {
		utils.ErrorP("Failed to get service: %v\n", err)
		return
	}
	tList := getTaskLists(srv)
	tID := tList.Id

	taskItems, err := api.GetTasks(srv, tID, infoTaskFlags.includeCompleted, 0)
	if err != nil {
		color.Red(err.Error())
		return
	}

	ind := getTaskIndex(args, taskItems, tList.Title)
	t := taskItems[ind]

	utils.Print("\n")
	utils.Print("Task: %s\n", t.Title)

	status := "Needs action"
	if t.Status == "completed" {
		status = "Completed"
	}
	utils.Print("Status: %s\n", status)

	if t.Due != "" {
		due, err := time.Parse(time.RFC3339, t.Due)
		if err == nil {
			utils.Print("Due: %s\n", due.Local().Format("02 January 2006"))
		} else {
			utils.Print("Due: Not set\n")
		}
	} else {
		utils.Print("Due: Not set\n")
	}

	if t.Notes != "" {
		utils.Print("Notes: %s\n", t.Notes)
	} else {
		utils.Print("Notes: None\n")
	}

	utils.Print("\n")
	if len(t.Links) > 0 {
		utils.Print("Links:\n")
		for _, link := range t.Links {
			utils.Print("  - %s\n", link.Link)
		}
	} else {
		utils.Print("Links: No links\n")
	}

	if t.WebViewLink != "" {
		utils.Print("\nView in Google Tasks: %s\n", t.WebViewLink)
	}
	utils.Print("\n")
}

func runUpdate(cmd *cobra.Command, args []string) {
	srv, err := api.GetService()
	if err != nil {
		utils.ErrorP("Failed to get service: %v\n", err)
		return
	}
	tList := getTaskLists(srv)
	tID := tList.Id

	taskItems, err := api.GetTasks(srv, tID, false, 0)
	if err != nil {
		color.Red(err.Error())
		return
	}

	ind := getTaskIndex(args, taskItems, tList.Title)
	t := taskItems[ind]

	utils.Info("Updating task: %s\n\n", t.Title)

	titleFlagSet := cmd.Flags().Changed("title")
	noteFlagSet := cmd.Flags().Changed("note")
	dueFlagSet := cmd.Flags().Changed("due")
	flagMode := titleFlagSet || noteFlagSet || dueFlagSet

	var newTitle, newNote, newDue string

	if flagMode {
		if titleFlagSet {
			newTitle = updateTaskFlags.title
		} else {
			newTitle = t.Title
		}
		if noteFlagSet {
			newNote = updateTaskFlags.note
		} else {
			newNote = t.Notes
		}
		if dueFlagSet {
			newDue = updateTaskFlags.due
		}
	} else {
		reader := bufio.NewReader(os.Stdin)

		currentTitle := t.Title
		utils.Print("Title [%s]: ", currentTitle)
		newTitle = getInput(reader)
		if newTitle == "" {
			newTitle = currentTitle
		}

		currentNote := t.Notes
		if currentNote == "" {
			utils.Print("Note []: ")
		} else {
			utils.Print("Note [%s]: ", currentNote)
		}
		newNote = getInput(reader)
		if newNote == "" {
			newNote = currentNote
		}

		currentDue := formatDueHuman(t.Due)
		if currentDue == "-" {
			utils.Print("Due []: ")
		} else {
			utils.Print("Due [%s]: ", currentDue)
		}
		newDue = getInput(reader)
	}

	t.Title = newTitle
	t.Notes = newNote

	if newDue != "" {
		parsedDue, err := dateparse.ParseAny(newDue)
		if err != nil {
			utils.ErrorP("Error: date format incorrect. Valid examples: https://github.com/araddon/dateparse#extended-example\n")
			return
		}
		t.Due = parsedDue.Format(time.RFC3339)
	} else if !flagMode && newDue == "" {
		// Keep existing due date in interactive mode when user presses Enter
	} else if flagMode && dueFlagSet && newDue == "" {
		t.Due = ""
	}

	_, err = api.UpdateTask(srv, t, tID)
	if err != nil {
		color.Red("Unable to update task: %v", err)
		return
	}
	utils.Info("\nUpdated: %s\n", t.Title)
}

var (
	viewTasksFlags struct {
		includeCompleted bool
		onlyCompleted    bool
		sort             string
		format           string
		max              int
	}
	taskListFlag string
	addTaskFlags struct {
		title       string
		note        string
		due         string
		repeat      string
		repeatCount int
		repeatUntil string
	}
	clearTasksFlags struct {
		force bool
	}
	updateTaskFlags struct {
		title string
		note  string
		due   string
	}
	infoTaskFlags struct {
		includeCompleted bool
	}
)

func init() {
	registerViewFlags(lsCmd)
	registerAddFlags(addCmd)
	registerClearFlags(clearCmd)
	registerUpdateFlags(updateCmd)
	registerInfoFlags(infoCmd)

	taskCmds := []*cobra.Command{lsCmd, addCmd, doneCmd, undoCmd, updateCmd, infoCmd, rmCmd, clearCmd}
	for _, c := range taskCmds {
		attachTasklistFlag(c)
		rootCmd.AddCommand(c)
	}

	registerLegacyTasksCommand()
}

func attachTasklistFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&taskListFlag, "tasklist", "l", "", "specify a tasklist")
}

func registerViewFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&viewTasksFlags.includeCompleted, "include-completed", "i", false, "include completed tasks")
	cmd.Flags().BoolVar(&viewTasksFlags.onlyCompleted, "completed", false, "only show completed tasks")
	cmd.Flags().StringVar(&viewTasksFlags.sort, "sort", "position", "sort by [due,title,position]")
	cmd.Flags().StringVar(&viewTasksFlags.format, "format", "table", "output format: table, json, csv")
	cmd.Flags().IntVar(&viewTasksFlags.max, "max", 0, "maximum number of tasks to return (0 = all)")
}

func registerAddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&addTaskFlags.title, "title", "t", "", "task title")
	cmd.Flags().StringVarP(&addTaskFlags.note, "note", "n", "", "task note")
	cmd.Flags().StringVarP(&addTaskFlags.due, "due", "d", "", "due date (e.g., '2024-12-25', 'Dec 25', 'tomorrow')")
	cmd.Flags().StringVarP(&addTaskFlags.repeat, "repeat", "r", "", "repeat pattern: daily, weekly, monthly, yearly")
	cmd.Flags().IntVar(&addTaskFlags.repeatCount, "repeat-count", 0, "number of occurrences for repeating task")
	cmd.Flags().StringVar(&addTaskFlags.repeatUntil, "repeat-until", "", "end date for repeating task (e.g., '2025-03-01')")
}

func registerClearFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&clearTasksFlags.force, "force", "f", false, "skip confirmation prompt")
}

func registerUpdateFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&updateTaskFlags.title, "title", "t", "", "new title for the task")
	cmd.Flags().StringVarP(&updateTaskFlags.note, "note", "n", "", "new note for the task")
	cmd.Flags().StringVarP(&updateTaskFlags.due, "due", "d", "", "new due date for the task")
}

func registerInfoFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&infoTaskFlags.includeCompleted, "include-completed", "i", false, "include completed tasks when selecting by number")
}

type TaskOutput struct {
	Number      int    `json:"number"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status"`
	Due         string `json:"due,omitempty"`
}

func outputTable(taskItems []*tasks.Task, listTitle string) {
	utils.Print("Tasks in %s:\n", listTitle)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"No", "Title", "Description", "Status", "Due"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetCenterSeparator("|")
	table.SetRowLine(false)
	table.SetRowSeparator("-")
	table.SetAutoWrapText(false)

	for ind, task := range taskItems {
		row := []string{
			fmt.Sprintf("%d", ind+1),
			truncate(task.Title, 30),
			truncate(task.Notes, 40),
			statusLabel(task.Status),
			formatDueHuman(task.Due),
		}
		table.Append(row)
	}

	table.Render()
}

func truncate(s string, maxLen int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", "")

	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

func outputJSON(taskItems []*tasks.Task) {
	var output []TaskOutput

	for ind, task := range taskItems {
		output = append(output, TaskOutput{
			Number:      ind + 1,
			Title:       task.Title,
			Description: task.Notes,
			Status:      statusLabel(task.Status),
			Due:         formatDueISO(task.Due),
		})
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(output)
}

func outputCSV(taskItems []*tasks.Task) {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	_ = writer.Write([]string{"No", "Title", "Description", "Status", "Due"})

	for ind, task := range taskItems {
		_ = writer.Write([]string{
			fmt.Sprintf("%d", ind+1),
			task.Title,
			task.Notes,
			statusLabel(task.Status),
			formatDueHuman(task.Due),
		})
	}
}

func statusLabel(status string) string {
	switch status {
	case "completed":
		return "completed"
	case "needsAction":
		return "pending"
	default:
		return status
	}
}

func formatDueHuman(due string) string {
	if due == "" {
		return "-"
	}
	parsed, err := time.Parse(time.RFC3339, due)
	if err != nil {
		return "-"
	}
	return parsed.Local().Format("02 January 2006")
}

func formatDueISO(due string) string {
	if due == "" {
		return ""
	}
	parsed, err := time.Parse(time.RFC3339, due)
	if err != nil {
		return ""
	}
	return parsed.Local().Format("2006-01-02")
}

func getInput(reader *bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	if runtime.GOOS == "windows" {
		line = strings.Replace(line, "\r\n", "", -1)
	} else {
		line = strings.Replace(line, "\n", "", -1)
	}
	return line
}

func getTaskIndex(args []string, taskItems []*tasks.Task, title string) int {
	if len(args) == 1 {
		index, err := parseTaskNumber(args[0], len(taskItems))
		if err != nil {
			utils.ErrorP("Error: %s\n", err.Error())
		}
		return index
	}

	utils.Print("Tasks in %s:\n", title)

	tString := []string{}
	for _, i := range taskItems {
		tString = append(tString, i.Title)
	}

	prompt := promptui.Select{
		Label: "Select Task",
		Items: tString,
	}

	option, _, err := prompt.Run()
	if err != nil {
		utils.ErrorP("Error: %s\n", err.Error())
	}

	return option
}

func getTaskLists(srv *tasks.Service) tasks.TaskList {
	list, err := api.GetTaskLists(srv)
	if err != nil {
		utils.ErrorP("Error: %v\n", err)
	}

	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Title <= list[j].Title
	})

	specified := effectiveTasklist(taskListFlag)
	resolved, needsPrompt, err := resolveTasklist(list, specified)
	if err != nil {
		utils.ErrorP("Error: %s\n", err.Error())
	}
	if !needsPrompt {
		return resolved
	}

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

	return list[option]
}
