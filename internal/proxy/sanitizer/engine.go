package sanitizer

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

type Pattern struct {
	Name    string
	Matcher *regexp.Regexp
}

type ConfigRule struct {
	Name       string `yaml:"name"`
	Expression string `yaml:"expression"`
}

type RuleConfig struct {
	CustomRules []ConfigRule `yaml:"custom_rules"`
}

type Engine struct {
	patterns []Pattern
}

func NewEngine() *Engine {
	return &Engine{
		patterns: []Pattern{
			{"RSA/OpenSSH Private Key", regexp.MustCompile(`-----BEGIN (RSA|DSA|EC|OPENSSH|PRIVATE) KEY-----[A-Za-z0-9+\/=\s]+-----END \1 KEY-----`)},
			{"OpenAI API Key", regexp.MustCompile(`sk-[a-zA-Z0-9]{20,}`)},
			{"AWS Access Key", regexp.MustCompile(`AKIA[0-9A-Z]{16}`)},
			{"GitHub PAT", regexp.MustCompile(`ghp_[a-zA-Z0-9]{36}|github_pat_[a-zA-Z0-9_]{22,}`)},
			{"JWT Token", regexp.MustCompile(`ey[A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*`)},
			{"Database Connection URI", regexp.MustCompile(`(postgres|mysql|mongodb|redis)://[a-zA-Z0-9_-]+:[^@\s]+@[a-zA-Z0-9.-]+:[0-9]+/[a-zA-Z0-9_-]+`)},
			{"Generic Bearer Token", regexp.MustCompile(`(?i)bearer\s+[a-zA-Z0-9_\-\.]{20,}`)},
			{"Stripe Live Key", regexp.MustCompile(`(rk|sk)_live_[0-9a-zA-Z]{24,}`)},
			{"Slack Token", regexp.MustCompile(`xox[baprs]-[0-9a-zA-Z]{10,48}`)},
			{"Email Address", regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)},
			{"IPv4 Address", regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)},
		},
	}
}

// LoadConfigFromFile reads a YAML rule file and injects custom regex patterns
func (e *Engine) LoadConfigFromFile(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Optional config, skip if not present
		}
		return err
	}

	var config RuleConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse rule yaml config: %v", err)
	}

	for _, rule := range config.CustomRules {
		re, err := regexp.Compile(rule.Expression)
		if err != nil {
			return fmt.Errorf("invalid regex for rule '%s': %v", rule.Name, err)
		}
		e.patterns = append(e.patterns, Pattern{Name: rule.Name, Matcher: re})
	}

	return nil
}

func (e *Engine) Sanitize(input []byte) []byte {
	output := string(input)
	for _, p := range e.patterns {
		output = p.Matcher.ReplaceAllString(output, fmt.Sprintf("[REDACTED_%s]", p.Name))
	}
	return []byte(output)
}
