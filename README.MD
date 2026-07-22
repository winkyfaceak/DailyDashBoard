# Daily Dashboard

Daily Dashboard is a small terminal app that serves up quick daily bits: a word,
a fact, math question, trivia question, tarot card, artwork, news, music, music
trivia, and a recipe. It is built with Bubble Tea and kept as a terminal-first
project.

The code is arranged so each part of the app has a clear home. Every feature
lives in exactly two files: an API client in `internal/daily` and a feature file
in `internal/tui/features`. The app shell is generic and picks features up from
a registry, so adding a feature never means editing the app package.

## Run It

```sh
go run ./cmd/dailydashboard
```

## Test It

```sh
go test ./...
```

## Project Structure

```text
.
|-- cmd/dailydashboard
|   `-- main.go
|-- internal
|   |-- daily
|   `-- tui
|       |-- app
|       |-- audio
|       |-- features
|       |-- platform
|       |-- terminalimage
|       |-- textutil
|       `-- theme
|-- go.mod
|-- go.sum
`-- README.MD
```

### `cmd/dailydashboard`

The executable entry point. This should stay tiny: create the Bubble Tea program,
run it, and report startup errors.

### `internal/daily`

The content/API layer. This package talks to external services and turns their
responses into clean Go structs. It should not know anything about Bubble Tea,
terminal styling, menus, or keyboard shortcuts.

HTTP plumbing is shared: `client.go` provides `getJSON`, which handles request
creation, headers, status checks, and consistent error messages. Feature files
only supply the endpoint and a decode function.

Good things to put here:

- API clients
- response decoding
- response cleanup
- random selection logic
- package-level tests for parsing and behavior

### `internal/tui/features`

Every dashboard feature in a single file each. `feature.go` defines the shared
types:

- `Feature` — menu label, loading message, error label, box width, and a
  `Fetch` function
- `Session` — a fetched result that renders itself with `View()`
- Optional session interfaces — `AudioSession`, `LinkSession`, `StorySession`,
  `ImageSession`, and `InteractiveSession` for quizzes that capture the
  keyboard

`Daily()` is the registry: the list of features, in menu order.

A feature file declares its `Feature`, its session type, and its formatting.
Interactive sessions (math, trivia, music trivia) also carry their quiz state
and key handling, which makes them directly unit-testable — see the
`*_session_test.go` files.

### `internal/tui/app`

The generic Bubble Tea shell. It renders menus, routes messages, plays audio,
opens links, and hosts whichever session is active. It never references an
individual feature; everything it needs comes from the `Feature` struct and the
session interfaces.

### `internal/tui/theme`

The shared color palette and base lipgloss styles. Both `app` and `features`
derive their styles from here, so the look stays consistent and colors are
defined once.

### `internal/tui/audio`

Audio playback support. This keeps process management and platform-specific
pause/resume behavior away from the app model.

### `internal/tui/terminalimage`

Terminal image rendering. Artwork, music, tarot, and recipes all use this
package instead of each feature carrying its own image code.

### `internal/tui/platform`

Small operating-system helpers, like opening an external URL.

### `internal/tui/textutil`

Shared text helpers for wrapping, truncating, and indenting terminal text.

## Adding A New Daily Feature

Here is the usual path.

### 1. Add the data client

Create a file in `internal/daily`, for example:

```text
internal/daily/weather.go
```

Keep this package focused on data. Define the structs the rest of the app should
use, fetch the upstream data with `getJSON`, clean it up, and return useful
errors.

Try to add tests beside it:

```text
internal/daily/weather_test.go
```

### 2. Add the feature file

Create a matching file in `internal/tui/features`, for example:

```text
internal/tui/features/weather.go
```

That file needs two things:

- a constructor returning a `*Feature` (menu label, loading message, error
  label, box width, and the `Fetch` function)
- a session type with a `View()` method that formats the result

Opt into extra behavior by implementing the optional interfaces: `Audio()` for
playback, `Links()` for key-triggered links, `RawImage()`/`ImageShown()` (or
embed `imageState`) for terminal graphics, and `Help()`/`HandleKey()` for
interactive quizzes.

### 3. Register it

Add the constructor to the `Daily()` registry in
`internal/tui/features/feature.go`. That is the only wiring step — the menu
item, fetch command, result handling, help text, and rendering all follow from
it.

### 4. Run the checks

```sh
gofmt -w ./cmd ./internal
go test ./...
```

If those pass, the new feature is probably sitting in the right places.
