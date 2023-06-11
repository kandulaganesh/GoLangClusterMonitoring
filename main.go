package main

import (
    "fmt"
    "sync"
    "log"
    "net/http"
    "io/ioutil"

    watchEvents "golangclustermonitoring/monitorEvents"
)

var wg sync.WaitGroup

const (
    root = "/"
    watchevent = "/watchevent"
    clusterstatus = "/clusterstatus"
)

func startWatchingEvents(namespace string) {
    var obj watchEvents.MonitorEventObject
    obj = watchEvents.MonitorEventObject {
        Namespace: namespace,
    }
    obj.WatchEvents()
}

func handler(w http.ResponseWriter, r *http.Request) {
   url := r.URL.String()
   body, err := ioutil.ReadAll(r.Body)
   if err != nil {
       http.Error(w, "Error reading request body", http.StatusInternalServerError)
       fmt.Fprint(w, "Failed to read input, retry")
       return
   }
   switch(url) {
       case root:
           fmt.Fprint(w, "Welcome to GoLang Monitoring Service")
           fmt.Fprint(w, "Supported operations are /watchevent, /clusterstatus")
           break
       case watchevent:
           data := string(body)
           wg.Add(1)
           go func() {
               startWatchingEvents(data)
               wg.Done()
           }()
           break
       default:
           fmt.Fprint(w, "Not supported operation")
           fmt.Fprint(w, "Supported operations are /watchevent, /clusterstatus")
   }
}

func startServer(port int) {

    // Register the handler function to handle requests.
    http.HandleFunc("/", handler)

    // Start the HTTP server
    listenSocket := fmt.Sprintf(":%d", port)
    log.Printf("Server listening on %s...", listenSocket)
    err := http.ListenAndServe(listenSocket, nil)
    if err != nil {
        log.Fatal(err)
    }    
}

func main() {
    fmt.Println("Starting Application")

    wg.Add(1)
    go func() {
        startServer(8081)
    }()

    wg.Wait()
    fmt.Println("Exiting application")
}

