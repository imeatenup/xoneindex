package entity

import "time"

type Children struct {
	Odata_count int64        `json:"@odata.count"`
	Value       []*DriveItem `json:"value"`
	Error       *Error       `json:"error"`
}

type DriveItem struct {
	DownloadURL    string `json:"@microsoft.graph.downloadUrl"`
	Name           string `json:"name"`
	Size           int    `json:"size"`
	FileSystemInfo struct {
		CreatedDateTime      time.Time `json:"createdDateTime"`
		LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
	} `json:"fileSystemInfo"`
	ParentReference *struct {
		Name string `json:"name"`
		Path string `json:"path"`
	} `json:"parentReference"`
	File *struct {
		MimeType string `json:"mimeType"`
	} `json:"file"`
	Folder *struct {
		ChildCount int64 `json:"childCount"`
	} `json:"folder"`
	Error *Error `json:"error"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
