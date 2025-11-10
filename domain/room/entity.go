package room

import (
	"fmt"

	"github.com/google/uuid"
)

type ID = uuid.UUID

type Room struct {
	ID   ID
	Type Type
}

func New(ID ID, roomType Type) Room {
	return Room{
		ID:   ID,
		Type: roomType,
	}
}

func NewWithTypeMapper(ID ID, roomType string) Room {
	var rType Type
	switch roomType {
	case "Standard":
		rType = Standard
	case "Deluxe":
		rType = Deluxe
	case "Suite":
		rType = Suite
	default:
		rType = Standard
	}
	return Room{
		ID:   ID,
		Type: rType,
	}
}

type Type struct {
	Name        string
	Price       float64
	Description []string
}

func NewRoomType(name string, price float64, desc ...string) Type {
	desc = append([]string{fmt.Sprintf("%s: %.2f euros par nuit", name, price)}, desc...)
	return Type{
		Name:        name,
		Price:       price,
		Description: desc,
	}
}

var (
	Standard = NewRoomType("Standard", 50, "Lit 1 place", "Wifi", "TV")
	Deluxe   = NewRoomType("Deluxe", 100, "Lit 2 places", "Wifi", "TV Ecran Plat", "Mini-bar", "Climatiseur")
	Suite    = NewRoomType("Suite", 200, "Lit 2 places", "Wifi", "TV Ecran Plat", "Mini-bar", "Climatiseur", "Baignore", "Terrasse")
)
