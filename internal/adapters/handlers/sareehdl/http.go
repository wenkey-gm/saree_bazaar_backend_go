package sareehdl

import (
	"product_api/internal/core/domain"
	"product_api/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type SareeHandler struct {
	sareeService ports.ISareeService
}

func NewSareeHandler(sareeService ports.ISareeService) *SareeHandler {
	return &SareeHandler{
		sareeService: sareeService,
	}
}

func (s *SareeHandler) FindAll(c *gin.Context) {
	sarees, err := s.sareeService.FindAll()
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, sarees)
}

func (s *SareeHandler) Find(c *gin.Context) {
	id := c.Param("id")
	saree, err := s.sareeService.Find(id)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, saree)
}

func (s *SareeHandler) Save(c *gin.Context) {
	var saree domain.Saree
	if err := c.ShouldBindJSON(&saree); err != nil {
		c.JSON(400, err)
		return
	}
	savedSaree, err := s.sareeService.Save(saree)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, savedSaree)
}

func (s *SareeHandler) Update(c *gin.Context) {
	var saree domain.Saree
	if err := c.ShouldBindJSON(&saree); err != nil {
		c.JSON(400, err)
		return
	}
	id := c.Param("id")
	updatedSaree, err := s.sareeService.Update(id, saree)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, updatedSaree)
}

func (s *SareeHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := s.sareeService.Delete(id)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, "Saree deleted successfully")

}
