package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dodizzle/chef"
	"github.com/dodizzle/ostack"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
)

var (
	flavorTypes = map[int][]string{}
	knifeCmd    = "/usr/local/bin/knife"
)

func main() {
	// check for KNIFE_PATH envar
	knifePath := getKnifePath()
	// check for SSH_KEY envar
	sshKey := getSSHKey()
	// check for SECRETS_PATH envar
	secretsPath := getSecretsPath()
	// get auth info from knife.rb
	chefConfigs := getChefCredentials(knifePath)
	//  authurl in knife.rb needs to be edited
	newAuthurl := strings.Replace(chefConfigs["openstackauthurl"], "auth/tokens", "", 1)
	nodeName := chefConfigs["nodename"]
	projectName := chefConfigs["openstackprojectname"]
	clientKey := os.Getenv("HOME") + "/.chef" + chefConfigs["clientkey"]
	chefServerURL := chefConfigs["chefserverurl"]
	// prompt for which chef environment to use
	chefEnvironment := chef.GetEnvironments(nodeName, clientKey, chefServerURL)

	opts := gophercloud.AuthOptions{
		IdentityEndpoint: newAuthurl,
		Username:         chefConfigs["openstackusername"],
		Password:         chefConfigs["openstackpassword"],
		DomainName:       chefConfigs["openstackdomainname"],
		TenantName:       projectName,
	}
	provider, err := openstack.AuthenticatedClient(opts)
	//fmt.Println(reflect.TypeOf(provider))
	if err != nil {
		fmt.Println("Error:", err)
	}
	// prompt for OS image to
	imageName := ostack.ListImages(provider)
	// prompt for which available floating ip to use
	ipAddress := ostack.ListIps(provider)
	// prompt for which flavor to use
	flavorChoice := ostack.ListFlavors(provider)
	// prompt for which keypair to use
	sshKeyName := ostack.ListSshkeys(provider)
	// use project name to get the network id
	netID := ostack.GetNetworkID(provider, projectName)
	// prompt for what encrypted_data_bag_secret file to use
	secretsFileName := ostack.GetSecretsFile(secretsPath)
	secretsFile := secretsPath + secretsFileName
	// prompt for the name used for chef and hostname
	hostname := getHostname()
	runKnife(projectName, flavorChoice, hostname, sshKey, netID, chefEnvironment, sshKeyName, imageName, ipAddress, secretsFile)
}

func getSecretsPath() string {
	secretsPath := os.Getenv("SECRETS_PATH")
	if len(secretsPath) == 0 {
		fmt.Println("You must set the SECRETS_PATH variable first")
		fmt.Println("ie. export SECRETS_PATH='/Users/daveo/.chef/APL/'")
		os.Exit(0)
	}
	return secretsPath
}

func getSSHKey() string {
	sshKey := os.Getenv("SSH_KEY")
	if len(sshKey) == 0 {
		fmt.Println("You must set the SSH_KEY variable first")
		fmt.Println("ie. export SSH_KEY='/Users/daveo/.ssh/APL/dodizzle.key'")
		os.Exit(0)
	}
	return sshKey
}

func getKnifePath() string {
	knifePath := os.Getenv("KNIFE_PATH")
	if len(knifePath) == 0 {
		fmt.Println("You must set the KNIFE_PATH variable first")
		fmt.Println("ie. export KNIFE_PATH='/Users/daveo/.chef/knife.rb'")
		os.Exit(0)
	}
	return knifePath
}

func runKnife(projectName string, flavorChoice string, hostname string,
	sshKey string, netID string, chefEnvironment string,
	sshKeyName string, imageName string, ipAddress string,
	secretsFile string) {
	fmt.Println("knife openstack server create -T", projectName,
		"--identity-file", sshKey,
		"-E", chefEnvironment,
		"--network-ids", netID,
		"--openstack-ssh-key-id", sshKeyName,
		"-f", flavorChoice,
		"-N", hostname,
		"--bootstrap-protocol ssh",
		"--bootstrap-version 12.17.44-1",
		"--secret-file", secretsFile,
		"-I", imageName,
		"-r role[common]",
		"--ssh-user ubuntu",
		"--sudo",
		"-G default",
		"-a", ipAddress,
		"-y")
}

func getHostname() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("hostname: ")
	hostname, _ := reader.ReadString('\n')
	return strings.TrimSpace(hostname)
}

func getChefCredentials(knifePath string) map[string]string {
	out, err := ioutil.ReadFile(knifePath)
	if err != nil {
		fmt.Println("Could not find the knife.rb file:", err)
		os.Exit(1)
	}
	dat := string(out)

	lines := strings.Split(dat, "\n")
	chefConfigs := make(map[string]string)
	for _, l := range lines {
		if strings.Contains(string(l), "openstack") {
			configs := strings.SplitAfterN(l, "=", 2)
			configKey := configs[0]
			configValue := strings.TrimSpace(strings.Replace(configs[1], "\"", "", -1))
			replaceKey := strings.NewReplacer("]", "", "knife[:", "", "_", "", "=", "")
			cleanConfigkey := strings.TrimSpace(replaceKey.Replace(configKey))
			chefConfigs[cleanConfigkey] = configValue
		} else if strings.Contains(string(l), "node_name") {
			words := strings.Fields(l)
			replaceKey := strings.NewReplacer("\"", "", "_", "")
			nodeNameKey := strings.TrimSpace(replaceKey.Replace(words[0]))
			nodeNameValue := strings.TrimSpace(replaceKey.Replace(words[1]))
			chefConfigs[nodeNameKey] = nodeNameValue
		} else if strings.Contains(string(l), "client_key") {
			words := strings.Fields(l)
			replaceKey := strings.NewReplacer("\"", "", "_", "", "#{current_dir}", "")
			clientKeyKey := strings.TrimSpace(replaceKey.Replace(words[0]))
			clientKeyValue := strings.TrimSpace(replaceKey.Replace(words[1]))
			chefConfigs[clientKeyKey] = clientKeyValue
		} else if strings.Contains(string(l), "chef_server_url") {
			words := strings.Fields(l)
			replaceKey := strings.NewReplacer("\"", "", "_", "")
			chefServerURLKey := strings.TrimSpace(replaceKey.Replace(words[0]))
			chefServerURLValue := strings.TrimSpace(replaceKey.Replace(words[1]))
			chefConfigs[chefServerURLKey] = chefServerURLValue
		}
	}

	return chefConfigs
}
