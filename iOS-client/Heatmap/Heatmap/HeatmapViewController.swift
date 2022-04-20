/* Copyright (c) 2017 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import Foundation
import GoogleMaps
import SwiftUI
import UIKit
import GoogleMapsUtils

struct HeatmapControllerRepresentable: UIViewControllerRepresentable {
    func makeCoordinator() -> Coordinator {
        Coordinator()
    }
    
    func makeUIViewController(context: Context) -> HeatmapViewController {
        let heatmapViewController = HeatmapViewController()
        return heatmapViewController
    }
    
    func updateUIViewController(_ uiViewController: HeatmapViewController, context: Context) {
        
    }
}

class HeatmapViewController: UIViewController, GMSMapViewDelegate, GMSIndoorDisplayDelegate {
    private var mapView: GMSMapView!
    private var heatmapLayer: GMUHeatmapTileLayer!
    private var colorMSLabel = UILabel(),
                opacityLabel = UILabel(),
                radiusLabel = UILabel(),
                minZILabel = UILabel(),
                maxZILabel = UILabel(),
                gStartLabel = UILabel(),
                gEndLabel = UILabel(),
                zoomLvlLabel = UILabel()
    
    private var gradientColors = [UIColor.green, UIColor.red]
    private var gradientStartPoints = [0.2, 0.6] as [NSNumber]
    
    private var gradientStart: Float = 0.2,
                gradientEnd: Float = 0.6,
                opacity: Float = 0.8,
                radius = UInt(300),
                colorMapSize = UInt(256)
    
    private var accesspoints: [GMUWeightedLatLng]!
    private var zoomLevel: Float = 16.0
    private var heatmapLayers: [GMUHeatmapTileLayer?]!
    
    private var backgroundView = UIView(frame: CGRect(x: 5, y: 5, width: 350, height: 400))
    
    override func loadView() {
        super.loadView()
        let camera = GMSCameraPosition.camera(withLatitude: 48.14957600438307, longitude: 11.567179933190348, zoom: zoomLevel)
        mapView = GMSMapView.map(withFrame: CGRect(x: 0, y: 0, width: 500, height: 800), camera: camera)
        mapView.delegate = self
        mapView.indoorDisplay.delegate = self
        self.view = mapView
        addSliders() // TODO: add slider for tileSize ?
        
        backgroundView.backgroundColor = .lightGray
        self.view.addSubview(backgroundView)
        
        //add markers
        let position = CLLocationCoordinate2D(latitude: 48.998, longitude: 11.668)
        let marker = GMSMarker(position: position)
        marker.title = "Access Point"
        marker.map = mapView
        
        let zoomTitle = UILabel(frame: CGRect(x: 5, y: 360, width: 150, height: 15))
        zoomTitle.text = "Zoom Lvl:"
        backgroundView.addSubview(zoomTitle)
        zoomLvlLabel.frame = CGRect(x: 160, y: 360, width: 80, height: 15)
        zoomLvlLabel.text = "0.0"
        backgroundView.addSubview(zoomLvlLabel)
    }
    
    override func viewDidLoad() {
        super.viewDidLoad()
        // TODO move this to a function and reinitalize heatmapLayer each time radius is changed for example. There might be a bug that radius change is not taken into consideration
        accesspoints = addHeatmap()
        setupMultipleLayers()
        initHeatmapLayer(radius: radius, opacity: opacity, colorMapSize: colorMapSize)
        print("zoom level is: \(mapView.camera.zoom)")
    }
    
    private func setupMultipleLayers() {
        heatmapLayers = [GMUHeatmapTileLayer?](repeating: nil, count: 23)
        for i in 1...23 {
            let hmLayer = GMUHeatmapTileLayer()
            if i <= 9 {
                hmLayer.radius = UInt(i)
                hmLayer.opacity = 0.7
            } else if i > 9 && i < 17 {
                hmLayer.radius = UInt(i * 5)
                hmLayer.opacity = 0.8
            } else {
                hmLayer.radius = UInt(i * 20)
                hmLayer.opacity = 0.95
            }
            hmLayer.weightedData = accesspoints
//            hmLayer.zIndex = Int32(i)
//            hmLayer.map = mapView
            heatmapLayers[i - 1] = hmLayer
        }
    }
    
    func mapView(_ mapView: GMSMapView, didChange position: GMSCameraPosition) {
        let currentZoomLevel = mapView.camera.zoom
        zoomLvlLabel.text = String(currentZoomLevel)
        if abs(currentZoomLevel - zoomLevel) > 1 {
            zoomLevel = currentZoomLevel //update
//            heatmapLayer = heatmapLayers[Int(zoomLevel)]
            let lvl = Int(zoomLevel)
            let hml = heatmapLayers[lvl]
            hml?.map = mapView
            
            // TODO: create more heatmap layers and update more often, not only in an interval of 1, but even 0.3 for example.
            // TODO: show markers on zoom >18 or >19 and then decide to remove or to still keep the heatmap layer
            
            //remove the other heatmaps
            for i in 0...22 {
                if i != lvl {
                    let otherHML = heatmapLayers[i]
                    otherHML?.map = nil
                }
            }
            print("zoom lvl: \(Int(zoomLevel)) & radius: \(hml?.radius)")
//            heatmapLayer.clearTileCache()
//            initHeatmapLayer(radius: radius, opacity: opacity, colorMapSize: colorMapSize)
        }
    }
    
    
    // TODO get building data. show heatmap based just on the building?
    func didChangeActiveBuilding(_ building: GMSIndoorBuilding?) {
        if let currentBuilding = building {
            let levels = currentBuilding.levels
            mapView.indoorDisplay.activeLevel = levels[2]
//            print("Changed to level 2")
        }
    }
    
    // TODO display different heatmap based on selected level
    func didChangeActiveLevel(_ level: GMSIndoorLevel?) {
//        print("level was changed")
    }
    
    private func initHeatmapLayer(radius: UInt, opacity: Float, colorMapSize: UInt) {
        heatmapLayer = GMUHeatmapTileLayer()
        heatmapLayer.radius = radius // old value = 300 , 50,
        heatmapLayer.opacity = opacity // default = 0.8 or 0.75
        heatmapLayer.fadeIn = true
        heatmapLayer.gradient = GMUGradient(colors: gradientColors,
                                            startPoints: gradientStartPoints,
                                            colorMapSize: colorMapSize) // 256
        heatmapLayer.weightedData = accesspoints

        //      heatmapLayer.zIndex // <- this is useful when you have multiple heatmap layers and want to decide which shows on top of which
        
        // Set the heatmap to the mapview.
//        heatmapLayer.map = mapView
    }
    
    // Parse JSON data and add it to the heatmap layer.
    func addHeatmap() -> [GMUWeightedLatLng] {
        var list = [GMUWeightedLatLng]()
        do {
            // Get the data: latitude/longitude positions of police stations.
            if let path = Bundle.main.url(forResource: "accessPoints", withExtension: "json") {
                let data = try Data(contentsOf: path)
                let json = try JSONSerialization.jsonObject(with: data, options: [])
                if let object = json as? [[String: Any]] {
                    for item in object {
                        let lat = item["Latitude"]
                        let lng = item["Longitude"]
                        var intensity = item["Intensity"] as? Double ?? 1.0
                        if intensity != 1.0 {
                            intensity *= 100000.0
                        }
                        //                        print("intensity is : \(Float(intensity))")
                        let coords = GMUWeightedLatLng(coordinate: CLLocationCoordinate2DMake(lat as! CLLocationDegrees, lng as! CLLocationDegrees), intensity: Float(intensity))
                        list.append(coords)
                    }
                } else {
                    print("Could not cast data from JSON")
                }
            } else {
                print("Could not read JSON data")
            }
        } catch {
            print(error.localizedDescription)
        }
        // Add the latlngs to the heatmap layer.
//      heatmapLayer.weightedData = list
        return list
    }
    
    private func addSliders() {
        let fromLeft = 5, width = 170, height = 15
        // MARK: colorMapSize
        addUISlider(title: "colorMS", xOrig: fromLeft, yOrig: 20, width: width, height: height, minVal: 3, maxVal: 1000.0, objcFunc: #selector(self.colorMapSizeSliderValueDidChange(_:)), &colorMSLabel)
        
        // MARK: Opacity
        addUISlider(title: "Opacity", xOrig: fromLeft, yOrig: 80, width: width, height: height, minVal: 0.05, maxVal: 10.0, objcFunc: #selector(self.opacitySliderValueDidChange(_:)), &opacityLabel)
        
        // MARK: Radius
        addUISlider(title: "Radius", xOrig: fromLeft, yOrig: 120, width: width, height: height, minVal: 1, maxVal: 500, objcFunc: #selector(self.radiusSliderValueDidChange(_:)), &radiusLabel)
        
        // MARK: Minimum Zone Intensity
        addUISlider(title: "MinZI", xOrig: fromLeft, yOrig: 180, width: width, height: height, minVal: 0, maxVal: 10, objcFunc: #selector(self.minZoneIntSliderValueDidChange(_:)), &minZILabel)
        
        // MARK: Maximum Zone Intensity
        addUISlider(title: "MaxZI", xOrig: fromLeft, yOrig: 220, width: width, height: height, minVal: 0, maxVal: 20, objcFunc: #selector(self.maxZoneIntSliderValueDidChange(_:)), &maxZILabel)
        
        // MARK: Gradient Start
        addUISlider(title: "G-start", xOrig: fromLeft, yOrig: 280, width: width, height: height, minVal: 0, maxVal: 0.59, objcFunc: #selector(self.gradientStartSliderValueDidChange(_:)), &gStartLabel)
        
        // MARK: Gradient End
        addUISlider(title: "G-end", xOrig: fromLeft, yOrig: 320, width: width, height: height, minVal: 0.6, maxVal: 1.0, objcFunc: #selector(self.gradientEndSliderValueDidChange(_:)), &gEndLabel)
    }
    
    private func addUISlider(title: String,
                             xOrig: Int,
                             yOrig: Int,
                             width: Int,
                             height: Int,
                             minVal: Float,
                             maxVal: Float,
                             objcFunc: Selector,
                             _ valueLabel: inout UILabel) {
        
        let slider = UISlider(frame:CGRect(x: xOrig + 90, y: yOrig, width: width, height: height))
        slider.minimumValue = minVal
        slider.maximumValue = maxVal
        slider.isContinuous = true
        slider.addTarget(self, action: objcFunc, for: .valueChanged)
        backgroundView.addSubview(slider)
        
        let titleLabel = UILabel(frame: CGRect(x: xOrig, y: yOrig, width: 70, height: height))
        titleLabel.text = title
        backgroundView.addSubview(titleLabel)
        
        valueLabel.frame = CGRect(x: xOrig + 270, y: yOrig, width: width, height: height)
        valueLabel.text = "0.0"
        backgroundView.addSubview(valueLabel)
    }
    
    @objc
    func colorMapSizeSliderValueDidChange(_ sender: UISlider!) {
        let step = 1
        let stepValue = Int(sender.value) / step * step
        sender.value = Float(stepValue)
        colorMapSize = UInt(stepValue)
        colorMSLabel.text = String(colorMapSize)
        heatmapLayer.gradient = GMUGradient(colors: gradientColors,
                                            startPoints: gradientStartPoints,
                                            colorMapSize: colorMapSize)
        heatmapLayer.clearTileCache()
    }
    
    @objc
    func opacitySliderValueDidChange(_ sender: UISlider!) {
        let step:Float = 0.05
        let stepValue = sender.value / step * step
        sender.value = stepValue
        heatmapLayer.opacity = stepValue
        opacityLabel.text = String(heatmapLayer.opacity)
//        heatmapLayer.clearTileCache() <-- it slows down the heatmap rendering when changing opacity
    }
    
    @objc
    func radiusSliderValueDidChange(_ sender: UISlider!) {
        let step = 50
        let stepValue = Int(sender.value) / step * step
        sender.value = Float(stepValue)
        radius = UInt(stepValue)
        heatmapLayer.radius = radius
        radiusLabel.text = String(radius)
        heatmapLayer.clearTileCache()
    }
    
    @objc
    func minZoneIntSliderValueDidChange(_ sender: UISlider!) {
        let step = 1
        let stepValue = Int(sender.value) / step * step
        sender.value = Float(stepValue)
        heatmapLayer.minimumZoomIntensity = UInt(stepValue)
        minZILabel.text = String(heatmapLayer.minimumZoomIntensity)
        heatmapLayer.clearTileCache()
    }
    
    @objc
    func maxZoneIntSliderValueDidChange(_ sender: UISlider!) {
        let step = 1
        let stepValue = Int(sender.value) / step * step
        sender.value = Float(stepValue)
        heatmapLayer.maximumZoomIntensity = UInt(stepValue)
        maxZILabel.text = String(heatmapLayer.maximumZoomIntensity)
        heatmapLayer.clearTileCache()
    }
    
    @objc
    func gradientStartSliderValueDidChange(_ sender: UISlider!) {
        let step:Float = 0.05
        let stepValue = sender.value / step * step
        sender.value = stepValue
        gradientStart = stepValue
        gradientStartPoints = [gradientStart, gradientEnd] as [NSNumber]
        heatmapLayer.gradient = GMUGradient(colors: gradientColors,
                                            startPoints: gradientStartPoints,
                                            colorMapSize: colorMapSize)
        gStartLabel.text = String(gradientStart)
        heatmapLayer.clearTileCache()
    }
    
    @objc
    func gradientEndSliderValueDidChange(_ sender: UISlider!) {
        let step:Float = 0.05
        let stepValue = sender.value / step * step
        sender.value = stepValue
        gradientEnd = stepValue
        gradientStartPoints = [gradientStart, gradientEnd] as [NSNumber]
        heatmapLayer.gradient = GMUGradient(colors: gradientColors,
                                            startPoints: gradientStartPoints,
                                            colorMapSize: colorMapSize)
        gEndLabel.text = String(gradientEnd)
        heatmapLayer.clearTileCache()
    }
}
