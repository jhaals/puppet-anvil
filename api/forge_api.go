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

type Owner struct {
	Uri      string `json:"uri"`
	Slug     string `json:"slug"`
	Username string `json:"username"`
}

type CurrentRelease struct {
	Uri               string `json:"uri"`
	Slug              string `json:"slug"`
	ModuleAbbreviated `json:"module"`
	Metadata          `json:"metadata"`
}

type ModuleAbbreviated struct {
	Uri   string `json:"uri"`
	Slug  string `json:"slug"`
	Name  string `json:"name"`
	Owner `json:"owner"`
}

type ModuleResult struct {
	Uri            string `json:"uri"`
	Slug           string `json:"slug"`
	Name           string `json:"name"`
	Owner          `json:"owner"`
	CurrentRelease `json:"current_release"`
	FileUri        string    `json:"file_uri"`
	Version        string    `json:"version"`
	Md5            string    `json:"file_md5"`
	Releases       []Release `json:"releases"`
}

type Release struct {
	Uri       string  `json:"uri"`
	Slug      string  `json:"slug"`
	FileUri   string  `json:"file_uri"`
	Version   string  `json:"version"`
	Deleted   *string `json:"deleted_at"`
	Created   *string `json:"created_at"`
	Supported *string `json:"supported"`
	FileSize  *string `json:"file_size"`
}

type Pagination struct {
	Next bool `json:"next"`
}
type Response struct {
	Results    []Result `json:"results"`
	Pagination `json:"pagination"`
}
