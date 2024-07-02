package external

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func RunExternalApiEmulation() {
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
		
		// Проверка на корректность серии и номера паспорта (можно добавить свои правила проверки)
		if len(passportSerie) != 4 || len(passportNumber) != 6 {
			log.Println("invalid passportSerie or passportNumber format")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid passportSerie or passportNumber format"})
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
		log.Fatalln("Failed to start external API emulation:", err)
	}
	log.Println("External API emulation started on http://localhost:8088")
}
