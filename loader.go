// Package dotenv provides functionality for loading environment variables
// from a `.env` file and setting them as environment variables in the system.
// It supports reading key-value pairs from a file, skipping invalid or comment lines,
// and allows users to customize the filename and logger through options.
//
// Example usage:
//
//   envVars, err := dotenv.Config(
//     dotenv.WithFilename(".env"), 
//     dotenv.WithLogger(log.New(os.Stdout, "dotenv: ", log.LstdFlags)),
//   )
//   if err != nil {
//     log.Fatal(err)
//   }
//   fmt.Println(envVars)
package dotenv

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func load(filename string) (map[string]string, error) {
	var key, value string
	envVars := make(map[string]string)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.SplitN(line, "=", 2)

		if len(tokens) != 2 || len(tokens[0]) == 0 || len(tokens[1]) == 0 {
			continue
		} else if string(tokens[0][0]) == "#" || string(tokens[0][0]) == "//" {
			continue
		} else if string(tokens[0]) == "" || string(tokens[1]) == "" {
			continue
		} else {
			value = strings.TrimSpace(tokens[1])
			key = strings.TrimSpace(tokens[0])

			keysize := strings.Split(key, " ")
			valuesize := strings.Split(value, " ")

			if len(keysize) != 1 {
				continue
			}

			if len(valuesize) > 1 {
				if string(valuesize[1][0]) == "#" || string(valuesize[1][:2]) == "//" {
					value = valuesize[0]
				}
			}

			if value[0] == '"' && value[len(value)-1] == '"' {
				value = strings.Trim(value, "\"")
			}

			envVars[key] = value
		}
	}

	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	return envVars, nil
}

type options struct {
	filename string
	logger   *log.Logger
}

type Option interface {
	apply(*options)
}

type filenameOption string

func (f filenameOption) apply(opts *options) {
	opts.filename = string(f)
}

func WithFilename(f string) Option {
	return filenameOption(f)
}

type loggerOption struct {
	Logger *log.Logger
}

func (l loggerOption) apply(opts *options) {
	opts.logger = l.Logger
}

func WithLogger(l *log.Logger) Option {
	return loggerOption{Logger: l}
}

func Config(params ...Option) (map[string]string, error) {
	opts := options{
		filename: ".env",
		logger:   log.New(os.Stderr, "dotenv:", log.LstdFlags),
	}

	for _, val := range params {
		val.apply(&opts)
	}

	envVars, err := load(opts.filename)
	if err != nil {
		if opts.logger != nil {
			opts.logger.Println(err)
		}
		return nil, err
	}
	for key, value := range envVars {
		err = os.Setenv(key, value)
		if err != nil {
			if opts.logger != nil {
				opts.logger.Println(err)
			}
			return nil, err
		}
	}

	return envVars, nil
}
