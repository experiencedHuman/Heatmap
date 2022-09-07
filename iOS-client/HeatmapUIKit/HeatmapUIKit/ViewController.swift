//
//  ViewController.swift
//  HeatmapUIKit
//
//  Created by Kamber Vogli on 15.03.22.
//

import SwiftUI
import UIKit
import AzureMapsControl

struct AccessPoint {
  let coord: CLLocationCoordinate2D
  let intensity: Double
}

class ViewController: UIViewController, AzureMapDelegate {
  private var azureMap: MapControl!
  private var heatmapSource, apSource: DataSource!
  private var popup = Popup()
  private var datePicker = UIDatePicker()
  private var backToCurrentButton = UIButton()
  private var timestamp = ""
  private var accessPoints: [Api_AccessPoint]!// = DataRepository.shared.getAPs()
  private let fromJSON = false
  private var maxIntensity = 0, minIntensity = 0
  
  func azureMap(_ map: AzureMap, didTapAt location: CLLocationCoordinate2D) {
    print("did tap at lat: \(location.latitude) , long: \(location.longitude)")
  }
  
  override func loadView() {
    super.loadView()
    azureMap = MapControl.init(frame: CGRect(x: 0, y: 0, width: 500, height: 800),
                               options: [
                                CameraOption.center(lat: 48.2692083204, lng: 11.6690079838),
                                CameraOption.zoom(15),
                                CameraOption.maxZoom(24)
                               ])
    
    heatmapSource = DataSource()
    apSource = DataSource(options: [.cluster(true)])
    
    if (fromJSON) {
      setupDataSourceFromJSON(heatmapSource)
      setupDataSourceFromJSON(apSource)
    } else {
      let date = Date()
      let dateFormatter: DateFormatter = DateFormatter()
      dateFormatter.dateFormat = "yyyy-MM-dd HH"
      timestamp = dateFormatter.string(from: date)
      accessPoints = DataRepository.shared.getAPs(timestamp: timestamp)
      for ap in accessPoints {
        print("Retrieved ap data: \(ap.name), \(ap.lat), \(ap.long), max: \(ap.max), min: \(ap.min), int: \(ap.intensity)")
      }
      
      setupDataSource(heatmapSource)
      setupDataSource(apSource)
    }
    
    azureMap.onReady { map in
      // Add two different sources, one for clustering, one for single Access Point Symbols
      map.sources.add(self.apSource)
      map.sources.add(self.heatmapSource)
      
      // Load a custom icon image into the image sprite of the map.
      map.images.add(UIImage(named: "3ap")!, withID: "cluster")
      map.images.add(UIImage(named: "single_ap")!, withID: "ap")
      
      map.events.addDelegate(self)
      map.popups.add(self.popup)
    }
    
    self.view.addSubview(azureMap)
  }
  
  // MARK: Popup annotation
  func azureMap(_ map: AzureMap, didTapOn features: [Feature]) {
    guard let feature = features.first else {
      // Popup has been released or no features provided
      return
    }
    
    // Create the custom view for the popup.
    let popupView = PopupTextView()
    popupView.layer.cornerRadius = 30
    popupView.layer.masksToBounds = true
    
    let date = Date()
    let calendar = Calendar.current
    let hour = calendar.component(.hour, from: date)
    let minutes = calendar.component(.minute, from: date)
    let time = "\(hour):\(minutes)"
    
    if let _ = feature.properties["point_count"] as? Int {
      let leaves = apSource.leaves(of: feature, offset: 0, limit: .max)
      var connectedDevices = 0
      for leaf in leaves {
        let apLoad = leaf.properties["intensity"] as? Int ?? 0
        connectedDevices += apLoad
      }
      popupView.setText("\(time) Uhr: Momentan stark besucht: \(connectedDevices) connected devices.")
      
    } else {
      let name = feature.properties["name"] as! String
      popupView.setText("\(time) Uhr: Momentan nicht so stark besucht! \(name)")
      DataRepository.shared.getAccessPointByName(name, timestamp: timestamp)
    }
    
    // Get the position of the tapped feature.
    let position = Math.positions(from: feature).first!
    
    // Set the options on the popup.
    popup.setOptions([
      // Set the popups position.
      .position(position),
      
      // Set the anchor point of the popup content.
      .anchor(.bottom),
      
      // Set the content of the popup.
      .content(popupView),
      
      // Offset the popup by a specified number of points.
      .pointOffset(CGPoint(x: 0, y: -20))
    ])
    
    // Open the popup.
    popup.open()
  }
  
  override func viewDidLoad() {
    super.viewDidLoad()
    
    azureMap.onReady { map in
      self.addHeatmapLayer(map)
      self.addAccessPointLayer(map)
      self.addClusterLayer(map)
    }
    
    addDatePickerView()
    addBackToCurrentView()
  }
  
  private func addDatePickerView() {
    datePicker.frame = CGRect(x: 150, y: 750, width: 170, height: 35)
    datePicker.minimumDate = dateBeforeDaysToSubstract(daysToSubstract: 31)
    datePicker.maximumDate = dateAfterDaysToAdd(daysToAdd: 15)
    datePicker.tintColor = .black
    datePicker.backgroundColor = .gray
    datePicker.layer.cornerRadius = 10
    datePicker.layer.masksToBounds = true
    datePicker.clipsToBounds = true
    datePicker.timeZone = NSTimeZone.local
    //      datePicker.datePickerMode = .date
    datePicker.addTarget(self, action: #selector(ViewController.datePickerValueChanged(_:)), for: .valueChanged)
    view.addSubview(datePicker)
  }
  
  private func addBackToCurrentView() {
    backToCurrentButton.frame = CGRect(x: 340, y: 750, width: 60, height: 35)
    backToCurrentButton.isEnabled = false // disable button
    backToCurrentButton.alpha = 0.0 // make button invisible
    backToCurrentButton.setTitle("Now", for: .normal)
    backToCurrentButton.backgroundColor = UIColor(red: 0/255, green: 123/255, blue: 224/255, alpha: 1.0)
    backToCurrentButton.layer.cornerRadius = 10
    backToCurrentButton.clipsToBounds = true
    backToCurrentButton.addTarget(self, action: #selector(revertHeatmapToNow), for: .touchUpInside)
    view.addSubview(backToCurrentButton)
  }
  
  @objc
  private func revertHeatmapToNow() {
    let dateFormatter: DateFormatter = DateFormatter()
    dateFormatter.dateFormat = "yyyy-MM-dd HH"
    timestamp = dateFormatter.string(from: Date())
    
    let apList = DataRepository.shared.getAPs(timestamp: timestamp)
    
    updateHeatmap(apList: apList)
  }
  
  private func dateBeforeDaysToSubstract(daysToSubstract: Int) -> Date? {
    let currentDate = Date()
    let earlierDate = Calendar.current.date(byAdding: .day, value: -daysToSubstract, to: currentDate)
    return earlierDate
  }
  
  private func dateAfterDaysToAdd(daysToAdd: Int) -> Date? {
    let currentDate = Date()
    let laterDate = Calendar.current.date(byAdding: .day, value: daysToAdd, to: currentDate)
    return laterDate
  }
  
  @objc
  func datePickerValueChanged(_ sender: UIDatePicker) {
    let dateFormatter: DateFormatter = DateFormatter()
    dateFormatter.dateFormat = "yyyy-MM-dd HH"
    timestamp = dateFormatter.string(from: sender.date)
    let timestampCurrent = dateFormatter.string(from: Date())
    
    if timestamp != timestampCurrent {
      // enable & show button
      backToCurrentButton.isEnabled = true
      backToCurrentButton.alpha = 1.0 //make button visible
    } else {
      backToCurrentButton.isEnabled = false
      backToCurrentButton.alpha = 0.0 //make button visible
    }
    
    let apList = DataRepository.shared.getAPs(timestamp: timestamp)
    
    updateHeatmap(apList: apList)
  }
  
  private func setupDataSource(_ dataSource: DataSource) {
    let features = updateIntensities(apList: accessPoints)
    dataSource.add(features: features)
  }
  
  private func updateHeatmap(apList: [Api_AccessPoint]) {
    let features = updateIntensities(apList: apList)
    heatmapSource.set(features: features)
    apSource.set(features: features)
  }
  
  private func updateIntensities(apList: [Api_AccessPoint]) -> [Feature] {
    var max = 0, min = 0
    var features: [Feature] = []
    
    for accessPoint in apList {
      let lat = Double(accessPoint.lat) ?? 0.0
      let long = Double(accessPoint.long) ?? 0.0
      let location = CLLocationCoordinate2D(latitude: lat, longitude: long)
      let apMax = Int(accessPoint.max)
      let apMin = Int(accessPoint.min)
      
      if accessPoint.intensity == 0 {
        continue
      }
      
      if apMax > max {
        max = apMax
      }
      
      if apMin < min {
        min = apMin
      }
      
      let feature = Feature(Point(location))
      // Add properties to the feature.
      feature.addProperty("name", value: "\(accessPoint.name)")
      feature.addProperty("intensity", value: accessPoint.intensity)
      features.append(feature)
    }
    
    maxIntensity = max
    minIntensity = min
    
    return features
  }
  
  private func setupDataSourceFromJSON(_ dataSource: DataSource) {
    let accessPoints = readCoordsFromJSON(file: "ap-2")
    for accessPoint in accessPoints {
      let feature = Feature(Point(accessPoint.coord))
      feature.addProperty("name", value: "test")
      feature.addProperty("intensity", value: accessPoint.intensity)
      dataSource.add(feature: feature)
    }
  }
  
  private func addHeatmapLayer(_ map: AzureMap) {
    let heatmapLayer = HeatMapLayer(
      source: heatmapSource,
      options: [
        .heatmapRadius(
          from: NSExpression(
            forAZMInterpolating: .zoomLevelAZMVariable,
            curveType: .exponential,
            parameters: NSExpression(forConstantValue: 7),
            stops: NSExpression(forConstantValue: [
              // Add multiple interpolation points, x:y
              // In the x-Axis the zoom levels (1-24)
              // In the y-Axis the radius values
              1: 2,
              15: 15,
              16: 20, //+1
              17: 30, //+1
              18: 35, //+1
              19: 40, //+1
              23: 1000
            ])
          )
        ),
        
          .heatmapWeight(
            from: NSExpression(
              forAZMInterpolating: NSExpression(forKeyPath: "intensity"),
              curveType: .linear,
              parameters: NSExpression(forConstantValue: 2),
              stops: NSExpression(forConstantValue: [
                minIntensity: 0, //0.0: minIntensity,
                //              0.01: 5,
                maxIntensity: 500 //0.1: maxIntensity
              ]))
          ),
        .heatmapIntensity(0.35),
        .heatmapOpacity(0.8),
        .minZoom(1.0),
        .maxZoom(24),
        //        .filter(from: NSPredicate(format: "point_count = NIL"))
        .filter(from: NSPredicate(format: "intensity > 0"))
      ]
    )
    map.layers.addLayer(heatmapLayer)
  }
  
  private func addClusterLayer(_ map: AzureMap) {
    let clusterLayer = SymbolLayer(source: apSource,
                                   options: [
                                    .iconImage("cluster"),
                                    .iconSize(0.25),
                                    .textField(from: NSExpression(forKeyPath: "point_count")),
                                    .textOffset(CGVector(dx: 0, dy: -1.5)),
                                    .iconOffset(CGVector(dx: 0, dy: 100)),
                                    .textSize(20),
                                    //                                    .textHaloBlur(12.0),
                                    .textFont(["StandardFont-Bold"]),
                                    .filter(from: NSPredicate(format: "point_count != NIL"))
                                   ]
    )
    map.layers.addLayer(clusterLayer)
  }
  
  private func addAccessPointLayer(_ map: AzureMap) {
    let accessPointLayer = SymbolLayer(source: apSource,
                                       options: [
                                        .iconImage("ap"),
                                        .iconSize(0.25),
                                        .filter(from: NSPredicate(format: "point_count = NIL"))
                                       ])
    map.layers.addLayer(accessPointLayer)
  }
  
  private func readCoordsFromJSON(file filename: String) -> [AccessPoint] {
    var accessPoints: [AccessPoint] = []
    do {
      if let path = Bundle.main.url(forResource: filename, withExtension: "json") {
        let data = try Data(contentsOf: path)
        let json = try JSONSerialization.jsonObject(with: data, options: [])
        if let object = json as? [[String: Any]] {
          for item in object {
            let lat  = item["Lat"] as? Double ?? 0.0
            let long = item["Long"] as? Double ?? 0.0
            let coord = CLLocationCoordinate2D(latitude: lat, longitude: long)
            let intensity = item["intensity"] as? Double ?? 0.0
            let ap = AccessPoint(coord: coord, intensity: intensity)
            accessPoints.append(ap)
          }
        }
      }
    } catch {
      print("Could not read json file!")
      return accessPoints
    }
    return accessPoints
  }
}
