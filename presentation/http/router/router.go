package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nuba55yo/go-101-CleanCRUD/application/usecase"
	v1 "github.com/nuba55yo/go-101-CleanCRUD/presentation/http/v1"
	v2 "github.com/nuba55yo/go-101-CleanCRUD/presentation/http/v2"
	"github.com/nuba55yo/go-101-CleanCRUD/presentation/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(bookUseCase usecase.BookUseCase) *gin.Engine {
	r := gin.New()
	_ = r.SetTrustedProxies(nil)
	r.Use(gin.Recovery(), middleware.AccessLog())

	// -------- v1 --------
	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/books", v1.ListBooks(bookUseCase))
		apiV1.GET("/books/:id", v1.GetBookByID(bookUseCase))
		apiV1.POST("/books", v1.CreateBook(bookUseCase))
		apiV1.PUT("/books/:id", v1.UpdateBook(bookUseCase))
		apiV1.DELETE("/books/:id", v1.DeleteBook(bookUseCase))
	}

	// -------- v2 --------
	apiV2 := r.Group("/api/v2")
	{
		apiV2.GET("/books", v2.ListBooks(bookUseCase))
		apiV2.GET("/books/:id", v2.GetBookByID(bookUseCase))
		apiV2.POST("/books", v2.CreateBook(bookUseCase))
		apiV2.PUT("/books/:id", v2.UpdateBook(bookUseCase))
		apiV2.DELETE("/books/:id", v2.DeleteBook(bookUseCase))
	}

	// -------- docs (???? gen ????) --------
	r.GET("/docs/v1/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.InstanceName("v1")))
	r.GET("/docs/v2/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.InstanceName("v2")))

	// -------- Swagger UI ??? (dropdown v1/v2) --------
	r.GET("/swagger", func(c *gin.Context) {
		html := `<!doctype html>
<html>
<head>
<meta charset="utf-8"/><title>Swagger UI</title>
<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
<style>body{margin:0}</style>
</head>
<body><div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
<script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
<script>
  const params  = new URLSearchParams(location.search);
  const primary = params.get('urls.primaryName') || 'v1';
  window.ui = SwaggerUIBundle({
    dom_id: '#swagger-ui',
    urls: [{url:'/docs/v1/doc.json',name:'v1'},{url:'/docs/v2/doc.json',name:'v2'}],
    'urls.primaryName': primary,
    deepLinking: true,
    presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
    layout: 'StandaloneLayout'
  });
</script>
</body></html>`
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	})

	return r
}
