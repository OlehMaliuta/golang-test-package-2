package test2

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type DVD struct {
	Name     string
	Director string
	Duration time.Duration
	Price    uint
}

func WriteDVDsToFile(filename string, dvds []DVD) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, dvd := range dvds {
		_, err := writer.WriteString(fmt.Sprintf("%s,%s,%d,%d\n", dvd.Name, dvd.Director, int(dvd.Duration.Seconds()), dvd.Price))
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func ReadDVDsFromFile(filename string) ([]DVD, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var dvds []DVD
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var name, director string
		var durationSeconds int
		var price uint

		parts := strings.Split(line, ",")
		if len(parts) != 4 {
			return nil, fmt.Errorf("неправильний формат даних у файлі: %s", line)
		}

		name = parts[0]
		director = parts[1]
		fmt.Sscanf(parts[2], "%d", &durationSeconds)
		fmt.Sscanf(parts[3], "%d", &price)

		dvds = append(dvds, DVD{
			Name:     name,
			Director: director,
			Duration: time.Duration(durationSeconds) * time.Second,
			Price:    price,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return dvds, nil
}

func RemoveDVDByDuration(filename string, targetDuration time.Duration) error {
	dvds, err := ReadDVDsFromFile(filename)
	if err != nil {
		return err
	}

	for i, dvd := range dvds {
		if dvd.Duration == targetDuration {
			dvds = append(dvds[:i], dvds[i+1:]...)
			break
		}
	}

	return WriteDVDsToFile(filename, dvds)
}

func AddDVDsAtIndex(filename string, newDVDs []DVD, index int) error {
	dvds, err := ReadDVDsFromFile(filename)
	if err != nil {
		return err
	}

	if index < 0 || index > len(dvds) {
		return fmt.Errorf("невірний індекс")
	}

	dvds = append(dvds[:index], append(newDVDs, dvds[index:]...)...)

	return WriteDVDsToFile(filename, dvds)
}
