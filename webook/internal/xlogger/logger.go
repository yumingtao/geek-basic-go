package xlogger

import "go.uber.org/zap"

// Logger 在main函数里初始化 Logger=xxxx，还是强耦合
var Logger *zap.Logger

var CommonLogger *zap.Logger
var SensitiveLogger *zap.Logger
