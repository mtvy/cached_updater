package cache

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"gotest.tools/assert"
)

type stubResponser struct {
	sync.RWMutex
	val interface{}
	err error
}

func (s *stubResponser) Get(ctx context.Context) (interface{}, error) {
	s.RLock()
	defer s.RUnlock()
	return s.val, s.err
}

type stubVariationAuctionResponser struct{}

func (s *stubVariationAuctionResponser) Get(ctx context.Context) (interface{}, error) {
	return &test{
		Advertisers: []testAdv{
			{200},
		},
	}, nil
}

type gotExpected struct {
	isRespNil bool
	gotErr    bool
}

type testAdv struct {
	ID int
}

type test struct {
	Advertisers []testAdv
}

// Проверяем Get с параллельным update
func TestGet(t *testing.T) {
	testCases := []struct {
		name           string
		stub           *stubResponser
		timeout        time.Duration
		updateInterval time.Duration
		want           gotExpected
	}{
		{
			// Кэш доступен
			name: "get_true_test_#1",
			stub: &stubResponser{
				val: &struct{}{},
				err: nil,
			},
			timeout:        time.Second * 1,
			updateInterval: time.Millisecond * 1,
			want: gotExpected{
				isRespNil: false,
				gotErr:    false,
			},
		},
		{
			// Кэш  недоступен из-за ошибки
			name: "get_false_test_#2",
			stub: &stubResponser{
				val: &struct{}{},
				err: errors.New("stub error"),
			},
			timeout:        time.Second * 1,
			updateInterval: time.Millisecond * 1,
			want: gotExpected{
				isRespNil: false,
				gotErr:    true,
			},
		},
		{
			// Кэш  недоступен из-за nil response
			name: "get_false_test_#2",
			stub: &stubResponser{
				val: nil,
				err: nil,
			},
			timeout:        time.Second * 1,
			updateInterval: time.Millisecond * 1,
			want: gotExpected{
				isRespNil: true,
				gotErr:    false,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
			defer cancel()

			c := NewCacheResponseDecorator(tc.updateInterval, tc.stub)
			c.Start(ctx)

			time.Sleep(tc.updateInterval)

			for {
				select {
				case <-ctx.Done():
					return
				default:
					resp, err := c.Get(ctx)
					assert.Equal(t, tc.want.isRespNil, resp == nil)
					assert.Equal(t, tc.want.gotErr, err != nil)
				}
			}
		})
	}

	// Кэш то недоступен, то доступен
	t.Run("get_variation_test_#3", func(t *testing.T) {
		moreCounter, lessCounter := 0, 0
		defer func() { fmt.Println("moreCounter: " + strconv.Itoa(moreCounter)) }()
		defer func() { fmt.Println("lessCounter: " + strconv.Itoa(lessCounter)) }()

		rand.Seed(time.Now().UnixNano())
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		defer cancel()
		c := NewCacheResponseDecorator(
			time.Millisecond*1,
			&stubVariationAuctionResponser{},
		)
		c.Start(ctx)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				resp, _ := c.Get(ctx)
				if resp.(test).Advertisers[0].ID == 200 {
					moreCounter++
					continue
				}
				lessCounter++
			}
		}
	})
}
