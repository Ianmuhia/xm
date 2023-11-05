package main

import (
	"context"
	"xmserver/db"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Company represents a company entity.
type Company struct {
	ID          string
	Name        string
	Description string
	Employees   int
	Registered  bool
	Type        string
}

// CompanyService is an interface for CRUD operations on the Company entity.
type CompanyService interface {
	// CreateCompany creates a new company and returns the created company or an error.
	CreateCompany(ctx context.Context, company *Company) (*Company, error)

	// GetCompany retrieves a company by its ID and returns it or an error.
	GetCompany(ctx context.Context, companyID string) (*Company, error)

	// UpdateCompany updates an existing company and returns the updated company or an error.
	UpdateCompany(ctx context.Context, company *Company) (*Company, error)

	// DeleteCompany deletes a company by its ID and returns an error if the deletion fails.
	DeleteCompany(ctx context.Context, companyID string) error

	ListCompanies(ctx context.Context) ([]*Company, error)
}

type companyService struct {
	db *db.Queries
}

// NewUsersService creates a new UsersService instance using the provided bun.DB database connection.
// It returns the newly created UsersService.
func NewCompanyService(q db.DBTX) CompanyService {
	return &companyService{db: db.New(q)}
}

func (s *companyService) CreateCompany(ctx context.Context, company *Company) (*Company, error) {
	data, err := s.db.CreateCompany(ctx, db.CreateCompanyParams{
		Name:        company.Name,
		Description: pgtype.Text{String: company.Description, Valid: true},
		Employees:   int32(company.Employees),
		Registered:  true,
		Type:        db.CompanyTypeCooperative,
	})
	if err != nil {
		return nil, err
	}

	return &Company{
		ID:          data.ID.String(),
		Name:        data.Name,
		Description: data.Description.String,
		Employees:   int(data.Employees),
		Registered:  data.Registered,
		Type:        string(data.Type),
	}, nil
}

func (s *companyService) GetCompany(ctx context.Context, companyID string) (*Company, error) {
	data, err := s.db.GetCompany(ctx, uuid.FromStringOrNil(companyID))
	if err != nil {
		return nil, err
	}

	return &Company{
		ID:          data.ID.String(),
		Name:        data.Name,
		Description: data.Description.String,
		Employees:   int(data.Employees),
		Registered:  data.Registered,
		Type:        string(data.Type),
	}, nil
}

func (s *companyService) UpdateCompany(ctx context.Context, company *Company) (*Company, error) {
	data, err := s.db.UpdateCompany(ctx, db.UpdateCompanyParams{
		Name:        company.Name,
		Description: pgtype.Text{String: company.Description, Valid: true},
		ID:          uuid.FromStringOrNil(company.ID),
		Employees:   int32(company.Employees),
		Type:        db.CompanyTypeCooperative,
	})
	if err != nil {
		return nil, err
	}

	return &Company{
		ID:          data.ID.String(),
		Name:        data.Name,
		Description: data.Description.String,
		Employees:   int(data.Employees),
		Registered:  data.Registered,
		Type:        string(data.Type),
	}, nil
}

func (s *companyService) DeleteCompany(ctx context.Context, companyID string) error {
	if err := s.db.DeleteCompany(ctx, uuid.FromStringOrNil(companyID)); err != nil {
		return err
	}
	return nil
}

func (s *companyService) ListCompanies(ctx context.Context) ([]*Company, error) {
	data, err := s.db.ListCompanies(ctx)
	if err != nil {
		return nil, err
	}
	companies := make([]*Company, len(data))

	for i, v := range data {
		companies[i] = &Company{
			ID:          v.ID.String(),
			Name:        v.Name,
			Description: v.Description.String,
			Employees:   int(v.Employees),
			Registered:  v.Registered,
			Type:        string(v.Type),
		}
	}

	return companies, nil
}
