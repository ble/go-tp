package game

import "encoding/json"


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

