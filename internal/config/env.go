package config

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type Env struct {
	Port         string
	JWTSecret    string
	DbLink       string
	CacheLink    string
	EmailService string
	EmailUser    string
	EmailFrom    string
	EmailPass    string
}

func ReadFile(filepath string) (map[string]string, error) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0o644)
	if err != nil {
		return map[string]string{}, errors.New("failed to open a file")
	}
	defer file.Close()
	envVariables := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lineSplitted := strings.SplitN(line, "=", 2)
		envVariables[lineSplitted[0]] = lineSplitted[1]
	}
	if err := scanner.Err(); err != nil {
		return map[string]string{}, errors.New("failed to scan a file")
	}
	return envVariables, nil
}

func SetConfig(filepath string) (*Env, error) {
	envVariables, err := ReadFile(filepath)
	if err != nil {
		return &Env{}, err
	}
	return &Env{
		Port:         envVariables["Port"],
		JWTSecret:    envVariables["JWTSecret"],
		DbLink:       envVariables["DbLink"],
		CacheLink:    envVariables["CacheLink"],
		EmailService: envVariables["EmailService"],
		EmailUser:    envVariables["EmailUser"],
		EmailFrom:    envVariables["EmailFrom"],
		EmailPass:    envVariables["EmailPass"],
	}, nil
}
