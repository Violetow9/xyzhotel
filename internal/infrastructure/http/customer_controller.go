package http

import (
	"math"
	"net/http"

	"xyzhotel/internal/application"
	"xyzhotel/internal/domain/money"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createCustomerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type creditWalletRequest struct {
	ID       string  `json:"id"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type CustomerController struct {
	CreateCustomer *application.CreateCustomerHandler
	ListCustomers  *application.ListCustomersHandler
	CreditWallet   *application.CreditWalletHandler
}

func (c *CustomerController) CreateCustomerHandler(ctx *gin.Context) {
	var in createCustomerRequest
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createCmd := &application.CreateCustomerCmd{
		FullName: in.Name,
		Email:    in.Email,
		Phone:    in.Phone,
	}

	newCustomer, err := c.CreateCustomer.Handle(ctx, createCmd)
	if err != nil {
		if err.Error() == "customer already exists" {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id": newCustomer.ID,
	})
}

func (c *CustomerController) CreditWalletHandler(ctx *gin.Context) {
	var in creditWalletRequest
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(in.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer ID"})
		return
	}

	cents := int(math.Round(in.Amount * 100))

	amount := money.Money{
		AmountCents: cents,
		Currency:    money.Currency(in.Currency),
	}

	creditWalletCmd := &application.CreditWalletCmd{
		CustomerID: id,
		Money:      amount,
	}

	err = c.CreditWallet.Handle(ctx, creditWalletCmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "wallet credited"})
}

func (c *CustomerController) ListCustomersHandler(ctx *gin.Context) {
	customers, err := c.ListCustomers.Handle(ctx, &application.ListCustomersCmd{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, customers)
}
