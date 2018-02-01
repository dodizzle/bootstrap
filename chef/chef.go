package chef

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/dodizzle/chef"
	"github.com/olekukonko/tablewriter"
)

// GetEnvironments get list of environments froh the chef server
func GetEnvironments(nodeName string, clientKey string, chefServerURL string) string {
	key, err := ioutil.ReadFile(clientKey)
	if err != nil {
		fmt.Println("Couldn't read key:", err)
		os.Exit(1)
	}

	// build a client
	client, err := chef.NewClient(&chef.Config{
		Name:    nodeName,
		Key:     string(key),
		BaseURL: chefServerURL,
	})
	if err != nil {
		fmt.Println("Issue setting up client:", err)
		os.Exit(1)
	}

	environments, err := client.Environments.List()

	if err != nil {
		fmt.Println("Environments.List returned error:")
	}
	environMap := map[int][]string{}
	startingNum := 0
	var sorted []string
	for k, v := range environments {
		_ = v
		name := []string{k}[0]
		sorted = append(sorted, name)
	}
	sort.Strings(sorted)
	for k, v := range sorted {
		startingNum++
		_ = k
		environMap[startingNum] = []string{v}
	}

	environTable := tablewriter.NewWriter(os.Stdout)
	environTable.SetHeader([]string{"", "environment"})
	var keys []int
	for k := range environMap {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		environTable.Append([]string{strconv.Itoa(k), environMap[k][0]})
	}

	environTable.Render()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Environment: ")
	environIndex, _ := reader.ReadString('\n')

	environChoice, erry := strconv.Atoi(strings.TrimSpace(environIndex))
	if erry != nil {
		fmt.Println("Error:", err)
	}
	return (environMap[environChoice][0])
}
