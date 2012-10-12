package main

import (
	. "ble/hash"
	"encoding/json"
	"fmt"
	"time"
)

func Marshal(o interface{}) ([]byte, error) {
	return json.MarshalIndent(o, "", "\t")
}

func main() {
	ch := make(chan GameEvent)
	g := NewGame(ch)
	gJson, _ := Marshal(g.view())
	fmt.Println(string(gJson))

	g.addArtist("asdf")
	g.addArtist("w0pak")
	g.addArtist("without_spaces")
	gJson, _ = Marshal(g.view())
	fmt.Println(string(gJson))

	g.start()
	gJson, _ = Marshal(g.view())
	fmt.Println(string(gJson))
}

func (g *Game) start() bool {
	//game is now started
	if g.Started {
		return false
	}
	g.Started = true

	//define predecessor, successor
	for i := range g.ArtistOrder {
		a0Id := g.ArtistOrder[i]
		a1Id := g.ArtistOrder[(i+1)%len(g.ArtistOrder)]
		g.NextArtist[a0Id] = a1Id
		g.PrevArtist[a1Id] = a0Id
	}

	//create one sequence per artist
	for _, id := range g.ArtistOrder {
		a := g.Artists[id]
		s := g.makeSequence(a)
		_ = g.addDrawingToSequence(s, a)
	}
	return true

}

func (g *Game) passSequence(id string) bool {
	seq := g.Sequences[id]
	if seq == nil {
		return false
	}

	return true
}

type GameEvent struct {
	time.Time
	Payload interface{}
}

type Artist struct {
	Id   string
	Name string
}

type Drawing struct {
	Id       string
	ArtistId string
	//....
}

type Sequence struct {
	Id       string
	Drawings []*Drawing
}

type artistMap map[string]Artist
type sequenceMap map[string]*Sequence

type Game struct {
	Artists     artistMap
	ArtistOrder []string

	Started           bool
	NextArtist        map[string]string
	PrevArtist        map[string]string
	Sequences         sequenceMap
	SequenceByStarter map[string]string
	SequenceByHolder  map[string][]string
	Drawings          map[string]*Drawing

	Events chan<- GameEvent
}

func NewGame(ch chan GameEvent) *Game {
	return &Game{
		Artists:           make(map[string]Artist),
		ArtistOrder:       make([]string, 0, 0),
		Started:           false,
		NextArtist:        make(map[string]string),
		PrevArtist:        make(map[string]string),
		Sequences:         make(map[string]*Sequence),
		SequenceByStarter: make(map[string]string),
		SequenceByHolder:  make(map[string][]string),
		Drawings:          make(map[string]*Drawing),
		Events:            ch}
}

func (g *Game) makeArtistId() string {
	return NewHashEasy().
		WriteStrAnd("artist").
		Nonce().
		WriteIntAnd(len(g.Artists)).String()
}

func (g *Game) addArtist(name string) *Artist {
	if g.Started {
		return nil
	}
	a := Artist{g.makeArtistId(), name}
	g.Artists[a.Id] = a
	g.ArtistOrder = append(g.ArtistOrder, a.Id)
	return &a
}

func (g *Game) makeSequenceId() string {
	return NewHashEasy().
		WriteStrAnd("sequence").
		Nonce().
		WriteIntAnd(len(g.Sequences)).String()
}

func (g *Game) makeSequence(firstArtist Artist) *Sequence {
	s := Sequence{g.makeSequenceId(), []*Drawing{}}
	g.Sequences[s.Id] = &s
	g.SequenceByStarter[firstArtist.Id] = s.Id
	g.SequenceByHolder[firstArtist.Id] = []string{s.Id}
	return &s
}

func (g *Game) makeDrawingId() string {
	return NewHashEasy().
		WriteStrAnd("drawing").
		Nonce().
		WriteIntAnd(len(g.Drawings)).String()
}

func (g *Game) addDrawingToSequence(s *Sequence, a Artist) *Drawing {
	d := Drawing{g.makeDrawingId(), a.Id}
	g.Drawings[d.Id] = &d
	s.Drawings = append(s.Drawings, &d)
	return &d
}

func (m artistMap) MarshalJSON() ([]byte, error) {
	obj := make(map[string]string)
	for id, artist := range m {
		obj[id] = artist.Name
	}
	return json.Marshal(obj)
}

func (s sequenceMap) MarshalJSON() ([]byte, error) {
	obj := make(map[string][]string)
	for id, seq := range s {
		drawingIds := make([]string, len(seq.Drawings), len(seq.Drawings))
		for i, d := range seq.Drawings {
			drawingIds[i] = d.Id
		}
		obj[id] = drawingIds
	}
	return json.Marshal(obj)
}

type gameView struct{ *Game }

func (g *Game) view() gameView {
	return gameView{g}
}

func (g gameView) MarshalJSON() ([]byte, error) {
	obj := make(map[string]interface{})
	obj["started"] = g.Started
	obj["artistNames"] = g.Artists
	obj["artistOrder"] = g.ArtistOrder
	if g.Started {
		obj["sequenceByHolder"] = g.SequenceByHolder
		obj["drawingsBySequence"] = g.Sequences
	}
	return json.Marshal(obj)
}
