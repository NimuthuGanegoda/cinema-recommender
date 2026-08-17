package main

import "fmt"

// Student demonstrates struct composition and methods
type Student struct {
	Name  string
	ID    int
	Marks []float64
}

// CalculateAverage calculates the mean grade
func (s Student) CalculateAverage() float64 {
	if len(s.Marks) == 0 {
		return 0.0
	}
	total := 0.0
	for _, m := range s.Marks {
		total += m
	}
	return total / float64(len(s.Marks))
}

// AddMark demonstrates pointer receiver method mutating struct fields
func (s *Student) AddMark(mark float64) {
	s.Marks = append(s.Marks, mark)
}

func main() {
	s := Student{Name: "Nimuthu", ID: 10705257, Marks: []float64{82.5, 88.0, 91.5}}
	fmt.Printf("Student: %s, ID: %d, Initial Average: %.2f\n", s.Name, s.ID, s.CalculateAverage())

	s.AddMark(95.0)
	fmt.Printf("After adding mark: New Average: %.2f (Total Marks: %d)\n", s.CalculateAverage(), len(s.Marks))
}
