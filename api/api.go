package api

import (
	"app-backend/setting"
	"github.com/gofiber/fiber/v2/middleware/cors"
	apiProduct "app-backend/api/product"
)

func Register() error {
	app := setting.GetCfg().Server.App()
	app.Use(cors.New())
	api := app.Group("/api")

	// ===== Product ======
	product := api.Group("product")

	product.Post("/add-prodcut", apiProduct.AddProduct)
	product.Post("/upload-img-product",apiProduct.ImgProductUpload)
	product.Delete("/remove-img-product",apiProduct.ImgProductRemove)

	return nil
}