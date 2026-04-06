package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	pb "github.com/RAF-SI-2025/EXBanka-4-Backend/shared/pb/securities"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type exchangeJSON struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Acronym  string `json:"acronym"`
	MICCode  string `json:"micCode"`
	Polity   string `json:"polity"`
	Currency string `json:"currency"`
	Timezone string `json:"timezone"`
}

func toExchangeJSON(e *pb.StockExchange) exchangeJSON {
	return exchangeJSON{
		ID:       e.Id,
		Name:     e.Name,
		Acronym:  e.Acronym,
		MICCode:  e.MicCode,
		Polity:   e.Polity,
		Currency: e.Currency,
		Timezone: e.Timezone,
	}
}

// GetStockExchanges godoc
// @Summary      List all stock exchanges
// @Tags         securities
// @Produce      json
// @Success      200  {array}   exchangeJSON
// @Router       /stock-exchanges [get]
func GetStockExchanges(client pb.SecuritiesServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		resp, err := client.GetStockExchanges(ctx, &pb.GetStockExchangesRequest{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch stock exchanges"})
			return
		}

		result := make([]exchangeJSON, 0, len(resp.Exchanges))
		for _, e := range resp.Exchanges {
			result = append(result, toExchangeJSON(e))
		}
		c.JSON(http.StatusOK, result)
	}
}

// GetStockExchangeByMIC godoc
// @Summary      Get stock exchange by MIC code
// @Tags         securities
// @Produce      json
// @Param        mic  path  string  true  "MIC code"
// @Success      200  {object}  exchangeJSON
// @Router       /stock-exchanges/{mic} [get]
func GetStockExchangeByMIC(client pb.SecuritiesServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		mic := c.Param("mic")
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		resp, err := client.GetStockExchangeByMIC(ctx, &pb.GetStockExchangeByMICRequest{MicCode: mic})
		if err != nil {
			if status.Code(err) == codes.NotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch exchange"})
			return
		}
		c.JSON(http.StatusOK, toExchangeJSON(resp.Exchange))
	}
}

// CreateStockExchange godoc
// @Summary      Create a new stock exchange
// @Tags         securities
// @Accept       json
// @Produce      json
// @Success      201  {object}  exchangeJSON
// @Router       /stock-exchanges [post]
func CreateStockExchange(client pb.SecuritiesServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			Name     string `json:"name"     binding:"required"`
			Acronym  string `json:"acronym"  binding:"required"`
			MICCode  string `json:"micCode"  binding:"required"`
			Polity   string `json:"polity"   binding:"required"`
			Currency string `json:"currency" binding:"required"`
			Timezone string `json:"timezone" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		resp, err := client.CreateStockExchange(ctx, &pb.CreateStockExchangeRequest{
			Name:     body.Name,
			Acronym:  body.Acronym,
			MicCode:  body.MICCode,
			Polity:   body.Polity,
			Currency: body.Currency,
			Timezone: body.Timezone,
		})
		if err != nil {
			if status.Code(err) == codes.AlreadyExists {
				c.JSON(http.StatusConflict, gin.H{"error": status.Convert(err).Message()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create exchange"})
			return
		}
		c.JSON(http.StatusCreated, toExchangeJSON(resp.Exchange))
	}
}

// UpdateStockExchange godoc
// @Summary      Update a stock exchange by MIC code
// @Tags         securities
// @Accept       json
// @Produce      json
// @Param        mic  path  string  true  "MIC code"
// @Success      200  {object}  exchangeJSON
// @Router       /stock-exchanges/{mic} [put]
func UpdateStockExchange(client pb.SecuritiesServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		mic := c.Param("mic")
		var body struct {
			Name     string `json:"name"     binding:"required"`
			Acronym  string `json:"acronym"  binding:"required"`
			Polity   string `json:"polity"   binding:"required"`
			Currency string `json:"currency" binding:"required"`
			Timezone string `json:"timezone" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		resp, err := client.UpdateStockExchange(ctx, &pb.UpdateStockExchangeRequest{
			MicCode:  mic,
			Name:     body.Name,
			Acronym:  body.Acronym,
			Polity:   body.Polity,
			Currency: body.Currency,
			Timezone: body.Timezone,
		})
		if err != nil {
			if status.Code(err) == codes.NotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update exchange"})
			return
		}
		c.JSON(http.StatusOK, toExchangeJSON(resp.Exchange))
	}
}

// DeleteStockExchange godoc
// @Summary      Delete a stock exchange by MIC code
// @Tags         securities
// @Param        mic  path  string  true  "MIC code"
// @Success      204
// @Router       /stock-exchanges/{mic} [delete]
func DeleteStockExchange(client pb.SecuritiesServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		mic := c.Param("mic")
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		_, err := client.DeleteStockExchange(ctx, &pb.DeleteStockExchangeRequest{MicCode: mic})
		if err != nil {
			if status.Code(err) == codes.NotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete exchange"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
