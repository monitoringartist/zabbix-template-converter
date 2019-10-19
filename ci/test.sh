
walk_dir () {
    shopt -s nullglob dotglob

    for pathname in "$1"/*; do
        if [ -d "$pathname" ]; then
            walk_dir "$pathname"
        else
	    go run ../main.go -in "$pathname" -out testoutput.xml
	    dif=$(diff --text --ignore-blank-lines --ignore-space-change "$pathname" testoutput.xml | wc -l)
	    printf '%s:\n' "$pathname"
	    if [[ ! $dif -eq 0 ]]; then
		echo "go run ../main.go -in "$pathname" -out testoutput.xml"
		echo "diff --text --ignore-blank-lines --ignore-space-change \"$pathname\" testoutput.xml"
                echo $dif
	    fi
        fi
    done
}

walk_dir "4.4"
#walk_dir "4.0"
