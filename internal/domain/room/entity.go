package room

import (
	"fmt"
)

type ID string

type Type string

const (
	Standard Type = "STANDARD"
	Superior Type = "SUPERIOR"
	Suite    Type = "SUITE"
)

type Room struct {
	ID   ID
	Type Type
}

func New(id ID, roomType Type) Room {
	return Room{
		ID:   id,
		Type: roomType,
	}
}

type RoomConfig struct {
	PriceCents  int
	Description []string
}

func (r Room) GetConfig() RoomConfig {
	switch r.Type {
	case Standard:
		return RoomConfig{
			PriceCents:  5000,
			Description: []string{"Lit 1 place", "Wifi", "TV"},
		}
	case Superior:
		return RoomConfig{
			PriceCents:  10000,
			Description: []string{"Lit 2 places", "Wifi", "TV écran plat", "Minibar", "Climatiseur"},
		}
	case Suite:
		return RoomConfig{
			PriceCents:  20000,
			Description: []string{"Lit 2 places", "Wifi", "TV écran plat", "Minibar", "Climatiseur", "Baignoire", "Terrasse"},
		}
	default:
		return RoomConfig{PriceCents: 0, Description: []string{}}
	}
}

func (r Room) String() string {
	conf := r.GetConfig()
	return fmt.Sprintf("Chambre %s (%s) - %.2f€", r.ID, r.Type, float64(conf.PriceCents)/100.0)
}
