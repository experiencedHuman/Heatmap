//
// DO NOT EDIT.
//
// Generated by the protocol buffer compiler.
// Source: AccessPoint.proto
//

//
// Copyright 2018, gRPC Authors All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
import GRPC
import NIO
import SwiftProtobuf


/// Usage: instantiate `Api_APServiceClient`, then call methods of this protocol to make API calls.
public protocol Api_APServiceClientProtocol: GRPCClient {
  var serviceName: String { get }
  var interceptors: Api_APServiceClientInterceptorFactoryProtocol? { get }

  func getAccessPoint(
    _ request: Api_Empty,
    callOptions: CallOptions?
  ) -> UnaryCall<Api_Empty, Api_AccessPoint>

  func listAccessPoints(
    _ request: Api_Empty,
    callOptions: CallOptions?,
    handler: @escaping (Api_AccessPoint) -> Void
  ) -> ServerStreamingCall<Api_Empty, Api_AccessPoint>
}

extension Api_APServiceClientProtocol {
  public var serviceName: String {
    return "api.APService"
  }

  /// Unary call to GetAccessPoint
  ///
  /// - Parameters:
  ///   - request: Request to send to GetAccessPoint.
  ///   - callOptions: Call options.
  /// - Returns: A `UnaryCall` with futures for the metadata, status and response.
  public func getAccessPoint(
    _ request: Api_Empty,
    callOptions: CallOptions? = nil
  ) -> UnaryCall<Api_Empty, Api_AccessPoint> {
    return self.makeUnaryCall(
      path: "/api.APService/GetAccessPoint",
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeGetAccessPointInterceptors() ?? []
    )
  }

  /// Server streaming call to ListAccessPoints
  ///
  /// - Parameters:
  ///   - request: Request to send to ListAccessPoints.
  ///   - callOptions: Call options.
  ///   - handler: A closure called when each response is received from the server.
  /// - Returns: A `ServerStreamingCall` with futures for the metadata and status.
  public func listAccessPoints(
    _ request: Api_Empty,
    callOptions: CallOptions? = nil,
    handler: @escaping (Api_AccessPoint) -> Void
  ) -> ServerStreamingCall<Api_Empty, Api_AccessPoint> {
    return self.makeServerStreamingCall(
      path: "/api.APService/ListAccessPoints",
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeListAccessPointsInterceptors() ?? [],
      handler: handler
    )
  }
}

public protocol Api_APServiceClientInterceptorFactoryProtocol {

  /// - Returns: Interceptors to use when invoking 'getAccessPoint'.
  func makeGetAccessPointInterceptors() -> [ClientInterceptor<Api_Empty, Api_AccessPoint>]

  /// - Returns: Interceptors to use when invoking 'listAccessPoints'.
  func makeListAccessPointsInterceptors() -> [ClientInterceptor<Api_Empty, Api_AccessPoint>]
}

public final class Api_APServiceClient: Api_APServiceClientProtocol {
  public let channel: GRPCChannel
  public var defaultCallOptions: CallOptions
  public var interceptors: Api_APServiceClientInterceptorFactoryProtocol?

  /// Creates a client for the api.APService service.
  ///
  /// - Parameters:
  ///   - channel: `GRPCChannel` to the service host.
  ///   - defaultCallOptions: Options to use for each service call if the user doesn't provide them.
  ///   - interceptors: A factory providing interceptors for each RPC.
  public init(
    channel: GRPCChannel,
    defaultCallOptions: CallOptions = CallOptions(),
    interceptors: Api_APServiceClientInterceptorFactoryProtocol? = nil
  ) {
    self.channel = channel
    self.defaultCallOptions = defaultCallOptions
    self.interceptors = interceptors
  }
}

