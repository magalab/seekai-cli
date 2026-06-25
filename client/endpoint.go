package client

import (
	"encoding/json"
	"fmt"
)

type EndpointConfig struct {
	AIModelName string `json:"ai_model_name"`
	URL         string `json:"url"`
	AccessKey   string `json:"access_key,omitempty"`
	Provider    string `json:"provider"`
}

type EndpointUpdate struct {
	URL       string `json:"url,omitempty"`
	AccessKey string `json:"access_key,omitempty"`
}

func (c *Client) ListEndpoints() (QueryResult, error) {
	return c.QueryResult("SELECT * FROM oceanbase.DBA_OB_AI_MODEL_ENDPOINTS")
}

func (c *Client) CreateEndpoint(name string, cfg EndpointConfig) error {
	payload, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return c.Exec("CALL DBMS_AI_SERVICE.CREATE_AI_MODEL_ENDPOINT(?, ?)", name, string(payload))
}

func (c *Client) UpdateEndpoint(name string, update EndpointUpdate) error {
	if update.URL == "" && update.AccessKey == "" {
		return fmt.Errorf("at least one of --url or --access-key is required")
	}
	payload, err := json.Marshal(update)
	if err != nil {
		return err
	}
	return c.Exec("CALL DBMS_AI_SERVICE.ALTER_AI_MODEL_ENDPOINT(?, ?)", name, string(payload))
}

func (c *Client) DeleteEndpoint(name string) error {
	return c.Exec("CALL DBMS_AI_SERVICE.DROP_AI_MODEL_ENDPOINT(?)", name)
}
