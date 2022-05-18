//
//  DataRepository.swift
//  HeatmapUIKit
//
//  Created by Kamber Vogli on 16.05.22.
//

import Foundation
import GRPC
import Logging

class DataRepository {
  static let shared = DataRepository()
  private var client: Helloworld_GreeterClient?
  
  private init() {
    let eventLoopGroup = PlatformSupport.makeEventLoopGroup(loopCount: 1)
    var logger = Logger(label: "gRPC", factory: StreamLogHandler.standardError(label:))
    logger.logLevel = .debug
    
    let channel = ClientConnection
//      .usingPlatformAppropriateTLS(for: eventLoopGroup)
//      .withBackgroundActivityLogger(logger)
      .insecure(group: eventLoopGroup)
      .connect(host: "192.168.0.109", port: 50051)
    
    let callOptions = CallOptions(logger: logger)
    client = Helloworld_GreeterClient(channel: channel, defaultCallOptions: callOptions)
    
    print("Connected to gRPC server")
  }
  
  func greetServer(message: String) {
    var req = Helloworld_HelloRequest()
    req.name = message
    let result = client?.sayHello(req, callOptions: .none)
    print("33")
    result?.response.whenComplete({ res in
      do {
        let reply = try res.get()
        print(reply.message)
      } catch {
        print("could not get reply")
      }
    })
  }
  
}
