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
        if let azureToken = Bundle.main.infoDictionary?["Azure_Token"] as? String {
            GMSServices.provideAPIKey(azureToken)
        }
    }
    
    var body: some Scene {
        WindowGroup {
            ContentView()
        }
    }
}
