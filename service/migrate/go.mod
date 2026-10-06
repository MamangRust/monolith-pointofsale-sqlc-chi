module github.com/MamangRust/monolith-point-of-sale-migrate

go 1.25.0

require (
	github.com/MamangRust/monolith-point-of-sale-pkg v1.0.2
	github.com/lib/pq v1.10.9
	github.com/pressly/goose/v3 v3.24.3
	github.com/spf13/viper v1.20.1
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.3.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/sagikazarmark/locafero v0.9.0 // indirect
	github.com/sethvargo/go-retry v0.3.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.14.0 // indirect
	github.com/spf13/cast v1.9.2 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/MamangRust/monolith-point-of-sale-pb => ../../pb
	github.com/MamangRust/monolith-point-of-sale-pkg => ../../pkg
	github.com/MamangRust/monolith-point-of-sale-shared => ../../shared
)
