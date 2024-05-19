package main

import (
	"io/fs"
	"net/http"
	"os"

	"github.com/Hari-Kiri/goalApplicationSettingsLoader"
	"github.com/Hari-Kiri/goalJson"
	"github.com/Hari-Kiri/goalMakeHandler"
	"github.com/Hari-Kiri/temboLog"
)

// Function constructor
func main() {
	// Load application settings parameter
	loadApplicationSettings, errorLoadApplicationSettings := goalApplicationSettingsLoader.LoadSettings()
	// If load application settings parameter return error handle it
	if errorLoadApplicationSettings != nil {
		temboLog.FatalLogging("Error opening application settings file", errorLoadApplicationSettings.Error())
		return
	}
	temboLog.InfoLogging("Starting", loadApplicationSettings.Settings.Name)
	temboLog.InfoLogging("Build and provided by:", loadApplicationSettings.Settings.Organisation)
	goalMakeHandler.HandleFileRequest("/", ",/")
	// List directory
	readDirectory, errorReadDir := os.ReadDir("./")
	if errorReadDir != nil {
		temboLog.FatalLogging("failed listing directory, error: ", errorReadDir)
	}
	// Handle directory request
	for i := 0; i < len(readDirectory); i++ {
		directoryHandler(readDirectory[i])
	}
	// Handle test page (its just for testing webserver online or not) request
	goalMakeHandler.HandleRequest(testHandler, "/test")
	// Run HTTP server
	goalMakeHandler.Serve(loadApplicationSettings.Settings.Name, loadApplicationSettings.Settings.Port)
}

func directoryHandler(directoryReadResult fs.DirEntry) {
	if directoryReadResult.IsDir() {
		http.HandleFunc("/"+directoryReadResult.Name()+"/", serveFiles)
	}
}

func serveFiles(responseWriter http.ResponseWriter, request *http.Request) {
	name := "." + request.URL.Path
	temboLog.InfoLogging("serving:", name[1:])
	http.ServeFile(responseWriter, request, name)
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
