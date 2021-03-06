# jo

jo is short for JSON Object and is convenience package for dealing
with JSON.

## Get Example

```
package main

import (
	"fmt"
	"github.com/mgutz/jo"
)

func main() {
	json, _ := jo.DecodeStr(`{
		"a": 1,
		"b": "cow",
		"c": 1.2,
		"d": {
			"e": true,
			"f": false,
			"g": {
				"h": [0,1,2,3],
				"i": [
					{"j": 1},
					{"k": "Moo"}
				]
			}
		}
	}`)

	// Get value in json path casted into different variable types

	a, _ := json.Int("a")
	fmt.Printf("a (%T) = %v\n", a, a)

	b, _ := json.String("b")
	fmt.Printf("a (%T) = %v\n", b, b)

	c, _ := json.Float("c")
	fmt.Printf("c (%T) = %v\n", c, c)

	de, _ := json.Bool("d.e")
	fmt.Printf("d.e (%T) = %v\n", de, de)

	dg, _ := json.Map("d.g")
	fmt.Printf("d.g (%T) = %v\n", dg, dg)

	dgh1, _ := json.Int("d.g.h[1]")
	fmt.Printf("d.g.h[1] (%T) = %v\n", dgh1, dgh1)

	dgi1k, _ := json.String("d.g.i[1].k")
	fmt.Printf("d.g.i[1].k (%T) = %v\n", dgi1k, dgi1k)
}
```

This example will generate the following output -

```
a (int) = 1
a (string) = cow
c (float64) = 1.2
d.e (bool) = true
d.g (map[string]interface {}) = map[h:[0 1 2 3] i:[map[j:1] map[k:Moo]]]
d.g.h[1] (int) = 1
d.g.i[1].k (string) = Moo
```

Set Example

```
func main() {
	json := nestedjson.New()

	json.Set("a.b.c", 1)

	json.Set("a.b.d", []interface{}{1, 2, 3})

	json.Set("a.b.e", map[string]interface{}{
		"f": "Hello", "g": "World",
	})

	json.Set("a.b.e.g", "Universe")

	json.Set("a.b.d[0]", 6.9)

	jsonStr, _ := json.EncodePrettyStr()
	fmt.Println(jsonStr)

}
```

This will generate the following JSON document -

```
{
  "a": {
    "b": {
      "c": 1,
      "d": [
        6.9,
        2,
        3
      ],
      "e": {
        "f": "Hello",
        "g": "Universe"
      }
    }
  }
}
```

## Misc

Remember a JSON object isn't necessarily a map. These are all valid JSON

```
"string"
1
[1, "s", true]
{"foo":"bar"}
```

## Credit

`jo` was originally a port from [nestedjson](https://github.com/wenxiang/go-nestedjson).
The main difference is `jo` is loosely typed and allows conversion to any type.
For example, a `Float64` can be retrieved as a `string`.

## TODO

- Proper godocs
- Delete function
- More tests

## LICENSE

MIT
