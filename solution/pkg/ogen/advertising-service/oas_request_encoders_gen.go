// Code generated by ogen, DO NOT EDIT.

package api

import (
	"bytes"
	"net/http"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"

	ht "github.com/ogen-go/ogen/http"
)

func encodeAdvanceDayRequest(
	req OptAdvanceDayReq,
	r *http.Request,
) error {
	const contentType = "application/json"
	if !req.Set {
		// Keep request with empty body if value is not set.
		return nil
	}
	e := new(jx.Encoder)
	{
		if req.Set {
			req.Encode(e)
		}
	}
	encoded := e.Bytes()
	ht.SetBody(r, bytes.NewReader(encoded), contentType)
	return nil
}

func encodeCreateCampaignRequest(
	req *CampaignCreate,
	r *http.Request,
) error {
	const contentType = "application/json"
	e := new(jx.Encoder)
	{
		req.Encode(e)
	}
	encoded := e.Bytes()
	ht.SetBody(r, bytes.NewReader(encoded), contentType)
	return nil
}

func encodeGenerateAdTextRequest(
	req *GenerateAdTextReq,
	r *http.Request,
) error {
	const contentType = "application/json"
	e := new(jx.Encoder)
	{
		req.Encode(e)
	}
	encoded := e.Bytes()
	ht.SetBody(r, bytes.NewReader(encoded), contentType)
	return nil
}

func encodeModerateAdTextRequest(
	req *ModerateAdTextReq,
	r *http.Request,
) error {
	const contentType = "application/json"
	e := new(jx.Encoder)
	{
		req.Encode(e)
	}
	encoded := e.Bytes()
	ht.SetBody(r, bytes.NewReader(encoded), contentType)
	return nil
}

func encodeRecordAdClickRequest(
	req *RecordAdClickReq,
	r *http.Request,
) error {
	const contentType = "application/json"
	e := new(jx.Encoder)
	{
		req.Encode(e)
	}
	encoded := e.Bytes()
	ht.SetBody(r, bytes.NewReader(encoded), contentType)
	return nil
}

func encodeUpdateCampaignRequest(
	req *CampaignUpdate,
	r *http.Request,
) error {
	const contentType = "application/json"
	e := new(jx.Encoder)
	{
		req.Encode(e)
	}
	encoded := e.Bytes()
	ht.SetBody(r, bytes.NewReader(encoded), contentType)
	return nil
}

func encodeUploadCampaignImageRequest(
	req UploadCampaignImageReq,
	r *http.Request,
) error {
	switch req := req.(type) {
	case *UploadCampaignImageReqImageJpeg:
		const contentType = "image/jpeg"
		body := req
		ht.SetBody(r, body, contentType)
		return nil
	case *UploadCampaignImageReqImagePNG:
		const contentType = "image/png"
		body := req
		ht.SetBody(r, body, contentType)
		return nil
	default:
		return errors.Errorf("unexpected request type: %T", req)
	}
}

func encodeUpsertAdvertisersRequest(
	req []AdvertiserUpsert,
	r *http.Request,
) error {
	const contentType = "application/json"
	e := new(jx.Encoder)
	{
		e.ArrStart()
		for _, elem := range req {
			elem.Encode(e)
		}
		e.ArrEnd()
	}
	encoded := e.Bytes()
	ht.SetBody(r, bytes.NewReader(encoded), contentType)
	return nil
}

func encodeUpsertClientsRequest(
	req []ClientUpsert,
	r *http.Request,
) error {
	const contentType = "application/json"
	e := new(jx.Encoder)
	{
		e.ArrStart()
		for _, elem := range req {
			elem.Encode(e)
		}
		e.ArrEnd()
	}
	encoded := e.Bytes()
	ht.SetBody(r, bytes.NewReader(encoded), contentType)
	return nil
}

func encodeUpsertMLScoreRequest(
	req *MLScore,
	r *http.Request,
) error {
	const contentType = "application/json"
	e := new(jx.Encoder)
	{
		req.Encode(e)
	}
	encoded := e.Bytes()
	ht.SetBody(r, bytes.NewReader(encoded), contentType)
	return nil
}
