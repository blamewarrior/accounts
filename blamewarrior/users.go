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

package blamewarrior

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type User struct {
	ID        int    `json:"-"`
	Token     string `json:"token"`
	UID       string `json:"uid"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
}

func (u User) Valid() *Validator {
	v := new(Validator)

	v.MustNotBeEmpty(u.Token, "token must not be empty")
	v.MustNotBeEmpty(u.UID, "uid must not be empty")
	v.MustNotBeEmpty(u.Nickname, "nickname must not be empty")
	v.MustNotBeEmpty(u.AvatarURL, "avatar_url must not be empty")
	v.MustNotBeEmpty(u.Name, "name must not be empty")

	return v
}

func SaveUser(db *sql.DB, u *User) (err error) {
	return db.QueryRow(
		SaveUserQuery,
		u.Token,
		u.UID,
		u.Nickname,
		u.AvatarURL,
		u.Name,
	).Scan(&u.ID)
}

const (
	SaveUserQuery = `INSERT INTO users(token, uid, nickname, avatar_url, name) VALUES ($1, $2, $3, $4, $5)
                        ON CONFLICT (nickname)
                        SET token = EXCLUDED.token, name = EXCLUDED.name, avatar_url = EXCLUDED.avatar_url
                        RETURNING id`
)
