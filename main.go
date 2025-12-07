package main

import (
	"log"
	"os"
	"xyzhotel/internal/application"
	"xyzhotel/internal/domain/money"
	"xyzhotel/internal/infrastructure/cli"
	httpInfra "xyzhotel/internal/infrastructure/http"
	"xyzhotel/internal/infrastructure/mysql"

	"github.com/gin-gonic/gin"
)

func main() {
	if os.Getenv("DB_DSN") == "" {
		os.Setenv("DB_DSN", "root:root@tcp(127.0.0.1:3306)/xyzhotel?parseTime=true")
	}

	db, err := mysql.OpenFromEnv()
	if err != nil {
		log.Fatalf("failed to open database connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("database is not reachable: %v", err)
	}

	customerRepo := mysql.NewCustomerRepository(db)
	roomRepo := mysql.NewRoomRepository(db)
	reservationRepo := mysql.NewReservationRepository(db)

	currencyConverter := money.NewConverter()

	debitWalletHandler := &application.DebitWalletHandler{
		CustomerRepo: customerRepo,
		Converter:    currencyConverter,
	}
	creditWalletHandler := &application.CreditWalletHandler{
		CustomerRepo: customerRepo,
		Converter:    currencyConverter,
	}

	createCustomerHandler := &application.CreateCustomerHandler{
		CustomerRepository: customerRepo,
	}
	listCustomersHandler := &application.ListCustomersHandler{
		CustomerRepository: customerRepo,
	}

	listRoomsHandler := &application.ListRoomsHandler{}

	bookRoomHandler := &application.BookRoomHandler{
		CustomerRepository: customerRepo,
		ReservationRepo:    reservationRepo,
		RoomRepository:     roomRepo,
		DebitWallet:        debitWalletHandler,
	}

	confirmResHandler := &application.ConfirmReservationHandler{
		ReservationRepo: reservationRepo,
		DebitWallet:     debitWalletHandler,
	}
	cancelResHandler := &application.CancelReservationHandler{
		ReservationRepo: reservationRepo,
	}
	checkoutResHandler := &application.CheckoutRoomHandler{
		ReservationRepo: reservationRepo,
	}
	historyHandler := &application.RoomHistoryHandler{
		ReservationRepo: reservationRepo,
	}
	occupiedHandler := &application.OccupiedRoomsHandler{
		ReservationRepo: reservationRepo,
	}

	if len(os.Args) > 1 && os.Args[1] == "cli" {
		log.Println("Starting in CLI Mode...")

		runner := cli.NewRunner(
			historyHandler,
			occupiedHandler,
			createCustomerHandler,
			creditWalletHandler,
			listCustomersHandler,
			bookRoomHandler,
			listRoomsHandler,
			confirmResHandler,
			cancelResHandler,
			checkoutResHandler,
		)

		runner.Start()
		return
	}

	customerController := &httpInfra.CustomerController{
		CreateCustomer: createCustomerHandler,
		ListCustomers:  listCustomersHandler,
		CreditWallet:   creditWalletHandler,
	}

	roomController := &httpInfra.RoomController{
		ListRooms: listRoomsHandler,
	}

	reservationController := &httpInfra.ReservationController{
		BookRoom:           bookRoomHandler,
		ConfirmReservation: confirmResHandler,
		CancelReservation:  cancelResHandler,
		CheckoutRoom:       checkoutResHandler,
	}

	adminController := &httpInfra.AdminController{
		RoomHistory:   historyHandler,
		OccupiedRooms: occupiedHandler,
	}

	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	customersGroup := r.Group("/customers")
	{
		customersGroup.POST("", customerController.CreateCustomerHandler)
		customersGroup.GET("", customerController.ListCustomersHandler)

		customersGroup.POST("/wallet/credit", customerController.CreditWalletHandler)
	}

	roomsGroup := r.Group("/rooms")
	{
		roomsGroup.GET("", roomController.ListRoomsHandler)
	}

	reservationGroup := r.Group("/reservations")
	{
		reservationGroup.POST("", reservationController.MakeReservationHandler)
		reservationGroup.POST("/:id/confirm", reservationController.ConfirmReservationHandler)
		reservationGroup.POST("/:id/cancel", reservationController.CancelReservationHandler)
		reservationGroup.POST("/:id/checkout", reservationController.CheckoutRoomHandler)
	}

	adminGroup := r.Group("/admin")
	{
		adminGroup.GET("/rooms/occupied", adminController.OccupiedHandler)
		adminGroup.GET("/rooms/:roomNumber/history", adminController.HistoryHandler)
	}

	log.Println("Server starting on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
