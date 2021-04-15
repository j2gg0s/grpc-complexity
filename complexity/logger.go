package complexity

import (
	"context"

	"google.golang.org/grpc/grpclog"
)

// LoggerV2 does underlying logging work for grpclog.
type Logger interface {
	// Info logs to INFO log. Arguments are handled in the manner of fmt.Print.
	Info(ctx context.Context, args ...interface{})
	// Infoln logs to INFO log. Arguments are handled in the manner of fmt.Println.
	Infoln(ctx context.Context, args ...interface{})
	// Infof logs to INFO log. Arguments are handled in the manner of fmt.Printf.
	Infof(ctx context.Context, format string, args ...interface{})
	// Warning logs to WARNING log. Arguments are handled in the manner of fmt.Print.
	Warning(ctx context.Context, args ...interface{})
	// Warningln logs to WARNING log. Arguments are handled in the manner of fmt.Println.
	Warningln(ctx context.Context, args ...interface{})
	// Warningf logs to WARNING log. Arguments are handled in the manner of fmt.Printf.
	Warningf(ctx context.Context, format string, args ...interface{})
	// Error logs to ERROR log. Arguments are handled in the manner of fmt.Print.
	Error(ctx context.Context, args ...interface{})
	// Errorln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
	Errorln(ctx context.Context, args ...interface{})
	// Errorf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
	Errorf(ctx context.Context, format string, args ...interface{})
	// Fatal logs to ERROR log. Arguments are handled in the manner of fmt.Print.
	// gRPC ensures that all Fatal logs will exit with os.Exit(1).
	// Implementations may also call os.Exit() with a non-zero exit code.
	Fatal(ctx context.Context, args ...interface{})
	// Fatalln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
	// gRPC ensures that all Fatal logs will exit with os.Exit(1).
	// Implementations may also call os.Exit() with a non-zero exit code.
	Fatalln(ctx context.Context, args ...interface{})
	// Fatalf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
	// gRPC ensures that all Fatal logs will exit with os.Exit(1).
	// Implementations may also call os.Exit() with a non-zero exit code.
	Fatalf(ctx context.Context, format string, args ...interface{})
	// V reports whether verbosity level l is at least the requested verbose level.
	V(l int) bool
}

type grpcLogger struct {
	logger grpclog.LoggerV2
}

var _ Logger = (*grpcLogger)(nil)

func (l *grpcLogger) Info(ctx context.Context, args ...interface{})   { l.logger.Info(args...) }
func (l *grpcLogger) Infoln(ctx context.Context, args ...interface{}) { l.logger.Infoln(args...) }
func (l *grpcLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}
func (l *grpcLogger) Warning(ctx context.Context, args ...interface{})   { l.logger.Warning(args...) }
func (l *grpcLogger) Warningln(ctx context.Context, args ...interface{}) { l.logger.Warningln(args...) }
func (l *grpcLogger) Warningf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Warningf(format, args...)
}
func (l *grpcLogger) Error(ctx context.Context, args ...interface{})   { l.logger.Error(args...) }
func (l *grpcLogger) Errorln(ctx context.Context, args ...interface{}) { l.logger.Errorln(args...) }
func (l *grpcLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}
func (l *grpcLogger) Fatal(ctx context.Context, args ...interface{})   { l.logger.Fatal(args...) }
func (l *grpcLogger) Fatalln(ctx context.Context, args ...interface{}) { l.logger.Fatalln(args...) }
func (l *grpcLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}
func (l *grpcLogger) V(lvl int) bool { return l.logger.V(lvl) }

func WrapGrpcLogger(logger grpclog.LoggerV2) Logger {
	return &grpcLogger{logger: logger}
}
