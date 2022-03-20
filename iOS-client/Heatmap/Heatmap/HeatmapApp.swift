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
        GMSServices.provideAPIKey("AIzaSyBC3NMmAu-jAl6fhAdWVDHCAQ0QpBI24iU")
    }
    
    var body: some Scene {
        WindowGroup {
            ContentView()
        }
    }
}
