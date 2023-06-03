package webui

import "github.com/valensto/ostraka/internal/workflow"

type workflowDTO struct {
	Name      string `json:"name"`
	NbInputs  int    `json:"nb_inputs"`
	NbOutputs int    `json:"nb_outputs"`
}

func mapWorkflowToDTO(workflows []*workflow.Workflow) []workflowDTO {
	var dtos []workflowDTO
	for _, wf := range workflows {
		dtos = append(dtos, workflowDTO{
			Name:      wf.Name,
			NbInputs:  len(wf.Inputs),
			NbOutputs: len(wf.Outputs),
		})
	}
	return dtos
}
