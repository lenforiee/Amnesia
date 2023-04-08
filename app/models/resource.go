package models

// custom model.
type Resource struct {
	FolderParentID string `json:"folder_parent_id"`
	Name           string `json:"name"`
	Username       string `json:"username"`
	URI            string `json:"uri"`
	Password       string `json:"password"`
	Description    string `json:"description"`
}

func NewResource(
	folderId string,
	name string,
	username string,
	uri string,
	password string,
	description string,
	err error,
) (Resource, error) {
	if err != nil {
		return Resource{}, err
	}

	return Resource{
		FolderParentID: folderId,
		Name:           name,
		Username:       username,
		URI:            uri,
		Password:       password,
		Description:    description,
	}, nil
}
