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
    private var clusterManager: GMUClusterManager!
    
    private var colorMSLabel = UILabel(),
                opacityLabel = UILabel(),
                radiusLabel  = UILabel(),
                minZILabel   = UILabel(),
                maxZILabel   = UILabel(),
                zoomLvlLabel = UILabel(),
                blueLabel    = UILabel(),
                cyanLabel    = UILabel(),
                greenLabel   = UILabel(),
                yellowLabel  = UILabel(),
                redLabel     = UILabel(),
                
                markersButton = UIButton(), clusterButton = UIButton(),
                show = true, cluster = true
    
    private var gradientColors = [UIColor(red: 0, green: 0, blue: 128/255, alpha: 1),
                                  UIColor.cyan,
                                  UIColor.green,
                                  UIColor.yellow,
                                  UIColor.red]
    private var gradientStartPoints = [0.05, 0.1, 0.15, 0.26,  0.5] as [NSNumber]
    
    private var blue: Float = 0.05,
                cyan: Float = 0.10,
                green: Float = 0.15,
                yellow: Float = 0.26,
                red: Float = 0.5,
                
                opacity: Float = 0.8,
                radius = UInt(87),
                colorMapSize = UInt(512) // default = 512
    
    private var apEG: [GMUWeightedLatLng]!,
                ap1OG: [GMUWeightedLatLng]!,
                ap2OG: [GMUWeightedLatLng]!,
                ap3OG: [GMUWeightedLatLng]!,
                allAPs: [GMUWeightedLatLng]!,
                allMarkers: [GMSMarker]!
    
    private var zoomLevel: Float = 16.0
    private var heatmapLayers: [GMUHeatmapTileLayer?]!
    
    private var interpolation = Interpolation()
    
    private var backgroundView = UIView(frame: CGRect(x: 5, y: 5, width: 350, height: 500))
    
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
        
        let zoomTitle = UILabel(frame: CGRect(x: 5, y: 480, width: 150, height: 15))
        zoomTitle.text = "Zoom Lvl:"
        backgroundView.addSubview(zoomTitle)
        zoomLvlLabel.frame = CGRect(x: 160, y: 480, width: 80, height: 15)
        zoomLvlLabel.text = "0.0"
        backgroundView.addSubview(zoomLvlLabel)
    }
    
    private func addMarker(coord: CLLocationCoordinate2D) {
        let marker = GMSMarker(position: coord)
        marker.title = "Access Point"
        marker.map = mapView
    }
    
    override func viewDidLoad() {
        super.viewDidLoad()
//        apEG = addHeatmap()
//        setupMultipleLayers()
        initHeatmapLayer(radius: radius, opacity: opacity, colorMapSize: colorMapSize)
        print("zoom level is: \(mapView.camera.zoom)")
        
        // show or not show markers in the map
        markersButton.frame = CGRect(x: 10, y: 720, width: 200, height: 40)
        markersButton.backgroundColor = .black
        markersButton.setTitle("Show markers", for: .normal)
        markersButton.addTarget(self, action: #selector(showMarkers), for: .touchUpInside)
        mapView.addSubview(markersButton)
        
        // Setup the cluster manager
        let iconGenerator = GMUDefaultClusterIconGenerator()
        let algo = GMUNonHierarchicalDistanceBasedAlgorithm()
        let renderer = GMUDefaultClusterRenderer(mapView: mapView,
                                                 clusterIconGenerator: iconGenerator)
        clusterManager = GMUClusterManager(map: mapView, algorithm: algo, renderer: renderer)
        clusterManager.setMapDelegate(self)
        
        clusterButton.frame = CGRect(x: 10, y: 675, width: 200, height: 40)
        clusterButton.backgroundColor = .black
        clusterButton.setTitle("Cluster markers", for: .normal)
        clusterButton.addTarget(self, action: #selector(clusterMarkers), for: .touchUpInside)
        mapView.addSubview(clusterButton)
    }
    
    @objc
    func clusterMarkers() {
        if cluster {
            for marker in allMarkers {
                clusterManager.add(marker)
            }
            clusterManager.cluster()
            cluster = false
        } else {
            for marker in allMarkers {
                clusterManager.remove(marker)
            }
            clusterManager.cluster()
            cluster = true
        }
    }
    
    @objc
    func showMarkers() {
        if show {
            for marker in allMarkers {
                marker.map = mapView
            }
            show = false
        } else {
            for marker in allMarkers {
                marker.map = nil
            }
            show = true
        }
    }
    
    func mapView(_ mapView: GMSMapView, didTap marker: GMSMarker) -> Bool {
        mapView.animate(toLocation: marker.position)
        
        if marker.userData is GMUCluster {
            NSLog("Did tap cluster")
            return true
        }
        
        NSLog("Did tap a normal marker")
        return false
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
//            hmLayer.weightedData = apEG
//            hmLayer.zIndex = Int32(i)
//            hmLayer.map = mapView
            heatmapLayers[i - 1] = hmLayer
        }
    }
    
    func mapView(_ mapView: GMSMapView, didChange position: GMSCameraPosition) {
        let currentZoomLevel = mapView.camera.zoom
        zoomLvlLabel.text = String(currentZoomLevel)
        if abs(currentZoomLevel - zoomLevel) > 0.5 {
            zoomLevel = currentZoomLevel //update
//            print("Update zl: \(zoomLevel)")
            let newRadius = lagrange(points: interpolation.radiusPoints, x: Double(currentZoomLevel))
            let newOpacity = lagrange(points: interpolation.opacityPoints, x: Double(currentZoomLevel))
            heatmapLayer.radius = UInt(newRadius)
            heatmapLayer.opacity = Float(newOpacity)
            heatmapLayer.map = mapView
            
            print("new rad: \(newRadius)")
            print("new opa: \(newOpacity)")
//            heatmapLayer = heatmapLayers[Int(zoomLevel)]
//            let lvl = Int(zoomLevel)
//            let hml = heatmapLayers[lvl]
//            hml?.map = mapView
//
//
//            //remove the other heatmaps
//            for i in 0...22 {
//                if i != lvl {
//                    let otherHML = heatmapLayers[i]
//                    otherHML?.map = nil
//                }
//            }
//            print("zoom lvl: \(Int(zoomLevel)) & radius: \(hml?.radius)")
//            heatmapLayer.clearTileCache()
//            initHeatmapLayer(radius: radius, opacity: opacity, colorMapSize: colorMapSize)
        }
    }
    
    
    // TODO display different heatmap based on selected level
    func didChangeActiveLevel(_ level: GMSIndoorLevel?) {
//        print("Showing heatmap layer for floor: \(level?.name)")
//        if level?.name == "EG" {
//            heatmapLayer.weightedData = apEG
//        } else if level?.name == "1OG" {
//            heatmapLayer.weightedData = ap1OG
//        } else if level?.name == "2OG" {
//            heatmapLayer.weightedData = ap2OG
//        } else if level?.name == "3OG" {
//            heatmapLayer.weightedData = ap3OG
//        }
//        heatmapLayer.clearTileCache()
//        heatmapLayer.map = heatmapLayer.map
    }
    
    private func initHeatmapLayer(radius: UInt, opacity: Float, colorMapSize: UInt) {
        heatmapLayer = GMUHeatmapTileLayer()
        heatmapLayer.radius = radius // old value = 300 , 50,
        heatmapLayer.opacity = opacity // default = 0.8 or 0.75
        heatmapLayer.fadeIn = true
        heatmapLayer.gradient = GMUGradient(colors: gradientColors,
                                            startPoints: gradientStartPoints,
                                            colorMapSize: colorMapSize) // 256

        addHeatmap()
        // Set the heatmap to the mapview.
        heatmapLayer.map = mapView
    }
    
    // Parse JSON data and add it to the heatmap layer.
    func addHeatmap() {
        var listEG = [GMUWeightedLatLng]()
        var list1OG = [GMUWeightedLatLng]()
        var list2OG = [GMUWeightedLatLng]()
        var list3OG = [GMUWeightedLatLng]()
        allAPs = [GMUWeightedLatLng]()
        allMarkers = [GMSMarker]()
        do {
            // Get the data: latitude/longitude positions of police stations.
            if let path = Bundle.main.url(forResource: "ap-2", withExtension: "json") {
                let data = try Data(contentsOf: path)
                let json = try JSONSerialization.jsonObject(with: data, options: [])
                if let object = json as? [[String: Any]] {
                    for item in object {
                        let lat = item["Lat"] as? Double
                        let lng = item["Long"] as? Double
                        var intensity = item["Intensity"] as? Double ?? 1.0
                        let floor = item["Floor"] as? String
                        
                        if intensity != 1.0 {
                            intensity *= 1000.0
                        }
                        
//                        if intensity < 5 || intensity > 20 {
//                            print("intensity is : \(Float(intensity))")
//                        }

                        let coords = GMUWeightedLatLng(
                            coordinate: CLLocationCoordinate2DMake(
                                lat!, lng!
                            ),
                            intensity: Float(intensity)
                        )
                        
                        if floor == "0" {
                            listEG.append(coords)
                        } else if floor == "1" {
                            list1OG.append(coords)
                        } else if floor == "2" {
                            list2OG.append(coords)
                        } else if floor == "3" {
                            list3OG.append(coords)
                        }
                        
                        allAPs.append(coords)
                        let marker = GMSMarker(position: CLLocationCoordinate2D(latitude: lat!, longitude: lng!))
                        allMarkers.append(marker)
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
        
        apEG  = listEG
        ap1OG = list1OG
        ap2OG = list2OG
        ap3OG = list3OG
        
        // Add the latlngs to the heatmap layer.
      heatmapLayer.weightedData = allAPs
    }
    
    private func addSliders() {
        let fromLeft = 5, width = 170, height = 15
        // MARK: colorMapSize
        addUISlider(title: "colorMS", xOrig: fromLeft, yOrig: 20, width: width, height: height, minVal: 2, maxVal: 512.0, objcFunc: #selector(self.colorMapSizeSliderValueDidChange(_:)), &colorMSLabel, Float(colorMapSize))
        
        // MARK: Opacity
        addUISlider(title: "Opacity", xOrig: fromLeft, yOrig: 80, width: width, height: height, minVal: 0.05, maxVal: 10.0, objcFunc: #selector(self.opacitySliderValueDidChange(_:)), &opacityLabel, opacity)
        
        // MARK: Radius
        addUISlider(title: "Radius", xOrig: fromLeft, yOrig: 120, width: width, height: height, minVal: 1, maxVal: 500, objcFunc: #selector(self.radiusSliderValueDidChange(_:)), &radiusLabel, Float(radius))
        
        // MARK: Minimum Zone Intensity
        addUISlider(title: "MinZI", xOrig: fromLeft, yOrig: 180, width: width, height: height, minVal: 0, maxVal: 4, objcFunc: #selector(self.minZoneIntSliderValueDidChange(_:)), &minZILabel, 0.0)
        
        // MARK: Maximum Zone Intensity
        addUISlider(title: "MaxZI", xOrig: fromLeft, yOrig: 220, width: width, height: height, minVal: 10, maxVal: 40, objcFunc: #selector(self.maxZoneIntSliderValueDidChange(_:)), &maxZILabel, 38.0)
        
        
        
        // MARK: Gradient Blue
        addUISlider(title: "Blue", xOrig: fromLeft, yOrig: 280, width: width, height: height, minVal: 0.0, maxVal: 0.1, objcFunc: #selector(self.blueSliderValueDidChange(_:)), &blueLabel, blue)
        
        // MARK: Gradient Cyan
        addUISlider(title: "Cyan", xOrig: fromLeft, yOrig: 320, width: width, height: height, minVal: 0.1, maxVal: 0.25, objcFunc: #selector(self.cyanSliderValueDidChange(_:)), &cyanLabel, cyan)
        
        // MARK: Gradient Green
        addUISlider(title: "Green", xOrig: fromLeft, yOrig: 360, width: width, height: height, minVal: 0.25, maxVal: 0.5, objcFunc: #selector(self.greenSliderValueDidChange(_:)), &greenLabel, green)
        
        // MARK: Gradient Yellow
        addUISlider(title: "Yellow", xOrig: fromLeft, yOrig: 400, width: width, height: height, minVal: 0.5, maxVal: 0.75, objcFunc: #selector(self.yellowSliderValueDidChange(_:)), &yellowLabel, yellow)
        
        // MARK: Gradient Red
        addUISlider(title: "Red", xOrig: fromLeft, yOrig: 440, width: width, height: height, minVal: 0.75, maxVal: 1.0, objcFunc: #selector(self.redSliderValueDidChange(_:)), &redLabel, red)
    }
    
    private func addUISlider(title: String,
                             xOrig: Int,
                             yOrig: Int,
                             width: Int,
                             height: Int,
                             minVal: Float,
                             maxVal: Float,
                             objcFunc: Selector,
                             _ valueLabel: inout UILabel,
                             _ initialValue: Float) {
        
        let slider = UISlider(frame:CGRect(x: xOrig + 90, y: yOrig, width: width, height: height))
        slider.minimumValue = minVal
        slider.maximumValue = maxVal
        slider.isContinuous = true
        slider.addTarget(self, action: objcFunc, for: .valueChanged)
        slider.value = initialValue
        backgroundView.addSubview(slider)
        
        let titleLabel = UILabel(frame: CGRect(x: xOrig, y: yOrig, width: 70, height: height))
        titleLabel.text = title
        backgroundView.addSubview(titleLabel)
        
        valueLabel.frame = CGRect(x: xOrig + 270, y: yOrig, width: width, height: height)
        valueLabel.text = "\(initialValue)"
        backgroundView.addSubview(valueLabel)
    }
    
    @objc
    func colorMapSizeSliderValueDidChange(_ sender: UISlider!) {
        let step: Float = 0.5
        let stepValue = sender.value / step * step
        sender.value = stepValue
        colorMapSize = UInt(stepValue)
        colorMSLabel.text = String(colorMapSize)
        heatmapLayer.gradient = GMUGradient(colors: gradientColors,
                                            startPoints: gradientStartPoints,
                                            colorMapSize: colorMapSize)
        heatmapLayer.map = heatmapLayer.map
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
        let step = 1
        let stepValue = Int(sender.value) / step * step
        sender.value = Float(stepValue)
        radius = UInt(stepValue)
//        heatmapLayer.map = nil
        heatmapLayer.radius = radius
        heatmapLayer.map = mapView
        radiusLabel.text = String(radius)
        heatmapLayer.clearTileCache()
    }
    
    @objc
    func minZoneIntSliderValueDidChange(_ sender: UISlider!) {
        let step: Float = 0.5
        let stepValue = sender.value / step * step
        sender.value = stepValue
        heatmapLayer.minimumZoomIntensity = UInt(stepValue)
        heatmapLayer.map = heatmapLayer.map
        minZILabel.text = String(heatmapLayer.minimumZoomIntensity)
        heatmapLayer.clearTileCache()
    }
    
    @objc
    func maxZoneIntSliderValueDidChange(_ sender: UISlider!) {
        let step = 1
        let stepValue = Int(sender.value) / step * step
        sender.value = Float(stepValue)
        heatmapLayer.maximumZoomIntensity = UInt(stepValue)
        heatmapLayer.map = heatmapLayer.map
        maxZILabel.text = String(heatmapLayer.maximumZoomIntensity)
        heatmapLayer.clearTileCache()
    }
    
    @objc
    func blueSliderValueDidChange(_ sender: UISlider!) {
        let step:Float = 0.001
        let stepValue = sender.value / step * step
        sender.value = stepValue
        blue = stepValue
        gradientStartPoints = [blue, cyan, green, yellow, red] as [NSNumber]
        heatmapLayer.gradient = GMUGradient(colors: gradientColors,
                                            startPoints: gradientStartPoints,
                                            colorMapSize: colorMapSize)
        heatmapLayer.map = heatmapLayer.map
        blueLabel.text = String(blue)
        heatmapLayer.clearTileCache()
    }
    
    @objc
    func cyanSliderValueDidChange(_ sender: UISlider!) {
        let step:Float = 0.001
        let stepValue = sender.value / step * step
        sender.value = stepValue
        cyan = stepValue
        gradientStartPoints = [blue, cyan, green, yellow, red] as [NSNumber]
        heatmapLayer.gradient = GMUGradient(colors: gradientColors,
                                            startPoints: gradientStartPoints,
                                            colorMapSize: colorMapSize)
        heatmapLayer.map = heatmapLayer.map
        cyanLabel.text = String(cyan)
        heatmapLayer.clearTileCache()
    }
    
    @objc
    func greenSliderValueDidChange(_ sender: UISlider!) {
        let step:Float = 0.001
        let stepValue = sender.value / step * step
        sender.value = stepValue
        green = stepValue
        gradientStartPoints = [blue, cyan, green, yellow, red] as [NSNumber]
        heatmapLayer.gradient = GMUGradient(colors: gradientColors,
                                            startPoints: gradientStartPoints,
                                            colorMapSize: colorMapSize)
        heatmapLayer.map = heatmapLayer.map
        greenLabel.text = String(green)
        heatmapLayer.clearTileCache()
    }
    
    @objc
    func yellowSliderValueDidChange(_ sender: UISlider!) {
        let step:Float = 0.001
        let stepValue = sender.value / step * step
        sender.value = stepValue
        yellow = stepValue
        gradientStartPoints = [blue, cyan, green, yellow, red] as [NSNumber]
        heatmapLayer.gradient = GMUGradient(colors: gradientColors,
                                            startPoints: gradientStartPoints,
                                            colorMapSize: colorMapSize)
        heatmapLayer.map = heatmapLayer.map
        yellowLabel.text = String(yellow)
        heatmapLayer.clearTileCache()
    }
    
    @objc
    func redSliderValueDidChange(_ sender: UISlider!) {
        let step:Float = 0.001
        let stepValue = sender.value / step * step
        sender.value = stepValue
        red = stepValue
        gradientStartPoints = [blue, cyan, green, yellow, red] as [NSNumber]
        heatmapLayer.gradient = GMUGradient(colors: gradientColors,
                                            startPoints: gradientStartPoints,
                                            colorMapSize: colorMapSize)
        heatmapLayer.map = heatmapLayer.map
        redLabel.text = String(red)
        heatmapLayer.clearTileCache()
    }
}
