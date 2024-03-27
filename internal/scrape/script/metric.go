package script

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

const metricRegex = `,`

func ParseStdout(r io.Reader) ([]Data, error) {
	out := make([]Data, 0)
	pr := NewParser()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		mt, err := pr.ParseMetric(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("unable to parse metric: %w", err)
		}
		out = append(out, mt)
	}
	return out, nil
}

func NewParser() Parser {
	re := regexp.MustCompile(metricRegex)
	p := Parser{re: re}
	return p
}

type Parser struct {
	re *regexp.Regexp
}

func (p Parser) ParseMetric(l string) (Data, error) {
	lbs := make(map[string]string)
	var v float64
	mtcs := p.re.Split(l, -1)
	for i, mt := range mtcs {
		if i == len(mtcs)-1 {
			vl, err := strconv.ParseFloat(mt, 64)
			if err != nil {
				return Data{}, fmt.Errorf("invalid value of metric, %s is not float64, %v", mt, err)
			}
			v = vl
			continue
		}
		spl := strings.Index(mt, "=")
		lbs[mt[:spl]] = cleanBracket(mt[spl+1:])
	}
	return Data{
		Value:  v,
		Labels: lbs,
	}, nil
}

func cleanBracket(input string) string {
	if strings.HasPrefix(input, `"`) || strings.HasPrefix(input, "'") {
		input = input[1:]
	}
	if strings.HasSuffix(input, `"`) || strings.HasSuffix(input, "'") {
		input = input[:len(input)-1]
	}
	return input
}
