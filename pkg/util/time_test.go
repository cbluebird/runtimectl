package util

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	timeStr := "2024-10-28-0326"

	// 解析字符串为 time.Time 类型
	parsedTime, err := time.Parse("2006-01-02-1504", timeStr)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return
	}

	// 打印解析后的时间
	fmt.Println("Parsed time:", parsedTime)
}
