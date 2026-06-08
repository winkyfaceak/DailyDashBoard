package daily

import "fmt"

type MathQuestion struct {
	Left     int
	Right    int
	Operator string
	Answer   int
}

func RandomMathQuestion() (MathQuestion, error) {
	operation, err := randomInt(4)
	if err != nil {
		return MathQuestion{}, err
	}

	switch operation {
	case 0:
		left, err := randomMathNumber(10, 89)
		if err != nil {
			return MathQuestion{}, err
		}

		right, err := randomMathNumber(10, 89)
		if err != nil {
			return MathQuestion{}, err
		}

		return MathQuestion{
			Left:     left,
			Right:    right,
			Operator: "+",
			Answer:   left + right,
		}, nil
	case 1:
		left, err := randomMathNumber(20, 99)
		if err != nil {
			return MathQuestion{}, err
		}

		right, err := randomMathNumber(1, left)
		if err != nil {
			return MathQuestion{}, err
		}

		return MathQuestion{
			Left:     left,
			Right:    right,
			Operator: "-",
			Answer:   left - right,
		}, nil
	case 2:
		left, err := randomMathNumber(2, 12)
		if err != nil {
			return MathQuestion{}, err
		}

		right, err := randomMathNumber(2, 12)
		if err != nil {
			return MathQuestion{}, err
		}

		return MathQuestion{
			Left:     left,
			Right:    right,
			Operator: "x",
			Answer:   left * right,
		}, nil
	default:
		right, err := randomMathNumber(2, 12)
		if err != nil {
			return MathQuestion{}, err
		}

		answer, err := randomMathNumber(2, 12)
		if err != nil {
			return MathQuestion{}, err
		}

		return MathQuestion{
			Left:     right * answer,
			Right:    right,
			Operator: "/",
			Answer:   answer,
		}, nil
	}
}

func randomMathNumber(minimum int, maximum int) (int, error) {
	if maximum < minimum {
		return 0, fmt.Errorf("math range maximum must be at least minimum")
	}

	value, err := randomInt(maximum - minimum + 1)
	if err != nil {
		return 0, err
	}

	return minimum + value, nil
}

func (q MathQuestion) Prompt() string {
	return fmt.Sprintf("%d %s %d = ?", q.Left, q.Operator, q.Right)
}

func (q MathQuestion) CheckAnswer(answer int) bool {
	return answer == q.Answer
}
