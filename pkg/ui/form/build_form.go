package form

import "miko.gs/gocrud/pkg/ui/input"

type FormOptions struct {
	input.Options
	Path string
}

func BuildForm(obj interface{}, options *FormOptions) (*Form, error) {
	inputOpts := &options.Options

	inputs, order := input.BuildInputs(obj, inputOpts)

	inputMap := make(map[string]input.Input, len(inputs))
	for _, inp := range inputs {
		inputMap[inp.FieldName] = inp
	}

	return &Form{
		Path:        options.Path,
		InputsOrder: order,
		Inputs:      inputMap,
	}, nil
}
