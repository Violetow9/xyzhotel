package customer

import (
	"net/http"
	application2 "xyzhotel/internal/application"
	money2 "xyzhotel/internal/domain/money"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type createCustomer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type creditWallet struct {
	ID       string  `json:"id"`
	Amount   float32 `json:"amount"`
	Currency string  `json:"currency"`
}

type Controller struct {
	CreateCustomer *application2.CreateCustomerHandler
	ListCustomers  *application2.ListCustomersHandler
	CreditWallet   *application2.CreditWalletHandler
}

func (c *Controller) CreateCustomerHandler(ctx *gin.Context) {
	var in createCustomer
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	createCmd := &application2.CreateCustomerCmd{
		FullName: in.Name,
		Email:    in.Email,
		Phone:    in.Phone,
	}

	newCustomer, err := c.CreateCustomer.Handle(ctx, createCmd)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{
		"id": newCustomer.ID,
	})
}

func (c *Controller) CreditWalletHandler(ctx *gin.Context) {
	var in creditWallet
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(in.ID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid customer ID"})
		return
	}

	currency, err := money2.CurrencyFromString(in.Currency)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid currency"})
		return
	}

	amount := money2.Money{
		Amount:   decimal.NewFromFloat32(in.Amount),
		Currency: currency,
	}

	creditWalletCmd := &application2.CreditWalletCmd{
		CustomerID: id,
		Money:      amount,
	}

	err = c.CreditWallet.Handle(ctx, creditWalletCmd)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{"status": "wallet credited"})
}

func (c *Controller) ListCustomersHandler(ctx *gin.Context) {
	customers, err := c.ListCustomers.Handle(ctx, &application2.ListCustomersCmd{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, customers)
}
