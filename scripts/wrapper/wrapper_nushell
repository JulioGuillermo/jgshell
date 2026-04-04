# Markers:
#   START: \x1b]123;START\x07<user> <pwd> >>>
#   END:   \x1b]123;<exit_code>;<uuid>;DONE\x07

def --env _jg_emit_end_marker [] {
    let exit_code = ($env.LAST_EXIT_CODE | default 0)
    printf "\e]123;($exit_code);DONE\a>>A"
}
