package provider

import "testing"

func TestDefaultURL(t *testing.T) {
	tests := []struct {
		name      string
		provider  string
		modelType string
		want      string
	}{
		{
			name:      "aliyun openai completion",
			provider:  "aliyun-openAI",
			modelType: "completion",
			want:      "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions",
		},
		{
			name:      "aliyun openai embedding",
			provider:  "aliyun-openAI",
			modelType: "dense_embedding",
			want:      "https://dashscope.aliyuncs.com/compatible-mode/v1/embeddings",
		},
		{
			name:      "aliyun dashscope completion",
			provider:  "aliyun-dashscope",
			modelType: "completion",
			want:      "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation",
		},
		{
			name:      "aliyun dashscope embedding",
			provider:  "aliyun-dashscope",
			modelType: "dense_embedding",
			want:      "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding",
		},
		{
			name:      "aliyun dashscope rerank",
			provider:  "aliyun-dashscope",
			modelType: "rerank",
			want:      "https://dashscope.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank",
		},
		{
			name:      "deepseek completion",
			provider:  "deepseek",
			modelType: "completion",
			want:      "https://api.deepseek.com/chat/completions",
		},
		{
			name:      "siliconflow completion",
			provider:  "siliconflow",
			modelType: "completion",
			want:      "https://api.siliconflow.cn/v1/chat/completions",
		},
		{
			name:      "siliconflow embedding",
			provider:  "siliconflow",
			modelType: "dense_embedding",
			want:      "https://api.siliconflow.cn/v1/embeddings",
		},
		{
			name:      "siliconflow rerank",
			provider:  "siliconflow",
			modelType: "rerank",
			want:      "https://api.siliconflow.cn/v1/rerank",
		},
		{
			name:      "hunyuan completion",
			provider:  "hunyuan-openAI",
			modelType: "completion",
			want:      "https://api.hunyuan.cloud.tencent.com/v1/chat/completions",
		},
		{
			name:      "hunyuan embedding",
			provider:  "hunyuan-openAI",
			modelType: "dense_embedding",
			want:      "https://api.hunyuan.cloud.tencent.com/v1/embeddings",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultURL(tt.provider, tt.modelType)
			if got != tt.want {
				t.Fatalf("DefaultURL(%q, %q) = %q, want %q", tt.provider, tt.modelType, got, tt.want)
			}
		})
	}
}

func TestNames(t *testing.T) {
	want := []string{
		"aliyun-dashscope",
		"aliyun-openAI",
		"deepseek",
		"hunyuan-openAI",
		"openAI",
		"siliconflow",
	}
	got := Names()
	if len(got) != len(want) {
		t.Fatalf("Names() length = %d, want %d: %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Names()[%d] = %q, want %q; all names: %v", i, got[i], want[i], got)
		}
	}
}
