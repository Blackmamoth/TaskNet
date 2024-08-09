package types

import script_model "github.com/blackmamoth/tasknet/pkg/models/script"

type ScriptService interface {
	CreateScript(*script_model.Script) error
	GetScriptByName(name string) (*script_model.Script, error)
}

type ScriptRepository interface {
	CreateScript(*script_model.Script) error
	GetScriptByName(name string) (*script_model.Script, error)
}
