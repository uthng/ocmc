package config

import (
    //"fmt"

    "github.com/spf13/viper"

    "github.com/mitchellh/mapstructure"

    "github.com/uthng/ocmc/types"

)

// ReadClusterConfigFromFile reads configurations of clusters from
// configuration file
func ReadClusterConfigFromFile() *types.PageClusterData {
    data := &types.PageClusterData{}

    configClusters := viper.Get("clusters").([]interface{})
    for _, cluster := range configClusters {
        m := make(map[string]interface{})
        for k, v := range cluster.(map[interface{}]interface{}) {
            m[k.(string)] = v
        }
        config := types.ClusterConfig{}
        mapstructure.Decode(m, &config)
        data.Configs = append(data.Configs, config)
    }

    //fmt.Printf("%v\n", data.Configs)
    return data
}

// GetClusterConfig returns configuration of the cluster corresponding
// to the given name
func GetClusterConfig(name string, data *types.PageClusterData) types.ClusterConfig {
    for _, c := range data.Configs {
        if c.Name == name {
            return c
        }
    }

    return types.ClusterConfig{}
}


