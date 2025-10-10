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
	"fmt"
	"os"
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

func Load(path string) (Database, error) {
	var db Database
	b, err := os.ReadFile(path)
	if err != nil {
		return db, err
	}
	if err := json.Unmarshal(b, &db); err != nil {
		return db, err
	}
	return db, nil
}

func Save(path string, db Database) error {
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
	return nil, fmt.Errorf("blueprint %q not found", name)
}
