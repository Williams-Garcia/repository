package warehouse

import (
	"context"
	"database/sql"
	"testing"

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
