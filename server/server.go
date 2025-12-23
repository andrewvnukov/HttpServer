package server

import (
	"log"
	"net/http"
	"restapi/handler"
	"restapi/middleware"
	"restapi/utils"

	_ "restapi/docs" // Импорт сгенерированной документации

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	apikey = "12345"
)

// Server структура HTTP сервера
// @Description Основной сервер приложения с маршрутизацией и middleware
type Server struct {
	port     string
	router   *mux.Router
	handlers handler.HandlerManager
}

// NewServer создает новый экземпляр сервера
// @Summary Создать новый сервер
// @Description Инициализирует новый HTTP сервер с указанным портом
// @Param port query string false "Порт для запуска сервера" default(8080)
// @Return *Server новый экземпляр сервера
func NewServer(port string) *Server {
	if port == "" {
		port = ":8080"
	}
	if port[0] != ':' {
		port = ":" + port
	}
	return &Server{
		port:     port,
		router:   mux.NewRouter(),
		handlers: handler.NewHandlerManager(),
	}
}

// Init инициализирует маршруты и middleware сервера
// @Summary Инициализировать сервер
// @Description Настраивает все маршруты API, middleware и Swagger документацию
func (s Server) Init() {
	s.router.NotFoundHandler = utils.ErrNotFoundApi

	// Swagger UI документация
	s.router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("#swagger-ui"),
		httpSwagger.UIConfig(map[string]string{
			"displayRequestDuration": "true",
			"filter":                 "true",
		}),
		httpSwagger.PersistAuthorization(true),
	))

	// Serve swagger.json напрямую
	s.router.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	// Стартовая страница
	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>Library REST API</title>
				<style>
					body { font-family: Arial, sans-serif; margin: 40px; line-height: 1.6; }
					h1 { color: #2c3e50; border-bottom: 2px solid #3498db; padding-bottom: 10px; }
					h2 { color: #34495e; margin-top: 30px; }
					.card { 
						background: #f8f9fa; 
						padding: 20px; 
						margin: 20px 0; 
						border-radius: 8px;
						border-left: 4px solid #3498db;
						box-shadow: 0 2px 4px rgba(0,0,0,0.1);
					}
					.api-version { 
						display: inline-block; 
						background: #e74c3c; 
						color: white; 
						padding: 3px 8px; 
						border-radius: 4px; 
						font-size: 12px; 
						margin-left: 10px; 
						font-weight: bold;
					}
					.v1 { background: #3498db; }
					.v2 { background: #2ecc71; }
					.endpoint { 
						background: white; 
						padding: 12px; 
						margin: 8px 0; 
						border-radius: 4px;
						border: 1px solid #ddd;
						font-family: monospace;
					}
					.method { 
						display: inline-block; 
						padding: 3px 8px; 
						border-radius: 3px; 
						font-weight: bold; 
						margin-right: 10px; 
						font-size: 12px;
						color: white;
					}
					.get { background: #2ecc71; }
					.post { background: #3498db; }
					.put { background: #f39c12; }
					.delete { background: #e74c3c; }
					a { 
						color: #2980b9; 
						text-decoration: none; 
						font-weight: bold;
					}
					a:hover { text-decoration: underline; color: #1a5276; }
					.api-key { 
						background: #fff3cd; 
						padding: 10px; 
						border-radius: 4px; 
						border: 1px solid #ffeaa7;
						margin: 10px 0;
						font-family: monospace;
					}
				</style>
			</head>
			<body>
				<h1>📚 Library Management REST API</h1>
				
				<div class="card">
					<h2>🚀 Quick Start</h2>
					<p>Это REST API для управления библиотекой книг, пользователями и историей покупок.</p>
					<p><strong>Base URL:</strong> <code>http://localhost` + s.port + `/api</code></p>
					<p><strong>API Key:</strong> <span class="api-key">12345</span> (используйте в заголовке X-API-Key)</p>
				</div>
				
				<div class="card">
					<h2>📖 Documentation</h2>
					<ul>
						<li><a href="/swagger/" target="_blank">📚 Swagger UI Documentation</a> - интерактивная документация</li>
						<li><a href="/swagger/doc.json" target="_blank">📄 Swagger JSON</a> - сырая документация в JSON</li>
					</ul>
				</div>
				
				<div class="card">
					<h2>🔐 Authentication</h2>
					<p>Все API endpoints требуют API ключ в заголовке:</p>
					<div class="endpoint">
						<strong>Header:</strong> X-API-Key: 12345
					</div>
					<p>Исключение: Swagger UI и главная страница не требуют аутентификации.</p>
				</div>
				
				<div class="card">
					<h2>API v1 <span class="api-version v1">v1</span></h2>
					<p><strong>Base URL:</strong> <code>/api/v1</code></p>
					
					<div class="endpoint">
						<span class="method get">GET</span> <strong>/users</strong> - получить всех пользователей
					</div>
					<div class="endpoint">
						<span class="method get">GET</span> <strong>/users/{id}</strong> - получить пользователя по ID
					</div>
					<div class="endpoint">
						<span class="method post">POST</span> <strong>/users/{action}</strong> - действия с пользователями (add/update)
					</div>
					
					<div class="endpoint">
						<span class="method get">GET</span> <strong>/books</strong> - получить все книги
					</div>
					<div class="endpoint">
						<span class="method get">GET</span> <strong>/books/{id}</strong> - получить книгу по ID
					</div>
					<div class="endpoint">
						<span class="method post">POST</span> <strong>/books/{action}</strong> - действия с книгами (add/update)
					</div>
					
					<div class="endpoint">
						<span class="method get">GET</span> <strong>/story</strong> - получить всю историю покупок
					</div>
					<div class="endpoint">
						<span class="method post">POST</span> <strong>/story</strong> - добавить покупку
					</div>
					<div class="endpoint">
						<span class="method get">GET</span> <span class="method put">PUT</span> <strong>/story/{action}/{id}</strong> - работа с конкретной покупкой
					</div>
				</div>
				
				<div class="card">
					<h2>API v2 <span class="api-version v2">v2</span></h2>
					<p><strong>Base URL:</strong> <code>/api/v2</code></p>
					<p>Версия v2 включает все возможности v1 плюс дополнительные операции удаления:</p>
					
					<div class="endpoint">
						<span class="method get">GET</span> <span class="method delete">DELETE</span> <strong>/users/{id}</strong> - получить/удалить пользователя
					</div>
					<div class="endpoint">
						<span class="method get">GET</span> <span class="method delete">DELETE</span> <strong>/books/{id}</strong> - получить/удалить книгу
					</div>
					<div class="endpoint">
						<span class="method get">GET</span> <span class="method put">PUT</span> <span class="method delete">DELETE</span> <strong>/story/{action}/{id}</strong> - полный CRUD для покупок
					</div>
				</div>
				
				<div class="card">
					<h2>📞 Примеры запросов</h2>
					<div class="endpoint">
						<strong>curl -X GET "http://localhost` + s.port + `/api/v1/books" -H "X-API-Key: 12345"</strong>
					</div>
					<div class="endpoint">
						<strong>curl -X POST "http://localhost` + s.port + `/api/v1/users/add" -d "name=John&surname=Doe" -H "X-API-Key: 12345"</strong>
					</div>
				</div>
				
				<div style="margin-top: 40px; padding-top: 20px; border-top: 1px solid #ddd; color: #7f8c8d; font-size: 14px;">
					<p>© 2024 Library Management API. Версия 1.0.0</p>
					<p>Все запросы к API должны содержать валидный API ключ в заголовке X-API-Key</p>
				</div>
			</body>
			</html>
		`))
	})

	var api = s.router.PathPrefix("/api").Subrouter()
	api.NotFoundHandler = utils.ErrNotFoundApi

	// Middleware для проверки API ключа
	api.Use(middleware.APIKeyMiddleware(apikey))

	// API Version 1
	var v1 = api.PathPrefix("/v1").Subrouter()
	v1.NotFoundHandler = utils.ErrNotFoundApi
	v1.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"version": "1.0", "status": "active", "message": "API v1 is running"}`))
	})

	{
		// Users endpoints v1
		v1.Handle("/users/{id}", s.handlers["users"]).Methods("GET")
		v1.Handle("/users/{action}", s.handlers["users"]).Methods("POST")

		// Books endpoints v1
		v1.Handle("/books/{id}", s.handlers["books"]).Methods("GET")
		v1.Handle("/books/{action}", s.handlers["books"]).Methods("POST")

		// Story endpoints v1
		v1.Handle("/story/{action}/{id}", s.handlers["story"]).Methods("GET", "PUT")

		// Collection endpoints v1
		v1.Handle("/story", s.handlers["story"]).Methods("GET", "POST")
		v1.Handle("/users", s.handlers["users"]).Methods("GET")
		v1.Handle("/books", s.handlers["books"]).Methods("GET")
	}

	// API Version 2
	var v2 = api.PathPrefix("/v2").Subrouter()
	v2.NotFoundHandler = utils.ErrNotFoundApi
	v2.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"version": "2.0", "status": "active", "message": "API v2 is running", "features": ["delete_operations"]}`))
	})
	{
		// Users endpoints v2
		v2.Handle("/users/{id}", s.handlers["users"]).Methods("GET", "DELETE")
		v2.Handle("/users/{action}", s.handlers["users"]).Methods("POST")

		// Books endpoints v2
		v2.Handle("/books/{id}", s.handlers["books"]).Methods("GET", "DELETE")
		v2.Handle("/books/{action}", s.handlers["books"]).Methods("POST")

		// Story endpoints v2
		v2.Handle("/story/{action}/{id}", s.handlers["story"]).Methods("GET", "PUT", "DELETE")

		// Collection endpoints v2
		v2.Handle("/story", s.handlers["story"]).Methods("GET", "POST")
		v2.Handle("/users", s.handlers["users"]).Methods("GET")
		v2.Handle("/books", s.handlers["books"]).Methods("GET")
	}
}

// StartServer запускает HTTP сервер
// @Summary Запустить сервер
// @Description Запускает HTTP сервер на указанном порту
func (s *Server) StartServer() {
	log.Printf("🚀 Server starting on http://localhost%s", s.port)
	log.Printf("📖 Swagger UI: http://localhost%s/swagger/", s.port)
	log.Printf("🔐 API Key required: %s", apikey)
	log.Printf("🌐 API v1: http://localhost%s/api/v1", s.port)
	log.Printf("🌐 API v2: http://localhost%s/api/v2", s.port)

	err := http.ListenAndServe(s.port, s.router)
	if err != nil {
		log.Fatalf("❌ Server error: %s", err)
	}
}
