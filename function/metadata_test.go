/*
 * Copyright 2018 the original author or authors.
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

package function_test

import (
	"fmt"
	"github.com/buildpack/libbuildpack/application"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack/build"
	"github.com/heroku/libfnbuildpack/function"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestMetadata(t *testing.T) {
	spec.Run(t, "Metadata", func(t *testing.T, _ spec.G, it spec.S) {

		g := NewGomegaWithT(t)
		var testBuild build.Build
		var testApplication application.Application

		it.Before(func() {
			testBuild = build.Build{}
			testApplication,  _ = application.DefaultApplication(testBuild.Logger)
			testBuild.Application = testApplication
		})

		it.After(func() {
			_ = os.Unsetenv(function.ArtifactEnv)
			_ = os.Unsetenv(function.HandlerEnv)
			_ = os.Unsetenv(function.OverrideEnv)

			_ = os.Remove(filepath.Join(testBuild.Application.Root, "metadata.toml"))
		})

		it("returns metadata if metadata.toml exists", func() {

			mdContent := `
artifact = "toml-artifact"
handler = "toml-handler"
override = "toml-override"
`
			filename := filepath.Join(testBuild.Application.Root, "metadata.toml")
			writeMetadataTestFile(t, filename, mdContent)

			actual, ok, err := function.NewMetadata(testBuild.Application, testBuild.Logger)
			g.Expect(ok).To(BeTrue())
			g.Expect(err).NotTo(HaveOccurred())

			g.Expect(actual).To(Equal(function.Metadata{
				Artifact: "toml-artifact",
				Handler:  "toml-handler",
				Override: "toml-override",
			}))
		})

		it("environment variables override metadata.toml", func() {
			mdContent := `
artifact = "toml-artifact"
handler = "toml-handler"
override = "toml-override"
`
			_ = os.Setenv("RIFF_ARTIFACT", "env-artifact")
			_ = os.Setenv("RIFF_OVERRIDE", "env-override")

			filename := filepath.Join(testBuild.Application.Root, "metadata.toml")
			writeMetadataTestFile(t, filename, mdContent)

			actual, ok, err := function.NewMetadata(testBuild.Application, testBuild.Logger)
			g.Expect(ok).To(BeTrue())
			g.Expect(err).NotTo(HaveOccurred())

			g.Expect(actual).To(Equal(function.Metadata{
				Artifact: "env-artifact",
				Handler:  "toml-handler",
				Override: "env-override",
			}))
		})

	}, spec.Report(report.Terminal{}))
}

func writeMetadataTestFile(t *testing.T, filename, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(filename, []byte(fmt.Sprintf(content)), 0644); err != nil {
		t.Fatal(err)
	}
}
