package ostack

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/floatingip"
	"github.com/rackspace/gophercloud/openstack/compute/v2/flavors"

	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/openstack/imageservice/v2/images"
	"github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/keypairs"
)

var (
	flavorTypes = map[int][]string{}
)

// ListIps get list of available floating ips
func ListIps(p *gophercloud.ProviderClient) string {
	osclient, err := openstack.NewComputeV2(p, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		fmt.Println("Error:", err)
	}
	pager := floatingip.List(osclient)
	ipsReturn := map[int][]string{}
	pager.EachPage(func(page pagination.Page) (bool, error) {
		ipsList, errs := floatingip.ExtractFloatingIPs(page)
		startingNum := 0
		for _, s := range ipsList {
			fixedIP := s.FixedIP
			// only include ip if it isn't already associatted
			if len(fixedIP) == 0 {
				startingNum++
				ipsReturn[startingNum] = []string{s.IP}
			}
		}
		if errs != nil {
			return false, nil
		}
		return true, nil
	})
	ipsTable := tablewriter.NewWriter(os.Stdout)
	ipsTable.SetHeader([]string{"", "IP"})
	// now
	var keys []int
	for k := range ipsReturn {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		ipsTable.Append([]string{strconv.Itoa(k), ipsReturn[k][0]})
	}

	ipsTable.Render()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("IP address: ")

	ipsIndex, _ := reader.ReadString('\n')

	ipsChoice, erry := strconv.Atoi(strings.TrimSpace(ipsIndex))
	if erry != nil {
		fmt.Println("Error:", err)
	}
	return (ipsReturn[ipsChoice][0])
}

// ListImages for user input
func ListImages(p *gophercloud.ProviderClient) string {
	osclient, err := openstack.NewComputeV2(p, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		fmt.Println("Error:", err)
	}
	imagesOpts := images.ListOpts{}
	pager := images.List(osclient, imagesOpts)
	imagesReturn := map[int][]string{}
	pager.EachPage(func(page pagination.Page) (bool, error) {
		imagesList, errs := images.ExtractImages(page)
		startingNum := 0
		for _, s := range imagesList {
			startingNum++

			imagesReturn[startingNum] = []string{s.Name}

		}
		if errs != nil {
			return false, nil
		}
		return true, nil
	})
	imagesTable := tablewriter.NewWriter(os.Stdout)
	imagesTable.SetHeader([]string{"", "Image"})
	// now
	var keys []int
	for k := range imagesReturn {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		imagesTable.Append([]string{strconv.Itoa(k), imagesReturn[k][0]})
	}

	imagesTable.Render()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Image: ")

	imagesIndex, _ := reader.ReadString('\n')

	imagesChoice, erry := strconv.Atoi(strings.TrimSpace(imagesIndex))
	if erry != nil {
		fmt.Println("Error:", err)
	}
	return (imagesReturn[imagesChoice][0])
}

// GetSecretsFile prompt for which secret file to use
func GetSecretsFile(secretsPath string) string {
	secretsReturn := map[int][]string{}
	startingNum := 0
	files, _ := ioutil.ReadDir(secretsPath)

	for _, f := range files {
		fname := f.Name()
		if strings.Contains(fname, "encrypted_data_bag_secret") {
			startingNum++
			secretsReturn[startingNum] = []string{fname}

		}
	}
	secretsTable := tablewriter.NewWriter(os.Stdout)
	secretsTable.SetHeader([]string{"", "SECRETS FILE"})
	// now
	var keys []int
	for k := range secretsReturn {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		secretsTable.Append([]string{strconv.Itoa(k), secretsReturn[k][0]})
	}

	secretsTable.Render()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("SECRETS_FILE: ")

	secretsIndex, _ := reader.ReadString('\n')

	secretsChoice, erry := strconv.Atoi(strings.TrimSpace(secretsIndex))
	if erry != nil {
		fmt.Println("Error:", erry)
	}
	return (secretsReturn[secretsChoice][0])

}

// ListSshkeys list available ssh keys
func ListSshkeys(p *gophercloud.ProviderClient) string {
	osclient, err := openstack.NewComputeV2(p, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		fmt.Println("Error:", err)
	}
	pager := keypairs.List(osclient)
	sshkeyReturn := map[int][]string{}
	pager.EachPage(func(page pagination.Page) (bool, error) {
		sshkeyList, errs := keypairs.ExtractKeyPairs(page)
		startingNum := 0
		for _, s := range sshkeyList {
			startingNum++
			sshkeyReturn[startingNum] = []string{s.Name}
		}
		if errs != nil {
			return false, nil
		}
		return true, nil
	})
	sshkeyTable := tablewriter.NewWriter(os.Stdout)
	sshkeyTable.SetHeader([]string{"", "Key Name"})
	var keys []int
	for k := range sshkeyReturn {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		sshkeyTable.Append([]string{strconv.Itoa(k), sshkeyReturn[k][0]})
	}
	sshkeyTable.Render()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Key Pair: ")
	sshkeyIndex, _ := reader.ReadString('\n')

	sshkeyChoice, erry := strconv.Atoi(strings.TrimSpace(sshkeyIndex))
	if erry != nil {
		fmt.Println("Error:", err)
	}
	return (sshkeyReturn[sshkeyChoice][0])
}

// GetNetworkID get id when you pass in the network name
func GetNetworkID(p *gophercloud.ProviderClient, projectName string) string {
	osclient, err := openstack.NewNetworkV2(p, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		fmt.Println("Error:", err)
	}
	netID, erry := networks.IDFromName(osclient, projectName+"_network")
	if erry != nil {
		fmt.Println("Error:", err)
	}
	return (netID)
}

// ListFlavors list flavors and get user input
func ListFlavors(p *gophercloud.ProviderClient) string {
	osclient, err := openstack.NewComputeV2(p, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		fmt.Println("Error:", err)
	}
	flavorOpts := flavors.ListOpts{}
	pager := flavors.ListDetail(osclient, flavorOpts)
	flavorReturn := map[int][]string{}
	pager.EachPage(func(page pagination.Page) (bool, error) {
		flavorList, errs := flavors.ExtractFlavors(page)
		startingNum := 0
		for _, s := range flavorList {
			startingNum++
			flavorReturn[startingNum] = []string{s.Name, strconv.Itoa(s.VCPUs), strconv.Itoa(s.Disk), strconv.Itoa(s.RAM)}
			//m["fish"] = []string{"orange", "red"}
		}
		if errs != nil {
			return false, nil
		}
		return true, nil
	})
	flavorTable := tablewriter.NewWriter(os.Stdout)
	flavorTable.SetHeader([]string{"", "Name", "Cores", "Disk", "RAM"})
	// now
	var keys []int
	for k := range flavorReturn {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		flavorTable.Append([]string{strconv.Itoa(k), flavorReturn[k][0], flavorReturn[k][1], flavorReturn[k][2], flavorReturn[k][3]})
	}

	flavorTable.Render()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Flavor: ")
	flavorIndex, _ := reader.ReadString('\n')

	flavorChoice, erry := strconv.Atoi(strings.TrimSpace(flavorIndex))
	if erry != nil {
		fmt.Println("Error:", err)
	}
	return (flavorReturn[flavorChoice][0])
}

// ListServers as stated and unused
func ListServers(p *gophercloud.ProviderClient) {

	osclient, err := openstack.NewComputeV2(p, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		fmt.Println("Error:", err)
	}
	serverOpts := servers.ListOpts{}
	pager := servers.List(osclient, serverOpts)

	pager.EachPage(func(page pagination.Page) (bool, error) {
		serverList, err := servers.ExtractServers(page)

		for _, s := range serverList {
			fmt.Println("name:", s.Name)
		}
		if err != nil {
			return false, nil
		}
		return true, nil
	})
}
