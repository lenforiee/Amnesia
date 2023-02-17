package passbolt

import (
	"github.com/lenforiee/AmnesiaGUI/internals/controllers"
	"github.com/lenforiee/AmnesiaGUI/models"
	"github.com/passbolt/go-passbolt/api"
	"github.com/passbolt/go-passbolt/helper"
)

func GetResources(app *controllers.AppContext, opts api.GetResourcesOptions) ([]api.Resource, error) {
	resources, err := app.PassboltClient.GetResources(*app.Context, &opts)
	if err != nil {
		return nil, err
	}

	return resources, nil
}

func GetResource(app *controllers.AppContext, id string) (*models.Resource, error) {
	resource, err := models.NewResource(helper.GetResource(*app.Context, app.PassboltClient, id))
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func CreateResource(app *controllers.AppContext, resource *models.Resource) error {
	_, err := helper.CreateResource(*app.Context, app.PassboltClient, resource.FolderParentID, resource.Name, resource.Username, resource.URI, resource.Password, resource.Description)
	if err != nil {
		return err
	}

	return nil
}

func UpdateResource(app *controllers.AppContext, id string, resource *models.Resource) error {
	err := helper.UpdateResource(*app.Context, app.PassboltClient, id, resource.Name, resource.Username, resource.URI, resource.Password, resource.Description)
	if err != nil {
		return err
	}

	return nil
}

func DeleteResource(app *controllers.AppContext, id string) error {
	err := helper.DeleteResource(*app.Context, app.PassboltClient, id)
	if err != nil {
		return err
	}

	return nil
}
