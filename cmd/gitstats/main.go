package main

import (
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/stefanoj3/gitstats/internal/presentation/cli/cmd"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	out := os.Stdout

	atomicLevel := zap.NewAtomicLevelAt(zap.InfoLevel)

	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.Lock(out),
		atomicLevel,
	))

	if err := buildCLICommand(out, &atomicLevel, logger).Execute(); err != nil {
		logger.Fatal("failed to run command", zap.Error(err))
	}
}

func buildCLICommand(out io.Writer, atomicLevel *zap.AtomicLevel, logger *zap.Logger) *cobra.Command {
	root := cmd.NewRootCommand(atomicLevel)
	root.AddCommand(cmd.NewCollectCommand(logger))
	root.AddCommand(cmd.NewVersionCommand(out))

	return root
}
