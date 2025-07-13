package service

type Service struct {
	Task
}

type Repo struct {
	TaskRepo taskRepo
}

func NewService(repos Repo) *Service {
	return &Service{
		Task: *NewTaskService(repos.TaskRepo),
	}
}
