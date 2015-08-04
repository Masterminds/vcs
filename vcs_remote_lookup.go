package vcs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type vcsInfo struct {
	host     string
	pattern  string
	vcs      VcsType
	addCheck func(m map[string]string) (VcsType, error)
	regex    *regexp.Regexp
}

var vcsList = []*vcsInfo{
	{
		host:    "github.com",
		vcs:     GitType,
		pattern: `^(github\.com/[A-Za-z0-9_.\-]+/[A-Za-z0-9_.\-]+)(/[A-Za-z0-9_.\-]+)*$`,
	},
	{
		host:     "bitbucket.org",
		pattern:  `^(bitbucket\.org/(?P<name>[A-Za-z0-9_.\-]+/[A-Za-z0-9_.\-]+))(/[A-Za-z0-9_.\-]+)*$`,
		addCheck: checkBitbucket,
	},
	{
		host:    "launchpad.net",
		pattern: `^(launchpad\.net/(([A-Za-z0-9_.\-]+)(/[A-Za-z0-9_.\-]+)?|~[A-Za-z0-9_.\-]+/(\+junk|[A-Za-z0-9_.\-]+)/[A-Za-z0-9_.\-]+))(/[A-Za-z0-9_.\-]+)*$`,
		vcs:     BzrType,
	},
	{
		host:    "git.launchpad.net",
		vcs:     GitType,
		pattern: `^(git\.launchpad\.net/(([A-Za-z0-9_.\-]+)|~[A-Za-z0-9_.\-]+/(\+git|[A-Za-z0-9_.\-]+)/[A-Za-z0-9_.\-]+))$`,
	},
}

func init() {
	// Precompile the regular expressions used to check VCS locations.
	for _, v := range vcsList {
		v.regex = regexp.MustCompile(v.pattern)
	}
}

// From a remote vcs url attempt to detect the VCS.
func detectVcsFromUrl(vcsUrl string) (VcsType, error) {
	u, err := url.Parse(vcsUrl)
	if err != nil {
		return "", err
	}

	// If there is no host found we cannot detect the VCS from
	// the url. Note, URIs beginning with git@github using the ssh
	// syntax fail this check.
	if u.Host == "" {
		return "", ErrCannotDetectVCS
	}

	// Try to detect from known hosts, such as Github
	for _, v := range vcsList {
		if v.host != u.Host {
			continue
		}

		// Make sure the pattern matches for an actual repo location. For example,
		// we should fail if the VCS listed is github.com/masterminds as that's
		// not actually a repo.
		uCheck := u.Host + u.Path
		m := v.regex.FindStringSubmatch(uCheck)
		if m == nil {
			return "", ErrCannotDetectVCS
		}

		// If we are here the host matches. If the host has a singular
		// VCS type, such as Github, we can return the type right away.
		if v.vcs != "" {
			return v.vcs, nil
		}

		// Run additional checks to determine try and determine the repo
		// for the matched service.
		info := make(map[string]string)
		for i, name := range v.regex.SubexpNames() {
			if name != "" {
				info[name] = m[i]
			}
		}
		t, err := v.addCheck(info)
		if err != nil {
			return "", ErrCannotDetectVCS
		}

		return t, nil

	}

	// If the url to the endpoint is not for one of the known services attempt
	// to figure it out.

	// Unable to determine the vcs from the url.
	return "", ErrCannotDetectVCS
}

// Bitbucket provides an API for checking the VCS.
func checkBitbucket(i map[string]string) (VcsType, error) {

	// The part of the response we care about.
	var response struct {
		SCM VcsType `json:"scm"`
	}

	u := expand(i, "https://api.bitbucket.org/1.0/repositories/{name}")
	data, err := get(u)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return "", fmt.Errorf("Decoding error %s: %v", u, err)
	}

	return response.SCM, nil

}

func get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s: %s", url, resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", url, err)
	}
	return b, nil
}

func expand(match map[string]string, s string) string {
	for k, v := range match {
		s = strings.Replace(s, "{"+k+"}", v, -1)
	}
	return s
}
