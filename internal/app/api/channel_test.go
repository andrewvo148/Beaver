// Copyright 2018 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"github.com/clivern/beaver/internal/pkg/utils"
	"github.com/nbio/st"
	"os"
	"strings"
	"testing"
	"time"
)

// init setup stuff
func init() {
	basePath := fmt.Sprintf("%s/src/github.com/clivern/beaver", os.Getenv("GOPATH"))
	configFile := fmt.Sprintf("%s/%s", basePath, "config.test.json")

	config := utils.Config{}
	ok, err := config.Load(configFile)

	if !ok || err != nil {
		panic(err.Error())
	}
	config.Cache()
	config.GinEnv()
	if !strings.Contains(os.Getenv("LogPath"), basePath) {
		os.Setenv("LogPath", fmt.Sprintf("%s/%s", basePath, os.Getenv("LogPath")))
	}
}

// TestChannelAPI test cases
func TestChannelAPI(t *testing.T) {

	createdAt := time.Now().Unix()
	updatedAt := time.Now().Unix()

	channelResult := ChannelResult{Name: "name", Type: "type", CreatedAt: createdAt, UpdatedAt: updatedAt}
	jsonValue, err := channelResult.ConvertToJSON()
	st.Expect(t, jsonValue, fmt.Sprintf(`{"name":"name","type":"type","created_at":%d,"updated_at":%d}`, createdAt, updatedAt))
	st.Expect(t, err, nil)

	ok, err := channelResult.LoadFromJSON([]byte(jsonValue))
	st.Expect(t, ok, true)
	st.Expect(t, err, nil)
	st.Expect(t, channelResult.Name, "name")
	st.Expect(t, channelResult.Type, "type")
	st.Expect(t, channelResult.CreatedAt, createdAt)
	st.Expect(t, channelResult.UpdatedAt, updatedAt)

	channelAPI := Channel{}
	st.Expect(t, channelAPI.Init(), true)

	//Clear
	channelAPI.DeleteChannelByName(channelResult.Name)

	ok, err = channelAPI.CreateChannel(channelResult)
	st.Expect(t, ok, true)
	st.Expect(t, err, nil)

	newChannelResult, err := channelAPI.GetChannelByName(channelResult.Name)
	st.Expect(t, channelResult.Name, newChannelResult.Name)
	st.Expect(t, channelResult.Type, newChannelResult.Type)
	st.Expect(t, channelResult.CreatedAt, newChannelResult.CreatedAt)
	st.Expect(t, channelResult.UpdatedAt, newChannelResult.UpdatedAt)
	st.Expect(t, err, nil)

	newChannelResult.Type = "new_type"

	ok, err = channelAPI.UpdateChannelByName(newChannelResult)
	st.Expect(t, ok, true)
	st.Expect(t, err, nil)

	newChannelResult, err = channelAPI.GetChannelByName(channelResult.Name)
	st.Expect(t, channelResult.Name, newChannelResult.Name)
	st.Expect(t, "new_type", newChannelResult.Type)
	st.Expect(t, channelResult.CreatedAt, newChannelResult.CreatedAt)
	st.Expect(t, channelResult.UpdatedAt, newChannelResult.UpdatedAt)
	st.Expect(t, err, nil)

	st.Expect(t, 0, int(channelAPI.CountSubscribers(channelResult.Name)))
	st.Expect(t, 0, int(channelAPI.CountListeners(channelResult.Name)))

	ok, err = channelAPI.DeleteChannelByName(newChannelResult.Name)
	st.Expect(t, ok, true)
	st.Expect(t, err, nil)
}
