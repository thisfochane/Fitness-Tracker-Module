package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
	runningMET                 = 8.0  // для бега
)

func parseTraining(data string) (int, string, time.Duration, error) {
	input := strings.Split(data, ",")
	if len(input) > 3 {
		return 0, "", 0, fmt.Errorf("неверное количество данных: ожидается 3, получено %d", len(input))
	}

	if strings.TrimSpace(input[0]) == "" {
		return 0, "", 0, fmt.Errorf("первое значение (количество шагов) пусто")
	}

	steps, err := strconv.Atoi(string(input[0]))
	if err != nil {
		return 0, "", 0, fmt.Errorf("не удалось преобразовать количество шагов: %v", err)
	}

	activity := strings.TrimSpace(input[1])

	durationStr := strings.TrimSpace(input[2])
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("не удалось преобразовать длительность: %v", err)
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	strideLength := stepLengthCoefficient * height
	totalDistance := float64(steps) * strideLength
	totalDistance = totalDistance / mInKm
	return totalDistance
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	speed := dist / duration.Hours()
	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainings, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}
	var calories float64
	var dist float64
	var avgSpeed float64

	switch trainings {
	case "Бег":
		dist = distance(steps, height)
		avgSpeed = meanSpeed(steps, height, duration)
		calories, err = RunningSpentCalories(steps, weight, height, duration)

	case "Ходьба":
		dist = distance(steps, height)
		avgSpeed = meanSpeed(steps, height, duration)
		calories, err = WalkingSpentCalories(steps, weight, height, duration)

	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", trainings)
	}

	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		trainings,
		duration.Minutes(),
		dist,
		avgSpeed,
		calories,
	), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("Ошибка: проверка количества шагов")
	}
	if weight <= 0 {
		return 0, errors.New("ошибка: некорректный вес")
	}
	if height <= 0 {
		return 0, errors.New("ошибка: некорректный рост")
	}
	if duration <= 0 {
		return 0, errors.New("ошибка: некорректный вывод времени")
	}

	averageSpeed := meanSpeed(steps, height, duration)

	minutes := duration.Minutes()

	caloriesBurned := (weight * averageSpeed * minutes) / minInH
	return caloriesBurned, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("Ошибка: проверка количества шагов")
	}
	if weight <= 0 {
		return 0, errors.New("ошибка: некорректный вес")
	}
	if height <= 0 {
		return 0, errors.New("ошибка: некорректный рост")
	}
	if duration <= 0 {
		return 0, errors.New("ошибка: некорректный вывод времени")
	}

	averageSpeed := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()
	person := (height * averageSpeed * minutes) / minInH
	person = person * walkingCaloriesCoefficient
	return person, nil
}
