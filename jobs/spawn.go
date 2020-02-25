package jobs

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/phomer/scheduler/accounts"
)

var output_flags = os.O_RDWR | os.O_CREATE
var output_perms os.FileMode = 0700

// TODO: Take a command and turn it into an active job
func Spawn(account *accounts.Account, job *ActiveJob) error {
	command := job.Cmd

	fmt.Printf("Spawn user: %s cmd: %s args: %v\n", account.Username, command.Cmd, command.Args)

	attributes := Attributes(account, command)

	pid, err := syscall.ForkExec(command.Cmd, command.Args, attributes)
	if err != nil {
		return err
	}

	fmt.Println("Spawned Job", pid)

	job = CheckStatus(pid, job) // Might be done already

	active := NewActive()
	active.AddJob(pid, job)

	return nil
}

// Set up the attributes for the process
func Attributes(account *accounts.Account, command *Command) *syscall.ProcAttr {

	//cwd := account.Directory // Set at registration
	cwd := ""              // Set at registration
	env_vars := []string{} // Pass through stuff for config reasons?

	//jobid := command.JobId

	filepath := command.Filepath
	output := OutputFile(filepath)

	if output == nil {
		// If we can't open the output file, log the contents to the server
		fmt.Println("Output file reset to Server's stderr")
		output = os.Stderr
	}

	// TODO: Not the server's stdin, but nil ...
	files := []uintptr{os.Stdin.Fd(), output.Fd(), output.Fd()} // Stop in, combine out+err

	sys := &syscall.SysProcAttr{}
	if account.Uid != uint32(os.Getuid()) || account.Gid != uint32(os.Getgid()) {
		fmt.Println("Running as user's account")
		sys.Credential = &syscall.Credential{Uid: account.Uid, Gid: account.Gid}
	}

	proc_attr := &syscall.ProcAttr{
		Dir:   cwd,
		Env:   env_vars,
		Files: files,
		Sys:   sys,
	}

	return proc_attr
}

func OutputFilepath(path string, username string, jobid int) string {
	return filepath.Join(path, fmt.Sprintf("%s-%d.output", username, jobid))
}

func OutputFile(filepath string) *os.File {
	file, err := os.OpenFile(filepath, output_flags, output_perms)
	if err != nil {
		fmt.Println("Failed to Openfile", filepath, err)
		return nil
	}

	return file
}
