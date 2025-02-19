package main

import (
	"errors"
	"fmt"
	"graid-tech-assignment/pkg/task1/question"
	"graid-tech-assignment/pkg/task1/student"
	"graid-tech-assignment/pkg/task1/teacher"
	"log"
	"math/rand"
	"sync"
	"time"
)

func generateQuestion(name string) (*question.Question, error) {
	a := rand.Int31n(100)
	b := rand.Int31n(100)
	operations := []string{"+", "-", "*", "/"}
	operation := operations[rand.Intn(len(operations))]
	return question.NewQuestion(name, a, b, operation)
}

func main() {
	t := teacher.NewTeacher()
	studentA := student.NewStudent("A")
	studentB := student.NewStudent("B")
	studentC := student.NewStudent("C")
	studentD := student.NewStudent("D")
	studentE := student.NewStudent("E")
	students := []*student.Student{
		studentA, studentB,
		studentC, studentD, studentE,
	}

	fmt.Println("Thinking question...")
	time.Sleep(3 * time.Second)

	questionCount := 1
	wg := sync.WaitGroup{}
	for questionCount < 10 {
		q, err := generateQuestion(fmt.Sprintf("Q%d", questionCount))
		if err != nil {
			if errors.Is(err, question.InvalidInputDivisionBy0) {
				continue
			}
			log.Fatalln(err)
		}
		t.SayAskingQuestion(q)

		wg.Add(1)
		go func(q *question.Question) {
			defer wg.Done()
			//fmt.Printf("%s, Thinking answer...\n", q.Name)
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second)

			answeredStudentMap := map[*student.Student]bool{}
			noCorrectAnswer := true
			for len(answeredStudentMap) < len(students) {
				selectedStudentIdx := rand.Intn(len(students))
				raisingHandStudent := students[selectedStudentIdx]

				if answeredStudentMap[raisingHandStudent] {
					continue
				}
				answeredStudentMap[raisingHandStudent] = true

				raisingHandStudent.LookupQuestion(q, false)
				raisingHandStudent.SayGuessAnswer(q)
				isCorrect := q.IsCorrect(raisingHandStudent.GuessAnswer(q))
				t.SayResponseToGuessAnswer(q, raisingHandStudent)
				if isCorrect {
					noCorrectAnswer = false
					for i, s := range students {
						if i == selectedStudentIdx {
							continue
						}
						s.Congratulate(q, raisingHandStudent)
					}
					break
				}
			}
			if noCorrectAnswer {
				t.SayResponseToNoCorrectAnswer(q)
			}
		}(q)
		questionCount++
		time.Sleep(time.Second)
	}
	wg.Wait()
}
