Write a program that allows a teacher and students named A, B, C, D, and E to play math questions. 

1. Run the following steps in a loop.
2. Teacher behavior: Ask a math question.
    1. Warm up for 3 seconds. 
    2. Ask a math question.
        1. Randomly generate "A c B".
        2. A and B are integers between 0 and 100. 3. c is a mathematical symbol (+, -, *, /).
   3. Wait for the answer, and say “XXX, you are right!” 
3. Student behavior:
   1. Wait for the question.
   2. See the question and think (randomly between 1 and 3 seconds).
   3. Raise hand and answer the question (only one student can answer) (assuming they are always correct).
   4. The other students may feel sad and say "XXX, you win".
4. Example output:
    ```text
    Teacher: Guys, are you ready?
    # Count 3
    Teacher: 1 + 1 = ?
    # May have a few seconds
    Student C: 1 + 1 = 2!
    Teacher: C, you are right!
    Student A: C, you win.
    Student B: C, you win.
    Student D: C, you win.
    Student E: C, you win.
    ```

5. Bonus
   1. Students may have wrong answer, the other students can try to raise hand and answer the question.
      ```text
      Teacher: Guys, are you ready?
      # Count 3
      Teacher: 1 + 1 = ?
      # May have a few seconds
      Student C: 1 + 1 = 3!
      Teacher: C, you are wrong.
      Student A: 1 + 1 = 4!
      Teacher: A, you are wrong.
      Student B: 1 + 1 = 2!
      Teacher: B, you are right!
      Student A: B, you win.
      Student C: B, you win.
      Student D: C, you win.
      Student E: C, you win.
      # Or there is no student has right answer.
      # (all 5 students have wrong answer).
      # Teacher feels sad and say the answer.
      Teacher: Boooo~ Answer is 2.
      ```
   2. Teacher writes the question on a board per second (means that teacher would not wait for students answer the question). Every questions is a independent process.
      ```text
      Teacher: Guys, are you ready?
      # Count 3
      Teacher: Q1: 1 + 1 = ?
      Student C: Q1: 1 + 1 = 2!
      Teacher: C, Q1 you are right!
      Student A: C, Q1 you win.
      # teacher ask 2nd question
      Teacher: Q2: 3 + 3 = ?
      # students can answer 2nd question
      # although teacher doesn't confirm 1st answer yet
      Student E: Q2: 3 + 3 = 6!
      Student B: C, Q1 you win.
      Student A: E, Q2 you win.
      Student B: E, Q2 you win.
      Student D: C, Q1 you win.
      Student D: E, Q2 you win.
      Student E: C, Q1 you win.
      Student C: E, Q2 you win.
      ```
