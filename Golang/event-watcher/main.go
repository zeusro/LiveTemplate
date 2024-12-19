package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Notification struct {
	EventType string `json:"event_type"`
	PodName   string `json:"pod_name"`
	Namespace string `json:"namespace"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func getKubernetesClient() (*kubernetes.Clientset, error) {
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = fmt.Sprintf("%s/.kube/config", home)
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func sendNotification(event Notification, webhookURL string) error {
	// Convert the event into JSON format
	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// Send the event to the external service via Webhook
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(eventData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send notification, status: %s", resp.Status)
	}

	return nil
}

func watchKubernetesEvents(clientset *kubernetes.Clientset, namespace string, webhookURL string) {
	watcher, err := clientset.CoreV1().Events(namespace).Watch(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error watching events: %v", err)
	}
	defer watcher.Stop()

	for event := range watcher.ResultChan() {
		// Check if the event is related to Pod creation, deletion, or failure
		if event.Type == "ADDED" || event.Type == "MODIFIED" || event.Type == "DELETED" {
			podName := event.Object.(*metav1.ObjectMeta).Name
			eventMessage := event.Object.(*metav1.ObjectMeta).Annotations["message"]

			// Prepare the Notification struct
			notification := Notification{
				EventType: string(event.Type),
				PodName:   podName,
				Namespace: namespace,
				Message:   eventMessage,
				Timestamp: time.Now().Format(time.RFC3339),
			}

			// Send the notification to an external Webhook
			err := sendNotification(notification, webhookURL)
			if err != nil {
				log.Printf("Failed to send notification for event %s: %v", event.Type, err)
			} else {
				log.Printf("Notification sent for event: %s", event.Type)
			}
		}
	}
}

func main() {
	kubeClient, err := getKubernetesClient()
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	// Define the namespace to watch (e.g., default)
	namespace := "default"

	// External Webhook URL
	webhookURL := "https://example.com/notification-webhook"

	// Watch Kubernetes events and send notifications
	go watchKubernetesEvents(kubeClient, namespace, webhookURL)

	// Set up a simple HTTP server
	r := mux.NewRouter()
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}).Methods("GET")

	http.Handle("/", r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting notification server on port %s", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
