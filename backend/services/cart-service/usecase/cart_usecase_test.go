package usecase_test

import (
	"context"
	"testing"
)

// TestExample はCIパイプラインの動作確認用のサンプルテスト
func TestExample(t *testing.T) {
	ctx := context.Background()

	// 基本的なテストの例
	t.Run("context should not be nil", func(t *testing.T) {
		if ctx == nil {
			t.Error("context should not be nil")
		}
	})

	t.Run("basic calculation", func(t *testing.T) {
		result := 2 + 2
		expected := 4
		if result != expected {
			t.Errorf("expected %d, got %d", expected, result)
		}
	})
}

// TestTableDriven はテーブル駆動テストの例
func TestTableDriven(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"positive number", 5, 5},
		{"zero", 0, 0},
		{"negative number", -3, -3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.input != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, tt.input)
			}
		})
	}
}
