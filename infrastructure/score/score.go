package score

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/opensourceways/xihe-script/domain/score"
	"github.com/opensourceways/xihe-script/infrastructure/message"
)

type calculateImpl struct {
	calculate string
}

func NewCalculateScore(calculate string) score.CalculateScore {
	return &calculateImpl{
		calculate: calculate,
	}
}

type evaluateImpl struct {
	evaluate string
}

func NewEvaluateScore(evaluate string) score.EvaluateScore {
	return &evaluateImpl{
		evaluate: evaluate,
	}
}

func (s *evaluateImpl) Evaluate(col *message.MatchFields) (data []byte, err error) {
	args := []string{s.evaluate, "--pred_path", col.Path, "--true_path", col.AnswerPath, "--cls", strconv.Itoa(col.Cls), "--pos", strconv.Itoa(col.Pos)}
	data, err = exec.Command("python3", args...).Output()

	if err != nil {
		return
	}
	data = bytes.ReplaceAll(bytes.TrimSpace(data), []byte(`'`), []byte(`"`))

	return
}

func (s *calculateImpl) Calculate(col *message.MatchFields) (data []byte, err error) {
	if len(col.Path) == 0 {
		return nil, fmt.Errorf("path is empty")
	}
	path := filepath.Join(os.Getenv("UPLOAD"), fmt.Sprintf("%d", time.Now().UnixMicro()))
	args := []string{s.calculate, "--user_result", col.Path, "--unzip_path", path, "--fid_weights_file", col.FidWeightsPath, "--real_result", col.RealPath}
	data, err = exec.Command("python3", args...).Output()

	if err != nil {
		if err2 := os.RemoveAll(path); err2 != nil {
			return nil, err2
		}

		return
	}

	data = bytes.ReplaceAll(bytes.TrimSpace(data), []byte(`'`), []byte(`"`))
	
	if err = os.RemoveAll(path); err != nil {
		return
	}

	return
}
