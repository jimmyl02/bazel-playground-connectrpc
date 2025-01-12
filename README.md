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

make sure we carefuully setup BUILD.bazel and gazelle with the correct `go_prefix`! for the instantiation, we just need to add a dependency in our `MODULE.bazel` then instantiate it in our root `BUILD.bazel`. exact documentation can be found [here](https://github.com/bazel-contrib/rules_go/blob/master/docs/go/core/bzlmod.md) and it's pretty accurate.

```
# gazelle:prefix github.com/jimmyl02/bazel-playground-connectrpc
```

we can then run gazelle with

```
bazelisk run //:gazelle
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

##### external dependency

when adding an external dependency, it is now recommended to use a go.mod which is parsed by the `go_deps` bazel extension. this means when adding a dependency, it should be through the standard `go get -u <package>` command.

```
go mod init github.com/jimmyl02/bazel-playground-connectrpc
bazelisk run @rules_go//go -- get -u github.com/moznion/go-optional
bazelisk run @rules_go//go -- mod tidy -e
bazelisk run //:gazelle
bazelisk mod tidy
```

it is required that the explicit dependencies are specified in `use_repo` of the root `MODULE.bazel`. the series go `go mod tidy`, then running gazelle and `bazel mod tidy` should be capable of automatically updating the `use_repo` list of dependencies. this allows us to just maintain the source of truth in one location which is `go.mod`. this has additional benefits in letting us use standard tooling go with this repository if we ever want to.

##### internal dependency

adding an internal dependency should be automatically handled by just calling the important in golang then running gazelle.

```
# after adding the dependency in code, there is a "metadata missing" error; anywhere run:
bazelisk run //:gazelle
```

#### run the cmd

```
bazel run //cmd/server
```

#### build and run the cmd

```
bazel build //cmd/testserver
./bazel-bin/cmd/server/server_/server
```

## connectrpc

with bazelmod being the new default, there are even less guides on how to properly configure it with gazelle. this is a walkthrough of how I've configured bazel for this playground.

in general, it seems we want to do the following:

1. keep the protobuf outputs of `protoc-gen-go`
2. generate the connectrpc service framework using `protoc-gen-connect-go`

#### setup

##### protoc-gen-go

the first step is to get protobuf working at all. this should be done by writing the proto file into a proto directory then run gazelle to generate the `BUILD.bazel`

notice that running `bazel build //...` fails because we are missing `@@com_google_protobuf`

we can fix this by adding add it to our MODULE.bazel thus adding the dependency

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

we find that we also can't just invoke the plugin as it generates in a subfolder with additional options. the buf team [recently added an option to generate in the same package](https://github.com/connectrpc/connect-go/discussions/310#discussioncomment-11765339) which means that by calling the plugin with the `package_suffix` we can get a super nicely working protobuf + connectrpc generation.

the awesome part of bazel is that it's a true monorepo! this means we can define our compiler / plugin in our root `BUILD.bazel` and reference it in the proto `BUILD.bazel`. you can find the final definition in the root `BUILD.bazel` and notice that we can add the gazelle directive `# gazelle:go_grpc_compilers @rules_go//proto:go_proto,//:connect_go_proto_compiler` which properly calls our newly defined compiler after the `go_proto` prefined compiler

## debugging

#### unexpected end of JSON input

this error occurs when something is wrong with the overall bazel configuration, the best way to debug is to attempt to build something and seeing what is wrong

## conclusion

this playground is always subject to change! but I will try to tag working versions at different stages. it has been a good (and to be honest, time consuming and quite painful) learning experience. I hope this repository is helpful to others who are trying to get started with bazel and especially with things like `protoc` plugins
