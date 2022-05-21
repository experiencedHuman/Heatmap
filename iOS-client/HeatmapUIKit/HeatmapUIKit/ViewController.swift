//
//  ViewController.swift
//  HeatmapUIKit
//
//  Created by Kamber Vogli on 15.03.22.
//

import SwiftUI
import UIKit
import AzureMapsControl

class ViewController: UIViewController, AzureMapDelegate {
  private var azureMap: MapControl!
  private var heatmapSource, apSource: DataSource!
  private var popup = Popup()
  private var datePicker: UIDatePicker = UIDatePicker()
  private var selectedTime: String = ""
  
  
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
    setupDataSource(heatmapSource)
    setupDataSource(apSource)
    
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
  
  // handle clicking on an access point icon
  func azureMap(_ map: AzureMap, didTapOn features: [Feature]) {
    guard let feature = features.first else {
      // Popup has been released or no features provided
      return
    }
    
    // Create the custom view for the popup.
    let popupView = PopupTextView()
    popupView.layer.cornerRadius = 30
    popupView.layer.masksToBounds = true
    
    selectedTime = datePicker.date.formatted(date: .omitted, time: .shortened)
    
    if let clusterCount = feature.properties["point_count"] as? Int {
      // let count = heatmapSource.leaves(of: feature, offset: 0, limit: .max)
      popupView.setText("\(selectedTime) Uhr: Cluster of \(clusterCount) APs")
//      DataRepository.shared.getAP()
      let res = DataRepository.shared.getAPs()
      for ap in res {
        print(ap.debugDescription)
      }
    } else {
      let text = feature.properties["title"] as! String
      popupView.setText("\(selectedTime) Uhr: Nicht so stark besucht! \(text)")
      DataRepository.shared.getAP()
    }
    
    // Set the text to the custom view.
    
    
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
    
    datePicker.frame = CGRect(x: 150, y: 750, width: 170, height: 35)
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
  
  @objc
  func datePickerValueChanged(_ sender: UIDatePicker) {
    let dateFormatter: DateFormatter = DateFormatter()
    dateFormatter.dateFormat = "MM/dd/yyyy hh:mm"
    let selectedDate: String = dateFormatter.string(from: sender.date)
    print("Selected hour \(selectedDate)")
    selectedTime = sender.date.formatted(date: .omitted, time: .shortened)
  }
  
  private func setupDataSource(_ dataSource: DataSource) {
    let locations = readCoordsFromJSON(file: "ap-2")
    
    var id = 1
    for location in locations {
      let feature = Feature(Point(location))
      
      id += 1
      // Add a property to the feature.
      feature.addProperty("title", value: "Access Point \(id)")
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
              16: 30,
              17: 40,
              18: 45,
              19: 50,
              23: 1000
            ])
          )
        ),
        .heatmapOpacity(0.8),
        .minZoom(1.0),
        .maxZoom(24),
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
                                    .textOffset(CGVector(dx: 0, dy: 1)),
                                    //                                    .iconOffset(CGVector(dx: 0, dy: -10)),
                                    .textSize(20),
                                    .textHaloBlur(1.0),
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
  
  private func readCoordsFromJSON(file filename: String) -> [CLLocationCoordinate2D] {
    var coordinates: [CLLocationCoordinate2D] = []
    do {
      if let path = Bundle.main.url(forResource: filename, withExtension: "json") {
        let data = try Data(contentsOf: path)
        let json = try JSONSerialization.jsonObject(with: data, options: [])
        if let object = json as? [[String: Any]] {
          for item in object {
            let lat  = item["Lat"] as? Double ?? 0.0
            let long = item["Long"] as? Double ?? 0.0
            let coord = CLLocationCoordinate2D(latitude: lat, longitude: long)
            coordinates.append(coord)
          }
        }
      }
    } catch {
      print("Could not read json file!")
      return coordinates
    }
    return coordinates
  }
}
