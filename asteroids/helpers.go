// File helpers.go provides cross-platform helper functions for reading and
// writing the player’s high score to the local file system.
package asteroids

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"strings"
)

// getHighScore reads the player's stored high score from a file,
// creating the directory and file if they do not yet exist.
//
// The save path differs per OS:
//   - macOS:   ~/Library/Application Support/Asteroids
//   - Windows: C:\Users\<user>\AppData
//   - Linux:   /users/<user> or /home/<user>/.asteroids
func getHighScore() (int, error) {
	// Resolve the current OS user.
	user, err := user.Current()
	if err != nil {
		return 0, err
	}

	// Build the appropriate platform path.
	var path string
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

	// Ensure the directory exists.
	if _, err := os.Stat(path); err != nil {
		if err := os.Mkdir(path, 0750); err != nil {
			return 0, err
		}
	}

	// Ensure the file exists with a default value.
	scoreFile := path + "/high-score.txt"
	if _, err := os.Stat(scoreFile); err != nil {
		if err := os.WriteFile(scoreFile, []byte("0"), 0750); err != nil {
			return 0, err
		}
	}

	// Read and parse score.
	data, err := os.ReadFile(scoreFile)
	if err != nil {
		return 0, err
	}
	s := strings.TrimSpace(string(data))
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return value, nil
}

// updateHighScore writes a new integer score value to the user’s
// high score file, overwriting any previous value.
func updateHighScore(score int) error {
	user, err := user.Current()
	if err != nil {
		return err
	}

	// Construct platform-specific file path.
	var path string
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

	// Write integer as plain text.
	return os.WriteFile(path, []byte(fmt.Sprintf("%d", score)), 0750)
}
