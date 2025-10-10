/*
 * // Copyright 2025 getDragon-dev
 * // Licensed under the Apache License, Version 2.0 (the "License");
 * // you may not use this file except in compliance with the License.
 * // You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 * // Unless required by applicable law or agreed to in writing, software
 * // distributed under the License is distributed on an "AS IS" BASIS,
 * // WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * // See the License for the specific language governing permissions and
 * // limitations under the License.
 */

package registry

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Blueprint struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Repo        string   `json:"repo"`
	Path        string   `json:"path"`
	DownloadURL string   `json:"download_url"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type Database struct {
	Blueprints []Blueprint `json:"blueprints"`
}

// Load reads registry.json from a local path or an HTTP(S) URL.
func Load(location string) (Database, error) {
	var db Database
	var data []byte
	var err error

	if strings.HasPrefix(location, "http://") || strings.HasPrefix(location, "https://") {
		data, err = fetch(location)
	} else {
		data, err = os.ReadFile(location)
	}
	if err != nil {
		return db, err
	}
	if err := json.Unmarshal(data, &db); err != nil {
		return db, err
	}
	if db.Blueprints == nil {
		db.Blueprints = []Blueprint{}
	}
	for i := range db.Blueprints {
		if db.Blueprints[i].Tags == nil {
			db.Blueprints[i].Tags = []string{}
		}
	}
	return db, nil
}

func fetch(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GET %s: %d: %s", url, resp.StatusCode, string(b))
	}
	return io.ReadAll(resp.Body)
}

func Save(path string, db Database) error {
	if db.Blueprints == nil {
		db.Blueprints = []Blueprint{}
	}
	for i := range db.Blueprints {
		if db.Blueprints[i].Tags == nil {
			db.Blueprints[i].Tags = []string{}
		}
	}
	b, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func Find(db Database, name string) (*Blueprint, error) {
	for i := range db.Blueprints {
		if db.Blueprints[i].Name == name {
			return &db.Blueprints[i], nil
		}
	}
	return nil, errors.New("blueprint not found: " + name)
}
