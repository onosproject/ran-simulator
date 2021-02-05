// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package cli

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// getBashCompletions returns a bash completion script from all dependencies
func getBashCompletions() string {
	completions := []string{}
	return strings.Join(completions, "\n")
}

func getCompletionCommand() *cobra.Command {
	return &cobra.Command{
		Use:       "completion <shell>",
		Short:     "Generated bash or zsh auto-completion script",
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"bash", "zsh"},
		Example: `For bash run the following command from the shell: eval $(simcli completion bash).
For zsh run the following command from the shell: source <(simcli completion zsh).`,
		Run: runCompletionCommand,
	}
}

func runCompletionCommand(cmd *cobra.Command, args []string) {
	if args[0] == "bash" {
		if err := runCompletionBash(os.Stdout, cmd.Parent()); err != nil {
			ExitWithError(ExitError, err)
		}
	} else if args[0] == "zsh" {
		if err := runCompletionZsh(os.Stdout, cmd.Parent()); err != nil {
			ExitWithError(ExitError, err)
		}

	} else {
		ExitWithError(ExitError, errors.New("unsupported shell type "+args[0]))
	}
}

func runCompletionBash(out io.Writer, cmd *cobra.Command) error {
	return cmd.GenBashCompletion(out)
}

func runCompletionZsh(out io.Writer, cmd *cobra.Command) error {
	header := "#compdef simcli\n"

	_, err := out.Write([]byte(header))
	if err != nil {
		ExitWithError(ExitError, err)
	}

	init := `
__simcli_bash_source() {
	alias shopt=':'
	alias _expand=_bash_expand
	alias _complete=_bash_comp
	emulate -L sh
	setopt kshglob noshglob braceexpand
	source "$@"
}
__simcli_type() {
	# -t is not supported by zsh
	if [ "$1" == "-t" ]; then
		shift
		# fake Bash 4 to disable "complete -o nospace". Instead
		# "compopt +-o nospace" is used in the code to toggle trailing
		# spaces. We don't support that, but leave trailing spaces on
		# all the time
		if [ "$1" = "__simcli_compopt" ]; then
			echo builtin
			return 0
		fi
	fi
	type "$@"
}
__simcli_compgen() {
	local completions w
	completions=( $(compgen "$@") ) || return $?
	# filter by given word as prefix
	while [[ "$1" = -* && "$1" != -- ]]; do
		shift
		shift
	done
	if [[ "$1" == -- ]]; then
		shift
	fi
	for w in "${completions[@]}"; do
		if [[ "${w}" = "$1"* ]]; then
			echo "${w}"
		fi
	done
}
__simcli_compopt() {
	true # don't do anything. Not supported by bashcompinit in zsh
}
__simcli_ltrim_colon_completions()
{
	if [[ "$1" == *:* && "$COMP_WORDBREAKS" == *:* ]]; then
		# Remove colon-word prefix from COMPREPLY items
		local colon_word=${1%${1##*:}}
		local i=${#COMPREPLY[*]}
		while [[ $((--i)) -ge 0 ]]; do
			COMPREPLY[$i]=${COMPREPLY[$i]#"$colon_word"}
		done
	fi
}
__simcli_get_comp_words_by_ref() {
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[${COMP_CWORD}-1]}"
	words=("${COMP_WORDS[@]}")
	cword=("${COMP_CWORD[@]}")
}
__simcli_filedir() {
	local RET OLD_IFS w qw
	__debug "_filedir $@ cur=$cur"
	if [[ "$1" = \~* ]]; then
		# somehow does not work. Maybe, zsh does not call this at all
		eval echo "$1"
		return 0
	fi
	OLD_IFS="$IFS"
	IFS=$'\n'
	if [ "$1" = "-d" ]; then
		shift
		RET=( $(compgen -d) )
	else
		RET=( $(compgen -f) )
	fi
	IFS="$OLD_IFS"
	IFS="," __debug "RET=${RET[@]} len=${#RET[@]}"
	for w in ${RET[@]}; do
		if [[ ! "${w}" = "${cur}"* ]]; then
			continue
		fi
		if eval "[[ \"\${w}\" = *.$1 || -d \"\${w}\" ]]"; then
			qw="$(__simcli_quote "${w}")"
			if [ -d "${w}" ]; then
				COMPREPLY+=("${qw}/")
			else
				COMPREPLY+=("${qw}")
			fi
		fi
	done
}
__simcli_quote() {
    if [[ $1 == \'* || $1 == \"* ]]; then
        # Leave out first character
        printf %q "${1:1}"
    else
    	printf %q "$1"
    fi
}
autoload -U +X bashcompinit && bashcompinit
# use word boundary patterns for BSD or GNU sed
LWORD='[[:<:]]'
RWORD='[[:>:]]'
if sed --help 2>&1 | grep -q GNU; then
	LWORD='\<'
	RWORD='\>'
fi
__simcli_convert_bash_to_zsh() {
	sed \
	-e 's/declare -F/whence -w/' \
	-e 's/_get_comp_words_by_ref "\$@"/_get_comp_words_by_ref "\$*"/' \
	-e 's/local \([a-zA-Z0-9_]*\)=/local \1; \1=/' \
	-e 's/flags+=("\(--.*\)=")/flags+=("\1"); two_word_flags+=("\1")/' \
	-e 's/must_have_one_flag+=("\(--.*\)=")/must_have_one_flag+=("\1")/' \
	-e "s/${LWORD}_filedir${RWORD}/__simcli_filedir/g" \
	-e "s/${LWORD}_get_comp_words_by_ref${RWORD}/__simcli_get_comp_words_by_ref/g" \
	-e "s/${LWORD}__ltrim_colon_completions${RWORD}/__simcli_ltrim_colon_completions/g" \
	-e "s/${LWORD}compgen${RWORD}/__simcli_compgen/g" \
	-e "s/${LWORD}compopt${RWORD}/__simcli_compopt/g" \
	-e "s/${LWORD}declare${RWORD}/builtin declare/g" \
	-e "s/\\\$(type${RWORD}/\$(__simcli_type/g" \
	<<'BASH_COMPLETION_EOF'
`
	_, err = out.Write([]byte(init))
	if err != nil {
		ExitWithError(ExitError, err)
	}

	buf := new(bytes.Buffer)
	err = cmd.GenBashCompletion(buf)
	if err != nil {
		ExitWithError(ExitError, err)
	}
	_, err = out.Write(buf.Bytes())
	if err != nil {
		ExitWithError(ExitError, err)
	}

	tail := `
BASH_COMPLETION_EOF
}
__simcli_bash_source <(__simcli_convert_bash_to_zsh)
_complete simcli 2>/dev/null
`
	_, err = out.Write([]byte(tail))
	if err != nil {
		ExitWithError(ExitError, err)
	}
	return nil
}
