package heartbeat

import (
	"fmt"
	"regexp"
	"strings"
)

// ParseSchedule converts natural language or raw cron expression to a cron expression.
// Supported patterns:
//   - "每天 09:00"         → "0 9 * * *"
//   - "每天 09:00 18:00"   → "0 9,18 * * *"
//   - "每小時"             → "0 * * * *"
//   - "每 30 分鐘"         → "*/30 * * * *"
//   - "每 N 小時"          → "0 */N * * *"
//   - "每週一 10:00"       → "0 10 * * 1"
//   - "0 9 * * *"          → "0 9 * * *" (passthrough)
func ParseSchedule(input string) (string, error) {
	s := strings.TrimSpace(input)
	if s == "" {
		return "", fmt.Errorf("empty schedule")
	}

	// Raw cron expression (5 fields)
	if matched, _ := regexp.MatchString(`^[0-9*/,-]+\s+[0-9*/,-]+\s+[0-9*/,-]+\s+[0-9*/,-]+\s+[0-9*/,-]+$`, s); matched {
		return s, nil
	}

	// 每 N 分鐘 / every N minutes
	if m := regexp.MustCompile(`(?:每|every)\s*(\d+)\s*(?:分鐘|分|min)`).FindStringSubmatch(s); m != nil {
		return fmt.Sprintf("*/%s * * * *", m[1]), nil
	}

	// 每小時 / every hour
	if regexp.MustCompile(`(?:每小時|every\s*hour)`).MatchString(s) {
		return "0 * * * *", nil
	}

	// 每 N 小時 / every N hours
	if m := regexp.MustCompile(`(?:每|every)\s*(\d+)\s*(?:小時|hour)`).FindStringSubmatch(s); m != nil {
		return fmt.Sprintf("0 */%s * * *", m[1]), nil
	}

	// 每週X HH:MM
	weekdays := map[string]string{
		"日": "0", "一": "1", "二": "2", "三": "3", "四": "4", "五": "5", "六": "6",
		"天": "0", "mon": "1", "tue": "2", "wed": "3", "thu": "4", "fri": "5", "sat": "6", "sun": "0",
	}
	weekRe := regexp.MustCompile(`(?:每週|every\s*)([日一二三四五六天]|mon|tue|wed|thu|fri|sat|sun)\s+(\d{1,2}):(\d{2})`)
	if m := weekRe.FindStringSubmatch(strings.ToLower(s)); m != nil {
		dow, ok := weekdays[m[1]]
		if !ok {
			return "", fmt.Errorf("unknown weekday: %s", m[1])
		}
		return fmt.Sprintf("%s %s * * %s", m[3], m[2], dow), nil
	}

	// 每天 HH:MM [HH:MM ...] / every day
	dayRe := regexp.MustCompile(`(?:每天|every\s*day)\s+(.+)`)
	if m := dayRe.FindStringSubmatch(s); m != nil {
		timeRe := regexp.MustCompile(`(\d{1,2}):(\d{2})`)
		times := timeRe.FindAllStringSubmatch(m[1], -1)
		if len(times) == 0 {
			return "", fmt.Errorf("no time found in: %s", s)
		}
		var hours, mins []string
		allSameMin := true
		firstMin := times[0][2]
		for _, t := range times {
			hours = append(hours, t[1])
			mins = append(mins, t[2])
			if t[2] != firstMin {
				allSameMin = false
			}
		}
		if allSameMin {
			return fmt.Sprintf("%s %s * * *", firstMin, strings.Join(hours, ",")), nil
		}
		// Different minutes — use first time only
		return fmt.Sprintf("%s %s * * *", times[0][2], times[0][1]), nil
	}

	return "", fmt.Errorf("無法解析排程: %s", s)
}
