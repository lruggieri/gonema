package cache

import (
	"fmt"
	"github.com/lruggieri/gonema/pkg/util"
	"strconv"
	"testing"
	"time"
)

func TestCache_InsertAndFetch(t *testing.T) {

	baseElementsToInsert := []interface{}{
		2,
		"ciao",
		"TokyoRulez",
		"SushiLover",
		"IActuallyDontLikeWasabi",
		false,
		3, 141592653589,
		"The answer to everything is ...",
		42,
	}

	c := NewCache()
	elementToInsert := 100000
	insertedElements := make([]interface{}, 0, elementToInsert)
	for i := 0; i < elementToInsert; i++ {
		elementToInsert := baseElementsToInsert[util.GetRandomPositiveInt(len(baseElementsToInsert))]
		c.Insert(CacheElementRoot(strconv.Itoa(i)), "key", elementToInsert)
		insertedElements = append(insertedElements, elementToInsert)
	}

	if len(c.elements) != elementToInsert {
		t.Error("inserted 1 element in cache, but there are ", len(c.elements), " elements in cache")
		t.FailNow()
	}

	for i := 0; i < elementToInsert; i++ {
		key := generateCacheKey(CacheElementRoot(strconv.Itoa(i)), "key")

		if c.elements[key].value != insertedElements[i] {
			t.Error("INSERT FAILURE: expecting to get element '", insertedElements[i], "' for key ", key, " but got '", c.elements[key], "'")
			t.FailNow()
		}

		fetchedElement := c.Fetch(CacheElementRoot(strconv.Itoa(i)), "key")
		if fetchedElement != c.elements[key].value {
			t.Error("FETCH FAILURE: expecting to get element '", c.elements[key].value, "' for key ", key, " but got '", fetchedElement, "'")
			t.FailNow()
		}

	}
}

func TestCacheElement_Deadline(t *testing.T) {
	cacheRefreshTime = 10 * time.Minute //so it doesn't get triggered during the test
	c := NewCache()
	cacheRefreshTime = 10 * time.Millisecond

	element1Duration := 100 * time.Millisecond
	element2Duration := 200 * time.Millisecond
	c.SetNewRootElementDuration("root1", element1Duration)
	c.SetNewRootElementDuration("root2", element2Duration)

	c.Insert("root1", "key", 1234)
	c.Insert("root2", "key", 1234)

	startInsertion := time.Now()
	element1Deadline := c.elements[generateCacheKey("root1", "key")].deadline
	element2Deadline := c.elements[generateCacheKey("root2", "key")].deadline
	elapsedInsertion := time.Since(startInsertion)

	if element1Deadline.After(element2Deadline) {
		t.Error("element 1 cannot have a deadline after element 2 one")
		t.FailNow()
	}

	if element1Deadline.After(startInsertion.Add(elapsedInsertion + element1Duration)) {
		t.Error("element 1 cannot have a deadline (", element1Deadline.String(), ") grater than startInsertion"+
			"(", startInsertion.String(), ") + insertion time (", elapsedInsertion.String(), ") + element1 duration (", element1Duration, ")")
		t.FailNow()
	}

	if element2Deadline.After(startInsertion.Add(elapsedInsertion + element2Duration)) {
		t.Error("element 2 cannot have a deadline (", element2Deadline.String(), ") grater than startInsertion"+
			"(", startInsertion.String(), ") + insertion time (", elapsedInsertion.String(), ") + element2 duration (", element2Duration, ")")
		t.FailNow()
	}

	if element2Deadline.After(
		element1Deadline.Add(elapsedInsertion + element2Duration)) {

		t.Error("element 2 cannot have a deadline (", element2Deadline.String(), ") grater than element 1 one"+
			"(", element1Deadline.String(), ") + insertion time (", elapsedInsertion.String(), ") + element2 duration (", element2Duration, ")")
		t.FailNow()
	}

	c = NewCache()
	element1Duration = 100 * time.Millisecond
	element2Duration = 200 * time.Millisecond
	c.SetNewRootElementDuration("root1", element1Duration)
	c.SetNewRootElementDuration("root2", element2Duration)

	c.Insert("root1", "key", 1234)
	c.Insert("root2", "key", 1234)

	element1 := c.elements[generateCacheKey("root1", "key")]
	element2 := c.elements[generateCacheKey("root2", "key")]
	time.Sleep(element1Duration)
	if !(element1.Deadline()) {
		t.Error("expecting element1 to be expired after its own root duration (", element1Duration, ")")
		t.FailNow()
	}
	if element2.Deadline() {
		t.Error("not expecting element2 to be expired after element 1 duration (", element1Duration, "), which is far lower")
		t.FailNow()
	}
	time.Sleep(element2Duration - element1Duration) //we already waited for element 1 duration
	if !element2.Deadline() {
		t.Error("expecting element2 to be expired after its own root duration (", element2Duration, ")")
		t.FailNow()
	}

}

func ExampleCache_Insert() {
	c := NewCache()
	c.Insert("root1", "key1", "value1")
	fmt.Println(c.Fetch("root1", "key1"))
	// Output: value1
}
func ExampleCache_Fetch() {
	c := NewCache()

	c.Insert("root1", "key1", "value1")
	c.Insert("root1", "key2", "value2")
	c.Insert("root2", "key1", "value3")

	fmt.Println(c.Fetch("root1", "key1"))
	fmt.Println(c.Fetch("root1", "key2"))
	fmt.Println(c.Fetch("root2", "key1"))
	// Output:
	// value1
	// value2
	// value3

}
func ExampleCache_Reset() {
	c := NewCache()
	c.Insert("root1", "key1", "value1")
	c.Reset()
	fmt.Println(c.Fetch("root1", "key1"))
	// Output: <nil>
}
func ExampleCache_SetNewRootElementDuration() {
	c := NewCache()
	c.SetNewRootElementDuration("root1", time.Second)

	fmt.Println(c.cacheRootElementsDuration["root1"])

	// Output: 1s
}
