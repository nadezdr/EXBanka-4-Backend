package handlers

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	pb "github.com/RAF-SI-2025/EXBanka-4-Backend/shared/pb/securities"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SecuritiesServer struct {
	pb.UnimplementedSecuritiesServiceServer
	DB *sql.DB
}

func (s *SecuritiesServer) Ping(_ context.Context, _ *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{Message: "securities-service ok"}, nil
}

func (s *SecuritiesServer) GetStockExchanges(ctx context.Context, _ *pb.GetStockExchangesRequest) (*pb.GetStockExchangesResponse, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, name, acronym, mic_code, polity, currency, timezone
		FROM stock_exchanges
		ORDER BY name`)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "query failed: %v", err)
	}
	defer rows.Close()

	var exchanges []*pb.StockExchange
	for rows.Next() {
		e := &pb.StockExchange{}
		if err := rows.Scan(&e.Id, &e.Name, &e.Acronym, &e.MicCode, &e.Polity, &e.Currency, &e.Timezone); err != nil {
			return nil, status.Errorf(codes.Internal, "scan failed: %v", err)
		}
		exchanges = append(exchanges, e)
	}
	return &pb.GetStockExchangesResponse{Exchanges: exchanges}, nil
}

func (s *SecuritiesServer) GetStockExchangeByMIC(ctx context.Context, req *pb.GetStockExchangeByMICRequest) (*pb.GetStockExchangeByMICResponse, error) {
	e := &pb.StockExchange{}
	err := s.DB.QueryRowContext(ctx, `
		SELECT id, name, acronym, mic_code, polity, currency, timezone
		FROM stock_exchanges
		WHERE mic_code = $1`, req.MicCode).
		Scan(&e.Id, &e.Name, &e.Acronym, &e.MicCode, &e.Polity, &e.Currency, &e.Timezone)
	if err == sql.ErrNoRows {
		return nil, status.Errorf(codes.NotFound, "exchange with MIC %q not found", req.MicCode)
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "query failed: %v", err)
	}
	return &pb.GetStockExchangeByMICResponse{Exchange: e}, nil
}

func (s *SecuritiesServer) CreateStockExchange(ctx context.Context, req *pb.CreateStockExchangeRequest) (*pb.CreateStockExchangeResponse, error) {
	e := &pb.StockExchange{}
	err := s.DB.QueryRowContext(ctx, `
		INSERT INTO stock_exchanges (name, acronym, mic_code, polity, currency, timezone)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, acronym, mic_code, polity, currency, timezone`,
		req.Name, req.Acronym, req.MicCode, req.Polity, req.Currency, req.Timezone).
		Scan(&e.Id, &e.Name, &e.Acronym, &e.MicCode, &e.Polity, &e.Currency, &e.Timezone)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, status.Errorf(codes.AlreadyExists, "exchange with MIC %q already exists", req.MicCode)
		}
		return nil, status.Errorf(codes.Internal, "insert failed: %v", err)
	}
	return &pb.CreateStockExchangeResponse{Exchange: e}, nil
}

func (s *SecuritiesServer) UpdateStockExchange(ctx context.Context, req *pb.UpdateStockExchangeRequest) (*pb.UpdateStockExchangeResponse, error) {
	e := &pb.StockExchange{}
	err := s.DB.QueryRowContext(ctx, `
		UPDATE stock_exchanges
		SET name=$1, acronym=$2, polity=$3, currency=$4, timezone=$5
		WHERE mic_code=$6
		RETURNING id, name, acronym, mic_code, polity, currency, timezone`,
		req.Name, req.Acronym, req.Polity, req.Currency, req.Timezone, req.MicCode).
		Scan(&e.Id, &e.Name, &e.Acronym, &e.MicCode, &e.Polity, &e.Currency, &e.Timezone)
	if err == sql.ErrNoRows {
		return nil, status.Errorf(codes.NotFound, "exchange with MIC %q not found", req.MicCode)
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update failed: %v", err)
	}
	return &pb.UpdateStockExchangeResponse{Exchange: e}, nil
}

func (s *SecuritiesServer) DeleteStockExchange(ctx context.Context, req *pb.DeleteStockExchangeRequest) (*pb.DeleteStockExchangeResponse, error) {
	res, err := s.DB.ExecContext(ctx, `DELETE FROM stock_exchanges WHERE mic_code = $1`, req.MicCode)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete failed: %v", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return nil, status.Errorf(codes.NotFound, "exchange with MIC %q not found", req.MicCode)
	}
	return &pb.DeleteStockExchangeResponse{}, nil
}
