package api_models

type IDParamRequest struct {
	ID uint `param:"id" bind:"required,min=1"`
}
