package external

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"timeTracker/app"
)

func RunExternalApiEmulation(myApp *app.App) {
	defer myApp.WaitGroup.Done()
	
	r := gin.Default()
	
	r.GET("/info", func(c *gin.Context) {
		passportSerie := c.Query("passportSerie")
		passportNumber := c.Query("passportNumber")
		
		// Проверка на обязательные параметры
		if passportSerie == "" || passportNumber == "" {
			log.Println("passportSerie and passportNumber are required")
			c.JSON(http.StatusBadRequest, gin.H{"error": "passportSerie and passportNumber are required"})
			return
		}
		
		// Проверка на корректность серии и номера паспорта
		if len(passportSerie) != 4 {
			log.Println("invalid passportSerie format")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid passportSerie format"})
			return
		}
		
		if len(passportNumber) != 6 {
			log.Println("invalid passportNumber format")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid passportNumber format"})
			return
		}
		
		// Возвращаем успешный ответ с данными
		c.JSON(http.StatusOK, gin.H{
			"surname":    "Иванов",
			"name":       "Иван",
			"patronymic": "Иванович",
			"address":    "г. Москва, ул. Ленина, д. 5, кв. 1",
		})
	})
	
	// Запуск сервера на порту 8088
	if err := r.Run(":8088"); err != nil {
		log.Println("Failed to start external API emulation:", err)
	}
	log.Println("External API emulation started on http://localhost:8088")
}
