package service

type AppResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Mount   string `json:"mount"`
}

type MountRequest struct {
	Mount string `json:"mount"`
}

type InstallRequest struct {
	AppInfo string `json:"appInfo"`
	Mount   string `json:"mount"`
}

type Config map[string]interface{}

type ErrorResponse struct {
	ErrorFlag bool   `json:"error"`
	Code      int    `json:"code"`
	ErrorNum  int    `json:"errorNum"`
	Message   string `json:"errorMessage"`
}

func (er ErrorResponse) Error() string {
	return er.Message
}
