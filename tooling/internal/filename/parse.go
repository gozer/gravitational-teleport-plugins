package filename

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/coreos/go-semver/semver"
	"github.com/gravitational/trace"
)

var (
	filenamePattern *regexp.Regexp = regexp.MustCompile(`^(?P<plugin>.*)-teleport-v(?P<version>.*)-(?P<os>linux|darwin|windows)-(?P<arch>amd64|arm|aarch)-bin.tar.gz$`)
)

// Info holds information about a plugin, deduced from from its Houston-compatible
// filename.
type Info struct {
	Type    string
	Version semver.Version
	OS      string
	Arch    string
}

// Parse attempts to deduce information about a staged plugin from its (assumed
// Houston-compatible) filename, returning an error if the filename can't be
// parsed.
func Parse(filename string) (Info, error) {
	filename = filepath.Base(filename)

	matches := filenamePattern.FindStringSubmatch(filename)
	if len(matches) == 0 {
		return Info{}, trace.Errorf("filename %q does not match required pattern", filename)
	}

	version, err := semver.NewVersion(matches[2])
	if err != nil {
		return Info{}, trace.Wrap(err, "failed parsing version as semver")
	}

	return Info{
		Type:    matches[1],
		Version: *version,
		OS:      matches[3],
		Arch:    matches[4],
	}, nil
}

// Filename generates a Houston-compatible filename for the Info block, with a
// given file extension (NB: the extension is expected to include the leading
// dot).
func (info *Info) Filename(extension string) string {
	return fmt.Sprintf("%s-teleport-v%s-%s-%s-bin%s", info.Type, info.Version, info.OS, info.Arch, extension)
}
