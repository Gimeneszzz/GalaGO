package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const sampleRate = 44100

type soundID int

const (
	soundShoot soundID = iota
	soundEnemyShoot
	soundDamage
	soundExplosion
	soundAppear
	soundVictory
	soundDefeat
)

type soundPlayer struct {
	context *audio.Context
	effects map[soundID][]byte
	music   *audio.Player
}

func newSoundPlayer() *soundPlayer {
	player := &soundPlayer{
		context: audio.NewContext(sampleRate),
		effects: map[soundID][]byte{
			soundShoot:      tone(720, 0.055, 0.25),
			soundEnemyShoot: tone(260, 0.12, 0.32),
			soundDamage:     tone(120, 0.16, 0.34),
			soundExplosion:  sweep(180, 55, 0.22, 0.44),
			soundAppear:     sweep(320, 720, 0.18, 0.28),
			soundVictory:    victoryTone(),
			soundDefeat:     sweep(180, 70, 0.38, 0.35),
		},
	}
	player.music = player.context.NewPlayerFromBytes(ambientLoop())
	player.music.SetVolume(0.36)
	return player
}

func (player *soundPlayer) play(id soundID) {
	data := player.effects[id]
	if len(data) == 0 {
		return
	}

	audioPlayer := player.context.NewPlayerFromBytes(data)
	audioPlayer.SetVolume(0.45)
	audioPlayer.Play()
}

func tone(frequency, seconds, volume float64) []byte {
	samples := int(sampleRate * seconds)
	data := make([]byte, samples*4)
	for sample := 0; sample < samples; sample++ {
		t := float64(sample) / sampleRate
		envelope := 1 - float64(sample)/float64(samples)
		value := int16(math.Sin(2*math.Pi*frequency*t) * volume * envelope * 32767)
		writeStereoSample(data, sample, value)
	}
	return data
}

func sweep(startFrequency, endFrequency, seconds, volume float64) []byte {
	samples := int(sampleRate * seconds)
	data := make([]byte, samples*4)
	for sample := 0; sample < samples; sample++ {
		progress := float64(sample) / float64(samples)
		frequency := startFrequency + (endFrequency-startFrequency)*progress
		t := float64(sample) / sampleRate
		envelope := 1 - progress
		value := int16(math.Sin(2*math.Pi*frequency*t) * volume * envelope * 32767)
		writeStereoSample(data, sample, value)
	}
	return data
}

func victoryTone() []byte {
	first := tone(520, 0.12, 0.24)
	second := tone(660, 0.12, 0.24)
	third := tone(880, 0.22, 0.28)
	data := make([]byte, 0, len(first)+len(second)+len(third))
	data = append(data, first...)
	data = append(data, second...)
	data = append(data, third...)
	return data
}

func ambientLoop() []byte {
	seconds := 8.0
	samples := int(sampleRate * seconds)
	data := make([]byte, samples*4)
	notes := []float64{220, 277, 330, 277, 196, 247, 294, 247}
	for sample := 0; sample < samples; sample++ {
		t := float64(sample) / sampleRate
		stepDuration := seconds / float64(len(notes))
		noteIndex := int(t/stepDuration) % len(notes)
		stepT := math.Mod(t, stepDuration)
		envelope := math.Min(1, stepT/0.08) * math.Min(1, (stepDuration-stepT)/0.16)
		bass := math.Sin(2*math.Pi*notes[noteIndex]*0.5*t) * 0.10
		lead := math.Sin(2*math.Pi*notes[noteIndex]*t) * 0.12 * envelope
		shimmer := math.Sin(2*math.Pi*notes[noteIndex]*2*t) * 0.035 * envelope
		value := bass + lead + shimmer
		writeStereoSample(data, sample, int16(value*32767))
	}
	return data
}

func (player *soundPlayer) updateMusic(muted bool) {
	if player.music == nil {
		return
	}
	if muted {
		player.music.Pause()
		return
	}
	if player.music.IsPlaying() {
		return
	}
	_ = player.music.Rewind()
	player.music.Play()
}

func writeStereoSample(data []byte, sample int, value int16) {
	offset := sample * 4
	data[offset] = byte(value)
	data[offset+1] = byte(value >> 8)
	data[offset+2] = byte(value)
	data[offset+3] = byte(value >> 8)
}

func (g *Game) playSound(id soundID) {
	if g.sounds == nil || g.muted {
		return
	}
	g.sounds.play(id)
}
