// main.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type NameCount struct {
	Name  string
	Count int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Использование: %s <имя_файла>\n", os.Args[0])
		os.Exit(1)
	}

	filename := os.Args[1]

	counts, err := countNames(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}

	result := make([]NameCount, 0, len(counts))
	for name, count := range counts {
		result = append(result, NameCount{Name: name, Count: count})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Count != result[j].Count {
			return result[i].Count > result[j].Count
		}
		return result[i].Name < result[j].Name
	})

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	for _, nc := range result {
		fmt.Fprintf(writer, "%s:%d\n", nc.Name, nc.Count)
	}
}

func countNames(filename string) (map[string]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл %q: %w", filename, err)
	}
	defer file.Close()

	counts := make(map[string]int)

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		name := strings.TrimSpace(scanner.Text())
		if name == "" {
			continue
		}
		counts[name]++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %w", err)
	}

	return counts, nil
}
