package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Hari-Kiri/goalApplicationSettingsLoader"
	"github.com/Hari-Kiri/goalJson"
	"github.com/Hari-Kiri/goalMakeHandler"
	"github.com/Hari-Kiri/goalRenderTemplate"
)

// HTML parser
var htmlTemplates = template.Must(template.ParseFiles("./html/index.html"))

// Function constructor
func main() {
	// Load application settings parameter
	log.Println("[info] Starting webserver!")
	loadApplicationSettings, error := goalApplicationSettingsLoader.LoadSettings()
	// If load application settings parameter return error handle it
	if error != nil {
		log.Fatal("[error main()] Error opening application settings file: " + error.Error())
		return
	}
	// Handle web application user interface components request
	goalMakeHandler.HandleFileRequest("/html/", "./html")
	// Handle web root request
	goalMakeHandler.HandleRequest(rootHandler, "/")
	// Handle test page (its just for testing webserver online or not) request
	goalMakeHandler.HandleRequest(testHandler, "/test")
	// Handle pokemon page request
	goalMakeHandler.HandleRequest(indexHandler, "/index")
	// // Handle catching request
	// goalMakeHandler.HandleRequest(catchHandler, "/pokemon/catch")
	// // Handle catching request
	// goalMakeHandler.HandleRequest(getNickNameHandler, "/pokemon/getnickname")
	// // Handle prime request
	// goalMakeHandler.HandleRequest(releaseHandler, "/pokemon/release")
	// Run HTTP server
	goalMakeHandler.Serve(loadApplicationSettings.Settings.Name, loadApplicationSettings.Settings.Port)
}

// Web root handler
func rootHandler(responseWriter http.ResponseWriter, request *http.Request) {
	// Redirect to home page
	http.Redirect(responseWriter, request, "/index", http.StatusFound)
	log.Println("[info] Webroot redirect to url path [ /index ], requested from", request.RemoteAddr)
}

// Test page handler
func testHandler(responseWriter http.ResponseWriter, request *http.Request) {
	// Http ok response
	okResponse, _ := goalJson.JsonEncode(map[string]interface{}{
		"response": true,
		"code":     200,
		"message":  "Go net/http webserver online"},
		false)
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write([]byte(okResponse))
	log.Println("[info] Serving test page [", request.URL.Path, "]")
}

// Index page handler
func indexHandler(responseWriter http.ResponseWriter, request *http.Request) {
	// Load application settings data
	appSettings, error := goalApplicationSettingsLoader.LoadSettings()
	// If load application settings data return error handle it
	if error != nil {
		// Http error response
		errorResponse, _ := goalJson.JsonEncode(map[string]interface{}{
			"response": false,
			"code":     500,
			"message":  "page failed to serve"},
			false)
		// Give 500 response code
		http.Error(responseWriter, errorResponse, http.StatusInternalServerError)
		log.Println("[error adminHandler()] Error opening application settings file: " + error.Error())
		return
	}
	// Open home page
	goalRenderTemplate.Process(htmlTemplates, responseWriter, "index", appSettings, request)
}
