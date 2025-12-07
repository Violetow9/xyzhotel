package application

import (
	"context"
	"testing"

	"xyzhotel/internal/domain/room"
)

func TestListRoomsHandler_Success(t *testing.T) {
	handler := &ListRoomsHandler{}

	infos, err := handler.Handle(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(infos) != 3 {
		t.Errorf("expected 3 room types, got %d", len(infos))
	}

	if infos[0].Type != string(room.Standard) {
		t.Errorf("expected first type to be %s, got %s", room.Standard, infos[0].Type)
	}
	if infos[0].PriceCents != 5000 {
		t.Errorf("expected price 5000, got %d", infos[0].PriceCents)
	}

	if infos[1].Type != string(room.Superior) {
		t.Errorf("expected second type to be %s, got %s", room.Superior, infos[1].Type)
	}
	if infos[1].PriceCents != 10000 {
		t.Errorf("expected price 10000, got %d", infos[1].PriceCents)
	}

	if infos[2].Type != string(room.Suite) {
		t.Errorf("expected third type to be %s, got %s", room.Suite, infos[2].Type)
	}
	if infos[2].PriceCents != 20000 {
		t.Errorf("expected price 20000, got %d", infos[2].PriceCents)
	}
}
