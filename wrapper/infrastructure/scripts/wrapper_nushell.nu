def --env _jg_emit_end_marker [] {
    let exit_code = ($env.LAST_EXIT_CODE | default 0)
    printf "\e]JGSHELL;($exit_code);DONE\a>>A"
}
