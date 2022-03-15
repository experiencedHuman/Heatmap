//
//  HeatmapApp.swift
//  Heatmap
//
//  Created by Kamber Vogli on 15.03.22.
//

import SwiftUI
import AzureMapsControl

@main
struct HeatmapApp: App {
    
    init() {
//        GMSServices.provideAPIKey("AIzaSyBC3NMmAu-jAl6fhAdWVDHCAQ0QpBI24iU")
        AzureMaps.configure(accessToken: "9g_Nd80y0xdLKe5Mv9ZsChjVEiTs0sTw8nHfFYE25_g")
    }
    
    var body: some Scene {
        WindowGroup {
            ContentView()
        }
    }
}
