package models

type CreateProblem struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Difficulty  string     `json:"difficulty"`
	Examples    []Example  `json:"examples"`
	TestCases   []TestCase `json:"test_cases"`
}

type Example struct {
	Input       string `json:"input"`
	Output      string `json:"output"`
	Explanation string `json:"explanation"`
}

type TestCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}
