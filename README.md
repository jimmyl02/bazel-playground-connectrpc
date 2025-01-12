# bazel-playground-connectrpc

this is a little playground to better learn how to work with bazel, protobuf, and connectrpc. this is based off of the [bazel playground for grpc](https://github.com/jimmyl02/bazel-playground)

## setup

#### install bazelisk and bazel

```
# install bazelisk and bazel
brew install bazelisk
bazel
```

#### setup bazel / golang in repository

optionally, we can define a strict bazel version to use with `.bazelliskrc` by setting `USE_BAZEL_VERSION=8.0.0`

setup MODULE.bazel with [docs](https://github.com/bazel-contrib/rules_go/blob/master/docs/go/core/bzlmod.md)

make sure we carefuully setup BUILD.bazel and gazelle with the correct `go_prefix`!

```
# gazelle:prefix github.com/jimmyl02/bazel-playground-connectrpc
```

we can then run gazelle with

```
bazel run //:gazelle
```

#### setup vscode with bazel

very helpful guide [here](https://github.com/bazelbuild/rules_go/issues/3014)

create scripts/gopackagesdriver.sh and make it executable

```
#!/bin/bash

exec bazel run -- @rules_go//go/tools/gopackagesdriver "$@"
```

edit the workspace preferences

```
{
    "go.toolsEnvVars": {
        "GOPACKAGESDRIVER": "${workspaceFolder}/scripts/gopackagesdriver.sh"
    }
}
```

## golang

#### import a new dependency with gazelle

when adding an external dependency, it is now recommended to use a go.mod which is parsed by the `go_deps` bazel extension. this means when adding a dependency, it should be through the standard `go get -u <package>` command.

```
go mod init github.com/jimmyl02/bazel-playground-connectrpc
bazelisk run @rules_go//go -- get -u github.com/moznion/go-optional
bazelisk run @rules_go//go -- mod tidy -e
bazelisk run //:gazelle
bazelisk mod tidy
```

after this, it is required the manually specify the package in `use_repo` of the root `MODULE.bazel`

add an internal dependency

```
# after adding the dependency in code, there is a "metadata missing" error; anywhere run:
gazelle
```

#### run the cmd

```
bazel run //cmd/testcmd
```

#### build and run the cmd

```
bazel build //cmd/testcmd
./bazel-bin/cmd/testcmd/testcmd_/testcmd
```

## connectrpc

with bazelmod being the new default, there are even less guides on how to properly configure it with gazelle. this is a walkthrough of how I've configured bazel for this playground.

in general, it seems we want to do the following:

1. keep the protobuf outputs of `protoc-gen-go`
2. generate the connectrpc service framework using `protoc-gen-connect-go`

#### setup

##### protoc-gen-go

first we need to update the gazelle goo grpc compilers by adding the directive `# gazelle:go_grpc_compilers @rules_go//proto:go_grpc` in our `BUILD.bazel`. this controls the `compilers` property of the `go_proto_library` and is the first step in getting it to output correctly with bazelmod.

write the proto file into a proto directory then run gazelle to generate the `BUILD.bazel`

notice that running `bazel build //...` fails because we are missing `@@com_google_protobuf`

we can add it to our MODULE.bazel by adding the dependency

```
bazel_dep(name = "protobuf", version = "29.3", repo_name = "com_google_protobuf")
```

there is another step after rules_go was updated to 48.0, we have to add rules_proto to the MODULE.bazel file
`bazel_dep(name = "rules_proto", version = "7.0.2")`

now running `bazel build //...` works!

##### protoc-gen-connect-go

now let's focus on getting `protoc-gen-connect-go` plugin to work with bazel / gazelle grpc generation.

for this, we need to change gazelle's `go_grpc_compiler` to use the `protoc-gen-connect-go` plugin from buf.

the quick summary of how to use these plugins with bazel is through definitions of `go_proto_compiler`. the language is pretty confusing because here is references a "compiler" while other places define this as a "plugin".

we need to define a compiler / plugin for `protoc-gen-connect-go` which is then invoked. however, one tricky part is that we also need to run `@rules_go//proto:go_proto` which is a predefined compiler rule because it seems like once there is no output from the compiled protobuf types without explicitly defining it as a plugin. for exact details about how to di this, take a look at the `proto/testproto/BUILD.bazel` file.

## debugging

#### unexpected end of JSON input

this error occurs when something is wrong with the overall bazel configuration, the best way to debug is to attempt to build something and seeing what is wrong
