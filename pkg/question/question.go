package question

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	InvalidInputDivisionBy0 = errors.New("invalid input: disallowed division by 0")
	InvalidOperation        = errors.New(fmt.Sprintf("invalid operation, only accept + - * /"))
)

type Question struct {
	Name      string
	a         int32
	b         int32
	operation string
	answer    float64
}

func NewQuestion(name string, a int32, b int32, operation string) (*Question, error) {
	q := &Question{
		Name:      name,
		a:         a,
		b:         b,
		operation: operation,
	}

	switch q.operation {
	case "+":
		q.answer = float64(q.a + q.b)
		return q, nil
	case "-":
		q.answer = float64(q.a - q.b)
	case "*":
		q.answer = float64(q.a * q.b)
	case "/":
		if q.b == 0 {
			return nil, InvalidInputDivisionBy0
		}
		q.answer = float64(q.a) / float64(q.b)
	default:
		return nil, InvalidOperation
	}
	return q, nil
}

func (q *Question) GetAnswer() float64 {
	return q.answer
}

func (q *Question) QuestionString() string {
	if len(q.Name) > 0 {
		return fmt.Sprintf("%s: %d %s %d = ?", q.Name, q.a, q.operation, q.b)
	}
	return fmt.Sprintf("%d %s %d = ?", q.a, q.operation, q.b)
}

func (q *Question) AnswerString() string {
	if len(q.Name) > 0 {
		return fmt.Sprintf("%s: %d %s %d = %d", q.Name, q.a, q.operation, q.b, ConvertFloatToString(q.answer))
	}
	return fmt.Sprintf("%d %s %d = %d", q.a, q.operation, q.b, ConvertFloatToString(q.answer))
}

func (q *Question) GuessedAnswerString(guess float64) string {
	if len(q.Name) > 0 {
		return fmt.Sprintf("%s: %d %s %d = %s", q.Name, q.a, q.operation, q.b, ConvertFloatToString(guess))
	}
	return fmt.Sprintf("%d %s %d = %s", q.a, q.operation, q.b, ConvertFloatToString(guess))
}

func (q *Question) IsCorrect(guess float64) bool {
	return guess == q.answer
}

func ConvertFloatToString(input float64) string {
	return strconv.FormatFloat(input, 'f', -1, 64)
}
