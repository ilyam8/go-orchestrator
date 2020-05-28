package file

import (
	"testing"

	"github.com/netdata/go-orchestrator/job/confgroup"
	"github.com/netdata/go-orchestrator/module"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	const (
		jobDef = 11
		cfgDef = 22
		modDef = 33
	)
	tests := map[string]func(t *testing.T, tmp *tmpDir){
		"static, default: +job +conf +module": func(t *testing.T, tmp *tmpDir) {
			reg := confgroup.Registry{
				"module": {
					UpdateEvery:        modDef,
					AutoDetectionRetry: modDef,
					Priority:           modDef,
				},
			}
			cfg := staticConfig{
				Default: confgroup.Default{
					UpdateEvery:        cfgDef,
					AutoDetectionRetry: cfgDef,
					Priority:           cfgDef,
				},
				Jobs: []confgroup.Config{
					{
						"name":                "name",
						"update_every":        jobDef,
						"autodetection_retry": jobDef,
						"priority":            jobDef,
					},
				},
			}
			filename := tmp.join("module.conf")
			tmp.writeYAML(filename, cfg)

			expected := &confgroup.Group{
				Source: filename,
				Configs: []confgroup.Config{
					{
						"name":                "name",
						"module":              "module",
						"update_every":        jobDef,
						"autodetection_retry": jobDef,
						"priority":            jobDef,
					},
				},
			}

			groups, err := parse(reg, filename)

			require.NoError(t, err)
			assert.Equal(t, expected, groups)
		},
		"static, default: +job +conf +module (merge all)": func(t *testing.T, tmp *tmpDir) {
			reg := confgroup.Registry{
				"module": {
					Priority: modDef,
				},
			}
			cfg := staticConfig{
				Default: confgroup.Default{
					AutoDetectionRetry: cfgDef,
				},
				Jobs: []confgroup.Config{
					{
						"name":         "name",
						"update_every": jobDef,
					},
				},
			}
			filename := tmp.join("module.conf")
			tmp.writeYAML(filename, cfg)

			expected := &confgroup.Group{
				Source: filename,
				Configs: []confgroup.Config{
					{
						"name":                "name",
						"module":              "module",
						"update_every":        jobDef,
						"autodetection_retry": cfgDef,
						"priority":            modDef,
					},
				},
			}

			groups, err := parse(reg, filename)

			require.NoError(t, err)
			assert.Equal(t, expected, groups)
		},
		"static, default: -job +conf +module": func(t *testing.T, tmp *tmpDir) {
			reg := confgroup.Registry{
				"module": {
					UpdateEvery:        modDef,
					AutoDetectionRetry: modDef,
					Priority:           modDef,
				},
			}
			cfg := staticConfig{
				Default: confgroup.Default{
					UpdateEvery:        cfgDef,
					AutoDetectionRetry: cfgDef,
					Priority:           cfgDef,
				},
				Jobs: []confgroup.Config{
					{
						"name": "name",
					},
				},
			}
			filename := tmp.join("module.conf")
			tmp.writeYAML(filename, cfg)

			expected := &confgroup.Group{
				Source: filename,
				Configs: []confgroup.Config{
					{
						"name":                "name",
						"module":              "module",
						"update_every":        cfgDef,
						"autodetection_retry": cfgDef,
						"priority":            cfgDef,
					},
				},
			}

			groups, err := parse(reg, filename)

			require.NoError(t, err)
			assert.Equal(t, expected, groups)
		},
		"static, default: -job -conf +module": func(t *testing.T, tmp *tmpDir) {
			reg := confgroup.Registry{
				"module": {
					UpdateEvery:        modDef,
					AutoDetectionRetry: modDef,
					Priority:           modDef,
				},
			}
			cfg := staticConfig{
				Jobs: []confgroup.Config{
					{
						"name": "name",
					},
				},
			}
			filename := tmp.join("module.conf")
			tmp.writeYAML(filename, cfg)

			expected := &confgroup.Group{
				Source: filename,
				Configs: []confgroup.Config{
					{
						"name":                "name",
						"module":              "module",
						"autodetection_retry": modDef,
						"priority":            modDef,
						"update_every":        modDef,
					},
				},
			}

			groups, err := parse(reg, filename)

			require.NoError(t, err)
			assert.Equal(t, expected, groups)
		},
		"static, default: -job -conf -module (+global)": func(t *testing.T, tmp *tmpDir) {
			reg := confgroup.Registry{
				"module": {},
			}
			cfg := staticConfig{
				Jobs: []confgroup.Config{
					{
						"name": "name",
					},
				},
			}
			filename := tmp.join("module.conf")
			tmp.writeYAML(filename, cfg)

			expected := &confgroup.Group{
				Source: filename,
				Configs: []confgroup.Config{
					{
						"name":                "name",
						"module":              "module",
						"autodetection_retry": module.AutoDetectionRetry,
						"priority":            module.Priority,
						"update_every":        module.UpdateEvery,
					},
				},
			}

			groups, err := parse(reg, filename)

			require.NoError(t, err)
			assert.Equal(t, expected, groups)
		},
		"sd, default: +job +module": func(t *testing.T, tmp *tmpDir) {
			reg := confgroup.Registry{
				"sd_module": {
					UpdateEvery:        modDef,
					AutoDetectionRetry: modDef,
					Priority:           modDef,
				},
			}
			cfg := sdConfig{
				{
					"name":                "name",
					"module":              "sd_module",
					"update_every":        jobDef,
					"autodetection_retry": jobDef,
					"priority":            jobDef,
				},
			}
			filename := tmp.join("module.conf")
			tmp.writeYAML(filename, cfg)

			expected := &confgroup.Group{
				Source: filename,
				Configs: []confgroup.Config{
					{
						"module":              "sd_module",
						"name":                "name",
						"update_every":        jobDef,
						"autodetection_retry": jobDef,
						"priority":            jobDef,
					},
				},
			}

			groups, err := parse(reg, filename)

			require.NoError(t, err)
			assert.Equal(t, expected, groups)
		},
		"sd, default: -job +module": func(t *testing.T, tmp *tmpDir) {
			reg := confgroup.Registry{
				"sd_module": {
					UpdateEvery:        modDef,
					AutoDetectionRetry: modDef,
					Priority:           modDef,
				},
			}
			cfg := sdConfig{
				{
					"name":   "name",
					"module": "sd_module",
				},
			}
			filename := tmp.join("module.conf")
			tmp.writeYAML(filename, cfg)

			expected := &confgroup.Group{
				Source: filename,
				Configs: []confgroup.Config{
					{
						"name":                "name",
						"module":              "sd_module",
						"update_every":        modDef,
						"autodetection_retry": modDef,
						"priority":            modDef,
					},
				},
			}

			groups, err := parse(reg, filename)

			require.NoError(t, err)
			assert.Equal(t, expected, groups)
		},
		"sd, default: -job -module (+global)": func(t *testing.T, tmp *tmpDir) {
			reg := confgroup.Registry{
				"sd_module": {},
			}
			cfg := sdConfig{
				{
					"name":   "name",
					"module": "sd_module",
				},
			}
			filename := tmp.join("module.conf")
			tmp.writeYAML(filename, cfg)

			expected := &confgroup.Group{
				Source: filename,
				Configs: []confgroup.Config{
					{
						"name":                "name",
						"module":              "sd_module",
						"update_every":        module.UpdateEvery,
						"autodetection_retry": module.AutoDetectionRetry,
						"priority":            module.Priority,
					},
				},
			}

			groups, err := parse(reg, filename)

			require.NoError(t, err)
			assert.Equal(t, expected, groups)
		},
		"sd, job has no 'module' or 'module' is empty": func(t *testing.T, tmp *tmpDir) {
			reg := confgroup.Registry{
				"sd_module": {},
			}
			cfg := sdConfig{
				{
					"name": "name",
				},
			}
			filename := tmp.join("module.conf")
			tmp.writeYAML(filename, cfg)

			expected := &confgroup.Group{
				Source:  filename,
				Configs: []confgroup.Config{},
			}

			groups, err := parse(reg, filename)

			require.NoError(t, err)
			assert.Equal(t, expected, groups)
		},
		"conf registry has no module": func(t *testing.T, tmp *tmpDir) {
			reg := confgroup.Registry{
				"sd_module": {},
			}
			cfg := sdConfig{
				{
					"name":   "name",
					"module": "module",
				},
			}
			filename := tmp.join("module.conf")
			tmp.writeYAML(filename, cfg)

			expected := &confgroup.Group{
				Source:  filename,
				Configs: []confgroup.Config{},
			}

			groups, err := parse(reg, filename)

			require.NoError(t, err)
			assert.Equal(t, expected, groups)
		},
		"empty file": func(t *testing.T, tmp *tmpDir) {
			reg := confgroup.Registry{
				"module": {},
			}

			filename := tmp.createFile("empty-*")
			groups, err := parse(reg, filename)

			require.NoError(t, err)
			assert.Nil(t, groups)
		},
		"unknown format": func(t *testing.T, tmp *tmpDir) {
			reg := confgroup.Registry{}

			filename := tmp.createFile("unknown-format-*")
			tmp.writeYAML(filename, "unknown")
			_, err := parse(reg, filename)

			assert.Error(t, err)
		},
	}

	for name, scenario := range tests {
		t.Run(name, func(t *testing.T) {
			tmp := newTmpDir(t, "parse-file-*")
			defer tmp.cleanup()
			scenario(t, tmp)
		})
	}
}