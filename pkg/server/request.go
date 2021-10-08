package server

type CreateWorkspacesRequestBody struct {
	WorkspaceName string `json:"workspaceName" validate:"required,min=2,max=250"`
}
