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

package semver

import (
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	Major int
	Minor int
	Patch int
	Pre   string
}

func Parse(v string) (Version, error) {
	var ver Version
	if v == "" {
		return ver, fmt.Errorf("empty version")
	}
	parts := strings.SplitN(v, "-", 2)
	core := parts[0]
	if len(parts) == 2 {
		ver.Pre = parts[1]
	}
	nums := strings.Split(core, ".")
	if len(nums) != 3 {
		return ver, fmt.Errorf("semver must have major.minor.patch: %s", v)
	}
	var err error
	if ver.Major, err = strconv.Atoi(nums[0]); err != nil {
		return ver, err
	}
	if ver.Minor, err = strconv.Atoi(nums[1]); err != nil {
		return ver, err
	}
	if ver.Patch, err = strconv.Atoi(nums[2]); err != nil {
		return ver, err
	}
	return ver, nil
}

func (v Version) String() string {
	if v.Pre != "" {
		return fmt.Sprintf("%d.%d.%d-%s", v.Major, v.Minor, v.Patch, v.Pre)
	}
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v Version) Bump(kind string) Version {
	nv := v
	switch kind {
	case "major":
		nv.Major++
		nv.Minor = 0
		nv.Patch = 0
		nv.Pre = ""
	case "minor":
		nv.Minor++
		nv.Patch = 0
		nv.Pre = ""
	default:
		nv.Patch++
		nv.Pre = ""
	}
	return nv
}
