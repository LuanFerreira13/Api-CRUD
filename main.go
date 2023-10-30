package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func routeHearth(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})

	c.Done()
}

type Estudante struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
	Age      int    `json:"age"`
}

var Estudantes = []Estudante{
	Estudante{ID: 1, FullName: "Luan", Age: 24},
	Estudante{ID: 2, FullName: "Joao", Age: 23},
}

func routeGetEstudantes(c *gin.Context) {
	c.JSON(http.StatusOK, Estudantes)
}

func routePostEstudantes(c *gin.Context) {
	var estudante Estudante

	err := c.Bind(&estudante)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Não foi possivel obter o payload",
		})
		return
	}

	estudante.ID = Estudantes[len(Estudantes)-1].ID + 1
	Estudantes = append(Estudantes, estudante)

	c.JSON(http.StatusCreated, estudante)
}

func routePutEstudantes(c *gin.Context) {
	var estudantePayload Estudante
	var estudanteLocal Estudante
	var newEstudantes []Estudante

	err := c.BindJSON(&estudantePayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Não foi possivel obter o payload",
		})
		return
	}

	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Não foi possivel obter o id",
		})
		return
	}

	for _, estudanteElemt := range Estudantes {
		if estudanteElemt.ID == id {
			estudanteLocal = estudanteElemt
		}
	}

	if estudanteLocal.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Não foi possivel encontrar estudante",
		})
		return
	}

	estudanteLocal.FullName = estudantePayload.FullName
	estudanteLocal.Age = estudantePayload.Age

	for _, estudanteElement := range Estudantes {
		if id == estudanteElement.ID {
			newEstudantes = append(newEstudantes, estudanteLocal)
		} else {
			newEstudantes = append(newEstudantes, estudanteElement)
		}
	}

	Estudantes = newEstudantes
	c.JSON(http.StatusOK, estudanteLocal)
}

func routeDeleteEstudantes(c *gin.Context) {
	var newEstudantes []Estudante

	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Não foi possivel obter o id",
		})
		return
	}

	for _, estudanteElement := range Estudantes {
		if estudanteElement.ID != id {
			newEstudantes = append(newEstudantes, estudanteElement)
		}
	}

	Estudantes = newEstudantes

	c.JSON(http.StatusOK, gin.H{
		"message": "Estudante excluido com sucesso",
	})
}

func routeGetEstudante(c *gin.Context) {
	var estudante Estudante
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Não foi possivel obter o id",
		})
		return
	}

	for _, estudanteElement := range Estudantes {
		if estudanteElement.ID == id {
			estudante = estudanteElement
		}
	}

	if estudante.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Não foi possivel encontrar o estudante",
		})
		return
	}

	c.JSON(http.StatusOK, estudante)

}

func getRoutes(c *gin.Engine) *gin.Engine {
	c.GET("/heart", routeHearth)

	groupEstudantes := c.Group("/estudantes")
	groupEstudantes.GET("/", routeGetEstudantes)
	groupEstudantes.POST("/", routePostEstudantes)
	groupEstudantes.PUT("/:id", routePutEstudantes)
	groupEstudantes.DELETE("/:id", routeDeleteEstudantes)
	groupEstudantes.GET("/:id", routeGetEstudante)

	return c
}

func main() {
	service := gin.Default()

	getRoutes(service)

	service.Run()
}
