package routes

import (
	"database/sql"
	"repository_class/cmd/server/handlers"
	"repository_class/internal/product"

	"github.com/gin-gonic/gin"
)

type Router interface {
	MapRoutes()
}

type router struct {
	eng *gin.Engine
	rg  *gin.RouterGroup
	db  *sql.DB
}

func NewRouter(eng *gin.Engine, db *sql.DB) Router {
	return &router{eng: eng, db: db}
}

func (r *router) MapRoutes() {
	r.setGroup()

	r.buildProductsRoutes()
}

func (r *router) setGroup() {
	r.rg = r.eng.Group("/api/v1")
}

func (r *router) buildProductsRoutes() {
	productRepository := product.NewRepository(r.db)
	productService := product.NewService(&productRepository)
	productHandler := handlers.NewProduct(productService)
	routerProduct := r.rg.Group("/products")

	// Products routes
	{
		routerProduct.GET("/", productHandler.GetAll())
		routerProduct.POST("", productHandler.Create())
		routerProduct.GET("/:id", productHandler.Get())
		routerProduct.DELETE("/:id", productHandler.Delete())
		routerProduct.PATCH("/:id", productHandler.Update())
	}
}
