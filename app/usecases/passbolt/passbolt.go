package passbolt

import (
	"context"
	"net/http"
	"os"

	"github.com/lenforiee/AmnesiaGUI/app/models"
	"github.com/passbolt/go-passbolt/helper"

	amnesiaApp "github.com/lenforiee/AmnesiaGUI/app"
	"github.com/lenforiee/AmnesiaGUI/app/mfa"
	"github.com/passbolt/go-passbolt/api"
)

func InitialisePassboltConnector(ctx amnesiaApp.AppContext, password string) error {
	// read the private key file
	privateKey, err := os.ReadFile(ctx.UserSettings.PrivateKeyPath)
	if err != nil {
		return err
	}

	client, err := api.NewClient(
		nil, ctx.UserSettings.UserAgent, ctx.UserSettings.ServerURI, string(privateKey), password,
	)

	client.MFACallback = func(respCtx context.Context, c *api.Client, res *api.APIResponse) (http.Cookie, error) {

		// Use channels to tranmit the cookie data as passbolt mfa callback is bit weird.
		mfaChan := make(chan http.Cookie)
		errChan := make(chan error)

		mfaView := mfa.NewMFAView(ctx, respCtx, c, res, mfaChan, errChan)
		mfaView.Window.Show()

		select {
		case mfaCookie := <-mfaChan:
			mfaView.Window.Close()
			return mfaCookie, nil
		case err := <-errChan:
			mfaView.Window.Close()
			return http.Cookie{}, err
		}
	}

	if err != nil {
		return err
	}

	err = client.Login(ctx.Context)
	if err != nil {
		return err
	}

	ctx.SetPassboltClient(client)
	return nil
}

func GetResources(ctx amnesiaApp.AppContext, opts api.GetResourcesOptions) ([]api.Resource, error) {
	resources, err := ctx.PassboltClient.GetResources(ctx.Context, &opts)
	if err != nil {
		return nil, err
	}

	return resources, nil
}

func GetResource(ctx amnesiaApp.AppContext, id string) (models.Resource, error) {
	resource := models.NewResource()

	folderId, name, username, uri, password, desc, err := helper.GetResource(
		ctx.Context, ctx.PassboltClient, id,
	)

	if err != nil {
		return resource, err
	}

	resource.SetFolderParentID(folderId)
	resource.SetName(name)
	resource.SetUsername(username)
	resource.SetURI(uri)
	resource.SetPassword(password)
	resource.SetDescription(desc)

	return resource, nil
}

func CreateResource(ctx amnesiaApp.AppContext, resource models.Resource) error {
	_, err := helper.CreateResource(
		ctx.Context,
		ctx.PassboltClient,
		resource.FolderParentID,
		resource.Name,
		resource.Username,
		resource.URI,
		resource.Password,
		resource.Description,
	)
	if err != nil {
		return err
	}

	return nil
}

func UpdateResource(ctx amnesiaApp.AppContext, id string, resource models.Resource) error {
	err := helper.UpdateResource(
		ctx.Context,
		ctx.PassboltClient,
		id,
		resource.Name,
		resource.Username,
		resource.URI,
		resource.Password,
		resource.Description,
	)
	if err != nil {
		return err
	}

	return nil
}

func DeleteResource(ctx amnesiaApp.AppContext, id string) error {
	err := helper.DeleteResource(ctx.Context, ctx.PassboltClient, id)
	if err != nil {
		return err
	}

	return nil
}
