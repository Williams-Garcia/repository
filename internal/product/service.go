package product

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"repository_class/internal/domain"

	"github.com/go-playground/validator/v10"
)

// Errors
var (
	ErrNotFound          = errors.New("product not found")
	ErrUniqueProduct     = errors.New("product code must be unique")
	ErrProductRegistered = errors.New("section number is already registered")
	ErrInvalidStruct     = errors.New("invalid input structure for section")
)

type Service interface {
	GetAll(ctx context.Context) ([]domain.Product, error)
	Get(ctx context.Context, id int) (domain.Product, error)
	Delete(ctx context.Context, id int) error
	Create(ctx context.Context, prod domain.Product) (domain.Product, error)
	Update(ctx context.Context, prod domain.Product, id int) (domain.Product, error)
}

type service struct {
	repo Repository
}

func validateUpdateFields(productDB domain.Product, productUpdate domain.Product) domain.Product {
	defaultTime := time.Time{}

	if productUpdate.Name != "" {
		productDB.Name = productUpdate.Name
	}
	if productUpdate.Quantity != 0 {
		productDB.Quantity = productUpdate.Quantity
	}
	if productUpdate.CodeValue != "" {
		productDB.CodeValue = productUpdate.CodeValue
	}
	if productUpdate.Expiration != defaultTime {
		productDB.Expiration = productUpdate.Expiration
	}
	if productUpdate.Price != 0 {
		productDB.Price = productUpdate.Price
	}

	return productDB
}

func (s *service) Update(ctx context.Context, prod domain.Product, id int) (domain.Product, error) {
	product, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Product{}, ErrNotFound
		}
		return domain.Product{}, err
	}
	if s.repo.Exists(ctx, prod.CodeValue) {
		return domain.Product{}, ErrProductRegistered
	}
	prod = validateUpdateFields(product, prod)
	err = s.repo.Update(ctx, prod)
	if err != nil {
		return domain.Product{}, err
	}
	return prod, nil
}

func (s *service) Create(ctx context.Context, prod domain.Product) (domain.Product, error) {
	validator := validator.New()
	if err := validator.Struct(&prod); err != nil {
		return domain.Product{}, ErrInvalidStruct
	}
	// TO DO:	Validate if product type exists

	// Method Exists return a true if prod exists in db
	if s.repo.Exists(ctx, prod.CodeValue) {
		return domain.Product{}, ErrUniqueProduct
	}
	idProd, err := s.repo.Save(ctx, prod)
	if err != nil {
		return domain.Product{}, err
	}
	prod, _ = s.repo.Get(ctx, idProd)
	return prod, nil
}

func (s *service) GetAll(ctx context.Context) ([]domain.Product, error) {
	products, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	if products == nil {
		return []domain.Product{}, nil
	}
	return products, nil
}

func (s *service) Get(ctx context.Context, id int) (domain.Product, error) {
	product, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Product{}, ErrNotFound
		}
		return domain.Product{}, err
	}
	return product, nil
}

func (s *service) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func NewService(repo *Repository) Service {
	return &service{repo: *repo}
}
