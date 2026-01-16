package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type apiServer struct {
	projectsDir string
}

type TaskInfo struct {
	ID          int               `json:"id"`
	Name        string            `json:"name"`
	Target      string            `json:"target"`
	Status      string            `json:"status"`
	CurrentStep int               `json:"current_step"`
	ErrorMsg    string            `json:"error_msg"`
	LastOutputs map[string]string `json:"last_outputs"`
	Workflow    []WorkflowNode    `json:"workflow"`
	UpdatedAt   string            `json:"updated_at"`
}

type WorkflowNode struct {
	Type       string   `json:"type"`
	Name       string   `json:"name"`
	Args       []string `json:"args,omitempty"`
	Command    string   `json:"command,omitempty"`
	InputFiles []string `json:"input_files,omitempty"`
	OutputFile string   `json:"output_file"`
}

func (a *apiServer) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tasks, err := a.listTasks()
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		writeJSON(w, tasks)
	case http.MethodPost:
		var payload TaskInfo
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeError(w, err, http.StatusBadRequest)
			return
		}
		if payload.Name == "" {
			writeError(w, errors.New("name is required"), http.StatusBadRequest)
			return
		}
		if payload.Status == "" {
			payload.Status = "Pending"
		}
		payload.UpdatedAt = time.Now().Format(time.RFC3339)
		if err := a.saveTask(payload.Name, payload); err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		writeJSON(w, payload)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a *apiServer) handleTask(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	if name == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		task, err := a.loadTask(name)
		if err != nil {
			writeError(w, err, http.StatusNotFound)
			return
		}
		writeJSON(w, task)
	case http.MethodPost:
		cmd := r.URL.Query().Get("cmd")
		task, err := a.loadTask(name)
		if err != nil {
			writeError(w, err, http.StatusNotFound)
			return
		}
		task.Status = "Running"
		task.ErrorMsg = ""
		task.UpdatedAt = time.Now().Format(time.RFC3339)
		if cmd != "" {
			task.LastOutputs = map[string]string{"last_cmd": cmd}
		}
		if err := a.saveTask(name, task); err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		writeJSON(w, task)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a *apiServer) handleLogs(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/api/logs/")
	if name == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	logPath := filepath.Join(a.projectsDir, name, "run.log")
	data, err := os.ReadFile(logPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			writeJSON(w, map[string]string{"content": ""})
			return
		}
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"content": string(data)})
}

func (a *apiServer) handleFiles(w http.ResponseWriter, r *http.Request) {
	rel := strings.TrimPrefix(r.URL.Path, "/api/files/")
	rel = filepath.Clean(rel)
	if rel == "." || strings.HasPrefix(rel, "..") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	full := filepath.Join(a.projectsDir, rel)
	file, err := os.Open(full)
	if err != nil {
		writeError(w, err, http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = io.Copy(w, file)
}

func (a *apiServer) listTasks() ([]TaskInfo, error) {
	entries, err := os.ReadDir(a.projectsDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []TaskInfo{}, nil
		}
		return nil, err
	}

	var tasks []TaskInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		info, err := a.loadTask(entry.Name())
		if err != nil {
			continue
		}
		tasks = append(tasks, info)
	}
	return tasks, nil
}

func (a *apiServer) loadTask(name string) (TaskInfo, error) {
	path := filepath.Join(a.projectsDir, name, "task_info.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return TaskInfo{}, err
	}
	var task TaskInfo
	if err := json.Unmarshal(data, &task); err != nil {
		return TaskInfo{}, err
	}
	if task.Name == "" {
		task.Name = name
	}
	return task, nil
}

func (a *apiServer) saveTask(name string, task TaskInfo) error {
	projectDir := filepath.Join(a.projectsDir, name)
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		return err
	}
	path := filepath.Join(projectDir, "task_info.json")
	return writeAtomicJSON(path, task)
}

func writeAtomicJSON(path string, payload any) error {
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}

func writeJSON(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	writeJSON(w, map[string]string{"error": err.Error()})
}
