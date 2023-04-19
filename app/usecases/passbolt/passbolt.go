package passbolt

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/lenforiee/Amnesia/app/internals/logger"
	"github.com/lenforiee/Amnesia/app/internals/settings"
	"github.com/lenforiee/Amnesia/app/models"
	"github.com/lenforiee/Amnesia/app/usecases/aes"
	"github.com/passbolt/go-passbolt/helper"

	amnesiaApp "github.com/lenforiee/Amnesia/app"
	"github.com/lenforiee/Amnesia/app/mfa"
	"github.com/passbolt/go-passbolt/api"
)

func WaitForCookie(
	ctx *amnesiaApp.AppContext,
	ctxResp context.Context,
	c *api.Client,
	res *api.APIResponse,
	password string,
) (http.Cookie, error) {
	// Use channels to tranmit the cookie data as passbolt mfa callback is bit weird.
	mfaChan := make(chan http.Cookie)
	errChan := make(chan error)

	userDir := settings.GetSaveDir()

	mfaView := mfa.NewMFAView(ctx, ctxResp, c, res, mfaChan, errChan)
	mfaView.Window.Show()

	select {
	case mfaCookie := <-mfaChan:
		mfaView.Window.Close()

		if mfaCookie.Expires.IsZero() {
			return mfaCookie, nil
		}

		// Save the cookie to file only if the user wants to remember the device.
		// This is because the cookie is valid for a month.
		cookieFile, err := os.Create(fmt.Sprintf("%s/amnesia/cookie.json", userDir))
		if err != nil {
			logger.LogErr.Printf("Failed to save cookie to file: %s", err)
			return mfaCookie, nil
		}

		cookieData, err := json.Marshal(mfaCookie)
		if err != nil {
			logger.LogErr.Printf("Failed to marshal cookie data: %s", err)
			return mfaCookie, nil
		}

		encryptPasswd := password
		if len(encryptPasswd) < 32 { // pad the password to 32 bytes
			encryptPasswd = password + "AMNESIAAMNESIAAMNESIAAMNESIA"
		}

		if len(encryptPasswd) > 32 {
			encryptPasswd = encryptPasswd[:32]
		}

		encryptedCookie, iv, err := aes.Aes256Encode(cookieData, []byte(encryptPasswd))

		if err != nil {
			logger.LogErr.Printf("Failed to encrypt cookie data: %s", err)
			return mfaCookie, nil
		}

		b64Content := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s|||%s", encryptedCookie, iv)))
		_, err = cookieFile.Write([]byte(b64Content))
		if err != nil {
			logger.LogErr.Printf("Failed to write cookie data to file: %s", err)
			return mfaCookie, nil
		}

		return mfaCookie, nil
	case err := <-errChan:
		mfaView.Window.Close()
		return http.Cookie{}, err
	}
}

func EnsureLoggedIn(ctx *amnesiaApp.AppContext) {
	logger.LogInfo.Print("Checking if passbolt client is logged in...")
	if ctx.PassboltClient.CheckSession(ctx.Context) {
		return
	}

	logger.LogInfo.Print("Passbolt client is not logged in, re-authenticating...")
	err := ctx.PassboltClient.Login(ctx.Context)
	if err != nil {
		logger.LogErr.Printf("Failed to re-authenticate passbolt client: %s", err)
		return
	}
}

func InitialisePassboltConnector(ctx *amnesiaApp.AppContext, password string) error {
	// read the private key file
	privateKey, err := os.ReadFile(ctx.UserSettings.PrivateKeyPath)
	if err != nil {
		return err
	}

	client, err := api.NewClient(
		nil, ctx.UserSettings.UserAgent, ctx.UserSettings.ServerURI, string(privateKey), password,
	)

	if err != nil {
		return err
	}

	client.MFACallback = func(respCtx context.Context, c *api.Client, res *api.APIResponse) (http.Cookie, error) {

		userDir := settings.GetSaveDir()
		_, err = os.Stat(fmt.Sprintf("%s/amnesia/cookie.json", userDir))
		if os.IsNotExist(err) {
			logger.LogWarn.Print("Cookie file does not exist, waiting for user input...")
			return WaitForCookie(ctx, respCtx, c, res, password)
		}

		cookieFile, err := os.Open(fmt.Sprintf("%s/amnesia/cookie.json", userDir))
		if err != nil {
			logger.LogErr.Printf("Failed to open cookie file: %s", err)
			return WaitForCookie(ctx, respCtx, c, res, password)
		}

		cookieData, err := io.ReadAll(cookieFile)
		if err != nil {
			logger.LogErr.Printf("Failed to read cookie file: %s", err)
			return WaitForCookie(ctx, respCtx, c, res, password)
		}

		b64Content, err := b64.StdEncoding.DecodeString(string(cookieData))
		if err != nil {
			logger.LogErr.Printf("Failed to decode cookie file (B64): %s", err)
			return WaitForCookie(ctx, respCtx, c, res, password)
		}

		cookieEncrypted := strings.Split(string(b64Content), "|||")

		decryptPasswd := password
		if len(decryptPasswd) < 32 { // pad the password to 32 bytes
			decryptPasswd = password + "AMNESIAAMNESIAAMNESIAAMNESIA"
		}

		if len(decryptPasswd) > 32 {
			decryptPasswd = decryptPasswd[:32]
		}

		cookieContent, err := aes.Aes256Decode([]byte(cookieEncrypted[0]), []byte(decryptPasswd), []byte(cookieEncrypted[1]))
		if err != nil {
			logger.LogErr.Printf("Failed to decrypt cookie file (AES): %s", err)
			return WaitForCookie(ctx, respCtx, c, res, password)
		}

		var cookie http.Cookie
		err = json.Unmarshal(cookieContent, &cookie)
		if err != nil {
			logger.LogErr.Printf("Failed to unmarshal cookie file: %s", err)
			return WaitForCookie(ctx, respCtx, c, res, password)
		}

		if cookie.Expires.Unix() < time.Now().Unix() {
			os.Remove(fmt.Sprintf("%s/amnesia/cookie.json", userDir))
			logger.LogWarn.Println("Cookie expired, waiting user input...")
			return WaitForCookie(ctx, respCtx, c, res, password)
		}

		return cookie, nil
	}

	logger.LogInfo.Print("Logging in to passbolt...")
	err = client.Login(ctx.Context)
	if err != nil {
		return err
	}

	ctx.SetPassboltClient(client)
	return nil
}

func GetResources(ctx *amnesiaApp.AppContext, opts api.GetResourcesOptions) ([]api.Resource, error) {
	EnsureLoggedIn(ctx)
	logger.LogInfo.Print("Getting resources...")

	resources, err := ctx.PassboltClient.GetResources(ctx.Context, &opts)
	if err != nil {
		return nil, err
	}

	return resources, nil
}

func GetResource(ctx *amnesiaApp.AppContext, id string) (models.Resource, error) {
	EnsureLoggedIn(ctx)
	resource := models.NewResource()

	logger.LogInfo.Printf("Getting resource of id: %s...", id)
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

func CreateResource(ctx *amnesiaApp.AppContext, resource models.Resource) error {
	EnsureLoggedIn(ctx)

	logger.LogInfo.Print("Creating resource...")
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

func UpdateResource(ctx *amnesiaApp.AppContext, id string, resource models.Resource) error {
	EnsureLoggedIn(ctx)

	logger.LogInfo.Print("Updating resource...")
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

func DeleteResource(ctx *amnesiaApp.AppContext, id string) error {
	EnsureLoggedIn(ctx)

	logger.LogInfo.Print("Deleting resource...")
	err := helper.DeleteResource(ctx.Context, ctx.PassboltClient, id)
	if err != nil {
		return err
	}

	return nil
}
