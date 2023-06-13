package clusterPods

import (
    "context"
    "log"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/api/core/v1"
)

func GetAllClusterPods() []v1.Pod {
    kubeconfig := "/root/dev_volume/golang_crash_course/k3s.yaml"
    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    if err != nil {
        log.Fatalf("Error building kubeconfig: %v", err)
    }

    // Create the Kubernetes clientset
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.Fatalf("Error creating clientset: %v", err)
    }

    // List all pods in all namespaces
    podList, err := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
    if err != nil {
        log.Fatalf("Error listing pods: %v", err)
    }
    return podList.Items
}

