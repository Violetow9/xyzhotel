package application

import (
	"context"
	"xyzhotel/internal/domain/room"
)

type RoomInfoDTO struct {
	Type        string   `json:"type"`
	PriceCents  int      `json:"price_cents"`
	Description []string `json:"description"`
}

type ListRoomsHandler struct {
	// Pas besoin de repository ici car les infos descriptives sont statiques dans le code (Config)
	// Si tu voulais lister les chambres physiques libres, on utiliserait room.Repository
}

func (h ListRoomsHandler) Handle(ctx context.Context) ([]RoomInfoDTO, error) {
	types := []room.Type{
		room.Standard,
		room.Superior,
		room.Suite,
	}

	var infos []RoomInfoDTO

	for _, t := range types {
		tempRoom := room.New("temp", t)
		config := tempRoom.GetConfig()

		infos = append(infos, RoomInfoDTO{
			Type:        string(t),
			PriceCents:  config.PriceCents,
			Description: config.Description,
		})
	}

	return infos, nil
}
