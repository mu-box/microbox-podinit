package main

import (
  "fmt"
  "os"
  "os/exec"
  "syscall"
  "time"
)

const (
  READY_FILE    = "/var/podinit/ready"
  LOOP_INTERVAL = time.Second
)

// The basic idea here, is that we need to wait to execute the provided
// command until the environment is ready. 
// 
// After validating the command, we'll look for the READY_FILE. If it
// doesn't exist, we'll just keep looping until it does.
// 
// Once the READY_FILE exists, we exec the provided command, using the
// exec syscall, essentially replacing this binary right in place.
func main() {
  // ensure this wasn't called without a command to run
  if len(os.Args) < 2 {
    fmt.Println("Sorry, you must provide a command to run")
    os.Exit(1)
  }
  
  // the binary to execute
  bin := os.Args[1]
  // the executable and all arguments
  args := os.Args[1:]

  // find the actual path of the specified binary
  path, err := exec.LookPath(bin)
  if err != nil {
    fmt.Printf("Error: '%s' not found in $PATH\n", bin)
    os.Exit(1)
  }
  
  // wait until the container is ready
  for {
    if _, err := os.Stat(READY_FILE); os.IsNotExist(err) {
      // keep sleeping
      <- time.After(LOOP_INTERVAL)
    } else {
      // we're ready, break out!
      break
    }
  }
  
  // todo: get secrets from vault, inject them into evars

  // exec the provided command, replacing this process
  // https://groob.io/posts/golang-execve/
  err = syscall.Exec(path, args, os.Environ())
  if err != nil {
    fmt.Printf("Failed to run '%s': %s\n", args, err)
    os.Exit(1)
  }
}
