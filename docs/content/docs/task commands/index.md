---
title: "Task commands"
description: "List, add, complete, undo, update, and delete tasks with gtasks. Full command reference with examples."
draft: false
weight: 4
sitemap:
  priority: 0.8
---

## Help

```
❯ gtasks --help

Available Commands:
  add         Add a task
  clear       Hide all completed tasks
  done        Mark a task as done
  info        View detailed information about a task
  ls          List tasks
  rm          Delete a task
  undo        Mark a completed task as incomplete
  update      Update an existing task
```

Use `-l` / `--tasklist` on any task command to select a list. If omitted, gtasks uses `GTASKS_DEFAULT_TASKLIST` or `default_task_list` from the config file. If only one tasklist exists, it is selected automatically. Otherwise you will be prompted to choose one.

```
gtasks ls -l "DSC VIT"
gtasks add "Buy milk" -l "DSC VIT"
```

Task numbers are the 1-based index shown by `gtasks ls`. They can change when tasks are added, deleted, or sorted.

## Add a task

```
❯ gtasks add "Buy milk"
Creating task in DSC VIT
Task created
```

Interactive mode (no title given):

```
❯ gtasks add
Creating task in DSC VIT
Title: testing
Note: testing
Due Date: 12 July 2021
```

Flags still work:

```
gtasks add -l "DSC VIT" --title <some title> [--note <some note> | --due <some due date>]
gtasks add "Call dentist" -d tomorrow
```

### Recurring Tasks

Create multiple tasks with a repeating schedule using the `--repeat` flag:

```
❯ gtasks add -l "DSC VIT" "Daily standup" -d "2025-02-10" --repeat daily --repeat-count 5
Creating task in DSC VIT
Creating 5 recurring tasks...
Created 5 tasks
```

This creates 5 tasks for Feb 10, 11, 12, 13, 14.

Available repeat patterns:
- `daily` or `day`
- `weekly` or `week`
- `monthly` or `month`
- `yearly` or `year`

You can use `--repeat-count` to specify the number of occurrences:

```
gtasks add "Weekly sync" -d "2025-02-10" --repeat weekly --repeat-count 4
```

Or use `--repeat-until` to specify an end date:

```
gtasks add "Weekly sync" -d "2025-02-10" --repeat weekly --repeat-until "2025-03-10"
```

Both can be combined — the command stops at whichever limit is reached first.

## List tasks

```
❯ gtasks ls
Tasks in DSC VIT:
| NO |        TITLE         |          DESCRIPTION           | STATUS |     DUE      |
|----|----------------------|--------------------------------|--------|--------------|
|  1 | testing              | testing                        | pending| 12 July 2021 |
|  2 | HopeHouse            | Checkout the app. Maybe        | pending| 06 July 2021 |
```

A specific tasklist:

```
gtasks ls -l "DSC VIT"
gtasks ls --tasklist "DSC VIT"
```

### Output formats (table, json, csv)

Use `--format` to change the output format. The default is `table`.

```
❯ gtasks ls --format table

❯ gtasks ls --format json

❯ gtasks ls --format csv
```

JSON example (pipe to `jq`):

```
❯ gtasks ls -l "DSC VIT" --format json | jq '.[] | {title, status, due}'
```

CSV example (redirect to a file):

```
❯ gtasks ls -l "DSC VIT" --format csv > tasks.csv
```

### Include completed tasks

```
❯ gtasks ls --include-completed

❯ gtasks ls -l "DSC VIT" -i
```

### Show only completed tasks

```
❯ gtasks ls --completed

❯ gtasks ls -l "DSC VIT" --completed
```

### Sort and limit

```
❯ gtasks ls --sort due

❯ gtasks ls -l "DSC VIT" --sort title

❯ gtasks ls --max 5

❯ gtasks ls -l "DSC VIT" --max 10
```

Sort options: `due`, `title`, `position` (default).

## Mark a task as done

```
❯ gtasks done 1
Marked as complete: testing
```

With a specific tasklist:

```
❯ gtasks ls -l "DSC VIT"
❯ gtasks done -l "DSC VIT" 1
Marked as complete: testing
```

If no task number is given, you will be prompted to select a task.

## Undo a completed task

```
❯ gtasks ls -l "DSC VIT" --include-completed
❯ gtasks undo -l "DSC VIT" 1
Marked as incomplete: testing
```

If no task number is given, you will be prompted to select from completed tasks.

## Clear completed tasks

Hide all completed tasks from the list. This marks completed tasks as hidden so they won't be returned by the API (primarily affects tasks completed via the CLI).

```
❯ gtasks clear -l "DSC VIT"
✔ Clear all completed tasks from 'DSC VIT'? [y/N]: y
Cleared completed tasks from DSC VIT
```

Use `--force` or `-f` to skip the confirmation prompt:

```
❯ gtasks clear -l "DSC VIT" --force
Cleared completed tasks from DSC VIT
```

## View detailed task information

The `info` command displays detailed information about a task, including links/URLs that may have been shared to Google Tasks (e.g., from Android's "Share With..." feature).

By default, `info` only considers pending tasks (matching `ls`). Use `-i` to include completed tasks.

```
❯ gtasks info -l "DSC VIT" 1

Task: testing
Status: Needs action
Due: 12 July 2021
Notes: testing

Links:
  - https://example.com/some-link

View in Google Tasks: https://tasks.google.com/...
```

To get info on a completed task, use `-i` (must match how you listed the tasks):

```
❯ gtasks ls -l "DSC VIT" -i
❯ gtasks info -l "DSC VIT" 3 -i
```

The info command is particularly useful for viewing:

- Full task notes (not truncated)
- Links/URLs attached to the task
- WebViewLink to open the task in Google Tasks web interface
- Complete due date information
- Task completion status

## Update a task

Update an existing task's title, note, or due date.

### Interactive mode

When no flags are provided, you'll be prompted for each field with the current value displayed. Press Enter to keep the current value, or type a new value.

```
❯ gtasks update 1
Updating task: testing

Title [testing]: new title
Note [testing notes]:
Due [12 July 2021]:

Updated: new title
```

### Flag mode

Use flags to update specific fields without prompts:

```
❯ gtasks update 1 --title "New title"
Updating task: testing

Updated: New title

❯ gtasks update 1 --note "Updated note" --due "tomorrow"
Updating task: New title

Updated: New title
```

Available flags:
- `-t, --title` - New title for the task
- `-n, --note` - New note for the task
- `-d, --due` - New due date for the task

## Delete a task

```
❯ gtasks rm 1
Deleted: testing
```

With a specific tasklist:

```
❯ gtasks ls -l "DSC VIT"
❯ gtasks rm -l "DSC VIT" 1
Deleted: testing
```

If no task number is given, you will be prompted to select a task.

## Legacy commands

The older `gtasks tasks …` namespace still works as a deprecated compatibility path:

| Legacy | Prefer |
|--------|--------|
| `gtasks tasks view` | `gtasks ls` |
| `gtasks tasks add` | `gtasks add` |
| `gtasks tasks done` | `gtasks done` |
| `gtasks tasks undo` | `gtasks undo` |
| `gtasks tasks update` | `gtasks update` |
| `gtasks tasks info` | `gtasks info` |
| `gtasks tasks rm` | `gtasks rm` |
| `gtasks tasks clear` | `gtasks clear` |
