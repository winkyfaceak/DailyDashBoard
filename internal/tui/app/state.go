package app

func (m Model) clearContent() Model {
	m.stopAudio()
	m.audio.RequestID++

	m = m.clearResultContent()
	m.audio.Starting = false

	return m
}

func (m Model) clearResultContent() Model {
	m.content = contentState{}
	m.links = linkState{}
	m.math = mathState{}
	m.trivia = triviaState{}
	m.musicTrivia = musicTriviaState{}
	m.audio.URL = ""
	m.audio.Label = ""
	m.audio.AllowDownload = false

	return m
}

func (m *Model) stopAudio() {
	if m.audio.Player != nil {
		m.audio.Player.Stop()
		m.audio.Player = nil
	}

	m.audio.Paused = false
	m.audio.Starting = false
}

func (m Model) currentAudioLabel() string {
	if m.audio.Label != "" {
		return m.audio.Label
	}

	return "audio"
}
