package gostd

import (
	"errors"
	"io/fs"
	"testing"
)

/*
! New 返回一个形式为给定文本 `text` 的 error。
! Unwrap 调用 `err` 的 `Unwrap` 方法（若实现，否则返回 `nil`）。它不会调用 `Join` 返回的 `joinErr` 的 `Unwrap` 方法。
! Join 包装给定的一组 `error`，并丢弃所有的 `nil`。返回的 `error` 实现了 Unwrap 方法。
! Is 报告 errs-tree 中是否有任何错误与 `target` 匹配。
! As 在 errs-tree 中查找与 `target` 类型匹配的第一个错误，找到时则将 `target` 设置为该错误值并返回 `true`。否则返回 `false`。
	如果 `target` 不是指向实现错误的类型或任何接口类型的非 `nil` 指针，则 panic。
*/

func TestErrors(t *testing.T) {
	var err1, err2, nilerr error = errors.New("Err1"), errors.New("Err2"), nil

	t.Run("Join err", func(t *testing.T) {
		errs := errors.Join(err1, err2, nilerr)
		checkErr(errs)
		errs = errors.Join(errs, errors.New("Join new Err3"))
		checkErr(errs)
	})

	t.Run("Unwrap", func(t *testing.T) {
		mErr := newMError(&fs.PathError{"Op", "Path", errors.New("test PathError")})
		checkErr(mErr)
		// errors.Unwrap use PathError.Unwrap
		checkErr(errors.Unwrap(mErr))
	})

	t.Run("As & Is", func(t *testing.T) {
		var perr error = &fs.PathError{"Op", "Path", errors.New("test PathError")}

		joins := errors.Join(newMError(errors.New("test MError")), err1, err2, perr)

		if errors.Is(joins, perr) {
			logf("joins is &fs.PathError{}")
		}

		var MErr *MError
		if errors.As(joins, MErr) {
			checkErr(MErr)
		}
	})
}

type MError struct {
	merr error
}

func newMError(err error) error {
	return &MError{err}
}

func (e *MError) Error() string {
	return "MError: " + e.merr.Error()
}

func (e *MError) Unwrap() error {
	logf("Call (*MError).Unwrap")
	return errors.Unwrap(e.merr)
}

func (e *MError) As(target any) bool {
	logf("Call (*MError).As")
	return errors.As(e.merr, target)
}
func (e *MError) Is(target error) bool {
	logf("Call (*MError).Is")
	return errors.Is(e.merr, target)
}
