// Package features contains every dashboard feature. Each feature lives in a
// single file that declares its menu entry, how its data is fetched, how the
// result is formatted, and (when interactive) how it handles key presses.
//
// The app package stays generic: it only knows the Feature and Session types
// declared here.
package features

import (
	tea "charm.land/bubbletea/v2"
)

// Feature describes one dashboard feature.
type Feature struct {
	// MenuLabel is the label shown in the daily menu.
	MenuLabel string
	// Loading is the status message shown while fetching.
	Loading string
	// ErrLabel prefixes fetch errors, e.g. "Fact" -> "Fact error: ...".
	ErrLabel string
	// BoxWidth is the width of the content box the result renders in.
	BoxWidth int
	// Fetch loads the feature's data and returns a Session for it. It runs
	// inside a tea.Cmd goroutine, so it may block on the network.
	Fetch func() (Session, error)
}

// FetchCmd wraps a feature's Fetch in a bubbletea command.
func (f *Feature) FetchCmd() tea.Cmd {
	return func() tea.Msg {
		session, err := f.Fetch()

		return ResultMsg{
			Feature: f,
			Session: session,
			Err:     err,
		}
	}
}

// ResultMsg carries the outcome of any feature fetch.
type ResultMsg struct {
	Feature *Feature
	Session Session
	Err     error
}

// Notice is an error that should be shown to the user as a plain status
// message instead of an "<ErrLabel> error: ..." message.
type Notice string

func (n Notice) Error() string { return string(n) }

// Session is a fetched feature result that can render itself. Sessions may
// additionally implement AudioSession, LinkSession, StorySession,
// ImageSession, or InteractiveSession to opt into extra behavior.
type Session interface {
	View() string
}

// Audio describes a playable audio clip attached to a session.
type Audio struct {
	URL           string
	Label         string
	AllowDownload bool
}

// AudioSession is implemented by sessions that provide audio playback.
type AudioSession interface {
	Session
	Audio() Audio
}

// LinkAction is an external link that can be opened with a single key.
type LinkAction struct {
	Key     string // key that triggers the open, e.g. "o"
	URL     string
	Help    string // help fragment, e.g. "o to open image"
	Message string // status message, e.g. "Opening full artwork image..."
}

// LinkSession is implemented by sessions that expose key-triggered links.
type LinkSession interface {
	Session
	Links() []LinkAction
}

// StorySession is implemented by sessions with numbered story links (1-9).
type StorySession interface {
	Session
	StoryLinks() []string
}

// ImageSession is implemented by sessions that flush a raw terminal image
// (kitty/iTerm2 graphics) once, then swap their view to the final text.
type ImageSession interface {
	Session
	// RawImage returns the raw escape sequence to write, or "" if none.
	RawImage() string
	// ImageShown is called after the raw sequence has been flushed.
	ImageShown()
}

// KeyResponse tells the app how an interactive session handled a key.
type KeyResponse struct {
	// Handled reports whether the session consumed the key; when true the
	// app sets its status message to Message (possibly clearing it).
	Handled bool
	Message string
	// Clear asks the app to stop audio and discard the session.
	Clear bool
	// ToggleAudio asks the app to toggle the session's audio playback.
	ToggleAudio bool
}

// InteractiveSession is implemented by sessions that capture the keyboard,
// such as quizzes with typed answers.
type InteractiveSession interface {
	Session
	// Help returns the full help line shown while the session is active.
	Help() string
	// HandleKey processes a key press, mutating the session as needed.
	HandleKey(msg tea.KeyPressMsg) KeyResponse
}

// Daily returns every feature in the daily menu, in menu order.
func Daily() []*Feature {
	return []*Feature{
		RandomWord(),
		RandomFact(),
		RandomMath(),
		RandomTrivia(),
		DailyTarot(),
		RandomArtwork(),
		DailyNews(),
		RandomMusic(),
		MusicTrivia(),
		RandomRecipe(),
	}
}
