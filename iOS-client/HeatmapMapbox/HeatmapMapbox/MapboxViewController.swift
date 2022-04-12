//
//  MapboxViewController.swift
//  HeatmapMapbox
//
//  Created by Kamber Vogli on 15.03.22.
//

import Foundation


import Mapbox
import SwiftUI
import UIKit

struct MapboxView: UIViewControllerRepresentable {
    
    func makeUIViewController(context: Context) -> MapboxViewController {
        MapboxViewController()
    }
    
    func updateUIViewController(_ uiViewController: MapboxViewController, context: Context) {
    }
}

class MapboxViewController: UIViewController, MGLMapViewDelegate {
    private var mapView: MGLMapView!
    private var heatmapLayer: MGLHeatmapStyleLayer!
    
    override func loadView() {
        super.loadView()
        
        // create and add a map view
        mapView = MGLMapView(frame: view.bounds, styleURL: MGLStyle.lightStyleURL)
        mapView.autoresizingMask = [.flexibleWidth, .flexibleHeight]
        mapView.setCenter(CLLocationCoordinate2D(latitude: 48.2692083204, longitude: 11.6690079838), zoomLevel: 9, animated: false)
        mapView.autoresizingMask = [.flexibleHeight, .flexibleWidth]
        mapView.delegate = self
        mapView.tintColor = .lightGray
//        view.addSubview(mapView)
         self.view = mapView
    }
    
    override func viewDidLoad() {
        super.viewDidLoad()
    }
    

    func mapView(_ mapView: MGLMapView, didFinishLoading style: MGLStyle) {
        let identifier = "accesspoints"
        var coordinates = readCoordsFromJSON(file: "accessPoints")
        let collection = MGLPointCollection(coordinates: &coordinates, count: UInt(coordinates.count))
        let source = MGLShapeSource(identifier: identifier, shape: collection, options: nil)
        style.addSource(source)
        
        // Create a heatmap layer.
        heatmapLayer = MGLHeatmapStyleLayer(identifier: identifier, source: source)
        
        // Adjust the color of the heatmap based on the point density.
        let colorDictionary: [NSNumber: UIColor] = [
            0.0: .clear,
            0.01: .white,
            0.15: UIColor(red: 0.19, green: 0.30, blue: 0.80, alpha: 1.0),
            0.5: UIColor(red: 0.73, green: 0.23, blue: 0.25, alpha: 1.0),
            1: .yellow
        ]
        
        heatmapLayer.heatmapColor =
            NSExpression(format: "mgl_interpolate:withCurveType:parameters:stops:($heatmapDensity, 'linear', nil, %@)", colorDictionary)
        
        // Heatmap weight measures how much a single data point impacts the layer's appearance.
        heatmapLayer.heatmapWeight =
            NSExpression(format: "mgl_interpolate:withCurveType:parameters:stops:(mag, 'linear', nil, %@)", [0: 0, 6: 1])
        
        // Heatmap intensity multiplies the heatmap weight based on zoom level.
        heatmapLayer.heatmapIntensity =
            NSExpression(format: "mgl_interpolate:withCurveType:parameters:stops:($zoomLevel, 'linear', nil, %@)", [0: 1, 9: 3])
        
        heatmapLayer.heatmapRadius =
            NSExpression(format: "mgl_interpolate:withCurveType:parameters:stops:($zoomLevel, 'linear', nil, %@)", [0: 4, 9: 30])
        
        // The heatmap layer should be visible up to zoom level 9.
        heatmapLayer.heatmapOpacity = NSExpression(format: "mgl_step:from:stops:($zoomLevel, 0.75, %@)", [0: 1.0, 20: 0])
        
        style.addLayer(heatmapLayer)
        
        // Add a circle layer to represent the earthquakes at higher zoom levels.
        let circleLayer = MGLCircleStyleLayer(identifier: "circle-layer", source: source)
        
        let magnitudeDictionary: [NSNumber: UIColor] = [
            0: .white,
            0.5: .yellow,
            2.5: UIColor(red: 0.73, green: 0.23, blue: 0.25, alpha: 1.0),
            5: UIColor(red: 0.19, green: 0.30, blue: 0.80, alpha: 1.0)
        ]
        circleLayer.circleColor = NSExpression(format: "mgl_interpolate:withCurveType:parameters:stops:(mag, 'linear', nil, %@)", magnitudeDictionary)
        
        // The heatmap layer will have an opacity of 0.75 up to zoom level 9, when the opacity becomes 0.
        circleLayer.circleOpacity = NSExpression(format: "mgl_step:from:stops:($zoomLevel, 0, %@)", [0: 0, 20: 1.0])
        circleLayer.circleRadius = NSExpression(forConstantValue: 20)
        style.addLayer(circleLayer)
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
