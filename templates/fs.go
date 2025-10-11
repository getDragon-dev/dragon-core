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
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

	sprig "github.com/Masterminds/sprig/v3"
)

type Context map[string]any

// RenderDir parses all files under src, applies text/template with Sprig funcs, writes to dst.
func RenderDir(src, dst string, ctx Context) error {
	// collect template files
	var files []string
	if err := filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		files = append(files, p)
		return nil
	}); err != nil {
		return err
	}

	// parse with Sprig
	t := template.New("").Funcs(sprig.TxtFuncMap())
	for _, f := range files {
		rel, _ := filepath.Rel(src, f)
		b, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		if _, err := t.New(rel).Parse(string(b)); err != nil {
			return err
		}
	}

	// execute each template into target tree (strip .tmpl suffix)
	for _, f := range files {
		rel, _ := filepath.Rel(src, f)
		outRel := rel
		if filepath.Ext(outRel) == ".tmpl" {
			outRel = outRel[:len(outRel)-len(".tmpl")]
		}

		// render path (use same Sprig funcs)
		pathTmpl, err := template.New("path").Funcs(sprig.TxtFuncMap()).Parse(outRel)
		if err != nil {
			return err
		}
		var pathBuf bytes.Buffer
		if err := pathTmpl.Execute(&pathBuf, ctx); err != nil {
			return fmt.Errorf("render path %s: %w", outRel, err)
		}
		outRel = pathBuf.String()

		// render file content
		var buf bytes.Buffer
		if err := t.ExecuteTemplate(&buf, rel, ctx); err != nil {
			return err
		}
		outPath := filepath.Join(dst, outRel)
		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(outPath, buf.Bytes(), 0o644); err != nil {
			return err
		}
	}

	return nil
}
