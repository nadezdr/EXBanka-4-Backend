package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	pb "github.com/RAF-SI-2025/EXBanka-4-Backend/shared/pb/account"
	"github.com/RAF-SI-2025/EXBanka-4-Backend/services/api-gateway/middleware"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateAccountRequest is the request body for creating a bank account.
type CreateAccountRequest struct {
	ClientID       int64       `json:"clientId"       binding:"required"`
	AccountType    string      `json:"accountType"    binding:"required"`
	CurrencyCode   string      `json:"currencyCode"   binding:"required"`
	InitialBalance float64     `json:"initialBalance"`
	AccountName    string      `json:"accountName"`
	CreateCard     bool        `json:"createCard"`
	CompanyData    *companyReq `json:"companyData"`
}

type companyReq struct {
	Name               string `json:"name"`
	RegistrationNumber string `json:"registrationNumber"`
	PIB                string `json:"pib"`
	ActivityCode       string `json:"activityCode"`
	Address            string `json:"address"`
}

// CreateAccount godoc
// @Summary      Create bank account
// @Description  Creates a new bank account for a client. Requires employee authentication.
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        body  body      CreateAccountRequest  true  "Account creation data"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/accounts/create [post]
func CreateAccount(accountClient pb.AccountServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateAccountRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		employeeID, err := middleware.GetUserIDFromToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "could not extract employee identity from token"})
			return
		}

		grpcReq := &pb.CreateAccountRequest{
			ClientId:       req.ClientID,
			AccountType:    req.AccountType,
			CurrencyCode:   req.CurrencyCode,
			InitialBalance: req.InitialBalance,
			AccountName:    req.AccountName,
			CreateCard:     req.CreateCard,
			EmployeeId:     employeeID,
		}
		if req.CompanyData != nil {
			grpcReq.CompanyData = &pb.CompanyData{
				Name:               req.CompanyData.Name,
				RegistrationNumber: req.CompanyData.RegistrationNumber,
				Pib:                req.CompanyData.PIB,
				ActivityCode:       req.CompanyData.ActivityCode,
				Address:            req.CompanyData.Address,
			}
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		resp, err := accountClient.CreateAccount(ctx, grpcReq)
		if err != nil {
			switch status.Code(err) {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		a := resp.Account
		c.JSON(http.StatusCreated, gin.H{
			"id":                a.Id,
			"accountNumber":     a.AccountNumber,
			"accountName":       a.AccountName,
			"ownerId":           a.OwnerId,
			"employeeId":        a.EmployeeId,
			"currencyCode":      a.CurrencyCode,
			"accountType":       a.AccountType,
			"status":            a.Status,
			"balance":           a.Balance,
			"availableBalance":  a.AvailableBalance,
			"createdDate":       a.CreatedDate,
		})
	}
}
