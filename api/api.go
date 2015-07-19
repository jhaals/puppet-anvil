package api

type Metadata struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Author       string `json:"author"`
	Licence      string `json:"license"`
	Dependencies []struct {
		Name               string `json:"name"`
		VersionRequirement string `json:"version_requirement,omitempty"`
	} `json:"dependencies"`
}

type Result struct {
	Uri      string `json:"uri"`
	FileUri  string `json:"file_uri"`
	Version  string `json:"version"`
	Md5      string `json:"file_md5"`
	Metadata `json:"metadata"`
}
type Pagination struct {
	Next bool `json:"next"`
}
type Response struct {
	Results    []Result `json:"results"`
	Pagination `json:"pagination"`
}
