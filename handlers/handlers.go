package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"timeTracker/app"
	"timeTracker/models"
)

// @Summary GetUsers
// @Description Возвращает список всех пользователей
// @Produce json
// @Param id query string false "ID пользователя"
// @Param passport_series query string false "Серия паспорта пользователя"
// @Param passport_number query string false "Номер паспорта пользователя"
// @Param surname query string false "Фамилия пользователя"
// @Param name query string false "Имя пользователя"
// @Param patronymic query string false "Отчество пользователя"
// @Param address query string false "Адрес пользователя"
// @Param limit query int false "Количество записей на странице" default(10)
// @Param page query int false "Номер страницы" default(1)
// @Success 200 {object} models.User
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Tags 1. Get
// @Router /users [get]
func GetUsers(c *gin.Context) {
	db := c.MustGet("app").(*app.App).DB
	
	id := c.Query("id")
	passportSeries := c.Query("passport_series")
	passportNumber := c.Query("passport_number")
	surname := c.Query("surname")
	name := c.Query("name")
	patronymic := c.Query("patronymic")
	address := c.Query("address")
	
	prepareDB := db
	
	if id != "" {
		prepareDB = prepareDB.Where("id = ?", id)
	}
	if passportSeries != "" {
		prepareDB = prepareDB.Where("passport_series = ?", passportSeries)
	}
	if passportNumber != "" {
		prepareDB = prepareDB.Where("passport_number = ?", passportNumber)
	}
	if surname != "" {
		prepareDB = prepareDB.Where("surname = ?", surname)
	}
	if name != "" {
		prepareDB = prepareDB.Where("name = ?", name)
	}
	if patronymic != "" {
		prepareDB = prepareDB.Where("patronymic = ?", patronymic)
	}
	if address != "" {
		prepareDB = prepareDB.Where("address = ?", address)
	}
	
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}
	
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit number"})
		return
	}
	
	var users []models.User
	offset := (page - 1) * limit
	
	prepareDB.Offset(offset).Limit(limit).Find(&users)
	
	var numPages int64
	prepareDB.Model(&models.User{}).Count(&numPages)
	
	c.JSON(http.StatusOK, gin.H{"current_page": page, "total_pages": numPages/int64(limit) + 1, "result": users})
}

type overviewRequest struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type overviewResponse struct {
	Tasks []struct {
		Task              models.Task `json:"task"`
		TimeSpentInPeriod string      `json:"time_spent_in_period"`
	} `json:"tasks"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// @Summary GetTasksOverview
// @Description Возвращает трудозатраты по пользователю
// @Produce json
// @Param id path int true "ID пользователя"
// @Param period body overviewRequest true "Период"
// @Success 200 {object} overviewResponse
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Tags 1. Get
// @Router /users/{id}/tasks_overview [get]
func GetTasksOverview(c *gin.Context) {
	db := c.MustGet("app").(*app.App).DB
	
	userId := c.Param("id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}
	
	var user models.User
	if err := db.Where("id = ?", userId).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	
	var req overviewRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if req.StartTime == "" || req.EndTime == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_time and end_time is required"})
		return
	}
	
	startTime, err := time.Parse(time.RFC822, req.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Time format must be RFC822, e.g. \"05 Jul 24 19:03 MSK\""})
		return
	}
	endTime, err := time.Parse(time.RFC822, req.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Time format must be RFC822, e.g. \"05 Jul 24 19:03 MSK\""})
		return
	}
	
	if startTime.After(endTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_time must be before end_time"})
		return
	}
	
	var tasks []models.Task
	dbTasks := db.Model(&models.Task{}).Where("user_id = ?", user.ID)
	dbTasks.Find(&tasks)
	
	result := overviewResponse{
		Tasks: make([]struct {
			Task              models.Task `json:"task"`
			TimeSpentInPeriod string      `json:"time_spent_in_period"`
		}, len(tasks)),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
	
	for idx, task := range tasks {
		result.Tasks[idx].Task = task
		
		var periods []models.Period
		db.Model(&models.Period{}).Where("task_id = ? AND user_id = ?", task.ID, userId).Find(&periods)
		
		var sumDuration time.Duration
		
		for _, period := range periods {
			intersectStart := period.StartTime
			if period.StartTime.Before(startTime) {
				intersectStart = startTime
			}
			intersectEnd := period.EndTime
			if task.IsActive && intersectEnd.IsZero() {
				intersectEnd = time.Now()
			}
			if intersectEnd.After(endTime) {
				intersectEnd = endTime
			}
			sumDuration += intersectEnd.Sub(intersectStart)
		}
		fmt.Println(sumDuration)
		
		result.Tasks[idx].TimeSpentInPeriod = fmt.Sprintf("%d:%d", int(sumDuration.Hours()), int(sumDuration.Minutes())%60)
	}
	
	c.JSON(http.StatusOK, result)
}

// @Summary StartTaskTimer
// @Description Начинает отсчет для задачи пользователя, создает задачу, если ее не существует
// @Produce json
// @Param id path int true "ID пользователя"
// @Param task_id path int true "Task ID"
// @Success 200 {object} models.Task
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Tags 2. Post
// @Router /users/{id}/tasks/{task_id}/start [post]
func StartTaskTimer(c *gin.Context) {
	db := c.MustGet("app").(*app.App).DB
	
	userId := c.Param("id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}
	
	userIdInt, err := strconv.Atoi(userId)
	if err != nil || userIdInt < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID must be an non-negative integer"})
	}
	
	var user models.User
	if err := db.Where("id = ?", userId).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	
	taskId := c.Param("task_id")
	if taskId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}
	
	taskIdInt, err := strconv.Atoi(taskId)
	if err != nil || taskIdInt < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID must be an non-negative integer"})
		return
	}
	
	var task models.Task
	if err := db.Model(&models.Task{}).Where("id = ? AND user_id = ?", taskId, userId).First(&task).Error; err != nil {
		task = models.Task{
			ID:               uint(taskIdInt),
			UserID:           user.ID,
			Name:             c.Query("name"),
			IsActive:         false,
			Periods:          []models.Period{},
			OverallTimeSpent: 0,
		}
		db.Create(&task)
	}
	
	if !task.IsActive {
		task.IsActive = true
		db.Save(&task)
		
		newPeriod := models.Period{
			TaskID:    task.ID,
			UserID:    uint(userIdInt),
			StartTime: time.Now(),
		}
		db.Create(&newPeriod)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task is already active"})
		return
	}
	
	db.Model(&models.Task{}).Where("id = ? AND user_id = ?", taskId, userId).First(&task)
	c.JSON(http.StatusOK, task)
}

// @Summary StopTaskTimer
// @Description Заканчивает отсчет для задачи пользователя
// @Produce json
// @Param id path int true "ID пользователя"
// @Param task_id path int true "ID задачи"
// @Success 200 {object} models.Task
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Tags 2. Post
// @Router /users/{id}/tasks/{task_id}/stop [post]
func StopTaskTimer(c *gin.Context) {
	db := c.MustGet("app").(*app.App).DB
	
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}
	
	var user models.User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	
	taskId := c.Param("task_id")
	if taskId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}
	
	taskIdInt, err := strconv.Atoi(taskId)
	if err != nil || taskIdInt < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID must be an non-negative integer"})
		return
	}
	
	var task models.Task
	if err := db.Where("id = ?", taskId).Last(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	
	var periodCount int
	if err := db.Model(&models.Period{}).Count(&periodCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get period count"})
		return
	}
	
	var periods []models.Period
	
	if err := db.Model(&models.Period{}).Where("task_id = ?", task.ID).First(&periods).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get task periods"})
		return
	}
	
	if task.IsActive {
		task.IsActive = false
		lastPeriod := &periods[len(periods)-1]
		lastPeriod.EndTime = time.Now()
		task.OverallTimeSpent = lastPeriod.EndTime.Sub(lastPeriod.StartTime)
		db.Save(&task)
		db.Save(lastPeriod)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task is not active"})
		return
	}
	
	db.Where("id = ?", taskId).Last(&task)
	db.Model(&models.Period{}).Where("task_id = ?", task.ID).First(&task.Periods)
	c.JSON(http.StatusOK, task)
}

// @Summary DeleteUser
// @Description Удаляет пользователя и все связанные с ними записи
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Tags 4. Delete
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	db := c.MustGet("app").(*app.App).DB
	
	userId := c.Param("id")
	var user models.User
	if err := db.Model(&models.User{}).Where("id = ?", userId).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	
	var tasks []models.Task
	db.Model(&models.Task{}).Where("id = ?", userId).Find(&tasks)
	
	var periods []models.Period
	db.Model(&models.Period{}).Where("id = ?", userId).Find(&periods)
	
	db.Delete(&user)
	db.Delete(&tasks)
	db.Delete(&periods)
	
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

type updateUserRequest struct {
	PassportSeries string `json:"passport_series"`
	PassportNumber string `json:"passport_number"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
	Patronymic     string `json:"patronymic"`
	Address        string `json:"address"`
}

// @Summary UpdateUser
// @Description Обновляет информацию о пользователе
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param passport_series body updateUserRequest false "Серия паспорта"
// @Param passport_number body updateUserRequest false "Номер паспорта"
// @Param surname body updateUserRequest false "Фамилия"
// @Param name body updateUserRequest false "Имя"
// @Param patronymic body updateUserRequest false "Отчество"
// @Param address body updateUserRequest false "Адрес"
// @Success 200 {object} models.User
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Tags 3. Put
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	db := c.MustGet("app").(*app.App).DB
	
	id := c.Param("id")
	
	idInt, err := strconv.Atoi(id)
	if err != nil || idInt < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID must be an non-negative integer"})
		return
	}
	
	var oldUser models.User
	if err := db.Where("id = ?", id).First(&oldUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	
	var req updateUserRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if req.PassportSeries == "" {
		req.PassportSeries = oldUser.PassportSeries
	}
	if req.PassportNumber == "" {
		req.PassportNumber = oldUser.PassportNumber
	}
	if req.Surname == "" {
		req.Surname = oldUser.Surname
	}
	if req.Name == "" {
		req.Name = oldUser.Name
	}
	if req.Patronymic == "" {
		req.Patronymic = oldUser.Patronymic
	}
	if req.Address == "" {
		req.Address = oldUser.Address
	}
	
	newUser := models.User{
		ID:             uint(idInt),
		PassportSeries: req.PassportSeries,
		PassportNumber: req.PassportNumber,
		Surname:        req.Surname,
		Name:           req.Name,
		Patronymic:     req.Patronymic,
		Address:        req.Address,
	}
	db.Save(&newUser)
	
	c.JSON(http.StatusOK, newUser)
}

type passportRequest struct {
	PassportNumber string `json:"passportNumber"`
}

// @Summary CreateUser
// @Description Создает пользователя
// @Accept json
// @Produce json
// @Param passportNumber body passportRequest true "Серия и номер паспорта"
// @Success 200 {object} models.User
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Tags 2. Post
// @Router /users [post]
func CreateUser(c *gin.Context) {
	db := c.MustGet("app").(*app.App).DB
	
	var req passportRequest
	
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	passport := req.PassportNumber
	if passport == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "passportNumber is required"})
		return
	}
	
	parts := strings.Split(passport, " ")
	if len(parts) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect passport format"})
		return
	}
	
	passportSeries, passportNumber := parts[0], parts[1]
	if passportSeries == "" || passportNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "passport_series and passport_number are required"})
		return
	}
	
	var cnt int
	db.Model(&models.User{}).Where("passport_series = ? AND passport_number = ?", passportSeries, passportNumber).Count(&cnt)
	if cnt != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		log.Println("User already exists")
		return
	}
	
	// Запрос к внешнему API
	url := fmt.Sprintf("http://localhost:8088/info?passportSerie=%s&passportNumber=%s", passportSeries, passportNumber)
	
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to external API"})
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to close response body"})
		}
	}(resp.Body)
	
	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response from external API"})
		return
	}
	
	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		errorBody := gin.H{}
		err := json.Unmarshal(body, &errorBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to parse response from external API: %e", err)})
			return
		}
		c.JSON(resp.StatusCode, errorBody)
		return
	}
	
	// Парсим JSON-ответ
	newUser := models.User{
		PassportSeries: passportSeries,
		PassportNumber: passportNumber,
	}
	
	if err = json.Unmarshal(body, &newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response from external API"})
		return
	}
	
	db.Create(&newUser)
	c.JSON(http.StatusOK, newUser)
}

type taskNameRequest struct {
	TaskName string `json:"task_name"`
}

// @Summary CreateTask
// @Description Создает задачу для пользователя
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param name body taskNameRequest false "Название задачи"
// @Success 200 {object} models.Task
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Tags 2. Post
// @Router /users/{id}/tasks [post]
func CreateTask(c *gin.Context) {
	db := c.MustGet("app").(*app.App).DB
	
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil || idInt < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID must be an non-negative integer"})
		return
	}
	
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}
	
	var user models.User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	
	var req taskNameRequest
	
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	task := models.Task{UserID: uint(idInt), Name: req.TaskName}
	db.Create(&task)
	
	c.JSON(http.StatusOK, &task)
}
