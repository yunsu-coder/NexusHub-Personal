package router

import (
	"nexushub-personal/internal/database"
	"nexushub-personal/internal/handler"
	"nexushub-personal/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORS())
	r.Use(middleware.RequestLogger())

	// Serve static files from uploads directory
	r.Static("/uploads", "./storage/uploads")

	// Root handler
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "online",
			"message": "NexusHub Backend API Server",
			"version": "v3.0.1",
		})
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Auth routes (public)
		authHandler := handler.NewAuthHandler()
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Health check (v1)
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Apply optional authentication - allows guest access
		v1.Use(middleware.OptionalAuthMiddleware())

		// User profile (authenticated users only)
		v1.GET("/profile", authHandler.GetProfile)

		// Notes
		noteHandler := handler.NewNoteHandler()
		notes := v1.Group("/notes")
		{
			notes.GET("", noteHandler.GetAll)
			notes.GET("/:id", noteHandler.GetByID)
			notes.POST("", noteHandler.Create)
			notes.PUT("/:id", noteHandler.Update)
			notes.DELETE("/:id", noteHandler.Delete)
		}

		// Tasks / Todos
		taskHandler := handler.NewTaskHandler()
		tasks := v1.Group("/tasks")
		{
			tasks.GET("", taskHandler.GetAll)
			tasks.GET("/:id", taskHandler.GetByID)
			tasks.POST("", taskHandler.Create)
			tasks.PUT("/:id", taskHandler.Update)
			tasks.DELETE("/:id", taskHandler.Delete)
		}
		// Alias for todos
		todos := v1.Group("/todos")
		{
			todos.GET("", taskHandler.GetAll)
			todos.GET("/:id", taskHandler.GetByID)
			todos.POST("", taskHandler.Create)
			todos.PUT("/:id", taskHandler.Update)
			todos.DELETE("/:id", taskHandler.Delete)
		}

		// Bookmarks
		bookmarkHandler := handler.NewBookmarkHandler()
		bookmarks := v1.Group("/bookmarks")
		{
			bookmarks.GET("", bookmarkHandler.GetAll)
			bookmarks.GET("/:id", bookmarkHandler.GetByID)
			bookmarks.POST("", bookmarkHandler.Create)
			bookmarks.PUT("/:id", bookmarkHandler.Update)
			bookmarks.DELETE("/:id", bookmarkHandler.Delete)
		}

		// Events
		eventHandler := handler.NewEventHandler()
		events := v1.Group("/events")
		{
			events.GET("", eventHandler.GetAllEvents)
			events.GET("/:id", eventHandler.GetEventByID)
			events.POST("", eventHandler.CreateEvent)
			events.PUT("/:id", eventHandler.UpdateEvent)
			events.DELETE("/:id", eventHandler.DeleteEvent)
		}

		// Files
		fileHandler := handler.NewFileHandler()
		files := v1.Group("/files")
		{
			files.GET("", fileHandler.GetAll)
			files.GET("/:id", fileHandler.GetByID)
			files.GET("/category/:category", fileHandler.GetByCategory)
			files.POST("/upload", fileHandler.Upload)
			files.GET("/download/:id", fileHandler.Download)
			files.PUT("/:id/rename", fileHandler.Rename)
			files.DELETE("/:id", fileHandler.Delete)
		}

		// Theme
		themeHandler := handler.NewThemeHandler()
		theme := v1.Group("/theme")
		{
			theme.GET("", themeHandler.Get)
			theme.PUT("", themeHandler.Update)
		}

		// Chat
		chatHandler := handler.NewChatHandler()
		chat := v1.Group("/chat")
		{
			chat.GET("/history", chatHandler.GetHistory)
			chat.POST("/message", chatHandler.SendMessage)
			chat.DELETE("/history", chatHandler.ClearHistory)
		}

		// Collection
		collectionHandler := handler.NewCollectionHandler()
		collection := v1.Group("/collections")
		{
			collection.GET("", collectionHandler.GetAllCollections)
			collection.GET("/:id", collectionHandler.GetCollectionByID)
			collection.POST("", collectionHandler.CreateCollection)
			collection.POST("/parse", collectionHandler.ParseURLHandler)
			collection.PUT("/:id", collectionHandler.UpdateCollection)
			collection.DELETE("/:id", collectionHandler.DeleteCollection)
		}

		// Code Arena
		codeHandler := handler.NewCodeHandler()
		code := v1.Group("/code")
		{
			code.POST("/run", codeHandler.RunCode)
		}

		// Blog
		blogHandler := handler.NewBlogHandler(database.DB)
		blog := v1.Group("/blog")
		{
			blog.GET("", blogHandler.GetAllPosts)
			blog.POST("", blogHandler.CreatePost)
			blog.PUT("/:id", blogHandler.UpdatePost)
			blog.DELETE("/:id", blogHandler.DeletePost)
		}
		// Monitor
		monitorHandler := handler.NewMonitorHandler()
		v1.POST("/monitor/check", monitorHandler.CheckHealth)

		// RSS
		rssHandler := handler.NewRSSHandler()
		v1.POST("/rss/feed", rssHandler.GetFeed)
	}

	return r
}
