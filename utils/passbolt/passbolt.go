package passbolt

import (
	"github.com/lenforiee/PassboltGUI/internals/controllers"
	"github.com/lenforiee/PassboltGUI/models"
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
