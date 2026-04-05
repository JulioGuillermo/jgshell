function _jg_emit_end_marker --on-event fish_prompt
    set -l exit_code $status
    printf '\033]JGSHELL;%s;DONE\007' $exit_code
end

function fish_prompt
    _jg_emit_end_marker
end
