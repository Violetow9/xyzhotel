package main

import (
	_ "fmt"
	"log"
	"xyzhotel/application"
	"xyzhotel/domain/customer"
	infraCustomer "xyzhotel/infrastructure/customer"
	"xyzhotel/infrastructure/mysql"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := mysql.OpenFromEnv()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	customerRepository := infraCustomer.NewCustomerRepository(db)
	customerService := &customer.Service{
		Repository: customerRepository,
	}
	creditWalletHandler := &application.CreditWalletHandler{
		CustomerService: customerService,
	}
	customerController := &infraCustomer.Controller{
		CustomerService: customerService,
		CreditWallet:    creditWalletHandler,
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
