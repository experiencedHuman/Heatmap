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
    private var heatmapLayer: HeatMapLayer!
    private var radius: Double = 20
    
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
            self.heatmapLayer = HeatMapLayer(source: source, options: [
                .heatmapRadius(from: NSExpression(
                    forAZMInterpolating: .zoomLevelAZMVariable,
                    curveType: ExpressionInterpolationMode.exponential,
                    parameters: NSExpression(forConstantValue: 2),
                    stops: NSExpression(forConstantValue: [

                        // For zoom level 1 set the radius to 2 points.
                        1: 5,

                        // Between zoom level 1 and 19, exponentially scale the radius from 2 points to 2 * 2^(maxZoom - minZoom) points.
                        22: 500
                    ])
                )),
                .heatmapOpacity(0.8),
//                .heatmapColor(
//                    from: // Stepped color expression
//                            NSExpression(
//                                forAZMStepping: .heatmapDensityAZMVariable,
//                                from: NSExpression(forConstantValue: UIColor.clear),
//                                stops: NSExpression(forConstantValue: [
//                                    0.01: UIColor(red: 0, green: 0, blue: 128/255, alpha: 1),
//                                    0.25: UIColor.cyan,
//                                    0.5: UIColor.green,
//                                    0.75: UIColor.yellow,
//                                    1: UIColor.red
//                                ])
//                            )
//                ),
                .minZoom(5.0)
            ])
            
            map.layers.insertLayer(self.heatmapLayer, below: "labels")
        }
        
        addUISlider(xOrig: 550, yOrig: 50, width: 150, height: 20, minVal: 1, maxVal: 1000, objcFunc: #selector(radiusSliderValueDidChange(_:)))
    }

    private func addUISlider(xOrig: Int, yOrig: Int, width: Int, height: Int, minVal: Float, maxVal: Float, objcFunc: Selector) {
        let slider = UISlider(frame:CGRect(x: xOrig, y: yOrig, width: width, height: height))
        slider.minimumValue = minVal
        slider.maximumValue = maxVal
        slider.isContinuous = true
        slider.addTarget(self, action: objcFunc, for: .valueChanged)
        view.addSubview(slider)
    }

    @objc
    func radiusSliderValueDidChange(_ sender:UISlider!) {
        let step = 1
        let stepValue = Int(sender.value) / step * step
        sender.value = Float(stepValue)
        radius = Double(stepValue)
        heatmapLayer.setOptions([.heatmapRadius(radius)])
        print("Radius: \(radius)")
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

