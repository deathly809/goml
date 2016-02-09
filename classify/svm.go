package classify

import (
	"math/rand"

	"github.com/deathly809/gomath"
	"github.com/deathly809/gotypes"
)

type smo struct {
	s      *svm
	errors []float32

	i1, i2 int

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
	wq           *WorkQueue
	workFunction func(int)
}

type svm struct {
	alpha          []float32
	supportVectors []Data
	kernel         Kernel
}

func (s *svm) Classify(data []gotypes.Value) float32 {
	result := float32(0.0)
	for i, sV := range s.supportVectors {
		alpha := s.alpha[i]
		if s.kernel == nil {
			panic("kernel is nil")
		}
		dot := s.kernel.Dot(sV.Value(), data)
		result += alpha * dot
	}
	return result
}

func (alg *smo) takeStep() bool {
	result := false
	if alg.i1 != alg.i2 {

		fZero := float32(0)

		alpha1 := alg.s.alpha[alg.i1]
		alpha2 := alg.a2

		alg.a1 = alpha1

		alg.v1 = alg.s.supportVectors[alg.i1].Value()
		alg.y1 = alg.s.supportVectors[alg.i1].Class()

		alg.E1 = alg.s.Classify(alg.v1) - alg.y1

		var L, H float32

		if alg.y1 == alg.y2 {
			sum := alpha2 + alpha1
			L = gomath.MaxFloat32(fZero, sum-alg.C)
			H = gomath.MinFloat32(alg.C, sum)
		} else {
			diff := alpha2 - alpha1
			L = gomath.MaxFloat32(fZero, diff)
			H = gomath.MinFloat32(alg.C, alg.C+diff)
		}

		if L != H {
			k11 := alg.s.kernel.Dot(alg.v1, alg.v1)
			k12 := alg.s.kernel.Dot(alg.v1, alg.v2)
			k22 := alg.s.kernel.Dot(alg.v2, alg.v2)

			eta := 2*k12 - k11 - k22

			if eta < 0 {
				alg.a2 = gomath.ClampFloat32(L, H, alpha1-alg.y2*(alg.E1-alg.E2)/eta)
			} else {
				Lobj := float32(5.0)
				Hobj := float32(5.0)
				if Lobj > Hobj+1E-3 {
					alg.a2 = L
				} else if Lobj < Hobj-1E-3 {
					alg.a2 = L
				} else {
					alg.a2 = alpha2
				}

				if alg.a2 < 1E-8 {
					alg.a2 = 0
				} else if alg.a2 > alg.C-1e-8 {
					alg.a2 = alg.C
				}

				if gomath.AbsFloat32(alg.a2-alpha2) >= alg.eps*(alg.a2+alpha2+alg.eps) {
					alg.s.alpha[alg.i1] = alg.a1
					alg.s.alpha[alg.i2] = alg.a2

					result = true
				}
			}
		}
	}
	return result
}

// Two heuristics:
func (alg *smo) examineExample() int {

	data := alg.s.supportVectors[alg.i1]
	alg.a2 = alg.s.alpha[alg.i2]
	alg.E2 = alg.errors[alg.i2]
	r2 := alg.E2 * data.Class()

	if (r2 < -alg.tolerance && alg.a2 < alg.C) || (r2 > alg.tolerance && alg.a2 > 0) {

		if alg.numNonBoundaryAlphas > 1 {

			// Second Choice Heuristic by biggest |E1 - E2|
			alg.i1 = 0
			prev := gomath.AbsFloat32(alg.errors[alg.i1] - alg.E2)

			for i := 1; i < len(alg.s.supportVectors); i++ {
				test := gomath.AbsFloat32(alg.errors[i] - alg.E2)
				if test > prev {
					prev = test
					alg.i1 = i
				}
			}

			alg.E1 = prev
			if alg.takeStep() {
				return 1
			}
		}

		for i, a := range alg.s.alpha {
			if a > 0 && a < alg.C {
				alg.i1 = i
				alg.E1 = alg.errors[i]
				if alg.takeStep() {
					return 1
				}
			}
		}

		randomStart := rand.Intn(len(alg.s.alpha))
		for i := 0; i < len(alg.s.alpha); i++ {
			alg.i1 = randomStart
			alg.E1 = alg.errors[randomStart]
			if alg.takeStep() {
				return 1
			}
			randomStart = (randomStart + 1) % len(alg.s.alpha)
		}
	}

	return 0
}

func (alg *smo) computeError() {
	alg.wq.Enqueue(0, len(alg.s.alpha)-1, alg.workFunction)
}

/*
   Fast Training of Support Vector Machines
   using Sequential Minimal Optimization

   max: W(alpha) = Sum_{i=0}^l alpha_i - \frac{1}{2} Sum_{i=0}^l Sum_{j=0}^l y_{i}y_{j}K(x_{i},x_{j})alpha_{i}alpha_{j}

   s.t.

   0 \lte alpha_i \lte C,
   Sum_{i=0}^l y_{i}*alpha_{i}

*/
func solve(data []Data, k Kernel, tol, C float32) (result *svm) {

	result = &svm{
		alpha:          make([]float32, len(data)),
		supportVectors: data,
		kernel:         k,
	}

	alg := &smo{
		errors:               make([]float32, len(data)),
		C:                    C,
		numNonBoundaryAlphas: 0,
		tolerance:            tol,
		threshold:            float32(0.0),
		E:                    make([]float32, len(data)),
		s:                    result,
	}

	// Initialize work queue
	alg.wq = &WorkQueue{}
	alg.wq.Init(10)
	alg.wq.Start()
	defer alg.wq.Stop()

	alg.workFunction = func(i int) {
		v := result.supportVectors[i]
		alg.E[i] = result.Classify(v.Value()) - v.Class()
	}

	alg.computeError()

	// Begin algorithm
	numChanged := 0
	examineAll := true
	for (numChanged > 0) || examineAll {

		numChanged = 0

		if examineAll {
			for i := range data {
				alg.i2 = i
				numChanged += alg.examineExample()
			}
			examineAll = !examineAll
		} else {
			for i := range data {
				a := result.alpha[i]
				if a > 0 && a < alg.threshold {
					alg.i2 = i
					numChanged += alg.examineExample()
				}
			}
		}

		if numChanged == 0 {
			examineAll = true
		} else {
			alg.computeError()
		}
	}

	// Compress (not in algorithm)
	pos := 0
	for i := range data {
		a := result.alpha[i]
		if a > 0 {
			result.alpha[pos] = a
			result.supportVectors[pos] = data[i]
			pos++
		}
	}

	// Make sure the garbage collector can release array
	result.alpha = append(([]float32)(nil), result.alpha[:pos]...)
	result.supportVectors = append(([]Data)(nil), result.supportVectors[:pos]...)

	return result
}

// NewSVMClassifer will create a new SVM classifier with the data given
func NewSVMClassifer(data []Data, k Kernel, alpha float32, C float32) Classifier {
	return solve(data, k, alpha, C)
}
