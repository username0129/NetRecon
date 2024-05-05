package util

import (
	"fmt"
	"time"
)

// TimeToCronSpec 将给定的时间转换为每天都执行的cron表达式
func TimeToCronSpec(t time.Time) string {
	// 返回格式: 秒 分钟 小时 日期 月份 星期
	return fmt.Sprintf("%d %d %d * * *", t.Second(), t.Minute(), t.Hour())
}
