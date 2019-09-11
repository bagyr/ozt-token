package repository_client

import "context"

const (
	EnvVal VariableType = "env_var"
	File   VariableType = "file"
)

type (
	Client interface {
		Variables(ctx context.Context, projectId int) ([]Variable, error)
		CreateVariable(ctx context.Context, projectId int, v Variable) error
		UpdateVariable(ctx context.Context, projectId int, v Variable) error
	}

	VariableType string

	Variable struct {
		VarType   VariableType `json:"variable_type"`
		Key       string       `json:"key"`
		Val       string       `json:"value"`
		Protected bool         `json:"protected"`
		Masked    bool         `json:"masked"`
		Scope     string       `json:"environment_scope"`
	}
)

func NewVar(key, val string, masked bool) Variable {
	return Variable{
		VarType: EnvVal,
		Key:     key,
		Val:     val,
		Masked:  masked,
		Scope:   `*`,
	}
}
