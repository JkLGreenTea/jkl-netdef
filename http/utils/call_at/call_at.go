package call_at

import (
	"time"
)

// CallAt - Вызов переданной функции раз в сутки в указанное время.
func CallAt(hour, min, sec int, f func() error) error {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return err
	}

	// Вычисляем время первого запуска.
	now := time.Now().Local()
	firstCallTime := time.Date(
		now.Year(), now.Month(), now.Day(), hour, min, sec, 0, loc)
	if firstCallTime.Before(now) {
		// Если получилось время раньше текущего, прибавляем сутки.
		firstCallTime = firstCallTime.Add(time.Hour * 24)
	}

	// Вычисляем временной промежуток до запуска.
	duration := firstCallTime.Sub(time.Now().Local())

	time.Sleep(duration)
	for {
		err = f()
		if err != nil {
			return err
		}

		// Следующий запуск через сутки.
		time.Sleep(time.Hour * 24)
	}
}
