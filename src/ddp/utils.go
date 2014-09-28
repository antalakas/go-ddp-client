package ddp

import (
  "fmt"
  "os"
)

func checkError(message string, err error) {
  if err != nil {
    fmt.Println("Fatal error in", message, ": " + err.Error())
    os.Exit(1)
  }
}