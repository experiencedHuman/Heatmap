//
//  Interpolation.swift
//  Heatmap
//
//  Created by Kamber Vogli on 26.04.22.
//

import Foundation

struct Point {
    var x: Double
    var y: Double
}

struct Interpolation {
    var radiusPoints: [Point] = [
        Point(x: 1, y: 5),
        Point(x: 3, y: 7),
        Point(x: 5, y: 8),
        Point(x: 7, y: 10),
        Point(x: 8, y: 12),
        Point(x: 10, y: 15),
        Point(x: 12, y: 25),
        Point(x: 13, y: 26),
        Point(x: 14, y: 27),
        Point(x: 15, y: 29),
        Point(x: 17, y: 30),
        Point(x: 20, y: 45),
        Point(x: 23, y: 40)
    ]
    /// a dictionary of `zoom level` key to `radius` value
    /// stores a radius value (for the heatmap) from zoom level 1 to 23 in 0.3 intervals
    var zoomToRadius = [Double : Double]()
    
    var opacityPoints: [Point] = [
        Point(x: 1, y: 0.8),
        Point(x: 2, y: 0.85),
        Point(x: 4, y: 0.87),
        Point(x: 5, y: 0.89),
        Point(x: 6, y: 0.92),
        Point(x: 8, y: 0.97),
        Point(x: 10, y: 1.0),
        Point(x: 12, y: 1.19),
        Point(x: 13, y: 2.0),
        Point(x: 14, y: 2.3),
        Point(x: 15, y: 2.5),
        Point(x: 16, y: 6),
        Point(x: 18, y: 7),
        Point(x: 20, y: 10),
        Point(x: 23, y: 9)
    ]
    
    var weightPoints: [Point] = [
        Point(x: 1, y: 0.001),
        Point(x: 2, y: 0.002),
        Point(x: 4, y: 0.003),
        Point(x: 5, y: 0.007),
        Point(x: 6, y: 0.009),
        Point(x: 8, y: 0.01),
        Point(x: 10, y: 0.04),
        Point(x: 12, y: 0.06),
        Point(x: 13, y: 0.1),
        Point(x: 14, y: 0.17),
        Point(x: 15, y: 0.23),
        Point(x: 16, y: 0.5),
        Point(x: 18, y: 0.6),
        Point(x: 20, y: 1.0),
        Point(x: 23, y: 1.2)
    ]
    /// a dictionary of `zoom level` key to `radius` value
    /// stores a radius value (for the heatmap) from zoom level 1 to 23 in 0.3 intervals
    var zoomToOpacity = [Double : Double]()
    
    mutating func precalculateRadius() {
        for i in stride(from: 1, to: 23, by: 0.3) {
            zoomToRadius[i] = lagrange(points: radiusPoints, x: i)
        }
    }
    
    mutating func precalculateOpacity() {
        for i in stride(from: 1, to: 23, by: 0.3) {
            zoomToOpacity[i] = lagrange(points: opacityPoints, x: i)
        }
    }
}

func lagrange(points: [Point], x: Double) -> Double {
    var res = 0.0
    let n = points.count - 1
    
    for i in 1...n {
        var prod = 1.0
        for j in 1...n {
            if i != j {
                prod *= (x - points[j].x) / (points[i].x - points[j].x)
            }
        }
        prod *= points[i].y
        res += prod
    }
    
    return res
}


// forward implementation
func newton(points: [Point], x: Double) -> Double {
    let n = points.count
    var y = Array(repeating: Array(repeating: 0.0, count: n), count: n)
    
    for i in 0..<n {
        y[i][0] = points[i].y
    }
    
    for i in 1..<n {
        for j in 0..<(n - i) {
            y[j][i] = y[j + 1][i - 1] - y[j][i - 1]
        }
    }
    
    var sum = y[0][0]
    let u = (x - points[0].x) / (points[1].x - points[0].x)
    for i in 1..<n {
        sum += (calcU(u, n: i) * y[0][i]) / Double(fac(i))
    }
    
    return sum
}

private func calcU(_ u: Double, n: Int) -> Double {
    var temp = u
    for i in 1..<n {
        temp *= (u - Double(i))
    }
    return temp
}

private func fac(_ n: Int) -> Int {
    if n <= 2 {
        return n
    }
    
    var f = 1
    for i in 2...n {
        f *= i
    }
    return f
}

func spline(points: [Point]) -> Float {
    //TODO: implement
    return 0.0
}
