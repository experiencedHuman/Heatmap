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

class HeatmapViewController: UIViewController, GMSMapViewDelegate {
    private var mapView: GMSMapView!
    private var heatmapLayer: GMUHeatmapTileLayer!
    private var button: UIButton!
    
    private var gradientColors = [UIColor.green, UIColor.red]
    private var gradientStartPoints = [0.2, 0.6] as [NSNumber]
    private var gradientStart:Float = 0.2, gradientEnd:Float = 0.6, colorMapSize = UInt(256), radius = UInt(20), opacity:Float = 0.8
    private var accesspoints: [GMUWeightedLatLng]!
    private var zoomLevel: Float = 16.0
    
    override func loadView() {
        let camera = GMSCameraPosition.camera(withLatitude: 48.14957600438307, longitude: 11.567179933190348, zoom: zoomLevel)
        mapView = GMSMapView.map(withFrame: CGRect.zero, camera: camera)
        mapView.delegate = self
        self.view = mapView
        addButtons()
        addSliders()
    }
    
//    func mapViewDidStartTileRendering(_ mapView: GMSMapView) {
//        print("did start")
//    }
//
//    func mapViewDidFinishTileRendering(_ mapView: GMSMapView) {
//        print("did finish")
//    }
    
    override func viewDidLoad() {
        print("viewDidLoad()")
        // TODO move this to a function and reinitalize heatmapLayer each time radius is changed for example. There might be a bug that radius change is not taken into consideration
        // try multiple layers with different zIndex ?
        accesspoints = addHeatmap()
        initHeatmapLayer(radius: radius, opacity: opacity, colorMapSize: colorMapSize)
        print("zoom level is: \(mapView.camera.zoom)")
//        heatmapLayer.tile
    }
    
    private func testFunc() {
        let hl = GMUHeatmapTileLayer()
    }
    
    private func initHeatmapLayer(radius: UInt, opacity: Float, colorMapSize: UInt) {
        heatmapLayer = GMUHeatmapTileLayer()
        heatmapLayer.radius = radius // old value = 300 , 50,
        heatmapLayer.opacity = opacity // default = 0.8 or 0.75
        heatmapLayer.fadeIn = true
        heatmapLayer.gradient = GMUGradient(colors: gradientColors,
                                            startPoints: gradientStartPoints,
                                            colorMapSize: colorMapSize) // 256
//      addHeatmap() // TODO store this in a class variable
        heatmapLayer.weightedData = accesspoints
//      heatmapLayer.zIndex // <- this is useful when you have multiple heatmap layers and want to decide which shows on top of which
        // Set the heatmap to the mapview.
        heatmapLayer.map = mapView
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
    
    func mapView(_ mapView: GMSMapView, didChange position: GMSCameraPosition) {
        let currentZoomLevel = mapView.camera.zoom
        if abs(currentZoomLevel - zoomLevel) > 1 {
            zoomLevel = currentZoomLevel //update
//            heatmapLayer.clearTileCache()
            initHeatmapLayer(radius: radius, opacity: opacity, colorMapSize: colorMapSize)
        }
    }
    
    private func addSliders() {
        // MARK: colorMapSize
        addUISlider(xOrig: 0, yOrig: 200, width: 200, height: 20, minVal: 3, maxVal: 1000.0, objcFunc: #selector(self.colorMapSizeSliderValueDidChange(_:)))
        
        // MARK: Opacity
        addUISlider(xOrig: 0, yOrig: 400, width: 200, height: 20, minVal: 0.05, maxVal: 10.0, objcFunc: #selector(self.opacitySliderValueDidChange(_:)))
        
        // MARK: Radius
        addUISlider(xOrig: 0, yOrig: 450, width: 200, height: 20, minVal: 1, maxVal: 500, objcFunc: #selector(self.radiusSliderValueDidChange(_:)))
        
        // MARK: Minimum Zone Intensity
        addUISlider(xOrig: 0, yOrig: 550, width: 200, height: 20, minVal: 0, maxVal: 10, objcFunc: #selector(self.minZoneIntSliderValueDidChange(_:)))
        
        // MARK: Maximum Zone Intensity
        addUISlider(xOrig: 0, yOrig: 600, width: 200, height: 20, minVal: 0, maxVal: 20, objcFunc: #selector(self.maxZoneIntSliderValueDidChange(_:)))
        
        // MARK: Gradient Start
        addUISlider(xOrig: 0, yOrig: 700, width: 200, height: 20, minVal: 0, maxVal: 0.59, objcFunc: #selector(self.gradientStartSliderValueDidChange(_:)))
        
        // MARK: Gradient End
        addUISlider(xOrig: 0, yOrig: 750, width: 200, height: 20, minVal: 0.6, maxVal: 1.0, objcFunc: #selector(self.gradientEndSliderValueDidChange(_:)))
    }
    
    private func addUISlider(xOrig: Int, yOrig: Int, width: Int, height: Int, minVal: Float, maxVal: Float, objcFunc: Selector) {
        let slider = UISlider(frame:CGRect(x: xOrig, y: yOrig, width: width, height: height))
        slider.minimumValue = minVal
        slider.maximumValue = maxVal
        slider.isContinuous = true
        slider.addTarget(self, action: objcFunc, for: .valueChanged)
        view.addSubview(slider)
    }
    
    func mapView(_ mapView: GMSMapView, didBeginDragging marker: GMSMarker) {
        print("started")
    }
    
    func mapView(_ mapView: GMSMapView, didEndDragging marker: GMSMarker) {
        print("ended")
    }
    
    @objc
    func colorMapSizeSliderValueDidChange(_ sender:UISlider!) {
        let step = 1
        let stepValue = Int(sender.value) / step * step
        sender.value = Float(stepValue)
        colorMapSize = UInt(stepValue)
        heatmapLayer.gradient = GMUGradient(colors: gradientColors,
                                            startPoints: gradientStartPoints,
                                            colorMapSize: colorMapSize)
        heatmapLayer.clearTileCache()
        print("Color map size \(stepValue)")
    }
    
    @objc
    func opacitySliderValueDidChange(_ sender:UISlider!) {
        let step:Float = 0.05
        let stepValue = sender.value / step * step
        sender.value = stepValue
        heatmapLayer.opacity = stepValue
        heatmapLayer.clearTileCache()
//        print("Opacity value \(stepValue)")
//        print("Radius value \(heatmapLayer.radius)")
    }
    
    @objc
    func radiusSliderValueDidChange(_ sender:UISlider!) {
        let step = 50
        let stepValue = Int(sender.value) / step * step
        sender.value = Float(stepValue)
        radius = UInt(stepValue)
//        heatmapLayer.radius = radius
        print("Radius: \(radius)")
//        heatmapLayer.clearTileCache()
        
//        initHeatmapLayer(radius: radius, opacity: opacity, colorMapSize: colorMapSize)
        
    }
    
    @objc
    func minZoneIntSliderValueDidChange(_ sender:UISlider!) {
        let step = 1
        let stepValue = Int(sender.value) / step * step
        sender.value = Float(stepValue)
        heatmapLayer.minimumZoomIntensity = UInt(stepValue)
        heatmapLayer.clearTileCache()
        print("Minimum Zoom Intensity value \(stepValue)")
    }
    
    @objc
    func maxZoneIntSliderValueDidChange(_ sender:UISlider!) {
        let step = 1
        let stepValue = Int(sender.value) / step * step
        sender.value = Float(stepValue)
        heatmapLayer.minimumZoomIntensity = UInt(stepValue)
        heatmapLayer.clearTileCache()
        print("Minimum Zoom Intensity value \(stepValue)")
    }
    
    @objc
    func gradientStartSliderValueDidChange(_ sender:UISlider!) {
        let step:Float = 0.05
        let stepValue = sender.value / step * step
        sender.value = stepValue
        gradientStart = stepValue
        gradientStartPoints = [gradientStart, gradientEnd] as [NSNumber]
        heatmapLayer.gradient = GMUGradient(colors: gradientColors,
                                            startPoints: gradientStartPoints,
                                            colorMapSize: colorMapSize)
        heatmapLayer.clearTileCache()
        print("Gradient start value \(stepValue)")
    }
    
    @objc
    func gradientEndSliderValueDidChange(_ sender:UISlider!) {
        let step:Float = 0.05
        let stepValue = sender.value / step * step
        sender.value = stepValue
        gradientEnd = stepValue
        gradientStartPoints = [gradientStart, gradientEnd] as [NSNumber]
        heatmapLayer.gradient = GMUGradient(colors: gradientColors,
                                            startPoints: gradientStartPoints,
                                            colorMapSize: colorMapSize)
        heatmapLayer.clearTileCache()
        print("Gradient end value \(stepValue)")
    }
    
    func mapView(_ mapView: GMSMapView, didTapAt coordinate: CLLocationCoordinate2D) {
        print("You tapped at \(coordinate.latitude), \(coordinate.longitude)")
    }
    
    @objc
    private func removeHeatmap() {
        heatmapLayer.map = nil
        heatmapLayer = nil
        // Disable the button to prevent subsequent calls, since heatmapLayer is now nil.
        button.isEnabled = false
    }
    
    // MARK: fadeIn
    @objc
    private func fadeIn() {
        if heatmapLayer.fadeIn {
            heatmapLayer.fadeIn = false
        } else {
            heatmapLayer.fadeIn = true
        }
        print("heatmap faded in \(heatmapLayer.fadeIn)")
    }
    
    // MARK: radius
    @objc
    private func increaseRadius() {
        heatmapLayer.radius += UInt(20.0)
        print("radius: \(heatmapLayer.radius)")
    }
    
    // MARK: opacity
    @objc
    private func increaseOpacity() {
        heatmapLayer.opacity += 0.05
        print("opacity: \(heatmapLayer.opacity)")
    }
    
    // MARK: gradient
    @objc
    private func changeGradient() {
        heatmapLayer.gradient = GMUGradient(colors: [UIColor.yellow, UIColor.green], startPoints: [0.1, 0.5] as [NSNumber], colorMapSize: 128)
        print("gradient: \(heatmapLayer.gradient.description)")
    }
    
    // MARK: tileSize
    @objc
    private func changeTileSize() {
        heatmapLayer.tileSize += 20
        print("tile size: \(heatmapLayer.tileSize)")
    }
    
    // MARK: zIndex
    @objc
    private func changeZindex() {
        heatmapLayer.zIndex += 10
        print("z index: \(heatmapLayer.zIndex)")
    }
    
    // MARK: maximumZoneIntensity
    @objc
    private func changeMaxZoneIntensity() {
        heatmapLayer.minimumZoomIntensity += UInt(0.1)
        print("Maximum zone intensity: \(heatmapLayer.minimumZoomIntensity)")
    }
    
    // MARK: minimumZoneIntensity
    @objc
    private func changeMinZoneIntensity() {
        heatmapLayer.minimumZoomIntensity += UInt(0.1)
        print("Minimum zone intensity: \(heatmapLayer.minimumZoomIntensity)")
    }
    
    private func addButtons() {
        addButton(title: "Remove",     xOrig: 0,     yOrig: 5,  width: 100, height: 35, objcFunc: #selector(removeHeatmap))
        addButton(title: "fadeIn",     xOrig: 110,   yOrig: 5,  width: 100, height: 35, objcFunc: #selector(fadeIn))
        addButton(title: "Radius",     xOrig: 220,   yOrig: 5,  width: 100, height: 35, objcFunc: #selector(increaseRadius))
        addButton(title: "Opacity",    xOrig: 330,   yOrig: 5,  width: 100, height: 35, objcFunc: #selector(increaseOpacity))
        addButton(title: "Gradient",   xOrig: 0,     yOrig: 45, width: 100, height: 35, objcFunc: #selector(changeGradient))
        addButton(title: "TileSize",   xOrig: 110,   yOrig: 45, width: 100, height: 35, objcFunc: #selector(changeTileSize))
        addButton(title: "zIndex",     xOrig: 220,   yOrig: 45, width: 100, height: 35, objcFunc: #selector(changeZindex))
        addButton(title: "MaxZoneInt", xOrig: 330,   yOrig: 45, width: 120, height: 35, objcFunc: #selector(changeMaxZoneIntensity))
        addButton(title: "MinZoneInt", xOrig: 460,   yOrig: 45, width: 120, height: 35, objcFunc: #selector(changeMinZoneIntensity))
    }
    
    private func addButton(title: String, xOrig:Int, yOrig: Int, width: Int, height: Int, objcFunc: Selector) {
        button = UIButton(frame: CGRect(x: xOrig, y: yOrig, width: width, height: height))
        button.backgroundColor = .blue
        button.setTitle(title, for: .normal)
        button.addTarget(self, action: objcFunc, for: .touchUpInside)
        self.mapView.addSubview(button)
    }
}
