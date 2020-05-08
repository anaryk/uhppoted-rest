package acl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	api "github.com/uhppoted/uhppoted-api/acl"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-rest/errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func Grant(impl *uhppoted.UHPPOTED, ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, *errors.IError) {
	url := r.URL.Path

	matches := regexp.MustCompile("^/uhppote/acl/card/([0-9]+)/doors/(\\S.*)$").FindStringSubmatch(url)
	if matches == nil || len(matches) < 3 {
		return nil, errors.Errorf(fmt.Errorf("%w: Missing card number/door", uhppoted.BadRequest), 0, "grant", "Missing card number/door")
	}

	cardID, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		return nil, errors.Errorf(fmt.Errorf("%w: Invalid card number (%s)", uhppoted.BadRequest, matches[1]), 0, "grant", "Invalid card number")
	}

	doors := []string{}
	tokens := strings.Split(matches[2], ",")
	for _, s := range tokens {
		if d := strings.TrimSpace(s); d != "" {
			doors = append(doors, d)
		}
	}

	if len(doors) == 0 {
		return nil, errors.Errorf(fmt.Errorf("%w: Invalid list of doors (%s)", uhppoted.BadRequest, matches[2]), 0, "grant", "Invalid list of doors")
	}

	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Errorf(fmt.Errorf("%w: Error reading request", err), 0, "grant", "Error reading request")
	}

	body := struct {
		From *types.Date `json:"start-date"`
		To   *types.Date `json:"end-date"`
	}{}

	if err = json.Unmarshal(blob, &body); err != nil {
		return nil, errors.Errorf(fmt.Errorf("%w: Invalid request format", err), 0, "grant", "Invalid request format")
	}

	if body.From == nil {
		return nil, errors.Errorf(fmt.Errorf("%w: Missing 'start-date'", uhppoted.BadRequest), 0, "grant", "Missing 'start-date'")
	}

	if body.To == nil {
		return nil, errors.Errorf(fmt.Errorf("%w: Missing 'end-date'", uhppoted.BadRequest), 0, "grant", "Missing 'end-date'")
	}

	u := ctx.Value("uhppote").(*uhppote.UHPPOTE)
	devices := ctx.Value("devices").([]*uhppote.Device)

	err = api.Grant(u, devices, uint32(cardID), *body.From, *body.To, doors)
	if err != nil {
		return nil, errors.Errorf(fmt.Errorf("%w: Error granting card access permissions", err), 0, "grant", "Error granting card access permissions")
	}

	return nil, nil
}