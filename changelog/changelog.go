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

package changelog

import (
	"bytes"
	"fmt"
	"time"
)

type Entry struct {
	Version string
	Date    time.Time
	Added   []string
	Changed []string
	Fixed   []string
}

func (e Entry) Render() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("## %s - %s\n\n", e.Version, e.Date.Format("2006-01-02")))
	if len(e.Added) > 0 {
		buf.WriteString("### Added\n")
		for _, a := range e.Added {
			buf.WriteString(fmt.Sprintf("- %s\n", a))
		}
		buf.WriteString("\n")
	}
	if len(e.Changed) > 0 {
		buf.WriteString("### Changed\n")
		for _, c := range e.Changed {
			buf.WriteString(fmt.Sprintf("- %s\n", c))
		}
		buf.WriteString("\n")
	}
	if len(e.Fixed) > 0 {
		buf.WriteString("### Fixed\n")
		for _, f := range e.Fixed {
			buf.WriteString(fmt.Sprintf("- %s\n", f))
		}
		buf.WriteString("\n")
	}
	return buf.String()
}
