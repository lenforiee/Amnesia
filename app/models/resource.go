package models

// Custom resource model because passbolt wrapper has all data
// all over the place. This is a bit more organized.
type Resource struct {
	FolderParentID string `json:"folder_parent_id"`
	Name           string `json:"name"`
	Username       string `json:"username"`
	URI            string `json:"uri"`
	Password       string `json:"password"`
	Description    string `json:"description"`
}

func NewResource() Resource {
	return Resource{}
}

func (r *Resource) SetFolderParentID(id string) {
	r.FolderParentID = id
}

func (r *Resource) SetName(name string) {
	r.Name = name
}

func (r *Resource) SetUsername(username string) {
	r.Username = username
}

func (r *Resource) SetURI(uri string) {
	r.URI = uri
}

func (r *Resource) SetPassword(password string) {
	r.Password = password
}

func (r *Resource) SetDescription(description string) {
	r.Description = description
}
