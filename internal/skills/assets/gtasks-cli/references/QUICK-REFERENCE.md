# Quick Reference Card

Fast lookup for common gtasks commands. For detailed explanations, see the main SKILL.md file.

## System Checks (Run First!)

```bash
# Check if gtasks is installed (cross-platform)
gtasks --version 2>/dev/null || gtasks.exe --version 2>/dev/null

# macOS/Linux: which command
which gtasks

# Windows: where command
where gtasks

# Check environment variables (macOS/Linux)
echo $GTASKS_CLIENT_ID
echo $GTASKS_CLIENT_SECRET

# Check environment variables (Windows PowerShell)
echo $env:GTASKS_CLIENT_ID
echo $env:GTASKS_CLIENT_SECRET

# Check authentication status
gtasks tasklists &>/dev/null && echo "✓ Authenticated" || echo "✗ Not authenticated"
```

## Authentication

```bash
gtasks login                    # Authenticate with Google
gtasks logout                   # Remove credentials
```

## Task Lists

```bash
gtasks tasklists                         # List all task lists
gtasks tasklists add "List Name"         # Create task list
gtasks tasklists rm Work                 # Delete task list
gtasks tasklists update Work -t "New"    # Rename task list
```

## View Tasks

```bash
gtasks ls                                # List tasks (default list)
gtasks ls -l "Work"                      # List tasks in a specific list
gtasks ls -i                             # Include completed tasks
gtasks ls --completed                    # Show only completed tasks
gtasks ls --sort=due                     # Sort by due date
gtasks ls --sort=title                   # Sort by title
gtasks ls --format=json                  # JSON output
gtasks ls --format=csv                   # CSV output
```

## Create Tasks

```bash
gtasks add                                            # Interactive mode
gtasks add "Title"                                    # With title only
gtasks add "Title" -n "Notes"                         # With notes
gtasks add "Title" -d "2024-12-25"                    # With due date
gtasks add "Title" -n "Notes" -d "tomorrow"           # All fields
gtasks add "Title" -l "Work"                          # Specify list
```

## Complete Tasks

```bash
gtasks done                     # Interactive selection
gtasks done 1                   # Complete task #1
gtasks done 3 -l "Work"         # Complete task #3 in Work list
```

## Delete Tasks

```bash
gtasks rm                       # Interactive selection
gtasks rm 2                     # Delete task #2
gtasks rm 1 -l "Personal"       # Delete task #1 in Personal list
```

## Task Details

```bash
gtasks info                     # Interactive selection
gtasks info 1                   # Show details for task #1
gtasks info 2 -l "Work"         # Show details for task #2 in Work list
```

## Date Format Examples

All these work with the `-d` flag:

```
2024-12-25          # ISO format
Dec 25, 2024        # Month day, year
December 25         # Month day (current year)
12/25/2024          # US format
tomorrow            # Relative day
next Friday         # Relative named day
in 3 days           # Relative duration
```

## Common Workflows

### Add task with deadline
```bash
gtasks add "Submit proposal" -l "Work" -d "next Friday"
```

### Check today's tasks
```bash
gtasks ls -l "Work" --sort=due
```

### Complete multiple tasks
```bash
gtasks done -l "Work"    # Select first task
gtasks done -l "Work"    # Select next task
```

### Export tasks
```bash
gtasks ls --format=json > tasks.json
gtasks ls --format=csv > tasks.csv
```

### Create shopping list
```bash
gtasks tasklists add "Shopping"
gtasks add "Milk" -l "Shopping"
gtasks add "Bread" -l "Shopping"
gtasks add "Eggs" -l "Shopping"
```

## Flags Reference

### Global Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--tasklist` | `-l` | Specify task list by name |

### View Tasks Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--include-completed` | `-i` | Include completed tasks |
| `--completed` | | Show only completed tasks |
| `--sort` | | Sort by: due, title, position |
| `--format` | | Output format: table, json, csv |

### Add Task Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--title` | `-t` | Task title (required in flag mode) |
| `--note` | `-n` | Task notes/description |
| `--due` | `-d` | Due date (flexible format) |

### Add Task List Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--title` | `-t` | Task list title (required) |

## Output Formats

### Table (Default)
```
No  Title              Description         Status     Due
1   Finish report      Q4 analysis         pending    25 December 2024
2   Team meeting       Weekly sync         pending    -
```

### JSON
```json
[
  {
    "number": 1,
    "title": "Finish report",
    "description": "Q4 analysis",
    "status": "pending",
    "due": "2024-12-25"
  }
]
```

### CSV
```
No,Title,Description,Status,Due
1,Finish report,Q4 analysis,pending,25 December 2024
2,Team meeting,Weekly sync,pending,-
```

## Task Numbering

- Tasks are numbered starting from 1
- Task numbers are shown in view/list output
- Task numbers can change when tasks are added/deleted/sorted
- Always view the list first to get current numbers

## Status Values

- `pending` - Task needs action (not completed)
- `completed` - Task is marked as done

## Error Messages

| Error | Cause | Solution |
|-------|-------|----------|
| "command not found: gtasks" | GTasks not installed | Run `curl -fsSL https://gtasks.sidv.dev/install \| bash` or download from [releases](https://github.com/BRO3886/gtasks/releases) |
| "Failed to get service" | Not authenticated or missing env vars | Check env vars, then run `gtasks login` |
| Missing GTASKS_CLIENT_ID/SECRET | Environment variables not set | Export GTASKS_CLIENT_ID and GTASKS_CLIENT_SECRET |
| "tasklist ... not found" | List doesn't exist | Check with `gtasks tasklists` |
| "task ... not found" | Invalid task number | Run `gtasks ls` to see valid numbers |
| "Date format incorrect" | Unparseable date | Use format like "2024-12-25" or "tomorrow" |

## Tips

1. **Task numbers change** - Always view list before using numbers
2. **Use -l flag** - Avoid interactive prompts in scripts
3. **Flexible dates** - Natural language works: "tomorrow", "next week"
4. **JSON for parsing** - Use `--format=json` when processing with jq
5. **Include completed** - Use `-i` to see full task history

## Getting Help

```bash
gtasks --help                  # General help
gtasks ls --help               # List-tasks help
gtasks add --help              # Add-task help
gtasks tasklists --help        # Task lists command help
```
