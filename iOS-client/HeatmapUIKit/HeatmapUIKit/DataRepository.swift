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
  private var apClient: Api_APServiceClient?
  
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
//    client = Helloworld_GreeterClient(channel: channel, defaultCallOptions: callOptions)
    
    apClient = Api_APServiceClient(channel: channel, defaultCallOptions: callOptions)
    print("Connected to gRPC server")
  }
  
  func greetServer(message: String) {
    var req = Helloworld_HelloRequest()
    req.name = message
    let result = client?.sayHello(req, callOptions: .none)
    
    result?.response.whenComplete({ res in
      do {
        let reply = try res.get()
        print(reply.message)
      } catch {
        print("could not get reply")
      }
    })
  }
  
  func getAP() {
    let result = apClient?.getAccessPoint(Api_Empty(), callOptions: .none)
    result?.response.whenComplete({ res in
      do {
        let reply = try res.get()
        print(reply.debugDescription)
      } catch {
        print("could not get access point")
      }
    })
  }
  
}
