package repository

import (
	"database/sql"

	"github.com/slodkiadrianek/octopus/internal/utils"
)

type DockerRepository struct {
	Db     *sql.DB
	Logger *utils.Logger
}

func NewDockerRepository(db *sql.DB, logger *utils.Logger) *DockerRepository {
	return &DockerRepository{
		Db:     db,
		Logger: logger,
	}
}

func (dr *DockerRepository) ImportContainers(ownerId int) error {
	// Implementation of importing docker containers
	return nil
}
