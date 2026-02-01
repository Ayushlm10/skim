# Phase 10: Fullscreen Preview Mode - Current State

**Branch:** `feature/fullscreen-preview`
**Status:** Work in Progress - Performance Issues Unresolved
**Date:** 2026-02-01

## Goal

Allow users to toggle the preview panel to fullscreen for distraction-free reading by pressing `f`.

## What Works

1. **Fullscreen toggle is implemented** - Press `f` when preview is focused
2. **Toggle is instant** when re-rendering is disabled
3. **File tree hides** in fullscreen mode
4. **Tab exits fullscreen** and returns to file tree
5. **Status bar hints** update to show `f fullscreen` / `f exit fullscreen`
6. **Help overlay** updated with `f` key documentation

## The Problem

When toggling fullscreen, the content needs to re-render with Glamour to use the new width for proper word-wrapping. Without re-rendering, the text stays wrapped at the narrow 25% width, wasting screen real estate in fullscreen mode.

### Approaches Tried

1. **Synchronous re-render in SetSize()** - Original approach
   - Result: Sluggish, blocking UI on every toggle
   - ~10ms render time per toggle, but feels much slower

2. **Async re-render via tea.Cmd** - Return command from Rerender()
   - Result: Still sluggish, command execution blocks before yielding
   
3. **Debounced re-render with tea.Tick** - 150ms delay before re-render
   - Result: Even worse - added delay plus still felt slow

4. **No re-render at all** - Just resize viewport
   - Result: Instant toggle! But text doesn't reflow, wastes space

5. **Delayed re-render with 1ms tea.Tick** - Let UI update first
   - Result: Still not instant, something still blocking

### Benchmark Results

Glamour rendering is actually fast in isolation:
- AGENTS.md (~1KB): <1ms
- implementation.md (~18KB): ~10ms
- viewport.SetContent: <0.1ms

The slowness is NOT in Glamour rendering itself.

## Suspected Root Cause

The issue may be in how Bubble Tea processes messages or how the view is rendered. Even though commands are "async", something in the update/view cycle is causing perceived lag.

Possible causes to investigate:
1. View() function re-rendering entire UI on every message
2. Lipgloss styling calculations on large content
3. Terminal output buffering/flushing
4. Something in the message dispatch loop

## Files Modified

- `internal/app/model.go` - Added `fullscreenPreview bool` field
- `internal/app/update.go` - Handle `f` key, resize logic, message forwarding
- `internal/app/view.go` - Conditional panel rendering, status bar hints
- `internal/components/preview/preview.go` - Rerender() method, message types
- `internal/components/help/help.go` - Added `f` key to Preview section

## Code State

Current code has:
- `RerenderCompleteMsg` and `rerenderTickMsg` message types
- `Rerender()` method using `tea.Tick` for delayed trigger
- `doRerender()` function for actual Glamour rendering
- Message forwarding to preview component for tick handling

## Next Steps to Try

1. **Profile the actual bottleneck** - Add timing logs to pinpoint exactly where time is spent

2. **Try goroutine with channel** - Bypass tea.Cmd entirely, use raw goroutine

3. **Cache rendered content by width** - Store renders at common widths (25%, 100%)

4. **Investigate Bubble Tea internals** - Check if there's a way to force immediate render

5. **Accept the trade-off** - Keep toggle instant, add manual "r" key to re-render

## How to Test

```bash
git checkout feature/fullscreen-preview
go build -o skim
./skim specs/
# Select a file, press Tab to focus preview, press f to toggle
```

## Related Specs

- `specs/plan.md` - Phase 10 design
- `specs/implementation.md` - Phase 10 tasks (marked complete but performance issue remains)
