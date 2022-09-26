package postquery

import (
	"errors"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

// isMultipartFormRequest returns true if the encoding type of
// the request is multipart/form-data, else it returns false
func isMultipartFormRequest(r *http.Request) bool {
	enctype := r.Header.Get("Content-Type")
	return strings.HasPrefix(enctype, "multipart/form-data")
}

// getValues returns the values from the request depending on encoding
// type received (MultipartForm or not)
func getValues(r *http.Request) (values url.Values) {
	if isMultipartFormRequest(r) {
		values = r.MultipartForm.Value
	} else {
		values = r.PostForm
	}
	return
}

// ParseRequestForm invokes ParseForm or ParseMultipartForm on the given request
// depending on the content-type received.
func ParseRequestForm(r *http.Request) (err error) {
	if isMultipartFormRequest(r) {
		err = r.ParseMultipartForm(1024 * 1024)
		return
	}
	err = r.ParseForm()
	if err != nil {
		return
	}
	return
}

// GetValuesByPrefix finds all keys that starts with the given prefix,
// puts each suffix as a key and the list of values as a value.
// It returns the resulting map.
func GetValuesByPrefix(r *http.Request, prefix string) (values url.Values, err error) {
	err = ParseRequestForm(r)
	if err != nil {
		return
	}
	values = make(map[string][]string)
	for k, v := range getValues(r) {
		if strings.HasPrefix(k, prefix) {
			values[strings.TrimPrefix(k, prefix)] = v
		}
	}
	return
}

// GetFileHeadersByPrefix finds all keys that starts with the given prefix,
// puts each suffix as a key and the list of file headers as a value.
// It returns the resulting map.
func GetFileHeadersByPrefix(r *http.Request, prefix string) (fileHeaders map[string][]*multipart.FileHeader, err error) {
	if !isMultipartFormRequest(r) {
		err = errors.New("Cannot get file from a request not using multipart/form-data encoding type.")
		return
	}
	err = ParseRequestForm(r)
	if err != nil {
		return
	}
	fileHeaders = make(map[string][]*multipart.FileHeader)
	for k, v := range r.MultipartForm.File {
		if strings.HasPrefix(k, prefix) {
			fileHeaders[strings.TrimPrefix(k, prefix)] = v
		}
	}
	return
}

// GetFile returns the multipart File object and its Header identified in the request by the given key.
// It can return an error if the key doesn't exist or the file cannot be retrieved.
func GetFile(r *http.Request, key string) (file multipart.File, fileHeader *multipart.FileHeader, err error) {
	if !isMultipartFormRequest(r) {
		err = errors.New("Cannot get file from a request not using multipart/form-data encoding type.")
		return
	}
	err = ParseRequestForm(r)
	if err != nil {
		return
	}
	file, fileHeader, err = r.FormFile(key)
	if err != nil {
		return
	}
	return
}

// GetAllFileHeadersByKey returns a slice of FileHeader related to the given key.
func GetAllFileHeadersByKey(r *http.Request, key string) (fileHeaders []*multipart.FileHeader, err error) {
	if !isMultipartFormRequest(r) {
		err = errors.New("Cannot get files from a request not using multipart/form-data encoding type.")
		return
	}
	err = ParseRequestForm(r)
	if err != nil {
		return
	}
	fileHeaders = r.MultipartForm.File[key]
	return
}

// GetFirstFileHeadersByPrefix uses the request to return only the first file header
// from file headers available for each key starting with the given prefix.
// The returned map can be empty if no key starts with the given prefix.
func GetFirstFileHeadersByPrefix(r *http.Request, prefix string) (fileHeaders map[string]*multipart.FileHeader, err error) {
	fileHeaders = make(map[string]*multipart.FileHeader)
	all, err := GetFileHeadersByPrefix(r, prefix)
	if err != nil {
		return
	}
	for k, v := range all {
		fileHeaders[k] = v[0]
	}
	return
}

// GetAllValuesByKey uses the request object to return all values
// associated to a given key.
func GetAllValuesByKey(r *http.Request, key string) (values []string, err error) {
	err = ParseRequestForm(r)
	if err != nil {
		return
	}
	data := getValues(r)
	for _, v := range data[key] {
		values = append(values, strings.TrimSpace(v))
	}
	return
}

// GetFirstValueByKey uses the data object to return only the first value
// from the values associated to a given key.
func GetFirstValueByKey(r *http.Request, key string) (value string, err error) {
	err = ParseRequestForm(r)
	if err != nil {
		return
	}
	values := getValues(r)
	if len(values[key]) <= 0 {
		value = ""
		return
	}
	value = strings.TrimSpace(values[key][0])
	return
}

// GetFirstValuesByPrefix uses the data object to return only the first value
// from values available for each key starting with the given prefix.
// The returned map can be empty if no key in querist starts with the given prefix.
func GetFirstValuesByPrefix(r *http.Request, prefix string) (values map[string]string, err error) {
	values = make(map[string]string)
	all, err := GetValuesByPrefix(r, prefix)
	if err != nil {
		return
	}
	for k, v := range all {
		values[k] = strings.TrimSpace(v[0])
	}
	return
}
