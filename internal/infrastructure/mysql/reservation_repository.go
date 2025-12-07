package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"xyzhotel/internal/infrastructure/mysql/sqldb"

	"xyzhotel/internal/domain/reservation"

	"github.com/google/uuid"
)

type ReservationRepository struct {
	db *sql.DB
	q  *sqldb.Queries
}

func NewReservationRepository(db *sql.DB) *ReservationRepository {
	return &ReservationRepository{
		db: db,
		q:  sqldb.New(db),
	}
}

func (r *ReservationRepository) IsRoomAvailable(ctx context.Context, roomID string, checkIn time.Time, nights int) (bool, error) {
	endDate := checkIn.AddDate(0, 0, nights)

	count, err := r.q.IsRoomAvailable(ctx, sqldb.IsRoomAvailableParams{
		RoomNumber: roomID,
		StartDate:  checkIn,
		EndDate:    endDate,
	})

	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (r *ReservationRepository) Save(ctx context.Context, res *reservation.Reservation) error {
	return r.q.CreateReservation(ctx, sqldb.CreateReservationParams{
		ID:               res.ID.String(),
		CustomerID:       res.CustomerID.String(),
		RoomNumber:       res.RoomNumber,
		CheckInDate:      res.CheckInDate,
		Nights:           int32(res.AmountOfNights),
		Status:           string(res.State),
		TotalAmountCents: int32(res.TotalAmountCents),
		PaidAmountCents:  int32(res.PaidAmountCents),
	})
}

func (r *ReservationRepository) FindByID(ctx context.Context, id reservation.ID) (*reservation.Reservation, error) {
	row, err := r.q.GetReservation(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("reservation not found: %w", err)
		}
		return nil, err
	}

	return r.mapToDomain(row)
}

func (r *ReservationRepository) ListAll(ctx context.Context) ([]reservation.Reservation, error) {
	rows, err := r.q.ListReservations(ctx)
	if err != nil {
		return nil, err
	}
	return r.mapRowsToDomain(rows)
}

func (r *ReservationRepository) FindHistoryByRoom(ctx context.Context, roomID string) ([]reservation.Reservation, error) {
	rows, err := r.q.GetReservationsByRoom(ctx, roomID)
	if err != nil {
		return nil, err
	}
	return r.mapRowsToDomain(rows)
}

func (r *ReservationRepository) mapRowsToDomain(rows []sqldb.Reservation) ([]reservation.Reservation, error) {
	var result []reservation.Reservation
	for _, row := range rows {
		res, err := r.mapToDomain(row)
		if err != nil {
			return nil, err
		}
		result = append(result, *res)
	}
	return result, nil
}

func (r *ReservationRepository) mapToDomain(row sqldb.Reservation) (*reservation.Reservation, error) {
	id, err := uuid.Parse(row.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid reservation id: %w", err)
	}
	custID, err := uuid.Parse(row.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("invalid customer id: %w", err)
	}

	return &reservation.Reservation{
		ID:               id,
		CustomerID:       custID,
		RoomNumber:       row.RoomNumber,
		CheckInDate:      row.CheckInDate,
		AmountOfNights:   int(row.Nights),
		State:            reservation.State(row.Status),
		TotalAmountCents: int(row.TotalAmountCents),
		PaidAmountCents:  int(row.PaidAmountCents),
	}, nil
}

func (r *ReservationRepository) Update(ctx context.Context, res *reservation.Reservation) error {
	return r.q.UpdateReservation(ctx, sqldb.UpdateReservationParams{
		Status:          string(res.State),
		PaidAmountCents: int32(res.PaidAmountCents),
		ID:              res.ID.String(),
	})
}

func (r *ReservationRepository) Delete(ctx context.Context, id reservation.ID) error {
	return nil
}
