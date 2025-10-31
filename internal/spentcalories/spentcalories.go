package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	lenStep                    = 0.65
	mInKm                      = 1000
	minInH                     = 60
	stepLengthCoefficient      = 0.45
	walkingCaloriesCoefficient = 0.5
)

// parseTraining - функция для парсинга строки с данными о тренировке.
// Она возвращает количество шагов, тип тренировки, продолжительность, или ошибку, если формат данных неверен.
// Формат данных должен быть следующим: "шаги, тип тренировки, продолжительность", где шаги - целое положительное число,
// а тип тренировки и продолжительность - строки.
// Если формат данных неверен, то функция возвращает ошибку.
func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат данных: %s", data)
	}

	stepsStr := strings.TrimSpace(parts[0])
	trainingType := strings.TrimSpace(parts[1])
	durationStr := strings.TrimSpace(parts[2])

	if stepsStr == "" {
		return 0, "", 0, fmt.Errorf("количество шагов не может быть пустым")
	}

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка парсинга шагов: %v", err)
	}

	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов должно быть положительным: %d", steps)
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка парсинга продолжительности: %v", err)
	}

	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("продолжительность должна быть положительной: %v", duration)
	}

	return steps, trainingType, duration, nil
}

// distance - функция для расчета дистанции, пройденной за время тренировки.
// Она принимает количество шагов и рост и возвращает пройденную дистанцию в километрах.
// Функция корректно работает только при положительных значениях входных параметров.
// В противном случае она возвращает ошибку.
func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distanceMeters := float64(steps) * stepLength
	return distanceMeters / mInKm
}

// meanSpeed - функция для расчета средней скорости при беге.
// Она принимает количество шагов, рост и продолжительность и возвращает среднюю скорость в километрах в час.
// Если продолжительность равна 0, то функция возвращает 0.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distanceKm := distance(steps, height)
	durationHours := duration.Hours()

	if durationHours == 0 {
		return 0
	}

	return distanceKm / durationHours
}

// RunningSpentCalories - функция для расчета сожженных калорий при беге.
// Она принимает количество шагов, вес, рост и продолжительность и возвращает количество сожженных калорий.
// Если при парсинге данных происходит ошибка, то функция возвращает ошибку.
// Функция корректно работает только при положительных значениях входных параметров.
// В противном случае она возвращает ошибку.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("некорректные входные параметры")
	}

	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()

	calories := (weight * speed * durationMinutes) / minInH
	return calories, nil
}

// WalkingSpentCalories - функция для расчета сожженных калорий при ходьбе.
// Она принимает количество шагов, вес, рост и продолжительность и возвращает количество сожженных калорий.
// Если при парсинге данных происходит ошибка, то функция возвращает ошибку.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("некорректные входные параметры")
	}

	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()

	calories := (weight * speed * durationMinutes) / minInH
	calories *= walkingCaloriesCoefficient
	return calories, nil
}

// TrainingInfo - функция для расчета информации о тренировке.
// Она принимает строку с данными о тренировке в формате "шаги, тип тренировки, продолжительность"
// и возвращает строку с информацией о типе тренировки, продолжительности, пройденном расстоянии, скорости и сожжении калорий.
// Если при парсинге данных происходит ошибка, то функция возвращает пустую строку.
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	var calories float64
	var calcErr error

	switch trainingType {
	case "Бег":
		calories, calcErr = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, calcErr = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", trainingType)
	}

	if calcErr != nil {
		return "", calcErr
	}

	distanceKm := distance(steps, height)
	speed := meanSpeed(steps, height, duration)
	durationHours := duration.Hours()

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		trainingType, durationHours, distanceKm, speed, calories), nil
}
