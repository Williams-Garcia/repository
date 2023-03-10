package product

import (
	"context"
	"database/sql"
	"log"

	"repository_class/internal/domain"
)

// Repository encapsulates the storage of a Product.
type Repository interface {
	GetAll(ctx context.Context) ([]domain.Product, error)
	Get(ctx context.Context, id int) (domain.Product, error)
	GetWithWarehouse(ctx context.Context, id int) (domain.ProductWithWarehouse, error)
	Exists(ctx context.Context, productCode string) bool
	Save(ctx context.Context, p domain.Product) (int, error)
	Update(ctx context.Context, p domain.Product) error
	Delete(ctx context.Context, id int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAll(ctx context.Context) ([]domain.Product, error) {
	query := "SELECT * FROM products;"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	var products []domain.Product

	for rows.Next() {
		p := domain.Product{}
		_ = rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.CodeValue, &p.IsPublished, &p.Expiration, &p.Price)
		products = append(products, p)
	}

	return products, nil
}

func (r *repository) Get(ctx context.Context, id int) (domain.Product, error) {
	query := "SELECT * FROM products WHERE id=?;"
	row := r.db.QueryRow(query, id)
	p := domain.Product{}
	err := row.Scan(&p.ID, &p.Name, &p.Quantity, &p.CodeValue, &p.IsPublished, &p.Expiration, &p.Price)
	if err != nil {
		return domain.Product{}, err
	}

	return p, nil
}

func (r *repository) GetWithWarehouse(ctx context.Context, id int) (domain.ProductWithWarehouse, error) {
	// query := "SELECT * FROM products WHERE id=?;"
	query := "SELECT p.id , p.name, p.quantity, p.code_value, p.is_published, p.expiration, p.price, w.id AS warehouseId, " +
		"w.name, w.adress, w.telephone, w.capacity " +
		"FROM products p " +
		"INNER JOIN warehouses w ON w.id = p.id_warehouse " +
		"WHERE p.id = ?"
	row := r.db.QueryRow(query, id)
	p := domain.ProductWithWarehouse{
		Product:   domain.Product{},
		Warehouse: domain.Warehouse{},
	}
	err := row.Scan(&p.Product.ID, &p.Product.Name, &p.Product.Quantity, &p.Product.CodeValue, &p.Product.IsPublished,
		&p.Product.Expiration, &p.Product.Price, &p.Warehouse.ID, &p.Warehouse.Name, &p.Warehouse.Address, &p.Warehouse.Telephone, &p.Warehouse.Capacity,
	)
	if err != nil {
		log.Fatal(err)
		return domain.ProductWithWarehouse{}, err
	}

	return p, nil
}

func (r *repository) Exists(ctx context.Context, productCode string) bool {
	query := "SELECT product_code FROM products WHERE product_code=?;"
	row := r.db.QueryRow(query, productCode)
	err := row.Scan(&productCode)
	return err == nil
}

func (r *repository) Save(ctx context.Context, p domain.Product) (int, error) {
	query := "INSERT INTO products(name,quantity,code_value,is_published,expiration,price,id_warehouse) VALUES (?,?,?,?,?,?,?)"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(p.Name, p.Quantity, p.CodeValue, p.IsPublished, p.Expiration, p.Price, p.IdWarehouse)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *repository) Update(ctx context.Context, p domain.Product) error {
	query := "UPDATE products SET name=?, quantity=?, code_value=?, is_published=?, expiration=?, price=? WHERE id=?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(p.Name, p.Quantity, p.CodeValue, p.IsPublished, p.Expiration, p.Price, p.ID)
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
	query := "DELETE FROM products WHERE id=?"
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
