//
//  ViewController.swift
//  HeatmapUIKit
//
//  Created by Kamber Vogli on 15.03.22.
//

import UIKit
import AzureMapsControl

class ViewController: UIViewController, AzureMapDelegate {
    private var azureMap: MapControl!
    private var heatmapLayer: HeatMapLayer!
    private var radius = 20.0, opacity = 0.8, intensity = 1.0, weight = 1.0
    
    private var radiusLabel  = UILabel(),
                opacityLabel = UILabel(),
                intensityLabel = UILabel(),
                weightLabel    = UILabel()
    
    func azureMapCameraIsMoving(_ map: AzureMap) {
        // TODO use NSExpression with Interpolation expression since i cannot access zoomLevel
        // use interpolation for making the heatmap smaller while zooming out :(
        
        //Option 2: PLay with minZoom, maxZoom and extra heatmap layers which we can add on the azureMap based on the min/maxZoom and visibility options
        
        // Option 3: Request feature as Github issue
        // github issue for subscription key being wrong, google maps for radius .clearTileCache() not working
    }
    
    override func loadView() {
        super.loadView()
        azureMap = MapControl.init(frame: CGRect(x: 0, y: 0, width: 500, height: 800),
                                   options: [
                                    CameraOption.center(lat: 48.2692083204, lng: 11.6690079838),
                                    CameraOption.zoom(9),
                                    // StyleOption.style(.grayscaleDark)
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
                        1: 10,
                        
                        // Between zoom level 1 and 19, exponentially scale the radius from 2 points to 2 * 2^(maxZoom - minZoom) points.
                        22: 100
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
        
        addControl(title: "Radius",  valueLabel: &radiusLabel,  distFromTop: 50,  min: 1, max: 50, #selector(updateRadius(_:)))
        addControl(title: "Opacity", valueLabel: &opacityLabel, distFromTop: 100, min: 0, max: 1,  #selector(updateOpacity(_:)))
        addControl(title: "Intencity", valueLabel: &intensityLabel, distFromTop: 150, min: 0, max: 5,  #selector(updateIntensity(_:)))
        addControl(title: "Weight", valueLabel: &weightLabel, distFromTop: 200, min: 0, max: 5,  #selector(updateWeight(_:)))
    }
    
    private func addControl(title: String, valueLabel: inout UILabel, distFromTop: Int, min: Float, max: Float ,_ objcFunc: Selector) {
        let titleLabel = UILabel(frame: CGRect(x: 510, y: distFromTop, width: 80, height: 20))
        titleLabel.text = title
        valueLabel = UILabel(frame: CGRect(x: 850, y: distFromTop, width: 100, height: 20))
        valueLabel.text = "0.0"
        addUISlider(xOrig: 600, yOrig: distFromTop, width: 200, height: 20, minVal: min, maxVal: max, objcFunc: objcFunc)
        view.addSubview(titleLabel)
        view.addSubview(valueLabel)
    }
    
    private func addUISlider(xOrig: Int, yOrig: Int, width: Int, height: Int, minVal: Float, maxVal: Float, objcFunc: Selector) {
        let slider = UISlider(frame: CGRect(x: xOrig, y: yOrig, width: width, height: height))
        slider.minimumValue = minVal
        slider.maximumValue = maxVal
        slider.isContinuous = true
        slider.addTarget(self, action: objcFunc, for: .valueChanged)
        view.addSubview(slider)
    }
    
    @objc
    func updateRadius(_ sender: UISlider!) {
        let step:Float = 0.05
        let stepValue = round(sender.value / step) * step
        sender.value = stepValue
        radius = Double(stepValue)
        radiusLabel.text = String(radius)
        heatmapLayer.setOptions([.heatmapRadius(radius)])
    }
    
    @objc
    func updateOpacity(_ sender: UISlider!) {
        let step:Float = 0.01
        let stepValue = round(sender.value / step) * step
        sender.value = stepValue
        opacity = Double(stepValue)
        opacityLabel.text = String(opacity)
        heatmapLayer.setOptions([.heatmapOpacity(opacity)])
    }
    
    @objc
    func updateIntensity(_ sender: UISlider!) {
        let step:Float = 0.01
        let stepValue = round(sender.value / step) * step
        sender.value = stepValue
        intensity = Double(stepValue)
        intensityLabel.text = String(intensity)
        heatmapLayer.setOptions([.heatmapIntensity(intensity)])
    }
    
    @objc
    func updateWeight(_ sender: UISlider!) {
        let step:Float = 0.01
        let stepValue = round(sender.value / step) * step
        sender.value = stepValue
        weight = Double(stepValue)
        weightLabel.text = String(weight)
        heatmapLayer.setOptions([.heatmapWeight(weight)])
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

