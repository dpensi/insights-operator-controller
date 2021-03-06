/*
Copyright © 2019, 2020 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

import (
	"github.com/redhatinsighs/insights-operator-controller/utils"
	"io/ioutil"
	"net/http"
	"strconv"
)

// ListConfigurationProfiles - read list of configuration profiles.
func (s Server) ListConfigurationProfiles(writer http.ResponseWriter, request *http.Request) {
	profiles, err := s.Storage.ListConfigurationProfiles()
	if err == nil {
		utils.SendResponse(writer, utils.BuildOkResponseWithData("profiles", profiles))
	} else {
		utils.SendError(writer, err.Error())
	}
}

// GetConfigurationProfile - read profile specified by its ID
func (s Server) GetConfigurationProfile(writer http.ResponseWriter, request *http.Request) {
	id, err := retrieveIDRequestParameter(request)
	if err != nil {
		utils.SendError(writer, "Error reading profile ID from request\n")
		return
	}

	profile, err := s.Storage.GetConfigurationProfile(int(id))
	if err == nil {
		utils.SendResponse(writer, utils.BuildOkResponseWithData("profile", profile))
	} else {
		utils.SendError(writer, err.Error())
	}
}

// NewConfigurationProfile - create new configuration profile
func (s Server) NewConfigurationProfile(writer http.ResponseWriter, request *http.Request) {
	username, foundUsername := request.URL.Query()["username"]
	description, foundDescription := request.URL.Query()["description"]

	if !foundUsername {
		utils.SendError(writer, "User name needs to be specified\n")
		return
	}

	if !foundDescription {
		utils.SendError(writer, "Description needs to be specified\n")
		return
	}

	configuration, err := ioutil.ReadAll(request.Body)
	if err != nil || len(configuration) == 0 {
		utils.SendError(writer, "Configuration needs to be provided in the request body")
		return
	}

	s.Splunk.LogAction("NewConfigurationProfile", username[0], string(configuration))
	profiles, err := s.Storage.StoreConfigurationProfile(username[0], description[0], string(configuration))
	if err != nil {
		utils.SendInternalServerError(writer, err.Error())
	} else {
		utils.SendCreated(writer, utils.BuildOkResponseWithData("profiles", profiles))
	}
}

// DeleteConfigurationProfile - delete configuration profile
func (s Server) DeleteConfigurationProfile(writer http.ResponseWriter, request *http.Request) {
	id, err := retrieveIDRequestParameter(request)
	if err != nil {
		utils.SendError(writer, "Error reading profile ID from request\n")
		return
	}

	s.Splunk.LogAction("DeleteConfigurationProfile", "tester", strconv.Itoa(int(id)))
	profiles, err := s.Storage.DeleteConfigurationProfile(int(id))
	if err != nil {
		utils.SendError(writer, err.Error())
	} else {
		utils.SendResponse(writer, utils.BuildOkResponseWithData("profiles", profiles))
	}
}

// ChangeConfigurationProfile - change configuration profile
func (s Server) ChangeConfigurationProfile(writer http.ResponseWriter, request *http.Request) {
	id, err := retrieveIDRequestParameter(request)
	username, foundUsername := request.URL.Query()["username"]
	description, foundDescription := request.URL.Query()["description"]

	if err != nil {
		utils.SendError(writer, "Error reading profile ID from request\n")
		return
	}

	if !foundUsername {
		utils.SendError(writer, "User name needs to be specified\n")
		return
	}

	if !foundDescription {
		utils.SendError(writer, "Description needs to be specified\n")
		return
	}

	configuration, err := ioutil.ReadAll(request.Body)
	if err != nil || len(configuration) == 0 {
		utils.SendError(writer, "Configuration needs to be provided in the request body")
		return
	}

	s.Splunk.LogAction("ChangeConfigurationProfile", username[0], string(configuration))
	profiles, err := s.Storage.ChangeConfigurationProfile(int(id), username[0], description[0], string(configuration))
	if err != nil {
		utils.SendError(writer, err.Error())
	} else {
		utils.SendAccepted(writer, utils.BuildOkResponseWithData("profiles", profiles))
	}
}
