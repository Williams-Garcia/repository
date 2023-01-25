package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"repository_class/internal/domain"
	"repository_class/internal/product"
	"repository_class/pkg/web"

	"github.com/gin-gonic/gin"
)

// Errors
var (
	ErrProductInternalServer = errors.New("internal server error")
)

type Product struct {
	service product.Service
}

func NewProduct(p product.Service) *Product {
	return &Product{
		service: p,
	}
}

func (prod *Product) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := prod.service.GetAll(c)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, ErrProductInternalServer.Error())
			return
		}
		web.Success(c, http.StatusOK, products)
	}
}

func (prod *Product) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		p, err := prod.service.Get(c, id)
		if err != nil {
			if errors.Is(err, product.ErrNotFound) {
				web.Error(c, http.StatusNotFound, err.Error())
			}
			web.Error(c, http.StatusInternalServerError, ErrProductInternalServer.Error())
			return
		}
		web.Success(c, http.StatusOK, p)
	}
}

func (prod *Product) GetWithWarehouse() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		p, err := prod.service.GetWithWarehouse(c, id)
		if err != nil {
			if errors.Is(err, product.ErrNotFound) {
				web.Error(c, http.StatusNotFound, err.Error())
			}
			web.Error(c, http.StatusInternalServerError, ErrProductInternalServer.Error())
			return
		}
		web.Success(c, http.StatusOK, p)
	}
}

func (p *Product) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var prod domain.Product
		// check json type
		if err := c.ShouldBindJSON(&prod); err != nil {
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		productCreated, err := p.service.Create(c, prod)
		if err != nil {
			if errors.Is(err, product.ErrProductRegistered) {
				web.Error(c, http.StatusConflict, err.Error())
				return
			} else if errors.Is(err, product.ErrInvalidStruct) {
				web.Error(c, http.StatusUnprocessableEntity, err.Error())
				return
			}
			web.Error(c, http.StatusInternalServerError, ErrProductInternalServer.Error())
			return
		}

		web.Success(c, http.StatusCreated, productCreated)
	}
}

func (p *Product) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		var prod domain.Product
		// check json type
		if err := c.ShouldBindJSON(&prod); err != nil {
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		prod, err = p.service.Update(c, prod, id)
		if err != nil {
			if errors.Is(err, product.ErrProductRegistered) {
				web.Error(c, http.StatusConflict, err.Error())
				return
			} else if errors.Is(err, product.ErrNotFound) {
				web.Error(c, http.StatusNotFound, err.Error())
				return
			}
			web.Error(c, http.StatusInternalServerError, ErrProductInternalServer.Error())
			return
		}
		web.Success(c, http.StatusOK, prod)
	}
}

func (prod *Product) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		err = prod.service.Delete(c, id)
		if err != nil {
			if errors.Is(err, product.ErrNotFound) {
				web.Error(c, http.StatusNotFound, err.Error())
				return
			}
			web.Error(c, http.StatusInternalServerError, ErrProductInternalServer.Error())
			return
		}
		web.Success(c, http.StatusNoContent, prod)
	}
}
