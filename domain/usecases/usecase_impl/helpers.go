package usecase_impl

import (
	"fmt"
	"time"
)

// formatDate formata uma data no padrão brasileiro
func formatDate(t time.Time) string {
	months := map[time.Month]string{
		time.January:   "jan",
		time.February:  "fev",
		time.March:     "mar",
		time.April:     "abr",
		time.May:       "mai",
		time.June:      "jun",
		time.July:      "jul",
		time.August:    "ago",
		time.September: "set",
		time.October:   "out",
		time.November:  "nov",
		time.December:  "dez",
	}

	return fmt.Sprintf("%02d de %s, %d", t.Day(), months[t.Month()], t.Year())
}

// truncateDescription trunca uma string se ultrapassar o tamanho máximo
func truncateDescription(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
