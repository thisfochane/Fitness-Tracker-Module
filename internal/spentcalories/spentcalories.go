package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

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
	if len(input) != 3 {
		return 0, "", 0, fmt.Errorf("неверное количество данных: ожидается 3, получено %d", len(input))
	}

	stepsStr := strings.TrimSpace(input[0])
	if stepsStr == "" {
		return 0, "", 0, fmt.Errorf("количество шагов не может быть пустым")
	}

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("не удалось преобразовать количество шагов '%s': %v", stepsStr, err)
	}

	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов должно быть больше нуля")
	}

	activity := strings.TrimSpace(input[1])

	durationStr := strings.TrimSpace(input[2])
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("не удалось преобразовать длительность '%s': %v", durationStr, err)
	}

	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("длительность должна быть больше нуля %v", duration)
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	strideLength := stepLengthCoefficient * height
	totalDistance := float64(steps) * strideLength / mInKm
	return totalDistance
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration.Seconds() <= 0 {
		log.Println("Ошибка: длительность должна быть больше нуля")
		return 0
	}
	dist := distance(steps, height)
	speed := dist / duration.Hours()
	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println("Ошибка при разборе данных:", err)
		return "", err
	}

	var calories float64
	var dist float64
	var avgSpeed float64

	dist = distance(steps, height)
	avgSpeed = meanSpeed(steps, height, duration)

	switch activity {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)

	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)

	default:
		errMsg := fmt.Sprintf("неизвестный тип тренировки: %s", activity)
		log.Println(errMsg)
		return "", errors.New(errMsg)
	}

	if err != nil {
		log.Println("Ошибка при расчете калорий:", err)
		return "", err
	}

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity,
		duration.Hours(),
		dist,
		avgSpeed,
		calories,
	), nil
}

func RunningSpentCalories(steps int, weight float64, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		log.Println("Ошибка: количество шагов должно быть больше нуля")
		return 0, errors.New("ошибка: проверка количества шагов")
	}
	if weight <= 0 {
		log.Println("Ошибка: некорректный вес")
		return 0, errors.New("ошибка: некорректный вес")
	}
	if height <= 0 {
		log.Println("Ошибка: некорректный рост")
		return 0, errors.New("ошибка: некорректный рост")
	}
	if duration.Seconds() <= 0 {
		log.Println("Ошибка: некорректная длительность времени")
		return 0, errors.New("ошибка: некорректный вывод времени")
	}

	averageSpeed := meanSpeed(steps, height, duration)

	minutes := duration.Minutes()

	caloriesBurned := (weight * averageSpeed * minutes) / minInH
	return caloriesBurned, nil
}

func WalkingSpentCalories(steps int, weight float64, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		log.Println("Ошибка: количество шагов должно быть больше нуля")
		return 0.0, errors.New("ошибка: проверка количества шагов")
	}
	if weight <= 0 {
		log.Println("Ошибка: некорректный вес")
		return 0.0, errors.New("ошибка: некорректный вес")
	}
	if height <= 0 {
		log.Println("Ошибка: некорректный рост")
		return 0.0, errors.New("ошибка: некорректный рост")
	}
	if duration.Seconds() <= 0 {
		log.Println("Ошибка: некорректная длительность времени")
		return 0.0, errors.New("ошибка: некорректный вывод времени")
	}

	averageSpeed := meanSpeed(steps, height, duration)

	minutes := duration.Minutes()

	caloriesBurned := (weight * averageSpeed * minutes) / minInH
	caloriesBurned *= walkingCaloriesCoefficient

	return caloriesBurned, nil
}
