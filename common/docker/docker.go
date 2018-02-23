package docker

import (
    "fmt"
    "errors"
    "strconv"
    //"bufio"
    //"io/ioutil"

    "golang.org/x/net/context"
    "github.com/docker/docker/client"
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/api/types/swarm"

    ocmc_types "github.com/uthng/ocmc/types"
    "github.com/uthng/common/ssh"
    "github.com/uthng/common/docker"
)

var (
    ErrNetworkIDNotFound = errors.New("Network ID not found")
    ErrServiceIDNotFound = errors.New("Service ID not found")
    ErrServiceNameNotFound = errors.New("Service Name not found")
    ErrNodeIDNotFound = errors.New("Node ID not found")
)

// NewDockerClient initializes a docker client to remote cluster
// following authentication configuration
func NewDockerClient(config ocmc_types.ConnConfig) (interface{}, error) {
    var client interface{}

    if config.Auth.Type == "ssh" {
        if config.Auth.Kind == "key" {
            sshConfig, err := ssh.NewClientConfigWithKeyFile(config.Auth.Username, config.Auth.SshKey, "", 0, false)
            if err != nil {
                return nil, err
            }

            client, err = docker.NewSSHClient(config.Host + ":" + strconv.Itoa(config.Port), "/var/run/docker.sock", "1.35", sshConfig.ClientConfig)
            if err != nil {
                return nil, err
            }
        }
    }
    return client, nil
}

// GetSwarmServices returns a list of swarm services created in the cluster
func GetSwarmServices (client *client.Client) ([]swarm.Service, error) {
    // Get swarm services
    return client.ServiceList(context.Background(), types.ServiceListOptions{})
}

// FindServiceByID searchs and returns the swarm service corresponding to
// the given ID
func FindSwarmServiceByID (id string, services []swarm.Service) (*swarm.Service, error) {
    for _, service := range services {
        if service.ID == id {
            return &service, nil
        }
    }

    return nil, ErrServiceIDNotFound
}

// FindServiceByName searchs and returns the swarm service corresponding to
// the given name
func FindSwarmServiceByName (name string, services []swarm.Service) (*swarm.Service, error) {
    for _, service := range services {
        if service.Spec.Annotations.Name == name {
            return &service, nil
        }
    }

    return nil, ErrServiceNameNotFound
}

// GetNetworks return a list of network defined in the cluster
func GetNetworks (client *client.Client) ([]types.NetworkResource, error) {
    return client.NetworkList(context.Background(), types.NetworkListOptions{})
}

// FindNetworkByID searchs and returns the NetworkResource corresponding to
// the given ID
func FindNetworkByID (id string, networks []types.NetworkResource) (*types.NetworkResource, error) {
    for _, net := range networks {
        if net.ID == id {
            return &net, nil
        }
    }

    return nil, ErrNetworkIDNotFound
}

// GetSwarmTasks return a list of swarm containers in the cluster
func GetSwarmTasks (client *client.Client) ([]swarm.Task, error) {
    return client.TaskList(context.Background(), types.TaskListOptions{})
}

// FindSwarmTasksByServiceID searchs and returns the tasks
// related to the given service ID
func FindSwarmTasksByServiceID (serviceId string, tasks []swarm.Task) ([]swarm.Task) {
    var result []swarm.Task

    for _, task := range tasks {
        if task.ServiceID == serviceId {
            result = append(result, task)
        }
    }

    return result
}

// GetSwarmNodes return a list of swarm nodes in the cluster
func GetSwarmNodes (client *client.Client) ([]swarm.Node, error) {
    return client.NodeList(context.Background(), types.NodeListOptions{})
}

// FindNodeByServiceID returns the node corresponding to the given ID
func FindSwarmNodeByID (id string, nodes []swarm.Node) (*swarm.Node, error) {
    for _, node := range nodes {
        if node.ID == id {
            return &node, nil
        }
    }

    return nil, ErrNodeIDNotFound
}

// GetContainers return a list of containers running on 
// the current host of the cluster (like docker ps)
func GetContainers (client *client.Client) ([]types.Container, error) {
    return client.ContainerList(context.Background(), types.ContainerListOptions{})
}

// CreateExec creates an exec instance and return exec id
//func CreateExec (client *client.Client, id string, cmd []string) (*string, context.Context, error) {
    //config := types.ExecConfig {
        ////AttachStdout: true,
        ////AttachStderr: true,
        ////Tty: false,
        ////Detach: false,
        ////DetachKeys: "ctrl-p,ctrl-q",
        //Cmd: cmd,
    //}
    //// Create a exec instance
    //ctx := context.Background()
    //res, err := client.ContainerExecCreate(context.Background(), id, config)
    //if err != nil {
        //return nil, ctx, err
    //}

    //fmt.Println("exec id ", res.ID)
    //return &res.ID, ctx, nil
//}


// StartExec start exec process and attach it to a reader.
// It returns a bufio reader for command output
//
// This function takes temporairly cmd in argument but
// need to be removed in the next release of docker
// because ContainerExecAttach does not take types.ExecConfig anymore
// Instead, it takes types.ExecStartCheck
//func StartExec (client *client.Client, ctx context.Context, execId string, cmd []string) (*bufio.Reader, error) {
    //config := types.ExecConfig {
        //AttachStdout: true,
        ////AttachStderr: true,
        //Tty: false,
        //Detach: false,
        ////DetachKeys: "ctrl-p,ctrl-q",
        ////Cmd: cmd,
    //}

    //// Create a exec instance
    ////err := client.ContainerExecStart(ctx, execId, types.ExecStartCheck{Detach: false, Tty: false})
    ////if err != nil {
        ////fmt.Println("error start exec")
        ////return nil, err
    ////}

    //// Attach an exec
    //res, err := client.ContainerExecAttach(ctx, execId, config)
    //if err != nil {
        //fmt.Println("error attach exec")
        //return nil, err
    //}
    //defer res.Close()

    //line, _, err := res.Reader.ReadLine()
    //fmt.Println("line ", string(line))
    //return res.Reader, nil
//}

func ExecCommand(client *client.Client, cid string, cmd []string) ([]byte, error) {
    ctx := context.Background()

    id, err := client.ContainerExecCreate(ctx, cid,
                types.ExecConfig{
                    //WorkingDir:   "/tmp",
                    //Env:          strslice.StrSlice([]string{"FOO=BAR"}),
                    AttachStdout: true,
                    Cmd: cmd,
                    //Cmd:          strslice.StrSlice([]string{"sh", "-c", "env"}),
                })
    if err != nil {
        fmt.Println("error create exec ", err)
        return nil, err
    }

    fmt.Println("id ", id)

    //insp, err := client.ContainerExecInspect(ctx, id.ID)
    //if err != nil {
         //fmt.Println("error inspect exec ", err)
        //return nil, err
    //}

    //err = client.ContainerExecStart(ctx, id.ID, types.ExecStartCheck{
                    //Detach: false,
                    //Tty:    true,
                //})
    //if err != nil {
         //fmt.Println("error start exec ", err)
        //return nil, err
    //}

    //fmt.Println("inspect ", insp)
    //fmt.Println("cid ", cid)
    resp, err := client.ContainerExecAttach(ctx, id.ID,
                types.ExecStartCheck{
                    Detach: false,
                    Tty:    false,
                })
    if err != nil {
        fmt.Println("error attach exec ", err)
        return nil, err
    }
    defer resp.Close()

    fmt.Println("response ", resp)
    line, _, err := resp.Reader.ReadLine()
    fmt.Println("response ", string(line), err)


    //r, err := ioutil.ReadAll(resp.Reader)
    //if err != nil {
        //fmt.Println("error readall ", err)
        //return nil, err
    //}

    return line, err
}
