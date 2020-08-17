package repo

import (
	"encoding/json"
	"errors"
	"github.com/tidwall/buntdb"
)

type Mask []string

type Change interface {
	Type() string
}

type SuiteAggUpdate struct {
	SuiteAgg
}

func (u SuiteAggUpdate) Type() string {
	return "suite_agg_update"
}

type SuiteUpsert struct {
	Suite
	Mask `json:"mask,omitempty"`
}

func (u SuiteUpsert) Type() string {
	return "suite_upsert"
}

type changeHandler interface {
	handleChanges(changes []Change)
}

type watcher struct {
	in  chan<- []Change
	out <-chan []Change
}

func newWatcher() *watcher {
	in := make(chan []Change)
	out := make(chan []Change)
	var buf [][]Change
	getNext := func() []Change {
		if len(buf) == 0 {
			return nil
		}
		return buf[0]
	}
	getOut := func() chan<- []Change {
		if len(buf) == 0 {
			return nil
		}
		return out
	}
	go func() {
		defer close(out)
		for {
			select {
			case changes, ok := <-in:
				if !ok {
					return
				}
				buf = append(buf, changes)
			case getOut() <- getNext():
				buf = buf[1:]
			}
		}
	}()
	return &watcher{in, out}
}

type SuiteWatcher struct {
	r *Repo
	*watcher

	id    string
	padLt int
	padGt int

	n   int
	min entry
	max entry
}

func (r *Repo) WatchSuites(id string, padLt, padGt int) (*SuiteWatcher, error) {
	w := SuiteWatcher{
		r:       r,
		watcher: newWatcher(),
	}
	r.addHandler(&w)
	return &w, r.db.View(func(tx *buntdb.Tx) error {
		var changes []Change
		add := func(k, v string) bool {
			var upsert SuiteUpsert
			if err := json.Unmarshal([]byte(v), &upsert.Suite); err != nil {
				panic(err)
			}
			changes = append(changes, upsert)
			return true
		}
		var n int
		lt := func(k, v string) bool {
			w.min = entry{k, v}
			n++
			return add(k, v)
		}
		eq := func(k, v string) bool {
			w.min = entry{k, v}
			w.max = entry{k, v}
			n++
			return add(k, v)
		}
		gt := func(k, v string) bool {
			w.max = entry{k, v}
			n++
			return add(k, v)
		}
		if err := itrAroundSuite(tx, id, padLt, padGt, lt, eq, gt); err != nil {
			return err
		}
		w.id = id
		w.padLt = padLt
		w.padGt = padGt
		w.n = n

		var agg SuiteAgg
		err := r.byId(SuiteAggColl, "", &agg)
		if err != nil && !errors.Is(err, errNotFound{}) {
			return err
		}
		w.in <- append(changes, SuiteAggUpdate{agg})
		return nil
	})
}

func (w *SuiteWatcher) Close() {
	w.r.removeHandler(w)
	close(w.in)
}

func (w *SuiteWatcher) Changes() <-chan []Change {
	return w.out
}

func itrAroundSuite(tx *buntdb.Tx, id string, padLt, padGt int,
	lt, eq, gt itr) error {
	firstEq := firstCond(eq)
	restLt := restCond(lt)
	restGt := restCond(gt)
	if id == "" {
		return tx.Descend(suiteIndexStartedAt,
			andCond(firstEq, restLt, limitCond(padGt+padLt)))
	}
	idVal, err := tx.Get(key(SuiteColl, id))
	if err == buntdb.ErrNotFound {
		return errNotFound{}
	} else if err != nil {
		return err
	}
	pivot := suiteIndexStartedAtPivot(idVal)
	err = tx.AscendGreaterOrEqual(suiteIndexStartedAt, pivot,
		andCond(firstEq, restGt, limitCond(padGt)))
	if err != nil {
		return err
	}
	if padLt == 0 {
		return nil
	}
	return tx.DescendLessOrEqual(suiteIndexStartedAt, pivot,
		andCond(restLt, limitCond(padLt)))
}

func (w *SuiteWatcher) handleChanges(changes []Change) {

}
