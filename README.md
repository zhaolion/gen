# gen
A package for generating things based on go files. In dev now :)

> inspired by [gengo](https://github.com/kubernetes/gengo)

# packages usage

## parser

demo:

```
builder := New()
builder.SetDebugLevel()
_ = builder.AddDirRecursive("./testpkg")
universe, _ := builder.FindTypes()
```
