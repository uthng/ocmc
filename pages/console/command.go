package console

import (
    "errors"
    //"strings"
    //"fmt"

    //"github.com/uthng/common/docker"

    "github.com/uthng/ocmc/types"
    //"github.com/uthng/common/ssh"

)

// CmdExec contains container's ID, command to execute
// and and its associated exec process ID
type CmdExec struct {
    // container id
    ContainerId         string
    Command             string
}

//var executedCommands    []*CmdExec

// execCommand executes a command following the type of Node (docker or kubernetes)
//
// It returns a bufio.Reader for command output
func execCommand(cid, cmd string, data *types.PageConsoleData) ([]byte, error) {
    var err error = nil

    //exec := &CmdExec {
        //ContainerId: cid,
        //Command: cmd,
    //}

    //executedCommands = append(executedCommands, exec)

    // if connection type is ssh, we have to launch a ssh shell command
    // if connection type is tls, in this case, probably, docker daemon
    // is configured to use tcp (remotely). So we can use docker command
    // such as ContainerExecCreate & ContainerExecAttach
    if data.Node.Config.Auth.Type == "ssh" {
        if sshClient != nil {
            return sshClient.ExecCommand("docker exec " + cid + " " + cmd)
        } else {
            return nil, errors.New("SSH Client is nil")
        }
    } else if data.Node.Config.Auth.Type == "tls" {
        //return execDockerCommand(exec, data.Node.Client.(*client.Client))
    }

    return nil, err
}

// execDockerCommand creates a docker exec instance and start it.
// If success, it adds the command to executed command list
//
//      - exec: command to execute inside specific container
//      - client: docker client (not ssh client)
//
// It returns a bufio.Reader for command output
//func execDockerCommand (exec *CmdExec, client *docker.Client) ([]byte, error) {
    // Create a exec instance
    //execId, ctx, err := docker.CreateExec(client, exec.ContainerId, strings.Split(exec.Command, " "))
    //if  err != nil {
        //return nil, err
    //}

    //exec.ExecId = *execId

    //// Start an exec
    //reader, err := docker.StartExec(client, ctx, exec.ExecId, strings.Split(exec.Command, " "))
    //if err != nil {
        //return nil, err
    //}

    //// Append this command to executed command list
    //executedCommands = append(executedCommands, exec)

    //res, err := client.ExecCommand(client, exec.ContainerId, strings.Split(exec.Command, " "))
    //return res, err
//}
