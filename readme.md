**Parser rewrite. Tag support is completely gone**

# Liquid Templates For Go

    template, err := liquid.ParseString("hello {{ name | upcase }}", nil)
    if err != nil { panic(err) }
    data := map[string]interface{
      "name": "leto",
    }
    println(template.Render(data))

Given a file path, liquid can also `ParseFile`. Give a `[]byte` it can also simply `Parse`.

## What's Missing
The following filters are missing:

- date
- sort
- map
- escape
- escape_once
- strip_html
- truncate
- truncatewords
- split
- modulo

The following tags are missing (all):

- assign
- capture
- case
- comment
- cycle
- for
- if
- include
- raw
- unless

Other:

- Support for 'now'
- Render can generate far less objects

## Template Cache
By default the templates are cached in a pretty dumb cache based. That is, once in the cache, items stay in the cache (there's no expiry). The cache can be disabled, on a per-template basis, via:

    template, _ := liquid.Parse(someByteTemplate, liquid.Configure().Cache(nil))
    //OR
    template, _ := liquid.ParseString(someStringTemplate, liquid.NoCache)

Alternatively, you can provide your own `liquid.Cache` implementation which
could implement expiry and other custom features.

## Configuration
As seen above, a configuration can be provided when generating a template. Configuration is achieved via a fluent-interface. Configurable options are:

- `Cache(cache liquid.Cache)`: the caching implementation to use. This defaults to `liquid.SimpleCache`, which is thread-safe.

## Data Binding
The template's `Render` method takes a `map[string]interface{}` as its argument. Beyond that, `Render` works on all built-in types, and will also reflect the exported fields of a struct.

    type User struct {
      Name  string
      Manager *User
    }

    t, _ := liquid.ParseString("{{ user.manager.name }}", nil)
    t.Render(map[string]interface{}{
      "user": &User{"Duncan", &User{"Leto", nil}},
    })

Notice that the template fields aren't case sensitive. If you're exporting fields such as `FirstName` and `Firstname` then shame on you. Make sure to downcase map keys.

Complex objects should implement the `fmt.Stringer` interface (which is Go's toString() equivalent):

    func (u *User) String() string {
      return u.Name
    }

    t, _ := liquid.ParseString("{{ user.manager }}", nil)

Failing this, `fmt.Sprintf("%v")` is used to generate a value. At this point, it's really more for debugging purposes.


## Filters
You can add custom filters by writing in to the `liquid.Filters` map. **This map is not thread safe; it is expected that you'll add filters on init and then leave it alone**.

It's best to look at the existing filters for ideas on how to proceed. Briefly, there are two types of filters: those with parameters and those without. To support both from a single interface, each filter has a factory. For filters without parameters, the factory is simple:

    func UpcaseFactory(parameters []string) Filter {
      return Upcase
    }
    func Upcase(input interface{}) interface{} {
      //todo
    }

For filters that expect parameters, a little more work is needed:

    func JoinFactory(parameters []string) Filter {
      return &JoinFilter{
        glue: []byte(parameters[0]),
      }
    }
    type JoinFilter struct {
      glue []byte
    }
    func (f *JoinFilter) Join(input interface{}) interface{} {

    }

It's a good idea to provide default values for parameters!

Finally, do note that Filters work with and return `interface{}`. Consider using a [type switch](http://golang.org/doc/effective_go.html#type_switch) with a `default` case which returns the input.

If you're filter works with `string` or `[]byte`, you should handle both `string` and `[]byte` types as you don't know what the user will provide nor what transformation previous filters might apply. Similarly, if you're expecting an array, you should handle both arrays and slices.

Again, look at existing filters for more insight.

## Errors
`Render` shouldn't fail, but it doesn't always stay silent about mistakes. For the sake of helping debug issues, it can inject data within the template. For example, using the output tag such as {{ user.name.first | upcase }} when *user.name.first* doesn't map to valid data will result in the literal "{{user.name.first}}"" being injecting in the template.

`Parse` and its variants will return an error if the template is not valid. The error message is meant to be helpful and can be shown as-is to users.
