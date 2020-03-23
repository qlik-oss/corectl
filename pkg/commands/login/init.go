package login

import (
	"bufio"
	"fmt"
	"github.com/qlik-oss/corectl/pkg/boot"
	"github.com/qlik-oss/corectl/pkg/dynconf"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"net/url"
	"os"
	"strings"
	"syscall"
)

// createInitCommand creates a command used for configuring access to
// QCS/QSEoK with an API-key
func CreateInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init <context name>",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Set up access to Qlik Sense Cloud",
		Long:  "Set up access to Qlik Sense on Cloud Services/Kubernetes by entering the domain name and the api key of the Qlik Sense instance. If no context name is supplied the domain name is used as context name",
		Run: func(cmd *cobra.Command, args []string) {

			tenant, _ := cmd.Flags().GetString("server")
			apiKey, _ := cmd.Flags().GetString("api-key")
			if len(args) == 1 {
				setupContext(tenant, apiKey, args[0])
			} else {
				setupContext(tenant, apiKey, "")
			}

			comm := boot.NewCommunicator(cmd)

			m := map[string]interface{}{}
			err := comm.RestCaller().CallStd("GET", "v1/users/me", "", nil, nil, &m)
			if err != nil {
				// TODO better error handling
				fmt.Println(err)
				fmt.Fprintln(os.Stderr, "Failed to validate context, it might be incorrect.")
				fmt.Fprintln(os.Stderr, "Perhaps the API-key supplied was faulty?")
				os.Exit(1)
			}
			name := m["name"].(string)
			fmt.Printf("Welcome %s, everything is now set up.\n", name)
		},
	}
	cmd.Flags().String("api-key", "", "API key of the tenant")
	return cmd
}

// askForConfirmation is borrowed from corectl. TODO
func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y/n]: ", s)
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))
	if response == "y" || response == "yes" {
		return true
	}
	return false
}

// setupContext sets up the configuration with tenant URL and API-key.
func setupContext(tenant, apikey, explicitContextName string) {
	if tenant == "" || apikey == "" {
		fmt.Println("Acquiring access to Qlik Cloud Services or Qlik Sense Enterprise on Kubernetes")
		fmt.Println("To complete the setup you have to have the 'developer' role and have")
		fmt.Println("API-keys enabled. If you're unsure, you can ask your tenant-admin.")
		fmt.Println()
	}
	if tenant == "" {
		fmt.Println("Specify your tenant URL, usually in the form: https://<tenant>.<region>.qlikcloud.com")
		fmt.Println("Where <tenant> is the name of the tenant and <region> is eu, us, ap, etc...")
		fmt.Print("Enter tenant url: ")
		reader := bufio.NewReader(os.Stdin)
		tenant, _ = reader.ReadString('\n')
		tenant = strings.TrimSpace(tenant)
	}
	var err error
	tenantUrl, err := parseTenantURL(tenant)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if apikey == "" {

		fmt.Printf("To generate a new API-key, go to %s/settings/api-keys\n", tenantUrl)
		fmt.Print("API-key: ")
		keyBytes, err := terminal.ReadPassword(syscall.Stdin)
		fmt.Println()
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to read API-key from input:", err)
			os.Exit(1)
		}
		apikey = strings.TrimSpace(string(keyBytes))
	}
	contextName := contextName(explicitContextName, tenantUrl)

	dynconf.CreateContext(contextName, map[string]interface{}{
		"server":  tenantUrl,
		"headers": map[string]string{"Authorization": "Bearer " + apikey},
	})
	dynconf.UseContext(contextName)
}

func parseTenantURL(rawURL string) (string, error) {
	if !strings.Contains(rawURL, "://") {
		rawURL = "https://" + rawURL
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse url: %w", err)
	}
	u.Path = ""
	return u.String(), nil
}

func contextName(explicitName, tenantUrl string) string {
	if explicitName != "" {
		return explicitName
	} else {
		res, _ := url.Parse(tenantUrl)
		return res.Hostname()
	}
}
