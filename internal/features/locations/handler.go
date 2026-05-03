package locations

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetRegions(c *gin.Context) {
	regions, err := h.service.GetRegions(c.Request.Context())
	if err != nil {
		fmt.Printf("[GetRegions] {GetRegions}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve regions",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, regions)
}

func (h *Handler) GetProvincesByRegion(c *gin.Context) {
	regionCode := c.Param("regionCode")
	if regionCode == "" {
		response.SendFail(c, gin.H{"error": "Invalid region code"})
		return
	}

	provinces, err := h.service.GetProvincesByRegion(
		c.Request.Context(),
		regionCode,
	)
	if err != nil {
		fmt.Printf("[GetProvincesByRegion] {GetProvinces}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve provinces",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, provinces)
}

func (h *Handler) GetCitiesByProvince(c *gin.Context) {
	provinceCode := c.Param("provinceCode")
	if provinceCode == "" {
		response.SendFail(c, gin.H{"error": "Invalid province code"})
		return
	}

	cities, err := h.service.GetCitiesByProvince(
		c.Request.Context(),
		provinceCode,
	)
	if err != nil {
		fmt.Printf("[GetCitiesByProvince] {GetCities}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve cities",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, cities)
}

func (h *Handler) GetCitiesByRegion(c *gin.Context) {
	regionCode := c.Param("regionCode")
	if regionCode == "" {
		response.SendFail(c, gin.H{"error": "Invalid region code"})
		return
	}

	cities, err := h.service.GetCitiesByRegion(c.Request.Context(), regionCode)
	if err != nil {
		fmt.Printf("[GetCitiesByRegion] {GetCities}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve cities",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, cities)
}

func (h *Handler) GetBarangaysByCity(c *gin.Context) {
	cityCode := c.Param("cityCode")
	if cityCode == "" {
		response.SendFail(c, gin.H{"error": "Invalid city code"})
		return
	}

	barangays, err := h.service.GetBarangaysByCity(
		c.Request.Context(),
		cityCode,
	)
	if err != nil {
		fmt.Printf("[GetBarangaysByCity] {GetBarangays}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve barangays",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, barangays)
}
