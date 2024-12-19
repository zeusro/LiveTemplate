package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	admissionv1 "k8s.io/api/admission/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Mutating webhook handler 在创建 Pod 时，自动给 Pod 的第一个容器添加一个环境变量。
func mutatePod(w http.ResponseWriter, r *http.Request) {
	var admissionReview admissionv1.AdmissionReview
	err := json.NewDecoder(r.Body).Decode(&admissionReview)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not read the admission review: %v", err), http.StatusBadRequest)
		return
	}

	// Get the Pod object from the admission request
	pod := v1.Pod{}
	err = json.Unmarshal(admissionReview.Request.Object.Raw, &pod)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not unmarshal pod object: %v", err), http.StatusInternalServerError)
		return
	}

	// Add an environment variable to the first container
	if len(pod.Spec.Containers) > 0 {
		container := &pod.Spec.Containers[0]
		container.Env = append(container.Env, v1.EnvVar{Name: "MY_ENV_VAR", Value: "mutated_value"})
	}

	// Create a patch to apply the changes to the Pod
	patch := []map[string]interface{}{
		{
			"op":    "add",
			"path":  "/spec/containers/0/env/-",
			"value": map[string]interface{}{"name": "MY_ENV_VAR", "value": "mutated_value"},
		},
	}

	patchBytes, err := json.Marshal(patch)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not marshal patch: %v", err), http.StatusInternalServerError)
		return
	}

	// Create the AdmissionResponse
	admissionResponse := admissionv1.AdmissionResponse{
		Allowed:   true,
		UID:       admissionReview.Request.UID,
		PatchType: func() *admissionv1.PatchType { pt := admissionv1.PatchTypeJSONPatch; return &pt }(),
		Patch:     patchBytes,
	}

	// Send the response
	admissionReviewResponse := admissionv1.AdmissionReview{
		Response: &admissionResponse,
	}

	responseBytes, err := json.Marshal(admissionReviewResponse)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not marshal AdmissionReview response: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}

// Validating webhook handler 验证 Pod 的环境变量，如果环境变量 MY_ENV_VAR 的值为 "invalid"，则拒绝创建 Pod。
func validatePod(w http.ResponseWriter, r *http.Request) {
	var admissionReview admissionv1.AdmissionReview
	err := json.NewDecoder(r.Body).Decode(&admissionReview)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not read the admission review: %v", err), http.StatusBadRequest)
		return
	}

	// Get the Pod object from the admission request
	pod := v1.Pod{}
	err = json.Unmarshal(admissionReview.Request.Object.Raw, &pod)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not unmarshal pod object: %v", err), http.StatusInternalServerError)
		return
	}

	// Check if the environment variable "MY_ENV_VAR" exists and its value is "invalid"
	var response admissionv1.AdmissionResponse
	allowed := true
	for _, container := range pod.Spec.Containers {
		for _, env := range container.Env {
			if env.Name == "MY_ENV_VAR" && env.Value == "invalid" {
				allowed = false
				break
			}
		}
	}

	if !allowed {
		response = admissionv1.AdmissionResponse{
			Allowed: false,
			UID:     admissionReview.Request.UID,
			Result: &metav1.Status{
				Message: "Pod creation failed due to invalid MY_ENV_VAR value",
			},
		}
	} else {
		response = admissionv1.AdmissionResponse{
			Allowed: true,
			UID:     admissionReview.Request.UID,
		}
	}

	// Send the response
	admissionReviewResponse := admissionv1.AdmissionReview{
		Response: &response,
	}

	responseBytes, err := json.Marshal(admissionReviewResponse)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not marshal AdmissionReview response: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}

func main() {
	r := mux.NewRouter()

	// Mutating webhook route
	r.HandleFunc("/mutate", mutatePod).Methods("POST")

	// Validating webhook route
	r.HandleFunc("/validate", validatePod).Methods("POST")

	// Run the server
	fmt.Println("Starting webhook server on https://localhost:443")
	err := http.ListenAndServeTLS(":443", "/path/to/cert.crt", "/path/to/cert.key", r)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
