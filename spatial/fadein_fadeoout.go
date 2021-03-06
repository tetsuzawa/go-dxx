package spatial

import (
	"fmt"
	"math"
	"os"

	"github.com/tetsuzawa/go-soundlib/dxx"
)

func FadeinFadeout(subject, soundName string, moveWidth, moveVelocity, endAngle int, outDir string) error {
	const (
		repeatTimes  = 1
		samplingFreq = 48 // [kHz]
	)

	// 移動時間 [ms]
	var moveTime float64 = float64(moveWidth) * 1000.0 / float64(moveVelocity)
	// 移動角度
	var moveAngle int = moveWidth*repeatTimes + 1

	// 1度動くのに必要なサンプル数
	// [ms]*[kHz] / [deg] = [sample/deg]
	var dwellingSamples int = int(moveTime) * samplingFreq / (moveWidth*repeatTimes*2 + 1)
	var durationSamples int = dwellingSamples * 63 / 64
	var overlapSamples int = dwellingSamples * 1 / 64

	fadeinFilter, fadeoutFilter := GenerateFadeinFadeoutFilt(overlapSamples)

	// 音データの読み込み
	sound, err := dxx.ReadFromFile(soundName)
	if err != nil {
		return err
	}

	for _, direction := range []string{"c", "cc"} {
		for _, LR := range []string{"L", "R"} {
			moveOut := make([]float64, dwellingSamples, int(moveTime)*samplingFreq)
			usedAngles := make([]int, moveAngle*2-1)

			for angle := 0; angle < (moveAngle*2 - 1); angle++ {
				// ノコギリ波の生成
				dataAngle := angle % ((moveWidth * 2) * 2)
				// ノコギリ波から三角波を生成
				if dataAngle > moveWidth*2 {
					dataAngle = (moveWidth*2)*2 - dataAngle
				}
				if direction == "cc" {
					dataAngle = -dataAngle
				}
				dataAngle = dataAngle / 2
				if dataAngle < 0 {
					dataAngle += 3600
				}

				// SLTFの読み込み
				SLTFName := fmt.Sprintf("%s/SLTF/SLTF_%d_%s.DDB", subject, (endAngle+dataAngle)%3600, LR)
				SLTF, err := dxx.ReadFromFile(SLTFName)
				if err != nil {
					return err
				}
				usedAngles[angle] = (endAngle + dataAngle) % 3600

				// Fadein-Fadeout
				// 音データと伝達関数の畳込み
				cutSound := sound[angle*(durationSamples+overlapSamples) : durationSamples*2+angle*(durationSamples+overlapSamples)+len(SLTF)*3+1]
				soundSLTF := LinearConvolutionTimeDomain(cutSound, SLTF)
				// 無音区間の切り出し
				soundSLTF = soundSLTF[len(SLTF)*2 : len(soundSLTF)-len(SLTF)*2]
				// 前の角度のfadeout部と現在の角度のfadein部の加算
				fadein := make([]float64, overlapSamples)
				for i := range fadein {
					fadein[i] = soundSLTF[i] * fadeinFilter[i]
					moveOut[(durationSamples+overlapSamples)*angle+i] += fadein[i]
				}

				// 持続時間
				moveOut = append(moveOut, soundSLTF[overlapSamples:len(soundSLTF)-overlapSamples]...)

				// fadeout
				fadeout := make([]float64, overlapSamples)
				for i := range fadein {
					fadeout[i] = soundSLTF[len(soundSLTF)-overlapSamples+i] * fadeoutFilter[i]
				}
				moveOut = append(moveOut, fadeout...)
			}

			// 先頭のFadein部をカット
			out := moveOut[overlapSamples:]

			// DDBへ出力
			outName := fmt.Sprintf("%s/move_judge_w%03d_mt%03d_%s_%d_%s.DDB", outDir, moveWidth, moveVelocity, direction, endAngle, LR)
			if err := dxx.WriteToFile(outName, out); err != nil {
				return err
			}
			_, err = fmt.Fprintf(os.Stderr, "%s: length=%d\n", outName, len(out))
			if err != nil {
				return err
			}
			_, err := fmt.Fprintf(os.Stderr, "used angle:%v\n", usedAngles)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GenerateFadeinFadeoutFilt(length int) (fadeinFilt, fadeoutFilt []float64) {
	// Fourier Series Window Coefficient
	a0 := (1 + math.Sqrt(2)) / 4
	a1 := 0.25 + 0.25*math.Sqrt((5-2*math.Sqrt(2))/2)
	a2 := (1 - math.Sqrt(2)) / 4
	a3 := 0.25 - 0.25*math.Sqrt((5-2*math.Sqrt(2))/2)

	// Fourier series window
	fadeinFilt = make([]float64, length)
	fadeoutFilt = make([]float64, length)
	flength := float64(length)
	for i := 0; i < length; i++ {
		f := float64(i)
		fadeinFilt[i] = a0 - a1*math.Cos(math.Pi/flength*f) + a2*math.Cos(2.0*math.Pi/flength*f) - a3*math.Cos(3.0*math.Pi/flength*f)
		fadeoutFilt[i] = a0 + a1*math.Cos(math.Pi/flength*f) + a2*math.Cos(2.0*math.Pi/flength*f) + a3*math.Cos(3.0*math.Pi/flength*f)
	}
	return fadeinFilt, fadeoutFilt
}
