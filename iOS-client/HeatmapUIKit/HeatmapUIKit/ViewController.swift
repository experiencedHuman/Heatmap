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
        let url = URL(string: "https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/all_week.geojson")!
        source.importData(fromURL: url)
        
        azureMap.onReady { map in
            map.sources.add(source)
            let layer = HeatMapLayer(source: source, options: [
                .heatmapRadius(10),
                .heatmapOpacity(0.8)
            ])
            map.layers.insertLayer(layer, below: "labels")
        }
        
    }
    
}

