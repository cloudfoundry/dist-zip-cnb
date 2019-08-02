/*
 * Copyright 2018-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package distribution

import (
	"path/filepath"
	"regexp"

	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/layers"
)

type Distribution struct {
	launcher string
	layers   layers.Layers
}

func (d Distribution) Contribute() error {
	return d.layers.WriteApplicationMetadata(layers.Metadata{
		Processes: layers.Processes{
			{"dist-zip", d.launcher},
			{"task", d.launcher},
			{"web", d.launcher},
		},
	})
}

func NewDistribution(build build.Build) (Distribution, bool, error) {
	l, err := launcher(build.Application.Root)
	if err != nil {
		return Distribution{}, false, err
	}

	if len(l) != 1 {
		return Distribution{}, false, nil
	}

	return Distribution{
		l[0],
		build.Layers,
	}, true, nil
}

func launcher(root string) ([]string, error) {
	c, err := filepath.Glob(filepath.Join(root, "*", "bin", "*"))
	if err != nil {
		return nil, err
	}

	bat := regexp.MustCompile(`.+\.bat`)

	i := 0
	for _, s := range c {
		if !bat.MatchString(s) {
			c[i] = s
			i++
		}
	}

	return c[:i], nil
}
