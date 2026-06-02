package terminalimage

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

const (
	kittyTerminalImageID     = 42
	kittyImagePlacementID    = 1
	terminalImageMaxByteSize = 10 << 20
)

type Image struct {
	Text string
	Raw  string
}

func RenderFromURL(imageURL string, maxCols int, maxRows int) (Image, error) {
	return RenderFromURLWithReferer(imageURL, maxCols, maxRows, "")
}

func RenderFromURLWithReferer(imageURL string, maxCols int, maxRows int, referer string) (Image, error) {
	req, err := http.NewRequest(http.MethodGet, imageURL, nil)
	if err != nil {
		return Image{}, fmt.Errorf("error creating image request: %w", err)
	}

	req.Header.Set("Accept", "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
	if referer != "" {
		req.Header.Set("Referer", referer)
	}
	req.Header.Set("User-Agent", "DailyDashBoard/1.0")

	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return Image{}, fmt.Errorf("error downloading image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Image{}, fmt.Errorf("image request failed: %s", resp.Status)
	}

	imageBytes, err := io.ReadAll(io.LimitReader(resp.Body, terminalImageMaxByteSize))
	if err != nil {
		return Image{}, fmt.Errorf("error reading image: %w", err)
	}

	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return Image{}, fmt.Errorf("error decoding image: %w", err)
	}

	if inlineImage := renderInlineTerminalImage(imageBytes, img, maxCols, maxRows); inlineImage.Text != "" || inlineImage.Raw != "" {
		return inlineImage, nil
	}

	return Image{Text: renderTerminalImage(img, maxCols, maxRows)}, nil
}

func RenderFromDataURL(dataURL string, maxCols int, maxRows int) (string, error) {
	_, encodedImage, found := strings.Cut(strings.TrimPrefix(dataURL, "data:"), ",")
	if !strings.HasPrefix(dataURL, "data:image/") || !found {
		return "", fmt.Errorf("thumbnail data was not available")
	}

	decodedImage, err := base64.StdEncoding.DecodeString(encodedImage)
	if err != nil {
		return "", fmt.Errorf("error decoding thumbnail data: %w", err)
	}

	img, _, err := image.Decode(bytes.NewReader(decodedImage))
	if err != nil {
		return "", fmt.Errorf("error decoding thumbnail image: %w", err)
	}

	return renderTerminalImage(img, maxCols, maxRows), nil
}

func renderInlineTerminalImage(imageBytes []byte, img image.Image, maxCols int, maxRows int) Image {
	switch inlineImageProtocol() {
	case "iterm2":
		return Image{Text: renderITerm2Image(imageBytes, maxCols, maxRows)}
	case "kitty":
		var pngBytes bytes.Buffer
		if err := png.Encode(&pngBytes, img); err != nil {
			return Image{}
		}

		return renderKittyImage(pngBytes.Bytes(), maxCols, maxRows)
	default:
		return Image{}
	}
}

func inlineImageProtocol() string {
	override := strings.ToLower(strings.TrimSpace(os.Getenv("DAILYDASHBOARD_IMAGE_PROTOCOL")))
	switch override {
	case "kitty", "iterm2":
		return override
	case "ansi", "none", "off":
		return ""
	}

	term := strings.ToLower(os.Getenv("TERM"))
	termProgram := strings.ToLower(os.Getenv("TERM_PROGRAM"))

	switch {
	case os.Getenv("KITTY_WINDOW_ID") != "",
		strings.Contains(term, "kitty"),
		strings.Contains(term, "ghostty"),
		strings.Contains(termProgram, "ghostty"),
		hasEnvWithPrefix("GHOSTTY_"):
		return "kitty"
	case os.Getenv("WEZTERM_EXECUTABLE") != "", strings.Contains(termProgram, "wezterm"), strings.Contains(termProgram, "iterm"):
		return "iterm2"
	default:
		return ""
	}
}

func hasEnvWithPrefix(prefix string) bool {
	for _, item := range os.Environ() {
		if strings.HasPrefix(item, prefix) {
			return true
		}
	}

	return false
}

func renderITerm2Image(imageBytes []byte, maxCols int, maxRows int) string {
	return terminalPassthrough(fmt.Sprintf(
		"\x1b]1337;File=inline=1;width=%d;height=%d;preserveAspectRatio=1:%s\a",
		maxCols,
		maxRows,
		base64.StdEncoding.EncodeToString(imageBytes),
	))
}

func renderKittyImage(imageBytes []byte, maxCols int, maxRows int) Image {
	encodedImage := base64.StdEncoding.EncodeToString(imageBytes)
	const chunkSize = 4096

	var graphics strings.Builder
	graphics.WriteString(fmt.Sprintf("\x1b_Ga=d,d=i,i=%d,q=2;\x1b\\", kittyTerminalImageID))

	firstChunk := true

	for len(encodedImage) > 0 {
		n := min(chunkSize, len(encodedImage))
		chunk := encodedImage[:n]
		encodedImage = encodedImage[n:]

		moreChunks := 0
		if len(encodedImage) > 0 {
			moreChunks = 1
		}

		if firstChunk {
			graphics.WriteString(fmt.Sprintf(
				"\x1b_Ga=T,f=100,i=%d,p=%d,c=%d,r=%d,U=1,q=2,m=%d;",
				kittyTerminalImageID,
				kittyImagePlacementID,
				maxCols,
				maxRows,
				moreChunks,
			))
			firstChunk = false
		} else {
			graphics.WriteString(fmt.Sprintf("\x1b_Gq=2,m=%d;", moreChunks))
		}

		graphics.WriteString(chunk)
		graphics.WriteString("\x1b\\")
	}

	return Image{
		Text: renderKittyPlaceholders(kittyTerminalImageID, maxCols, maxRows),
		Raw:  terminalPassthrough(graphics.String()),
	}
}

func renderKittyPlaceholders(imageID int, cols int, rows int) string {
	var builder strings.Builder

	for row := 0; row < rows; row++ {
		builder.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm", (imageID>>16)&0xff, (imageID>>8)&0xff, imageID&0xff))

		for col := 0; col < cols; col++ {
			builder.WriteRune('\U0010EEEE')
			builder.WriteRune(kittyPlaceholderDiacritic(row))
			builder.WriteRune(kittyPlaceholderDiacritic(col))
		}

		builder.WriteString("\x1b[0m")

		if row < rows-1 {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

func kittyPlaceholderDiacritic(index int) rune {
	diacritics := []rune{
		'\u0305', '\u030D', '\u030E', '\u0310', '\u0312', '\u033D', '\u033E', '\u033F',
		'\u0346', '\u034A', '\u034B', '\u034C', '\u0350', '\u0351', '\u0352', '\u0357',
		'\u035B', '\u0363', '\u0364', '\u0365', '\u0366', '\u0367', '\u0368', '\u0369',
		'\u036A', '\u036B', '\u036C', '\u036D', '\u036E', '\u036F', '\u0483', '\u0484',
		'\u0485', '\u0486', '\u0487', '\u0592', '\u0593', '\u0594', '\u0595', '\u0597',
		'\u0598', '\u0599', '\u059C', '\u059D', '\u059E', '\u059F', '\u05A0', '\u05A1',
		'\u05A8', '\u05A9', '\u05AB', '\u05AC', '\u05AF', '\u05C4', '\u0610', '\u0611',
		'\u0612', '\u0613', '\u0614', '\u0615', '\u0616', '\u0617', '\u0657', '\u0658',
		'\u0659', '\u065A', '\u065B', '\u065D', '\u065E', '\u06D6', '\u06D7', '\u06D8',
	}

	if index < 0 || index >= len(diacritics) {
		return diacritics[0]
	}

	return diacritics[index]
}

func terminalPassthrough(sequence string) string {
	if os.Getenv("TMUX") == "" {
		return sequence
	}

	return "\x1bPtmux;" + strings.ReplaceAll(sequence, "\x1b", "\x1b\x1b") + "\x1b\\"
}

func renderTerminalImage(img image.Image, maxCols int, maxRows int) string {
	bounds := img.Bounds()
	sourceWidth := bounds.Dx()
	sourceHeight := bounds.Dy()

	if sourceWidth == 0 || sourceHeight == 0 || maxCols == 0 || maxRows == 0 {
		return ""
	}

	maxPixelHeight := maxRows * 2
	scale := math.Min(
		float64(maxCols)/float64(sourceWidth),
		float64(maxPixelHeight)/float64(sourceHeight),
	)

	pixelWidth := max(1, int(math.Round(float64(sourceWidth)*scale)))
	pixelHeight := max(2, int(math.Round(float64(sourceHeight)*scale)))

	if pixelWidth > maxCols {
		pixelWidth = maxCols
	}

	if pixelHeight > maxPixelHeight {
		pixelHeight = maxPixelHeight
	}

	if pixelHeight%2 != 0 {
		pixelHeight--
	}

	var builder strings.Builder
	leftPadding := strings.Repeat(" ", max(0, (maxCols-pixelWidth)/2))

	for y := 0; y < pixelHeight; y += 2 {
		builder.WriteString(leftPadding)

		for x := 0; x < pixelWidth; x++ {
			top := sampleImageColor(img, bounds, x, y, pixelWidth, pixelHeight)
			bottom := sampleImageColor(img, bounds, x, y+1, pixelWidth, pixelHeight)

			builder.WriteString(fmt.Sprintf(
				"\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀",
				top.R,
				top.G,
				top.B,
				bottom.R,
				bottom.G,
				bottom.B,
			))
		}

		builder.WriteString("\x1b[0m")

		if y+2 < pixelHeight {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

func sampleImageColor(img image.Image, bounds image.Rectangle, x int, y int, width int, height int) color.NRGBA {
	sourceX := (float64(x)+0.5)*float64(bounds.Dx())/float64(width) - 0.5
	sourceY := (float64(y)+0.5)*float64(bounds.Dy())/float64(height) - 0.5

	return bilinearSample(img, bounds, sourceX, sourceY)
}

func bilinearSample(img image.Image, bounds image.Rectangle, sourceX float64, sourceY float64) color.NRGBA {
	x0 := clamp(int(math.Floor(sourceX)), 0, bounds.Dx()-1)
	y0 := clamp(int(math.Floor(sourceY)), 0, bounds.Dy()-1)
	x1 := clamp(x0+1, 0, bounds.Dx()-1)
	y1 := clamp(y0+1, 0, bounds.Dy()-1)

	xWeight := sourceX - math.Floor(sourceX)
	yWeight := sourceY - math.Floor(sourceY)

	topLeft := imageColorAt(img, bounds, x0, y0)
	topRight := imageColorAt(img, bounds, x1, y0)
	bottomLeft := imageColorAt(img, bounds, x0, y1)
	bottomRight := imageColorAt(img, bounds, x1, y1)

	top := mixColors(topLeft, topRight, xWeight)
	bottom := mixColors(bottomLeft, bottomRight, xWeight)

	return mixColors(top, bottom, yWeight)
}

func imageColorAt(img image.Image, bounds image.Rectangle, x int, y int) color.NRGBA {
	pixel := color.NRGBAModel.Convert(img.At(bounds.Min.X+x, bounds.Min.Y+y)).(color.NRGBA)

	if pixel.A == 255 {
		return pixel
	}

	alpha := float64(pixel.A) / 255
	return color.NRGBA{
		R: uint8(float64(pixel.R)*alpha + 255*(1-alpha)),
		G: uint8(float64(pixel.G)*alpha + 255*(1-alpha)),
		B: uint8(float64(pixel.B)*alpha + 255*(1-alpha)),
		A: 255,
	}
}

func mixColors(left color.NRGBA, right color.NRGBA, weight float64) color.NRGBA {
	return color.NRGBA{
		R: mixChannel(left.R, right.R, weight),
		G: mixChannel(left.G, right.G, weight),
		B: mixChannel(left.B, right.B, weight),
		A: 255,
	}
}

func mixChannel(left uint8, right uint8, weight float64) uint8 {
	return uint8(math.Round(float64(left)*(1-weight) + float64(right)*weight))
}

func clamp(value int, minimum int, maximum int) int {
	if value < minimum {
		return minimum
	}

	if value > maximum {
		return maximum
	}

	return value
}
