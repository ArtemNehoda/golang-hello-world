package ports

// Logger is the logging interface.
// Any component that needs to emit log output depends on this interface,
// I put it in the ports dir because prefered golang approach is to put interfaces in the packages that depend on them,
// but at the same time I want to avoid replicating logger interface on all parts of the application.
// So, this is a bit of a compromise.
type Logger interface {
	Printf(format string, v ...any)
	Println(v ...any)
	Fatalf(format string, v ...any)
	Fatalln(v ...any)
}
