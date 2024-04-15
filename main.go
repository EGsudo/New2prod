package main

import (
	"fmt"
	"math"
	"time"
)

const (
	MInKm      = 1000
	MinInHours = 60
	LenStep    = 0.65
	CmInM      = 100
)

type Training struct {
	TrainingType string
	Action       int
	LenStep      float64
	Duration     time.Duration
	Weight       float64
}

func (t Training) distance() float64 {

	return float64(t.Action) * t.LenStep / MInKm
}
func (t Training) meanSpeed() float64 {
	return (float64(t.Action) * t.LenStep / MInKm) / t.Duration.Hours()
}
func (t Training) Calories() float64 {
	return 0
}

type InfoMessage struct {
	TrainingType string
	Duration     time.Duration
	Distance     float64
	Speed        float64
	Calories     float64
}

func (t Training) TrainingInfo() InfoMessage {
	return InfoMessage{}
}
func (i InfoMessage) String() string {
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %v мин\nДистанция: %.2f км.\nСр. скорость: %.2f км/ч\nПотрачено ккал: %.2f\n",
		i.TrainingType,
		i.Duration.Minutes(),
		i.Distance,
		i.Speed,
		i.Calories,
	)
}

type CaloriesCalculator interface {
	Calories() float64
	TrainingInfo() InfoMessage
}

const (
	CaloriesMeanSpeedMultiplier = 18
	CaloriesMeanSpeedShift      = 1.79
)

type Running struct {
	Training
}

func (r Running) Calories() float64 {
	return ((CaloriesMeanSpeedMultiplier*r.meanSpeed() + CaloriesMeanSpeedShift) * r.Weight / MInKm * float64(r.Training.Duration.Hours()) * MinInHours)
}
func (r Running) TrainingInfo() InfoMessage {
	info := r.Training.TrainingInfo()
	info.TrainingType = r.TrainingType
	info.Duration = r.Duration
	info.Distance = r.Training.distance()
	info.Speed = r.meanSpeed()
	info.Calories = r.Calories()
	return info
}

const (
	CaloriesWeightMultiplier      = 0.035
	CaloriesSpeedHeightMultiplier = 0.029
	KmHInMsec                     = 0.278
)

type Walking struct {
	Training
	Height float64
}

func (w Walking) Calories() float64 {
	return (CaloriesWeightMultiplier*w.Training.Weight + math.Pow(w.meanSpeed()*KmHInMsec, 2)) / (w.Height / CmInM) * CaloriesSpeedHeightMultiplier * w.Training.Weight * w.Training.Duration.Hours() * MinInHours
}
func (w Walking) TrainingInfo() InfoMessage {
	info := w.Training.TrainingInfo()
	info.TrainingType = w.TrainingType
	info.Duration = w.Duration
	info.Distance = w.Training.distance()
	info.Speed = w.meanSpeed()
	info.Calories = w.Calories()
	return info
}

const (
	SwimmingLenStep                  = 1.38
	SwimmingCaloriesMeanSpeedShift   = 1.1
	SwimmingCaloriesWeightMultiplier = 2
)

type Swimming struct {
	Training
	LengthPool int
	CountPool  int
}

func (s Swimming) meanSpeed() float64 {
	return float64(s.LengthPool) * float64(s.CountPool) / MInKm / float64(s.Training.Duration.Hours())
}
func (s Swimming) Calories() float64 {
	return (s.meanSpeed() + SwimmingCaloriesMeanSpeedShift) * SwimmingCaloriesWeightMultiplier * s.Training.Weight * (s.Duration.Hours())
}
func (s Swimming) TrainingInfo() InfoMessage {
	info := s.Training.TrainingInfo()
	info.TrainingType = s.TrainingType
	info.Duration = s.Duration
	info.Distance = s.Training.distance()
	info.Speed = s.meanSpeed()
	info.Calories = s.Calories()
	return info
}
func ReadData(training CaloriesCalculator) string {
	training.Calories()
	info := training.TrainingInfo()
	return fmt.Sprint(info)
}
func main() {
	swimming := Swimming{
		Training: Training{
			TrainingType: "Плавание",
			Action:       2000,
			LenStep:      SwimmingLenStep,
			Duration:     90 * time.Minute,
			Weight:       85,
		},
		LengthPool: 50,
		CountPool:  5,
	}
	fmt.Println(ReadData(swimming))
	walking := Walking{
		Training: Training{
			TrainingType: "Ходьба",
			Action:       20000,
			LenStep:      LenStep,
			Duration:     3*time.Hour + 45*time.Minute,
			Weight:       85,
		},
		Height: 185,
	}
	fmt.Println(ReadData(walking))
	running := Running{
		Training: Training{
			TrainingType: "Бег",
			Action:       5000,
			LenStep:      LenStep,
			Duration:     30 * time.Minute,
			Weight:       85,
		},
	}
	fmt.Println(ReadData(running))
}
