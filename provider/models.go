package provider

import "sort"

type Provider struct {
	Name string
	URLs map[string]string
}

var known = []Provider{
	{
		Name: "aliyun-dashscope",
		URLs: map[string]string{
			"completion":      "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation",
			"dense_embedding": "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding",
			"rerank":          "https://dashscope.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank",
		},
	},
	{
		Name: "aliyun-openAI",
		URLs: map[string]string{
			"completion":      "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions",
			"dense_embedding": "https://dashscope.aliyuncs.com/compatible-mode/v1/embeddings",
		},
	},
	{
		Name: "deepseek",
		URLs: map[string]string{
			"completion": "https://api.deepseek.com/chat/completions",
		},
	},
	{
		Name: "hunyuan-openAI",
		URLs: map[string]string{
			"completion":      "https://api.hunyuan.cloud.tencent.com/v1/chat/completions",
			"dense_embedding": "https://api.hunyuan.cloud.tencent.com/v1/embeddings",
		},
	},
	{
		Name: "openAI",
		URLs: map[string]string{
			"completion":      "https://api.openai.com/v1/chat/completions",
			"dense_embedding": "https://api.openai.com/v1/embeddings",
		},
	},
	{
		Name: "siliconflow",
		URLs: map[string]string{
			"completion":      "https://api.siliconflow.cn/v1/chat/completions",
			"dense_embedding": "https://api.siliconflow.cn/v1/embeddings",
			"rerank":          "https://api.siliconflow.cn/v1/rerank",
		},
	},
}

func Names() []string {
	names := make([]string, 0, len(known))
	for _, p := range known {
		names = append(names, p.Name)
	}
	sort.Strings(names)
	return names
}

func DefaultURL(name, modelType string) string {
	for _, p := range known {
		if p.Name == name {
			return p.URLs[modelType]
		}
	}
	return ""
}
