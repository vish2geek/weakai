package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/unixpickle/weakai/boosting"
)

const Dimensions = 5
const Samples = 80

func main() {
	rand.Seed(time.Now().UnixNano())

	problem := &boosting.Problem{}
	for i := 0; i < Samples; i++ {
		positive := (rand.Intn(2) == 1)
		sample := make([]float64, Dimensions)
		for d := 0; d < Dimensions; d++ {
			sample[d] = rand.Float64()
		}
		problem.Samples = append(problem.Samples, sample)
		problem.Classifications = append(problem.Classifications, positive)
	}
	stumps := boosting.AllTreeStumps(problem.Samples, Dimensions)
	for _, s := range stumps {
		problem.Classifiers = append(problem.Classifiers, s)
	}

	solve("AdaBoost", boosting.AdaboostSolver{MaxReuse: 10}, problem)
	solve("Gradient", &boosting.GradientSolver{
		Iterations: 1000,
		StepSize:   0.01,
		Attempts:   5,
	}, problem)

	problem = pathologicalProblem()
	solve("AdaBoost (pathological)", boosting.AdaboostSolver{MaxReuse: 10}, problem)
	solve("Gradient (pathological)", &boosting.GradientSolver{
		Iterations: 1000,
		StepSize:   0.01,
		Attempts:   5,
	}, problem)
}

func solve(name string, solver boosting.Solver, p *boosting.Problem) {
	solution := solver.Solve(p)
	fmt.Println("Using", name, "we get", errorCount(solution, p), "wrong.")
}

func pathologicalProblem() *boosting.Problem {
	prob := &boosting.Problem{
		Samples: []boosting.Sample{[]float64{0, 0}, []float64{1, 0}, []float64{0, 1},
			[]float64{1, 1}},
		Classifications: []bool{true, false, false, true},
	}
	stumps := boosting.AllTreeStumps(prob.Samples, 2)
	for _, s := range stumps {
		prob.Classifiers = append(prob.Classifiers, s)
	}
	return prob
}

func errorCount(s *boosting.Solution, p *boosting.Problem) int {
	var count int
	for i, sam := range p.Samples {
		if s.Classify(sam) != p.Classifications[i] {
			count++
		}
	}
	return count
}
