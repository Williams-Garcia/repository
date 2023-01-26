package warehouse

import (
	"context"
	"database/sql"
	"log"

	"repository_class/internal/domain"
)

// Repository encapsulates the storage of a warehouse.
type Repository interface {
	GetAll(ctx context.Context) ([]domain.Warehouse, error)
	Get(ctx context.Context, id int) (domain.Warehouse, error)
	Exists(ctx context.Context, warehouseCode string) bool
	Save(ctx context.Context, w domain.Warehouse) (int, error)
	Update(ctx context.Context, w domain.Warehouse) error
	Delete(ctx context.Context, id int) error
	ReportProducts(ctx context.Context, id int) (domain.WarehouseReport, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) ReportProducts(ctx context.Context, id int) (domain.WarehouseReport, error) {
	query := "SELECT w.name AS warehouseName, count(p.id_warehouse) AS totalProducts FROM warehouses w " +
		"INNER JOIN products p ON p.id_warehouse = w.id " +
		"WHERE w.id = ? GROUP BY w.name;"

	row := r.db.QueryRow(query, id)

	w := domain.WarehouseReport{}

	err := row.Scan(&w.WarehouseName, &w.ProductCount)
	if err != nil {
		return domain.WarehouseReport{}, err
	}

	return w, nil
}

func (r *repository) GetAll(ctx context.Context) ([]domain.Warehouse, error) {
	query := "SELECT * FROM warehouses"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	var warehouses []domain.Warehouse

	for rows.Next() {
		w := domain.Warehouse{}
		_ = rows.Scan(&w.ID, &w.Name, &w.Address, &w.Telephone, &w.Capacity)
		warehouses = append(warehouses, w)
	}

	return warehouses, nil
}

func (r *repository) Get(ctx context.Context, id int) (domain.Warehouse, error) {
	query := "SELECT * FROM warehouses WHERE id=?;"
	row, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		log.Fatal(err)
		return domain.Warehouse{}, err
	}
	w := domain.Warehouse{}
	err = row.Scan(&w.ID, &w.Name, &w.Address, &w.Telephone, &w.Capacity)
	if err != nil {
		return domain.Warehouse{}, err
	}

	return w, nil
}

func (r *repository) Exists(ctx context.Context, warehouseCode string) bool {
	query := "SELECT warehouse_code FROM warehouses WHERE warehouse_code=?;"
	row := r.db.QueryRow(query, warehouseCode)
	err := row.Scan(&warehouseCode)
	return err == nil
}

func (r *repository) Save(ctx context.Context, w domain.Warehouse) (int, error) {
	query := "INSERT INTO warehouses (name, adress, telephone, capacity) VALUES (?, ?, ?, ?)"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(&w.Name, &w.Address, &w.Telephone, &w.Capacity)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *repository) Update(ctx context.Context, w domain.Warehouse) error {
	query := "UPDATE warehouses SET name=?, adress=?, telephone=?, capacity=? WHERE id=?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(&w.Name, &w.Address, &w.Telephone, &w.Capacity, &w.ID)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM warehouses WHERE id=?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affect < 1 {
		return ErrNotFound
	}

	return nil
}
