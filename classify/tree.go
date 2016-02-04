package classify

// Decision Tree

type dTree struct {
}

const (
	// InformationGain means that to build the Decision Tree we
	// use information gain to determine the split.  This will
	// use the C4.5 algorithm.
	InformationGain = iota
	// GiniImpurity is a measure of how often a randomly chosen
	// element from the set would be incorrectly labeled (Wikipedia)
	GiniImpurity = iota
	// Variance is for continuous target values
	Variance = iota
	// CART means Classfiication and Regression Tree
	CART = iota
	// Boost means we use boosting
	Boost = iota
	//
)

// DecisionTreeParameters is the parameters to generate
// a decision tree
type DecisionTreeParameters struct {
}

func (dt *dTree) Classify([]float32) float32 {
	panic("Unimplemented")
}

// NewDecisionTree constructs a new decision tree
func NewDecisionTree(data []Data, param DecisionTreeParameters) Classifier {
	return &dTree{}
}
