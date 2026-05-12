# Pull Request: Refactor Base TUI Components to Bubbletea v2 and Introduce Async Loading Pattern

## Overview
This PR migrates the core TUI components of the Harbor CLI to use the modern `Bubbletea v2` and `lipgloss v2` ecosystem. It also introduces a standard asynchronous loading state pattern (with animated spinners) for all base list and grid views to vastly improve the user experience when performing network requests.

This addresses **Issue #821** (Bubbletea Model Refactor) and provides the foundation for closing **Issue #756**.

## Changes Made
1. **Dependency Upgrades**: Migrated all import paths from `github.com/charmbracelet/*` (v1) to `charm.land/*` (v2) vanity paths across 68+ files.
2. **Mechanical API Fixes**:
   - `View() string` updated to `View() tea.View`
   - `tea.KeyMsg` updated to `tea.KeyPressMsg`
   - `case " ":` updated to `case "space":`
   - Removed `tea.WithAltScreen()` from 17 `NewProgram` initialization calls.
   - Fixed `list.DefaultStyles()` and `viewport` initialization to match new functional options API.
3. **Async Loading State Pattern**: 
   - Refactored `tablelist`, `selection`, `multiselect`, and `tablegrid` models to support asynchronous data fetching.
   - Introduced `loading`, `err`, `spinner`, and `fetchCmd` state variables into each Model.
   - Added `NewModelWithFetch()` / `NewWithFetch()` constructors that accept a `FetchFn` closure containing the API call.
   - When using the fetch constructor, `Init()` immediately starts a batch command with the data fetch and spinner animation. `View()` properly branches to show the spinner and hides the UI until the `DataLoadedMsg` triggers `Update()`.
   - **Backward Compatibility**: Fully maintained the existing synchronous `NewModel()` signatures to avoid breaking the 53+ command handlers that currently rely on them.

*Note: We maintained the `lipgloss` v1 import for `instance/update` and `replication/policies/create` due to their internal dependency on `huh` v1.0.0, which still wraps the legacy API.*

## How to Test
1. Compile the project using `go build ./...` and verify zero errors.
2. Run standard commands (e.g. `harbor user list`) to confirm synchronous rendering still works perfectly with no regressions.
3. Use the new pattern in a future or test command by swapping `tablelist.NewModel` with `tablelist.NewModelWithFetch(func() ([]table.Row, error) { ... })` and confirm the `bubbles/spinner` cleanly loads before displaying the table.
