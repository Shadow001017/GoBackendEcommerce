package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/quangdangfit/gocommon/logger"

	"goshop/app/serializers"
	"goshop/app/services"
	"goshop/pkg/response"
	"goshop/pkg/utils"
	"goshop/pkg/validation"
)

type Product struct {
	validator validation.Validation
	service   services.IProductService
}

func NewProductAPI(service services.IProductService) *Product {
	return &Product{
		validator: validation.New(),
		service:   service,
	}
}

// GetProductByID godoc
// @Summary Get get product by uuid
// @Produce json
// @Param uuid path string true "Product UUID"
// @Security ApiKeyAuth
// @Success 200 {object} serializers.Product
// @Router /api/v1/products/{id} [get]
func (p *Product) GetProductByID(c *gin.Context) {
	productId := c.Param("uuid")

	ctx := c.Request.Context()
	product, err := p.service.GetProductByID(ctx, productId)
	if err != nil {
		logger.Error("Failed to get product: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res serializers.Product
	copier.Copy(&res, &product)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

// ListProducts godoc
// @Summary Get list products
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} serializers.ListProductRes
// @Router /api/v1/products [get]
func (p *Product) ListProducts(c *gin.Context) {
	var req serializers.ListProductReq
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error("Failed to parse request query: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	products, pagination, err := p.service.ListProducts(c, req)
	if err != nil {
		logger.Error("Failed to get products: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res serializers.ListProductRes
	err = copier.Copy(&res.Products, &products)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	res.Pagination = pagination
	response.JSON(c, http.StatusOK, res)
}

// CreateProduct godoc
// @Summary create product
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} serializers.Product
// @Router /api/v1/products [post]
func (p *Product) CreateProduct(c *gin.Context) {
	var req serializers.CreateProductReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	if err := p.validator.ValidateStruct(req); err != nil {
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	ctx := c.Request.Context()
	product, err := p.service.Create(ctx, &req)
	if err != nil {
		logger.Error("Failed to create product", err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res []serializers.Product
	err = copier.Copy(&res, &product)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, res)
}

// UpdateProduct godoc
// @Summary update product
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} serializers.Product
// @Router /api/v1/products/{id} [put]
func (p *Product) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var req serializers.UpdateProductReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	if err := p.validator.ValidateStruct(req); err != nil {
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	ctx := c.Request.Context()
	product, err := p.service.Update(ctx, id, &req)
	if err != nil {
		logger.Error("Failed to create product", err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res []serializers.Product
	err = copier.Copy(&res, &product)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, res)
}
