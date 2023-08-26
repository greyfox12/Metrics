package logmy

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger = zap.NewNop()

// Initialize инициализирует синглтон логера с необходимым уровнем логирования.
func Initialize(level string) error {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	// создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout", "./agent_logs.txt"}
	cfg.DisableCaller = false
	cfg.DisableStacktrace = false
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// устанавливаем уровень
	cfg.Level = lvl
	//	cfg.EncoderConfig = //"\"callerKey\": \"caller\""
	// создаём логер на основе конфигурации
	zl, err := cfg.Build( /*zap.AddCaller()*/ )
	if err != nil {
		return err
	}
	// устанавливаем синглтон
	Log = zl
	return nil
}

func OutLog(errorStr error) {
	Log.Info("Error:",
		zap.String("Message", fmt.Sprint(errorStr)),
	)
	Log.Sync()
}
