function _jg_emit_end_marker --on-event fish_prompt
    printf '\033]JGSHELL;%s;DONE\007' "$status"
end

function fish_prompt
    _jg_emit_end_marker
end
