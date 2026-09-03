---
name: gtasks-cli
description: Manage Google Tasks from the command line - view, create, update, delete tasks and task lists. Use when the user asks to interact with Google Tasks, manage to-do items, create task lists, mark tasks complete, or check their Google Tasks.
homepage: https://github.com/BRO3886/gtasks
license: MIT
compatibility: Requires gtasks CLI tool to be installed and authenticated
metadata:
  author: BRO3886
  version: "1.0"
required-env:
  - name: GTASKS_CLIENT_ID
    description: Google OAuth2 client ID — can also be set in config file under [credentials]
  - name: GTASKS_CLIENT_SECRET
    description: Google OAuth2 client secret — can also be set in config file under [credentials]
allowed-tools: Bash(gtasks:*)
---

# Google Tasks CLI Skill

This skill enables you to manage Google Tasks directly from the command line using the `gtasks` CLI tool.

It can also be installed into supported AI agent skill directories with the built-in `gtasks skills` commands.

## Prerequisites

Before using any commands, ensure the following requirements are met:

### 1. GTasks Installation

Check if gtasks is installed on the system:

```bash
# Cross-platform check (works on macOS, Linux, Windows Git Bash)
gtasks --version 2>/dev/null || gtasks.exe --version 2>/dev/null || echo "gtasks not found"

# Or use which/where commands
# macOS/Linux:
which gtasks

# Windows (Command Prompt):
where gtasks

# Windows (PowerShell):
Get-Command gtasks
```

**If gtasks is not installed:**

**macOS/Linux (recommended):**
```bash
curl -fsSL https://gtasks.sidv.dev/install | bash
```
Installs to `~/.local/bin` by default. Set `INSTALL_DIR` to override:
```bash
INSTALL_DIR=/usr/local/bin curl -fsSL https://gtasks.sidv.dev/install | bash
```

**Manual install:**
1. Download the binary for your system from [GitHub Releases](https://github.com/BRO3886/gtasks/releases)
2. Move to a directory in your PATH (e.g. `~/.local/bin` or `/usr/local/bin`)
3. `chmod +x gtasks`

**Windows:** Download the binary from [GitHub Releases](https://github.com/BRO3886/gtasks/releases) and add to PATH.

Verify installation: `gtasks --version`

**IMPORTANT for Agents:** Always check if gtasks is installed before attempting to use it. If the command is not found, inform the user and provide installation instructions.

### 2. Environment Variables

Set up Google OAuth2 credentials as environment variables:

```bash
export GTASKS_CLIENT_ID="your-client-id.apps.googleusercontent.com"
export GTASKS_CLIENT_SECRET="your-client-secret"
```

**How to get credentials:**
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the Google Tasks API
4. Create OAuth2 credentials (Application type: "Desktop app")
5. Note the authorized redirect URIs that gtasks uses:
   - `http://localhost:8080/callback`
   - `http://localhost:8081/callback`
   - `http://localhost:8082/callback`
   - `http://localhost:9090/callback`
   - `http://localhost:9091/callback`

**For persistent setup**, the recommended approach is to add credentials to the gtasks config file:

```toml
# ~/.config/gtasks/config.toml  (new installs)
# ~/.gtasks/config.toml          (legacy installs)
[credentials]
client_id     = "your-client-id.apps.googleusercontent.com"
client_secret = "your-client-secret"
```

Set permissions: `chmod 600 ~/.config/gtasks/config.toml`

Alternatively, export from your shell profile — do not commit these values to version control:

```bash
export GTASKS_CLIENT_ID="your-client-id"
export GTASKS_CLIENT_SECRET="your-client-secret"
```

### 3. Authentication

Once environment variables are set, authenticate with Google:

```bash
gtasks login
```

This will open a browser for OAuth2 authentication. The token is stored in the system keyring (macOS Keychain, Linux Secret Service, Windows Credential Manager). On headless systems where no keyring is available, it falls back to a token file in the config directory. If you no longer need access, run `gtasks logout` to remove the stored token.

### 4. Optional: Install This Skill for Supported Agents

If `gtasks` is already installed, you can install this skill into supported agent directories with:

```bash
gtasks skills install
```

Supported agent targets:

- `claude` -> `~/.claude/skills/gtasks-cli/`
- `codex` -> `~/.agents/skills/gtasks-cli/`
- `openclaw` -> `~/.openclaw/skills/gtasks-cli/`

Useful commands:

```bash
gtasks skills status
gtasks skills install --agent codex
gtasks skills install --agent all
gtasks skills uninstall --agent codex
```

For automated setups, prefer `--agent <name>` or `--agent all` to avoid interactive prompts.

## Core Concepts

- **Task Lists**: Containers that hold tasks (like "Work", "Personal", "Shopping")
- **Tasks**: Individual to-do items within a task list
- **Task Properties**: Title (required), notes/description (optional), due date (optional), status (pending/completed)

## Command Structure

Common operations are top-level commands:

```
gtasks ls
gtasks add "Buy milk"
gtasks done 1
gtasks tasklists
```

Use `-l` / `--tasklist` to select a list. If omitted, gtasks uses the configured default, auto-selects when only one list exists, or prompts interactively.

## Authentication

### Login
```bash
gtasks login
```
Opens browser for Google OAuth2 authentication. Required before using any other commands.

### Logout
```bash
gtasks logout
```
Removes stored credentials from the system keyring (and token file if present).

## Skill Management

These commands manage installation of the `gtasks-cli` skill itself for supported AI agents.

### Check Skill Status

```bash
gtasks skills status
```

Shows whether the skill is installed for Claude, Codex, and OpenClaw, along with the installed version when available.

### Install the Skill

```bash
gtasks skills install
gtasks skills install --agent claude
gtasks skills install --agent codex
gtasks skills install --agent openclaw
gtasks skills install --agent all
```

Installs the embedded `gtasks-cli` skill files into the selected agent skill directory.

### Uninstall the Skill

```bash
gtasks skills uninstall
gtasks skills uninstall --agent codex
gtasks skills uninstall --agent all
```

Removes the installed `gtasks-cli` skill from the selected agent skill directory.

## Task List Management

### View All Task Lists
```bash
gtasks tasklists
```
Displays all task lists with numbered indices.

**Output Example:**
```
[1] My Tasks
[2] Work
[3] Personal
```

### Create a Task List
```bash
gtasks tasklists add "Work Projects"
gtasks tasklists add --title "Shopping List"
```
Creates a new task list with the specified title.

**Flags:**
- `-t, --title`: Task list title (optional if a positional title is given)

### Delete a Task List
```bash
gtasks tasklists rm Work
gtasks tasklists rm
```
Delete by name, or omit the name for an interactive prompt.

### Update Task List Title
```bash
gtasks tasklists update Work -t "New Title"
gtasks tasklists update -t "New Title"
```
Rename a task list. `--title` / `-t` is required for the new name. If no current name is given, you will be prompted to select one.

**Flags:**
- `-t, --title`: New title for the task list (required)

## Task Management

All task commands can optionally specify a task list using the `-l` flag. If omitted, gtasks uses `GTASKS_DEFAULT_TASKLIST` env var or `tasks.default_task_list` from the config file. If only one list exists, it is selected automatically. Otherwise you'll be prompted interactively.

### View Tasks

**Basic view:**
```bash
gtasks ls
gtasks ls -l "Work"
```

**Include completed tasks:**
```bash
gtasks ls --include-completed
gtasks ls -i
```

**Show only completed tasks:**
```bash
gtasks ls --completed
```

**Sort tasks:**
```bash
gtasks ls --sort=due        # Sort by due date
gtasks ls --sort=title      # Sort by title
gtasks ls --sort=position   # Sort by position (default)
```

**Output formats:**
```bash
gtasks ls --format=table    # Table format (default)
gtasks ls --format=json     # JSON output
gtasks ls --format=csv      # CSV output
```

**Table Output Example:**
```
Tasks in Work:
No  Title              Description         Status     Due
1   Finish report      Q4 analysis         pending    25 December 2024
2   Team meeting       Weekly sync         pending    -
3   Code review        PR #123            completed  20 December 2024
```

**JSON Output Example:**
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

### Create a Task

**Interactive mode:**
```bash
gtasks add
gtasks add -l "Work"
```
Prompts for title, notes, and due date.

**Positional / flag mode:**
```bash
gtasks add "Buy groceries"
gtasks add "Finish report" -n "Q4 analysis" -d "2024-12-25"
gtasks add "Call dentist" -d "tomorrow"
gtasks add "Team meeting" -d "Dec 25" -l "Work"
```

**Flags:**
- `-t, --title`: Task title (optional if a positional title is given)
- `-n, --note`: Task notes/description (optional)
- `-d, --due`: Due date (optional, flexible format)
- `-l, --tasklist`: Task list name

**Date Format Examples:**
The date parser supports many formats:
- `2024-12-25` (ISO format)
- `Dec 25, 2024`
- `December 25`
- `tomorrow`
- `next Friday`
- `12/25/2024`

See [dateparse examples](https://github.com/araddon/dateparse#extended-example) for all supported formats.

### Mark Task as Complete

**With task number:**
```bash
gtasks done 1
gtasks done 3 -l "Work"
```

**Interactive mode:**
```bash
gtasks done
gtasks done -l "Personal"
```
Prompts to select a task from the list.

### Delete a Task

**With task number:**
```bash
gtasks rm 2
gtasks rm 1 -l "Shopping"
```

**Interactive mode:**
```bash
gtasks rm
gtasks rm -l "Work"
```
Prompts to select a task to delete.

### View Task Details

**With task number:**
```bash
gtasks info 1
gtasks info 3 -l "Work"
```

**Interactive mode:**
```bash
gtasks info
gtasks info -l "Personal"
```

**Output Example:**
```
Task: Finish report
Status: Needs action
Due: 25 December 2024
Notes: Complete Q4 analysis and submit to manager

Links:
  - https://docs.google.com/document/d/...

View in Google Tasks: https://tasks.google.com/...
```

### Undo a Completed Task

```bash
gtasks undo 1
gtasks undo 1 -l "Work"
```

### Update a Task

```bash
gtasks update 1
gtasks update 1 --title "New title"
gtasks update 1 --note "Updated note" --due tomorrow -l "Work"
```

### Clear Completed Tasks

```bash
gtasks clear
gtasks clear -l "Work" --force
```

## Common Workflows

### Quick Task Creation
When a user says "add a task to my work list":
```bash
gtasks add "Task title" -l "Work"
```

### Check Today's Tasks
```bash
gtasks ls --sort=due
```

### Complete Multiple Tasks
```bash
gtasks done -l "Work"
# Interactive prompt appears, select task
gtasks done -l "Work"
# Repeat as needed
```

### View All Tasks Across Lists
Run ls for each list, or first list all task lists:
```bash
gtasks tasklists
gtasks ls -l "Work"
gtasks ls -l "Personal"
```

### Export Tasks
```bash
gtasks ls --format=json > tasks.json
gtasks ls --format=csv > tasks.csv
```

## Best Practices

1. **Always check authentication first**: If commands fail with authentication errors, run `gtasks login`

2. **Use task list flag for automation**: When scripting or when the user specifies a list name, use `-l` flag to avoid interactive prompts

3. **Leverage flexible date parsing**: The `--due` flag accepts natural language dates like "tomorrow", "next week", etc.

4. **Use appropriate output format**:
   - Table format for human-readable output
   - JSON for parsing/integration with other tools
   - CSV for spreadsheet import

5. **Task numbers are ephemeral**: Task numbers change when tasks are added, completed, or deleted. Always view the list first to get current numbers.

6. **Handle missing lists gracefully**: If a user specifies a non-existent list name, the command will error. Always verify list names first with `gtasks tasklists`.

## Error Handling

Common errors and solutions:

- **"Failed to get service"** or **Authentication errors**:
  - First, ensure environment variables are set: `echo $GTASKS_CLIENT_ID`
  - If variables are not set, export them (see Prerequisites section)
  - Then run `gtasks login` to authenticate
- **"tasklist ... not found"**: The specified list name doesn't exist. Use `gtasks tasklists` to see available lists
- **"task ... not found"**: The task number is invalid. Use `gtasks ls` to see current task numbers
- **"Date format incorrect"**: The date string couldn't be parsed. Use formats like "2024-12-25", "tomorrow", or "Dec 25"

## Examples

### Example 1: Create a shopping list and add items
```bash
gtasks tasklists add "Shopping"
gtasks add "Milk" -l "Shopping"
gtasks add "Bread" -l "Shopping"
gtasks add "Eggs" -l "Shopping"
```

### Example 2: Review and complete work tasks
```bash
gtasks ls -l "Work" --sort=due
gtasks done 1 -l "Work"
```

### Example 3: Add task with deadline
```bash
gtasks add "Submit proposal" -l "Work" -n "Include budget and timeline" -d "next Friday"
```

### Example 4: Export completed tasks
```bash
gtasks ls --completed --format=json -l "Work" > completed_work.json
```

## Tips for Agents

### Before Running Any Commands

1. **Check gtasks installation first**:
   ```bash
   # Try to run gtasks version check
   gtasks --version 2>/dev/null || gtasks.exe --version 2>/dev/null
   ```
   If this fails, inform the user that gtasks is not installed and provide installation instructions from the Prerequisites section.

2. **Verify environment variables are set**:
   ```bash
   # Check if variables exist (macOS/Linux)
   [ -n "$GTASKS_CLIENT_ID" ] && echo "GTASKS_CLIENT_ID is set" || echo "GTASKS_CLIENT_ID is not set"
   [ -n "$GTASKS_CLIENT_SECRET" ] && echo "GTASKS_CLIENT_SECRET is set" || echo "GTASKS_CLIENT_SECRET is not set"

   # Windows PowerShell
   if ($env:GTASKS_CLIENT_ID) { "GTASKS_CLIENT_ID is set" } else { "GTASKS_CLIENT_ID is not set" }
   if ($env:GTASKS_CLIENT_SECRET) { "GTASKS_CLIENT_SECRET is set" } else { "GTASKS_CLIENT_SECRET is not set" }
   ```

3. **Check authentication status**:
   ```bash
   gtasks tasklists &>/dev/null && echo "Authenticated" || echo "Not authenticated - run 'gtasks login'"
   ```

### General Tips

- When the user mentions "tasks" without specifying a tool, ask if they want to use Google Tasks
- If the user asks about their tasks, first run `gtasks tasklists` to see available lists
- Always confirm which task list to use if not specified by the user
- When creating tasks with dates, prefer explicit date formats (YYYY-MM-DD) over relative terms for clarity
- Remember that task numbers are 1-indexed and change after modifications
- If a command requires interaction but you're running non-interactively, use flags to provide all required information
