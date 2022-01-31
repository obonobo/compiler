#!/bin/bash

files=(
    'resources/src/*'
)

main() (
    outdir="lex-demo"
    cd_project_root

    echo
    print_header "Building ESAC..."
    build
    echo

    print_header "Clearing out demo files..."
    clear_demo_files "$outdir"
    echo

    print_header "Running CLI on files..."
    lex_files "$outdir"
)

clear_demo_files() {
    echo "rm -rf $1"
    rm -rf $1
}

lex_files() (
    outdir="${1:-lex-demo}"

    echo -e "./esacc lex --outdir \"$outdir\" $(expand_files)"
    echo
    ./esacc lex --outdir "$outdir" $(expand_files)
    echo "Done..."
    echo "Output Files:"

    i=1
    for file in ${outdir}/*; do
        echo -e "\t${file}"
        (($i % 2 == 0)) && echo
        ((i++))
    done
)

expand_files() {
    for file in "${files[@]}"; do
        echo -n "$(eval echo $file) "
    done
}

print_header() {
    echo -e "\033[1;34m$@\033[0m"
    print_hr
}

print_hr() {
    echo -e "\033[1;33m----------------------------------------------------------------------\033[0m"
}

cd_project_root() {
    [[ -f go.mod ]] && return
    cd "$(dirname $(find ../ -name go.mod -type f))"
}

build() {
    make
}

[[ ${BASH_SOURCE[0]} == $0 ]] && main "$@"
