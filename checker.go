/**
 * @Author: ChenJunJi
 * @Date: 2023/12/21
 * @Desc:
**/

package timehelper

import "time"

type Checker struct {
	interval time.Duration
	next     time.Time
}

type Option func(checker *Checker)

func WithDelay(delay time.Duration) Option {
	return func(checker *Checker) {
		checker.next = time.Now().Add(delay)
	}
}
func NewChecker(interval time.Duration, opts ...Option) *Checker {
	checker := &Checker{}
	for _, opt := range opts {
		opt(checker)
	}

	checker.interval = interval
	if checker.next.Before(time.Now()) {
		checker.next = time.Now().Add(interval)
	}
	return checker
}

func (checker *Checker) Check() bool {
	return checker.next.After(time.Now())
}

func (checker *Checker) Next() time.Time {
	return checker.next
}

func (checker *Checker) CheckAndSet(ignore bool) bool {
	now := time.Now()
	if checker.next.After(now) {
		return false
	}
	if ignore {
		checker.next = now.Add(checker.interval)
	} else {
		checker.next = checker.next.Add(checker.interval)
	}
	return true
}
