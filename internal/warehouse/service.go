package warehouse

import (
	"context"
	"errors"

	"repository_class/internal/domain"
)

// Errors
var (
	ErrNotFound            = errors.New("warehouse not found")
	ErrWarehouseRegistered = errors.New("warehouse number is already registered")
	ErrInvalidStruct       = errors.New("invalid input structure for section")
	ErrInvalidId           = errors.New("invalid id")
)

type Service interface {
	//read
	GetAll(ctx context.Context) ([]domain.Warehouse, error)
	Get(ctx context.Context, id int) (domain.Warehouse, error)
	Create(ctx context.Context, w domain.Warehouse) (domain.Warehouse, error)
	Update(ctx context.Context, w domain.Warehouse, id int) (domain.Warehouse, error)
	Delete(ctx context.Context, id int) error
	ReportProducts(ctx context.Context, id int) (domain.WarehouseReport, error)
}

type service struct {
	repo Repository
}

func (s *service) ReportProducts(ctx context.Context, id int) (domain.WarehouseReport, error) {
	warehouse, err := s.repo.ReportProducts(ctx, id)
	if err != nil {
		return domain.WarehouseReport{}, err
	}
	if warehouse == (domain.WarehouseReport{}) {
		warehouse = domain.WarehouseReport{}
	}
	return warehouse, nil
}

func NewService(repo *Repository) Service {
	return &service{repo: *repo}
}

/*
	func validateUpdateFields(w domain.Warehouse,wu domain.WarehouseUpdate) domain.Warehouse{
		if wu.ID!=0{}
		if wu.Address!=""{

		}
		if wu.Telephone!=""{}
		if wu.WarehouseCode!=""{}
		if wu.MinimumCapacity!=0{}
		if wu.MinimumTemperature!=0{}
		if wu.ID!=
		return
	}
*/
func (s *service) GetAll(ctx context.Context) ([]domain.Warehouse, error) {
	warehouse, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	if warehouse == nil {
		warehouse = []domain.Warehouse{}
	}
	return warehouse, nil
}

func (s *service) Get(ctx context.Context, id int) (domain.Warehouse, error) {
	warehouse, err := s.repo.Get(ctx, id)
	if err != nil {
		return domain.Warehouse{}, err
	}
	if warehouse == (domain.Warehouse{}) {
		warehouse = domain.Warehouse{}
	}
	return warehouse, nil
}

func (s *service) Create(ctx context.Context, w domain.Warehouse) (domain.Warehouse, error) {
	// exist := s.repo.Exists(ctx, w.WarehouseCode)
	// if exist {
	// 	return domain.Warehouse{}, ErrWarehouseRegistered
	// }
	id, err := s.repo.Save(ctx, w)

	if err != nil {
		return domain.Warehouse{}, err
	}
	warehouse, err := s.repo.Get(ctx, id)
	if err != nil {
		return domain.Warehouse{}, err
	}
	return warehouse, nil
}

func (s *service) Update(ctx context.Context, w domain.Warehouse, id int) (domain.Warehouse, error) {
	_, err := s.repo.Get(ctx, id)
	if err != nil {
		return domain.Warehouse{}, err
	}
	// if s.repo.Exists(ctx, w.WarehouseCode) {
	// 	return domain.Warehouse{}, ErrWarehouseRegistered
	// }
	errs := s.repo.Update(ctx, w)
	if errs != nil {
		return domain.Warehouse{}, err
	}
	return w, nil
}

func (s *service) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
