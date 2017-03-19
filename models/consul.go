package models

type ConsulPayload struct {
	HealthChecks []ConsulHealthCheck `json:"health_checks"`
	Environment  string              `json:"environment"`
}

type ConsulHealthCheck struct {
	Node        string `json:"node"`
	CheckID     string `json:"check_id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Notes       string `json:"notes"`
	Output      string `json:"output"`
	ServiceID   string `json:"service_id"`
	ServiceName string `json:"service_name"`
}
