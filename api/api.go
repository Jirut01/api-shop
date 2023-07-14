package api

import (
	apiProduct "app-backend/api/product"
	apiLogin "app-backend/api/login"
	"app-backend/middleware"
	"app-backend/setting"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/patcharp/golib/v2/server"
)

func Register() error {
	app := setting.GetCfg().Server.App()
	skipper := server.NewSkipperPath("/api")
	app.Use(cors.New())
	api := app.Group("/api")
	api.Post("/login", apiLogin.Login)
	api.Get("/logout", apiLogin.Logout)
	api.Get("/getuser", middleware.AuthSystem(&skipper), apiLogin.GetUser)


	admin := api.Group("/admin", middleware.AdminAuth(&skipper))
	admin.Post("/register", apiLogin.RegistorUser)

	// middleware.ClearCheckApiCache()   เคลีย chache
 
	// ===== Product ======
	product := api.Group("product", middleware.AuthSystem(&skipper))
	product.Post("/add-prodcut", apiProduct.AddProduct)
	product.Post("/upload-img-product",apiProduct.ImgProductUpload, middleware.ClearCheckApiCache())
	product.Delete("/remove-img-product",apiProduct.ImgProductRemove)

	return nil
}