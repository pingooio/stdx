package email

import (
	"bufio"
	"bytes"
	_ "embed"
	"strings"
)

//go:embed email_domains_blocklist.txt
var domainsBlocklistBytes []byte

//go:embed emails_tlds_blocklist.txt
var tldBlocklistBytes []byte

var domainsBlocklist map[string]bool
var tldBlocklist []string

func IsInBlocklist(email string) bool {
	if domainsBlocklist == nil {
		mailBlocklistFileReader := bytes.NewReader(domainsBlocklistBytes)
		mailBlocklistScanner := bufio.NewScanner(mailBlocklistFileReader)
		domainsBlocklist = map[string]bool{}

		for mailBlocklistScanner.Scan() {
			line := strings.TrimSpace(mailBlocklistScanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			domainsBlocklist[line] = true
		}
	}

	return domainsBlocklist[email]
}

func IsBlockedTld(email string) bool {
	if tldBlocklist == nil {
		mailBlocklistFileReader := bytes.NewReader(tldBlocklistBytes)
		mailBlocklistScanner := bufio.NewScanner(mailBlocklistFileReader)
		tldBlocklist = make([]string, 0, 10)

		for mailBlocklistScanner.Scan() {
			line := strings.TrimSpace(mailBlocklistScanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			tldBlocklist = append(tldBlocklist, line)
		}
	}

	for _, blockedTld := range tldBlocklist {
		if strings.HasSuffix(email, blockedTld) {
			return true
		}
	}

	return false
}
