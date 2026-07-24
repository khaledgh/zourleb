// Package router registers all HTTP routes onto an Echo instance.
package router

import (
	"github.com/labstack/echo/v4"
	emw "github.com/labstack/echo/v4/middleware"

	"github.com/zourleb/zourleb-api/internal/app"
	"github.com/zourleb/zourleb-api/internal/middleware"
)

// Register mounts global middleware and all /api/v1 routes.
func Register(e *echo.Echo, c *app.Container) {
	e.Use(emw.Recover())
	e.Use(emw.CORSWithConfig(emw.CORSConfig{
		AllowOrigins: c.Cfg.HTTP.AllowOrigins,
		AllowHeaders: []string{"Authorization", "Content-Type", "Accept-Language", "X-Agency-ID", "X-Request-ID"},
		AllowMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
	}))
	e.Use(middleware.RequestIDMiddleware())
	e.Use(middleware.LocaleMiddleware(c.Locale))

	// static uploads (local storage driver)
	e.Static("/uploads", c.Cfg.Storage.LocalPath)

	e.GET("/health", c.Health.Live)
	e.GET("/ready", c.Health.Ready)

	v1 := e.Group("/api/v1")
	registerAuth(v1, c)
	registerPublic(v1, c)
	registerMe(v1, c)
	registerBookings(v1, c)
	registerEngagement(v1, c)
	registerAgency(v1, c)
	registerShop(v1, c)
	registerAdmin(v1, c)
}

// registerShop mounts the toggleable shop module behind the shop.enabled flag.
func registerShop(g *echo.Group, c *app.Container) {
	gate := middleware.RequireFeature(c.Settings, "shop.enabled")

	shop := g.Group("/shop", gate)
	shop.GET("/products", c.Shop.ListProducts)
	shop.GET("/products/:slug", c.Shop.GetProduct)
	shop.POST("/orders", c.Shop.CreateOrder, c.Auth.Required())
	shop.GET("/orders", c.Shop.ListOrders, c.Auth.Required())

	// agency product management (gated + agency-scoped + permission)
	ap := g.Group("/agency/products", gate, c.Auth.Required(), c.Auth.RequireAgencyScope(), c.Auth.RequirePermission("product.manage"))
	ap.GET("", c.Shop.AgencyProducts)
	ap.POST("", c.Shop.CreateProduct)
	ap.PUT("/:id", c.Shop.UpdateProduct)
}

func registerAuth(g *echo.Group, c *app.Container) {
	a := g.Group("/auth")
	a.POST("/register", c.AuthH.Register)
	a.POST("/login", c.AuthH.Login)
	a.POST("/google", c.AuthH.Google)
	a.POST("/refresh", c.AuthH.Refresh)
	a.POST("/logout", c.AuthH.Logout)

	// OTP — verify is auth-optional so it can attach the phone to a user
	g.POST("/otp/request", c.OTP.Request)
	g.POST("/otp/verify", c.OTP.Verify, c.Auth.Optional())
}

func registerPublic(g *echo.Group, c *app.Container) {
	g.GET("/languages", c.I18nH.Languages)
	g.GET("/i18n/:locale", c.I18nH.Bundle)

	g.GET("/home", c.Catalog.Home)
	g.GET("/tours", c.Catalog.ListTours)
	g.GET("/tours/:slug", c.Catalog.GetTour)
	g.GET("/categories", c.Catalog.Categories)
	g.GET("/regions", c.Catalog.Regions)
	g.GET("/agencies/:slug", c.Catalog.GetAgency)
	g.GET("/reviews", c.Engage.TourReviews)
}

func registerMe(g *echo.Group, c *app.Container) {
	me := g.Group("/me", c.Auth.Required())
	me.GET("", c.Account.Me)
	me.PATCH("", c.Account.Update)
	me.POST("/devices", c.Account.RegisterDevice)

	n := g.Group("/notifications", c.Auth.Required())
	n.GET("", c.Notif.List)
	n.POST("/:id/read", c.Notif.MarkRead)
}

func registerBookings(g *echo.Group, c *app.Container) {
	b := g.Group("/bookings", c.Auth.Required())
	b.POST("", c.Booking.Create)
	b.GET("", c.Booking.List)
	b.GET("/:code", c.Booking.Get)
	b.POST("/:code/pay", c.Booking.Pay)
	b.POST("/:code/cancel", c.Booking.Cancel)
}

func registerEngagement(g *echo.Group, c *app.Container) {
	auth := c.Auth.Required()
	g.GET("/favorites", c.Engage.ListFavorites, auth)
	g.POST("/favorites/:tourId", c.Engage.AddFavorite, auth)
	g.DELETE("/favorites/:tourId", c.Engage.RemoveFavorite, auth)
	g.POST("/reviews", c.Engage.CreateReview, auth)
}

// registerAgency mounts the agency portal. Apply is open to any authenticated
// user; the rest require an agency-scoped role + the relevant permission.
func registerAgency(g *echo.Group, c *app.Container) {
	g.POST("/agency/apply", c.Agency.Apply, c.Auth.Required())

	ag := g.Group("/agency", c.Auth.Required(), c.Auth.RequireAgencyScope())

	ag.GET("/profile", c.Agency.Profile)
	ag.PATCH("/profile", c.Agency.UpdateProfile, c.Auth.RequirePermission("agency.update"))
	ag.GET("/members", c.Agency.Members)

	ag.GET("/tours", c.Agency.ListTours)
	ag.POST("/tours", c.Agency.CreateTour, c.Auth.RequirePermission("tour.create"))
	ag.GET("/tours/:id", c.Agency.GetTour)
	ag.PUT("/tours/:id", c.Agency.UpdateTour, c.Auth.RequirePermission("tour.update"))
	ag.POST("/tours/:id/publish", c.Agency.Publish, c.Auth.RequirePermission("tour.publish"))
	ag.POST("/tours/:id/departures", c.Agency.AddDeparture, c.Auth.RequirePermission("departure.manage"))
	ag.POST("/tours/:id/prices", c.Agency.AddPrice, c.Auth.RequirePermission("price.manage"))
	ag.POST("/tours/:id/images", c.Agency.AddImage, c.Auth.RequirePermission("tour.image.manage"))
	ag.DELETE("/images/:imageId", c.Agency.DeleteImage, c.Auth.RequirePermission("tour.image.manage"))

	ag.GET("/bookings", c.Agency.Bookings, c.Auth.RequirePermission("booking.view"))

	ag.POST("/uploads/image", c.Upload.Image, c.Auth.RequirePermission("tour.image.manage"))

	// boosts
	ag.GET("/boost-packages", c.Boost.Packages, c.Auth.RequirePermission("boost.view"))
	ag.GET("/boosts", c.Boost.List, c.Auth.RequirePermission("boost.view"))
	ag.POST("/boosts", c.Boost.Create, c.Auth.RequirePermission("boost.purchase"))
	ag.POST("/boosts/:id/pay", c.Boost.Pay, c.Auth.RequirePermission("boost.purchase"))
}

// registerAdmin mounts the super-admin surface behind permission checks.
func registerAdmin(g *echo.Group, c *app.Container) {
	ad := g.Group("/admin", c.Auth.Required())

	ad.GET("/agencies", c.Admin.ListAgencies, c.Auth.RequirePermission("admin.agency.approve"))
	ad.POST("/agencies/:id/approve", c.Admin.ApproveAgency, c.Auth.RequirePermission("admin.agency.approve"))

	ad.GET("/languages", c.Admin.Languages, c.Auth.RequirePermission("admin.language.manage"))
	ad.POST("/languages", c.Admin.CreateLanguage, c.Auth.RequirePermission("admin.language.manage"))
	ad.PUT("/languages/:id", c.Admin.UpdateLanguage, c.Auth.RequirePermission("admin.language.manage"))

	ad.POST("/translations", c.Admin.SaveTranslation, c.Auth.RequirePermission("admin.translation.manage"))
	ad.POST("/translations/bulk", c.Admin.BulkTranslations, c.Auth.RequirePermission("admin.translation.manage"))

	ad.GET("/settings", c.Admin.Settings, c.Auth.RequirePermission("admin.settings.manage"))
	ad.POST("/settings", c.Admin.SaveSetting, c.Auth.RequirePermission("admin.settings.manage"))

	ad.GET("/banners", c.Admin.Banners, c.Auth.RequirePermission("admin.banner.manage"))
	ad.POST("/banners", c.Admin.CreateBanner, c.Auth.RequirePermission("admin.banner.manage"))
	ad.DELETE("/banners/:id", c.Admin.DeleteBanner, c.Auth.RequirePermission("admin.banner.manage"))

	ad.POST("/categories", c.Admin.CreateCategory, c.Auth.RequirePermission("admin.category.manage"))
	ad.POST("/regions", c.Admin.CreateRegion, c.Auth.RequirePermission("admin.region.manage"))

	ad.GET("/reviews/pending", c.Admin.PendingReviews, c.Auth.RequirePermission("admin.review.moderate"))
	ad.POST("/reviews/:id/moderate", c.Admin.ModerateReview, c.Auth.RequirePermission("admin.review.moderate"))

	ad.GET("/users", c.Admin.ListUsers, c.Auth.RequirePermission("admin.user.manage"))
	ad.POST("/users/:id/status", c.Admin.SetUserStatus, c.Auth.RequirePermission("admin.user.manage"))

	ad.GET("/boosts/pending", c.Admin.PendingBoosts, c.Auth.RequirePermission("admin.boost.manage"))
	ad.POST("/payments/confirm", c.Admin.ConfirmPayment, c.Auth.RequirePermission("admin.boost.manage"))

	ad.GET("/analytics", c.Admin.Analytics, c.Auth.RequirePermission("admin.analytics.view"))
}
