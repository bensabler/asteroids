package asteroids

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"strings"
)

func getHighScore() (int, error) {
	// Get the user name
	user, err := user.Current()
	if err != nil {
		return 0, err
	}

	path := ""
	switch runtime.GOOS {
	case "darwin":
		path = fmt.Sprintf("/Users/%s/Library/Application Support/Asteroids", user.Username)
	case "windows":
		path = fmt.Sprintf("C:\\Users\\%s\\AppData", user.Username)
	case "linux":
		path = fmt.Sprintf("/users/%s", user.Username)
	default:
		path = fmt.Sprintf("/home/%s/.asteroids", user.Username)
	}

	if _, err := os.Stat(path); err != nil {
		if err := os.Mkdir(path, 0750); err != nil {
			return 0, err
		}
	}

	if _, err := os.Stat(path + "/high-score.txt"); err != nil {
		err := os.WriteFile(path+"/high-score.txt", []byte("0"), 0750)
		if err != nil {
			return 0, err
		}
	}

	contents, err := os.ReadFile(path + "/high-score.txt")
	if err != nil {
		return 0, err
	}

	score := string(contents)
	score = strings.TrimSpace(score)
	s, err := strconv.Atoi(string(score))
	if err != nil {
		return 0, err
	}
	return s, nil
}

func updateHighScore(score int) error {
	// Get the user name
	user, err := user.Current()
	if err != nil {
		return err
	}

	path := ""
	switch runtime.GOOS {
	case "darwin":
		path = fmt.Sprintf("/Users/%s/Library/Application Support/Asteroids/high-score.txt", user.Username)
	case "windows":
		path = fmt.Sprintf("C:\\Users\\%s\\AppData\\high-score.txt", user.Username)
	case "linux":
		path = fmt.Sprintf("/users/%s/high-score.txt", user.Username)
	default:
		path = fmt.Sprintf("/home/%s/.asteroids/high-score.txt", user.Username)
	}

	s := fmt.Sprintf("%d", score)
	if err := os.WriteFile(path, []byte(s), 0750); err != nil {
		return err
	}
	return nil
}
