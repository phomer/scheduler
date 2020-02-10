package jobs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/phomer/scheduler/accounts"
)

var output_flags = os.O_RDWR | os.O_CREATE
var output_perms os.FileMode = 0777

func Spawn(root interface{}, account *accounts.Account, command string, args []string) {
	fmt.Println("Spawning off Process", command, " for ", account.Username)

	attributes := Attributes(root, account)

	pid, err := syscall.ForkExec(command, args, attributes)
	if err != nil {
		log.Fatal("Spawn Died ", err)
	}

	fmt.Println("Spawn Completed as ", pid)
}

// Set up the attributes for the process
func Attributes(root interface{}, account *accounts.Account) *syscall.ProcAttr {

	//cwd := account.Directory // Set at registration
	cwd := "" // Set at registration

	env_vars := []string{} // Pass through stuff for config reasons?

	jobid := account.NextId
	account.IncrementId()

	output := OutputFile(root, account.Username, jobid)

	// TODO: Not the server's stdin, but nil ...
	files := []uintptr{os.Stdin.Fd(), output.Fd(), output.Fd()} // Stop in, combine out+err

	sys := &syscall.SysProcAttr{}
	if account.Uid != uint32(os.Getuid()) || account.Gid != uint32(os.Getgid()) {
		fmt.Println("Running as user's account")
		sys.Credential = &syscall.Credential{Uid: account.Uid, Gid: account.Gid}
	}

	return &syscall.ProcAttr{
		Dir:   cwd,
		Env:   env_vars,
		Files: files,
		Sys:   sys,
	}
}

func OutputFile(root interface{}, username string, jobid uint) *os.File {

	data_path := "data" // TODO: Pick this out of root

	path := filepath.Join(data_path, fmt.Sprintf("%s-%d.output", username, jobid))

	file, err := os.OpenFile(path, output_flags, output_perms)
	if err != nil {
		log.Fatal("Output File", err)
	}

	return file
}
