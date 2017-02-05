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
	"errors"
)

var UserNotFound = errors.New("User not found")

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

func SaveUser(q Queryer, u *User) (err error) {
	return q.QueryRow(
		SaveUserQuery,
		u.Token,
		u.UID,
		u.Nickname,
		u.AvatarURL,
		u.Name,
	).Scan(&u.ID)
}

func GetUserByNickname(q Queryer, nickname string) (u *User, err error) {
	err = q.QueryRow(GetUserByNicknameQuery, nickname).Scan(
		&u.Token,
		&u.UID,
		&u.Nickname,
		&u.AvatarURL,
		&u.Name,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, UserNotFound
		}
		return nil, err
	}

	return u, nil
}

const (
	GetUserByNicknameQuery = "SELECT token, uid, nickname, avatar_url, name FROM users WHERE nickname = $1"
	SaveUserQuery          = `INSERT INTO users(token, uid, nickname, avatar_url, name) VALUES ($1, $2, $3, $4, $5)
                        ON CONFLICT (nickname) DO UPDATE
                        SET token = EXCLUDED.token, name = EXCLUDED.name, avatar_url = EXCLUDED.avatar_url
                        RETURNING id`
)
