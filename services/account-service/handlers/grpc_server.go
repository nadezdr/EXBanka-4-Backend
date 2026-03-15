package handlers

import (
	"context"
	"database/sql"
	"time"

	pb "github.com/RAF-SI-2025/EXBanka-4-Backend/shared/pb/account"
	"github.com/RAF-SI-2025/EXBanka-4-Backend/services/account-service/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccountServer struct {
	pb.UnimplementedAccountServiceServer
	DB         *sql.DB
	ClientDB   *sql.DB
	ExchangeDB *sql.DB
}

// accountTypeCode maps account type string to 2-digit code used in account number generation.
func accountTypeCode(accountType string) string {
	switch accountType {
	case "CURRENT":
		return "01"
	case "SAVINGS":
		return "02"
	case "FOREIGN_CURRENCY":
		return "03"
	case "BUSINESS":
		return "04"
	default:
		return "00"
	}
}

func (s *AccountServer) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	// 1. Validate client exists
	var clientID int64
	err := s.ClientDB.QueryRowContext(ctx,
		`SELECT id FROM clients WHERE id = $1`, req.ClientId).Scan(&clientID)
	if err == sql.ErrNoRows {
		return nil, status.Errorf(codes.NotFound, "client with id %d not found", req.ClientId)
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to verify client: %v", err)
	}

	// 2. Validate currency exists and get its id
	var currencyID int64
	var currencyCode string
	err = s.ExchangeDB.QueryRowContext(ctx,
		`SELECT id, code FROM currencies WHERE code = $1`, req.CurrencyCode).Scan(&currencyID, &currencyCode)
	if err == sql.ErrNoRows {
		return nil, status.Errorf(codes.NotFound, "currency with code %q not found", req.CurrencyCode)
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to verify currency: %v", err)
	}

	// 3. Generate account number
	accountNumber := utils.GenerateAccountNumber("265", "0001", accountTypeCode(req.AccountType))

	// 4. Set expiration date 5 years from now
	expirationDate := time.Now().AddDate(5, 0, 0).Format("2006-01-02")

	// 5. Insert account
	var accountID int64
	var createdDate string
	err = s.DB.QueryRowContext(ctx,
		`INSERT INTO accounts
			(account_number, account_name, owner_id, employee_id, currency_id,
			 account_type, status, balance, available_balance, expiration_date)
		VALUES ($1, $2, $3, $4, $5, $6, 'ACTIVE', $7, $7, $8)
		RETURNING id, created_date`,
		accountNumber,
		req.AccountName,
		req.ClientId,
		req.EmployeeId,
		currencyID,
		req.AccountType,
		req.InitialBalance,
		expirationDate,
	).Scan(&accountID, &createdDate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create account: %v", err)
	}

	// 6. Insert company if provided
	if req.CompanyData != nil && req.CompanyData.Name != "" {
		_, err = s.DB.ExecContext(ctx,
			`INSERT INTO companies
				(name, registration_number, pib, activity_code, address, owner_client_id)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			req.CompanyData.Name,
			req.CompanyData.RegistrationNumber,
			req.CompanyData.Pib,
			req.CompanyData.ActivityCode,
			req.CompanyData.Address,
			req.ClientId,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "account created but failed to create company: %v", err)
		}
	}

	return &pb.CreateAccountResponse{
		Account: &pb.AccountResponse{
			Id:               accountID,
			AccountNumber:    accountNumber,
			AccountName:      req.AccountName,
			OwnerId:          req.ClientId,
			EmployeeId:       req.EmployeeId,
			CurrencyCode:     currencyCode,
			AccountType:      req.AccountType,
			Status:           "ACTIVE",
			Balance:          req.InitialBalance,
			AvailableBalance: req.InitialBalance,
			CreatedDate:      createdDate,
		},
	}, nil
}
