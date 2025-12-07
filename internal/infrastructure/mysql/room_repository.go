package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"xyzhotel/internal/domain/room"
	"xyzhotel/internal/infrastructure/mysql/sqldb"
)

type RoomRepository struct {
	db *sql.DB
	q  *sqldb.Queries
}

func NewRoomRepository(db *sql.DB) *RoomRepository {
	return &RoomRepository{
		db: db,
		q:  sqldb.New(db),
	}
}

func (r *RoomRepository) FindByID(ctx context.Context, id string) (*room.Room, error) {
	row, err := r.q.GetRoom(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("room not found: %w", err)
		}
		return nil, err
	}

	return r.mapToDomain(row)
}

func (r *RoomRepository) ListAll(ctx context.Context) ([]*room.Room, error) {
	rows, err := r.q.ListRooms(ctx)
	if err != nil {
		return nil, err
	}

	var rooms []*room.Room
	for _, row := range rows {
		domainRoom, err := r.mapToDomain(row)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, domainRoom)
	}
	return rooms, nil
}

func (r *RoomRepository) Save(ctx context.Context, rm room.Room) error {
	return r.q.CreateRoom(ctx, sqldb.CreateRoomParams{
		RoomNumber: string(rm.ID),
		Type:       string(rm.Type),
	})
}

func (r *RoomRepository) mapToDomain(row sqldb.Room) (*room.Room, error) {
	var roomType room.Type
	switch row.Type {
	case "STANDARD":
		roomType = room.Standard
	case "SUPERIOR":
		roomType = room.Superior
	case "SUITE":
		roomType = room.Suite
	default:
		return nil, fmt.Errorf("unknown room type in db: %s", row.Type)
	}

	entity := room.New(room.ID(row.RoomNumber), roomType)
	return &entity, nil
}

func (r *RoomRepository) Update(ctx context.Context, rm room.Room) error {
	return nil
}

func (r *RoomRepository) Delete(ctx context.Context, id room.ID) error {
	return nil
}
