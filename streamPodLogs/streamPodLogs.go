package streamPodLogs

import (
    "fmt"
    "log"
    "net/http"
    "time"

    corev1 "k8s.io/api/core/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

func CapturePodLogs(w http.ResponseWriter, r *http.Request, podName string, namespace string) {
    // Specify the kubeconfig file path
    kubeconfigPath := "/root/dev_volume/golang_crash_course/k3s.yaml"

    // Build the configuration from the provided kubeconfig file
    config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
    if err != nil {
        log.Fatalf("Error building kubeconfig: %v", err)
    }

    // Create the Kubernetes clientset
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.Fatalf("Error creating clientset: %v", err)
    }

    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
        return
    }

        w.Header().Set("Content-Type", "text/event-stream")
        w.Header().Set("Cache-Control", "no-cache")
        w.Header().Set("Connection", "keep-alive")
        w.Header().Set("Transfer-Encoding", "chunked")

        // Retrieve the live logs of the pod
        stream, err := clientset.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{
            Follow: true,
        }).Stream(r.Context())
        if err != nil {
            log.Printf("Error retrieving pod logs: %v", err)
            http.Error(w, "Error retrieving logs", http.StatusInternalServerError)
            return
        }
        defer stream.Close()

        buf := make([]byte, 4096)
        for {
            bytesRead, err := stream.Read(buf)
            if err != nil {
                log.Printf("Error reading pod logs: %v", err)
                return
            }

            if bytesRead > 0 {
                fmt.Fprintf(w, "data: %s", string(buf[:bytesRead]))
                flusher.Flush()
            }

            // Add a small delay to control the streaming rate
            time.Sleep(100 * time.Millisecond)
        }
}

