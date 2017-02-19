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

package main_test

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/blamewarrior/users"
	"github.com/blamewarrior/users/blamewarrior"
)

func TestGetUserByNickname(t *testing.T) {
	db, teardown := setup()

	_, err := db.Exec("TRUNCATE users;")

	require.NoError(t, err)

	handlers := main.NewUserHandlers(db)
	defer teardown()

	user := &blamewarrior.User{
		Token:     "test_token",
		UID:       "133445",
		Nickname:  "blamewarrior_test",
		AvatarURL: "http://example.com/12345.jpg",
		Name:      "Blamewarrior Test",
	}

	err = blamewarrior.SaveUser(db, user)

	require.NoError(t, err)

	req, err := http.NewRequest("GET", "/users?:nickname=blamewarrior_test", nil)

	require.NoError(t, err)

	w := httptest.NewRecorder()

	handlers.GetUserByNickname(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(
		t,
		"{\"token\":\"test_token\",\"uid\":\"133445\",\"nickname\":\"blamewarrior_test\",\"avatar_url\":\"http://example.com/12345.jpg\",\"name\":\"Blamewarrior Test\"}\n",
		fmt.Sprintf("%v", w.Body),
	)
}

func TestSaveUser(t *testing.T) {

	db, teardown := setup()

	_, err := db.Exec("TRUNCATE users;")

	require.NoError(t, err)

	defer teardown()

	handlers := main.NewUserHandlers(db)

	results := []struct {
		RequestBody  string
		ResponseCode int
		ResponseBody string
	}{
		{
			RequestBody:  `{}`,
			ResponseCode: http.StatusUnprocessableEntity,
			ResponseBody: "[\"token must not be empty\",\"uid must not be empty\",\"nickname must not be empty\",\"avatar_url must not be empty\"]\n",
		},
		{
			RequestBody:  testLoginResponse,
			ResponseCode: http.StatusCreated,
			ResponseBody: "",
		},
	}

	for _, result := range results {
		req, err := http.NewRequest("POST", "/users", strings.NewReader(result.RequestBody))

		require.NoError(t, err)

		w := httptest.NewRecorder()

		handlers.SaveUser(w, req)

		assert.Equal(t, result.ResponseCode, w.Code)
		assert.Equal(t, result.ResponseBody, fmt.Sprintf("%v", w.Body))
	}

}

func setup() (db *sql.DB, teardownFn func()) {
	dbName := os.Getenv("DB_NAME")

	if dbName == "" {
		log.Fatal("missing test database name (expected to be passed via ENV['DB_NAME'])")
	}

	opts := &blamewarrior.DatabaseOptions{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	db, err := blamewarrior.ConnectDatabase(dbName, opts)
	if err != nil {
		log.Fatalf("failed to establish connection with test db %s using connection string %s: %s", dbName, opts.ConnectionString(), err)
	}

	if err != nil {
		log.Fatal("failed to create transaction, %s", err)
	}

	return db, func() {

		if err := db.Close(); err != nil {
			log.Printf("failed to close database connection: %s", err)
		}
	}
}

const testLoginResponse = `
  {
   "provider":"github",
   "uid":"583231",
   "info":{
      "nickname":"octocat",
      "email":"octocat@github.com",
      "name":"The Octocat",
      "image":"https://avatars.githubusercontent.com/u/583231?v=3",
      "urls":{
         "GitHub":"https://github.com/octocat",
         "Blog":null
      }
   },
   "credentials":{
      "token":"test_token",
      "expires":false
   },
   "extra":{
      "raw_info":{
         "login":"octocat",
         "id":583231,
         "avatar_url":"https://avatars.githubusercontent.com/u/583231?v=3",
         "gravatar_id":"",
         "url":"https://api.github.com/users/octocat",
         "html_url":"https://github.com/octocat",
         "followers_url":"https://api.github.com/users/octocat/followers",
         "following_url":"https://api.github.com/users/octocat/following{/other_user}",
         "gists_url":"https://api.github.com/users/octocat/gists{/gist_id}",
         "starred_url":"https://api.github.com/users/octocat/starred{/owner}{/repo}",
         "subscriptions_url":"https://api.github.com/users/octocat/subscriptions",
         "organizations_url":"https://api.github.com/users/octocat/orgs",
         "repos_url":"https://api.github.com/users/octocat/repos",
         "events_url":"https://api.github.com/users/octocat/events{/privacy}",
         "received_events_url":"https://api.github.com/users/octocat/received_events",
         "type":"User",
         "site_admin":false,
         "name":"The Octocat",
         "company":null,
         "blog":null,
         "location":null,
         "email":"octocat@github.com",
         "hireable":null,
         "bio":null,
         "public_repos":16,
         "public_gists":5,
         "followers":3,
         "following":0,
         "created_at":"2011-09-21T20:26:50Z",
         "updated_at":"2017-01-30T13:19:47Z"
      }
   }
}
`
