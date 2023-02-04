package handlers

import (
	"net/http"

	"github.com/EvertonTomalok/go-template/internal/app/domain/dto"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// GO TEMPLATE API
// @Summary Get Person
// @Description Get Person
// @Tags People
// @Router /person/{person_id} [get]
// @Produce json
// @Success 200 {object} dto.Person
func GetPersonById(c *gin.Context) {

	personId := c.Param("personId")

	log.Infof("Fetching person id: %s", personId)
	person := dto.Person{
		Name: "My Name",
		CPF:  "SOME CPF",
	}
	c.JSON(http.StatusOK, person)
}
