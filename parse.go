package dispatch

import (
	"fmt"
	"strings"
)

// parseArgs validates & parses required arguments/argument combinations
func parseArgs(cmd *Command, data *Data, args []string) error {
	// If no combos defined, use positional parsing
	if len(cmd.ArgumentCombos) == 0 {
		return parsePositionalArgs(cmd, data, args)
	}

	return parseComboArgs(cmd, data, args)
}

// parsePositionalArgs handles the simple case where no argument combos are defined and we just have a list of args
func parsePositionalArgs(cmd *Command, data *Data, args []string) error {
	if len(args) < cmd.ArgsRequired {
		return fmt.Errorf("Missing required arguments:\n```%s %s```", cmd.Command, BuildPositionalDisplay(cmd))
	}

	// Reject extra args unless last arg is string
	// TODO: Add a field to the arg definition to specify that it should absorb remaining args, instead of just assuming this for the last string-type arg?
	if len(args) > len(cmd.Args) {
		last := cmd.Args[len(cmd.Args)-1]
		if last.Type != String {
			return fmt.Errorf("Too many arguments provided.\n```%s %s```", cmd.Command, BuildPositionalDisplay(cmd))
		}
	}

	parsedArgs := []*ParsedArg{}
	for i, arg := range cmd.Args {
		var value string
		if i < len(args) {
			value = args[i]
			// If the last arg is of type String, it absorbs the remaining args, so we join them together with spaces. This allows for multi-word string arguments without needing to use quotes.
			// TODO: Add a field to the arg definition to specify that it should absorb remaining args, instead of just assuming this for the last string-type arg?
			if i == len(cmd.Args)-1 && arg.Type == String && len(args) > i {
				value = strings.Join(args[i:], " ")
			}
		} else {
			parsedArgs = append(parsedArgs, &ParsedArg{Def: arg})
			continue
		}

		parsedArgs = append(parsedArgs, &ParsedArg{
			Def: arg,
			Raw: value,
		})
	}

	// Validate supplied arguments
	// We only validate the supplied args, so if the user doesn't provide an optional arg, they won't get an error for it
	for _, pArg := range parsedArgs {
		if pArg.Raw == "" {
			continue
		}

		value, ok := pArg.Def.Type.ValidateArg(pArg, data)
		if !ok {
			return fmt.Errorf("Invalid `%s` argument. Expected: `%s`", pArg.Def.Name, pArg.Def.Type.Help())
		}

		pArg.Value = value
	}

	data.ParsedArgs = parsedArgs
	return nil
}

// BuildPositionalDisplay returns the usage for positional args
// Exported for use in custom help functions
func BuildPositionalDisplay(cmd *Command) string {
	var display strings.Builder
	for i, a := range cmd.Args {
		if i < cmd.ArgsRequired {
			display.WriteString(" <" + a.Name + ":" + a.Type.Help() + ">")
		} else {
			display.WriteString(" [" + a.Name + ":" + a.Type.Help() + "]")
		}
	}
	return display.String()
}

// parseComboArgs handles the case where argument combos are defined and we need to determine which combo the user is trying to use, if any
func parseComboArgs(cmd *Command, data *Data, args []string) error {
	var lastErr error
	for i, combo := range cmd.ArgumentCombos {
		parsed, err := parseCombo(cmd, data, combo, args)
		if err == nil {
			data.ParsedArgs = parsed
			return nil
		}

		lastErr = fmt.Errorf("Combination %d failed: %v", i+1, err)
	}

	return fmt.Errorf("No matching argument combination found.\n%s", lastErr.Error())
}

func parseCombo(cmd *Command, data *Data, combo []int, args []string) ([]*ParsedArg, error) {
	parsedArgs := make([]*ParsedArg, len(cmd.Args))
	argIdx := 0

	for i, defPos := range combo {
		def := cmd.Args[defPos]
		var value string

		argsLeft := len(args) - argIdx
		remainingComboArgs := len(combo) - (i + 1)

if argsLeft <= 0 {
    return nil, fmt.Errorf(
        "Missing argument `%s`. Expected: `%s`",
        def.Name,
        def.Type.Help(),
    )
}

		// String argument logic
		if def.Type == String {
			if i == len(combo)-1 {
				// Last string in combo absorbs all remaining args
				value = strings.Join(args[argIdx:], " ")
				argIdx = len(args)
			} else {
				// Absorb enough args to leave remaining combo args untouched
				if argsLeft < remainingComboArgs {
					return nil, fmt.Errorf("Not enough arguments for remaining combos\nExpected one of:```%s```", BuildComboDisplay(cmd))
				}
				absorbCount := argsLeft - remainingComboArgs
				if absorbCount > 0 {
					value = strings.Join(args[argIdx:argIdx+absorbCount], " ")
					argIdx += absorbCount
				} else {
					value = args[argIdx]
					argIdx++
				}
			}
		} else {
			if argsLeft <= 0 {
				return nil, fmt.Errorf("Missing argument `%s`. Expected: `%s`",
					def.Name, def.Type.Help())
			}

			value = args[argIdx]
			argIdx++
		}

		// Validate if value is provided
		p := &ParsedArg{Def: def, Raw: value}
		if value != "" {
			value, ok := def.Type.ValidateArg(p, data)
			if !ok {
				return nil, fmt.Errorf("Invalid `%s` argument. Expected: `%s`", def.Name, def.Type.Help())
			}

			p.Value = value
		}

		parsedArgs[defPos] = p
	}

	// Fail if extra user args left after parsing combo
	if argIdx < len(args) {
		return nil, fmt.Errorf("Too many arguments provided\nExpected one of:```%s```", BuildComboDisplay(cmd))
	}

	return parsedArgs, nil
}

// BuildComboDisplay returns the usage for argument combos
// Exported for use in custom help functions
func BuildComboDisplay(cmd *Command) string {
	parts := []string{}
	for _, combo := range cmd.ArgumentCombos {
		var s strings.Builder
		s.WriteString(cmd.Command)
		for _, idx := range combo {
			if idx < 0 || idx >= len(cmd.Args) {
				continue
			}
			argHelp := cmd.Args[idx].Name + ":" + cmd.Args[idx].Type.Help()
			s.WriteString(" <" + argHelp + ">")
		}
		parts = append(parts, s.String())
	}
	return strings.Join(parts, "\n")
}
