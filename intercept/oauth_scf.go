package intercept

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/xnuc/xoneindex/config"
	"github.com/xnuc/xoneindex/log"

	"github.com/pkg/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

var Token *oauth2.Token

type OauthScf struct{}

func (i *OauthScf) PreHandle(w http.ResponseWriter, r *http.Request) bool {
	var err error
	defer func() {
		if err != nil {
			log.Errorf(r.Context(), "OauthScf.PreHandle err{%+v}", err)
		}
	}()
	log.Debugf(r.Context(), "OauthScf.PreHandle")
	err = i.untilSync(r.Context())
	if err != nil {
		return false
	}

	if r.URL.Path == config.CALLBACK { // callback update token
		r.ParseForm()
		if !(len(r.Form.Get("code")) != 0 && r.Form.Get("state") == config.STATE) {
			return false
		}
		Token, err = i.callback(r.Context(), r.Form.Get("code"))
		http.Redirect(w, r, "/", http.StatusFound)
		return false
	}

	token := os.Getenv("XINDEX_ONEDRIVE_TOKEN") // token is nil
	if token == "" {
		redirect := (&oauth2.Config{
			ClientID:    config.CLIENT_ID,
			Scopes:      config.Scopes,
			RedirectURL: config.REDIRECT_URI,
			Endpoint:    microsoft.AzureADEndpoint("consumers"),
		}).AuthCodeURL(config.STATE)
		http.Redirect(w, r, redirect, http.StatusFound)
		return false
	}

	Token = &oauth2.Token{}
	err = json.Unmarshal([]byte(token), Token)
	if err != nil {
		err = errors.WithStack(err)
		return false
	}
	if Token.Valid() { // valid token
		return true
	}

	Token, err = i.reuse(r.Context(), Token) // invalid token
	return err == nil
}

func (i *OauthScf) PostHandle(_ http.ResponseWriter, r *http.Request) bool {
	log.Debugf(r.Context(), "OauthScf.PostHandle")
	return true
}

func (i *OauthScf) untilSync(ctx context.Context) (err error) {
	secretId := config.XINDEX_TENCENTCLOUD_SECRETID
	SecretKey := config.XINDEX_TENCENTCLOUD_SECRETKEY
	credential := common.NewCredential(
		secretId,
		SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "scf.tencentcloudapi.com"
	client, _ := scf.NewClient(credential, "ap-hongkong", cpf)
	getFunctionRequest := scf.NewGetFunctionRequest()
	getFunctionRequest.FunctionName = common.StringPtr("xindex")
	var status string
	for !(status == "Active" || status == "UpdateFailed") {
		response, err := client.GetFunction(getFunctionRequest)
		if err != nil {
			return errors.WithStack(err)
		}
		status = *response.Response.Status
	}
	if status == "UpdateFailed" {
		return errors.WithStack(fmt.Errorf("UpdateFailed"))
	}
	return nil
}

func (i *OauthScf) callback(ctx context.Context, code string) (token *oauth2.Token, err error) {
	token, err = (&oauth2.Config{
		ClientID:     config.CLIENT_ID,
		ClientSecret: config.CLIENT_SECRET,
		Scopes:       config.Scopes,
		RedirectURL:  config.REDIRECT_URI,
		Endpoint:     microsoft.AzureADEndpoint("consumers"),
	}).Exchange(ctx, code)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var tokenByte []byte
	tokenByte, err = json.Marshal(token)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	tokenJSON := string(tokenByte)
	err = i.syncUpdateScfEnv(ctx, "XINDEX_ONEDRIVE_TOKEN", tokenJSON)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// syncUpdateScfEnv https://cloud.tencent.com/document/product/583/18580
func (i *OauthScf) syncUpdateScfEnv(ctx context.Context, key, val string) (err error) {
	secretId := config.XINDEX_TENCENTCLOUD_SECRETID
	SecretKey := config.XINDEX_TENCENTCLOUD_SECRETKEY
	credential := common.NewCredential(
		secretId,
		SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "scf.tencentcloudapi.com"
	client, _ := scf.NewClient(credential, "ap-hongkong", cpf)
	updateFunctionConfigurationRequest := scf.NewUpdateFunctionConfigurationRequest()
	updateFunctionConfigurationRequest.FunctionName = common.StringPtr("xindex")
	updateFunctionConfigurationRequest.Environment = &scf.Environment{
		Variables: []*scf.Variable{{
			Key:   common.StringPtr(key),
			Value: common.StringPtr(val),
		}},
	}
	_, err = client.UpdateFunctionConfiguration(updateFunctionConfigurationRequest)
	if err != nil {
		return errors.WithStack(err)
	}
	return i.untilSync(ctx)
}

func (i *OauthScf) reuse(ctx context.Context, invalid *oauth2.Token) (token *oauth2.Token, err error) {
	token, err = (&oauth2.Config{
		ClientID:     config.CLIENT_ID,
		ClientSecret: config.CLIENT_SECRET,
		Scopes:       config.Scopes,
		RedirectURL:  config.REDIRECT_URI,
		Endpoint:     microsoft.AzureADEndpoint("consumers"),
	}).TokenSource(ctx, invalid).Token()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var tokenByte []byte
	tokenByte, err = json.Marshal(token)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	tokenJSON := string(tokenByte)
	err = i.syncUpdateScfEnv(ctx, "XINDEX_ONEDRIVE_TOKEN", tokenJSON)
	if err != nil {
		return nil, err
	}
	return token, nil
}
