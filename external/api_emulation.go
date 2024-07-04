package external

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"timeTracker/app"
)

func RunExternalApiEmulation(myApp *app.App) {
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
		
		// Делаем запрос к randomdatatools для генерации случайных данных
		// Лимит - 1 запрос в секунду
		url := "https://api.randomdatatools.ru/?count=1&typeName=classic&params=LastName,FirstName,FatherName,Address"
		
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error while executing GET request on https://api.randomdatatools.ru: %v\n", err)
			return
		}
		defer resp.Body.Close()
		
		person := struct {
			LastName   string `json:"LastName"`
			FirstName  string `json:"FirstName"`
			FatherName string `json:"FatherName"`
			Address    string `json:"Address"`
		}{}
		
		if resp.StatusCode == http.StatusOK {
			// Читаем тело ответа
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error while reading response body: %v\n", err)
				return
			}
			
			// Декодируем JSON
			err = json.Unmarshal([]byte(string(body)), &person)
			if err != nil {
				log.Printf("Error while decoding JSON: %v\n", err)
				return
			}
		}
		
		// Возвращаем успешный ответ с данными
		c.JSON(http.StatusOK, gin.H{
			"surname":    person.LastName,
			"name":       person.FirstName,
			"patronymic": person.FatherName,
			"address":    person.Address,
		})
	})
	
	myApp.WaitGroup.Done()
	log.Println("External API emulation started on http://localhost:8088")
	
	// Запуск сервера на порту 8088
	if err := r.Run(":8088"); err != nil {
		log.Println("Failed to start external API emulation:", err)
	}
}
