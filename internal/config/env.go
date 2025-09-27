package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Env struct {
	Port       string
	JWTSecret  string
	DbLink     string
	CacheLink  string
	DockerHost string
}

func ReadFile(filepath string) (map[string]string, error) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filepath, err)
	}
	defer file.Close()
	envVariables := make(map[string]string)
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid format at line %d: %s", lineNum, line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" {
			return nil, fmt.Errorf("empty key at line %d", lineNum)
		}
		envVariables[key] = value
	}
	return envVariables, scanner.Err()
}

func (e *Env) Validate() error {
	if e.Port == "" {
		return errors.New("PORT is required")
	}
	if e.JWTSecret == "" {
		return errors.New("JWT_SECRET is required")
	}
	if e.DbLink == "" {
		return errors.New("DbLink is required")
	}
	if e.CacheLink == "" {
		return errors.New("CacheLink is required")
	}
	if e.DockerHost == "" {
		return errors.New("DockerHost is required")
	}
	return nil
}

func SetConfig(filepath string) (*Env, error) {
	envVariables, err := ReadFile(filepath)
	if err != nil {
		return &Env{}, err
	}
	return &Env{
		Port:       envVariables["Port"],
		JWTSecret:  envVariables["JWTSecret"],
		DbLink:     envVariables["DbLink"],
		CacheLink:  envVariables["CacheLink"],
		DockerHost: envVariables["DockerHost"],
	}, nil
}
