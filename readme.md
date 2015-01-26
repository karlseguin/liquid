# Liquid Templates For Go

```go
template, err := liquid.ParseString("hello {{ name | upcase }}", nil)
if err != nil { panic(err) }
data := map[string]interface{}{
  "name": "leto",
}
writer := new(bytes.Buffer)
template.Render(writer, data)
return writer.String()
```

Given a file path, liquid can also `ParseFile`. Give a `[]byte` it can also simply `Parse`.

## What's Missing
The following filters are missing:

- map

The following tags are missing:

- cycle

**As per the original Liquid template, there's no automatic protection against XSS. This might change if I decide to turn this into something more than just a clone**

## Template Cache
By default the templates are cached in a pretty dumb cache. That is, once in the cache, items stay in the cache (there's no expiry). The cache can be disabled, on a per-template basis, via:

```go
template, _ := liquid.Parse(someByteTemplate, liquid.Configure().Cache(nil))
//OR
template, _ := liquid.ParseString(someStringTemplate, liquid.NoCache)
```

Alternatively, you can provide your own `core.Cache` implementation which
could implement expiry and other custom features.

## Configuration
As seen above, a configuration can be provided when generating a template. Configuration is achieved via a fluent-interface. Configurable options are:

- `Cache(cache core.Cache)`: the caching implementation to use. This defaults to `liquid.SimpleCache`, which is thread-safe.
- `IncludeHandler(handler core.IncludeHandler)`: the callback used for handling includes. By default, includes are ignored. See below for more information.
- `PreserveWhitespace()`: By default, Liquid will slightly compact whitespace around tags. It doesn't do a perfect job, but it does reduce the whitespace noise. This method lets you skip this whitespace compaction

## Global Configuration
A few global configuration methods are available:

- `SetInternalBuffer(count, size int)`: the number of internal buffers and the maximum size of each buffer to use. Defaults to 512 and 4KB. This is currently only used for the capture tag. If you need to capture more than 4KB, increase the 2nd value.

## Data Binding
The template's `Render` method takes a `map[string]interface{}` as its argument. Beyond that, `Render` works on all built-in types, and will also reflect the exported fields of a struct.

```go
type User struct {
  Name  string
  Manager *User
}

t, _ := liquid.ParseString("{{ user.manager.name }}", nil)
t.Render(os.Stdout, map[string]interface{}{
  "user": &User{"Duncan", &User{"Leto", nil}},
})
```

or can be a **map**,

```go
t, _ := liquid.ParseString("{{ user.manager.name }}", nil)
t.Render(os.Stdout, map[string]interface{}{
	"user": map[string]interface{}{
		"manager": map[string]interface{}{"name": "Leto"} ,
		},
})
```

Notice that the template fields aren't case sensitive. If you're exporting fields such as `FirstName` and `Firstname` then shame on you. Make sure to downcase map keys.

Complex objects should implement the `fmt.Stringer` interface (which is Go's toString() equivalent):

```go
func (u *User) String() string {
  return u.Name
}

t, _ := liquid.ParseString("{{ user.manager }}", nil)
```

Failing this, `fmt.Sprintf("%v")` is used to generate a value. At this point, it's really more for debugging purposes.


## Filters
You can add custom filters by calling `core.RegisterFilter`. **The filter lookup is not thread safe; it is expected that you'll add filters on init and then leave it alone**.

It's best to look at the existing filters for ideas on how to proceed. Briefly, there are two types of filters: those with parameters and those without. To support both from a single interface, each filter has a factory. For filters without parameters, the factory is simple:

```go
func UpcaseFactory(parameters []string) core.Filter {
  return Upcase
}
func Upcase(input interface{}, data map[string]interface{}) interface{} {
  //todo
}
```

For filters that expect parameters, a little more work is needed:

```go
func JoinFactory(parameters []core.Value) core.Filter {
  if len(parameters) == 0 {
    return defaultJoin.Join
  }
  return (&JoinFilter{parameters[0]}).Join
}
type JoinFilter struct {
  glue core.Value
}
func (f *JoinFilter) Join(input interface{}, data map[string]interface{}) interface{} {

}
```

It's a good idea to provide default values for parameters!

Finally, do note that Filters work with and return `interface{}`. Consider using a [type switch](http://golang.org/doc/effective_go.html#type_switch) with a `default` case which returns the input.

If you're filter works with `string` or `[]byte`, you should handle both `string` and `[]byte` types as you don't know what the user will provide nor what transformation previous filters might apply. Similarly, if you're expecting an array, you should handle both arrays and slices.

Again, look at existing filters for more insight.

## Include
The include tag is supported by configuring a custom `IncludeHandler`. The handler is responsible for resolving the include. This provides the greatest amount of flexibility: included templates can be loaded from the file system, a database or some other location.

For example, an include handler which loads templates from the filsystem, might look like:

```go
var config = liquid.Configure().IncludeHandler(includeHandler)

func main {
  template, err := liquid.ParseString(..., config)
  //...
}

// name is equal to the paramter passed to include
// data is the data available to the template
func includeHandler(name string, writer io.Writer, data map[string]interface{}) {
  // not sure if this is good enough, but do be mindful of directory traversal attacks
  fileName := path.Join("./templates", string.Replace(name, "..", ""))
  template, _ := liquid.ParseFile(fileName, config)
  template.Render(writer, data)
}
```

Sadly, include doesn't currently support the more advanced variations, such as specifying a specific value for the include or automatically including in a loop. However, the flexibility provided hopefully suffices for now.

## Errors
`Render` shouldn't fail, but it doesn't always stay silent about mistakes. For the sake of helping debug issues, it can inject data within the template. For example, using the output tag such as {{ user.name.first | upcase }} when *user.name.first* doesn't map to valid data will result in the literal "{{user.name.first}}"" being injecting in the template.

`Parse` and its variants will return an error if the template is not valid. The error message is meant to be helpful and can be shown as-is to users.
