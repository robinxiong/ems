package utils

import "os"

var AppRoot, _ = os.Getwd()

type ContextKey string