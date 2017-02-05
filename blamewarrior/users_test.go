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

package blamewarrior_test

import (
	"database/sql"
	"log"
	"os"

	"github.com/blamewarrior/users/blamewarrior"

	"github.com/stretchr/testify/require"
	"testing"
)

func TestSaveUser(t *testing.T) {
	tx, teardown := setup()
	defer teardown()

	_, err := tx.Exec("TRUNCATE users;")

	require.NoError(t, err)

	user := &blamewarrior.User{
		Token:     "test token",
		UID:       "123",
		Nickname:  "blamewarrior_test",
		AvatarURL: "https://avatars1.githubusercontent.com/u/788766655341678980?v=3&s=40",
		Name:      "Blamewarrior Test",
	}

	err = blamewarrior.SaveUser(tx, user)

	require.NoError(t, err)
}

func TestGetUserByNickname_userExists(t *testing.T) {
	tx, teardown := setup()
	defer teardown()

	_, err := tx.Exec("TRUNCATE users;")

	require.NoError(t, err)

	_, err = tx.Exec(
		blamewarrior.SaveUserQuery,
		"test_token",
		"uid123",
		"test_user",
		"https://avatars1.githubusercontent.com/u/788766655341678980?v=3&s=40",
		"Blamewarrior Test",
	)

	require.NoError(t, err)

	_, err = blamewarrior.GetUserByNickname(tx, "test_user")

	require.NoError(t, err)

}

func TestGetUserByNickname_userDoNotExist(t *testing.T) {
	tx, teardown := setup()
	defer teardown()

	_, err := tx.Exec("TRUNCATE users;")

	require.NoError(t, err)

	_, err = blamewarrior.GetUserByNickname(tx, "test_user")

	require.Error(t, blamewarrior.UserNotFound)

}

func setup() (tx *sql.Tx, teardownFn func()) {
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
	tx, err = db.Begin()

	if err != nil {
		log.Fatal("failed to create transaction, %s", err)
	}

	return tx, func() {
		tx.Rollback()
		if err := db.Close(); err != nil {
			log.Printf("failed to close database connection: %s", err)
		}
	}
}
