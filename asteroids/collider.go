// collider.go provides collision queries using the resolv library.
// The helper wraps IntersectionTest to support broad-phase queries against
// nearby cells or a specific target collider.
package asteroids

import "github.com/solarlune/resolv"

// checkCollision reports whether obj intersects anything relevant.
//
// If against is nil, the test runs against shapes from the neighboring
// grid cells (broad-phase). If against is non-nil, the test is restricted
// to shapes near that specific collider to reduce search scope.
//
// The OnIntersect callback returns true immediately on the first hit,
// which short-circuits the query for efficiency.
func (g *GameScene) checkCollision(obj, against *resolv.Circle) bool {
	if against == nil {
		// Broad-phase: query nearby cells around obj and test against all shapes.
		return obj.IntersectionTest(resolv.IntersectionTestSettings{
			TestAgainst: obj.SelectTouchingCells(1).FilterShapes(),
			OnIntersect: func(set resolv.IntersectionSet) bool {
				// Early exit on first intersection.
				return true
			},
		})
	}

	// Narrow-phase: query cells near the target collider only.
	return obj.IntersectionTest(resolv.IntersectionTestSettings{
		TestAgainst: against.SelectTouchingCells(1).FilterShapes(),
		OnIntersect: func(set resolv.IntersectionSet) bool {
			return true
		},
	})
}
