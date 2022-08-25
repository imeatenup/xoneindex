package log

import (
	"context"
	"fmt"
	"log"
	"os"
)

type Trace struct{}

var tracelog *log.Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

// Infof info...
func Infof(ctx context.Context, format string, v ...any) {
	tracelog.Output(2, fmt.Sprintf("INFO "+format+"	trace:%+v", append(v, ctx.Value(Trace{}))...))
}

// Errorf error...
func Errorf(ctx context.Context, format string, v ...any) {
	tracelog.Output(2, fmt.Sprintf("ERROR "+format+"	trace:%+v", append(v, ctx.Value(Trace{}))...))
}

// Debugf debug...
func Debugf(ctx context.Context, format string, v ...any) {
	tracelog.Output(2, fmt.Sprintf("DEBUG "+format+"	trace:%+v", append(v, ctx.Value(Trace{}))...))
}
