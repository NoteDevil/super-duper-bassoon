package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	stepLength = 0.65
	mInKm      = 1000
)

// parsePackage - функция для парсинга строки с данными о тренировке в формате "шаги, продолжительность".
// Она возвращает количество шагов и продолжительность, или ошибку, если формат данных неверен.
//
// Формат данных должен быть следующим: "шаги, продолжительность", где шаги - целое положительное число,
// а продолжительность - строка в формате "XhYm", где X - часы, а Y - минуты.
// Если продолжительность не указана, то она считается равной 0.
//
// Если формат данных неверен, то функция возвращает ошибку.
func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("неверный формат данных: %s", data)
	}

	stepsStr := parts[0]
	durationStr := parts[1]

	if strings.TrimSpace(stepsStr) == "" {
		return 0, 0, fmt.Errorf("количество шагов не может быть пустым")
	}

	if stepsStr != strings.TrimSpace(stepsStr) {
		return 0, 0, fmt.Errorf("неверный формат шагов: пробелы не допускаются")
	}

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга шагов: %v", err)
	}

	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть положительным: %d", steps)
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга продолжительности: %v", err)
	}

	if duration <= 0 {
		return 0, 0, fmt.Errorf("продолжительность должна быть положительной: %v", duration)
	}

	return steps, duration, nil
}

// DayActionInfo - функция для расчета информации о дневной активности.
// Она принимает строку с данными о тренировке в формате "шаги, продолжительность"
// и возвращает строку с информацией о количестве шагов, пройденном расстоянии и сожжении калорий.
// Если при парсинге данных происходит ошибка, то функция возвращает пустую строку.
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distanceKm, calories)
}
