package layout

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"

	"github.com/xnuc/xoneindex/config"
	"github.com/xnuc/xoneindex/entity"
	"github.com/xnuc/xoneindex/intercept"
	"github.com/xnuc/xoneindex/log"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

// Index handle.Handle
type Index struct{}

func (h *Index) Handle(w http.ResponseWriter, r *http.Request) (err error) {
	defer func() {
		if err != nil {
			log.Errorf(r.Context(), "Index.Handle err{%+v}", err)
		}
	}()
	cfile, cfileerr := h.tryFile(r.Context(), r.URL.Path)
	cfolder, cfoldererr := h.tryFolder(r.Context(), r.URL.Path)

	file, err := <-cfile, <-cfileerr
	if err != nil {
		return err
	}
	if file.File != nil {
		http.Redirect(w, r, file.DownloadURL, http.StatusFound)
		return nil
	}

	folder, err := <-cfolder, <-cfoldererr
	if err != nil {
		return err
	}
	document, err := template.ParseFiles("layout/document.html", "layout/index.html")
	if err != nil {
		return errors.WithStack(err)
	}
	mp := map[string]interface{}{
		"Title":      "Index of " + r.URL.Path,
		"Path":       r.URL.Path,
		"ParentPath": file.ParentReference.Path,
		"Folder":     folder,
	}
	err = document.Execute(w, mp)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (h *Index) tryFile(ctx context.Context, uri string) (cfile chan *entity.DriveItem, cerr chan error) {
	cfile, cerr = make(chan *entity.DriveItem, 1), make(chan error, 1)
	go func() {
		file, err := h.driveItem(ctx, uri)
		if err != nil {
			cfile <- nil
			cerr <- err
			return
		}
		cfile <- file
		cerr <- nil
	}()
	return cfile, cerr
}

func (h *Index) tryFolder(ctx context.Context, uri string) (cfolder chan *entity.Children, cerr chan error) {
	cfolder, cerr = make(chan *entity.Children, 1), make(chan error, 1)
	go func() {
		folder, err := h.children(ctx, uri)
		if err != nil {
			cfolder <- nil
			cerr <- err
			return
		}
		cfolder <- folder
		cerr <- nil
	}()
	return cfolder, cerr
}

func (h *Index) children(ctx context.Context, uri string) (children *entity.Children, err error) {
	var url string
	if len(config.ROOT_PATH) == 0 && uri == "/" {
		url = "https://graph.microsoft.com/v1.0/me/drive/root/children"
	} else {
		url = "https://graph.microsoft.com/v1.0/me/drive/root:" + config.ROOT_PATH + uri + ":/children"
	}
	resp, err := (&oauth2.Config{
		ClientID:     config.CLIENT_ID,
		ClientSecret: config.CLIENT_SECRET,
		Scopes:       config.Scopes,
		RedirectURL:  config.REDIRECT_URI,
		Endpoint:     microsoft.AzureADEndpoint("consumers"),
	}).Client(ctx, intercept.Token).Get(url)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()
	children = &entity.Children{}
	json.Unmarshal(body, children)
	if children.Error != nil {
		return nil, errors.WithStack(errors.New(children.Error.Message))
	}
	for idx := range children.Value {
		children.Value[idx].ParentReference.Path = strings.TrimPrefix(children.Value[idx].ParentReference.Path,
			"/drive/root:"+config.ROOT_PATH)
	}
	return children, nil
}

func (h *Index) driveItem(ctx context.Context, uri string) (driveItem *entity.DriveItem, err error) {
	var url string
	if len(config.ROOT_PATH) == 0 && uri == "/" {
		url = "https://graph.microsoft.com/v1.0/me/drive/root"
	} else {
		url = "https://graph.microsoft.com/v1.0/me/drive/root:" + config.ROOT_PATH + uri
	}
	resp, err := (&oauth2.Config{
		ClientID:     config.CLIENT_ID,
		ClientSecret: config.CLIENT_SECRET,
		Scopes:       config.Scopes,
		RedirectURL:  config.REDIRECT_URI,
		Endpoint:     microsoft.AzureADEndpoint("consumers"),
	}).Client(ctx, intercept.Token).Get(url)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()
	driveItem = &entity.DriveItem{}
	json.Unmarshal(body, driveItem)
	if driveItem.Error != nil {
		return nil, errors.WithStack(errors.New(driveItem.Error.Message))
	}
	driveItem.ParentReference.Path = strings.TrimPrefix(driveItem.ParentReference.Path,
		"/drive/root:"+config.ROOT_PATH)
	return driveItem, nil
}
