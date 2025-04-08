package onnx

import ort "github.com/yalue/onnxruntime_go"

type Model struct {
}

func NewModel() *Model {
	return &Model{}
}

func (m *Model) Load(path string) error {

	if err := ort.InitializeEnvironment(); err != nil {
		return err
	}

	// TODO destroy runtime

	return nil
}
