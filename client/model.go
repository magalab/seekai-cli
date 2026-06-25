package client

import "encoding/json"

type ModelConfig struct {
	Type      string `json:"type"`
	ModelName string `json:"model_name"`
}

func (c *Client) ListModels() (QueryResult, error) {
	return c.QueryResult("SELECT * FROM oceanbase.DBA_OB_AI_MODELS")
}

func (c *Client) CreateModel(name, modelType, modelName string) error {
	payload, err := json.Marshal(ModelConfig{Type: modelType, ModelName: modelName})
	if err != nil {
		return err
	}
	return c.Exec("CALL DBMS_AI_SERVICE.CREATE_AI_MODEL(?, ?)", name, string(payload))
}

func (c *Client) DeleteModel(name string) error {
	return c.Exec("CALL DBMS_AI_SERVICE.DROP_AI_MODEL(?)", name)
}
