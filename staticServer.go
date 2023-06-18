package main

import (
	"html/template"
	"net/http"

	"github.com/Hari-Kiri/goalApplicationSettingsLoader"
	"github.com/Hari-Kiri/goalJson"
	"github.com/Hari-Kiri/goalMakeHandler"
	"github.com/Hari-Kiri/goalRenderTemplate"
	"github.com/Hari-Kiri/temboLog"
)

// HTML parser
var htmlTemplates = template.Must(template.ParseFiles("./html/index.html"))

// var uriName string

// Function constructor
func main() {
	// Load application settings parameter
	loadApplicationSettings, errorLoadApplicationSettings := goalApplicationSettingsLoader.LoadSettings()
	// If load application settings parameter return error handle it
	if errorLoadApplicationSettings != nil {
		temboLog.FatalLogging("Error opening application settings file", errorLoadApplicationSettings.Error())
		return
	}
	// Load static server settings parameter
	// loadStaticServerSettings, errorLoadStaticServerSettings := goalApplicationSettingsLoader.LoadStaticServerConfiguration()
	// // If load static server settings parameter return error handle it
	// if errorLoadStaticServerSettings != nil {
	// 	temboLog.FatalLogging("Error opening application settings file", errorLoadStaticServerSettings.Error())
	// 	return
	// }
	// uriName = loadStaticServerSettings.StaticServer.Uri
	temboLog.InfoLogging("Starting", loadApplicationSettings.Settings.Name, "version", loadApplicationSettings.Settings.Version)
	temboLog.InfoLogging("Build and provided by:", loadApplicationSettings.Settings.Organisation)
	// Handle web application user interface components request
	goalMakeHandler.HandleFileRequest("/html/", "./html")
	// Handle .well-known folder
	goalMakeHandler.HandleFileRequest("/.well-known", "./html/.well-known")
	// Handle web root request
	goalMakeHandler.HandleRequest(indexHandler, "/")
	// Handle test page (its just for testing webserver online or not) request
	goalMakeHandler.HandleRequest(testHandler, "/test")
	// Run HTTP server
	goalMakeHandler.Serve(loadApplicationSettings.Settings.Name, loadApplicationSettings.Settings.Port)
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
		temboLog.ErrorLogging("Error opening application settings file:", error.Error())
		return
	}
	// Open home page
	goalRenderTemplate.Process(htmlTemplates, responseWriter, "index", appSettings, request)
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
	temboLog.InfoLogging("Serving test page [", request.URL.Path, "]")
}
