package game

import (
	. "ble/hash"
	"time"
)

func (g *Game) addArtist(name string) *Artist {
	if g.Started {
		return nil
	}
	a := Artist{g.makeArtistId(), name}
	g.Artists[a.Id] = a
	g.ArtistOrder = append(g.ArtistOrder, a.Id)
	return &a
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

func (g *Game) passSequence(artistId string) bool {
	//determine which sequence is being passed
	sequences := g.SequenceByHolder[artistId]
	if len(sequences) < 1 {
		return false
	}
	sequenceId := sequences[0]
	sequence := g.Sequences[sequenceId]
	if sequence == nil {
		return false
	}

	//remove the sequence from its current holder
	g.SequenceByHolder[artistId] = sequences[1:]

	//if everyone has drawn, the sequence is done
	if len(sequence.Drawings) >= len(g.Artists) {
		g.SequencesComplete[sequenceId] = true
		if len(g.SequencesComplete) == len(g.Artists) {
			g.Finished = true
		}
	} else {

		//otherwise, it's on to the next artist.
		nextArtistId := g.NextArtist[artistId]
		g.SequenceByHolder[nextArtistId] = append(g.SequenceByHolder[nextArtistId], sequenceId)

		//add a new drawing for the new artist
		g.addDrawingToSequence(sequence, g.Artists[nextArtistId])
	}
	return true
}

type GameEvent struct {
	time.Time
	Payload interface{}
}

type Game struct {
	Started     bool
	Finished    bool
	Artists     artistMap
	ArtistOrder []string

	NextArtist map[string]string
	PrevArtist map[string]string

	Sequences         sequenceMap
	SequenceByStarter map[string]string
	SequenceByHolder  map[string][]string
	SequencesComplete map[string]bool
	Drawings          map[string]*Drawing

	Events chan<- GameEvent
}

func newGame(ch chan GameEvent) *Game {
	return &Game{
		Started:           false,
		Finished:          false,
		Artists:           make(map[string]Artist),
		ArtistOrder:       make([]string, 0, 0),
		NextArtist:        make(map[string]string),
		PrevArtist:        make(map[string]string),
		Sequences:         make(map[string]*Sequence),
		SequenceByStarter: make(map[string]string),
		SequenceByHolder:  make(map[string][]string),
		SequencesComplete: make(map[string]bool),
		Drawings:          make(map[string]*Drawing),
		Events:            ch}
}

func (g *Game) makeArtistId() string {
	return NewHashEasy().
		WriteStrAnd("artist").
		Nonce().
		WriteIntAnd(len(g.Artists)).String()
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

func (g *Game) viewJSON() interface{} {
	obj := make(map[string]interface{})
	obj["started"] = g.Started
	obj["finished"] = g.Finished
	obj["artistNames"] = g.Artists
	obj["artistOrder"] = g.ArtistOrder
	if g.Started {
		obj["sequenceByHolder"] = g.SequenceByHolder
		obj["drawingsBySequence"] = g.Sequences
		obj["sequencesComplete"] = g.SequencesComplete
	}
	return obj
}
