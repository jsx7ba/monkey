package token

import (
	"fmt"
)

type LineInfo struct {
	FileIndex SourceHandle
	Line      uint16
	Char      uint16
}

func (li LineInfo) String() string {
	return fmt.Sprintf("%s: line %d, char %d", registry.sourceFiles[li.FileIndex], li.Line, li.Char)
}

type SourceHandle uint32

type SourceRegistry struct {
	index       int
	sourceFiles []string
}

var registry *SourceRegistry

func init() {
	registry = &SourceRegistry{0, make([]string, 0)}
}

func (sr *SourceRegistry) addSource(fileName string) SourceHandle {
	sr.sourceFiles = append(sr.sourceFiles, fileName)
	handle := SourceHandle(sr.index)
	sr.index++
	return handle
}

func (sr *SourceRegistry) reset() {
	sr.index = 0
	sr.sourceFiles = sr.sourceFiles[:0]
}

// ResetForTesting Only call this if multiple tests need to use the registry in the same package.
func ResetForTesting() {
	registry.reset()
}

func AddSource(fileName string) SourceHandle {
	return registry.addSource(fileName)
}

func (h SourceHandle) LineInfo(line, char uint16) LineInfo {
	return LineInfo{h, line, char}
}
