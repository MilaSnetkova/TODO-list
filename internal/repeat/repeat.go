package repeat

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/MilaSnetkova/TODO-list/internal/constants"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	parsedDate, err := time.Parse(constants.DateFormat, date)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %v", err)
	}
	if repeat == "" {
		return "", fmt.Errorf("no repeat rule mentioned")
	}

	// Учитываем только день
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var nextDate time.Time

	repeatSplit := strings.Split(repeat, " ")

	switch repeatSplit[0] {
	// Задача переносится на указанное число дней
	case "d":
		if len(repeatSplit) != 2 {
			return "", fmt.Errorf("invalid repeat format: %v", repeat)
		}

		days, err := strconv.Atoi(repeatSplit[1])
		if err != nil || days < 1 || days > 400 {
			return "", fmt.Errorf("invalid day range: %v", repeatSplit[1])
		}

		if parsedDate.Equal(now) {
			nextDate = now
		} else {
			nextDate = parsedDate.AddDate(0, 0, days)
			for nextDate.Before(now) {
				nextDate = nextDate.AddDate(0, 0, days)
			}
		}
		// Задача выполняется ежегодно
	case "y":
		nextDate = parsedDate.AddDate(1, 0, 0)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
		// TODO: задача назначается в указанные дни недели
	case "w":
		return "", fmt.Errorf("")
		// TODO: задача назначается в указанные дни месяца
	case "m":
		return "", fmt.Errorf("")
	default:
		return "", fmt.Errorf("invalid repeat rule: %v", repeat)
	}
	return nextDate.Format(constants.DateFormat), nil
}