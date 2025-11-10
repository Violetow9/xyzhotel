package main

import (
	_ "fmt"
	"log"
	"xyzhotel/internal/application"
	"xyzhotel/internal/domain/customer"
	customer2 "xyzhotel/internal/infrastructure/customer"
	"xyzhotel/internal/infrastructure/mysql"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := mysql.OpenFromEnv()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	customerRepository := customer2.NewCustomerRepository(db)
	customerService := &customer.Service{
		Repository: customerRepository,
	}
	createCustomerHandler := &application.CreateCustomerHandler{
		CustomerRepository: customerRepository,
	}
	listCustomersHandler := &application.ListCustomersHandler{
		CustomerRepository: customerRepository,
	}
	creditWalletHandler := &application.CreditWalletHandler{
		CustomerService: customerService,
	}
	customerController := &customer2.Controller{
		CreateCustomer: createCustomerHandler,
		ListCustomers:  listCustomersHandler,
		CreditWallet:   creditWalletHandler,
	}

	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	customersGroup := r.Group("/customers")
	customersGroup.POST("", customerController.CreateCustomerHandler)
	customersGroup.GET("", customerController.ListCustomersHandler)
	walletGroup := customersGroup.Group("/wallet")
	walletGroup.POST("/credit", customerController.CreditWalletHandler)

	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
