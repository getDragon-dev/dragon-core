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

package templates

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Context map[string]any

// RenderDir processes *.tmpl files in srcDir into dstDir using the context.
func RenderDir(srcDir, dstDir string, ctx Context) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(srcDir, path)
		out := filepath.Join(dstDir, strings.TrimSuffix(rel, ".tmpl"))
		if info.IsDir() {
			return os.MkdirAll(out, 0o755)
		}
		if strings.HasSuffix(path, ".tmpl") {
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			t, err := template.New(rel).Parse(string(b))
			if err != nil {
				return err
			}
			f, err := os.Create(out)
			if err != nil {
				return err
			}
			defer f.Close()
			return t.Execute(f, ctx)
		}
		// copy non-template files as-is
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(out, b, 0o644)
	})
}
