package ginItem

import (
	"app/common"
	"app/modules/item/biz"
	"app/modules/item/model"
	"app/modules/item/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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

		var Filter model.Filter
		if err := ctx.ShouldBind(&Filter); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		store := storage.NewSQLStore(db)
		business := biz.NewListItemBiz(store)

		result, err := business.ListItemById(ctx.Request.Context(), &Filter, &Paging)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, common.NewSuccessResponse(result, Paging, Filter))
	}
}
