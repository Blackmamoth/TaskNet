package script_service

import (
	script_model "github.com/blackmamoth/tasknet/pkg/models/script"
	"github.com/blackmamoth/tasknet/pkg/types"
)

type Service struct {
	repository types.ScriptRepository
}

func New(repository types.ScriptRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateScript(script *script_model.Script) error {
	return s.repository.CreateScript(script)
}

func (s *Service) GetScriptByName(name string) (*script_model.Script, error) {
	return s.repository.GetScriptByName(name)
}
