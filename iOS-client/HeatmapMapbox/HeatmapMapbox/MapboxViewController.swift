//
//  MapboxViewController.swift
//  HeatmapMapbox
//
//  Created by Kamber Vogli on 15.03.22.
//

import Foundation


import Mapbox
import SwiftUI

struct MapboxView: UIViewControllerRepresentable {
    
    func makeUIViewController(context: Context) -> MapboxViewController {
        MapboxViewController()
    }
    
    func updateUIViewController(_ uiViewController: MapboxViewController, context: Context) {
        
    }
}
 
class MapboxViewController: UIViewController {
    override func viewDidLoad() {
    super.viewDidLoad()
     
    let url = URL(string: "mapbox://styles/mapbox/streets-v11")
    let mapView = MGLMapView(frame: view.bounds, styleURL: url)
    mapView.autoresizingMask = [.flexibleWidth, .flexibleHeight]
    mapView.setCenter(CLLocationCoordinate2D(latitude: 59.31, longitude: 18.06), zoomLevel: 9, animated: false)
    view.addSubview(mapView)
    }
}
