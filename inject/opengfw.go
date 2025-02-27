package inject

import (
	"context"
	"github.com/apernet/OpenGFW/engine"
	"github.com/apernet/OpenGFW/ruleset"
	"github.com/apernet/OpenGFW/ruleset/builtins/geo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"lcf-controller/inject/opengfw"
	"lcf-controller/logger"
	"lcf-controller/pkg/config"
)

var cfg = config.ReadCfg()

var (
	configFile = cfg.OpenGFWConfig.ConfigFilePath
	rulesFile  = cfg.OpenGFWConfig.RulesetFilePath
)

// RunOpenGFW 运行 OpenGFW 引擎
func RunOpenGFW(ctx context.Context) {
	// Config
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		logger.Logger.Fatal("failed to read OpenGFW opengfwCliCfg", zap.Error(err))
	}
	var opengfwCliCfg opengfw.CliConfig
	if err := viper.Unmarshal(&opengfwCliCfg); err != nil {
		logger.Logger.Fatal("failed to parse OpenGFW opengfwCliCfg", zap.Error(err))
	}
	engineConfig, err := opengfwCliCfg.Config()
	if err != nil {
		logger.Logger.Fatal("failed to parse OpenGFW opengfwCliCfg", zap.Error(err))
	}
	defer engineConfig.IO.Close() // Make sure to close IO on exit

	// Ruleset
	rawRs, err := ruleset.ExprRulesFromYAML(rulesFile)
	if err != nil {
		logger.Logger.Fatal("failed to load OpenGFW rules", zap.Error(err))
	}
	rsConfig := &ruleset.BuiltinConfig{
		Logger:               &opengfw.RulesetLogger{},
		GeoMatcher:           geo.NewGeoMatcher(opengfwCliCfg.Ruleset.GeoSite, opengfwCliCfg.Ruleset.GeoIp),
		ProtectedDialContext: engineConfig.IO.ProtectedDialContext,
	}
	rs, err := ruleset.CompileExprRules(rawRs, opengfw.Analyzers, opengfw.Modifiers, rsConfig)
	if err != nil {
		logger.Logger.Fatal("failed to compile OpenGFW rules", zap.Error(err))
	}
	engineConfig.Ruleset = rs

	// Engine
	en, err := engine.NewEngine(*engineConfig)
	if err != nil {
		logger.Logger.Fatal("failed to initialize OpenGFW engine", zap.Error(err))
	}

	logger.Logger.Info("OpenGFW engine started")
	err = en.Run(ctx)
	if err != nil {
		logger.Logger.Fatal("error while running OpenGFW engine", zap.Error(err))
	}
}
