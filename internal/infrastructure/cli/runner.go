package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"xyzhotel/internal/application"
	"xyzhotel/internal/domain/money"

	"github.com/google/uuid"
)

type Runner struct {
	HistoryHandler  *application.RoomHistoryHandler
	OccupiedHandler *application.OccupiedRoomsHandler

	CreateCustomer *application.CreateCustomerHandler
	CreditWallet   *application.CreditWalletHandler
	ListCustomers  *application.ListCustomersHandler
	BookRoom       *application.BookRoomHandler
	ListRooms      *application.ListRoomsHandler

	ConfirmRes  *application.ConfirmReservationHandler
	CancelRes   *application.CancelReservationHandler
	CheckoutRes *application.CheckoutRoomHandler
}

func NewRunner(
	history *application.RoomHistoryHandler,
	occupied *application.OccupiedRoomsHandler,
	createCust *application.CreateCustomerHandler,
	creditWallet *application.CreditWalletHandler,
	listCust *application.ListCustomersHandler,
	bookRoom *application.BookRoomHandler,
	listRooms *application.ListRoomsHandler,
	confirm *application.ConfirmReservationHandler,
	cancel *application.CancelReservationHandler,
	checkout *application.CheckoutRoomHandler,
) *Runner {
	return &Runner{
		HistoryHandler:  history,
		OccupiedHandler: occupied,
		CreateCustomer:  createCust,
		CreditWallet:    creditWallet,
		ListCustomers:   listCust,
		BookRoom:        bookRoom,
		ListRooms:       listRooms,
		ConfirmRes:      confirm,
		CancelRes:       cancel,
		CheckoutRes:     checkout,
	}
}

func (r *Runner) Start() {
	reader := bufio.NewScanner(os.Stdin)
	ctx := context.Background()

	fmt.Println("XYZ Hotel CLI")
	r.printHelp()

	for {
		fmt.Print("\n> ")
		if !reader.Scan() {
			break
		}
		line := reader.Text()
		parts := strings.Fields(line)

		if len(parts) == 0 {
			continue
		}

		command := parts[0]

		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("error: %v\n", r)
				}
			}()

			switch command {
			case "exit", "quit":
				os.Exit(0)
			case "help":
				r.printHelp()
			case "occupied":
				r.handleOccupied(ctx)
			case "history":
				if len(parts) < 2 {
					fmt.Println("usage: history <room_number>")
					return
				}
				r.handleHistory(ctx, parts[1])
			case "customers":
				r.handleListCustomers(ctx)
			case "create-customer":
				r.handleCreateCustomer(ctx, reader)
			case "credit":
				r.handleCreditWallet(ctx, reader)
			case "rooms":
				r.handleListRooms(ctx)
			case "book":
				r.handleBookRoom(ctx, reader)
			case "confirm":
				if len(parts) < 2 {
					fmt.Println("usage: confirm <reservation_uuid>")
					return
				}
				r.handleConfirm(ctx, parts[1])
			case "cancel":
				if len(parts) < 2 {
					fmt.Println("usage: cancel <reservation_uuid>")
					return
				}
				r.handleCancel(ctx, parts[1])
			case "checkout":
				if len(parts) < 2 {
					fmt.Println("usage: checkout <reservation_uuid>")
					return
				}
				r.handleCheckout(ctx, parts[1])
			default:
				fmt.Printf("unknown command: %s\n", command)
			}
		}()
	}
}

func (r *Runner) printHelp() {
	fmt.Println("commands:")
	fmt.Println("occupied, history <room>")
	fmt.Println("customers, create-customer, credit")
	fmt.Println("rooms, book")
	fmt.Println("confirm <id>, cancel <id>, checkout <id>")
	fmt.Println("exit")
}

func (r *Runner) handleListCustomers(ctx context.Context) {
	cmd := &application.ListCustomersCmd{}
	customers, err := r.ListCustomers.Handle(ctx, cmd)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("%-36s | %-20s | %-10s\n", "ID", "Name", "Balance")
	for _, c := range customers {
		fmt.Printf("%-36s | %-20s | %.2f\n", c.ID, c.FullName, float64(c.Wallet.Balance())/100.0)
	}
}

func (r *Runner) handleCreateCustomer(ctx context.Context, reader *bufio.Scanner) {
	fmt.Print("name: ")
	reader.Scan()
	name := reader.Text()

	fmt.Print("email: ")
	reader.Scan()
	email := reader.Text()

	fmt.Print("phone: ")
	reader.Scan()
	phone := reader.Text()

	cmd := &application.CreateCustomerCmd{
		FullName: name,
		Email:    email,
		Phone:    phone,
	}

	cust, err := r.CreateCustomer.Handle(ctx, cmd)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("customer created: %s\n", cust.ID)
}

func (r *Runner) handleCreditWallet(ctx context.Context, reader *bufio.Scanner) {
	fmt.Print("customer id: ")
	reader.Scan()
	idStr := reader.Text()
	uid, err := uuid.Parse(idStr)
	if err != nil {
		fmt.Println("invalid uuid")
		return
	}

	fmt.Print("amount: ")
	reader.Scan()
	amountStr := reader.Text()
	amountFloat, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		fmt.Println("invalid amount")
		return
	}

	fmt.Print("currency (EUR, USD...): ")
	reader.Scan()
	currency := strings.ToUpper(reader.Text())

	cents := int(amountFloat * 100)

	cmd := &application.CreditWalletCmd{
		CustomerID: uid,
		Money: money.Money{
			AmountCents: cents,
			Currency:    money.Currency(currency),
		},
	}

	err = r.CreditWallet.Handle(ctx, cmd)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Println("wallet credited")
}

func (r *Runner) handleListRooms(ctx context.Context) {
	infos, err := r.ListRooms.Handle(ctx)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	for _, i := range infos {
		fmt.Printf("%s : %.2f\n", i.Type, float64(i.PriceCents)/100.0)
	}
}

func (r *Runner) handleBookRoom(ctx context.Context, reader *bufio.Scanner) {
	fmt.Print("customer id: ")
	reader.Scan()
	custID, _ := uuid.Parse(reader.Text())

	fmt.Print("date (YYYY-MM-DD): ")
	reader.Scan()
	dateStr := reader.Text()
	checkIn, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		fmt.Println("invalid date")
		return
	}

	fmt.Print("nights: ")
	reader.Scan()
	nights, _ := strconv.Atoi(reader.Text())

	fmt.Print("rooms (space separated): ")
	reader.Scan()
	roomsStr := reader.Text()
	rooms := strings.Fields(roomsStr)

	cmd := &application.BookRoomCmd{
		CustomerID:     custID,
		CheckInDate:    checkIn,
		AmountOfNights: nights,
		RoomNumbers:    rooms,
	}

	resIDs, err := r.BookRoom.Handle(ctx, cmd)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Println("reservation created. ids:")
	for _, id := range resIDs {
		fmt.Println(id)
	}
}

func (r *Runner) handleConfirm(ctx context.Context, idStr string) {
	uid, _ := uuid.Parse(idStr)
	cmd := &application.ConfirmReservationCmd{ReservationID: uid}
	if err := r.ConfirmRes.Handle(ctx, cmd); err != nil {
		fmt.Printf("error: %v\n", err)
	} else {
		fmt.Println("confirmed")
	}
}

func (r *Runner) handleCancel(ctx context.Context, idStr string) {
	uid, _ := uuid.Parse(idStr)
	cmd := &application.CancelReservationCmd{ReservationID: uid}
	if err := r.CancelRes.Handle(ctx, cmd); err != nil {
		fmt.Printf("error: %v\n", err)
	} else {
		fmt.Println("cancelled")
	}
}

func (r *Runner) handleCheckout(ctx context.Context, idStr string) {
	uid, _ := uuid.Parse(idStr)
	cmd := &application.CheckoutRoomCmd{ReservationID: uid}
	if err := r.CheckoutRes.Handle(ctx, cmd); err != nil {
		fmt.Printf("error: %v\n", err)
	} else {
		fmt.Println("checked out")
	}
}

func (r *Runner) handleOccupied(ctx context.Context) {
	query := application.OccupiedRoomsQuery{}
	reservations, err := r.OccupiedHandler.Handle(ctx, query)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	if len(reservations) == 0 {
		fmt.Println("no occupied rooms")
		return
	}
	for _, res := range reservations {
		fmt.Printf("%s | %s\n", res.RoomNumber, res.CustomerID)
	}
}

func (r *Runner) handleHistory(ctx context.Context, roomNumber string) {
	query := application.RoomHistoryQuery{RoomNumber: roomNumber}
	reservations, err := r.HistoryHandler.Handle(ctx, query)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	for _, res := range reservations {
		fmt.Printf("%s | %s | %s\n", res.CheckInDate.Format("2006-01-02"), res.State, res.ID)
	}
}
