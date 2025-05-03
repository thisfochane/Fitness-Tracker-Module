package daysteps

import (
	"Fitness-Tracker-Module/internal/spentcalories"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	input := strings.Split(data, ",")

	if len(input) != 2 {
		return 0, 0, fmt.Errorf("неверное количество данных: ожидается 2, получено %d", len(input))
	}
	if strings.TrimSpace(input[0]) == "" {
		return 0, 0, fmt.Errorf("первое значение (количество шагов) пусто")
	}
	steps, err := strconv.Atoi(string(input[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("не удалось преобразовать '%s' в число: %v", input, err)
	}

	if steps <= 0 {
		return 0, 0, fmt.Errorf("ошибка колличество шагов не может равняться 0: %v", err)
	}

	duration, err := time.ParseDuration(input[1])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка преобразования времени: %v", err)
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("ошибка продолжительности: продолжительность не может быть равна нулю")
	}
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	fmt.Println("Входные данные:", data)

	if strings.TrimSpace(data) == "" {
		log.Println("Ошибка: входные данные пусты")
		return ""
	}

	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println("Ошибка: получение функции ", err)
		return ""
	}
	if steps <= 0 {
		log.Println("steps is less or equal to zero:", steps)
		return ""
	}

	distance := float64(steps) * stepLength / mInKm

	caloriesBurned, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println("Ошибка при расчете калорий:", err)
		return ""
	}
	info := fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps,
		distance,
		caloriesBurned,
	)
	return info
}
