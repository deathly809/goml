package classify

import (
	"math/rand"

	"github.com/deathly809/gomath"
	"github.com/deathly809/gotypes"
)

type svm struct {
	alpha          []float32
	supportVectors []Data
	kernel         Kernel
}

func (s *svm) Classify(data []gotypes.Value) float32 {
	panic("unimplemented")
}

type smoState struct {
	// Parameters
	tolerance float32
	eps       float32
	C         float32

	// State
	E                    []float32
	threshold            float32
	numNonBoundaryAlphas int

	// Current build up
	a1, a2 float32
	v1, v2 []gotypes.Value
	E1, E2 float32
	y1, y2 float32

	// Misc
	wq           WorkQueue
	workFunction func(int)
}

func takeStep(s *svm, i1, i2 int, state *smoState) bool {
	result := false
	if i1 != i2 {

		fZero := float32(0)

		alpha1 := s.alpha[i1]
		alpha2 := state.a2

		state.a1 = alpha1

		state.v1 = s.supportVectors[i1].Value()
		state.y1 = s.supportVectors[i1].Class()

		state.E1 = s.Classify(state.v1) - state.y1

		var L, H float32

		if state.y1 == state.y2 {
			sum := alpha2 + alpha1
			L = gomath.MaxFloat32(fZero, sum-state.C)
			H = gomath.MinFloat32(state.C, sum)
		} else {
			diff := alpha2 - alpha1
			L = gomath.MaxFloat32(fZero, diff)
			H = gomath.MinFloat32(state.C, state.C+diff)
		}

		if L != H {
			k11 := s.kernel.Dot(state.v1, state.v1)
			k12 := s.kernel.Dot(state.v1, state.v2)
			k22 := s.kernel.Dot(state.v2, state.v2)

			eta := 2*k12 - k11 - k22

			if eta < 0 {
				state.a2 = gomath.ClampFloat32(L, H, alpha1-state.y2*(state.E1-state.E2)/eta)
			} else {
				Lobj := float32(5.0)
				Hobj := float32(5.0)
				if Lobj > Hobj+1E-3 {
					state.a2 = L
				} else if Lobj < Hobj-1E-3 {
					state.a2 = L
				} else {
					state.a2 = alpha2
				}

				if state.a2 < 1E-8 {
					state.a2 = 0
				} else if state.a2 > state.C-1e-8 {
					state.a2 = state.C
				}

				if gomath.AbsFloat32(state.a2-alpha2) >= state.eps*(state.a2+alpha2+state.eps) {
					s.alpha[i1] = state.a1
					s.alpha[i2] = state.a2

					// Update error
					state.wq.Enqueue(0, len(s.alpha), state.workFunction)

					result = true
				}
			}
		}
	}
	return result
}

// Two heuristics:
func examineExample(s *svm, target int, state *smoState) int {

	data := s.supportVectors[target]
	a2 := s.alpha[target]
	E2 := s.Classify(data.Value()) - data.Class()
	r2 := E2 * data.Class()
	if (r2 < -state.tolerance && a2 < state.C) || (r2 > state.tolerance && a2 > 0) {
		if state.numNonBoundaryAlphas > 1 {
			// Second Choice Heuristic
			i := 9
			if takeStep(s, i, target, state) {
				return 1
			}
		}

		for i, a := range s.alpha {
			if a > 0 && a < state.C {
				if takeStep(s, i, target, state) {
					return 1
				}
			}
		}

		randomStart := rand.Intn(len(s.alpha))
		for i := 0; i < len(s.alpha); i++ {
			if takeStep(s, randomStart, target, state) {
				return 1
			}
			randomStart = (randomStart + 1) % len(s.alpha)
		}
	}

	return 0
}

/*
   Fast Training of Support Vector Machines
   using Sequential Minimal Optimization

   max: W(alpha) = Sum_{i=0}^l alpha_i - \frac{1}{2} Sum_{i=0}^l Sum_{j=0}^l y_{i}y_{j}K(x_{i},x_{j})alpha_{i}alpha_{j}

   s.t.

   0 \lte alpha_i \lte C,
   Sum_{i=0}^l y_{i}*alpha_{i}

*/
func smo(data []Data, k Kernel, tol, C float32) *svm {

	result := &svm{
		alpha:          make([]float32, len(data)),
		supportVectors: data,
	}

	state := &smoState{
		C:                    C,
		numNonBoundaryAlphas: 0,
		tolerance:            tol,
		threshold:            float32(0.0),
		E:                    make([]float32, len(data)),
	}

	state.wq.Init(10)
	state.wq.Start()
	defer state.wq.Stop()

	state.workFunction = func(i int) {
		v := result.supportVectors[i]
		state.E[i] = result.Classify(v.Value()) - v.Class()
	}

	numChanged := 0
	examineAll := true
	for (numChanged > 0) || examineAll {
		numChanged = 0
		if examineAll {
			for i := range data {
				numChanged += examineExample(result, i, state)
			}
			examineAll = !examineAll
		} else {
			for i := range data {
				a := result.alpha[i]
				if a > 0 && a < state.threshold {
					numChanged += examineExample(result, i, state)
				}
			}
		}

		if numChanged == 0 {
			examineAll = true
		}
	}

	// Compress
	pos := 0
	for i := range data {
		a := result.alpha[i]
		if a > 0 {
			result.alpha[pos] = a
			result.supportVectors[pos] = data[i]
			pos++
		}
	}

	result.alpha = append(([]float32)(nil), result.alpha[:pos]...)
	result.supportVectors = append(([]Data)(nil), result.supportVectors[:pos]...)

	return result
}

// NewSVMClassifer will create a new SVM classifier with the data given
func NewSVMClassifer(data []Data, k Kernel, alpha float32, C float32) Classifier {
	return smo(data, k, alpha, C)
}
