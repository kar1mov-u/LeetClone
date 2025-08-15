package services

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kar1mov-u/LeetClone/internal/models"
	"github.com/kar1mov-u/LeetClone/internal/repo"
)

type ProblemService struct {
	problemRepo *repo.ProblemRepo
	txStarter   TxStarter
}

func (s *ProblemService) CreateProblem(ctx context.Context, problemData models.CreateProblem) (int, error) {
	//create a TX
	tx, err := s.txStarter.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return -1, fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	//create new repo that uses tx as a connection to DB, instead of pool
	rp := s.problemRepo.WithTx(tx)

	id, err := rp.CreateProblem(ctx, &problemData)
	if err != nil {
		return -1, fmt.Errorf("create problem: %w", err)
	}
	err = rp.CreateExamples(ctx, problemData.Examples, id)
	if err != nil {
		return -1, fmt.Errorf("create examples: %w", err)
	}
	err = rp.CreateTestCases(ctx, problemData.TestCases, id)
	if err != nil {
		return -1, fmt.Errorf("craete test cases: %w", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit tx:%w", err)
	}
	return id, nil

}
