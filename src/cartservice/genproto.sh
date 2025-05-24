#!/bin/bash -eu
#
# Copyright 2018 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# [START gke_cartservice_genproto]

# Generate C# protobuf and gRPC files
# This script assumes protoc and grpc_tools_node_protoc_plugin are available

protodir=../../protos
cartprotodir=./protos

# Generate from global protos
protoc --proto_path=$protodir --csharp_out=. --grpc_out=. --plugin=protoc-gen-grpc=grpc_csharp_plugin $protodir/demo.proto

# Generate from cart-specific protos  
protoc --proto_path=$cartprotodir --csharp_out=. --grpc_out=. --plugin=protoc-gen-grpc=grpc_csharp_plugin $cartprotodir/Cart.proto

# [END gke_cartservice_genproto]
