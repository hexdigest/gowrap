// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package imports

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/token"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

// Debug controls verbose logging.
var Debug = false

// LocalPrefix is a comma-separated string of import path prefixes, which, if
// set, instructs Process to sort the import paths with the given prefixes
// into another group after 3rd-party packages.
var LocalPrefix string

var goPackagesDir string
var go111ModuleEnv string

func localPrefixes() []string {
	if LocalPrefix != "" {
		return strings.Split(LocalPrefix, ",")
	}
	return nil
}

// importToGroup is a list of functions which map from an import path to
// a group number.
var importToGroup = []func(importPath string) (num int, ok bool){
	func(importPath string) (num int, ok bool) {
		for _, p := range localPrefixes() {
			if strings.HasPrefix(importPath, p) || strings.TrimSuffix(p, "/") == importPath {
				return 3, true
			}
		}
		return
	},
	func(importPath string) (num int, ok bool) {
		if strings.HasPrefix(importPath, "appengine") {
			return 2, true
		}
		return
	},
	func(importPath string) (num int, ok bool) {
		if strings.Contains(importPath, ".") {
			return 1, true
		}
		return
	},
}

func importGroup(importPath string) int {
	for _, fn := range importToGroup {
		if n, ok := fn(importPath); ok {
			return n
		}
	}
	return 0
}

// importInfo is a summary of information about one import.
type importInfo struct {
	Name  string
	Path  string // full import path (e.g. "crypto/rand")
	Alias string // import alias, if present (e.g. "crand")
}

// packageInfo is a summary of features found in a package.
type packageInfo struct {
	Globals map[string]bool     // symbol => true
	Imports map[importInfo]bool // Every import encountered, even if they conflict.
	// refs are a set of package references currently satisfied by imports.
	// first key: package path
	// second key: referenced package symbol (e.g. "Println")
	Refs map[string]map[string]bool
}

// dirPackageInfoFile gets information from other files in the package.
func dirPackageInfo(src []byte, absFilename string) (*packageInfo, error) {
	cfg := &packages.Config{
		Dir:     goPackagesDir,
		Mode:    packages.LoadSyntax,
		Tests:   strings.HasSuffix(absFilename, "_test.go"),
		Env:     appendModulesEnv(append(os.Environ(), "GOROOT="+build.Default.GOROOT, "GOPATH="+build.Default.GOPATH)),
		Overlay: map[string][]byte{absFilename: src},
	}

	pkgs, err := packages.Load(cfg, "file="+absFilename)
	if err != nil {
		return nil, fmt.Errorf("loading package for file %q: %v", absFilename, err)
	}

	info := &packageInfo{
		Globals: make(map[string]bool),
		Imports: make(map[importInfo]bool),
		Refs:    make(map[string]map[string]bool),
	}

	for _, pkg := range pkgs {
		for _, syntax := range pkg.Syntax {
			importMap := make(map[string]importInfo)
			for _, decl := range syntax.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}

				for _, spec := range genDecl.Specs {
					valueSpec, ok := spec.(*ast.ValueSpec)
					if !ok {
						continue
					}
					info.Globals[valueSpec.Names[0].Name] = true
				}

				for _, imp := range syntax.Imports {
					impInfo := importInfo{Path: strings.Trim(imp.Path.Value, `"`)}
					impInfo.Name = path.Base(impInfo.Path)
					if imp.Name != nil {
						impInfo.Alias = imp.Name.Name
						importMap[impInfo.Alias] = impInfo
					}
					importMap[impInfo.Name] = impInfo

					info.Imports[impInfo] = true
				}
			}
			visitor := collectReferences(info, importMap)
			ast.Walk(visitor, syntax)
		}
	}
	return info, nil
}

func appendModulesEnv(e []string) []string {
	if go111ModuleEnv != "" {
		return append(e, "GO111MODULE="+go111ModuleEnv)
	}
	return e
}

// collectReferences returns a visitor that collects all exported package
// references
func collectReferences(info *packageInfo, importMap map[string]importInfo) visitFn {
	var visitor visitFn
	visitor = func(node ast.Node) ast.Visitor {
		if node == nil {
			return visitor
		}
		switch v := node.(type) {
		case *ast.SelectorExpr:
			xident, ok := v.X.(*ast.Ident)
			if !ok {
				break
			}
			if xident.Obj != nil {
				// if the parser can resolve it, it's not a package ref
				break
			}
			pkgName := xident.Name
			pkg, ok := importMap[pkgName]
			if !ok {
				// Not a known package/alias.
				break
			}
			r := info.Refs[pkg.Path]
			if r == nil {
				r = make(map[string]bool)
				info.Refs[pkg.Path] = r
			}
			if ast.IsExported(v.Sel.Name) {
				r[v.Sel.Name] = true
			}
		}
		return visitor
	}
	return visitor
}

func fixImports(src []byte, fset *token.FileSet, f *ast.File, filename string) (added []string, err error) {
	// refs are a set of possible package references currently unsatisfied by imports.
	// first key: either base package (e.g. "fmt") or renamed package
	// second key: referenced package symbol (e.g. "Println")
	refs := make(map[string]map[string]bool)

	// decls are the current package imports. key is base package or renamed package.
	decls := make(map[string]*ast.ImportSpec)

	abs, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}
	srcDir := filepath.Dir(abs)
	if Debug {
		log.Printf("fixImports(filename=%q), abs=%q, srcDir=%q ...", filename, abs, srcDir)
	}

	packageInfo, err := dirPackageInfo(src, abs)
	if err != nil {
		return nil, err
	}

	// collect potential uses of packages.
	var visitor visitFn
	visitor = visitFn(func(node ast.Node) ast.Visitor {
		if node == nil {
			return visitor
		}
		switch v := node.(type) {
		case *ast.ImportSpec:
			if v.Name != nil {
				decls[v.Name.Name] = v
				break
			}
			ipath := strings.Trim(v.Path.Value, `"`)
			if ipath == "C" {
				break
			}
			local := path.Base(ipath)
			decls[local] = v
		case *ast.SelectorExpr:
			xident, ok := v.X.(*ast.Ident)
			if !ok {
				break
			}
			if xident.Obj != nil {
				// if the parser can resolve it, it's not a package ref
				break
			}
			pkgName := xident.Name
			if refs[pkgName] == nil {
				refs[pkgName] = make(map[string]bool)
			}
			if decls[pkgName] == nil && (packageInfo == nil || !packageInfo.Globals[pkgName]) {
				refs[pkgName][v.Sel.Name] = true
			}
		}
		return visitor
	})
	ast.Walk(visitor, f)

	// Nil out any unused ImportSpecs, to be removed in following passes
	unusedImport := map[string]string{}
	for pkg, is := range decls {
		if refs[pkg] == nil && pkg != "_" && pkg != "." {
			name := ""
			if is.Name != nil {
				name = is.Name.Name
			}
			unusedImport[strings.Trim(is.Path.Value, `"`)] = name
		}
	}
	for ipath, name := range unusedImport {
		if ipath == "C" {
			// Don't remove cgo stuff.
			continue
		}
		astutil.DeleteNamedImport(fset, f, name, ipath)
	}

	for pkgName, symbols := range refs {
		if len(symbols) == 0 {
			// skip over packages already imported
			delete(refs, pkgName)
		}
	}

	// Fast path, all references already imported.
	if len(refs) == 0 {
		return nil, nil
	}

	var loadQueries []string
	for pkgName := range refs {
		loadQueries = append(loadQueries, "name="+pkgName)
	}
	cfg := &packages.Config{
		Dir:  goPackagesDir,
		Mode: packages.LoadTypes,
		Env:  appendModulesEnv(append(os.Environ(), "GOROOT="+build.Default.GOROOT, "GOPATH="+build.Default.GOPATH)),
	}
	pkgs, err := packages.Load(cfg, loadQueries...)
	pkgsByName := make(map[string][]*packages.Package)
	for _, pkg := range pkgs {
		pkgsByName[pkg.Name] = append(pkgsByName[pkg.Name], pkg)
	}

outer:
	for pkgName, symbols := range refs {
		for sibling := range packageInfo.Imports {
			if sibling.Name != pkgName && sibling.Alias != pkgName {
				continue
			}

			refs := packageInfo.Refs[sibling.Path]
			allFound := true
			for symbol := range symbols {
				if !refs[symbol] {
					allFound = false
					break
				}
			}
			if allFound {
				if sibling.Name == pkgName {
					astutil.AddImport(fset, f, sibling.Path)
				} else {
					astutil.AddNamedImport(fset, f, sibling.Alias, sibling.Path)
				}
				continue outer
			}
		}

		ipath, rename, err := findImport(pkgName, pkgsByName[pkgName], symbols, filename)
		if err != nil {
			return nil, err
		}

		if ipath == "" {
			continue // No matching package.
		}

		if rename {
			astutil.AddNamedImport(fset, f, pkgName, ipath)
		} else {
			astutil.AddImport(fset, f, ipath)
		}
	}
	return added, nil
}

var stdImportPackage = map[string]string{} // "net/http" => "http"

func init() {
	// Nothing in the standard library has a package name not
	// matching its import base name.
	for _, pkg := range stdlib {
		if _, ok := stdImportPackage[pkg]; !ok {
			stdImportPackage[pkg] = path.Base(pkg)
		}
	}
}

type pkg struct {
	*packages.Package
	dir             string // absolute file path to pkg directory ("/usr/lib/go/src/net/http")
	importPath      string // full pkg import path ("net/http", "foo/bar/vendor/a/b")
	importPathShort string // vendorless import path ("net/http", "a/b")
}

type pkgDistance struct {
	pkg      *pkg
	distance int // relative distance to target
}

// byDistanceOrImportPathShortLength sorts by relative distance breaking ties
// on the short import path length and then the import string itself.
type byDistanceOrImportPathShortLength []pkgDistance

func (s byDistanceOrImportPathShortLength) Len() int { return len(s) }
func (s byDistanceOrImportPathShortLength) Less(i, j int) bool {
	di, dj := s[i].distance, s[j].distance
	if di == -1 {
		return false
	}
	if dj == -1 {
		return true
	}
	if di != dj {
		return di < dj
	}

	vi, vj := s[i].pkg.importPathShort, s[j].pkg.importPathShort
	if len(vi) != len(vj) {
		return len(vi) < len(vj)
	}
	return vi < vj
}
func (s byDistanceOrImportPathShortLength) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func distance(basepath, targetpath string) int {
	p, err := filepath.Rel(basepath, targetpath)
	if err != nil {
		return -1
	}
	if p == "." {
		return 0
	}
	return strings.Count(p, string(filepath.Separator)) + 1
}

// VendorlessPath returns the devendorized version of the import path ipath.
// For example, VendorlessPath("foo/bar/vendor/a/b") returns "a/b".
func VendorlessPath(ipath string) string {
	// Devendorize for use in import statement.
	if i := strings.LastIndex(ipath, "/vendor/"); i >= 0 {
		return ipath[i+len("/vendor/"):]
	}
	if strings.HasPrefix(ipath, "vendor/") {
		return ipath[len("vendor/"):]
	}
	return ipath
}

// findImport searches for a package with the given symbols.
// If no package is found, findImport returns ("", false, nil)
//
// This is declared as a variable rather than a function so goimports
// can be easily extended by adding a file with an init function.
//
// The rename value tells goimports whether to use the package name as
// a local qualifier in an import. For example, if findImports("pkg",
// "X") returns ("foo/bar", rename=true), then goimports adds the
// import line:
// 	import pkg "foo/bar"
// to satisfy uses of pkg.X in the file.
func findImport(pkgName string, candidatePkgs []*packages.Package, symbols map[string]bool, filename string) (foundPkg string, rename bool, err error) {
	pkgDir, err := filepath.Abs(filename)
	if err != nil {
		return "", false, err
	}
	pkgDir = filepath.Dir(pkgDir)

	// Fast path for the standard library.
	// In the common case we hopefully never have to scan the GOPATH, which can
	// be slow with moving disks.
	if pkg, ok := findImportStdlib(pkgName, symbols); ok {
		return pkg, false, nil
	}
	if pkgName == "rand" && symbols["Read"] {
		// Special-case rand.Read.
		//
		// If findImportStdlib didn't find it above, don't go
		// searching for it, lest it find and pick math/rand
		// in GOROOT (new as of Go 1.6)
		//
		// crypto/rand is the safer choice.
		return "", false, nil
	}

	// Find candidate packages, looking only at their directory names first.
	var candidates []pkgDistance
	for _, packagePkg := range candidatePkgs {
		pkg := &pkg{
			Package:         packagePkg,
			dir:             filepath.Dir(packagePkg.CompiledGoFiles[0]),
			importPath:      packagePkg.PkgPath,
			importPathShort: VendorlessPath(packagePkg.PkgPath),
		}

		if pkgIsCandidate(filename, pkgName, pkg) {
			candidates = append(candidates, pkgDistance{
				pkg:      pkg,
				distance: distance(pkgDir, pkg.dir),
			})
		}
	}

	// Sort the candidates by their import package length,
	// assuming that shorter package names are better than long
	// ones.  Note that this sorts by the de-vendored name, so
	// there's no "penalty" for vendoring.
	sort.Sort(byDistanceOrImportPathShortLength(candidates))
	if Debug {
		for i, c := range candidates {
			log.Printf("%s candidate %d/%d: %v in %v", pkgName, i+1, len(candidates), c.pkg.importPathShort, c.pkg.dir)
		}
	}

	var selected *pkg
outer:
	for _, candidate := range candidates {
		exports := make(map[string]bool)
		for _, name := range candidate.pkg.Types.Scope().Names() {
			if ast.IsExported(name) {
				exports[name] = true
			}
		}
		if Debug {
			exportList := make([]string, 0, len(exports))
			for k := range exports {
				exportList = append(exportList, k)
			}
			sort.Strings(exportList)
			log.Printf("loaded exports in package %v: %v", candidate.pkg.Name, strings.Join(exportList, ", "))
		}

		for sym := range symbols {
			if !exports[sym] {
				continue outer
			}
		}
		selected = candidate.pkg
		break
	}

	if selected == nil {
		return "", false, nil
	}
	// If the package name in the source doesn't match the import path,
	// return true so the rewriter adds a name (import foo "github.com/bar/go-foo")
	// Module-style version suffixes are allowed.
	lastSeg := path.Base(selected.importPath)
	if isVersionSuffix(lastSeg) {
		lastSeg = path.Base(path.Dir(selected.importPath))
	}
	needsRename := lastSeg != pkgName
	return selected.importPathShort, needsRename, nil
}

// isVersionSuffix reports whether the path segment seg is a semantic import
// versioning style major version suffix.
func isVersionSuffix(seg string) bool {
	if seg == "" {
		return false
	}
	if seg[0] != 'v' {
		return false
	}
	if _, err := strconv.Atoi(seg[1:]); err != nil {
		return false
	}
	return true
}

// pkgIsCandidate reports whether pkg is a candidate for satisfying the
// finding which package pkgIdent in the file named by filename is trying
// to refer to.
//
// This check is purely lexical and is meant to be as fast as possible
// because it's run over all $GOPATH directories to filter out poor
// candidates in order to limit the CPU and I/O later parsing the
// exports in candidate packages.
//
// filename is the file being formatted.
// pkgIdent is the package being searched for, like "client" (if
// searching for "client.New")
func pkgIsCandidate(filename, pkgIdent string, pkg *pkg) bool {
	// Check "internal" and "vendor" visibility:
	if !canUse(filename, pkg.dir) {
		return false
	}

	// Speed optimization to minimize disk I/O:
	// the last two components on disk must contain the
	// package name somewhere.
	//
	// This permits mismatch naming like directory
	// "go-foo" being package "foo", or "pkg.v3" being "pkg",
	// or directory "google.golang.org/api/cloudbilling/v1"
	// being package "cloudbilling", but doesn't
	// permit a directory "foo" to be package
	// "bar", which is strongly discouraged
	// anyway. There's no reason goimports needs
	// to be slow just to accommodate that.
	lastTwo := lastTwoComponents(pkg.importPathShort)
	if strings.Contains(lastTwo, pkgIdent) {
		return true
	}
	if hasHyphenOrUpperASCII(lastTwo) && !hasHyphenOrUpperASCII(pkgIdent) {
		lastTwo = lowerASCIIAndRemoveHyphen(lastTwo)
		if strings.Contains(lastTwo, pkgIdent) {
			return true
		}
	}

	return false
}

func hasHyphenOrUpperASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		b := s[i]
		if b == '-' || ('A' <= b && b <= 'Z') {
			return true
		}
	}
	return false
}

func lowerASCIIAndRemoveHyphen(s string) (ret string) {
	buf := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		b := s[i]
		switch {
		case b == '-':
			continue
		case 'A' <= b && b <= 'Z':
			buf = append(buf, b+('a'-'A'))
		default:
			buf = append(buf, b)
		}
	}
	return string(buf)
}

// canUse reports whether the package in dir is usable from filename,
// respecting the Go "internal" and "vendor" visibility rules.
func canUse(filename, dir string) bool {
	// Fast path check, before any allocations. If it doesn't contain vendor
	// or internal, it's not tricky:
	// Note that this can false-negative on directories like "notinternal",
	// but we check it correctly below. This is just a fast path.
	if !strings.Contains(dir, "vendor") && !strings.Contains(dir, "internal") {
		return true
	}

	dirSlash := filepath.ToSlash(dir)
	if !strings.Contains(dirSlash, "/vendor/") && !strings.Contains(dirSlash, "/internal/") && !strings.HasSuffix(dirSlash, "/internal") {
		return true
	}
	// Vendor or internal directory only visible from children of parent.
	// That means the path from the current directory to the target directory
	// can contain ../vendor or ../internal but not ../foo/vendor or ../foo/internal
	// or bar/vendor or bar/internal.
	// After stripping all the leading ../, the only okay place to see vendor or internal
	// is at the very beginning of the path.
	absfile, err := filepath.Abs(filename)
	if err != nil {
		return false
	}
	absdir, err := filepath.Abs(dir)
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(absfile, absdir)
	if err != nil {
		return false
	}
	relSlash := filepath.ToSlash(rel)
	if i := strings.LastIndex(relSlash, "../"); i >= 0 {
		relSlash = relSlash[i+len("../"):]
	}
	return !strings.Contains(relSlash, "/vendor/") && !strings.Contains(relSlash, "/internal/") && !strings.HasSuffix(relSlash, "/internal")
}

// lastTwoComponents returns at most the last two path components
// of v, using either / or \ as the path separator.
func lastTwoComponents(v string) string {
	nslash := 0
	for i := len(v) - 1; i >= 0; i-- {
		if v[i] == '/' || v[i] == '\\' {
			nslash++
			if nslash == 2 {
				return v[i:]
			}
		}
	}
	return v
}

type visitFn func(node ast.Node) ast.Visitor

func (fn visitFn) Visit(node ast.Node) ast.Visitor {
	return fn(node)
}

func findImportStdlib(shortPkg string, symbols map[string]bool) (importPath string, ok bool) {
	for symbol := range symbols {
		key := shortPkg + "." + symbol
		path := stdlib[key]
		if path == "" {
			if key == "rand.Read" {
				continue
			}
			return "", false
		}
		if importPath != "" && importPath != path {
			// Ambiguous. Symbols pointed to different things.
			return "", false
		}
		importPath = path
	}
	if importPath == "" && shortPkg == "rand" && symbols["Read"] {
		return "crypto/rand", true
	}
	return importPath, importPath != ""
}
