package main

import (
    "fmt"
    "sync"
    "log"
    "encoding/json"
    "net/http"
    "io/ioutil"

    watchEvents "golangclustermonitoring/monitorEvents"
    captureClusterPods "golangclustermonitoring/clusterPods"
    podLogs "golangclustermonitoring/streamPodLogs"
)

var wg sync.WaitGroup

const (
    root        = "/"
    watchevent  = "/watchevent"
    clusterpods = "/clusterpods"
    podlogs     = "/podlogs"
)

func startWatchingEvents(namespace string) {
    var obj watchEvents.MonitorEventObject
    obj = watchEvents.MonitorEventObject {
        Namespace: namespace,
    }
    obj.WatchEvents()
}

func getClusterPods(w http.ResponseWriter) {
    allPods := captureClusterPods.GetAllClusterPods()
    for _, pod := range(allPods) {
        fmt.Fprintln(w, fmt.Sprintf("Pod %s found in namespace %s", pod.Name, pod.Namespace))
    }
}

func streamPodLogs(w http.ResponseWriter, r *http.Request, podName string, namespace string) {
    podLogs.CapturePodLogs(w, r, podName, namespace)
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
           fmt.Fprint(w, "Welcome to GoLang Monitoring Service, ")
           fmt.Fprint(w, "Supported operations are /watchevent, /clusterpods, /podlogs\n")
           break
       case watchevent:
           data := string(body)
           wg.Add(1)
           go func() {
               startWatchingEvents(data)
               wg.Done()
           }()
           break
       case podlogs:
           jsonMap := map[string]string{}
           err := json.Unmarshal(body, &jsonMap)
           if err != nil {
               log.Fatal(err)
           }
           fmt.Fprint(w, fmt.Sprintf("Podname is %s and namespace is %s\n", jsonMap["podname"], jsonMap["namespace"]))
           streamPodLogs(w, r, jsonMap["podname"], jsonMap["namespace"])
       case clusterpods:
           getClusterPods(w)
           break
       default:
           fmt.Fprint(w, "Not supported operation, ")
           fmt.Fprint(w, "Supported operations are /watchevent, /clusterpods, /podlogs\n")
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

