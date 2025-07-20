package config

import (
	"bufio"
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

func (e *Env) ReadFile() map[string]string {
	file, err := os.OpenFile(".env", os.O_RDONLY, 0o644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	envVariables := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lineSplitted := strings.Split(line, "=")
		envVariables[lineSplitted[0]] = lineSplitted[1]
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return envVariables
}

func (e *Env) Set() *Env {
	envVariables := e.ReadFile()
	return &Env{
		Port:         envVariables["Port"],
		JWTSecret:    envVariables["JWTSecret"],
		DbLink:       envVariables["DbLink"],
		CacheLink:    envVariables["CacheLink"],
		EmailService: envVariables["EmailService"],
		EmailUser:    envVariables["EmailUser"],
		EmailFrom:    envVariables["EmailFrom"],
		EmailPass:    envVariables["EmailPass"],
	}
}
