// main_test.go
package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
)

func createTempFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "names.txt")
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("не удалось создать временный файл: %v", err)
	}
	return path
}

func TestCountNames_Basic(t *testing.T) {
	path := createTempFile(t, "Алёна\nМиша\nАлёна\nДима\n")

	counts, err := countNames(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := map[string]int{
		"Алёна": 2,
		"Миша":  1,
		"Дима":  1,
	}

	if len(counts) != len(expected) {
		t.Fatalf("expected %d unique names, got %d", len(expected), len(counts))
	}

	for name, want := range expected {
		got, ok := counts[name]
		if !ok {
			t.Errorf("name %q not found in result", name)
			continue
		}
		if got != want {
			t.Errorf("name %q: expected count %d, got %d", name, want, got)
		}
	}
}

func TestCountNames_EmptyFile(t *testing.T) {
	path := createTempFile(t, "")

	counts, err := countNames(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(counts) != 0 {
		t.Fatalf("expected 0 names, got %d", len(counts))
	}
}

func TestCountNames_TrimSpaces(t *testing.T) {
	path := createTempFile(t, "  Олег  \nОлег\n  Олег\n")

	counts, err := countNames(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if counts["Олег"] != 3 {
		t.Errorf("expected Олег:3, got Олег:%d", counts["Олег"])
	}
}

func TestCountNames_SingleName(t *testing.T) {
	path := createTempFile(t, "Иван\n")

	counts, err := countNames(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if counts["Иван"] != 1 {
		t.Errorf("expected Иван:1, got Иван:%d", counts["Иван"])
	}
}

func TestCountNames_FileNotFound(t *testing.T) {
	_, err := countNames("/nonexistent/path/file.txt")
	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}
}

func TestCountNames_HugeFile(t *testing.T) {
	names := []string{
		"Александр", "Мария", "Дмитрий", "Анна", "Максим",
		"Елена", "Иван", "Ольга", "Артём", "Наталья",
		"Сергей", "Екатерина", "Андрей", "Татьяна", "Алексей",
		"Юлия", "Михаил", "Ирина", "Никита", "Светлана",
		"Павел", "Анастасия", "Владимир", "Дарья", "Роман",
		"Валентина", "Егор", "Людмила", "Тимофей", "Галина",
	}

	const totalLines = 10_000_000
	expectedCounts := make(map[string]int)

	rng := rand.New(rand.NewSource(42))

	dir := t.TempDir()
	path := filepath.Join(dir, "huge_names.txt")

	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("не удалось создать файл: %v", err)
	}

	for i := 0; i < totalLines; i++ {
		name := names[rng.Intn(len(names))]
		expectedCounts[name]++
		fmt.Fprintln(file, name)
	}
	file.Close()

	counts, err := countNames(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(counts) != len(expectedCounts) {
		t.Fatalf("expected %d unique names, got %d", len(expectedCounts), len(counts))
	}

	totalCounted := 0
	for name, want := range expectedCounts {
		got, ok := counts[name]
		if !ok {
			t.Errorf("name %q not found in result", name)
			continue
		}
		if got != want {
			t.Errorf("name %q: expected %d, got %d", name, want, got)
		}
		totalCounted += got
	}

	if totalCounted != totalLines {
		t.Errorf("total count mismatch: expected %d, got %d", totalLines, totalCounted)
	}
}
