package handlers

import (
	"encoding/json"
	"net/http"
	"repository_class/internal/domain"
	"repository_class/internal/warehouse"
	"repository_class/pkg/web"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Struct for warehouse with service
type Warehouse struct {
	warehouseService warehouse.Service
}

// Constructor for warehouse with service
func NewWarehouse(w warehouse.Service) *Warehouse {
	return &Warehouse{
		warehouseService: w,
	}
}

// GetOne for warehouse
//
// @Summary		GetOne for warehouse
// @Description	Get a warehouse by id
// @Tags		Warehouse
// @Produce		json
// @Param		id	path		int	true	"warehouse ID"
// @Success		200	{object}	domain.Warehouse
// @Failure		400	{string}	string	"Bad request"
// @Failure		404 {string}	string	"warehouse not found"
// @Failure		500	{string}	string	"Internal server error"
// @Router		/warehouse/{id} [get]
func (w *Warehouse) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, warehouse.ErrInvalidId.Error())
			return
		}
		warehouse, err := w.warehouseService.Get(c, id)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		if warehouse == (domain.Warehouse{}) {
			web.Error(c, http.StatusInternalServerError, err.Error())
		}
		web.Success(c, http.StatusOK, warehouse)
	}
}

func (w *Warehouse) ReportProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Query("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, warehouse.ErrInvalidId.Error())
			return
		}
		warehouse, err := w.warehouseService.ReportProducts(c, id)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		if warehouse == (domain.WarehouseReport{}) {
			web.Error(c, http.StatusInternalServerError, err.Error())
		}
		web.Success(c, http.StatusOK, warehouse)
	}
}

// GetAll for warehouse
//
// @Summary		GetAll for warehouse
// @Description	Get a list of all created warehouse
// @Tags		Warehouse
// @Produce		json
// @Success		200	{object}	[]domain.Warehouse
// @Failure		500	{string}	string	"Internal server error"
// @Router		/warehouse [get]
func (w *Warehouse) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		warehouse, err := w.warehouseService.GetAll(c)

		if err != nil {
			web.Error(c, http.StatusInternalServerError, err.Error())
		}
		web.Success(c, http.StatusOK, warehouse)
	}
}

// Create warehouse
//
// @Summary		Create a warehouse
// @Description	Create a warehouse with the struct domain.warehouse
// @Tags		Warehouse
// @Accept 		json
// @Produce		json
// @Param		warehouse	body	domain.Warehouse	true	"Add warehouse"
// @Success		201	{object}	domain.Warehouse
// @Failure		400	{string}	string	"Bad request"
// @Failure		409 {string}	string	"warehouse number is already registered"
// @Failure		422 {string}	string	"Invalid input structure for warehouse"
// @Failure		500	{string}	string	"Internal server error"
// @Router		/warehouse [post]
func (w *Warehouse) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var war domain.Warehouse
		err := c.ShouldBindJSON(&war)
		if err != nil {
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		// if war.WarehouseCode == "" {
		// 	web.Error(c, http.StatusUnprocessableEntity, err.Error())
		// 	return
		// }
		warehouse, err := w.warehouseService.Create(c, war)
		if err != nil {
			web.Error(c, http.StatusConflict, err.Error())
			return
		}
		web.Success(c, http.StatusCreated, warehouse)
	}
}

// Update a warehouse
//
// @Summary		Update a warehouse
// @Description	Update a warehouse or some fields of it
// @Tags		Warehouse
// @Accept 		json
// @Produce		json
// @Param		id	path	int	true	"Warehouse ID"
// @Param		warehouseUpdate	body	domain.Warehouse	true	"Update Warehouse"
// @Success		200	{object}	domain.Warehouse
// @Failure		400	{string}	string	"Bad request"
// @Failure		404 {string}	string	"Warehouse not found"
// @Failure		409 {string}	string	"Warehouse number is already registered"
// @Failure		500	{string}	string	"Internal server error"
// @Router		/Warehouse/{id} [patch]
func (w *Warehouse) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, warehouse.ErrInvalidId.Error())
			return
		}
		warehouse, err := w.warehouseService.Get(c, id)
		if err != nil {
			web.Error(c, http.StatusNotFound, err.Error())
			return
		}
		if warehouse == (domain.Warehouse{}) {
			web.Error(c, http.StatusNotFound, err.Error())
		}
		err = json.NewDecoder(c.Request.Body).Decode(&warehouse)
		if err != nil {
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		war, eror := w.warehouseService.Update(c, warehouse, id)
		if eror != nil {
			web.Error(c, http.StatusNotFound, eror.Error())
		}
		web.Success(c, http.StatusOK, war)
	}
}

// Delete a warehouse
//
// @Summary		Delete a warehouse
// @Description	Delete a warehouse by id
// @Tags		Warehouse
// @Produce		json
// @Param		id	path	int	true	"Delete Warehouse"
// @Success		204	{string}	string ""
// @Failure		400	{string}	string	"Bad request"
// @Failure		404 {string}	string	"Warehouse not found"
// @Failure		500	{string}	string	"Internal server error"
// @Router		/warehouse/{id} [delete]
func (w *Warehouse) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, warehouse.ErrInvalidId.Error())
			return
		}
		err = w.warehouseService.Delete(c, id)
		if err != nil {
			web.Error(c, http.StatusNotFound, err.Error())
			return
		}
		web.Success(c, http.StatusNoContent, "")
	}
}
