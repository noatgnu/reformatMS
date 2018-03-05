package input

import (
	"bufio"
	"fmt"
	"os"
)

func Input(message string) (filename string, err error) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(message)
	if scanner.Scan() {
		filename = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return filename, nil
}
