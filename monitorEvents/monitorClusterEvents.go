package monitorEvents

import (
    "context"
    "log"
    
    v1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/watch"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

type MonitorEventObject struct {
    Namespace string
}

/*
  Depandants:
    1. Kubernetes namespace
    2. logFileName
  Description:
    This function monitor the Kubernetes events of a namespace and log them to a file.
*/
func (obj MonitorEventObject) WatchEvents() {
    k8sNamespace := obj.Namespace
    //logFilePath := fmt.Sprintf("/mnt/var/log/%s/events.log", k8sNamespace)
    kubeconfig := "/root/dev_volume/golang_crash_course/k3s.yaml"

    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    if err != nil {
        log.Fatalf("Failed to build config: %v", err)
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    eventClient := clientset.CoreV1().Events(k8sNamespace)
    // Capture all pod events of k8sNamespace
    listOptions := metav1.ListOptions{}

    watcher, err := eventClient.Watch(context.TODO(), listOptions)
    if err != nil {
        log.Fatalf("Failed to create watcher: %v", err)
    }

    ch := watcher.ResultChan()
    for event := range ch {
        if event.Type == watch.Error {
            log.Printf("Received error event: %v", event.Object)
                continue
        }

        eventObj, ok := event.Object.(*v1.Event)
        if !ok {
            log.Printf("Failed to convert event object: %v", event.Object)
            continue
        }

        log.Printf("Received event: %s/%s - %s: %s", eventObj.Namespace, eventObj.InvolvedObject.Name, eventObj.Reason, eventObj.Message)
    }
}
