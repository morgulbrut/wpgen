package sketch

import (
	"image"
	"image/color"
	"math/rand"
	"strings"

	"github.com/fogleman/gg"
)

type UserParams struct {
	StrokeRatio              float64
	DestWidth                int
	DestHeight               int
	InitialAlpha             float64
	StrokeReduction          float64
	AlphaIncrease            float64
	StrokeInversionThreshold float64
	StrokeJitter             int
	MinEdgeCount             int
	MaxEdgeCount             int
	RotationJitter           float64
	Shape                    string
	Fill                     bool
	Stroke                   bool
}

type Sketch struct {
	UserParams        // embed for easier access
	source            image.Image
	dc                *gg.Context
	sourceWidth       int
	sourceHeight      int
	strokeSize        float64
	initialStrokeSize float64
}

func NewSketch(source image.Image, userParams UserParams) *Sketch {
	s := &Sketch{UserParams: userParams}
	bounds := source.Bounds()
	s.sourceWidth, s.sourceHeight = bounds.Max.X, bounds.Max.Y
	s.initialStrokeSize = s.StrokeRatio * float64(s.DestWidth)
	s.strokeSize = s.initialStrokeSize

	canvas := gg.NewContext(s.DestWidth, s.DestHeight)
	canvas.SetColor(color.Black)
	canvas.DrawRectangle(0, 0, float64(s.DestWidth), float64(s.DestHeight))
	canvas.FillPreserve()

	s.source = source
	s.dc = canvas
	return s
}

func (s *Sketch) Update() {

	// fmt.Print(".") // uncomment for some simple update bar

	// Step 1
	// choose random point on the source image and get the colour of that point
	rndX := rand.Float64() * float64(s.sourceWidth)
	rndY := rand.Float64() * float64(s.sourceHeight)
	r, g, b := rgb255(s.source.At(int(rndX), int(rndY)))

	// Step 2
	// set the coordinates of the shape to draw to more or less those from Step 1
	destX := rndX * float64(s.DestWidth) / float64(s.sourceWidth)
	destX += float64(randRange(s.StrokeJitter))
	destY := rndY * float64(s.DestHeight) / float64(s.sourceHeight)
	destY += float64(randRange(s.StrokeJitter))

	// Step 3
	// Draw the shape
	switch strings.ToLower(s.Shape) {
	case "circle":
		s.dc.DrawCircle(destX, destY, s.strokeSize)
	case "roundedsquare":
		s.dc.DrawRoundedRectangle(destX, destY, s.strokeSize, s.strokeSize, s.strokeSize/12)
	case "square":
		edges := 4
		s.dc.DrawRegularPolygon(edges, destX, destY, s.strokeSize, rand.Float64()*s.RotationJitter)
	case "hexagon":
		edges := 6
		s.dc.DrawRegularPolygon(edges, destX, destY, s.strokeSize, rand.Float64()*s.RotationJitter)
	default:
		edges := s.MinEdgeCount + rand.Intn(s.MaxEdgeCount-s.MinEdgeCount+1)
		s.dc.DrawRegularPolygon(edges, destX, destY, s.strokeSize, rand.Float64()*s.RotationJitter)
	}

	// Step 4
	// fill the shape with the chosen colour if -nofill isn't set
	if !s.Fill {
		// set the drawing colour to the one from Step 1
		s.dc.SetRGBA255(r, g, b, int(s.InitialAlpha))
		// fill the shape
		s.dc.FillPreserve()
	}

	// Step 5
	// draw the outline
	// if -nostroke isn't set choose black or white depending on the colour from Step 1
	if !s.Stroke {
		if s.strokeSize <= s.StrokeInversionThreshold*s.initialStrokeSize {
			if (r+g+b)/3 < 128 {
				s.dc.SetRGBA255(255, 255, 255, int(s.InitialAlpha*2))
			} else {
				s.dc.SetRGBA255(0, 0, 0, int(s.InitialAlpha*2))
			}
		}
		// if -nostroke is set, set the colour to alpha 0
		// (just not calling s.dc.Stroke() makes the program slow as hell)
	} else {
		s.dc.SetRGBA255(0, 0, 0, 0)
	}
	// finally draw the outline
	s.dc.Stroke()

	// Step 6
	// adapt the strokesize and alpha for the next run
	s.strokeSize -= s.StrokeReduction * s.strokeSize
	s.InitialAlpha += s.AlphaIncrease
}

func (s *Sketch) Output() image.Image {
	return s.dc.Image()
}

func rgb255(c color.Color) (r, g, b int) {
	r0, g0, b0, _ := c.RGBA()
	return int(r0 / 257), int(g0 / 257), int(b0 / 257)
}

func randRange(max int) int {
	return -max + rand.Intn(2*max)
}
