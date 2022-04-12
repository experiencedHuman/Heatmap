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
    
    override func loadView() {
        //    let camera = GMSCameraPosition.camera(withLatitude: -37.848, longitude: 145.001, zoom: 10)
        let camera = GMSCameraPosition.camera(withLatitude: 48.14957600438307, longitude: 11.567179933190348, zoom: 16)
        mapView = GMSMapView.map(withFrame: CGRect.zero, camera: camera)
        mapView.delegate = self
        self.view = mapView
        removeButton()
        fadeInButton()
        opacityButton()
        gradientButton()
        tileSizeButton()
        zIndexButton()
        maxZoneIntensityButton()
        minZoneIntensityButton()
    }
    
    override func viewDidLoad() {
        // Set heatmap options.
        heatmapLayer = GMUHeatmapTileLayer()
        heatmapLayer.radius = 300
        heatmapLayer.opacity = 0.8
        heatmapLayer.gradient = GMUGradient(colors: gradientColors,
                                            startPoints: gradientStartPoints,
                                            colorMapSize: 256)
        addHeatmap()
        
        // Set the heatmap to the mapview.
        heatmapLayer.map = mapView
    }
    
    // Parse JSON data and add it to the heatmap layer.
    func addHeatmap()  {
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
        heatmapLayer.weightedData = list
    }
    
    @objc
    func removeHeatmap() {
        heatmapLayer.map = nil
        heatmapLayer = nil
        // Disable the button to prevent subsequent calls, since heatmapLayer is now nil.
        button.isEnabled = false
    }
    
    func mapView(_ mapView: GMSMapView, didTapAt coordinate: CLLocationCoordinate2D) {
        print("You tapped at \(coordinate.latitude), \(coordinate.longitude)")
    }
    
    // Add a button to the view.
    func removeButton() {
        // A button to test removing the heatmap.
        button = UIButton(frame: CGRect(x: 0, y: 0, width: 100, height: 35))
        button.backgroundColor = .blue
//        button.alpha = 0.5
        button.setTitle("Remove map", for: .normal)
        button.addTarget(self, action: #selector(removeHeatmap), for: .touchUpInside)
        self.mapView.addSubview(button)

    }
    
    /*
     TODO properties & methods of mapView
     
     properties of heatmapLayer
     heatmapLayer.map
     heatmapLayer.fadeIn //bool. specifies whether the tiles should fade in. default yes.
     heatmapLayer.radius
     heatmapLayer.opacity
     heatmapLayer.gradient
     heatmapLayer.tileSize
     heatmapLayer.zIndex
     heatmapLayer.weightedData
     heatmapLayer.maximumZoomIntensity
     heatmapLayer.minimumZoomIntensity
     
     methods
     heatmapLayer.clearTileCache()
     heatmapLayer.requestTileFor(x: <#T##UInt#>, y: <#T##UInt#>, zoom: <#T##UInt#>, receiver: <#T##GMSTileReceiver#>)
     heatmapLayer.tileFor(x: <#T##UInt#>, y: <#T##UInt#>, zoom: <#T##UInt#>)
     */
    
    // MARK: fadeIn
    @objc
    func fadeIn() {
        if heatmapLayer.fadeIn {
            heatmapLayer.fadeIn = false
        } else {
            heatmapLayer.fadeIn = true
        }
        print("heatmap faded in \(heatmapLayer.fadeIn)")
    }
    
    func fadeInButton() {
        button = UIButton(frame: CGRect(x: 110, y: 0, width: 100, height: 35))
        button.backgroundColor = .blue
        button.setTitle("fade in", for: .normal)
        button.addTarget(self, action: #selector(fadeIn), for: .touchUpInside)
        self.mapView.addSubview(button)
    }
    
    // MARK: radius
    @objc
    func increaseRadius() {
        heatmapLayer.radius += UInt(20.0)
        print("radius: \(heatmapLayer.radius)")
    }
    
    func radiusButton() {
        button = UIButton(frame: CGRect(x: 220, y: 0, width: 100, height: 35))
        button.backgroundColor = .blue
        button.setTitle("rad +20", for: .normal)
        button.addTarget(self, action: #selector(increaseRadius), for: .touchUpInside)
        self.mapView.addSubview(button)
    }
    
    // MARK: opacity
    @objc
    func increaseOpacity() {
        heatmapLayer.opacity += 0.05
        print("opacity: \(heatmapLayer.opacity)")
    }
    
    func opacityButton() {
        button = UIButton(frame: CGRect(x: 330, y: 5, width: 100, height: 35))
        button.backgroundColor = .blue
        button.setTitle("opacity +0.05", for: .normal)
        button.addTarget(self, action: #selector(increaseOpacity), for: .touchUpInside)
        self.mapView.addSubview(button)
    }
    
    // MARK: gradient
    @objc
    func changeGradient() {
        heatmapLayer.gradient = GMUGradient(colors: [UIColor.yellow, UIColor.green], startPoints: [0.1, 0.5] as [NSNumber], colorMapSize: 128)
        print("gradient: \(heatmapLayer.gradient.description)")
    }
    
    func gradientButton() {
        button = UIButton(frame: CGRect(x: 0, y: 45, width: 100, height: 35))
        button.backgroundColor = .blue
        button.setTitle("gradient", for: .normal)
        button.addTarget(self, action: #selector(changeGradient), for: .touchUpInside)
        self.mapView.addSubview(button)
    }
    
    // MARK: tileSize
    @objc
    func changeTileSize() {
        heatmapLayer.tileSize += 20
        print("tile size: \(heatmapLayer.tileSize)")
    }
    
    func tileSizeButton() {
        button = UIButton(frame: CGRect(x: 110, y: 45, width: 100, height: 35))
        button.backgroundColor = .blue
        button.setTitle("tileSize", for: .normal)
        button.addTarget(self, action: #selector(changeTileSize), for: .touchUpInside)
        self.mapView.addSubview(button)
    }
    
    // MARK: zIndex
    @objc
    func changeZindex() {
        heatmapLayer.zIndex += 10
        print("z index: \(heatmapLayer.zIndex)")
    }
    
    func zIndexButton() {
        button = UIButton(frame: CGRect(x: 220, y: 45, width: 100, height: 35))
        button.backgroundColor = .blue
        button.setTitle("fade in", for: .normal)
        button.addTarget(self, action: #selector(changeZindex), for: .touchUpInside)
        self.mapView.addSubview(button)
    }
    
    // MARK: maximumZoneIntensity
    @objc
    func changeMaxZoneIntensity() {
        heatmapLayer.minimumZoomIntensity += UInt(0.1)
        print("Maximum zone intensity: \(heatmapLayer.minimumZoomIntensity)")
    }
    
    func maxZoneIntensityButton() {
        button = UIButton(frame: CGRect(x: 330, y: 45, width: 100, height: 35))
        button.backgroundColor = .blue
        button.setTitle("MaxZone", for: .normal)
        button.addTarget(self, action: #selector(changeMaxZoneIntensity), for: .touchUpInside)
        self.mapView.addSubview(button)
    }
    
    // MARK: minimumZoneIntensity
    @objc
    func changeMinZoneIntensity() {
        heatmapLayer.minimumZoomIntensity += UInt(0.1)
        print("Minimum zone intensity: \(heatmapLayer.minimumZoomIntensity)")
    }
    
    func minZoneIntensityButton() {
        button = UIButton(frame: CGRect(x: 440, y: 45, width: 100, height: 35))
        button.backgroundColor = .blue
        button.setTitle("MinZone", for: .normal)
        button.addTarget(self, action: #selector(changeMinZoneIntensity), for: .touchUpInside)
        self.mapView.addSubview(button)
    }
}
