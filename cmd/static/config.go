package main

import env "github.com/Netflix/go-env"

type config struct {
	ArchiveBucket    string `env:"CRANE_ARCHIVE_BUCKET"`
	ArchiveFolder    string `env:"CRANE_ARCHIVE_FOLDER"`
	DeployBucket     string `env:"CRANE_DEPLOY_BUCKET"`
	StageBucket      string `env:"CRANE_STAGE_BUCKET"`
	ProductionBucket string `env:"CRANE_PRODUCTION_BUCKET"`
	Region           string `env:"CRANE_REGION,default=ap-southeast-2"`
}

func loadConfig() (*config, error) {
	cfg := &config{}
	_, err := env.UnmarshalFromEnviron(cfg)
	return cfg, err
}
