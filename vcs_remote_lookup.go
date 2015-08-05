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
	{
		host:    "go.googlesource.com",
		vcs:     GitType,
		pattern: `^(go\.googlesource\.com/[A-Za-z0-9_.\-]+/?)$`,
	},
	// TODO: Once Google Code becomes fully deprecated this can be removed.
	{
		host:     "code.google.com",
		addCheck: checkGoogle,
		pattern:  `^(code\.google\.com/[pr]/(?P<project>[a-z0-9\-]+)(\.(?P<repo>[a-z0-9\-]+))?)(/[A-Za-z0-9_.\-]+)*$`,
	},
	// Alternative Google setup. This is the previous structure but it still works... until Google Code goes away.
	{
		addCheck: checkUrl,
		pattern:  `^([a-z0-9_\-.]+)\.googlecode\.com/(?P<type>git|hg|svn)(/.*)?$`,
	},
	// If none of the previous detect the type they will fall to this looking for the type in a generic sense
	// by the extension to the path.
	{
		addCheck: checkUrl,
		pattern:  `\.(?P<type>git|hg|svn|bzr)$`,
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
		if v.host != "" && v.host != u.Host {
			continue
		}

		// Make sure the pattern matches for an actual repo location. For example,
		// we should fail if the VCS listed is github.com/masterminds as that's
		// not actually a repo.
		uCheck := u.Host + u.Path
		m := v.regex.FindStringSubmatch(uCheck)
		if m == nil {
			if v.host != "" {
				return "", ErrCannotDetectVCS
			}

			continue
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

// Google supports Git, Hg, and Svn. The SVN style is only
// supported through their legacy setup at <project>.googlecode.com.
// I wonder if anyone is actually using SVN support.
func checkGoogle(i map[string]string) (VcsType, error) {

	// To figure out which of the VCS types is used in Google Code you need
	// to parse a web page and find it. Ugh. I mean... ugh.
	var hack = regexp.MustCompile(`id="checkoutcmd">(hg|git|svn)`)

	d, err := get(expand(i, "https://code.google.com/p/{project}/source/checkout?repo={repo}"))
	if err != nil {
		return "", err
	}

	if m := hack.FindSubmatch(d); m != nil {
		if vcs := string(m[1]); vcs != "" {
			if vcs == "svn" {
				// While Google supports SVN it can only be used with the legacy
				// urls of <project>.googlecode.com. I considered creating a new
				// error for this problem but Google Code is going away and there
				// is support for the legacy structure.
				return "", ErrCannotDetectVCS
			}

			return VcsType(vcs), nil
		}
	}

	return "", ErrCannotDetectVCS
}

// Expect a type key on i with the exact type detected from the regex.
func checkUrl(i map[string]string) (VcsType, error) {
	return VcsType(i["type"]), nil
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
