//
//  ViewController.swift
//  HeatmapUIKit
//
//  Created by Kamber Vogli on 15.03.22.
//

import UIKit
import AzureMapsControl

class ViewController: UIViewController {
    private var azureMap: MapControl!
    
    override func loadView() {
        super.loadView()
        azureMap = MapControl.init(frame: CGRect(x: 0, y: 0, width: 500, height: 800),
                                options: [
                                    CameraOption.center(lat: 48.2692083204, lng: 11.6690079838),
                                    CameraOption.zoom(9),
//                                    StyleOption.style(.grayscaleDark)
                                ])
        self.view.addSubview(azureMap)
    }
    
    override func viewDidLoad() {
        super.viewDidLoad()
        // Do any additional setup after loading the view.
        let source = DataSource()
        let locations = readCoordsFromJSON(file: "accessPoints")
        let pointCollection = PointCollection(locations)
        source.add(geometry: pointCollection)
        
        azureMap.onReady { map in
            map.sources.add(source)
            let layer = HeatMapLayer(source: source, options: [
                .heatmapRadius(10),
                .heatmapOpacity(0.8)
            ])
            map.layers.insertLayer(layer, below: "labels")
        }
        
    }
    
    private func readCoordsFromJSON(file filename: String) -> [CLLocationCoordinate2D] {
        var coordinates: [CLLocationCoordinate2D] = []
        do {
            if let path = Bundle.main.url(forResource: filename, withExtension: "json") {
                let data = try Data(contentsOf: path)
                let json = try JSONSerialization.jsonObject(with: data, options: [])
                if let object = json as? [[String: Any]] {
                    for item in object {
                        let lat  = item["Latitude"] as? Double ?? 0.0
                        let long = item["Longitude"] as? Double ?? 0.0
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

