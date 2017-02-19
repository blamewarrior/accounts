/*
   Copyright (C) 2017 The BlameWarrior Authors.
   This file is a part of BlameWarrior service.
   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/blamewarrior/users/blamewarrior"
)

type LoginResult struct {
	UID string `json:"uid"`

	Info struct {
		Nickname string `json:"nickname"`
		Image    string `json:"image"`
		Name     string `json:"name"`
	} `json:"info"`

	Credentials struct {
		Token string `json:token`
	} `json:"credentials"`
}

type UserHandlers struct {
	db *sql.DB
}

func NewUserHandlers(db *sql.DB) *UserHandlers {
	return &UserHandlers{db}
}

func (handler *UserHandlers) GetUserByNickname(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	nickname := req.URL.Query().Get(":nickname")

	if nickname == "" {
		http.Error(w, "Specify correct username", http.StatusNotFound)
		return
	}

	user, err := blamewarrior.GetUserByNickname(handler.db, nickname)

	if err != nil {
		if err == blamewarrior.UserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s\t%s\t%v\t%s", "GET", req.RequestURI, http.StatusInternalServerError, err)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s\t%s\t%v\t%s", "GET", req.RequestURI, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *UserHandlers) SaveUser(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	body, err := requestBody(req)

	if err != nil {
		http.Error(w, "Unable to parse incoming json", http.StatusBadRequest)
		return
	}

	loginResult := &LoginResult{}

	if err = json.Unmarshal(body, &loginResult); err != nil {
		http.Error(w, "Unable to parse incoming json", http.StatusBadRequest)
		return
	}

	user := &blamewarrior.User{
		Token:     loginResult.Credentials.Token,
		UID:       loginResult.UID,
		Nickname:  loginResult.Info.Nickname,
		AvatarURL: loginResult.Info.Image,
		Name:      loginResult.Info.Name,
	}

	validator := user.Valid()

	if valid := validator.IsValid(); !valid {
		messages, err := json.Marshal(validator.ErrorMessages())

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("%s\t%s\t%v\t%s", "POST", req.RequestURI, http.StatusInternalServerError, err)
			return
		}

		http.Error(w, string(messages), http.StatusUnprocessableEntity)
		return
	}

	if err = blamewarrior.SaveUser(handler.db, user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s\t%s\t%v\t%s", "POST", req.RequestURI, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
}

func requestBody(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		return nil, err
	}
	if err := r.Body.Close(); err != nil {
		return nil, err
	}
	return body, err
}
