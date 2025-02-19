package teacher

import (
	"fmt"
	"graid-tech-assignment/pkg/task1/question"
	"graid-tech-assignment/pkg/task1/student"
)

type Teacher struct {
}

func NewTeacher() *Teacher {
	t := &Teacher{}
	t.Say("Guys, are you ready?")
	return t
}

func (t *Teacher) Say(message string) {
	fmt.Printf("Teacher: %s\n", message)
}

func (t *Teacher) SayAskingQuestion(q *question.Question) {
	t.Say(q.QuestionString())
}
func (t *Teacher) SayResponseToGuessAnswer(q *question.Question, s *student.Student) {
	t.Say(t.RespondAnswer(q, s))
}

func (t *Teacher) RespondAnswer(q *question.Question, student *student.Student) string {
	temp := ""
	if q.IsCorrect(student.GuessAnswer(q)) {
		temp = fmt.Sprintf("%s, you are right!", student.Name)
	} else {
		temp = fmt.Sprintf("%s, you are wrong.", student.Name)
	}

	if len(q.Name) > 0 {
		return fmt.Sprintf("%s: %s", q.Name, temp)
	} else {
		return temp
	}
}

func (t *Teacher) SayResponseToNoCorrectAnswer(q *question.Question) {
	if len(q.Name) > 0 {
		t.Say(fmt.Sprintf("%s: Boooo~ Answer is %s.", q.Name, question.ConvertFloatToString(q.GetAnswer())))
	} else {
		t.Say(fmt.Sprintf("Boooo~ Answer is %s.", question.ConvertFloatToString(q.GetAnswer())))

	}

}
