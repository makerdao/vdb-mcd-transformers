package test_data

type MockInitializerGenerator struct {
	GenerateFlipInitializerCalled bool
	FlipInitializerErr error
	GenerateMedianInitializerCalled bool
	MedianInitializerErr error
}

func (i *MockInitializerGenerator) GenerateFlipInitializer() error {
	i.GenerateFlipInitializerCalled = true
	return i.FlipInitializerErr
}

func (i *MockInitializerGenerator) GenerateMedianInitializer() error {
	i.GenerateMedianInitializerCalled = true
	return i.MedianInitializerErr
}
