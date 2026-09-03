---
title: "Tasklist commands"
description: "Manage Google task lists from the terminal with gtasks: list, create, update, and delete tasklists, and see inline help for every tasklist command you can run."
draft: false
weight: 3
sitemap:
  priority: 0.8
---

## Help

```
❯ gtasks tasklists --help

List and manage tasklists for the currently signed-in account.

Running this command with no subcommand lists all tasklists.

Usage:
  gtasks tasklists [flags]
  gtasks tasklists [command]

Available Commands:
  add         Add a tasklist
  rm          Delete a tasklist
  update      Update a tasklist
```

## List tasklists

```
❯ gtasks tasklists
[1] DSC VIT
[2] Daily todo
[3] Life
```

`gtasks tasklists view` still works as a deprecated alias.

## Create a tasklist

```
❯ gtasks tasklists add "Work"

❯ gtasks tasklists add --title "some title"

❯ gtasks tasklists add -t "some title"
```

## Update a tasklist title

```
❯ gtasks tasklists update Work -t "Personal"

❯ gtasks tasklists update --title "some title"
```

If no tasklist name is given, you will be prompted to select one. `--title` / `-t` is required for the new name.

## Delete a tasklist

```
❯ gtasks tasklists rm Work
```

If no name is given, you will be prompted to select one:

```
❯ gtasks tasklists rm
Use the arrow keys to navigate: ↓ ↑ → ←
? Select Tasklist:
  ▸ VIT
    Daily todo
    personal projects
    To watch
↓   DSC VIT
```
