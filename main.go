package main

import (
    "flag"

    "github.com/stephenneal/go-start/ds"

    "github.com/op/go-logging"
)

func main() {
	// Parse CLI parameters
    var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "Turn verbose logging on/off")
    flag.Parse()

    defaultLevel := logging.INFO;
    if verbose {
        defaultLevel = logging.DEBUG
    }
    logging.SetLevel(defaultLevel, "")
    var log = logging.MustGetLogger("main")

    log.Debug("Connect to the database")
    mgr := ds.ConnectDb("/tmp/GoStartDb")
    mgr.PrintCols()
    defer mgr.CloseDb()
}

