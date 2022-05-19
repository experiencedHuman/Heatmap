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
    
    apClient = Api_APServiceClient(channel: channel, defaultCallOptions: callOptions)
    print("Connected to gRPC server")
  }
  
  // TODO: pass ap name or id as parameter
  func getAP() {
    let result = apClient?.getAccessPoint(Api_Empty(), callOptions: .none)
    result?.response.whenComplete({ res in
      do {
        let reply = try res.get()
        print(reply.debugDescription)
      } catch {
        print("Could not get the access point!")
      }
    })
  }
  
  func getAPs() -> [Api_AccessPoint] {
    var apList: [Api_AccessPoint] = []
    let result = apClient?.listAccessPoints(Api_Empty(), callOptions: .none, handler: { api_AccessPoint in
      apList.append(api_AccessPoint)
    })
    do {
      _ = try result?.status.wait()
    } catch {
      print("Could not get the list of access points!")
    }
    return apList
  }
  
}
