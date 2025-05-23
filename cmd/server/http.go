package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PluginParameters struct {
	LabelSelector metav1.LabelSelector `json:"labelSelector,omitempty"`
}

type PluginInput struct {
	Parameters *PluginParameters `json:"parameters,omitempty"`
}

type ServiceRequest struct {
	ApplicationSetName *string      `json:"applicationSetName,omitempty"`
	Input              *PluginInput `json:"input,omitempty"`
}

type ResponseParameters struct {
	Namespace *string `json:"namespace,omitempty"`
}

type ResponseOutput struct {
	Parameters []*ResponseParameters `json:"parameters,omitempty"`
}

type ResponseBody struct {
	Output *ResponseOutput `json:"output,omitempty"`
}

func (c *ServerConfig) secretsHandler(ctx context.Context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Received request", "address", r.RemoteAddr, "method", r.Method, "url", r.URL.String(), "content-type", r.Header.Get("Content-Type"))
		if r.Method != http.MethodPost {
			slog.Debug("Method not allowed", "method", r.Method, "address", r.RemoteAddr, "url", r.URL.String())
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte("Method not allowed"))
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			slog.Debug("Unsupported media type", "media-type", r.Header.Get("Content-Type"), "address", r.RemoteAddr, "url", r.URL.String())
			w.WriteHeader(http.StatusUnsupportedMediaType)
			_, _ = w.Write([]byte("Unsupported media type"))
			return
		}
		if c.ListenToken != "" && r.Header.Get("Authorization") != "Bearer "+c.ListenToken {
			slog.Debug("Unauthorized", "address", r.RemoteAddr, "url", r.URL.String())
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Unauthorized"))
			return
		}

		input := ServiceRequest{}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			slog.Debug("Unable to read input json", "address", r.RemoteAddr, "error", err)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Bad request"))
			return
		}

		if input.Input == nil {
			slog.Debug("No input provided", "address", r.RemoteAddr)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Bad request"))
			return
		}

		if input.Input.Parameters == nil {
			slog.Debug("No parameters provided", "address", r.RemoteAddr)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Bad request"))
			return
		}

		slog.Debug("Received input", "input", input.Input.Parameters, "address", r.RemoteAddr)

		k8sClient, err := c.GetClient(input.Input.Parameters)
		if err != nil {
			slog.Error("Failed to get k8s client", "address", r.RemoteAddr, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Internal server error"))
			return
		}

		listOptions := metav1.ListOptions{}
		if input.Input.Parameters != nil {
			labelSelector := input.Input.Parameters.LabelSelector
			if labelSelector.MatchLabels != nil || labelSelector.MatchExpressions != nil {
				listOptions.LabelSelector = metav1.FormatLabelSelector(&labelSelector)
				slog.Debug("Using label selector", "labelSelector", listOptions.LabelSelector, "address", r.RemoteAddr)
			}
		}

		namespaces, err := k8sClient.CoreV1().Namespaces().List(ctx, listOptions)
		if err != nil {
			slog.Error("Failed to list namespaces", "address", r.RemoteAddr, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Internal server error"))
			return
		}

		output := ResponseBody{
			Output: &ResponseOutput{
				Parameters: []*ResponseParameters{},
			},
		}

		for _, ns := range namespaces.Items {
			output.Output.Parameters = append(output.Output.Parameters, &ResponseParameters{
				Namespace: &ns.Name,
			})
		}

		slog.Debug("Returning response", "address", r.RemoteAddr, "output", output)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(output); err != nil {
			slog.Error("Failed to encode response", "address", r.RemoteAddr, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Internal server error"))
		}
	}
}
