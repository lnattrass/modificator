package api

type PutRequestOptions struct {
	File     []byte `json:"file"`
	Message  string `json:"commit_message"`
	CreatePR bool   `json:"create_pr"`
	MergePR  bool   `json:"merge_pr"`

	Token string `json:"token"`
}

type PatchRequestOptions struct {
	Message string `json:"message"`

	Patch []byte `json:"patch"`

	CreatePR bool `json:"create_pr"`
	MergePR  bool `json:"merge_pr"`

	Token string `json:"token"`
}
