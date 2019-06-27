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

package distribution_test

import (
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/dist-zip-cnb/distribution"
	"github.com/cloudfoundry/jvm-application-cnb/jvmapplication"
	"github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/test"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestBuild(t *testing.T) {
	spec.Run(t, "Build", func(t *testing.T, _ spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		var f *test.BuildFactory

		it.Before(func() {
			f = test.NewBuildFactory(t)
		})

		it("returns false with no jvm-application dependency", func() {
			_, ok, err := distribution.NewDistribution(f.Build)

			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(ok).To(BeFalse())
		})

		it("returns false with zero scripts", func() {
			test.TouchFile(t, f.Build.Application.Root, "application-0.0.1", "bin", "bravo.bat")

			f.AddBuildPlan(jvmapplication.Dependency, buildplan.Dependency{})

			_, ok, err := distribution.NewDistribution(f.Build)

			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(ok).To(BeFalse())
		})

		it("returns false with two scripts", func() {
			f.AddBuildPlan(jvmapplication.Dependency, buildplan.Dependency{})
			test.TouchFile(t, f.Build.Application.Root, "application-0.0.1", "bin", "alpha")
			test.TouchFile(t, f.Build.Application.Root, "application-0.0.1", "bin", "bravo.bat")
			test.TouchFile(t, f.Build.Application.Root, "application-0.0.1", "bin", "charlie")

			_, ok, err := distribution.NewDistribution(f.Build)

			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(ok).To(BeFalse())
		})

		it("returns true with one script", func() {
			f.AddBuildPlan(jvmapplication.Dependency, buildplan.Dependency{})
			test.TouchFile(t, f.Build.Application.Root, "application-0.0.1", "bin", "alpha")
			test.TouchFile(t, f.Build.Application.Root, "application-0.0.1", "bin", "bravo.bat")

			_, ok, err := distribution.NewDistribution(f.Build)

			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(ok).To(BeTrue())
		})

		it("contributes command", func() {
			f.AddBuildPlan(jvmapplication.Dependency, buildplan.Dependency{})
			test.TouchFile(t, f.Build.Application.Root, "application-0.0.1", "bin", "alpha")
			test.TouchFile(t, f.Build.Application.Root, "application-0.0.1", "bin", "bravo.bat")

			d, _, err := distribution.NewDistribution(f.Build)

			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(d.Contribute()).To(Succeed())

			command := filepath.Join(f.Build.Application.Root, "application-0.0.1", "bin", "alpha")

			g.Expect(f.Build.Layers).To(test.HaveApplicationMetadata(layers.Metadata{
				Processes: layers.Processes{
					{"dist-zip", command},
					{"task", command},
					{"web", command},
				},
			}))
		})
	}, spec.Report(report.Terminal{}))
}
