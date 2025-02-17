package student

import (
	"fmt"
	"graid-tech-assignment/pkg/question"
	"math/rand"
	"sync"
)

type Student struct {
	Name        string
	questionMap map[*question.Question]float64
	mux         sync.Mutex
}

func NewStudent(name string) *Student {
	s := &Student{
		Name:        name,
		questionMap: make(map[*question.Question]float64),
		mux:         sync.Mutex{},
	}
	return s
}
func (s *Student) Say(message string) {
	fmt.Printf("Student %s: %s\n", s.Name, message)
}

func (s *Student) SayAnswer(q *question.Question) {
	s.Say(q.GuessedAnswerString(s.questionMap[q]))
}

func (s *Student) LookupQuestion(q *question.Question, tooEasy bool) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if tooEasy || rand.Intn(2) == 1 {
		s.questionMap[q] = q.GetAnswer()
	} else {
		s.questionMap[q] = float64(rand.Int31n(100))
	}
}

func (s *Student) GuessAnswer(q *question.Question) float64 {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.questionMap[q]
}

func (s *Student) Congratulate(q *question.Question, other *Student) {
	if len(q.Name) > 0 {
		s.Say(fmt.Sprintf("%s: %s you win.", q.Name, other.Name))
		return
	}
	s.Say(fmt.Sprintf("%s you win.", other.Name))
}
