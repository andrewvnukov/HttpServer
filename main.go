package main

import (
	"restapi/server"
)

// @title Library Management REST API
// @version 2.0.0
// @description Полное REST API для управления библиотекой с поддержкой версий v1 и v2
//
// ### Возможности:
// - **v1**: Базовые операции (GET, POST, PUT)
// - **v2**: Расширенные операции (+ DELETE)
//
// ### Аутентификация:
// Все запросы требуют API ключ в заголовке X-API-Key
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@libraryapi.com
// @contact.url https://www.libraryapi.com/support
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
// @description API Key Authentication
func main() {
	server := server.NewServer("8080")
	server.Init()
	server.StartServer()
}
