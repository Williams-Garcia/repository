package warehouse

import (
	"context"
	"database/sql"
	"repository_class/internal/domain"
	"testing"
	"time"

	"github.com/DATA-DOG/go-txdb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func init() {
	txdb.Register("txdb", "mysql", "root@/my_db?allowNativePasswords=false&checkConnLiveness=false&parseTime=true&maxAllowedPacket=0")
}

func TestGet(t *testing.T) {
	db, err := sql.Open("txdb", "identifier")
	assert.NoError(t, err)
	defer db.Close()

	rp := NewRepository(db)

	w, err := rp.GetAll(context.Background())

	assert.NoError(t, err)
	assert.NotEmpty(t, w)
}

func TestGetOne(t *testing.T) {
	db, err := sql.Open("txdb", "identifier")
	id := 2
	wr := domain.Warehouse{
		ID:        2,
		Name:      "x",
		Address:   "x",
		Telephone: "x",
		Capacity:  10,
	}
	assert.NoError(t, err)
	defer db.Close()

	rp := NewRepository(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	w, err := rp.Get(ctx, id)

	assert.NoError(t, err)
	assert.NotEmpty(t, w)
	assert.Equal(t, wr.Name, w.Name)
}

func TestGetOneTimeOut(t *testing.T) {
	db, err := sql.Open("txdb", "identifier")
	id := 2
	wr := domain.Warehouse{
		ID:        2,
		Name:      "x",
		Address:   "x",
		Telephone: "x",
		Capacity:  10,
	}
	assert.NoError(t, err)
	defer db.Close()

	rp := NewRepository(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	w, err := rp.Get(ctx, id)

	assert.NoError(t, err)
	assert.NotEmpty(t, w)
	assert.Equal(t, wr.Name, w.Name)
}
