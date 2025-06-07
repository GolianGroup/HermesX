package services

import (
	"context"
	"golang_template/internal/repositories"
)

type RepoStatus struct {
	Healthy bool
	Error   error
}

type SystemService interface {
	ReadyCheck(ctx context.Context) (map[string]RepoStatus, []error)
}

type systemService struct {
	systemRepository repositories.SystemRepository
}

func NewSystemService(systemRepository repositories.SystemRepository) SystemService {
	return &systemService{systemRepository: systemRepository}
}

func (s *systemService) ReadyCheck(ctx context.Context) (map[string]RepoStatus, []error) {

	statuses := make(map[string]RepoStatus)
	var errors []error

	// Check Postgres
	if err := s.systemRepository.DBPing(ctx); err != nil {
		statuses["postgres"] = RepoStatus{Healthy: false, Error: err}
		errors = append(errors, err)
	} else {
		statuses["postgres"] = RepoStatus{Healthy: true, Error: nil}
	}

	// Check Clickhouse
	if err := s.systemRepository.ClickhousePing(ctx); err != nil {
		statuses["clickhouse"] = RepoStatus{Healthy: false, Error: err}
		errors = append(errors, err)
	} else {
		statuses["clickhouse"] = RepoStatus{Healthy: true, Error: nil}
	}

	// Check ArangoDB
	if err := s.systemRepository.ArangoPing(ctx); err != nil {
		statuses["arango"] = RepoStatus{Healthy: false, Error: err}
		errors = append(errors, err)
	} else {
		statuses["arango"] = RepoStatus{Healthy: true, Error: nil}
	}

	// Check Redis
	if err := s.systemRepository.RedisPing(ctx); err != nil {
		statuses["redis"] = RepoStatus{Healthy: false, Error: err}
		errors = append(errors, err)
	} else {
		statuses["redis"] = RepoStatus{Healthy: true, Error: nil}
	}

	if len(errors) > 0 {
		return statuses, errors
	}

	return statuses, nil

}
