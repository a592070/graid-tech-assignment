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

func main() {
	t := teacher.NewTeacher()
	studentA := student.NewStudent("A")
	studentB := student.NewStudent("B")
	studentC := student.NewStudent("C")
	studentD := student.NewStudent("D")
	studentE := student.NewStudent("E")
	students := []*student.Student{
		studentA, studentB, studentC, studentD, studentE,
	}

	for {
		a := rand.Int31n(100)
		b := rand.Int31n(100)
		operations := []string{"+", "-", "*", "/"}
		operation := operations[rand.Intn(len(operations))]
		q, err := question.NewQuestion("", a, b, operation)
		if err != nil {
			if errors.Is(err, question.InvalidInputDivisionBy0) {
				continue
			}
			log.Fatalln(err)
		}

		fmt.Println("Thinking question...")
		time.Sleep(3 * time.Second)
		t.SayAskingQuestion(q)

		selectedStudentIdx := rand.Intn(len(students))
		raisingHandStudent := students[selectedStudentIdx]

		fmt.Println("Thinking answer...")
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)

		raisingHandStudent.LookupQuestion(q, true)
		raisingHandStudent.SayGuessAnswer(q)
		isCorrect := q.IsCorrect(raisingHandStudent.GuessAnswer(q))
		t.SayResponseToGuessAnswer(q, raisingHandStudent)
		if isCorrect {
			for i, s := range students {
				if i == selectedStudentIdx {
					continue
				}
				s.Congratulate(q, raisingHandStudent)
			}
		}

		fmt.Println("==========================")
	}

}
