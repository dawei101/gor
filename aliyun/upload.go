package aliyun

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/rs/xid"
	"roo.bo/rlib/rhttp"
	"roo.bo/rlib/rlog"
)

type FileUploadHandler struct {
	*Bucket
	MaxSize       int64
	Filekey       string
	UploadPath    string
	AllowFileExts map[string]bool
	DenyFileExts  map[string]bool
}

func get_file_ext(fn string) string {
	exts := strings.Split(fn, ".")
	ext := exts[len(exts)-1]
	return strings.ToLower(ext)
}

// 创建一个图片上传的handler
// 路由中添加以下即可
// r.HandleFunc(path, NewImageUploadHandler().Handle).Methods("POST")
func NewImageUploadHandler(bucket *Bucket, path string) *FileUploadHandler {
	h := NewFileUploadHandler(bucket, path)
	h.AllowFileExts = map[string]bool{
		"jpg":  true,
		"jpeg": true,
		"png":  true,
		"gif":  true,
	}
	return h
}

// 创建一个文件上传的handler
// 路由中添加以下即可
// h := NewFileUploadHandler()
// h.AllowFileExts = map[string]bool{} #bool值无实际作用
// h.MaxSize = 1024*1024*4 # 4M
// r.HandleFunc(path, h.Handle).Methods("POST")
func NewFileUploadHandler(bucket *Bucket, path string) *FileUploadHandler {
	return &FileUploadHandler{
		MaxSize:       32 << 20,
		Filekey:       "file",
		Bucket:        bucket,
		UploadPath:    path,
		AllowFileExts: map[string]bool{},
		DenyFileExts: map[string]bool{
			"jsp":  true,
			"php":  true,
			"sh":   true,
			"bash": true,
		},
	}
}

func (h *FileUploadHandler) AllowFileTypes(exts ...string) {
	for _, ext := range exts {
		ext = strings.ToLower(ext)
		delete(h.DenyFileExts, ext)
		h.AllowFileExts[ext] = true
	}
}

func (h *FileUploadHandler) DenyFileTypes(exts ...string) {
	for _, ext := range exts {
		ext = strings.ToLower(ext)
		delete(h.AllowFileExts, ext)
		h.DenyFileExts[ext] = true
	}
}

func (h *FileUploadHandler) GenerateObjectKey(ext string) string {
	t := time.Now()
	date := fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
	uniq := xid.New().String()
	return path.Join(h.UploadPath, date, uniq+"."+ext)
}

func (h *FileUploadHandler) checkFileExtension(filename string) error {

	ext := get_file_ext(filename)

	if _, ok := h.DenyFileExts[ext]; ok {
		msg := fmt.Sprintf("upload file type:%s is in deny list!", ext)
		return errors.New(msg)
	}

	if len(h.AllowFileExts) > 0 {
		if _, allow := h.AllowFileExts[ext]; !allow {
			msg := fmt.Sprintf("upload file type:%s is not in allow list!", ext)
			return errors.New(msg)
		}
	}
	return nil
}

func (h *FileUploadHandler) Handle(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(h.MaxSize)
	f, header, err := r.FormFile(h.Filekey)
	if err != nil {
		errmsg := fmt.Sprintf("can't find key:%s in form request", h.Filekey)
		rlog.Error(r.Context(), errmsg)
		rhttp.NewErrResp(422, errmsg, err.Error()).Flush(w)
		return
	}
	defer f.Close()

	err = h.checkFileExtension(header.Filename)
	if err != nil {
		rlog.Info(r.Context(), err.Error())
		rhttp.NewErrResp(403, err.Error(), "").Flush(w)
		return
	}

	ext := get_file_ext(header.Filename)

	objkey := h.GenerateObjectKey(ext)
	rlog.Debug(r.Context(), "oss:generate object key:", objkey)
	err = h.PutObject(objkey, f)
	if err != nil {
		errmsg := fmt.Sprintf("put file:%s to aliyun bucket:%s failed", objkey, h.BucketName)
		rlog.Error(r.Context(), errmsg)
		rhttp.NewErrResp(500, errmsg, err.Error()).Flush(w)
		return
	}

	url := h.FileUrl(objkey)
	rhttp.NewResp(map[string]string{"url": url}).Flush(w)
}
