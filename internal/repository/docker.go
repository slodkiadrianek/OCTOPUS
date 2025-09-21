package repository

import (
	"database/sql"

	"github.com/slodkiadrianek/octopus/internal/utils/logger"
)

type DockerRepository struct {
	Db     *sql.DB
	Logger *logger.Logger
}

func NewDockerRepository(db *sql.DB, logger *logger.Logger) *DockerRepository {
	return &DockerRepository{
		Db:     db,
		Logger: logger,
	}
}


func (dr *DockerRepository) ImportContainers(ownerId int) error {
	// Implementation of importing docker containers
	return nil
}
