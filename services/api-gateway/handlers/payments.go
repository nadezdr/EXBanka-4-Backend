package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	pb "github.com/RAF-SI-2025/EXBanka-4-Backend/shared/pb/payment"
	"github.com/RAF-SI-2025/EXBanka-4-Backend/services/api-gateway/middleware"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreatePaymentRequest struct {
	RecipientName    string  `json:"recipientName"    binding:"required"`
	RecipientAccount string  `json:"recipientAccount" binding:"required"`
	Amount           float64 `json:"amount"           binding:"required,gt=0"`
	PaymentCode      string  `json:"paymentCode"`
	ReferenceNumber  string  `json:"referenceNumber"`
	Purpose          string  `json:"purpose"`
	FromAccount      string  `json:"fromAccount"      binding:"required"`
}

// CreatePayment godoc
// @Summary      Create a new payment
// @Description  Initiates a payment from client's account to a recipient account.
// @Tags         payments
// @Accept       json
// @Produce      json
// @Param        body  body      CreatePaymentRequest  true  "Payment data"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]string
// @Failure      403   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/payments/create [post]
func CreatePayment(paymentClient pb.PaymentServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreatePaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		clientID, err := middleware.GetUserIDFromToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "could not extract identity from token"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		resp, err := paymentClient.CreatePayment(ctx, &pb.CreatePaymentRequest{
			ClientId:        clientID,
			FromAccount:     req.FromAccount,
			RecipientName:   req.RecipientName,
			RecipientAccount: req.RecipientAccount,
			Amount:          req.Amount,
			PaymentCode:     req.PaymentCode,
			ReferenceNumber: req.ReferenceNumber,
			Purpose:         req.Purpose,
		})
		if err != nil {
			switch status.Code(err) {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
			case codes.PermissionDenied:
				c.JSON(http.StatusForbidden, gin.H{"error": status.Convert(err).Message()})
			case codes.FailedPrecondition:
				c.JSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":            resp.Id,
			"orderNumber":   resp.OrderNumber,
			"fromAccount":   resp.FromAccount,
			"toAccount":     resp.ToAccount,
			"initialAmount": resp.InitialAmount,
			"finalAmount":   resp.FinalAmount,
			"fee":           resp.Fee,
			"status":        resp.Status,
			"timestamp":     resp.Timestamp,
		})
	}
}
