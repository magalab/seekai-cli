package client

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func (c *Client) Complete(modelKey, prompt, parameters string) (string, time.Duration, error) {
	start := time.Now()
	var row []map[string]any
	var err error
	if parameters == "" {
		row, err = c.Query("SELECT AI_COMPLETE(?, ?) AS result", modelKey, prompt)
	} else {
		row, err = c.Query("SELECT AI_COMPLETE(?, ?, ?) AS result", modelKey, prompt, parameters)
	}
	if err != nil {
		return "", 0, err
	}
	return firstString(row, "result"), time.Since(start), nil
}

func (c *Client) Embed(modelKey, text string, dim int) (any, error) {
	var row []map[string]any
	var err error
	if dim > 0 {
		row, err = c.Query("SELECT AI_EMBED(?, ?, ?) AS embedding", modelKey, text, dim)
	} else {
		row, err = c.Query("SELECT AI_EMBED(?, ?) AS embedding", modelKey, text)
	}
	if err != nil {
		return "", err
	}
	return parseJSONValue(firstString(row, "embedding")), nil
}

func (c *Client) EmbedBatch(modelKey string, texts []string, dim int) ([]any, error) {
	out := make([]any, 0, len(texts))
	for _, text := range texts {
		embedding, err := c.Embed(modelKey, text, dim)
		if err != nil {
			return nil, err
		}
		out = append(out, embedding)
	}
	return out, nil
}

func (c *Client) Rerank(modelKey, query, documentsJSON string) (string, error) {
	row, err := c.Query("SELECT AI_RERANK(?, ?, ?) AS result", modelKey, query, documentsJSON)
	if err != nil {
		return "", err
	}
	return firstString(row, "result"), nil
}

func (c *Client) Prompt(template string, args []string) (string, time.Duration, error) {
	start := time.Now()
	params := append([]string{template}, args...)
	placeholders := strings.TrimRight(strings.Repeat("?,", len(params)), ",")
	values := make([]any, len(params))
	for i, value := range params {
		values[i] = value
	}
	row, err := c.Query(fmt.Sprintf("SELECT AI_PROMPT(%s) AS result", placeholders), values...)
	if err != nil {
		return "", 0, err
	}
	return firstString(row, "result"), time.Since(start), nil
}

func firstString(rows []map[string]any, key string) string {
	if len(rows) == 0 {
		return ""
	}
	value, ok := rows[0][key]
	if !ok || value == nil {
		return ""
	}
	return fmt.Sprint(value)
}

func parseJSONValue(value string) any {
	var parsed any
	if json.Unmarshal([]byte(value), &parsed) == nil {
		return parsed
	}
	return value
}
