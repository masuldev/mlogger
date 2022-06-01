package mlogger

import (
	"fmt"
	"github.com/lestrrat-go/strftime"
	"github.com/masuldev/mlogger/internal/util"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

func (c clockFn) Now() time.Time {
	return c()
}

func New(p string, options ...Option) (*RotateLog, error) {
	globalPattern := p
	for _, re := range patternConversionRegexps {
		globalPattern = re.ReplaceAllString(globalPattern, "*")
	}

	pattern, err := strftime.New(p)
	if err != nil {
		return nil, errors.Wrap(err, "invalid strftime pattern")
	}

	var clock Clock = Local
	rotationTime := 24 * time.Hour
	var rotationSize int64
	var rotationCount uint
	var linkName string
	var maxAge time.Duration
	var handler Handler
	var forceNewFile bool

	for _, o := range options {
		switch o.Name() {
		case optionClock:
			clock = o.Value().(Clock)
		case optionLinkName:
			linkName = o.Value().(string)
		case optionMaxAge:
			maxAge = o.Value().(time.Duration)
			if maxAge < 0 {
				maxAge = 0
			}
		case optionRotationTime:
			rotationTime = o.Value().(time.Duration)
			if rotationTime < 0 {
				rotationTime = 0
			}
		case optionRotationSize:
			rotationSize = o.Value().(int64)
			if rotationSize < 0 {
				rotationSize = 0
			}
		case optionRotationCount:
			rotationCount = o.Value().(uint)
		case optionHandler:
			handler = o.Value().(Handler)
		case optionForceNewFile:
			forceNewFile = true
		}
	}

	if maxAge > 0 && rotationCount > 0 {
		return nil, errors.New("MaxAge and RotationCount cannot be both set")
	}

	if maxAge == 0 && rotationCount == 0 {
		// default 7days
		maxAge = 7 * 24 * time.Hour
	}

	return &RotateLog{
		clock:         clock,
		eventHandler:  handler,
		globalPattern: globalPattern,
		linkName:      linkName,
		maxAge:        maxAge,
		rotationCount: rotationCount,
		rotationSize:  rotationSize,
		rotationTime:  rotationTime,
		forceNewFile:  forceNewFile,
		pattern:       pattern,
	}, nil
}

// for interface io.Writer
func (rl *RotateLog) Write(p []byte) (n int, err error) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	out, err := rl.get
}

func (rl *RotateLog) getWriterNolock(bailOnRotateFail, useGenerationalNames bool) (io.Writer, error) {
	generation := rl.generation
	previousFn := rl.curFn

	baseFn := util.GenerateFn(rl.pattern, rl.clock, rl.rotationTime)
	filename := baseFn
	var forceNewFile bool

	// check file size
	fi, err := os.Stat(rl.curFn)
	sizeRotation := false
	if err == nil && rl.rotationSize > 0 && rl.rotationSize <= fi.Size() {
		forceNewFile = true
		sizeRotation = true
	}

	// check is first file
	if baseFn != rl.curBaseFn {
		generation = 0
		if rl.forceNewFile {
			forceNewFile = true
		}
	} else {
		if !useGenerationalNames && !sizeRotation {
			return rl.outFh, nil
		}
		forceNewFile = true
		generation++
	}

	if forceNewFile {
		var name string
		for {
			if generation == 0 {
				name = filename
			} else {
				name = fmt.Sprintf("%s.%d", filename, generation)
			}
			if _, err := os.Stat(name); err != nil {
				filename = name
				break
			}
			generation++
		}
	}

	fh, err := util.CreateFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create a new file %v", filename)
	}

	if err := rl.rotateNolock(filename); err != nil {
		err = errors.Wrap(err, "failed to rotate")
		if bailOnRotateFail {
			if fh != nil {
				fh.Close()
			}
			return nil, err
		}
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}

	rl.outFh.Close()
	rl.outFh = fh
	rl.curBaseFn = baseFn
	rl.curFn = filename
	rl.generation = generation

	if h := rl.eventHandler; h != nil {
		go h.Handle(&FileRotatedEvent{
			prev:    previousFn,
			current: filename,
		})
	}

	return fh, nil
}

func (rl *RotateLog) rotateNolock(filename string) error {
	lockfn := filename + "_lock"
	fh, err := os.OpenFile(lockfn, os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}

	// remove lock file
	var gurad cleanupGuard
	gurad.fn = func() {
		fh.Close()
		os.Remove(lockfn)
	}
	defer gurad.Run()

	if rl.linkName != "" {
		tmpLinkName := filename + "_symlink"
		linkDest := filename
		linkDir := filepath.Dir(rl.linkName)

		baseDir := filepath.Dir(filename)
		if strings.Contains(rl.linkName, baseDir) {
			tmp, err := filepath.Rel(linkDir, filename)
			if err != nil {
				return errors.Wrapf(err, "failed to evaluate relative path from %#v to %#v", baseDir, rl.linkName)
			}

			linkDest = tmp
		}

		if err := os.Symlink(linkDest, tmpLinkName); err != nil {
			return errors.Wrap(err, "failed to create new symlink")
		}

		_, err := os.Stat(linkDir)
		if err != nil {
			if err := os.MkdirAll(linkDir, 0755); err != nil {
				return errors.Wrapf(err, "failed to create directory %s", linkDir)
			}
		}

		if err := os.Rename(tmpLinkName, rl.linkName); err != nil {
			return errors.Wrap(err, "failed to rename new symlink")
		}
	}

	if rl.maxAge <= 0 && rl.rotationCount <= 0 {
		return errors.New("panic: maxAge and rotationCount are both set")
	}

	matches, err := filepath.Glob(rl.globalPattern)
	if err != nil {
		return err
	}

	cutoff := rl.clock.Now().Add(-1 * rl.maxAge)

	toUnlink := make([]string, 0, len(matches))
	for _, path := range matches {
		if strings.HasSuffix(path, "_lock") || strings.HasSuffix(path, "_symlink") {
			continue
		}

		fi, err := os.Stat(path)
		if err != nil {
			continue
		}

		fl, err := os.Lstat(path)
		if err != nil {
			continue
		}

		if rl.maxAge > 0 && fl.ModTime().After(cutoff) {
			continue
		}

		if rl.rotationCount > 0 && fl.Mode()&os.ModeSymlink == os.ModeSymlink {
			continue
		}
		toUnlink = append(toUnlink, path)
	}

	if rl.rotationCount > 0 {
		if rl.rotationCount >= uint(len(toUnlink)) {
			return nil
		}

		toUnlink = toUnlink[:len(toUnlink)-int(rl.rotationCount)]
	}

	if len(toUnlink) <= 0 {
		return nil
	}

	gurad.Enable()
	go func() {
		for _, path := range toUnlink {
			os.Remove(path)
		}
	}()

	return nil
}

var patternConversionRegexps = []*regexp.Regexp{
	regexp.MustCompile(`%[%+A-Za-z]`),
	regexp.MustCompile(`\*+`),
}

func (rl *RotateLog) Rotate() error {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	_, err := rl.getWriterNolock(true, true)

	return err
}

func (rl *RotateLog) Close() error {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	if rl.outFh == nil {
		return nil
	}

	rl.outFh.Close()
	rl.outFh = nil

	return nil
}

func (rl *RotateLog) CurrentFileName() string {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()

	return rl.curFn
}

type cleanupGuard struct {
	enable bool
	fn     func()
	mutex  sync.Mutex
}

func (g *cleanupGuard) Enable() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.enable = true
}

func (g *cleanupGuard) Run() {
	g.fn()
}