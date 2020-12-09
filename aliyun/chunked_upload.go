package aliyun

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gorilla/mux"
	"github.com/rs/xid"

	"roo.bo/rlib/rhttp"
	"roo.bo/rlib/rlog"
	"roo.bo/rlib/rredis"
)

type BigFileUpload struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	MD5      string `json:"md5"`
}

func (u BigFileUpload) totalParts(chunkSize int64) int {
	return int(1 + (u.Size-1)/chunkSize)
}

type ChunkedInfo struct {
	UUID      string                            `json:"uuid"`
	Filename  string                            `json:"filename"`
	Size      int64                             `json:"size"`
	ChunkSize int64                             `json:"chunkSize"`
	Parts     int                               `json:"parts"`
	ObjectKey string                            `json:"objectKey"`
	Imur      oss.InitiateMultipartUploadResult `json:"imur"`
}

type ChunkedUploadHandler struct {
	*FileUploadHandler
	ChunkSize int64 `json:"chuckSize"`
}

func NewChunkedFileHandler(bucket *Bucket, path string) *ChunkedUploadHandler {
	h := &ChunkedUploadHandler{NewFileUploadHandler(bucket, path), 1024 * 1024 * 2}
	h.MaxSize = 1024 * 1024 * 1024
	return h
}

func (h *ChunkedUploadHandler) ChuckFile(u *BigFileUpload) (*ChunkedInfo, error) {
	exts := strings.Split(u.Filename, ".")
	ext := exts[len(exts)-1]
	objkey := h.GenerateObjectKey(ext)
	imur, err := h.InitiateMultipartUpload(objkey)
	if err != nil {
		return nil, err
	}

	return &ChunkedInfo{
		Filename:  u.Filename,
		Size:      u.Size,
		ChunkSize: h.ChunkSize,
		Parts:     u.totalParts(h.ChunkSize),
		UUID:      xid.New().String(),
		ObjectKey: objkey,
		Imur:      imur,
	}, nil
}

func uploading_file_key(uuid string) string {
	return fmt.Sprintf("uploading:%s", uuid)
}

func uploading_file_part_key(uuid string) string {
	return fmt.Sprintf("uploading:%s:part", uuid)
}

func (h *ChunkedUploadHandler) Route(r *mux.Router) {
	r.HandleFunc("/chunked-upload", h.CreateChunkedUploadHandle).Methods("POST")
	r.HandleFunc("/chunked-upload/{uuid}/{chunk:[0-9]+}", h.UploadOneChunk).Methods("POST")
}

func (h *ChunkedUploadHandler) CreateChunkedUploadHandle(w http.ResponseWriter, r *http.Request) {
	s, err := rhttp.JsonBody(r)
	if err != nil {
		rhttp.NewErrResp(422, "body is not json", "").Flush(w)
		return
	}
	upload := &BigFileUpload{}
	s.DataAssignTo(upload)
	err = h.checkFileExtension(upload.Filename)
	if err != nil {
		rlog.Info(r.Context(), err.Error())
		rhttp.NewErrResp(403, err.Error(), "").Flush(w)
		return
	}
	if upload.Size > h.MaxSize || upload.Size <= 0 {
		rhttp.NewErrResp(422, fmt.Sprintf("file size is not ok, max size is: %d", h.MaxSize), "").Flush(w)
		return
	}
	if upload.totalParts(h.ChunkSize) <= 1 {
		rhttp.NewErrResp(422, fmt.Sprintf("file size is less than one chunk size:%d", h.ChunkSize), "").Flush(w)
		return
	}

	info, err := h.ChuckFile(upload)
	if err != nil {
		rhttp.NewErrResp(422, fmt.Sprintf("file size is less than one chunk size:%d", h.ChunkSize), "")
		return
	}
	info_s, _ := json.Marshal(info)
	err = rredis.DefaultRedis().Set(r.Context(), uploading_file_key(info.UUID), info_s, 30*time.Minute).Err()
	if err != nil {
		rlog.Error(r.Context(), "save multipart uploading info failed:", err)
		rhttp.NewErrResp(500, "prepare upload failed, please retry!", "save multipart info to redis failed").Flush(w)
		return
	}
	uuid := info.UUID
	time.AfterFunc(30*time.Minute, func() {
		h.clearUpload(uuid)
	})
	rhttp.NewResp(info).Flush(w)
}

func (h *ChunkedUploadHandler) UploadOneChunk(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]
	chunk_s := mux.Vars(r)["chunk"]
	chunk, _ := strconv.Atoi(chunk_s)

	info_s, err := rredis.DefaultRedis().Get(r.Context(), uploading_file_key(uuid)).Result()
	if err == rredis.Nil {
		rhttp.NewErrResp(422, "no upload action found!", "").Flush(w)
		return
	} else if err != nil {
		rlog.Error(r.Context(), "get multipart uploading info failed:", err)
		rhttp.NewErrResp(500, "upload failed, please retry!", "get multipart info to redis failed").Flush(w)
		return
	}
	var info ChunkedInfo
	json.Unmarshal([]byte(info_s), &info)
	if chunk > info.Parts {
		rhttp.NewErrResp(422, "chunked parts not found!", "").Flush(w)
		return
	}
	r.ParseMultipartForm(h.ChunkSize)
	f, header, err := r.FormFile(h.Filekey)
	if err != nil {
		errmsg := fmt.Sprintf("can't find key:%s in form request", h.Filekey)
		rlog.Error(r.Context(), errmsg)
		rhttp.NewErrResp(422, errmsg, err.Error()).Flush(w)
		return
	}

	respect_size := info.ChunkSize
	if chunk == info.Parts {
		respect_size = info.Size % info.ChunkSize
	}
	if header.Size != respect_size {
		errmsg := fmt.Sprintf("chunk=%d size=%d is not correct, respect=%d", chunk, header.Size, respect_size)
		rhttp.NewErrResp(422, errmsg, "").Flush(w)
		return
	}
	part, err := h.UploadPart(info.Imur, f, respect_size, chunk)
	if err != nil {
		rlog.Error(r.Context(), "upload chunk=", chunk, "failed with err:", err)
		rhttp.NewErrResp(500, err.Error(), "").Flush(w)
		return
	}
	completed_c, err := h.completePart(r.Context(), info, chunk, part)
	if err != nil {
		rlog.Error(r.Context(), "upload chunk=", chunk, "failed with err:", err)
		rhttp.NewErrResp(500, err.Error(), "").Flush(w)
		return
	}
	rhttp.NewResp(map[string]interface{}{
		"chunk":     chunk,
		"completed": completed_c,
		"total":     info.Parts,
		"done":      completed_c == info.Parts,
		"url":       h.FileUrl(info.ObjectKey),
	}).Flush(w)
}

func (h *ChunkedUploadHandler) completePart(ctx context.Context, info ChunkedInfo, chunk int, part oss.UploadPart) (completed int, err error) {
	part_s, _ := json.Marshal(part)
	err = rredis.DefaultRedis().HSet(ctx, uploading_file_part_key(info.UUID), string(chunk), part_s).Err()
	if err != nil {
		return 0, err
	}
	result, err := rredis.DefaultRedis().HGetAll(ctx, uploading_file_part_key(info.UUID)).Result()
	if err != nil {
		return 0, err
	}
	parts := []oss.UploadPart{}
	for _, part_s := range result {
		one := oss.UploadPart{}
		json.Unmarshal([]byte(part_s), &one)
		parts = append(parts, one)
	}
	if len(parts) == info.Parts {
		_, err := h.CompleteMultipartUpload(info.Imur, parts)
		if err != nil {
			return len(parts), err
		}
		go h.clearUpload(info.UUID)
	}
	return len(parts), nil
}

func (h *ChunkedUploadHandler) clearUpload(uuid string) {
	rredis.DefaultRedis().Del(context.Background(), uploading_file_part_key(uuid))
	rredis.DefaultRedis().Del(context.Background(), uploading_file_key(uuid))
}
