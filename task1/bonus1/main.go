package main

import (
	"errors"
	"fmt"
	"graid-tech-assignment/pkg/task1/question"
	"graid-tech-assignment/pkg/task1/student"
	"graid-tech-assignment/pkg/task1/teacher"
	"log"
	"math/rand"
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

	for {
		q, err := generateQuestion("")
		if err != nil {
			if errors.Is(err, question.InvalidInputDivisionBy0) {
				continue
			}
			log.Fatalln(err)
		}

		fmt.Println("Thinking question...")
		time.Sleep(3 * time.Second)
		t.SayAskingQuestion(q)

		fmt.Println("Thinking answer...")
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

		fmt.Println("==========================")
	}
}
