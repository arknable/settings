# settings

A very simple Go yaml configuration file loader.

## How to Use

Assuming you want to load `database.yaml` in `myapp` folder:

```
collection, err := settings.NewCollection("database", "myapp")
if err != nil {
    // do something    
}

model := &struct {
    Address string `yaml:"address"`
    Port    int    `yaml:"port"`
}{}

if err := collection.Load(model); err != nil {
    // use default or do something.
}
```

If you want a different extension:
```
collection, err := settings.NewCollection("database", "myapp")
if err != nil {
    // do something    
}
collection.Extension = "conf" // will look for database.conf instead.

model := &struct {
    Address string `yaml:"address"`
    Port    int    `yaml:"port"`
}{}

if err := collection.Load(model); err != nil {
    // use default or do something.
}
```

If you want to add additional search path:
```
collection, err := settings.NewCollection("database", "myapp")
if err != nil {
    // do something    
}
collection.SearchPaths = append(collection.SearchPaths, "/path/to/somewhere")

model := &struct {
    Address string `yaml:"address"`
    Port    int    `yaml:"port"`
}{}

if err := collection.Load(model); err != nil {
    // use default or do something.
}
```

default search paths are
1. `$HOME/.config`
2. `/etc`				(Non Windows)
3. `/usr/local/etc`	(Non Windows)
