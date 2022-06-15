// DO NOT EDIT.
// swift-format-ignore-file
//
// Generated by the Swift generator plugin for the protocol buffer compiler.
// Source: api/AccessPoint.proto
//
// For information on using the generated types, please see the documentation:
//   https://github.com/apple/swift-protobuf/

import Foundation
import SwiftProtobuf

// If the compiler emits an error on this type, it is because this file
// was generated by a version of the `protoc` Swift plug-in that is
// incompatible with the version of SwiftProtobuf to which you are linking.
// Please ensure that you are building against the same version of the API
// that was used to generate this file.
fileprivate struct _GeneratedWithProtocGenSwiftVersion: SwiftProtobuf.ProtobufAPIVersionCheck {
  struct _2: SwiftProtobuf.ProtobufAPIVersion_2 {}
  typealias Version = _2
}

public struct Api_AccessPoint {
  // SwiftProtobuf.Message conformance is added in an extension below. See the
  // `Message` and `Message+*Additions` files in the SwiftProtobuf library for
  // methods supported on all messages.

  public var name: String = String()

  public var lat: String = String()

  public var long: String = String()

  public var intensity: String = String()

  public var max: Int64 = 0

  public var min: Int64 = 0

  public var unknownFields = SwiftProtobuf.UnknownStorage()

  public init() {}
}

public struct Api_APRequest {
  // SwiftProtobuf.Message conformance is added in an extension below. See the
  // `Message` and `Message+*Additions` files in the SwiftProtobuf library for
  // methods supported on all messages.

  public var name: String = String()

  public var timestamp: String = String()

  public var unknownFields = SwiftProtobuf.UnknownStorage()

  public init() {}
}

public struct Api_APResponse {
  // SwiftProtobuf.Message conformance is added in an extension below. See the
  // `Message` and `Message+*Additions` files in the SwiftProtobuf library for
  // methods supported on all messages.

  public var accesspoint: Api_AccessPoint {
    get {return _accesspoint ?? Api_AccessPoint()}
    set {_accesspoint = newValue}
  }
  /// Returns true if `accesspoint` has been explicitly set.
  public var hasAccesspoint: Bool {return self._accesspoint != nil}
  /// Clears the value of `accesspoint`. Subsequent reads from it will return its default value.
  public mutating func clearAccesspoint() {self._accesspoint = nil}

  public var unknownFields = SwiftProtobuf.UnknownStorage()

  public init() {}

  fileprivate var _accesspoint: Api_AccessPoint? = nil
}

#if swift(>=5.5) && canImport(_Concurrency)
extension Api_AccessPoint: @unchecked Sendable {}
extension Api_APRequest: @unchecked Sendable {}
extension Api_APResponse: @unchecked Sendable {}
#endif  // swift(>=5.5) && canImport(_Concurrency)

// MARK: - Code below here is support for the SwiftProtobuf runtime.

fileprivate let _protobuf_package = "api"

extension Api_AccessPoint: SwiftProtobuf.Message, SwiftProtobuf._MessageImplementationBase, SwiftProtobuf._ProtoNameProviding {
  public static let protoMessageName: String = _protobuf_package + ".AccessPoint"
  public static let _protobuf_nameMap: SwiftProtobuf._NameMap = [
    1: .same(proto: "name"),
    2: .same(proto: "lat"),
    3: .same(proto: "long"),
    4: .same(proto: "intensity"),
    5: .same(proto: "max"),
    6: .same(proto: "min"),
  ]

  public mutating func decodeMessage<D: SwiftProtobuf.Decoder>(decoder: inout D) throws {
    while let fieldNumber = try decoder.nextFieldNumber() {
      // The use of inline closures is to circumvent an issue where the compiler
      // allocates stack space for every case branch when no optimizations are
      // enabled. https://github.com/apple/swift-protobuf/issues/1034
      switch fieldNumber {
      case 1: try { try decoder.decodeSingularStringField(value: &self.name) }()
      case 2: try { try decoder.decodeSingularStringField(value: &self.lat) }()
      case 3: try { try decoder.decodeSingularStringField(value: &self.long) }()
      case 4: try { try decoder.decodeSingularStringField(value: &self.intensity) }()
      case 5: try { try decoder.decodeSingularInt64Field(value: &self.max) }()
      case 6: try { try decoder.decodeSingularInt64Field(value: &self.min) }()
      default: break
      }
    }
  }

  public func traverse<V: SwiftProtobuf.Visitor>(visitor: inout V) throws {
    if !self.name.isEmpty {
      try visitor.visitSingularStringField(value: self.name, fieldNumber: 1)
    }
    if !self.lat.isEmpty {
      try visitor.visitSingularStringField(value: self.lat, fieldNumber: 2)
    }
    if !self.long.isEmpty {
      try visitor.visitSingularStringField(value: self.long, fieldNumber: 3)
    }
    if !self.intensity.isEmpty {
      try visitor.visitSingularStringField(value: self.intensity, fieldNumber: 4)
    }
    if self.max != 0 {
      try visitor.visitSingularInt64Field(value: self.max, fieldNumber: 5)
    }
    if self.min != 0 {
      try visitor.visitSingularInt64Field(value: self.min, fieldNumber: 6)
    }
    try unknownFields.traverse(visitor: &visitor)
  }

  public static func ==(lhs: Api_AccessPoint, rhs: Api_AccessPoint) -> Bool {
    if lhs.name != rhs.name {return false}
    if lhs.lat != rhs.lat {return false}
    if lhs.long != rhs.long {return false}
    if lhs.intensity != rhs.intensity {return false}
    if lhs.max != rhs.max {return false}
    if lhs.min != rhs.min {return false}
    if lhs.unknownFields != rhs.unknownFields {return false}
    return true
  }
}

extension Api_APRequest: SwiftProtobuf.Message, SwiftProtobuf._MessageImplementationBase, SwiftProtobuf._ProtoNameProviding {
  public static let protoMessageName: String = _protobuf_package + ".APRequest"
  public static let _protobuf_nameMap: SwiftProtobuf._NameMap = [
    1: .same(proto: "name"),
    2: .same(proto: "timestamp"),
  ]

  public mutating func decodeMessage<D: SwiftProtobuf.Decoder>(decoder: inout D) throws {
    while let fieldNumber = try decoder.nextFieldNumber() {
      // The use of inline closures is to circumvent an issue where the compiler
      // allocates stack space for every case branch when no optimizations are
      // enabled. https://github.com/apple/swift-protobuf/issues/1034
      switch fieldNumber {
      case 1: try { try decoder.decodeSingularStringField(value: &self.name) }()
      case 2: try { try decoder.decodeSingularStringField(value: &self.timestamp) }()
      default: break
      }
    }
  }

  public func traverse<V: SwiftProtobuf.Visitor>(visitor: inout V) throws {
    if !self.name.isEmpty {
      try visitor.visitSingularStringField(value: self.name, fieldNumber: 1)
    }
    if !self.timestamp.isEmpty {
      try visitor.visitSingularStringField(value: self.timestamp, fieldNumber: 2)
    }
    try unknownFields.traverse(visitor: &visitor)
  }

  public static func ==(lhs: Api_APRequest, rhs: Api_APRequest) -> Bool {
    if lhs.name != rhs.name {return false}
    if lhs.timestamp != rhs.timestamp {return false}
    if lhs.unknownFields != rhs.unknownFields {return false}
    return true
  }
}

extension Api_APResponse: SwiftProtobuf.Message, SwiftProtobuf._MessageImplementationBase, SwiftProtobuf._ProtoNameProviding {
  public static let protoMessageName: String = _protobuf_package + ".APResponse"
  public static let _protobuf_nameMap: SwiftProtobuf._NameMap = [
    1: .same(proto: "accesspoint"),
  ]

  public mutating func decodeMessage<D: SwiftProtobuf.Decoder>(decoder: inout D) throws {
    while let fieldNumber = try decoder.nextFieldNumber() {
      // The use of inline closures is to circumvent an issue where the compiler
      // allocates stack space for every case branch when no optimizations are
      // enabled. https://github.com/apple/swift-protobuf/issues/1034
      switch fieldNumber {
      case 1: try { try decoder.decodeSingularMessageField(value: &self._accesspoint) }()
      default: break
      }
    }
  }

  public func traverse<V: SwiftProtobuf.Visitor>(visitor: inout V) throws {
    // The use of inline closures is to circumvent an issue where the compiler
    // allocates stack space for every if/case branch local when no optimizations
    // are enabled. https://github.com/apple/swift-protobuf/issues/1034 and
    // https://github.com/apple/swift-protobuf/issues/1182
    try { if let v = self._accesspoint {
      try visitor.visitSingularMessageField(value: v, fieldNumber: 1)
    } }()
    try unknownFields.traverse(visitor: &visitor)
  }

  public static func ==(lhs: Api_APResponse, rhs: Api_APResponse) -> Bool {
    if lhs._accesspoint != rhs._accesspoint {return false}
    if lhs.unknownFields != rhs.unknownFields {return false}
    return true
  }
}
