package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type PodRequest struct {
	Name           string `json:"name"`
	Type           string `json:"type"` // claude, codex, gemini, bob
	Repo           string `json:"repo"` // optional override
	APIKey         string `json:"apiKey"`
	GoogleAuthJSON string `json:"googleAuthJSON"`
	UseVertexAI    bool   `json:"useVertexAI"`
	UseGCA         bool   `json:"useGCA"`
	MountKeys      bool   `json:"mountKeys"`
}

var clientset *kubernetes.Clientset
var namespace string
var defaultRepo = "ai-cli-suite"

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Error getting in-cluster config: %v", err)
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating clientset: %v", err)
	}

	namespace = os.Getenv("NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}

	if repo := os.Getenv("DEFAULT_REPO"); repo != "" {
		defaultRepo = repo
	}

	http.HandleFunc("/pods", podsHandler)
	http.HandleFunc("/pods/", podDetailHandler)

	fmt.Println("Manager server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func podsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request for /pods", r.Method)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req PodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Creating pod: Name=%s, Type=%s, MountKeys=%v", req.Name, req.Type, req.MountKeys)

	image := fmt.Sprintf("%s/%s-cli:latest", defaultRepo, req.Type)
	if req.Repo != "" {
		image = fmt.Sprintf("%s/%s-cli:latest", req.Repo, req.Type)
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: req.Name,
			Labels: map[string]string{
				"app": "ai-cli",
				"cli": req.Type,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "cli",
					Image:           image,
					ImagePullPolicy: corev1.PullIfNotPresent,
					TTY:             true,
					Stdin:           true,
				},
			},
			RestartPolicy: corev1.RestartPolicyAlways,
		},
	}

	if req.MountKeys {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Printf("Warning: Could not determine home dir for SSH mount: %v", err)
			homeDir = "/root"
		}
		sshPath := homeDir + "/.ssh"

		pod.Spec.Volumes = append(pod.Spec.Volumes, corev1.Volume{
			Name: "ssh-keys",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: sshPath,
				},
			},
		})

		pod.Spec.Containers[0].VolumeMounts = append(pod.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
			Name:      "ssh-keys",
			MountPath: "/root/.ssh",
			ReadOnly:  true,
		})
		log.Printf("Mounted %s to /root/.ssh (readonly)", sshPath)
	}

	if req.APIKey != "" {
		envVarName := ""
		switch req.Type {
		case "gemini":
			envVarName = "GEMINI_API_KEY"
		case "claude":
			envVarName = "ANTHROPIC_API_KEY"
		case "codex":
			envVarName = "OPENAI_API_KEY"
		case "bob":
			envVarName = "BOBSHELL_API_KEY"
		}

		if envVarName != "" {
			pod.Spec.Containers[0].Env = append(pod.Spec.Containers[0].Env, corev1.EnvVar{
				Name:  envVarName,
				Value: req.APIKey,
			})
			log.Printf("Injected API key env var: %s", envVarName)
		}
	}

	if req.Type == "gemini" {
		// Trust the workspace for automated environments
		pod.Spec.Containers[0].Env = append(pod.Spec.Containers[0].Env, corev1.EnvVar{
			Name:  "GEMINI_CLI_TRUST_WORKSPACE",
			Value: "true",
		})
		log.Printf("Enabled GEMINI_CLI_TRUST_WORKSPACE=true")

		if req.UseVertexAI {
			pod.Spec.Containers[0].Env = append(pod.Spec.Containers[0].Env, corev1.EnvVar{
				Name:  "GOOGLE_GENAI_USE_VERTEXAI",
				Value: "true",
			})
			log.Printf("Injected GOOGLE_GENAI_USE_VERTEXAI=true")
		}
		if req.UseGCA {
			pod.Spec.Containers[0].Env = append(pod.Spec.Containers[0].Env, corev1.EnvVar{
				Name:  "GOOGLE_GENAI_USE_GCA",
				Value: "true",
			})
			log.Printf("Injected GOOGLE_GENAI_USE_GCA=true")
		}
		if req.GoogleAuthJSON != "" {
			pod.Spec.Containers[0].Env = append(pod.Spec.Containers[0].Env, corev1.EnvVar{
				Name:  "GOOGLE_AUTH_JSON",
				Value: req.GoogleAuthJSON,
			})
			log.Printf("Injected GOOGLE_AUTH_JSON")
		}
	}

	_, err := clientset.CoreV1().Pods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		log.Printf("Failed to create pod: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create pod: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Pod %s created successfully", req.Name)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Pod created successfully", "name": req.Name})
}

func podDetailHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/pods/"), "/")
	name := parts[0]
	log.Printf("Received %s request for /pods/%s", r.Method, name)

	switch r.Method {
	case http.MethodDelete:
		err := clientset.CoreV1().Pods(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
		if err != nil {
			log.Printf("Failed to delete pod %s: %v", name, err)
			http.Error(w, fmt.Sprintf("Failed to delete pod: %v", err), http.StatusInternalServerError)
			return
		}
		log.Printf("Pod %s deleted successfully", name)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Pod deleted successfully"})

	case http.MethodPut:
		if len(parts) > 1 && parts[1] == "restart" {
			log.Printf("Restarting pod %s", name)
			restartPod(w, r, name)
		} else {
			log.Printf("Updating pod %s", name)
			updatePod(w, r, name)
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func restartPod(w http.ResponseWriter, r *http.Request, name string) {
	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		log.Printf("Failed to get pod %s for restart: %v", name, err)
		http.Error(w, fmt.Sprintf("Failed to get pod: %v", err), http.StatusInternalServerError)
		return
	}

	// Simple restart: Delete and recreate
	err = clientset.CoreV1().Pods(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		log.Printf("Failed to delete pod %s for restart: %v", name, err)
		http.Error(w, fmt.Sprintf("Failed to delete pod for restart: %v", err), http.StatusInternalServerError)
		return
	}

	newPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: pod.Labels,
		},
		Spec: pod.Spec,
	}
	newPod.ResourceVersion = "" // Clear for creation

	_, err = clientset.CoreV1().Pods(namespace).Create(context.TODO(), newPod, metav1.CreateOptions{})
	if err != nil {
		log.Printf("Failed to recreate pod %s: %v", name, err)
		http.Error(w, fmt.Sprintf("Failed to recreate pod: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Pod %s restarted successfully", name)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Pod restarted successfully"})
}

func updatePod(w http.ResponseWriter, r *http.Request, name string) {
	var req PodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding update request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		log.Printf("Failed to get pod %s for update: %v", name, err)
		http.Error(w, fmt.Sprintf("Failed to get pod: %v", err), http.StatusInternalServerError)
		return
	}

	image := fmt.Sprintf("%s/%s-cli:latest", defaultRepo, req.Type)
	if req.Repo != "" {
		image = fmt.Sprintf("%s/%s-cli:latest", req.Repo, req.Type)
	}

	pod.Spec.Containers[0].Image = image

	_, err = clientset.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		log.Printf("Failed to update pod %s: %v", name, err)
		http.Error(w, fmt.Sprintf("Failed to update pod: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Pod %s updated successfully", name)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Pod updated successfully"})
}
