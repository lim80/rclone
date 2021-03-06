// Package api has type definitions for pcloud
//
// Converted from the API docs with help from https://mholt.github.io/json-to-go/
package api

import (
	"fmt"
	"time"
)

const (
	// Sun, 16 Mar 2014 17:26:04 +0000
	timeFormat = `"` + time.RFC1123Z + `"`
)

// Time represents represents date and time information for the
// pcloud API, by using RFC1123Z
type Time time.Time

// MarshalJSON turns a Time into JSON (in UTC)
func (t *Time) MarshalJSON() (out []byte, err error) {
	timeString := (*time.Time)(t).Format(timeFormat)
	return []byte(timeString), nil
}

// UnmarshalJSON turns JSON into a Time
func (t *Time) UnmarshalJSON(data []byte) error {
	newT, err := time.Parse(timeFormat, string(data))
	if err != nil {
		return err
	}
	*t = Time(newT)
	return nil
}

// Error is returned from pcloud when things go wrong
//
// If result is 0 then everything is OK
type Error struct {
	Result      int    `json:"result"`
	ErrorString string `json:"error"`
}

// Error returns a string for the error and statistifes the error interface
func (e *Error) Error() string {
	return fmt.Sprintf("pcloud error: %s (%d)", e.ErrorString, e.Result)
}

// Update returns err directly if it was != nil, otherwise it returns
// an Error or nil if no error was detected
func (e *Error) Update(err error) error {
	if err != nil {
		return err
	}
	if e.Result == 0 {
		return nil
	}
	return e
}

// Check Error statisfies the error interface
var _ error = (*Error)(nil)

// Item describes a folder or a file as returned by Get Folder Items and others
type Item struct {
	Path           string `json:"path"`
	Name           string `json:"name"`
	Created        Time   `json:"created"`
	IsMine         bool   `json:"ismine"`
	Thumb          bool   `json:"thumb"`
	Modified       Time   `json:"modified"`
	Comments       int    `json:"comments"`
	ID             string `json:"id"`
	IsShared       bool   `json:"isshared"`
	IsDeleted      bool   `json:"isdeleted"`
	Icon           string `json:"icon"`
	IsFolder       bool   `json:"isfolder"`
	ParentFolderID int64  `json:"parentfolderid"`
	FolderID       int64  `json:"folderid,omitempty"`
	Height         int    `json:"height,omitempty"`
	FileID         int64  `json:"fileid,omitempty"`
	Width          int    `json:"width,omitempty"`
	Hash           uint64 `json:"hash,omitempty"`
	Category       int    `json:"category,omitempty"`
	Size           int64  `json:"size,omitempty"`
	ContentType    string `json:"contenttype,omitempty"`
	Contents       []Item `json:"contents"`
}

// ModTime returns the modification time of the item
func (i *Item) ModTime() (t time.Time) {
	t = time.Time(i.Modified)
	if t.IsZero() {
		t = time.Time(i.Created)
	}
	return t
}

// ItemResult is returned from the /listfolder, /createfolder, /deletefolder, /deletefile etc methods
type ItemResult struct {
	Error
	Metadata Item `json:"metadata"`
}

// Hashes contains the supported hashes
type Hashes struct {
	SHA1 string `json:"sha1"`
	MD5  string `json:"md5"`
}

// UploadFileResponse is the response from /uploadfile
type UploadFileResponse struct {
	Error
	Items     []Item   `json:"metadata"`
	Checksums []Hashes `json:"checksums"`
	Fileids   []int64  `json:"fileids"`
}

// GetFileLinkResult is returned from /getfilelink
type GetFileLinkResult struct {
	Error
	Dwltag  string   `json:"dwltag"`
	Hash    uint64   `json:"hash"`
	Size    int64    `json:"size"`
	Expires Time     `json:"expires"`
	Path    string   `json:"path"`
	Hosts   []string `json:"hosts"`
}

// IsValid returns whether the link is valid and has not expired
func (g *GetFileLinkResult) IsValid() bool {
	if g == nil {
		return false
	}
	if len(g.Hosts) == 0 {
		return false
	}
	return time.Time(g.Expires).Sub(time.Now()) > 30*time.Second
}

// URL returns a URL from the Path and Hosts.  Check with IsValid
// before calling.
func (g *GetFileLinkResult) URL() string {
	// FIXME rotate the hosts?
	return "https://" + g.Hosts[0] + g.Path
}

// ChecksumFileResult is returned from /checksumfile
type ChecksumFileResult struct {
	Error
	Hashes
	Metadata Item `json:"metadata"`
}
