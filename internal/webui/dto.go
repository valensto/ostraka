package webui

import (
	"github.com/valensto/ostraka/internal/workflow"
)

type workflowDTO struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	NbInputs  int    `json:"nb_inputs"`
	NbOutputs int    `json:"nb_outputs"`
}

func mapWorkflowToDTO(workflows []*workflow.Workflow) []workflowDTO {
	var dtos []workflowDTO
	for _, wf := range workflows {
		dtos = append(dtos, workflowDTO{
			Name:      wf.Name,
			Slug:      wf.Slug,
			NbInputs:  len(wf.Inputs),
			NbOutputs: len(wf.Outputs),
		})
	}
	return dtos
}
