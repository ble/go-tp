package main


type Room struct {
  Id, Name string
}

type Nexus struct {
  Rooms map[string]*Room
}

func (n *Nexus) GetRoom(id String) {
  return n.Rooms[id]
}
