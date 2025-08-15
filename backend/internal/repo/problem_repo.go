package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/kar1mov-u/LeetClone/internal/models"
)

type ProblemRepo struct {
	conn Queryer
}

func NewProblemRepo(conn Queryer) *ProblemRepo {
	return &ProblemRepo{conn: conn}
}

func (r *ProblemRepo) WithTx(tx pgx.Tx) *ProblemRepo {
	return NewProblemRepo(tx)
}

func (r *ProblemRepo) CreateProblem(context context.Context, problemData *models.CreateProblem) (int, error) {
	id := -1
	query := `INSERT INTO problems (name, description, difficulty) VALUES ($1, $2, $3) RETURNING id `
	row := r.conn.QueryRow(context, query, problemData.Name, problemData.Description, problemData.Difficulty)
	err := row.Scan(&id)

	return id, err
}

func (r *ProblemRepo) CreateExamples(context context.Context, examples []models.Example, problemID int) error {
	batch := &pgx.Batch{}
	for _, example := range examples {
		batch.Queue("INSERT INTO example (input, output, explanation, problem_id) VALUES ($1, $2, $3, $4)", example.Input, example.Output, example.Explanation, problemID)
	}
	br := r.conn.SendBatch(context, batch)
	return br.Close()
}

func (r *ProblemRepo) CreateTestCases(context context.Context, testcases []models.TestCase, problemID int) error {
	batch := &pgx.Batch{}
	for _, testcase := range testcases {
		batch.Queue("INSERT INTO testcase (input, output, problem_id) VALUES ($1, $2, $3)", testcase.Input, testcase.Output, problemID)
	}
	br := r.conn.SendBatch(context, batch)
	return br.Close()
}
