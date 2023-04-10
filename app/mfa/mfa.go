package mfa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	amnesiaApp "github.com/lenforiee/AmnesiaGUI/app"
	"github.com/passbolt/go-passbolt/api"
)

// MFA is bit of a schmuck and needs to be in separate package
// to avoid import cycle.

type MFAView struct {
	Window    fyne.Window
	Container *fyne.Container
}

func NewMFAView(
	ctx *amnesiaApp.AppContext,
	ctxResp context.Context,
	c *api.Client,
	res *api.APIResponse,
	mfaChan chan http.Cookie,
	mfaErr chan error,
) MFAView {

	window := ctx.App.NewWindow(fmt.Sprintf("%s :: MFA", ctx.AppName))
	view := MFAView{
		Window: window,
	}

	mfaLabel := widget.NewLabelWithStyle(
		"Enter MFA Code",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	mfaText := widget.NewLabelWithStyle(
		"Enter the six digit number as presented on your phone or tablet.",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	mfaText.Wrapping = fyne.TextWrapWord

	mfaCode := widget.NewEntry()
	mfaCode.SetPlaceHolder("eg. 794732")

	rememberCheck := widget.NewCheckWithData(
		"Remember this device for a month.",
		binding.NewBool(),
	)

	mfaBtn := widget.NewButton("Submit", func() {
		VerifyOTPCode(ctxResp, c, res, mfaCode.Text, rememberCheck.Checked, mfaChan, mfaErr)
	})

	containerBox := container.NewBorder(
		mfaLabel,
		container.New(
			layout.NewVBoxLayout(),
			mfaCode,
			rememberCheck,
			widget.NewSeparator(),
			mfaBtn,
		),
		nil,
		nil,
		container.New(
			layout.NewVBoxLayout(),
			mfaText,
		),
	)
	view.Container = containerBox

	view.Window.SetContent(view.Container)
	view.Window.Resize(fyne.NewSize(300, 200))
	view.Window.CenterOnScreen()
	return view
}

type MFAChallangeResponse struct {
	TOTP     string `json:"totp,omitempty"`
	Remember bool   `json:"remember,omitempty"`
}

func VerifyOTPCode(
	ctx context.Context,
	c *api.Client,
	res *api.APIResponse,
	code string,
	remember bool,
	mfaChan chan http.Cookie,
	mfaErr chan error,
) {
	challange := api.MFAChallange{}
	err := json.Unmarshal(res.Body, &challange)
	if err != nil {
		mfaErr <- err
		return
	}
	if challange.Provider.TOTP == "" {
		mfaErr <- fmt.Errorf("server Provided no TOTP Provider")
		return
	}

	req := MFAChallangeResponse{
		TOTP:     code,
		Remember: remember,
	}
	var raw *http.Response
	raw, _, err = c.DoCustomRequestAndReturnRawResponse(ctx, "POST", "mfa/verify/totp.json", "v2", req, nil)
	if err != nil {
		mfaErr <- err
		return
	}

	// MFA worked so lets find the cookie and return it
	for _, cookie := range raw.Cookies() {
		if cookie.Name == "passbolt_mfa" {
			mfaChan <- *cookie
			return
		}
	}
	mfaErr <- fmt.Errorf("unable to find Passbolt MFA Cookie")
}
