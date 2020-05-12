package acl

import (
	"context"
	"encoding/json"
	"github.com/uhppoted/uhppote-core/uhppote"
	api "github.com/uhppoted/uhppoted-api/acl"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-rest/errors"
	"io/ioutil"
	"net/http"
)

func PutACL(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("put-acl", "Error reading request"),
			err
	}

	body := struct {
		ACL []permission `json:"acl"`
	}{}

	if err = json.Unmarshal(blob, &body); err != nil {
		return http.StatusBadRequest,
			errors.NewRESTError("put-acl", "Invalid request format"),
			err
	}

	table, err := PermissionsToTable(body.ACL)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("put-acl", "Error parsing request"),
			err
	}

	u := ctx.Value("uhppote").(*uhppote.UHPPOTE)
	devices := ctx.Value("devices").([]*uhppote.Device)

	acl, err := api.ParseTable(*table, devices)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("put-acl", "Error processing access control list"),
			err
	}

	rpt, err := api.PutACL(u, *acl)
	if err != nil {
		return http.StatusInternalServerError,
			errors.NewRESTError("put-table", "Error storing access control list"),
			err
	}

	report := []struct {
		DeviceID  uint32 `json:"device-id"`
		Unchanged int    `json:"unchanged"`
		Updated   int    `json:"updated"`
		Added     int    `json:"added"`
		Deleted   int    `json:"deleted"`
		Failed    int    `json:"failed"`
	}{}

	for k, v := range rpt {
		report = append(report, struct {
			DeviceID  uint32 `json:"device-id"`
			Unchanged int    `json:"unchanged"`
			Updated   int    `json:"updated"`
			Added     int    `json:"added"`
			Deleted   int    `json:"deleted"`
			Failed    int    `json:"failed"`
		}{
			DeviceID:  k,
			Unchanged: v.Unchanged,
			Updated:   v.Updated,
			Added:     v.Added,
			Deleted:   v.Deleted,
			Failed:    v.Failed,
		})
	}

	return http.StatusOK, &struct {
		Report []struct {
			DeviceID  uint32 `json:"device-id"`
			Unchanged int    `json:"unchanged"`
			Updated   int    `json:"updated"`
			Added     int    `json:"added"`
			Deleted   int    `json:"deleted"`
			Failed    int    `json:"failed"`
		} `json:"report"`
	}{
		Report: report,
	}, nil
}
