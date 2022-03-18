//
//  HeatmapApp.swift
//  Heatmap
//
//  Created by Kamber Vogli on 15.03.22.
//

import SwiftUI
import GoogleMaps

@main
struct HeatmapApp: App {
    
    init() {
        if let googleMapsToken = Bundle.main.infoDictionary?["GoogleMaps"] as? String {
            GMSServices.provideAPIKey(googleMapsToken)
        }
    }
    
    var body: some Scene {
        WindowGroup {
            ContentView()
        }
    }
}
