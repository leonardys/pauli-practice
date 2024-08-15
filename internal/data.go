package internal

import (
	"encoding/json"
	"os"
)

func loadData(filename string) []Session {
	var sessions []Session

	data, err := os.ReadFile(filename)
	if err != nil {
		sessions = make([]Session, 0)
	} else {
		json.Unmarshal(data, &sessions)
	}

	return sessions
}

func saveData(filename string, thisSession Session) error {
	sessions := loadData(filename)
	sessions = append(sessions, thisSession)

	json, _ := json.Marshal(sessions)
	err := os.WriteFile(filename, json, 0644)

	return err
}
