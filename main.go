package main

import (
	"app/common"
	"app/modules/item/model"
	ginItem "app/modules/item/transport/gin"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:Kaitran2003@@tcp(127.0.0.1:3306)/social_todo_list?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln((err))
	}
	fmt.Println(db)
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		items := v1.Group("/items")
		{
			items.POST("", ginItem.CreateItem(db))
			items.GET("", ListItem(db))
			items.GET("/:id", ginItem.GetItem(db))
			items.PATCH("/:id", ginItem.UpdateItem(db))
			items.DELETE("/:id", ginItem.DeleteItem(db))
		}
	}
	r.Run()

}

func ListItem(db *gorm.DB) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var Paging common.Paging
		if err := ctx.ShouldBind(&Paging); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		Paging.Process()

		var result []model.TodoItem

		db = db.Where("status <> ?", "Deleted")

		if err := db.Table(model.TodoItem{}.TableName()).Count(&Paging.Total).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := db.Order("id desc").Offset((Paging.Page - 1) * Paging.Limit).Limit(Paging.Limit).Find(&result).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, common.NewSuccessResponse(result, Paging, nil))
	}
}
