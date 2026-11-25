package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func loadSession() (Session, error) {
	var session Session
	data, err := os.ReadFile(activeSessionPath)
	if err != nil {
		return session, err
	}
	err = json.Unmarshal(data, &session)
	if err != nil {
		return session, err
	}
	return session, nil
}

func migrateCSV() error {
	file, err := os.Open(completedSessionsPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	closeErr := file.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}
	if len(records) == 0 {
		return nil
	}
	header := strings.Join(records[0], ",")
	switch header {
	case "module,start_time,end_time,duration_seconds,note":
		return nil
	case "module,start_time,end_time,duration_seconds":
		for i := range records {
			records[i] = append(records[i], "")
		}
		records[0][4] = "note"
	default:
		return fmt.Errorf("unsupported CSV schema")
	}
	temp, err := os.CreateTemp(filepath.Dir(completedSessionsPath), "sessions-*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(temp.Name())
	defer temp.Close()
	writer := csv.NewWriter(temp)
	writer.WriteAll(records)
	if err := writer.Error(); err != nil {
		return err
	}
	if err := temp.Sync(); err != nil {
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	return os.Rename(temp.Name(), completedSessionsPath)
}

func saveSession(session CompletedSession) error {
	if err := migrateCSV(); err != nil {
		return err
	}
	file, err := os.OpenFile(completedSessionsPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return err
	}
	writer := csv.NewWriter(file)
	if info.Size() == 0 {
		err = writer.Write([]string{"module", "start_time", "end_time", "duration_seconds", "note"})
		if err != nil {
			return err
		}
	}
	record := []string{
		session.Module,
		session.StartTime.Format(time.RFC3339Nano),
		session.EndTime.Format(time.RFC3339Nano),
		strconv.FormatInt(int64(session.Duration/time.Second), 10),
		session.Note,
	}
	err = writer.Write(record)
	if err != nil {
		return err
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}
	return file.Sync()
}

func loadCompletedSessions() ([]CompletedSession, error) {
	if err := migrateCSV(); err != nil {
		return nil, err
	}
	file, err := os.Open(completedSessionsPath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 5
	_, err = reader.Read()
	if err == io.EOF {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var sessions []CompletedSession
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		startTime, err := time.Parse(time.RFC3339Nano, record[1])
		if err != nil {
			return nil, err
		}
		endTime, err := time.Parse(time.RFC3339Nano, record[2])
		if err != nil {
			return nil, err
		}
		seconds, err := strconv.ParseInt(record[3], 10, 64)
		if err != nil {
			return nil, err
		}
		session := CompletedSession{Module: record[0], StartTime: startTime, EndTime: endTime, Duration: time.Duration(seconds) * time.Second, Note: record[4]}
		sessions = append(sessions, session)
	}
	return sessions, nil
}
