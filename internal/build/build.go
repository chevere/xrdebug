package build

import (
	"embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"regexp"
	"strings"
)

const (
	nonceLength = 12
	tagLength   = 16
)

// assetType represents supported asset types for replacement
type assetType struct {
	extension string
	mimeType  string
}

var (
	svgAsset  = assetType{"svg", "image/svg+xml"}
	pngAsset  = assetType{"png", "image/png"}
	woffAsset = assetType{"woff", "font/woff"}
)

// Build represents a compiled xrDebug interface with all its assets embedded.
type Build struct {
	// path is the base path for asset resolution
	path string
	// version represents the xrDebug version
	version string
	// sessionName is the name of the debugging session
	sessionName string
	// editor specifies the preferred editor for file opening
	editor string
	// isEncryptionEnabled determines if e2e encryption is active
	isEncryptionEnabled bool
	// isSignVerificationEnabled determines if signature verification is active
	isSignVerificationEnabled bool
	// filesystem contains the embedded assets
	filesystem embed.FS
	// content holds the processed template content
	content string
}

// Replacements defines the template variables for the xrDebug interface.
type Replacements struct {
	// Version is the xrDebug version
	Version string
	// IsEncryptionEnabled indicates if encryption is active
	IsEncryptionEnabled bool
	// NonceLength is the length of the encryption nonce
	NonceLength int
	// TagLength is the length of the encryption tag
	TagLength int
	// SessionName is the name of the debugging session
	SessionName string
	// Editor is the preferred editor for file opening
	Editor string
	// Security describes the active security features
	Security string
}

// New creates a new Build instance with the provided configuration and processes all embedded assets.
func New(source []byte, filesystem embed.FS, path, version, sessionName, editor string,
	isEncryptionEnabled, isSignVerificationEnabled bool) (*Build, error) {
	b := &Build{
		path:                      path,
		version:                   version,
		sessionName:               sessionName,
		editor:                    editor,
		isEncryptionEnabled:       isEncryptionEnabled,
		isSignVerificationEnabled: isSignVerificationEnabled,
		filesystem:                filesystem,
		content:                   string(source),
	}
	if err := b.processTemplate(); err != nil {
		return nil, fmt.Errorf("processing template: %w", err)
	}
	if err := b.embedAssets(); err != nil {
		return nil, fmt.Errorf("embedding assets: %w", err)
	}
	return b, nil
}

func (b *Build) processTemplate() error {
	t, err := template.New("xrdebug").Parse(b.content)
	if err != nil {
		return err
	}
	replacements := Replacements{
		Version:             b.version,
		IsEncryptionEnabled: b.isEncryptionEnabled,
		NonceLength:         nonceLength,
		TagLength:           tagLength,
		SessionName:         b.sessionName,
		Editor:              b.editor,
		Security:            b.security(),
	}
	w := &strings.Builder{}
	if err := t.Execute(w, replacements); err != nil {
		return err
	}
	b.content = w.String()
	return nil
}

func (b *Build) embedAssets() error {
	for _, asset := range []assetType{svgAsset, pngAsset} {
		if err := b.replaceIcons(asset); err != nil {
			return err
		}
	}
	if err := b.replaceStyles(); err != nil {
		return err
	}
	if err := b.replaceFont(woffAsset); err != nil {
		return err
	}
	return b.replaceScripts()
}

// String returns the processed content as a string.
func (b *Build) String() string {
	return b.content
}

// Bytes returns the processed content as a byte slice.
func (b *Build) Bytes() []byte {
	return []byte(b.content)
}

// replaceStyles embeds CSS files inline in the HTML.
func (b *Build) replaceStyles() error {
	return b.replaceMatches(`<link rel="stylesheet".*href="(.*)".*>`, func(path string) string {
		return fmt.Sprintf("<style media=\"all\">\n%s\n</style>", path)
	})
}

// replaceScripts embeds JavaScript files inline in the HTML.
func (b *Build) replaceScripts() error {
	return b.replaceMatches(`<script .*src="(.*)".*></script>`, func(path string) string {
		return fmt.Sprintf("<script>%s</script>", path)
	})
}

// replaceIcons embeds icon files as base64-encoded data URLs.
func (b *Build) replaceIcons(asset assetType) error {
	return b.replaceMatches(fmt.Sprintf(`="(icon\.%s)"`, asset.extension), func(content string) string {
		return fmt.Sprintf(`="data:%s;base64,%s"`, asset.mimeType, base64.StdEncoding.EncodeToString([]byte(content)))
	})
}

// replaceMatches is a helper function to handle common replacement patterns
func (b *Build) replaceMatches(pattern string, formatFn func(string) string) error {
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(b.content, -1)
	for _, match := range matches {
		fileContent, err := b.filesystem.ReadFile(b.path + match[1])
		if err != nil {
			return fmt.Errorf("reading file %s: %w", match[1], err)
		}
		b.replace(match[0], formatFn(string(fileContent)))
	}
	return nil
}

// replaceFont embeds font files as base64-encoded data URLs.
func (b *Build) replaceFont(asset assetType) error {
	return b.replaceMatches(fmt.Sprintf(`url\(['"]?(.*\.%s)['"]?\)`, asset.extension), func(content string) string {
		return fmt.Sprintf(`url(data:%s;base64,%s)`, asset.mimeType, base64.StdEncoding.EncodeToString([]byte(content)))
	})
}

// replace performs a string replacement in the content.
func (b *Build) replace(search, replace string) {
	b.content = strings.ReplaceAll(b.content, search, replace)
}

// security returns a string describing the active security features.
func (b *Build) security() string {
	switch {
	case b.isEncryptionEnabled && b.isSignVerificationEnabled:
		return "End-to-end encrypted and sign verified"
	case b.isEncryptionEnabled:
		return "End-to-end encrypted"
	case b.isSignVerificationEnabled:
		return "Sign verified"
	default:
		return ""
	}
}
