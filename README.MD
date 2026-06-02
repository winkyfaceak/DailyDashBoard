# Daily Dashboard

Daily Dashboard is a small terminal app that serves up quick daily bits: a word,
a fact, tarot card, artwork, news, music, and a recipe. It is built with Bubble
Tea and kept as a terminal-first project.

The code is arranged so each part of the app has a clear home. If you want to add
another daily feature, you should not need to thread your way through one giant
package.

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
|       `-- textutil
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

Good things to put here:

- API clients
- response decoding
- response cleanup
- random selection logic
- package-level tests for parsing and behavior

### `internal/tui/app`

The Bubble Tea app shell. This package owns the model, menus, key handling,
state transitions, and the final view composition.

Good things to put here:

- app state
- menu actions
- keyboard shortcuts
- result message handling
- deciding which content box appears on screen

### `internal/tui/features`

The terminal-facing feature layer. Each feature gets a file here with its fetch
command, result message, and formatter.

For example:

- `fact.go`
- `dictionary.go`
- `artwork.go`
- `music.go`
- `news.go`
- `recipe.go`

This package is the bridge between `internal/daily` and `internal/tui/app`: it
fetches data from `daily`, formats it for the terminal, and returns Bubble Tea
messages to the app.

### `internal/tui/audio`

Audio playback support. This keeps process management and platform-specific
pause/resume behavior away from the app model.

### `internal/tui/terminalimage`

Terminal image rendering. Artwork, music, and recipes all use this package
instead of each feature carrying its own image code.

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
use, fetch the upstream data, clean it up, and return useful errors.

Try to add tests beside it:

```text
internal/daily/weather_test.go
```

### 2. Add the terminal feature

Create a matching file in `internal/tui/features`, for example:

```text
internal/tui/features/weather.go
```

That file usually needs three things:

- a result message type, like `WeatherResultMsg`
- a fetch command, like `FetchWeatherCmd`
- a formatter, like `FormatWeather`

The command calls `internal/daily`; the formatter turns the result into a string
ready for the terminal.

### 3. Wire it into the app

Update `internal/tui/app`:

- Add a new `menuAction` in `model.go`.
- Add the visible menu item in `InitialModel`.
- Add any content or link state the feature needs.
- Handle the menu action in `interaction.go`.
- Handle the result message in `update.go`.
- Render the content in `view.go`.

This is the only place that should know how the new feature fits into the app
workflow.

### 4. Run the checks

```sh
gofmt -w ./cmd ./internal
go test ./...
```

If those pass, the new feature is probably sitting in the right places.
