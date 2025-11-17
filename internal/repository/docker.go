package repository

import (
	"database/sql"

	"github.com/slodkiadrianek/octopus/internal/utils"
)

type DockerRepository struct {
	db            *sql.DB
	loggerService utils.LoggerService
}

func NewDockerRepository(db *sql.DB, loggerService utils.LoggerService) *DockerRepository {
	return &DockerRepository{
		db:            db,
		loggerService: loggerService,
	}
}

func (dr *DockerRepository) ImportContainers(ownerId int) error {
	// Implementation of importing docker containers
	return nil
}
